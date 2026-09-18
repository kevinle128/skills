package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	maxReleaseResponse = 2 << 20
	maxChecksumFile    = 1 << 20
	maxArchive         = 64 << 20
)

type Asset struct {
	Name string
	URL  string
	Size int64
}

type Release struct {
	Version  string
	Archive  Asset
	Checksum Asset
}

type Config struct {
	APIURL         string
	CurrentVersion string
	Executable     string
	StateDir       string
	Target         string
	Force          bool
	LockToken      string
	HTTPClient     *http.Client
	Stdout         io.Writer
	Stderr         io.Writer
	Env            []string
	Now            func() time.Time
}

type Updater struct {
	config Config
}

type Result struct {
	Version string
	Updated bool
	Backup  string
	Archive string
}

func New(config Config) (*Updater, error) {
	if config.APIURL == "" {
		return nil, fmt.Errorf("release API URL is required")
	}
	if err := validateSecureURL(config.APIURL); err != nil {
		return nil, fmt.Errorf("release API URL: %w", err)
	}
	if config.Executable == "" {
		return nil, fmt.Errorf("current executable path is required")
	}
	if config.StateDir == "" {
		return nil, fmt.Errorf("state directory is required")
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	if config.Stdout == nil {
		config.Stdout = io.Discard
	}
	if config.Stderr == nil {
		config.Stderr = io.Discard
	}
	if config.Now == nil {
		config.Now = time.Now
	}
	return &Updater{config: config}, nil
}

func (u *Updater) Check(ctx context.Context) (Release, bool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.config.APIURL, nil)
	if err != nil {
		return Release{}, false, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "KevinKit/"+u.config.CurrentVersion)
	response, err := u.config.HTTPClient.Do(request)
	if err != nil {
		return Release{}, false, fmt.Errorf("fetch latest release: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Release{}, false, fmt.Errorf("fetch latest release: HTTP %d", response.StatusCode)
	}

	var body struct {
		TagName    string `json:"tag_name"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
		Assets     []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxReleaseResponse))
	if err := decoder.Decode(&body); err != nil {
		return Release{}, false, fmt.Errorf("decode latest release: %w", err)
	}
	if body.Draft || body.Prerelease {
		return Release{}, false, fmt.Errorf("latest release %q is not stable", body.TagName)
	}
	version := strings.TrimPrefix(body.TagName, "v")
	releaseVersion, err := parseVersion(version)
	if err != nil {
		return Release{}, false, fmt.Errorf("latest release version: %w", err)
	}
	archiveName := archiveName(version, runtime.GOOS, runtime.GOARCH)
	var release Release
	release.Version = version
	for _, asset := range body.Assets {
		item := Asset{Name: asset.Name, URL: asset.BrowserDownloadURL, Size: asset.Size}
		switch asset.Name {
		case archiveName:
			if release.Archive.Name != "" {
				return Release{}, false, fmt.Errorf("release contains duplicate asset %q", archiveName)
			}
			release.Archive = item
		case "checksums.txt":
			if release.Checksum.Name != "" {
				return Release{}, false, fmt.Errorf("release contains duplicate checksums.txt")
			}
			release.Checksum = item
		}
	}
	if release.Archive.Name == "" {
		return Release{}, false, fmt.Errorf("release does not contain %s", archiveName)
	}
	if release.Checksum.Name == "" {
		return Release{}, false, fmt.Errorf("release does not contain checksums.txt")
	}
	if release.Archive.Size < 0 || release.Archive.Size > maxArchive {
		return Release{}, false, fmt.Errorf("release archive size %d is outside the allowed range", release.Archive.Size)
	}
	current := strings.TrimPrefix(u.config.CurrentVersion, "v")
	if current == "dev" {
		return release, true, nil
	}
	currentVersion, err := parseVersion(current)
	if err != nil {
		return Release{}, false, fmt.Errorf("current version: %w", err)
	}
	return release, compareVersion(releaseVersion, currentVersion) > 0, nil
}

func (u *Updater) Apply(ctx context.Context, release Release) (Result, error) {
	if u.config.LockToken == "" {
		return Result{}, fmt.Errorf("lifecycle lock token is required")
	}
	if release.Archive.Name == "" || release.Checksum.Name == "" {
		return Result{}, fmt.Errorf("release assets are incomplete")
	}
	tempDir, err := os.MkdirTemp("", "kevinkit-update-*")
	if err != nil {
		return Result{}, err
	}
	defer os.RemoveAll(tempDir)

	checksums, err := u.download(ctx, release.Checksum, maxChecksumFile)
	if err != nil {
		return Result{}, err
	}
	archive, err := u.download(ctx, release.Archive, maxArchive)
	if err != nil {
		return Result{}, err
	}
	wantHash, err := checksumFor(checksums, release.Archive.Name)
	if err != nil {
		return Result{}, err
	}
	gotHash := sha256.Sum256(archive)
	if hex.EncodeToString(gotHash[:]) != wantHash {
		return Result{}, fmt.Errorf("checksum mismatch for %s", release.Archive.Name)
	}
	fmt.Fprintf(u.config.Stdout, "Verified checksum for %s.\n", release.Archive.Name)

	newBinary := filepath.Join(tempDir, executableName(runtime.GOOS))
	if err := extractBinary(release.Archive.Name, archive, newBinary); err != nil {
		return Result{}, err
	}
	if err := verifyExecutableDirectory(u.config.Executable); err != nil {
		return Result{}, err
	}
	backup, err := u.backupExecutable()
	if err != nil {
		return Result{}, err
	}
	if err := installExecutable(newBinary, u.config.Executable); err != nil {
		if restoreErr := restoreExecutable(backup, u.config.Executable); restoreErr != nil {
			return Result{}, fmt.Errorf("replace executable: %v; restore backup: %w", err, restoreErr)
		}
		return Result{}, fmt.Errorf("replace executable: %w", err)
	}
	fmt.Fprintf(u.config.Stdout, "Activated kk %s.\n", release.Version)

	args := []string{"install", "--target", defaultTarget(u.config.Target), "--yes"}
	if u.config.Force {
		args = append(args, "--force")
	}
	command := exec.CommandContext(ctx, u.config.Executable, args...)
	command.Stdout = u.config.Stdout
	command.Stderr = u.config.Stderr
	command.Env = append(append([]string{}, u.config.Env...), "KEVINKIT_LOCK_TOKEN="+u.config.LockToken)
	if err := command.Run(); err != nil {
		return Result{Version: release.Version, Updated: true, Backup: backup, Archive: release.Archive.Name}, fmt.Errorf("new kk binary could not synchronize skills: %w", err)
	}
	return Result{Version: release.Version, Updated: true, Backup: backup, Archive: release.Archive.Name}, nil
}

func (u *Updater) download(ctx context.Context, asset Asset, limit int64) ([]byte, error) {
	if asset.URL == "" {
		return nil, fmt.Errorf("asset %q has no download URL", asset.Name)
	}
	if err := validateSecureURL(asset.URL); err != nil {
		return nil, fmt.Errorf("asset %q URL: %w", asset.Name, err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.URL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "KevinKit/"+u.config.CurrentVersion)
	response, err := u.config.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", asset.Name, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: HTTP %d", asset.Name, response.StatusCode)
	}
	if response.ContentLength > limit {
		return nil, fmt.Errorf("download %s exceeds %d bytes", asset.Name, limit)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", asset.Name, err)
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("download %s exceeds %d bytes", asset.Name, limit)
	}
	return data, nil
}

type semanticVersion [3]uint64

func parseVersion(value string) (semanticVersion, error) {
	var result semanticVersion
	if strings.ContainsAny(value, "+-") {
		return result, fmt.Errorf("%q is not a stable semantic version", value)
	}
	parts := strings.Split(value, ".")
	if len(parts) != len(result) {
		return result, fmt.Errorf("%q is not a semantic version", value)
	}
	for index, part := range parts {
		if part == "" || (len(part) > 1 && part[0] == '0') {
			return result, fmt.Errorf("%q is not a semantic version", value)
		}
		number, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return result, fmt.Errorf("%q is not a semantic version", value)
		}
		result[index] = number
	}
	return result, nil
}

func compareVersion(left, right semanticVersion) int {
	for index := range left {
		if left[index] < right[index] {
			return -1
		}
		if left[index] > right[index] {
			return 1
		}
	}
	return 0
}

func validateSecureURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Hostname() == "" {
		return fmt.Errorf("invalid URL")
	}
	if parsed.Scheme == "https" {
		return nil
	}
	host := parsed.Hostname()
	if parsed.Scheme == "http" && (host == "localhost" || host == "127.0.0.1" || host == "::1") {
		return nil
	}
	return fmt.Errorf("URL must use HTTPS")
}

func (u *Updater) backupExecutable() (string, error) {
	backup := filepath.Join(u.config.StateDir, "backups", u.config.Now().UTC().Format("20060102T150405.000000000Z"), "binary", filepath.Base(u.config.Executable))
	if err := copyExecutable(u.config.Executable, backup); err != nil {
		return "", fmt.Errorf("back up current executable: %w", err)
	}
	return backup, nil
}

func archiveName(version, goos, goarch string) string {
	extension := ".tar.gz"
	if goos == "windows" {
		extension = ".zip"
	}
	return fmt.Sprintf("kk_%s_%s_%s%s", version, goos, goarch, extension)
}

func executableName(goos string) string {
	if goos == "windows" {
		return "kk.exe"
	}
	return "kk"
}

func defaultTarget(target string) string {
	if target == "" {
		return "all"
	}
	return target
}

func checksumFor(data []byte, wanted string) (string, error) {
	found := ""
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || fields[1] != wanted {
			continue
		}
		candidate := strings.ToLower(fields[0])
		decoded, err := hex.DecodeString(candidate)
		if err != nil || len(decoded) != sha256.Size {
			return "", fmt.Errorf("invalid checksum for %s", wanted)
		}
		if found != "" {
			return "", fmt.Errorf("duplicate checksum for %s", wanted)
		}
		found = candidate
	}
	if found == "" {
		return "", fmt.Errorf("missing checksum for %s", wanted)
	}
	return found, nil
}

func extractBinary(archiveName string, data []byte, destination string) error {
	if strings.HasSuffix(archiveName, ".zip") {
		return extractZip(data, destination)
	}
	if strings.HasSuffix(archiveName, ".tar.gz") {
		return extractTarGZ(data, destination)
	}
	return fmt.Errorf("unsupported release archive %q", archiveName)
}

func extractTarGZ(data []byte, destination string) error {
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("open release archive: %w", err)
	}
	defer gzipReader.Close()
	reader := tar.NewReader(gzipReader)
	found := false
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read release archive: %w", err)
		}
		if header.Typeflag != tar.TypeReg || filepath.Base(header.Name) != "kk" || header.Name != "kk" || found {
			return fmt.Errorf("release archive contains unexpected entry %q", header.Name)
		}
		content, err := io.ReadAll(io.LimitReader(reader, maxArchive+1))
		if err != nil {
			return fmt.Errorf("read kk from release archive: %w", err)
		}
		if len(content) > maxArchive {
			return fmt.Errorf("kk exceeds %d bytes", maxArchive)
		}
		if err := atomicExecutableWrite(destination, content); err != nil {
			return err
		}
		found = true
	}
	if !found {
		return fmt.Errorf("release archive does not contain kk")
	}
	return nil
}

func extractZip(data []byte, destination string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("open release archive: %w", err)
	}
	found := false
	for _, file := range reader.File {
		if file.Name != "kk.exe" || !file.Mode().IsRegular() || found {
			return fmt.Errorf("release archive contains unexpected entry %q", file.Name)
		}
		opened, err := file.Open()
		if err != nil {
			return err
		}
		content, readErr := io.ReadAll(io.LimitReader(opened, maxArchive+1))
		closeErr := opened.Close()
		if readErr != nil {
			return readErr
		}
		if closeErr != nil {
			return closeErr
		}
		if len(content) > maxArchive {
			return fmt.Errorf("kk.exe exceeds %d bytes", maxArchive)
		}
		if err := atomicExecutableWrite(destination, content); err != nil {
			return err
		}
		found = true
	}
	if !found {
		return fmt.Errorf("release archive does not contain kk.exe")
	}
	return nil
}

func atomicExecutableWrite(name string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(name), 0o700); err != nil {
		return err
	}
	return os.WriteFile(name, data, 0o700)
}

func verifyExecutableDirectory(executable string) error {
	dir := filepath.Dir(executable)
	temp, err := os.CreateTemp(dir, ".kevinkit-write-test-*")
	if err != nil {
		return fmt.Errorf("executable directory is not writable: %w", err)
	}
	name := temp.Name()
	temp.Close()
	os.Remove(name)
	return nil
}

func copyExecutable(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return err
	}
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o700)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	syncErr := output.Sync()
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}

func restoreExecutable(backup, destination string) error {
	staged, err := stageExecutable(backup, destination)
	if err != nil {
		return err
	}
	return replaceExecutable(staged, destination)
}

func installExecutable(source, destination string) error {
	staged, err := stageExecutable(source, destination)
	if err != nil {
		return err
	}
	return replaceExecutable(staged, destination)
}

func stageExecutable(source, destination string) (string, error) {
	input, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer input.Close()
	temp, err := os.CreateTemp(filepath.Dir(destination), ".kevinkit-binary-*")
	if err != nil {
		return "", err
	}
	name := temp.Name()
	defer func() {
		if temp != nil {
			temp.Close()
		}
	}()
	if err := temp.Chmod(0o755); err != nil {
		os.Remove(name)
		return "", err
	}
	if _, err := io.Copy(temp, input); err != nil {
		os.Remove(name)
		return "", err
	}
	if err := temp.Sync(); err != nil {
		os.Remove(name)
		return "", err
	}
	if err := temp.Close(); err != nil {
		os.Remove(name)
		return "", err
	}
	temp = nil
	return name, nil
}
