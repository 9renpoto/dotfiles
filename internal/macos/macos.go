package macos

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Asset struct {
	Name   string `json:"name"`
	URL    string `json:"browser_download_url"`
	Digest string `json:"digest"`
}

type RunnerOptions struct{ Repository, Name, Labels, RunnerDir string }
type RunnerDependencies struct {
	GOOS, ARCH string
	UID        int
	Lookup     func(string) error
	Run        func(context.Context, string, string, ...string) ([]byte, error)
	Download   func(context.Context, string, string) error
	Extract    func(context.Context, string, string) error
}

func SetupRunner(ctx context.Context, options RunnerOptions, home string, deps RunnerDependencies, output io.Writer) error {
	if deps.GOOS != "darwin" {
		return fmt.Errorf("runner setup is only supported on macOS")
	}
	if deps.UID == 0 {
		return fmt.Errorf("do not configure a GitHub Actions runner as root")
	}
	arch := deps.ARCH
	switch arch {
	case "x86_64", "amd64":
		arch = "x64"
	case "arm64":
	default:
		return fmt.Errorf("unsupported architecture: %s", deps.ARCH)
	}
	owner, repo, ok := strings.Cut(options.Repository, "/")
	if !ok || owner == "" || repo == "" || strings.Contains(repo, "/") {
		return fmt.Errorf("repository must be OWNER/PRIVATE_REPOSITORY")
	}
	if options.RunnerDir == "" {
		options.RunnerDir = os.Getenv("ACTIONS_RUNNER_DIR")
	}
	if options.RunnerDir == "" {
		options.RunnerDir = filepath.Join(home, "actions-runner")
	}
	if options.Name == "" {
		if name, err := deps.Run(ctx, "", "scutil", "--get", "LocalHostName"); err == nil && strings.TrimSpace(string(name)) != "" {
			options.Name = strings.TrimSpace(string(name))
		} else if name, err := deps.Run(ctx, "", "hostname", "-s"); err == nil {
			options.Name = strings.TrimSpace(string(name))
		}
	}
	if options.Name == "" {
		return fmt.Errorf("could not determine runner name")
	}
	if options.Labels == "" {
		options.Labels = "macos,ai-development"
	}
	if err := deps.Lookup("gh"); err != nil {
		return fmt.Errorf("GitHub CLI is required")
	}
	if _, err := deps.Run(ctx, "", "gh", "auth", "status"); err != nil {
		return fmt.Errorf("authenticate GitHub CLI with gh auth login first")
	}
	visibility, err := deps.Run(ctx, "", "gh", "repo", "view", options.Repository, "--json", "visibility", "--jq", ".visibility")
	if err != nil || strings.TrimSpace(string(visibility)) != "PRIVATE" {
		return fmt.Errorf("%s is not a private repository", options.Repository)
	}
	if _, err := os.Stat(filepath.Join(options.RunnerDir, ".runner")); err == nil {
		fmt.Fprintf(output, "A runner is already configured in %s.\n", options.RunnerDir)
		status, err := deps.Run(ctx, options.RunnerDir, "./svc.sh", "status")
		if len(status) > 0 {
			_, _ = output.Write(status)
		}
		return err
	}
	if entries, err := os.ReadDir(options.RunnerDir); err == nil && len(entries) > 0 {
		return fmt.Errorf("runner directory is not empty; choose another ACTIONS_RUNNER_DIR")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	data, err := deps.Run(ctx, "", "gh", "api", "repos/actions/runner/releases/latest")
	if err != nil {
		return fmt.Errorf("could not read the latest runner release")
	}
	var release struct {
		Tag    string  `json:"tag_name"`
		Assets []Asset `json:"assets"`
	}
	if err := json.Unmarshal(data, &release); err != nil {
		return fmt.Errorf("invalid runner release metadata: %w", err)
	}
	asset, err := SelectRunnerAsset(release.Tag, arch, release.Assets)
	if err != nil {
		return err
	}
	temp, err := os.MkdirTemp("", "dotfiles-runner-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	archive := filepath.Join(temp, "actions-runner.tar.gz")
	if err := deps.Download(ctx, asset.URL, archive); err != nil {
		return fmt.Errorf("could not download runner archive: %w", err)
	}
	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	verifyErr := VerifyArchive(f, asset.SHA256)
	closeErr := f.Close()
	if verifyErr != nil {
		return verifyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if err := os.MkdirAll(options.RunnerDir, 0o755); err != nil {
		return err
	}
	if err := deps.Extract(ctx, archive, options.RunnerDir); err != nil {
		return fmt.Errorf("could not extract verified runner archive: %w", err)
	}
	tokenData, err := deps.Run(ctx, "", "gh", "api", "--method", "POST", fmt.Sprintf("repos/%s/actions/runners/registration-token", options.Repository), "--jq", ".token")
	if err != nil {
		return fmt.Errorf("could not obtain runner registration token")
	}
	token := strings.TrimSpace(string(tokenData))
	if token == "" {
		return fmt.Errorf("GitHub returned an empty runner registration token")
	}
	configArgs := []string{"./config.sh", "--unattended", "--url", "https://github.com/" + options.Repository, "--token", token, "--name", options.Name, "--labels", options.Labels, "--work", "_work"}
	_, configErr := deps.Run(ctx, options.RunnerDir, configArgs[0], configArgs[1:]...)
	token = ""
	if configErr != nil {
		return fmt.Errorf("runner configuration failed")
	}
	for _, args := range [][]string{{"./svc.sh", "install"}, {"./svc.sh", "start"}, {"./svc.sh", "status"}} {
		if _, err := deps.Run(ctx, options.RunnerDir, args[0], args[1:]...); err != nil {
			return fmt.Errorf("runner service command %s failed: %w", args[1], err)
		}
	}
	fmt.Fprintf(output, "GitHub Actions runner %s is registered for %s.\n", options.Name, options.Repository)
	return nil
}

func DownloadArchive(ctx context.Context, url, destination string) error {
	client := &http.Client{Timeout: 5 * time.Minute}
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		response, err := client.Do(req)
		if err != nil {
			last = err
			continue
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			response.Body.Close()
			last = fmt.Errorf("download returned HTTP %s", response.Status)
			continue
		}
		file, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
		if err != nil {
			response.Body.Close()
			return err
		}
		_, copyErr := io.Copy(file, response.Body)
		bodyErr := response.Body.Close()
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		if bodyErr != nil {
			return bodyErr
		}
		if closeErr != nil {
			return closeErr
		}
		return nil
	}
	return last
}

func ExtractRunner(ctx context.Context, archive, destination string) error {
	cmd := exec.CommandContext(ctx, "tar", "-xzf", archive, "-C", destination)
	return cmd.Run()
}

type CommandRunner func(context.Context, string, map[string]string, ...string) ([]byte, error)

func SetupObservability(ctx context.Context, goos, home, configHome string, lookPath func(string) error, run CommandRunner, output io.Writer) error {
	if goos != "darwin" {
		return fmt.Errorf("observability setup is only supported on macOS")
	}
	for _, command := range []string{"brew", "alloy"} {
		if err := lookPath(command); err != nil {
			return fmt.Errorf("%s is required", command)
		}
	}
	if configHome == "" {
		configHome = filepath.Join(home, ".config")
	}
	sourceDir := filepath.Join(configHome, "alloy")
	sourceConfig := filepath.Join(sourceDir, "config.alloy")
	sourceEnv := filepath.Join(sourceDir, "config.env")
	sourceArgs := filepath.Join(sourceDir, "extra-args.txt")
	for _, path := range []string{sourceConfig, sourceEnv, sourceArgs} {
		info, err := os.Stat(path)
		if err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("required Alloy configuration file is missing: %s", path)
		}
	}
	file, err := os.Open(sourceEnv)
	if err != nil {
		return err
	}
	credentials, parseErr := ParseCredentialExports(file)
	closeErr := file.Close()
	if parseErr != nil {
		return parseErr
	}
	if closeErr != nil {
		return closeErr
	}
	if _, err := run(ctx, "alloy", credentials, "validate", sourceConfig); err != nil {
		return fmt.Errorf("Alloy configuration validation failed: %w", err)
	}
	prefix, err := run(ctx, "brew", nil, "--prefix")
	if err != nil {
		return fmt.Errorf("could not find Homebrew prefix: %w", err)
	}
	targetDir := filepath.Join(strings.TrimSpace(string(prefix)), "etc", "alloy")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}
	if err := os.Chmod(sourceEnv, 0o600); err != nil {
		return fmt.Errorf("could not restrict credential file permissions: %w", err)
	}
	for _, pair := range [][2]string{{sourceConfig, filepath.Join(targetDir, "config.alloy")}, {sourceEnv, filepath.Join(targetDir, "config.env")}, {sourceArgs, filepath.Join(targetDir, "extra-args.txt")}} {
		if err := linkConfig(pair[0], pair[1]); err != nil {
			return err
		}
	}
	if _, err := run(ctx, "brew", nil, "services", "restart", "grafana/grafana/alloy"); err != nil {
		return fmt.Errorf("could not restart Grafana Alloy: %w", err)
	}
	fmt.Fprintln(output, "Grafana Alloy is configured and running.")
	fmt.Fprintln(output, "Check readiness at http://127.0.0.1:12345/-/ready")
	return nil
}

func linkConfig(source, target string) error {
	info, err := os.Lstat(target)
	exists := err == nil
	symlink := exists && info.Mode()&os.ModeSymlink != 0
	current := ""
	if symlink {
		current, err = os.Readlink(target)
		if err != nil {
			return err
		}
	}
	create, err := LinkDecision(exists, symlink, current, source)
	if err != nil {
		return fmt.Errorf("%s already exists; move it aside before running setup", target)
	}
	if !create {
		return nil
	}
	if symlink {
		if err := os.Remove(target); err != nil {
			return err
		}
	}
	if err := os.Symlink(source, target); err != nil {
		return err
	}
	return nil
}

type RunnerAsset struct{ Version, Name, URL, SHA256 string }

func SelectRunnerAsset(tag, arch string, assets []Asset) (RunnerAsset, error) {
	version := strings.TrimPrefix(tag, "v")
	if version == "" || (arch != "x64" && arch != "arm64") {
		return RunnerAsset{}, fmt.Errorf("invalid runner release metadata")
	}
	name := fmt.Sprintf("actions-runner-osx-%s-%s.tar.gz", arch, version)
	for _, asset := range assets {
		if asset.Name != name {
			continue
		}
		if asset.URL == "" || !strings.HasPrefix(asset.Digest, "sha256:") {
			return RunnerAsset{}, fmt.Errorf("runner asset is missing its URL or SHA-256 digest")
		}
		digest := strings.TrimPrefix(asset.Digest, "sha256:")
		decoded, err := hex.DecodeString(digest)
		if err != nil || len(decoded) != sha256.Size {
			return RunnerAsset{}, fmt.Errorf("runner asset has an invalid SHA-256 digest")
		}
		return RunnerAsset{Version: version, Name: name, URL: asset.URL, SHA256: strings.ToLower(digest)}, nil
	}
	return RunnerAsset{}, fmt.Errorf("runner release asset %s was not found", name)
}

func VerifyArchive(archive io.Reader, expectedSHA256 string) error {
	expected, err := hex.DecodeString(expectedSHA256)
	if err != nil || len(expected) != sha256.Size {
		return fmt.Errorf("invalid expected SHA-256 digest")
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, archive); err != nil {
		return err
	}
	if !strings.EqualFold(hex.EncodeToString(hash.Sum(nil)), expectedSHA256) {
		return fmt.Errorf("runner archive checksum verification failed")
	}
	return nil
}

func LinkDecision(exists, symlink bool, currentTarget, source string) (create bool, err error) {
	if symlink && currentTarget == source {
		return false, nil
	}
	if exists && !symlink {
		return false, fmt.Errorf("target already exists and is not a symlink")
	}
	return true, nil
}

func ParseCredentialExports(input io.Reader) (map[string]string, error) {
	allowed := map[string]bool{
		"GRAFANA_CLOUD_PROMETHEUS_URL":      true,
		"GRAFANA_CLOUD_PROMETHEUS_USERNAME": true,
		"GRAFANA_CLOUD_API_KEY":             true,
	}
	values := make(map[string]string, len(allowed))
	scanner := bufio.NewScanner(input)
	for line := 1; scanner.Scan(); line++ {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		name, value, found := strings.Cut(strings.TrimPrefix(text, "export "), "=")
		if !strings.HasPrefix(text, "export ") || !found || !allowed[name] {
			return nil, fmt.Errorf("invalid credential environment line %d", line)
		}
		if _, duplicate := values[name]; duplicate {
			return nil, fmt.Errorf("duplicate credential environment variable %s", name)
		}
		decoded, err := unquoteShellValue(value)
		if err != nil {
			return nil, fmt.Errorf("invalid quoted value on credential environment line %d", line)
		}
		values[name] = decoded
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for key := range allowed {
		if strings.TrimSpace(values[key]) == "" {
			return nil, fmt.Errorf("missing required credential environment variable %s", key)
		}
	}
	return values, nil
}

func unquoteShellValue(value string) (string, error) {
	var result strings.Builder
	quoted := false
	started := false
	for i := 0; i < len(value); i++ {
		if quoted {
			if value[i] == '\'' {
				quoted = false
			} else {
				result.WriteByte(value[i])
			}
			continue
		}
		if value[i] == '\'' {
			quoted = true
			started = true
			continue
		}
		if value[i] == '\\' && i+1 < len(value) && value[i+1] == '\'' {
			result.WriteByte('\'')
			i++
			started = true
			continue
		}
		return "", fmt.Errorf("unquoted shell character")
	}
	if quoted || !started {
		return "", fmt.Errorf("unbalanced shell quotes")
	}
	return result.String(), nil
}
