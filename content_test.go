package kevinkit

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/kevinle128/skills/internal/catalog"
)

func TestEmbeddedCatalogMatchesSourceTree(t *testing.T) {
	got, err := catalog.Load(SkillFS, SkillRoot)
	if err != nil {
		t.Fatal(err)
	}

	var want []string
	err = filepath.WalkDir("skills", func(name string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel("skills", name)
		if err != nil {
			return err
		}
		if !strings.HasPrefix(filepath.ToSlash(rel), "kk-") {
			return nil
		}
		want = append(want, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(want)

	gotPaths := make([]string, len(got.Files))
	for i, file := range got.Files {
		gotPaths[i] = file.Path
	}
	if !equalStrings(gotPaths, want) {
		t.Fatalf("embedded catalog differs from source\nembedded only: %v\nsource only: %v", difference(gotPaths, want), difference(want, gotPaths))
	}
	if len(got.Skills) != 34 {
		t.Fatalf("got %d embedded skills, want 34", len(got.Skills))
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func difference(a, b []string) []string {
	set := make(map[string]bool, len(b))
	for _, item := range b {
		set[item] = true
	}
	var result []string
	for _, item := range a {
		if !set[item] {
			result = append(result, item)
		}
	}
	return result
}

var _ fs.FS = SkillFS
