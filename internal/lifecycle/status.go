package lifecycle

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/kevinle128/skills/internal/catalog"
)

func (m *Manager) Status(selector string) (Status, error) {
	targets, err := m.targets(selector)
	if err != nil {
		return Status{}, err
	}
	current, err := loadManifest(m.ManifestPath(), m.Version)
	if err != nil {
		return Status{}, err
	}
	currentTargets := manifestTargetMap(current)
	desired := make(map[string]catalog.File, len(m.Catalog.Files))
	for _, file := range m.Catalog.Files {
		desired[file.Path] = file
	}

	result := Status{Version: current.KitVersion}
	for _, target := range targets {
		installed, ok := currentTargets[target.ID]
		if !ok {
			result.Targets = append(result.Targets, TargetStatus{ID: target.ID, Root: target.Root})
			result.Drift = append(result.Drift, Drift{Target: target.ID, Kind: "not-installed"})
			continue
		}
		if filepath.Clean(installed.Root) != filepath.Clean(target.Root) {
			return Status{}, fmt.Errorf("manifest target %s root changed from %q to %q", target.ID, installed.Root, target.Root)
		}
		owned := manifestFileMap(installed)
		result.Targets = append(result.Targets, TargetStatus{ID: target.ID, Root: target.Root, Installed: true, Owned: len(owned)})

		for name, installedHash := range owned {
			fullName, err := m.managedPath(target, name)
			if err != nil {
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "unsafe"})
				continue
			}
			currentHash, exists, err := fileHash(fullName)
			if err != nil {
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "unsafe"})
				continue
			}
			want, wanted := desired[name]
			switch {
			case !exists:
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "missing"})
			case currentHash != installedHash:
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "modified"})
			case !wanted:
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "stale"})
			case installedHash != want.SHA256:
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "outdated"})
			}
		}

		for name := range desired {
			if _, ok := owned[name]; ok {
				continue
			}
			fullName, err := m.managedPath(target, name)
			if err != nil {
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "unsafe"})
				continue
			}
			_, exists, err := fileHash(fullName)
			if err != nil {
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "unsafe"})
			} else if exists {
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "unmanaged"})
			} else {
				result.Drift = append(result.Drift, Drift{Target: target.ID, Path: name, Kind: "missing"})
			}
		}
	}
	sort.Slice(result.Drift, func(i, j int) bool {
		if result.Drift[i].Target != result.Drift[j].Target {
			return result.Drift[i].Target < result.Drift[j].Target
		}
		if result.Drift[i].Path != result.Drift[j].Path {
			return result.Drift[i].Path < result.Drift[j].Path
		}
		return result.Drift[i].Kind < result.Drift[j].Kind
	})
	return result, nil
}
