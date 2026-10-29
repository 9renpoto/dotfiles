package skills

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Command func(context.Context, string, ...string) ([]byte, error)

func ReadInstalled(ctx context.Context, command Command) ([]Installed, error) {
	data, err := command(ctx, "gh", "skill", "list", "--scope", "user", "--json", "skillName,sourceURL,path")
	if err != nil {
		return nil, err
	}
	var items []Installed
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	for i := range items {
		items[i].Scope = "user"
	}
	return items, nil
}

func Sync(ctx context.Context, home, statePath, lockPath string, desired []Desired, command Command, output io.Writer) error {
	if _, err := command(ctx, "gh", "skill", "list", "--help"); err != nil {
		fmt.Fprintln(output, "'gh skill' is unavailable; skipping skill sync.")
		return nil
	}
	owned := readOwned(statePath)
	ownedSet := make(map[string]bool, len(owned))
	for _, key := range owned {
		ownedSet[key] = true
	}
	if len(owned) > 0 {
		if installed, err := ReadInstalled(ctx, command); err == nil {
			for _, item := range RemovalCandidates(home, owned, desired, installed) {
				if !OwnedSkillPath(home, item.Path) {
					continue
				}
				fmt.Fprintf(output, "Removing skill dropped from Skillsfile: %s (%s)\n", item.SkillName, item.Path)
				if err := os.RemoveAll(item.Path); err != nil {
					fmt.Fprintf(output, "  could not remove skill: %v\n", err)
					continue
				}
				delete(ownedSet, OwnedKey(mustRepository(item.SourceURL), pathSkillName(item.SkillName)))
				if data, err := os.ReadFile(lockPath); err == nil {
					if updated, err := RemoveLockEntry(data, mustRepository(item.SourceURL), pathSkillName(item.SkillName)); err == nil {
						_ = writeAtomic(lockPath, updated, 0o600)
					}
				}
			}
		}
	}
	for _, entry := range desired {
		args := []string{"skill", "install", entry.Repository, entry.Skill}
		args = append(args, entry.Flags...)
		key := OwnedKey(entry.Repository, entry.Name)
		if hasInstallDestination(entry.Flags) {
			if _, err := command(ctx, "gh", args...); err == nil {
				ownedSet[key] = true
			} else {
				fmt.Fprintf(output, "Could not install skill %s %s: %v\n", entry.Repository, entry.Skill, err)
			}
			continue
		}
		success := false
		for _, agent := range []string{"codex", "antigravity"} {
			agentArgs := append(append([]string(nil), args...), "--agent", agent, "--scope", "user")
			if _, err := command(ctx, "gh", agentArgs...); err != nil {
				fmt.Fprintf(output, "Could not install skill for %s: %s %s: %v\n", agent, entry.Repository, entry.Skill, err)
			} else {
				success = true
			}
		}
		if success {
			ownedSet[key] = true
		}
	}
	applied := make([]string, 0, len(ownedSet))
	for key := range ownedSet {
		applied = append(applied, key)
	}
	sort.Strings(applied)
	return writeAtomic(statePath, []byte(strings.Join(applied, "\n")+func() string {
		if len(applied) > 0 {
			return "\n"
		}
		return ""
	}()), 0o600)
}

func hasInstallDestination(flags []string) bool {
	for _, flag := range flags {
		name, _, _ := strings.Cut(flag, "=")
		switch name {
		case "--agent", "--scope", "--dir":
			return true
		}
	}
	return false
}

func readOwned(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var result []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "\t") {
			result = append(result, line)
		}
	}
	return result
}

func mustRepository(source string) string { repo, _ := SourceRepository(source); return repo }

func writeAtomic(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".dotfiles-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	if err := f.Chmod(mode); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}

type Desired struct {
	Repository string
	Skill      string
	Name       string
	Flags      []string
}

func ParseDesired(r io.Reader) ([]Desired, error) {
	var entries []Desired
	scanner := bufio.NewScanner(r)
	for line := 1; scanner.Scan(); line++ {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		if len(fields) < 2 || strings.Contains(fields[0], "/") == false {
			return nil, fmt.Errorf("Skillsfile line %d: expected <owner/repo> <skill>", line)
		}
		skill := fields[1]
		name := skill
		if i := strings.IndexByte(name, '@'); i >= 0 {
			name = name[:i]
		}
		if i := strings.LastIndexByte(name, '/'); i >= 0 {
			name = name[i+1:]
		}
		if name == "" {
			return nil, fmt.Errorf("Skillsfile line %d: empty skill name", line)
		}
		entries = append(entries, Desired{Repository: fields[0], Skill: skill, Name: name, Flags: append([]string(nil), fields[2:]...)})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func SourceRepository(source string) (string, error) {
	value := strings.TrimSuffix(strings.TrimSpace(source), ".git")
	value = strings.TrimPrefix(value, "https://")
	value = strings.TrimPrefix(value, "http://")
	value = strings.TrimPrefix(value, "ssh://")
	value = strings.TrimPrefix(value, "git@")
	value = strings.Replace(value, ":", "/", 1)
	parts := strings.Split(value, "/")
	if len(parts) > 2 && strings.Contains(parts[0], ".") {
		parts = parts[1:]
	}
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("invalid skill source URL %q", source)
	}
	return parts[0] + "/" + parts[1], nil
}

func OwnedSkillPath(home, path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == ".." {
			return false
		}
	}
	home, err := filepath.Abs(home)
	if err != nil {
		return false
	}
	home, err = filepath.EvalSymlinks(home)
	if err != nil {
		return false
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return false
	}
	resolved, resolveErr := filepath.EvalSymlinks(path)
	if resolveErr == nil {
		path = resolved
	} else {
		parent, parentErr := filepath.EvalSymlinks(filepath.Dir(path))
		if parentErr != nil {
			return false
		}
		path = filepath.Join(parent, filepath.Base(path))
	}
	if !strings.HasPrefix(path, home+string(filepath.Separator)) {
		return false
	}
	rel, err := filepath.Rel(home, path)
	if err != nil {
		return false
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	for i, part := range parts {
		if part == "skills" && i+1 < len(parts) && parts[i+1] != "" && parts[i+1] != "." && parts[i+1] != ".." {
			return true
		}
	}
	return false
}

type Installed struct {
	SkillName string `json:"skillName"`
	SourceURL string `json:"sourceURL"`
	Path      string `json:"path"`
	Scope     string `json:"scope"`
}

func OwnedKey(repository, name string) string { return repository + "\t" + name }

func RemovalCandidates(home string, owned []string, desired []Desired, installed []Installed) []Installed {
	wanted := make(map[string]bool, len(desired))
	for _, entry := range desired {
		wanted[OwnedKey(entry.Repository, entry.Name)] = true
	}
	var remove []Installed
	for _, key := range owned {
		if wanted[key] {
			continue
		}
		parts := strings.SplitN(key, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		for _, item := range installed {
			if item.Scope != "user" || pathSkillName(item.SkillName) != parts[1] || !OwnedSkillPath(home, item.Path) {
				continue
			}
			repo, err := SourceRepository(item.SourceURL)
			if err == nil && repo == parts[0] {
				remove = append(remove, item)
			}
		}
	}
	return remove
}

func pathSkillName(value string) string {
	if i := strings.LastIndexByte(value, '/'); i >= 0 {
		return value[i+1:]
	}
	return value
}

func RemoveLockEntry(data []byte, repository, name string) ([]byte, error) {
	var document map[string]json.RawMessage
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, err
	}
	var skills map[string]json.RawMessage
	if raw := document["skills"]; len(raw) > 0 {
		if err := json.Unmarshal(raw, &skills); err != nil {
			return nil, err
		}
	}
	if skill, ok := skills[name]; ok {
		var fields map[string]any
		if err := json.Unmarshal(skill, &fields); err != nil {
			return nil, err
		}
		if source, _ := fields["source"].(string); source == repository {
			delete(skills, name)
		}
	}
	updated, err := json.Marshal(skills)
	if err != nil {
		return nil, err
	}
	document["skills"] = updated
	return json.MarshalIndent(document, "", "  ")
}
