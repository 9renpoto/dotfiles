package macos

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectRunnerAsset(t *testing.T) {
	for _, tc := range []struct {
		name, arch, tag string
		assets          []Asset
		wantErr         bool
	}{
		{name: "selects matching asset", arch: "arm64", tag: "v2.330.0", assets: []Asset{{Name: "actions-runner-osx-arm64-2.330.0.tar.gz", URL: "https://example.invalid/runner.tar.gz", Digest: "sha256:" + strings.Repeat("a", 64)}}},
		{name: "missing asset", arch: "x64", tag: "v2.330.0", wantErr: true},
		{name: "missing digest", arch: "arm64", tag: "v2.330.0", assets: []Asset{{Name: "actions-runner-osx-arm64-2.330.0.tar.gz", URL: "https://example.invalid/runner.tar.gz"}}, wantErr: true},
		{name: "unsupported digest", arch: "arm64", tag: "v2.330.0", assets: []Asset{{Name: "actions-runner-osx-arm64-2.330.0.tar.gz", URL: "https://example.invalid/runner.tar.gz", Digest: "sha512:bad"}}, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SelectRunnerAsset(tc.tag, tc.arch, tc.assets)
			if (err != nil) != tc.wantErr {
				t.Fatalf("SelectRunnerAsset() = %#v, %v", got, err)
			}
		})
	}
}

func TestLinkDecision(t *testing.T) {
	for _, tc := range []struct {
		name                string
		exists, symlink     bool
		current, source     string
		wantCreate, wantErr bool
	}{
		{name: "missing target", source: "/src/config", wantCreate: true},
		{name: "same symlink", exists: true, symlink: true, current: "/src/config", source: "/src/config"},
		{name: "different symlink", exists: true, symlink: true, current: "/other/config", source: "/src/config", wantCreate: true},
		{name: "regular file", exists: true, source: "/src/config", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			create, err := LinkDecision(tc.exists, tc.symlink, tc.current, tc.source)
			if (err != nil) != tc.wantErr || create != tc.wantCreate {
				t.Fatalf("LinkDecision() = %v, %v", create, err)
			}
		})
	}
}

func TestVerifyArchive(t *testing.T) {
	data := "synthetic runner archive"
	sum := sha256.Sum256([]byte(data))
	want := hex.EncodeToString(sum[:])
	for _, tc := range []struct {
		name, expected string
		wantErr        bool
	}{
		{name: "matching digest", expected: want},
		{name: "mismatching digest", expected: strings.Repeat("0", 64), wantErr: true},
		{name: "malformed digest", expected: "not-a-digest", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyArchive(strings.NewReader(data), tc.expected)
			if (err != nil) != tc.wantErr {
				t.Fatalf("VerifyArchive() error = %v", err)
			}
		})
	}
}

func TestParseCredentialExports(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		wantErr     bool
	}{
		{name: "valid exports", input: "export GRAFANA_CLOUD_PROMETHEUS_URL='https://metrics.example/api'\nexport GRAFANA_CLOUD_PROMETHEUS_USERNAME='tenant'\\''s'\nexport GRAFANA_CLOUD_API_KEY='secret-value'\n"},
		{name: "rejects shell expression", input: "export GRAFANA_CLOUD_API_KEY=$(touch /tmp/unexpected)\n", wantErr: true},
		{name: "rejects malformed line", input: "GRAFANA_CLOUD_API_KEY=\"secret\"\n", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseCredentialExports(strings.NewReader(tc.input))
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseCredentialExports() = %#v, %v", got, err)
			}
			if tc.name == "valid exports" && got["GRAFANA_CLOUD_API_KEY"] != "secret-value" {
				t.Fatal("secret value was not parsed")
			}
			if tc.name == "valid exports" && got["GRAFANA_CLOUD_PROMETHEUS_USERNAME"] != "tenant's" {
				t.Fatalf("shell-escaped value = %q", got["GRAFANA_CLOUD_PROMETHEUS_USERNAME"])
			}
		})
	}
}

func TestSetupObservabilityUsesTempConfigurationAndRestrictsCredentials(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home space")
	config := filepath.Join(root, "redirected config")
	source := filepath.Join(config, "alloy")
	if err := os.MkdirAll(source, 0o700); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{"config.alloy": "logging {}\n", "config.env": "export GRAFANA_CLOUD_PROMETHEUS_URL='https://metrics.example/api'\nexport GRAFANA_CLOUD_PROMETHEUS_USERNAME='tenant'\nexport GRAFANA_CLOUD_API_KEY='fixture-secret'\n", "extra-args.txt": "--disable-reporting\n"}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(source, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	prefix := filepath.Join(root, "brew prefix")
	var calls []string
	run := func(_ context.Context, name string, env map[string]string, args ...string) ([]byte, error) {
		calls = append(calls, name+" "+strings.Join(args, " "))
		if name == "alloy" && env["GRAFANA_CLOUD_API_KEY"] != "fixture-secret" {
			t.Fatal("credential environment was not passed to alloy validation")
		}
		if name == "brew" && len(args) == 1 && args[0] == "--prefix" {
			return []byte(prefix + "\n"), nil
		}
		return nil, nil
	}
	if err := SetupObservability(context.Background(), "darwin", home, config, func(string) error { return nil }, run, &strings.Builder{}); err != nil {
		t.Fatal(err)
	}
	targetDir := filepath.Join(prefix, "etc", "alloy")
	for _, name := range []string{"config.alloy", "config.env", "extra-args.txt"} {
		target := filepath.Join(targetDir, name)
		info, err := os.Lstat(target)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("%s is not a symlink", target)
		}
	}
	info, err := os.Stat(filepath.Join(source, "config.env"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("credential mode = %o", info.Mode().Perm())
	}
	if len(calls) != 3 {
		t.Fatalf("commands = %v", calls)
	}
}

func TestSetupObservabilityRejectsExistingRegularFile(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	config := filepath.Join(root, "config")
	source := filepath.Join(config, "alloy")
	if err := os.MkdirAll(source, 0o700); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{"config.alloy": "x", "config.env": "export GRAFANA_CLOUD_PROMETHEUS_URL='u'\nexport GRAFANA_CLOUD_PROMETHEUS_USERNAME='n'\nexport GRAFANA_CLOUD_API_KEY='k'\n", "extra-args.txt": "x"} {
		if err := os.WriteFile(filepath.Join(source, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	prefix := filepath.Join(root, "brew")
	targetDir := filepath.Join(prefix, "etc", "alloy")
	if err := os.MkdirAll(targetDir, 0o700); err != nil {
		t.Fatal(err)
	}
	existing := filepath.Join(targetDir, "config.alloy")
	if err := os.WriteFile(existing, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	run := func(_ context.Context, name string, _ map[string]string, args ...string) ([]byte, error) {
		if name == "brew" && len(args) == 1 {
			return []byte(prefix), nil
		}
		return nil, nil
	}
	err := SetupObservability(context.Background(), "darwin", home, config, func(string) error { return nil }, run, &strings.Builder{})
	if err == nil {
		t.Fatal("expected existing regular file to be rejected")
	}
	content, readErr := os.ReadFile(existing)
	if readErr != nil || string(content) != "keep" {
		t.Fatalf("existing content = %q, %v", content, readErr)
	}
}

func TestLinkConfigHandlesSymlinksAndRegularFiles(t *testing.T) {
	for _, tc := range []struct {
		name    string
		kind    string
		wantErr bool
	}{
		{name: "missing link"},
		{name: "same link", kind: "same"},
		{name: "replace old link", kind: "different"},
		{name: "reject regular file", kind: "file", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			source := filepath.Join(root, "source")
			target := filepath.Join(root, "target")
			if err := os.WriteFile(source, []byte("source"), 0o600); err != nil {
				t.Fatal(err)
			}
			if tc.kind == "same" {
				if err := os.Symlink(source, target); err != nil {
					t.Fatal(err)
				}
			}
			if tc.kind == "different" {
				if err := os.Symlink(filepath.Join(root, "old"), target); err != nil {
					t.Fatal(err)
				}
			}
			if tc.kind == "file" {
				if err := os.WriteFile(target, []byte("keep"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			err := linkConfig(source, target)
			if (err != nil) != tc.wantErr {
				t.Fatalf("linkConfig() error = %v", err)
			}
			if tc.kind != "file" {
				got, err := os.Readlink(target)
				if err != nil || got != source {
					t.Fatalf("link target = %q, %v", got, err)
				}
			}
		})
	}
}

func TestSetupRunnerVerifiesArchiveBeforeExtractionAndNeverLogsToken(t *testing.T) {
	for _, tc := range []struct {
		name    string
		digest  string
		wantErr bool
	}{
		{name: "verified archive", digest: "valid"},
		{name: "checksum mismatch", digest: "invalid", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			runnerDir := filepath.Join(home, "runner with spaces")
			archiveData := []byte("synthetic runner archive")
			sum := sha256.Sum256(archiveData)
			digest := hex.EncodeToString(sum[:])
			if tc.digest == "invalid" {
				digest = strings.Repeat("0", 64)
			}
			release, _ := json.Marshal(map[string]any{"tag_name": "v2.330.0", "assets": []Asset{{Name: "actions-runner-osx-arm64-2.330.0.tar.gz", URL: "https://example.invalid/runner.tar.gz", Digest: "sha256:" + digest}}})
			var calls []string
			var output strings.Builder
			deps := RunnerDependencies{GOOS: "darwin", ARCH: "arm64", UID: 501, Lookup: func(string) error { return nil },
				Run: func(_ context.Context, dir, name string, args ...string) ([]byte, error) {
					calls = append(calls, name+" "+strings.Join(args, " "))
					if name == "gh" && len(args) > 0 && args[0] == "repo" {
						return []byte("PRIVATE\n"), nil
					}
					if name == "gh" && len(args) > 0 && args[0] == "api" && strings.Contains(strings.Join(args, " "), "actions/runner/releases/latest") {
						return release, nil
					}
					if name == "gh" && len(args) > 0 && args[0] == "api" && strings.Contains(strings.Join(args, " "), "registration-token") {
						return []byte("registration-secret\n"), nil
					}
					return nil, nil
				},
				Download: func(_ context.Context, _, dest string) error { return os.WriteFile(dest, archiveData, 0o600) },
				Extract: func(_ context.Context, _, dest string) error {
					if tc.digest == "invalid" {
						t.Fatal("extraction ran before checksum passed")
					}
					return os.MkdirAll(dest, 0o700)
				},
			}
			deps.Lookup = func(string) error { return nil }
			err := SetupRunner(context.Background(), RunnerOptions{Repository: "owner/private", Name: "test-runner", RunnerDir: runnerDir}, home, deps, &output)
			if (err != nil) != tc.wantErr {
				t.Fatalf("SetupRunner() error = %v", err)
			}
			if strings.Contains(output.String(), "registration-secret") {
				t.Fatal("registration token was logged")
			}
			if tc.wantErr {
				for _, call := range calls {
					if strings.Contains(call, "registration-token") {
						t.Fatal("registration started before archive verification")
					}
				}
			}
		})
	}
}

func TestSetupRunnerRejectsUnsafePrerequisites(t *testing.T) {
	for _, tc := range []struct {
		name, goos, arch, failAt string
		uid                      int
		populate                 bool
	}{
		{name: "wrong operating system", goos: "linux", arch: "arm64", uid: 501, failAt: "early"},
		{name: "root user", goos: "darwin", arch: "arm64", uid: 0, failAt: "early"},
		{name: "unsupported architecture", goos: "darwin", arch: "i386", uid: 501, failAt: "early"},
		{name: "authentication required", goos: "darwin", arch: "arm64", uid: 501, failAt: "auth"},
		{name: "private repository required", goos: "darwin", arch: "arm64", uid: 501, failAt: "visibility"},
		{name: "nonempty runner directory", goos: "darwin", arch: "arm64", uid: 501, populate: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			runnerDir := filepath.Join(home, "runner")
			if tc.populate {
				if err := os.MkdirAll(runnerDir, 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(runnerDir, "keep"), []byte("x"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			var calls []string
			deps := RunnerDependencies{GOOS: tc.goos, ARCH: tc.arch, UID: tc.uid, Lookup: func(string) error { return nil }, Download: func(context.Context, string, string) error { t.Fatal("unexpected download"); return nil }, Extract: func(context.Context, string, string) error { t.Fatal("unexpected extraction"); return nil }, Run: func(_ context.Context, _, name string, args ...string) ([]byte, error) {
				calls = append(calls, name+" "+strings.Join(args, " "))
				if tc.failAt == "auth" && name == "gh" && len(args) > 0 && args[0] == "auth" {
					return nil, os.ErrPermission
				}
				if name == "gh" && len(args) > 0 && args[0] == "repo" {
					if tc.failAt == "visibility" {
						return []byte("PUBLIC"), nil
					}
					return []byte("PRIVATE"), nil
				}
				return nil, nil
			}}
			err := SetupRunner(context.Background(), RunnerOptions{Repository: "owner/private", Name: "test", RunnerDir: runnerDir}, home, deps, &strings.Builder{})
			if err == nil {
				t.Fatal("expected prerequisite error")
			}
			if tc.failAt == "early" && len(calls) > 0 {
				t.Fatalf("ran commands before platform/user checks: %v", calls)
			}
		})
	}
}
