package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestRefreshWindowsPathPreservesUnicode(t *testing.T) {
	for _, native := range []bool{false, true} {
		name := "UTF-8 output contract"
		if native {
			name = "Windows PowerShell legacy code page"
		}
		t.Run(name, func(t *testing.T) {
			if native && runtime.GOOS != "windows" {
				t.Skip("requires Windows PowerShell")
			}
			powershell := ""
			if native {
				var err error
				powershell, err = exec.LookPath("powershell.exe")
				if err != nil {
					t.Fatal(err)
				}
			}
			t.Setenv("PATH", "")
			err := refreshWindowsPath(context.Background(), func(ctx context.Context, command string, args ...string) ([]byte, error) {
				if command != "powershell.exe" || len(args) != 3 {
					t.Fatalf("unexpected command: %s %q", command, args)
				}
				script := args[2]
				if !native {
					encoding := strings.Index(script, "[Console]::OutputEncoding")
					if encoding < 0 || encoding > strings.Index(script, "ConvertTo-Json") {
						t.Fatal("PowerShell output encoding must be set before JSON serialization")
					}
					return []byte(`{"Machine":"C:\\表","User":"D:\\利用者"}`), nil
				}
				script = strings.ReplaceAll(script, "[Environment]::GetEnvironmentVariable('Path','Machine')", "('C:\\'+[char]0x8868)")
				script = strings.ReplaceAll(script, "[Environment]::GetEnvironmentVariable('Path','User')", "('D:\\'+[char]0x5229+[char]0x7528+[char]0x8005)")
				script = "[Console]::OutputEncoding=[Text.Encoding]::GetEncoding(932); " + script
				return exec.CommandContext(ctx, powershell, "-NoProfile", "-Command", script).CombinedOutput()
			})
			if err != nil {
				t.Fatal(err)
			}
			want := strings.Join([]string{`C:\表`, `D:\利用者`}, string(os.PathListSeparator))
			if got := os.Getenv("PATH"); got != want {
				t.Fatalf("PATH = %q, want %q", got, want)
			}
		})
	}
}

func TestRunPackages(t *testing.T) {
	root := t.TempDir()
	manifest := "Brewfile"
	commandName := "brew"
	wantArgs := []string{"bundle", "--file=" + filepath.Join(root, manifest)}
	if runtime.GOOS == "windows" {
		manifest = "winget-packages.json"
		commandName = "winget.exe"
		wantArgs = []string{"import", "--import-file", filepath.Join(root, manifest), "--ignore-unavailable", "--ignore-versions", "--accept-package-agreements", "--accept-source-agreements"}
	}
	if err := os.WriteFile(filepath.Join(root, manifest), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	var gotName string
	var gotArgs []string
	var stdout, stderr bytes.Buffer
	wantErr := errors.New("installer failed")
	err := run(context.Background(), []string{"packages"}, root, &stdout, &stderr,
		func(_ context.Context, name string, args ...string) error {
			gotName, gotArgs = name, args
			return wantErr
		})
	if !errors.Is(err, wantErr) {
		t.Fatalf("run() error = %v, want %v", err, wantErr)
	}
	if gotName != commandName || !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("command = %q %q", gotName, gotArgs)
	}
	if runtime.GOOS == "windows" && (stdout.Len() != 0 || stderr.Len() != 0) {
		t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	if runtime.GOOS != "windows" && (stdout.String() != "Installing Homebrew packages from Brewfile.\n" || stderr.Len() == 0) {
		t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunPackagesRejectsMissingManifest(t *testing.T) {
	err := run(context.Background(), []string{"packages"}, t.TempDir(), &bytes.Buffer{}, &bytes.Buffer{},
		func(context.Context, string, ...string) error { t.Fatal("unexpected command"); return nil })
	if err == nil {
		t.Fatal("expected missing manifest error")
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	err := run(context.Background(), []string{"unknown"}, t.TempDir(), &bytes.Buffer{}, &bytes.Buffer{},
		func(context.Context, string, ...string) error { t.Fatal("unexpected command"); return nil })
	if err == nil {
		t.Fatal("expected unknown command error")
	}
}

func TestGitCredentialsSkipsWhenGitHubCLIIsMissing(t *testing.T) {
	var stdout bytes.Buffer
	err := gitCredentials(context.Background(), &stdout,
		func(context.Context, string, ...string) error { t.Fatal("unexpected command"); return nil },
		func(string) (string, error) { return "", errors.New("not found") })
	if err != nil || stdout.String() != "gh not found; skipping 'gh auth setup-git'.\n" {
		t.Fatalf("gitCredentials() = %v, output %q", err, stdout.String())
	}
}

func TestGitCredentialsSkipsWhenUnauthenticated(t *testing.T) {
	var stdout bytes.Buffer
	var calls [][]string
	err := gitCredentials(context.Background(), &stdout,
		func(_ context.Context, name string, args ...string) error {
			calls = append(calls, append([]string{name}, args...))
			return errors.New("not authenticated")
		}, func(string) (string, error) { return "/usr/bin/gh", nil })
	if err != nil || !reflect.DeepEqual(calls, [][]string{{"gh", "auth", "status"}}) {
		t.Fatalf("gitCredentials() = %v, calls %q", err, calls)
	}
	if !strings.Contains(stdout.String(), "gh is not authenticated") {
		t.Fatalf("output = %q", stdout.String())
	}
}

func TestGitCredentialsConfiguresAuthenticatedGitHubCLI(t *testing.T) {
	var stdout bytes.Buffer
	var calls [][]string
	err := gitCredentials(context.Background(), &stdout,
		func(_ context.Context, name string, args ...string) error {
			calls = append(calls, append([]string{name}, args...))
			return nil
		}, func(string) (string, error) { return "/usr/bin/gh", nil })
	want := [][]string{{"gh", "auth", "status"}, {"gh", "auth", "setup-git"}}
	if err != nil || !reflect.DeepEqual(calls, want) {
		t.Fatalf("gitCredentials() = %v, calls %q", err, calls)
	}
}
