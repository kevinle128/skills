package lifecycle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func validateRelativePath(name string) error {
	if !fs.ValidPath(name) || name == "." || path.Clean(name) != name {
		return fmt.Errorf("invalid relative path %q", name)
	}
	if !strings.HasPrefix(strings.SplitN(name, "/", 2)[0], "kk-") {
		return fmt.Errorf("path is outside a kk-* skill: %q", name)
	}
	return nil
}

func safeJoin(root, name string) (string, error) {
	if err := validateRelativePath(name); err != nil {
		return "", err
	}
	joined := filepath.Join(root, filepath.FromSlash(name))
	rel, err := filepath.Rel(root, joined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes target root: %q", name)
	}
	return joined, nil
}

func rejectSymlinkComponents(base, name string) error {
	base = filepath.Clean(base)
	name = filepath.Clean(name)
	rel, err := filepath.Rel(base, name)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path %q escapes base %q", name, base)
	}
	current := base
	if info, err := os.Lstat(current); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("symlink is not allowed in managed path: %s", current)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		if part == "." || part == "" {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink is not allowed in managed path: %s", current)
		}
	}
	return nil
}

func fileHash(name string) (string, bool, error) {
	info, err := os.Lstat(name)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", true, fmt.Errorf("managed path is a symlink: %s", name)
	}
	if !info.Mode().IsRegular() {
		return "", true, fmt.Errorf("managed path is not a regular file: %s", name)
	}
	file, err := os.Open(name)
	if err != nil {
		return "", true, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", true, err
	}
	return hex.EncodeToString(hash.Sum(nil)), true, nil
}

func atomicWrite(name string, data []byte, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(name), ".kevinkit-*")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if err := temp.Chmod(mode); err != nil {
		temp.Close()
		return err
	}
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return replacePath(tempName, name)
}
