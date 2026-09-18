package e2e

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
)

func TestUserLifecycleAcrossRealBinaryUpdate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("self-update is not published for Windows until replacement is proven there")
	}
	repoRoot := findRepoRoot(t)
	work := t.TempDir()
	home := filepath.Join(work, "home")
	if err := os.MkdirAll(home, 0o700); err != nil {
		t.Fatal(err)
	}
	installedBinary := filepath.Join(work, "bin", "kk")
	releaseBinary := filepath.Join(work, "release", "kk")
	if err := os.MkdirAll(filepath.Dir(installedBinary), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(releaseBinary), 0o755); err != nil {
		t.Fatal(err)
	}

	assetName := fmt.Sprintf("kk_2.0.0_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	var archive []byte
	var corrupt atomic.Bool
	corrupt.Store(true)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/latest":
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"tag_name": "v2.0.0",
				"assets": []map[string]any{
					{"name": assetName, "browser_download_url": server.URL + "/archive", "size": len(archive)},
					{"name": "checksums.txt", "browser_download_url": server.URL + "/checksums", "size": 80},
				},
			})
		case "/archive":
			_, _ = writer.Write(archive)
		case "/checksums":
			hash := sha256.Sum256(archive)
			encoded := hex.EncodeToString(hash[:])
			if corrupt.Load() {
				encoded = strings.Repeat("0", 64)
			}
			fmt.Fprintf(writer, "%s  %s\n", encoded, assetName)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	buildBinary(t, repoRoot, "e2ea", "1.0.0", server.URL+"/latest", installedBinary)
	buildBinary(t, repoRoot, "e2eb", "2.0.0", server.URL+"/latest", releaseBinary)
	archive = tarGZBinary(t, releaseBinary)
	environment := isolatedEnvironment(home)

	output, code := runKK(t, installedBinary, environment, "version")
	assertExitAndText(t, output, code, 0, "kk 1.0.0")
	output, code = runKK(t, installedBinary, environment, "install", "--yes")
	assertExitAndText(t, output, code, 0, "Applied changes: create=8")
	output, code = runKK(t, installedBinary, environment, "status")
	assertExitAndText(t, output, code, 0, "status: clean")

	agentsSkill := filepath.Join(home, ".agents", "skills", "kk-fixture")
	claudeSkill := filepath.Join(home, ".claude", "skills", "kk-fixture")
	assertFile(t, filepath.Join(agentsSkill, "clean.txt"), "clean-a\n")
	assertFile(t, filepath.Join(claudeSkill, "clean.txt"), "clean-a\n")
	if err := os.WriteFile(filepath.Join(agentsSkill, "user.txt"), []byte("user-edited"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(agentsSkill, "unknown.txt"), []byte("unknown"), 0o644); err != nil {
		t.Fatal(err)
	}
	binaryBeforePreview := readFile(t, installedBinary)
	manifestPath := filepath.Join(home, ".kevinkit", "install-manifest.json")
	manifestBeforePreview := readFile(t, manifestPath)

	output, code = runKK(t, installedBinary, environment, "update", "--check")
	assertExitAndText(t, output, code, 0, "Update available: 1.0.0 -> 2.0.0")
	assertFileBytes(t, installedBinary, binaryBeforePreview)
	assertFileBytes(t, manifestPath, manifestBeforePreview)
	assertFile(t, filepath.Join(agentsSkill, "user.txt"), "user-edited")
	assertFile(t, filepath.Join(agentsSkill, "unknown.txt"), "unknown")

	output, code = runKK(t, installedBinary, environment, "update", "--dry-run")
	assertExitAndText(t, output, code, 0, "Would update: 1.0.0 -> 2.0.0")
	assertFileBytes(t, installedBinary, binaryBeforePreview)
	assertFileBytes(t, manifestPath, manifestBeforePreview)
	assertFile(t, filepath.Join(agentsSkill, "user.txt"), "user-edited")
	assertFile(t, filepath.Join(agentsSkill, "unknown.txt"), "unknown")

	output, code = runKK(t, installedBinary, environment, "update", "--yes")
	if code != 1 || !strings.Contains(output, "checksum mismatch") {
		t.Fatalf("corrupt update exit=%d output=%s", code, output)
	}
	if strings.Contains(output, "Activated kk") {
		t.Fatalf("corrupt release was activated: %s", output)
	}
	output, code = runKK(t, installedBinary, environment, "version")
	assertExitAndText(t, output, code, 0, "kk 1.0.0")
	assertFile(t, filepath.Join(agentsSkill, "clean.txt"), "clean-a\n")

	corrupt.Store(false)
	output, code = runKK(t, installedBinary, environment, "update", "--yes")
	if code != 0 {
		t.Fatalf("update exit=%d output=%s", code, output)
	}
	assertOutputOrder(t, output,
		"Verified checksum for "+assetName,
		"Activated kk 2.0.0",
		"Applied changes:",
		"KevinKit updated to 2.0.0",
	)

	output, code = runKK(t, installedBinary, environment, "version")
	assertExitAndText(t, output, code, 0, "kk 2.0.0")
	assertFile(t, filepath.Join(agentsSkill, "clean.txt"), "clean-b\n")
	assertFile(t, filepath.Join(agentsSkill, "user.txt"), "user-edited")
	assertFile(t, filepath.Join(agentsSkill, "unknown.txt"), "unknown")
	assertFile(t, filepath.Join(agentsSkill, "new.txt"), "new-b\n")
	assertMissing(t, filepath.Join(agentsSkill, "stale.txt"))
	assertFile(t, filepath.Join(claudeSkill, "clean.txt"), "clean-b\n")
	assertFile(t, filepath.Join(claudeSkill, "user.txt"), "user-b\n")
	assertMissing(t, filepath.Join(claudeSkill, "stale.txt"))

	output, code = runKK(t, installedBinary, environment, "status")
	assertExitAndText(t, output, code, 1, "agents:kk-fixture/user.txt: modified")
	output, code = runKK(t, installedBinary, environment, "uninstall", "--yes")
	assertExitAndText(t, output, code, 0, "preserve=2")
	assertFile(t, filepath.Join(agentsSkill, "user.txt"), "user-edited")
	assertFile(t, filepath.Join(agentsSkill, "unknown.txt"), "unknown")
	assertMissing(t, filepath.Join(agentsSkill, "clean.txt"))
	assertMissing(t, claudeSkill)
	assertMissing(t, filepath.Join(home, ".kevinkit", "install-manifest.json"))
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate E2E test source")
	}
	root, err := filepath.Abs(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func buildBinary(t *testing.T, repoRoot, tag, version, releaseURL, destination string) {
	t.Helper()
	ldflags := strings.Join([]string{
		"-X", "github.com/kevinle128/skills/internal/buildinfo.Version=" + version,
		"-X", "github.com/kevinle128/skills/internal/buildinfo.Commit=e2e",
		"-X", "github.com/kevinle128/skills/internal/buildinfo.BuildDate=2026-09-18",
		"-X", "github.com/kevinle128/skills/internal/buildinfo.ReleaseAPIURL=" + releaseURL,
	}, " ")
	command := exec.Command("go", "build", "-tags", tag, "-ldflags", ldflags, "-o", destination, "./cmd/kk")
	command.Dir = repoRoot
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build %s binary: %v\n%s", tag, err, output)
	}
}

func tarGZBinary(t *testing.T, binary string) []byte {
	t.Helper()
	data, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	var result bytes.Buffer
	gzipWriter := gzip.NewWriter(&result)
	tarWriter := tar.NewWriter(gzipWriter)
	header := &tar.Header{Name: "kk", Mode: 0o755, Size: int64(len(data)), Typeflag: tar.TypeReg}
	if err := tarWriter.WriteHeader(header); err != nil {
		t.Fatal(err)
	}
	if _, err := tarWriter.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return result.Bytes()
}

func isolatedEnvironment(home string) []string {
	result := make([]string, 0, len(os.Environ())+3)
	for _, item := range os.Environ() {
		if strings.HasPrefix(item, "HOME=") || strings.HasPrefix(item, "USERPROFILE=") || strings.HasPrefix(item, "KEVINKIT_LOCK_TOKEN=") {
			continue
		}
		result = append(result, item)
	}
	return append(result, "HOME="+home, "USERPROFILE="+home)
}

func runKK(t *testing.T, binary string, environment []string, args ...string) (string, int) {
	t.Helper()
	command := exec.Command(binary, args...)
	command.Env = environment
	output, err := command.CombinedOutput()
	if err == nil {
		return string(output), 0
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("run kk %v: %v\n%s", args, err, output)
	}
	return string(output), exitError.ExitCode()
}

func assertExitAndText(t *testing.T, output string, gotCode, wantCode int, text string) {
	t.Helper()
	if gotCode != wantCode || !strings.Contains(output, text) {
		t.Fatalf("exit=%d want=%d output=%q missing=%q", gotCode, wantCode, output, text)
	}
}

func assertOutputOrder(t *testing.T, output string, values ...string) {
	t.Helper()
	position := -1
	for _, value := range values {
		next := strings.Index(output, value)
		if next <= position {
			t.Fatalf("%q is out of order in output:\n%s", value, output)
		}
		position = next
	}
}

func assertFile(t *testing.T, name, want string) {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", name, data, want)
	}
}

func readFile(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func assertFileBytes(t *testing.T, name string, want []byte) {
	t.Helper()
	got := readFile(t, name)
	if !bytes.Equal(got, want) {
		t.Fatalf("%s changed unexpectedly", name)
	}
}

func assertMissing(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Lstat(name); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s exists or returned %v", name, err)
	}
}
