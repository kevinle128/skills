package lifecycle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kevinle128/skills/internal/catalog"
)

func TestLifecyclePreservesUserFilesAcrossUpdateAndUninstall(t *testing.T) {
	home := t.TempDir()
	managerA := testManager(t, home, "v1.0.0", map[string]string{
		"kk-fixture/SKILL.md":  "fixture-a",
		"kk-fixture/clean.txt": "clean-a",
		"kk-fixture/user.txt":  "user-a",
		"kk-fixture/stale.txt": "stale-a",
	})

	installed, err := managerA.Sync(Options{Target: "all"})
	if err != nil {
		t.Fatal(err)
	}
	if installed.Count(ActionCreate) != 8 {
		t.Fatalf("created %d files, want 8", installed.Count(ActionCreate))
	}
	assertFile(t, filepath.Join(home, ".agents/skills/kk-fixture/clean.txt"), "clean-a")
	assertFile(t, filepath.Join(home, ".claude/skills/kk-fixture/clean.txt"), "clean-a")

	idempotent, err := managerA.Sync(Options{Target: "all"})
	if err != nil {
		t.Fatal(err)
	}
	if idempotent.Count(ActionCreate)+idempotent.Count(ActionUpdate)+idempotent.Count(ActionRemove) != 0 {
		t.Fatalf("idempotent sync changed files: %#v", idempotent.Changes)
	}

	agentsRoot := filepath.Join(home, ".agents/skills/kk-fixture")
	if err := os.WriteFile(filepath.Join(agentsRoot, "user.txt"), []byte("user-edited"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(agentsRoot, "unknown.txt"), []byte("unknown"), 0o644); err != nil {
		t.Fatal(err)
	}

	managerB := testManager(t, home, "v2.0.0", map[string]string{
		"kk-fixture/SKILL.md":  "fixture-b",
		"kk-fixture/clean.txt": "clean-b",
		"kk-fixture/user.txt":  "user-b",
		"kk-fixture/new.txt":   "new-b",
	})
	updated, err := managerB.Sync(Options{Target: "all"})
	if err != nil {
		t.Fatal(err)
	}
	if updated.BackupDir == "" {
		t.Fatal("update did not create a backup")
	}
	if updated.Count(ActionPreserve) != 2 {
		t.Fatalf("preserved %d files, want the user edit and unknown file", updated.Count(ActionPreserve))
	}
	assertFile(t, filepath.Join(agentsRoot, "clean.txt"), "clean-b")
	assertFile(t, filepath.Join(agentsRoot, "user.txt"), "user-edited")
	assertFile(t, filepath.Join(agentsRoot, "unknown.txt"), "unknown")
	assertMissing(t, filepath.Join(agentsRoot, "stale.txt"))
	assertFile(t, filepath.Join(home, ".claude/skills/kk-fixture/user.txt"), "user-b")

	status, err := managerB.Status("all")
	if err != nil {
		t.Fatal(err)
	}
	if status.Clean() {
		t.Fatal("status should report the preserved user edit")
	}
	if !hasDrift(status, "agents", "kk-fixture/user.txt", "modified") {
		t.Fatalf("missing modified drift: %#v", status.Drift)
	}

	removed, err := managerB.Uninstall(Options{Target: "all"})
	if err != nil {
		t.Fatal(err)
	}
	if removed.Count(ActionPreserve) != 2 {
		t.Fatalf("preserved %d files, want the user edit and unknown file", removed.Count(ActionPreserve))
	}
	assertFile(t, filepath.Join(agentsRoot, "user.txt"), "user-edited")
	assertFile(t, filepath.Join(agentsRoot, "unknown.txt"), "unknown")
	assertMissing(t, filepath.Join(agentsRoot, "clean.txt"))
	assertMissing(t, filepath.Join(home, ".claude/skills/kk-fixture"))
	assertMissing(t, managerB.ManifestPath())
}

func TestLiveOldLockCannotBeStolen(t *testing.T) {
	home := t.TempDir()
	manager := testManager(t, home, "v1", map[string]string{"kk-fixture/SKILL.md": "x"})
	oldTime := time.Now().Add(-24 * time.Hour)
	manager.Now = func() time.Time { return oldTime }
	lock, err := manager.AcquireLock("")
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Release()
	manager.Now = time.Now
	if _, err := manager.AcquireLock(""); err == nil {
		t.Fatal("a live lock was stolen because it was old")
	}
}

func TestLockReleaseVerifiesOwnership(t *testing.T) {
	home := t.TempDir()
	manager := testManager(t, home, "v1", map[string]string{"kk-fixture/SKILL.md": "x"})
	lock, err := manager.AcquireLock("")
	if err != nil {
		t.Fatal(err)
	}
	ownerPath := filepath.Join(manager.StateDir, "lifecycle.lock", "owner.json")
	data, err := os.ReadFile(ownerPath)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), lock.Token(), "replacement-token", 1))
	if err := os.WriteFile(ownerPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := lock.Release(); err == nil {
		t.Fatal("lock release ignored changed ownership")
	}
	if _, err := os.Stat(ownerPath); err != nil {
		t.Fatalf("changed lock was removed: %v", err)
	}
}

func TestInterruptedSyncConvergesOnRetry(t *testing.T) {
	home := t.TempDir()
	managerA := testManager(t, home, "v1", map[string]string{
		"kk-fixture/SKILL.md": "a",
		"kk-fixture/one.txt":  "one-a",
		"kk-fixture/two.txt":  "two-a",
	})
	if _, err := managerA.Sync(Options{Target: "agents"}); err != nil {
		t.Fatal(err)
	}

	managerB := testManager(t, home, "v2", map[string]string{
		"kk-fixture/SKILL.md": "b",
		"kk-fixture/one.txt":  "one-b",
		"kk-fixture/two.txt":  "two-b",
	})
	count := 0
	managerB.BeforeChange = func(Change) error {
		count++
		if count == 2 {
			return errors.New("injected interruption")
		}
		return nil
	}
	if _, err := managerB.Sync(Options{Target: "agents"}); err == nil {
		t.Fatal("expected injected interruption")
	}

	managerB.BeforeChange = nil
	if _, err := managerB.Sync(Options{Target: "agents"}); err != nil {
		t.Fatal(err)
	}
	status, err := managerB.Status("agents")
	if err != nil {
		t.Fatal(err)
	}
	if !status.Clean() {
		t.Fatalf("retry did not converge: %#v", status.Drift)
	}
}

func TestSyncPreservesUnownedConflictUnlessForced(t *testing.T) {
	home := t.TempDir()
	manager := testManager(t, home, "v1", map[string]string{"kk-fixture/SKILL.md": "canonical"})
	name := filepath.Join(home, ".agents/skills/kk-fixture/SKILL.md")
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte("local"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := manager.Sync(Options{Target: "agents"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Count(ActionPreserve) != 1 {
		t.Fatalf("preserved %d files, want 1", result.Count(ActionPreserve))
	}
	assertFile(t, name, "local")

	result, err = manager.Sync(Options{Target: "agents", Force: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.BackupDir == "" {
		t.Fatal("forced replacement did not create backup")
	}
	assertFile(t, name, "canonical")
}

func TestSyncKeepsOwnershipOfStaleModifiedFile(t *testing.T) {
	home := t.TempDir()
	managerA := testManager(t, home, "v1", map[string]string{
		"kk-fixture/SKILL.md":  "fixture",
		"kk-fixture/stale.txt": "original",
	})
	if _, err := managerA.Sync(Options{Target: "agents"}); err != nil {
		t.Fatal(err)
	}
	stalePath := filepath.Join(home, ".agents", "skills", "kk-fixture", "stale.txt")
	if err := os.WriteFile(stalePath, []byte("user edit"), 0o644); err != nil {
		t.Fatal(err)
	}

	managerB := testManager(t, home, "v2", map[string]string{"kk-fixture/SKILL.md": "fixture"})
	result, err := managerB.Sync(Options{Target: "agents"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Count(ActionPreserve) != 1 {
		t.Fatalf("preserved %d files, want stale user edit", result.Count(ActionPreserve))
	}
	assertFile(t, stalePath, "user edit")
	status, err := managerB.Status("agents")
	if err != nil {
		t.Fatal(err)
	}
	if !hasDrift(status, "agents", "kk-fixture/stale.txt", "modified") {
		t.Fatalf("stale user edit lost ownership: %#v", status.Drift)
	}
	removed, err := managerB.Uninstall(Options{Target: "agents"})
	if err != nil {
		t.Fatal(err)
	}
	if removed.Count(ActionPreserve) != 1 {
		t.Fatalf("uninstall preserved %d files, want stale user edit", removed.Count(ActionPreserve))
	}
	assertFile(t, stalePath, "user edit")
}

func TestDryRunDoesNotWrite(t *testing.T) {
	home := t.TempDir()
	manager := testManager(t, home, "v1", map[string]string{"kk-fixture/SKILL.md": "x"})
	result, err := manager.Sync(Options{Target: "all", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.Count(ActionCreate) != 2 {
		t.Fatalf("dry-run planned %d creates, want 2", result.Count(ActionCreate))
	}
	assertMissing(t, filepath.Join(home, ".kevinkit"))
	assertMissing(t, filepath.Join(home, ".agents"))
}

func TestSyncRejectsManagedSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink setup requires privileges on Windows")
	}
	home := t.TempDir()
	manager := testManager(t, home, "v1", map[string]string{"kk-fixture/SKILL.md": "x"})
	target := filepath.Join(home, ".agents/skills/kk-fixture/SKILL.md")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(home, "outside")
	if err := os.WriteFile(outside, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, target); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Sync(Options{Target: "agents", Force: true}); err == nil {
		t.Fatal("expected symlink rejection")
	}
	assertFile(t, outside, "outside")
}

func testManager(t *testing.T, home, version string, files map[string]string) *Manager {
	t.Helper()
	payload := catalog.Catalog{}
	for name, value := range files {
		sum := sha256.Sum256([]byte(value))
		payload.Files = append(payload.Files, catalog.File{Path: name, Data: []byte(value), SHA256: hex.EncodeToString(sum[:])})
	}
	manager, err := NewManager(payload, version, home)
	if err != nil {
		t.Fatal(err)
	}
	manager.Now = func() time.Time { return time.Date(2026, 9, 18, 10, 0, 0, 0, time.UTC) }
	return manager
}

func assertFile(t *testing.T, name, want string) {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", name, data, want)
	}
}

func assertMissing(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Lstat(name); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s exists or returned unexpected error: %v", name, err)
	}
}

func hasDrift(status Status, target, name, kind string) bool {
	for _, drift := range status.Drift {
		if drift.Target == target && drift.Path == name && drift.Kind == kind {
			return true
		}
	}
	return false
}
