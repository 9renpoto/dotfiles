package windows

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestMergePathEntries(t *testing.T) {
	for _, tc := range []struct {
		name, machine, user, process string
		want                         []string
	}{
		{name: "merges unique entries", machine: `C:\Machine;C:\Shared`, user: `C:\User`, process: `C:\Machine;C:\Process`, want: []string{`C:\Machine`, `C:\Shared`, `C:\User`, `C:\Process`}},
		{name: "deduplicates case-insensitively", machine: `C:\Tools`, process: `c:\tools;C:\Other`, want: []string{`C:\Tools`, `C:\Other`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := MergePathEntries(tc.machine, tc.user, tc.process, ";")
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("MergePathEntries() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestWingetImportArgs(t *testing.T) {
	got := WingetImportArgs(`C:\work tree\winget-packages.json`)
	want := []string{"import", "--import-file", `C:\work tree\winget-packages.json`, "--ignore-unavailable", "--ignore-versions", "--accept-package-agreements", "--accept-source-agreements"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("WingetImportArgs() = %#v, want %#v", got, want)
	}
}

func TestProfileLoaderPreservesCustomizedProfiles(t *testing.T) {
	loader := ProfileLoader(`C:\Users\Test\.config\powershell\profile.ps1`)
	for _, tc := range []struct {
		name          string
		exists        bool
		content, want string
		write         bool
	}{
		{name: "missing", want: "\ufeff" + loader + "\n", write: true},
		{name: "empty", exists: true, content: "  ", want: "\ufeff" + loader + "\n", write: true},
		{name: "loader", exists: true, content: loader, want: "\ufeff" + loader + "\n", write: true},
		{name: "BOM loader", exists: true, content: "\ufeff" + loader + "\n", want: "\ufeff" + loader + "\n", write: true},
		{name: "custom", exists: true, content: "# customized", want: "# customized", write: false},
		{name: "BOM custom", exists: true, content: "\ufeff# customized", want: "\ufeff# customized", write: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, write := MergeProfile([]byte(tc.content), tc.exists, loader)
			if string(got) != tc.want || write != tc.write {
				t.Fatalf("MergeProfile() = %q, %v", got, write)
			}
		})
	}
}

func TestInstallProfilesPreservesUnicodePathsAndReapplication(t *testing.T) {
	documents := filepath.Join(t.TempDir(), "Documents")
	home := filepath.Join(t.TempDir(), "利用者")
	loader := ProfileLoader(filepath.Join(home, ".config", "powershell", "profile.ps1"))
	for _, name := range []string{"WindowsPowerShell", "PowerShell"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(documents, name, "profile.ps1")
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(loader+"\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			for run := 0; run < 2; run++ {
				preserved, err := InstallProfiles(documents, home)
				if err != nil {
					t.Fatal(err)
				}
				if len(preserved) != 0 {
					t.Fatalf("managed profiles treated as customized: %v", preserved)
				}
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.HasPrefix(content, []byte{0xef, 0xbb, 0xbf}) || !utf8.Valid(content) || string(content) != "\ufeff"+loader+"\n" {
					t.Fatalf("loader did not preserve UTF-8 BOM and Unicode path: %q", content)
				}
			}
		})
	}
}

func TestInstallProfilesIsIdempotentAndPreservesCustomizedFiles(t *testing.T) {
	documents := filepath.Join(t.TempDir(), "redirected Documents")
	home := filepath.Join(t.TempDir(), "user home")
	if _, err := InstallProfiles(documents, home); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{"WindowsPowerShell", "PowerShell"} {
		t.Run(dir, func(t *testing.T) {
			path := filepath.Join(documents, dir, "profile.ps1")
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(content), filepath.Join(home, ".config", "powershell", "profile.ps1")) {
				t.Fatalf("loader = %q", content)
			}
		})
	}
	custom := filepath.Join(documents, "PowerShell", "profile.ps1")
	if err := os.WriteFile(custom, []byte("# keep custom profile"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := InstallProfiles(documents, home); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(custom)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "# keep custom profile" {
		t.Fatalf("custom profile changed: %q", content)
	}
}

func TestSyncWslGhqRoot(t *testing.T) {
	for _, tc := range []struct {
		name        string
		failAt      string
		convert     string
		wantWarning string
	}{
		{name: "success", convert: "/mnt/c/user/ghq root\n"},
		{name: "wsl missing", failAt: "lookup", wantWarning: "WSL was not found"},
		{name: "conversion fails", failAt: "convert", wantWarning: "Could not convert"},
		{name: "mkdir fails", failAt: "mkdir", convert: "/mnt/c/ghq\n", wantWarning: "Could not create"},
		{name: "git fails", failAt: "git", convert: "/mnt/c/ghq\n", wantWarning: "Git was not available"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls [][]string
			var output strings.Builder
			lookup := func(string) error {
				if tc.failAt == "lookup" {
					return os.ErrNotExist
				}
				return nil
			}
			run := func(_ context.Context, name string, args ...string) ([]byte, error) {
				call := append([]string{name}, args...)
				calls = append(calls, call)
				if (tc.failAt == "convert" && len(args) > 1 && args[1] == "wslpath") || (tc.failAt == "mkdir" && len(args) > 1 && args[1] == "mkdir") || (tc.failAt == "git" && len(args) > 1 && args[1] == "git") {
					return nil, os.ErrNotExist
				}
				if len(args) > 1 && args[1] == "wslpath" {
					return []byte(tc.convert), nil
				}
				return nil, nil
			}
			if err := SyncWslGhqRoot(context.Background(), `C:\Users\Test\ghq root`, lookup, run, &output); err != nil {
				t.Fatal(err)
			}
			if tc.wantWarning != "" && !strings.Contains(output.String(), tc.wantWarning) {
				t.Fatalf("output = %q", output.String())
			}
			if tc.wantWarning == "" && (len(calls) != 3 || calls[1][4] != "/mnt/c/user/ghq root" || calls[2][6] != "/mnt/c/user/ghq root") {
				t.Fatalf("calls = %#v", calls)
			}
		})
	}
}

func TestInitializeHonorsSkipFlagsAndPersistsGhqRoot(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		skipPackages, skipApply bool
		want                    []string
	}{
		{name: "full", want: []string{"mise trust", "mise config", "winget.exe import", "refresh", "mise bootstrap", "wsl.exe --exec wslpath"}},
		{name: "skip packages", skipPackages: true, want: []string{"mise trust", "mise config", "mise bootstrap", "wsl.exe --exec wslpath"}},
		{name: "skip apply", skipApply: true, want: []string{"mise trust", "mise config", "winget.exe import", "refresh"}},
		{name: "skip both", skipPackages: true, skipApply: true, want: []string{"mise trust", "mise config"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			home := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "winget-packages.json"), []byte("{}"), 0o600); err != nil {
				t.Fatal(err)
			}
			var calls []string
			run := func(_ context.Context, name string, args ...string) ([]byte, error) {
				calls = append(calls, strings.Join(append([]string{name}, args...), " "))
				if name == "wsl.exe" && len(args) > 1 && args[1] == "wslpath" {
					return []byte("/mnt/c/ghq\n"), nil
				}
				return nil, nil
			}
			refresh := func(context.Context) error { calls = append(calls, "refresh"); return nil }
			err := RunInitialize(context.Background(), InitializeOptions{Checkout: root, GhqRoot: filepath.Join(home, "ghq root"), SkipPackages: tc.skipPackages, SkipApply: tc.skipApply}, home, func(string) error { return nil }, run, refresh, &strings.Builder{})
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range tc.want {
				found := false
				for _, got := range calls {
					if strings.HasPrefix(got, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing %q in calls %v", want, calls)
				}
			}
			if !strings.Contains(calls[1], filepath.Join(root, "mise.local.toml")) {
				t.Fatalf("mise config did not use local config: %v", calls)
			}
		})
	}
}
