package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestCheckSelectsCurrentPlatformRelease(t *testing.T) {
	assetName := archiveName("2.0.0", runtime.GOOS, runtime.GOARCH)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, `{"tag_name":"v2.0.0","assets":[{"name":%q,"browser_download_url":"%s/archive","size":123},{"name":"checksums.txt","browser_download_url":"%s/checksums","size":80}]}`, assetName, serverURL(request), serverURL(request))
	}))
	defer server.Close()

	updater := newTestUpdater(t, Config{
		APIURL:         server.URL,
		CurrentVersion: "1.0.0",
		Executable:     filepath.Join(t.TempDir(), executableName(runtime.GOOS)),
		StateDir:       filepath.Join(t.TempDir(), ".kevinkit"),
	})
	release, available, err := updater.Check(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !available || release.Version != "2.0.0" || release.Archive.Name != assetName {
		t.Fatalf("unexpected release: available=%v release=%+v", available, release)
	}

	updater.config.CurrentVersion = "v2.0.0"
	_, available, err = updater.Check(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if available {
		t.Fatal("current release reported as an update")
	}

	updater.config.CurrentVersion = "3.0.0"
	_, available, err = updater.Check(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if available {
		t.Fatal("older release reported as an update")
	}
}

func TestCheckRejectsIncompleteAndUnstableReleases(t *testing.T) {
	assetName := archiveName("2.0.0", runtime.GOOS, runtime.GOARCH)
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "prerelease", body: `{"tag_name":"v2.0.0","prerelease":true}`, want: "not stable"},
		{name: "missing archive", body: `{"tag_name":"v2.0.0","assets":[{"name":"checksums.txt","browser_download_url":"https://example.test/checksums"}]}`, want: assetName},
		{name: "missing checksum", body: fmt.Sprintf(`{"tag_name":"v2.0.0","assets":[{"name":%q,"browser_download_url":"https://example.test/archive"}]}`, assetName), want: "checksums.txt"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				io.WriteString(writer, test.body)
			}))
			defer server.Close()
			updater := newTestUpdater(t, Config{
				APIURL:         server.URL,
				CurrentVersion: "1.0.0",
				Executable:     filepath.Join(t.TempDir(), executableName(runtime.GOOS)),
				StateDir:       filepath.Join(t.TempDir(), ".kevinkit"),
			})
			_, _, err := updater.Check(context.Background())
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("got error %v, want text %q", err, test.want)
			}
		})
	}
}

func TestApplyVerifiesReplacesAndExecutesNewBinary(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the process fixture is a POSIX shell script")
	}
	home := t.TempDir()
	executable := filepath.Join(home, "bin", "kk")
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	oldBinary := []byte("#!/bin/sh\nexit 99\n")
	if err := os.WriteFile(executable, oldBinary, 0o755); err != nil {
		t.Fatal(err)
	}
	newBinary := []byte("#!/bin/sh\nprintf '%s\\n%s\\n' \"$*\" \"$KEVINKIT_LOCK_TOKEN\" > \"$KK_TEST_OUTPUT\"\n")
	archive := testArchive(t, executableName(runtime.GOOS), newBinary)
	assetName := archiveName("2.0.0", runtime.GOOS, runtime.GOARCH)
	server := assetServer(t, assetName, archive, checksumLine(assetName, archive))
	defer server.Close()
	output := filepath.Join(home, "child-output")
	fixedTime := time.Date(2026, 9, 18, 10, 0, 0, 0, time.UTC)

	updater := newTestUpdater(t, Config{
		APIURL:         server.URL + "/latest",
		CurrentVersion: "1.0.0",
		Executable:     executable,
		StateDir:       filepath.Join(home, ".kevinkit"),
		Target:         "agents",
		LockToken:      "borrowed-token",
		Env:            append(os.Environ(), "KK_TEST_OUTPUT="+output),
		Now:            func() time.Time { return fixedTime },
	})
	release, available, err := updater.Check(context.Background())
	if err != nil || !available {
		t.Fatalf("check: available=%v err=%v", available, err)
	}
	result, err := updater.Apply(context.Background(), release)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Updated || result.Version != "2.0.0" {
		t.Fatalf("unexpected result: %+v", result)
	}
	assertFileBytes(t, executable, newBinary)
	assertFileBytes(t, result.Backup, oldBinary)
	assertFileBytes(t, output, []byte("install --target agents --yes\nborrowed-token\n"))
}

func TestApplyRejectsBadChecksumBeforeReplacement(t *testing.T) {
	home := t.TempDir()
	executable := filepath.Join(home, "kk")
	oldBinary := []byte("old executable")
	if err := os.WriteFile(executable, oldBinary, 0o755); err != nil {
		t.Fatal(err)
	}
	archive := testArchive(t, executableName(runtime.GOOS), []byte("new executable"))
	assetName := archiveName("2.0.0", runtime.GOOS, runtime.GOARCH)
	server := assetServer(t, assetName, archive, strings.Repeat("0", 64)+"  "+assetName+"\n")
	defer server.Close()
	updater := newTestUpdater(t, Config{
		APIURL:         server.URL + "/latest",
		CurrentVersion: "1.0.0",
		Executable:     executable,
		StateDir:       filepath.Join(home, ".kevinkit"),
		LockToken:      "borrowed-token",
	})
	release, _, err := updater.Check(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := updater.Apply(context.Background(), release); err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("got error %v, want checksum mismatch", err)
	}
	assertFileBytes(t, executable, oldBinary)
}

func TestExtractBinaryRejectsUnexpectedEntries(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fixture exercises tar.gz validation")
	}
	archive := testTarGZ(t, []tarEntry{
		{name: "../kk", data: []byte("bad")},
		{name: "kk", data: []byte("good")},
	})
	err := extractBinary("kk.tar.gz", archive, filepath.Join(t.TempDir(), "kk"))
	if err == nil || !strings.Contains(err.Error(), "unexpected entry") {
		t.Fatalf("got error %v, want unexpected entry", err)
	}
}

func TestChecksumForRejectsAmbiguousEntries(t *testing.T) {
	hash := strings.Repeat("a", 64)
	for _, test := range []struct {
		name string
		data string
		want string
	}{
		{name: "missing", data: hash + "  other\n", want: "missing checksum"},
		{name: "duplicate", data: hash + "  kk.tar.gz\n" + hash + "  kk.tar.gz\n", want: "duplicate checksum"},
		{name: "invalid", data: "bad  kk.tar.gz\n", want: "invalid checksum"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := checksumFor([]byte(test.data), "kk.tar.gz")
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("got error %v, want text %q", err, test.want)
			}
		})
	}
}

func TestNewRejectsInsecureReleaseURL(t *testing.T) {
	_, err := New(Config{
		APIURL:     "http://example.com/latest",
		Executable: filepath.Join(t.TempDir(), "kk"),
		StateDir:   filepath.Join(t.TempDir(), ".kevinkit"),
	})
	if err == nil || !strings.Contains(err.Error(), "HTTPS") {
		t.Fatalf("got error %v, want HTTPS requirement", err)
	}
}

func TestCheckHonorsHTTPTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		<-request.Context().Done()
	}))
	defer server.Close()
	updater := newTestUpdater(t, Config{
		APIURL:         server.URL,
		CurrentVersion: "1.0.0",
		Executable:     filepath.Join(t.TempDir(), executableName(runtime.GOOS)),
		StateDir:       filepath.Join(t.TempDir(), ".kevinkit"),
		HTTPClient:     &http.Client{Timeout: 10 * time.Millisecond},
	})
	if _, _, err := updater.Check(context.Background()); err == nil || !strings.Contains(err.Error(), "fetch latest release") {
		t.Fatalf("got error %v, want release fetch timeout", err)
	}
}

func TestDownloadRejectsOversizedBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Length", fmt.Sprint(maxArchive+1))
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	updater := newTestUpdater(t, Config{
		APIURL:     server.URL,
		Executable: filepath.Join(t.TempDir(), executableName(runtime.GOOS)),
		StateDir:   filepath.Join(t.TempDir(), ".kevinkit"),
	})
	_, err := updater.download(context.Background(), Asset{Name: "oversized", URL: server.URL}, maxArchive)
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("got error %v, want size rejection", err)
	}
}

func TestExecutableRecoveryPrimitives(t *testing.T) {
	directory := t.TempDir()
	current := filepath.Join(directory, executableName(runtime.GOOS))
	backup := filepath.Join(directory, "backup")
	assertWriteExecutable(t, current, []byte("current"))
	assertWriteExecutable(t, backup, []byte("backup"))

	if err := installExecutable(filepath.Join(directory, "missing"), current); err == nil {
		t.Fatal("missing staged binary did not fail")
	}
	assertFileBytes(t, current, []byte("current"))

	if err := restoreExecutable(backup, current); err != nil {
		t.Fatal(err)
	}
	assertFileBytes(t, current, []byte("backup"))
}

func newTestUpdater(t *testing.T, config Config) *Updater {
	t.Helper()
	updater, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	return updater
}

func serverURL(request *http.Request) string {
	return "http://" + request.Host
}

func assetServer(t *testing.T, assetName string, archive []byte, checksums string) *httptest.Server {
	t.Helper()
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/latest":
			fmt.Fprintf(writer, `{"tag_name":"v2.0.0","assets":[{"name":%q,"browser_download_url":%q,"size":%d},{"name":"checksums.txt","browser_download_url":%q,"size":%d}]}`, assetName, server.URL+"/archive", len(archive), server.URL+"/checksums", len(checksums))
		case "/archive":
			writer.Write(archive)
		case "/checksums":
			io.WriteString(writer, checksums)
		default:
			http.NotFound(writer, request)
		}
	}))
	return server
}

func checksumLine(name string, data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]) + "  " + name + "\n"
}

func testArchive(t *testing.T, name string, data []byte) []byte {
	t.Helper()
	if runtime.GOOS != "windows" {
		return testTarGZ(t, []tarEntry{{name: name, data: data}})
	}
	var result bytes.Buffer
	writer := zip.NewWriter(&result)
	header := &zip.FileHeader{Name: name, Method: zip.Deflate}
	header.SetMode(0o755)
	entry, err := writer.CreateHeader(header)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return result.Bytes()
}

type tarEntry struct {
	name string
	data []byte
}

func testTarGZ(t *testing.T, entries []tarEntry) []byte {
	t.Helper()
	var result bytes.Buffer
	gzipWriter := gzip.NewWriter(&result)
	tarWriter := tar.NewWriter(gzipWriter)
	for _, item := range entries {
		header := &tar.Header{Name: item.name, Mode: 0o755, Size: int64(len(item.data)), Typeflag: tar.TypeReg}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatal(err)
		}
		if _, err := tarWriter.Write(item.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return result.Bytes()
}

func assertFileBytes(t *testing.T, name string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s contains %q, want %q", name, got, want)
	}
}

func assertWriteExecutable(t *testing.T, name string, data []byte) {
	t.Helper()
	if err := os.WriteFile(name, data, 0o755); err != nil {
		t.Fatal(err)
	}
}
