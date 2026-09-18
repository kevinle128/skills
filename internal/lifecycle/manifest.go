package lifecycle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type manifest struct {
	SchemaVersion int              `json:"schema_version"`
	KitVersion    string           `json:"kit_version"`
	Targets       []manifestTarget `json:"targets"`
}

type manifestTarget struct {
	ID    string         `json:"id"`
	Root  string         `json:"root"`
	Files []manifestFile `json:"files"`
}

type manifestFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

func emptyManifest(version string) manifest {
	return manifest{SchemaVersion: manifestSchemaVersion, KitVersion: version}
}

func loadManifest(name, version string) (manifest, error) {
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return emptyManifest(version), nil
	}
	if err != nil {
		return manifest{}, fmt.Errorf("open manifest: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(io.LimitReader(file, 8<<20))
	decoder.DisallowUnknownFields()
	var result manifest
	if err := decoder.Decode(&result); err != nil {
		return manifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	if result.SchemaVersion != manifestSchemaVersion {
		return manifest{}, fmt.Errorf("unsupported manifest schema %d", result.SchemaVersion)
	}
	if err := validateManifest(result); err != nil {
		return manifest{}, err
	}
	return result, nil
}

func validateManifest(value manifest) error {
	targets := make(map[string]bool)
	for _, target := range value.Targets {
		if target.ID != "agents" && target.ID != "claude-code" {
			return fmt.Errorf("manifest contains unknown target %q", target.ID)
		}
		if targets[target.ID] {
			return fmt.Errorf("manifest contains duplicate target %q", target.ID)
		}
		targets[target.ID] = true
		paths := make(map[string]bool)
		for _, file := range target.Files {
			if err := validateRelativePath(file.Path); err != nil {
				return fmt.Errorf("manifest target %s: %w", target.ID, err)
			}
			if paths[file.Path] {
				return fmt.Errorf("manifest target %s contains duplicate path %q", target.ID, file.Path)
			}
			paths[file.Path] = true
			if len(file.SHA256) != 64 {
				return fmt.Errorf("manifest target %s has invalid hash for %q", target.ID, file.Path)
			}
		}
	}
	return nil
}

func writeManifest(name string, value manifest) error {
	sort.Slice(value.Targets, func(i, j int) bool { return value.Targets[i].ID < value.Targets[j].ID })
	for i := range value.Targets {
		sort.Slice(value.Targets[i].Files, func(a, b int) bool {
			return value.Targets[i].Files[a].Path < value.Targets[i].Files[b].Path
		})
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(name), 0o700); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}
	return atomicWrite(name, data, 0o600)
}

func removeManifest(name string) error {
	err := os.Remove(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("remove manifest: %w", err)
	}
	return nil
}

func manifestTargetMap(value manifest) map[string]manifestTarget {
	result := make(map[string]manifestTarget, len(value.Targets))
	for _, target := range value.Targets {
		result[target.ID] = target
	}
	return result
}

func manifestFileMap(value manifestTarget) map[string]string {
	result := make(map[string]string, len(value.Files))
	for _, file := range value.Files {
		result[file.Path] = file.SHA256
	}
	return result
}
