package skills

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseDesired(t *testing.T) {
	got, err := ParseDesired(strings.NewReader("# note\ngithub/repo path/skill@v1 --agent codex\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Repository != "github/repo" || got[0].Name != "skill" || got[0].Skill != "path/skill@v1" || strings.Join(got[0].Flags, " ") != "--agent codex" {
		t.Fatalf("ParseDesired() = %#v", got)
	}
}

func TestOwnedSkillPathRejectsOutsideAndRootPaths(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".agents", "skills", "owned"), 0o700); err != nil {
		t.Fatal(err)
	}
	if got := OwnedSkillPath(home, filepath.Join(home, ".agents", "skills", "owned")); !got {
		t.Fatal("expected user skill path to be accepted")
	}
	for _, tc := range []struct{ name, path string }{
		{name: "outside", path: filepath.Join(home, "elsewhere", "skill")},
		{name: "skills root", path: filepath.Join(home, ".agents", "skills")},
		{name: "traversal", path: filepath.Join(home, ".agents", "skills", "..", "sensitive", "skills", "nested")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if OwnedSkillPath(home, tc.path) {
				t.Fatalf("OwnedSkillPath(%q) accepted", tc.path)
			}
		})
	}
}

func TestOwnedSkillPathRejectsSymlinkEscape(t *testing.T) {
	home, outside := t.TempDir(), t.TempDir()
	root := filepath.Join(home, ".agents", "skills")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "escape")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	if OwnedSkillPath(home, link) {
		t.Fatal("accepted symlink escaping the skills root")
	}
}

func TestSourceRepositoryMatchesSupportedURLs(t *testing.T) {
	for _, tc := range []struct{ name, source string }{
		{name: "https", source: "https://github.com/github/repo.git"},
		{name: "ssh", source: "git@github.com:github/repo.git"},
		{name: "host path", source: "github.com/github/repo"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SourceRepository(tc.source)
			if err != nil || got != "github/repo" {
				t.Fatalf("SourceRepository(%q) = %q, %v", tc.source, got, err)
			}
		})
	}
}

func TestRemovalCandidatesRequireOwnedUserCopyAndMatchingSource(t *testing.T) {
	for _, tc := range []struct {
		name, source, scope string
		want                bool
	}{
		{name: "owned user copy", source: "https://github.com/github/repo", scope: "user", want: true},
		{name: "different source", source: "https://github.com/other/repo", scope: "user"},
		{name: "project scope", source: "https://github.com/github/repo", scope: "project"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			path := filepath.Join(home, ".agents", "skills", "old")
			if err := os.MkdirAll(path, 0o700); err != nil {
				t.Fatal(err)
			}
			installed := []Installed{{SkillName: "old", SourceURL: tc.source, Path: path, Scope: tc.scope}}
			got := RemovalCandidates(home, []string{"github/repo\told"}, nil, installed)
			if (len(got) == 1) != tc.want {
				t.Fatalf("RemovalCandidates() = %#v, want match %v", got, tc.want)
			}
		})
	}
}

func TestRemoveLockEntryPreservesOtherEntriesAndMismatchedSources(t *testing.T) {
	input := []byte(`{"other":true,"skills":{"old":{"source":"github/repo"},"keep":{"source":"other/repo"}}}`)
	got, err := RemoveLockEntry(input, "github/repo", "old")
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatal(err)
	}
	skills := decoded["skills"].(map[string]any)
	if _, ok := skills["old"]; ok {
		t.Fatal("matching lock entry retained")
	}
	if _, ok := skills["keep"]; !ok || decoded["other"] != true {
		t.Fatal("unrelated lock data was lost")
	}
}

func TestSyncRecordsOnlySuccessfulInstalls(t *testing.T) {
	for _, tc := range []struct {
		name string
		fail string
		want string
	}{
		{name: "success", want: "github/repo\tskill\n"},
		{name: "failure", fail: "all", want: ""},
		{name: "codex succeeds", fail: "antigravity", want: "github/repo\tskill\n"},
		{name: "antigravity succeeds", fail: "codex", want: "github/repo\tskill\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			state := filepath.Join(home, "state", "installed.tsv")
			desired := []Desired{{Repository: "github/repo", Skill: "skill", Name: "skill"}}
			cmd := func(_ context.Context, name string, args ...string) ([]byte, error) {
				if len(args) == 3 && args[0] == "skill" && args[1] == "list" && args[2] == "--help" {
					return nil, nil
				}
				if tc.fail == "all" || (tc.fail != "" && args[5] == tc.fail) {
					return nil, os.ErrPermission
				}
				return nil, nil
			}
			if err := Sync(context.Background(), home, state, filepath.Join(home, ".agents", ".skill-lock.json"), desired, cmd, &strings.Builder{}); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(state)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tc.want {
				t.Fatalf("state = %q, want %q", data, tc.want)
			}
		})
	}
}

func TestSyncInstallDestinations(t *testing.T) {
	for _, tc := range []struct {
		name     string
		flags    []string
		explicit bool
	}{
		{name: "force", flags: []string{"--force"}},
		{name: "pin", flags: []string{"--pin", "v1"}},
		{name: "agent", flags: []string{"--agent", "codex"}, explicit: true},
		{name: "scope", flags: []string{"--scope", "user"}, explicit: true},
		{name: "directory", flags: []string{"--dir", "custom"}, explicit: true},
		{name: "agent equals", flags: []string{"--agent=codex"}, explicit: true},
		{name: "scope equals", flags: []string{"--scope=user"}, explicit: true},
		{name: "directory equals", flags: []string{"--dir=custom"}, explicit: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			desired := []Desired{{Repository: "github/repo", Skill: "skill", Name: "skill", Flags: tc.flags}}
			var calls [][]string
			command := func(_ context.Context, _ string, args ...string) ([]byte, error) {
				if len(args) == 3 && args[1] == "list" {
					return nil, nil
				}
				calls = append(calls, append([]string(nil), args...))
				return nil, nil
			}
			if err := Sync(context.Background(), home, filepath.Join(home, "state"), filepath.Join(home, "lock"), desired, command, &strings.Builder{}); err != nil {
				t.Fatal(err)
			}
			base := append([]string{"skill", "install", "github/repo", "skill"}, tc.flags...)
			want := [][]string{base}
			if !tc.explicit {
				want = nil
				for _, agent := range []string{"codex", "antigravity"} {
					want = append(want, append(append([]string(nil), base...), "--agent", agent, "--scope", "user"))
				}
			}
			if !reflect.DeepEqual(calls, want) {
				t.Fatalf("calls = %q, want %q", calls, want)
			}
		})
	}
}

func TestSyncPreservesExplicitInstallFlags(t *testing.T) {
	home := t.TempDir()
	desired := []Desired{{Repository: "github/repo", Skill: "sub/skill@v2", Name: "skill", Flags: []string{"--agent", "codex", "--scope", "user"}}}
	var got []string
	command := func(_ context.Context, _ string, args ...string) ([]byte, error) {
		if len(args) == 3 && args[0] == "skill" && args[1] == "list" && args[2] == "--help" {
			return nil, nil
		}
		got = append([]string(nil), args...)
		return nil, nil
	}
	if err := Sync(context.Background(), home, filepath.Join(home, "state"), filepath.Join(home, "lock"), desired, command, &strings.Builder{}); err != nil {
		t.Fatal(err)
	}
	want := []string{"skill", "install", "github/repo", "sub/skill@v2", "--agent", "codex", "--scope", "user"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("args = %q", got)
	}
}

func TestSyncRemovesOnlyPreviouslyOwnedUserSkill(t *testing.T) {
	home := t.TempDir()
	ownedPath := filepath.Join(home, ".agents", "skills", "owned")
	handPath := filepath.Join(home, ".agents", "skills", "manual")
	for _, path := range []string{ownedPath, handPath} {
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	state := filepath.Join(home, "state", "installed.tsv")
	if err := os.MkdirAll(filepath.Dir(state), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(state, []byte("github/repo\towned\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	lock := filepath.Join(home, ".agents", ".skill-lock.json")
	if err := os.MkdirAll(filepath.Dir(lock), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lock, []byte(`{"skills":{"owned":{"source":"github/repo"},"manual":{"source":"other/repo"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	listing, err := json.Marshal([]Installed{
		{SkillName: "owned", SourceURL: "https://github.com/github/repo", Path: ownedPath},
		{SkillName: "manual", SourceURL: "https://github.com/other/repo", Path: handPath},
	})
	if err != nil {
		t.Fatal(err)
	}
	command := func(_ context.Context, _ string, args ...string) ([]byte, error) {
		if len(args) == 3 && args[0] == "skill" && args[1] == "list" && args[2] == "--help" {
			return nil, nil
		}
		if len(args) > 2 && args[0] == "skill" && args[1] == "list" {
			return listing, nil
		}
		return nil, nil
	}
	if err := Sync(context.Background(), home, state, lock, nil, command, &strings.Builder{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ownedPath); !os.IsNotExist(err) {
		t.Fatalf("owned skill remains, stat error %v", err)
	}
	if _, err := os.Stat(handPath); err != nil {
		t.Fatalf("manual skill removed: %v", err)
	}
	updated, err := os.ReadFile(lock)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(updated), `"owned"`) || !strings.Contains(string(updated), `"manual"`) {
		t.Fatalf("lock entries = %s", updated)
	}
}

func TestSyncRetainsOwnershipWhenInstalledListFails(t *testing.T) {
	home := t.TempDir()
	state := filepath.Join(home, "state", "installed.tsv")
	if err := os.MkdirAll(filepath.Dir(state), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(state, []byte("github/repo\tremoved\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	command := func(_ context.Context, _ string, args ...string) ([]byte, error) {
		if len(args) == 3 && args[0] == "skill" && args[1] == "list" && args[2] == "--help" {
			return nil, nil
		}
		return nil, os.ErrPermission
	}
	if err := Sync(context.Background(), home, state, filepath.Join(home, "lock"), nil, command, &strings.Builder{}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(state)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "github/repo\tremoved\n" {
		t.Fatalf("ownership state = %q", data)
	}
}
