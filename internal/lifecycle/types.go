package lifecycle

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/kevinle128/skills/internal/catalog"
)

const manifestSchemaVersion = 1

type Target struct {
	ID   string
	Root string
}

type Action string

const (
	ActionCreate    Action = "create"
	ActionUpdate    Action = "update"
	ActionRemove    Action = "remove"
	ActionPreserve  Action = "preserve"
	ActionUnchanged Action = "unchanged"
)

type Change struct {
	Target string
	Path   string
	Action Action
	Reason string
	Data   []byte
}

type Result struct {
	Changes      []Change
	BackupDir    string
	ManifestPath string
}

func (r Result) Count(action Action) int {
	count := 0
	for _, change := range r.Changes {
		if change.Action == action {
			count++
		}
	}
	return count
}

type Drift struct {
	Target string
	Path   string
	Kind   string
}

type Status struct {
	Version string
	Drift   []Drift
	Targets []TargetStatus
}

type TargetStatus struct {
	ID        string
	Root      string
	Installed bool
	Owned     int
}

func (s Status) Clean() bool { return len(s.Drift) == 0 }

type Options struct {
	Target    string
	Force     bool
	DryRun    bool
	LockToken string
}

type Manager struct {
	Catalog      catalog.Catalog
	Version      string
	HomeDir      string
	StateDir     string
	Now          func() time.Time
	BeforeChange func(Change) error
}

func NewManager(payload catalog.Catalog, version, home string) (*Manager, error) {
	if !filepath.IsAbs(home) {
		return nil, fmt.Errorf("home directory must be absolute: %q", home)
	}
	return &Manager{
		Catalog:  payload,
		Version:  version,
		HomeDir:  filepath.Clean(home),
		StateDir: filepath.Join(filepath.Clean(home), ".kevinkit"),
		Now:      time.Now,
	}, nil
}

func (m *Manager) targets(selector string) ([]Target, error) {
	all := []Target{
		{ID: "agents", Root: filepath.Join(m.HomeDir, ".agents", "skills")},
		{ID: "claude-code", Root: filepath.Join(m.HomeDir, ".claude", "skills")},
	}
	switch selector {
	case "", "all":
		return all, nil
	case "agents", "claude-code":
		for _, target := range all {
			if target.ID == selector {
				return []Target{target}, nil
			}
		}
	}
	return nil, fmt.Errorf("unknown target %q; want all, agents, or claude-code", selector)
}

func sortedManifestFiles(files map[string]string) []manifestFile {
	result := make([]manifestFile, 0, len(files))
	for name, hash := range files {
		result = append(result, manifestFile{Path: name, SHA256: hash})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Path < result[j].Path })
	return result
}
