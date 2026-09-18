package catalog

import (
	"testing"
	"testing/fstest"
)

func TestLoadIncludesDotfilesAndSortsFiles(t *testing.T) {
	source := fstest.MapFS{
		"skills/kk-plan/SKILL.md":    {Data: []byte("plan")},
		"skills/kk-plan/.example":    {Data: []byte("hidden")},
		"skills/kk-debug/SKILL.md":   {Data: []byte("debug")},
		"skills/kk-debug/nested.txt": {Data: []byte("nested")},
	}

	got, err := Load(source, "skills")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Files) != 4 {
		t.Fatalf("got %d files, want 4", len(got.Files))
	}
	if got.Files[0].Path != "kk-debug/SKILL.md" || got.Files[2].Path != "kk-plan/.example" {
		t.Fatalf("files are not sorted or dotfile is missing: %#v", got.Files)
	}
	if len(got.Skills) != 2 || got.Skills[0] != "kk-debug" || got.Skills[1] != "kk-plan" {
		t.Fatalf("unexpected skills: %#v", got.Skills)
	}
}

func TestLoadRejectsSkillWithoutSkillMarkdown(t *testing.T) {
	_, err := Load(fstest.MapFS{"skills/kk-plan/readme.txt": {Data: []byte("x")}}, "skills")
	if err == nil {
		t.Fatal("expected missing SKILL.md error")
	}
}

func TestLoadRejectsPathOutsideSkill(t *testing.T) {
	_, err := Load(fstest.MapFS{"skills/readme.txt": {Data: []byte("x")}}, "skills")
	if err == nil {
		t.Fatal("expected invalid top-level path error")
	}
}
