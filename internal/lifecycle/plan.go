package lifecycle

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kevinle128/skills/internal/catalog"
)

type changePlan struct {
	changes []Change
	next    manifest
	targets map[string]Target
}

func (m *Manager) buildSyncPlan(current manifest, targets []Target, force bool) (changePlan, error) {
	currentTargets := manifestTargetMap(current)
	nextTargets := manifestTargetMap(current)
	selected := make(map[string]Target, len(targets))
	desired := make(map[string]catalog.File, len(m.Catalog.Files))
	for _, file := range m.Catalog.Files {
		desired[file.Path] = file
	}

	var changes []Change
	for _, target := range targets {
		selected[target.ID] = target
		previous := currentTargets[target.ID]
		if previous.Root != "" && filepath.Clean(previous.Root) != filepath.Clean(target.Root) {
			return changePlan{}, fmt.Errorf("manifest target %s root changed from %q to %q", target.ID, previous.Root, target.Root)
		}
		oldFiles := manifestFileMap(previous)
		nextFiles := make(map[string]string, len(desired))
		for name, hash := range oldFiles {
			nextFiles[name] = hash
		}

		for _, file := range m.Catalog.Files {
			fullName, err := m.managedPath(target, file.Path)
			if err != nil {
				return changePlan{}, err
			}
			currentHash, exists, err := fileHash(fullName)
			if err != nil {
				return changePlan{}, fmt.Errorf("inspect %s:%s: %w", target.ID, file.Path, err)
			}
			oldHash, owned := oldFiles[file.Path]
			switch {
			case !exists:
				changes = append(changes, Change{Target: target.ID, Path: file.Path, Action: ActionCreate, Reason: "missing", Data: file.Data})
				nextFiles[file.Path] = file.SHA256
			case currentHash == file.SHA256:
				changes = append(changes, Change{Target: target.ID, Path: file.Path, Action: ActionUnchanged, Reason: "already current"})
				nextFiles[file.Path] = file.SHA256
			case owned && currentHash == oldHash:
				changes = append(changes, Change{Target: target.ID, Path: file.Path, Action: ActionUpdate, Reason: "clean owned file", Data: file.Data})
				nextFiles[file.Path] = file.SHA256
			case force:
				reason := "forced unowned conflict"
				if owned {
					reason = "forced modified file"
				}
				changes = append(changes, Change{Target: target.ID, Path: file.Path, Action: ActionUpdate, Reason: reason, Data: file.Data})
				nextFiles[file.Path] = file.SHA256
			default:
				reason := "unowned conflict"
				if owned {
					reason = "user-modified file"
				}
				changes = append(changes, Change{Target: target.ID, Path: file.Path, Action: ActionPreserve, Reason: reason})
				if !owned {
					delete(nextFiles, file.Path)
				}
			}
		}

		for name, oldHash := range oldFiles {
			if _, ok := desired[name]; ok {
				continue
			}
			fullName, err := m.managedPath(target, name)
			if err != nil {
				return changePlan{}, err
			}
			currentHash, exists, err := fileHash(fullName)
			if err != nil {
				return changePlan{}, fmt.Errorf("inspect stale %s:%s: %w", target.ID, name, err)
			}
			delete(nextFiles, name)
			switch {
			case !exists:
				changes = append(changes, Change{Target: target.ID, Path: name, Action: ActionUnchanged, Reason: "stale file already absent"})
			case currentHash == oldHash:
				changes = append(changes, Change{Target: target.ID, Path: name, Action: ActionRemove, Reason: "stale clean file"})
			case force:
				changes = append(changes, Change{Target: target.ID, Path: name, Action: ActionRemove, Reason: "forced stale modified file"})
			default:
				changes = append(changes, Change{Target: target.ID, Path: name, Action: ActionPreserve, Reason: "stale user-modified file"})
				nextFiles[name] = oldHash
			}
		}
		unknown, err := m.unknownFiles(target, desired, oldFiles)
		if err != nil {
			return changePlan{}, err
		}
		changes = append(changes, unknown...)

		if len(nextFiles) == 0 {
			delete(nextTargets, target.ID)
		} else {
			nextTargets[target.ID] = manifestTarget{ID: target.ID, Root: target.Root, Files: sortedManifestFiles(nextFiles)}
		}
	}

	next := emptyManifest(m.Version)
	for _, target := range nextTargets {
		next.Targets = append(next.Targets, target)
	}
	sortChanges(changes)
	return changePlan{changes: changes, next: next, targets: selected}, nil
}

func (m *Manager) buildUninstallPlan(current manifest, targets []Target) (changePlan, error) {
	currentTargets := manifestTargetMap(current)
	nextTargets := manifestTargetMap(current)
	selected := make(map[string]Target, len(targets))
	var changes []Change

	for _, target := range targets {
		selected[target.ID] = target
		previous, ok := currentTargets[target.ID]
		if !ok {
			continue
		}
		if filepath.Clean(previous.Root) != filepath.Clean(target.Root) {
			return changePlan{}, fmt.Errorf("manifest target %s root changed from %q to %q", target.ID, previous.Root, target.Root)
		}
		for _, file := range previous.Files {
			fullName, err := m.managedPath(target, file.Path)
			if err != nil {
				return changePlan{}, err
			}
			currentHash, exists, err := fileHash(fullName)
			if err != nil {
				return changePlan{}, fmt.Errorf("inspect %s:%s: %w", target.ID, file.Path, err)
			}
			switch {
			case !exists:
				changes = append(changes, Change{Target: target.ID, Path: file.Path, Action: ActionUnchanged, Reason: "already absent"})
			case currentHash == file.SHA256:
				changes = append(changes, Change{Target: target.ID, Path: file.Path, Action: ActionRemove, Reason: "clean owned file"})
			default:
				changes = append(changes, Change{Target: target.ID, Path: file.Path, Action: ActionPreserve, Reason: "user-modified file"})
			}
		}
		known := manifestFileMap(previous)
		unknown, err := m.unknownFiles(target, nil, known)
		if err != nil {
			return changePlan{}, err
		}
		changes = append(changes, unknown...)
		delete(nextTargets, target.ID)
	}

	next := emptyManifest(current.KitVersion)
	for _, target := range nextTargets {
		next.Targets = append(next.Targets, target)
	}
	sortChanges(changes)
	return changePlan{changes: changes, next: next, targets: selected}, nil
}

func (m *Manager) unknownFiles(target Target, desired map[string]catalog.File, owned map[string]string) ([]Change, error) {
	known := make(map[string]bool, len(desired)+len(owned))
	skills := make(map[string]bool)
	for name := range desired {
		known[name] = true
		skills[strings.SplitN(name, "/", 2)[0]] = true
	}
	for name := range owned {
		known[name] = true
		skills[strings.SplitN(name, "/", 2)[0]] = true
	}
	var changes []Change
	for skill := range skills {
		root := filepath.Join(target.Root, filepath.FromSlash(skill))
		if err := rejectSymlinkComponents(m.HomeDir, root); err != nil {
			return nil, err
		}
		err := filepath.WalkDir(root, func(name string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				if os.IsNotExist(walkErr) && name == root {
					return nil
				}
				return walkErr
			}
			if name == root || entry.IsDir() {
				return nil
			}
			relative, err := filepath.Rel(target.Root, name)
			if err != nil {
				return err
			}
			relative = filepath.ToSlash(relative)
			if !fs.ValidPath(relative) || path.Clean(relative) != relative {
				return fmt.Errorf("invalid existing path %q", relative)
			}
			if !known[relative] {
				changes = append(changes, Change{Target: target.ID, Path: relative, Action: ActionPreserve, Reason: "unknown file"})
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("inspect unknown files in %s:%s: %w", target.ID, skill, err)
		}
	}
	return changes, nil
}

func (m *Manager) managedPath(target Target, name string) (string, error) {
	fullName, err := safeJoin(target.Root, name)
	if err != nil {
		return "", err
	}
	if err := rejectSymlinkComponents(m.HomeDir, fullName); err != nil {
		return "", err
	}
	return fullName, nil
}

func sortChanges(changes []Change) {
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Target != changes[j].Target {
			return changes[i].Target < changes[j].Target
		}
		if changes[i].Path != changes[j].Path {
			return changes[i].Path < changes[j].Path
		}
		return changes[i].Action < changes[j].Action
	})
}
