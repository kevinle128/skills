package catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

type File struct {
	Path   string
	Data   []byte
	SHA256 string
}

type Catalog struct {
	Files  []File
	Skills []string
}

func Load(source fs.FS, root string) (Catalog, error) {
	if source == nil {
		return Catalog{}, fmt.Errorf("skill filesystem is nil")
	}
	if !fs.ValidPath(root) || root == "." {
		return Catalog{}, fmt.Errorf("invalid skill root %q", root)
	}

	sub, err := fs.Sub(source, root)
	if err != nil {
		return Catalog{}, fmt.Errorf("open skill root %q: %w", root, err)
	}

	var files []File
	skills := make(map[string]bool)
	err = fs.WalkDir(sub, ".", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if name == "." {
			return nil
		}
		if !fs.ValidPath(name) || path.Clean(name) != name {
			return fmt.Errorf("invalid embedded path %q", name)
		}
		parts := strings.Split(name, "/")
		if !strings.HasPrefix(parts[0], "kk-") || len(parts[0]) == len("kk-") {
			return fmt.Errorf("embedded path is outside a kk-* skill: %q", name)
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			return fmt.Errorf("embedded symlink is not allowed: %q", name)
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect embedded file %q: %w", name, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("embedded entry is not a regular file: %q", name)
		}
		data, err := fs.ReadFile(sub, name)
		if err != nil {
			return fmt.Errorf("read embedded file %q: %w", name, err)
		}
		sum := sha256.Sum256(data)
		files = append(files, File{Path: name, Data: data, SHA256: hex.EncodeToString(sum[:])})
		if len(parts) == 2 && parts[1] == "SKILL.md" {
			skills[parts[0]] = true
		}
		return nil
	})
	if err != nil {
		return Catalog{}, fmt.Errorf("walk embedded skills: %w", err)
	}
	if len(files) == 0 {
		return Catalog{}, fmt.Errorf("embedded skill catalog is empty")
	}

	top := make(map[string]bool)
	for _, file := range files {
		top[strings.SplitN(file.Path, "/", 2)[0]] = true
	}
	for skill := range top {
		if !skills[skill] {
			return Catalog{}, fmt.Errorf("skill %q does not contain SKILL.md", skill)
		}
	}

	skillNames := make([]string, 0, len(skills))
	for skill := range skills {
		skillNames = append(skillNames, skill)
	}
	sort.Strings(skillNames)
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return Catalog{Files: files, Skills: skillNames}, nil
}
