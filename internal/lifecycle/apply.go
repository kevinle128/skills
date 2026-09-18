package lifecycle

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (m *Manager) ManifestPath() string {
	return filepath.Join(m.StateDir, "install-manifest.json")
}

func (m *Manager) AcquireLock(token string) (*Lock, error) {
	return acquireLock(m.StateDir, token, m.now())
}

func (m *Manager) Sync(options Options) (Result, error) {
	targets, err := m.targets(options.Target)
	if err != nil {
		return Result{}, err
	}
	if options.DryRun {
		current, err := loadManifest(m.ManifestPath(), m.Version)
		if err != nil {
			return Result{}, err
		}
		plan, err := m.buildSyncPlan(current, targets, options.Force)
		return Result{Changes: plan.changes, ManifestPath: m.ManifestPath()}, err
	}

	lock, err := m.AcquireLock(options.LockToken)
	if err != nil {
		return Result{}, err
	}
	defer lock.Release()

	current, err := loadManifest(m.ManifestPath(), m.Version)
	if err != nil {
		return Result{}, err
	}
	plan, err := m.buildSyncPlan(current, targets, options.Force)
	if err != nil {
		return Result{}, err
	}
	return m.apply(plan)
}

func (m *Manager) Uninstall(options Options) (Result, error) {
	targets, err := m.targets(options.Target)
	if err != nil {
		return Result{}, err
	}
	if options.DryRun {
		current, err := loadManifest(m.ManifestPath(), m.Version)
		if err != nil {
			return Result{}, err
		}
		plan, err := m.buildUninstallPlan(current, targets)
		return Result{Changes: plan.changes, ManifestPath: m.ManifestPath()}, err
	}

	lock, err := m.AcquireLock(options.LockToken)
	if err != nil {
		return Result{}, err
	}
	defer lock.Release()

	current, err := loadManifest(m.ManifestPath(), m.Version)
	if err != nil {
		return Result{}, err
	}
	plan, err := m.buildUninstallPlan(current, targets)
	if err != nil {
		return Result{}, err
	}
	return m.apply(plan)
}

func (m *Manager) apply(plan changePlan) (Result, error) {
	result := Result{Changes: plan.changes, ManifestPath: m.ManifestPath()}
	for _, change := range plan.changes {
		if change.Action != ActionCreate && change.Action != ActionUpdate && change.Action != ActionRemove {
			continue
		}
		if m.BeforeChange != nil {
			if err := m.BeforeChange(change); err != nil {
				return result, fmt.Errorf("before %s %s:%s: %w", change.Action, change.Target, change.Path, err)
			}
		}
		target := plan.targets[change.Target]
		fullName, err := m.managedPath(target, change.Path)
		if err != nil {
			return result, err
		}

		if change.Action == ActionUpdate || change.Action == ActionRemove {
			if result.BackupDir == "" {
				result.BackupDir = filepath.Join(m.StateDir, "backups", m.now().UTC().Format("20060102T150405.000000000Z"))
			}
			backupName, err := safeJoin(filepath.Join(result.BackupDir, change.Target), change.Path)
			if err != nil {
				return result, err
			}
			data, err := os.ReadFile(fullName)
			if err != nil {
				return result, fmt.Errorf("read backup source %s:%s: %w", change.Target, change.Path, err)
			}
			if err := atomicWrite(backupName, data, 0o600); err != nil {
				return result, fmt.Errorf("back up %s:%s: %w", change.Target, change.Path, err)
			}
		}

		switch change.Action {
		case ActionCreate, ActionUpdate:
			if err := atomicWrite(fullName, change.Data, 0o644); err != nil {
				return result, fmt.Errorf("write %s:%s: %w", change.Target, change.Path, err)
			}
		case ActionRemove:
			if err := os.Remove(fullName); err != nil && !errors.Is(err, os.ErrNotExist) {
				return result, fmt.Errorf("remove %s:%s: %w", change.Target, change.Path, err)
			}
			pruneEmptyParents(filepath.Dir(fullName), target.Root)
		}
	}

	if len(plan.next.Targets) == 0 {
		if err := removeManifest(m.ManifestPath()); err != nil {
			return result, err
		}
		return result, nil
	}
	if err := writeManifest(m.ManifestPath(), plan.next); err != nil {
		return result, err
	}
	return result, nil
}

func pruneEmptyParents(start, root string) {
	root = filepath.Clean(root)
	current := filepath.Clean(start)
	for current != root && strings.HasPrefix(current, root+string(filepath.Separator)) {
		if err := os.Remove(current); err != nil {
			return
		}
		current = filepath.Dir(current)
	}
}

func (m *Manager) now() time.Time {
	if m.Now == nil {
		return time.Now()
	}
	return m.Now()
}
