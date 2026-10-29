package tests

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repository(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	return root
}
func TestIsolatedEnvKey(t *testing.T) {
	for _, key := range []string{"APPDATA", "AppData", "appdata"} {
		t.Run(key, func(t *testing.T) {
			if !isIsolatedEnvKey(key) {
				t.Fatalf("%q should be removed from the isolated environment", key)
			}
		})
	}
}
func TestIsolateHomePaths(t *testing.T) {
	config := []byte("\"~/.config/tool/config.toml\" = { source = \"home/config.toml\" }\n")
	got := isolateHomePaths(config, "C:\\Users\\runner\\Temp\\home")
	want := "\"C:/Users/runner/Temp/home/.config/tool/config.toml\" = { source = \"home/config.toml\" }\n"
	if string(got) != want {
		t.Fatalf("isolateHomePaths() = %q, want %q", got, want)
	}
}
func isolated(t *testing.T) []string {
	t.Helper()
	home := t.TempDir()
	var env []string
	for _, entry := range os.Environ() {
		key := strings.SplitN(entry, "=", 2)[0]
		if !isIsolatedEnvKey(key) {
			env = append(env, entry)
		}
	}
	return append(env, "HOME="+home, "USERPROFILE="+home, "XDG_CONFIG_HOME="+home+"/.config", "APPDATA="+home+"/AppData/Roaming", "MISE_CONFIG_DIR="+home+"/.config/mise", "MISE_DATA_DIR="+home+"/data", "MISE_CACHE_DIR="+home+"/cache", "MISE_STATE_DIR="+home+"/state", "MISE_SYSTEM_CONFIG_DIR="+home+"/system", "MISE_GLOBAL_CONFIG_FILE="+home+"/global.toml", "MISE_OFFLINE=1")
}
func isIsolatedEnvKey(key string) bool {
	key = strings.ToUpper(key)
	return strings.HasPrefix(key, "MISE_") || strings.HasPrefix(key, "__MISE_") || key == "HOME" || key == "USERPROFILE" || key == "APPDATA" || strings.HasPrefix(key, "XDG_")
}
func isolateHomePaths(config []byte, home string) []byte {
	home = strings.ReplaceAll(home, "\\", "/")
	return bytes.ReplaceAll(config, []byte("~"), []byte(home))
}
func isolatedMiseConfig(t *testing.T, source, destination, home string) {
	t.Helper()
	config := read(t, source)
	if runtime.GOOS == "windows" {
		config = isolateHomePaths(config, home)
	}
	if err := os.WriteFile(destination, config, 0600); err != nil {
		t.Fatal(err)
	}
}
func command(t *testing.T, dir string, env []string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
	return string(out)
}
func read(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func contains(t *testing.T, path, value string) {
	t.Helper()
	if !bytes.Contains(read(t, path), []byte(value)) {
		t.Fatalf("%s does not contain expected value", filepath.Base(path))
	}
}
func homeOf(env []string) string {
	for _, e := range env {
		if strings.HasPrefix(e, "HOME=") {
			return strings.TrimPrefix(e, "HOME=")
		}
	}
	panic("missing HOME")
}
func copyTree(t *testing.T, from, to string) {
	t.Helper()
	err := filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(from, path)
		target := filepath.Join(to, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, 0600)
	})
	if err != nil {
		t.Fatal(err)
	}
}
func TestDeployment(t *testing.T) {
	sourceRoot := repository(t)
	env := isolated(t)
	home := homeOf(env)
	root := filepath.Join(home, "checkout")
	if err := os.Mkdir(root, 0755); err != nil {
		t.Fatal(err)
	}
	copyTree(t, filepath.Join(sourceRoot, "home"), filepath.Join(root, "home"))
	for _, name := range []string{"mise.toml", "Brewfile", "winget-packages.json"} {
		t.Run("checkout_file/"+name, func(t *testing.T) {
			destination := filepath.Join(root, name)
			if name == "mise.toml" {
				isolatedMiseConfig(t, filepath.Join(sourceRoot, name), destination, home)
			} else if err := os.WriteFile(destination, read(t, filepath.Join(sourceRoot, name)), 0600); err != nil {
				t.Fatal(err)
			}
		})
	}
	env = append(env, "MISE_TRUSTED_CONFIG_PATHS="+root+string(os.PathListSeparator)+home, "WAKATIME_API_KEY=test-key", "GRAFANA_CLOUD_API_KEY=test-token '$HOME $(false)")
	applyOutput := command(t, root, env, "mise", "dot", "apply")
	biome := filepath.Join(home, ".config/biome/biome.json")
	if runtime.GOOS == "darwin" {
		biome = filepath.Join(home, "Library/Application Support/biome/biome.json")
	}
	if runtime.GOOS == "windows" {
		biome = filepath.Join(home, "AppData/Roaming/biome/config/biome.json")
	}
	if _, err := os.Stat(biome); err != nil {
		var deployed []string
		_ = filepath.WalkDir(home, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() {
				rel, err := filepath.Rel(home, path)
				if err != nil {
					return err
				}
				deployed = append(deployed, filepath.ToSlash(rel))
			}
			return nil
		})
		var biomeOutput []string
		for _, line := range strings.Split(applyOutput, "\n") {
			if strings.Contains(strings.ToLower(line), "biome") {
				biomeOutput = append(biomeOutput, strings.TrimSpace(line))
			}
		}
		t.Fatalf("Biome config was not deployed to %s; files under isolated home: %v; mise output: %v", filepath.ToSlash(strings.TrimPrefix(biome, home+string(os.PathSeparator))), deployed, biomeOutput)
	}
	var config struct {
		Formatter struct {
			IndentWidth int `json:"indentWidth"`
		} `json:"formatter"`
	}
	if err := json.Unmarshal(read(t, biome), &config); err != nil {
		t.Fatal(err)
	}
	if config.Formatter.IndentWidth != 2 {
		t.Fatal("unexpected indent width")
	}
	contains(t, filepath.Join(home, ".wakatime.cfg"), "test-key")
	contains(t, filepath.Join(home, ".config/alloy/config.env"), "test-token")
	if runtime.GOOS != "windows" {
		contains(t, filepath.Join(home, ".gitconfig"), "~/.config/git/config")
		info, err := os.Stat(filepath.Join(home, ".wakatime.cfg"))
		if err != nil || info.Mode().Perm() != 0600 {
			t.Fatal("secret permissions", err)
		}
		value := command(t, root, env, "sh", "-c", `. "$1"; printf %s "$GRAFANA_CLOUD_API_KEY"`, "sh", filepath.Join(home, ".config/alloy/config.env"))
		if value != "test-token '$HOME $(false)" {
			t.Fatal("secret shell round trip failed")
		}
	}
	before := map[string][]byte{}
	err := filepath.WalkDir(home, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == "state" || d.Name() == "cache") {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			before[path] = read(t, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	command(t, root, env, "mise", "dot", "apply")
	for path, want := range before {
		rel, err := filepath.Rel(home, path)
		if err != nil {
			t.Fatal(err)
		}
		t.Run("repeat_apply/"+filepath.ToSlash(rel), func(t *testing.T) {
			if !bytes.Equal(read(t, path), want) {
				t.Fatalf("repeat apply changed %s", filepath.Base(path))
			}
		})
	}
	command(t, root, env, "mise", "tasks", "ls")
	command(t, root, env, "mise", "bootstrap", "repos", "apply", "--dry-run")
	command(t, root, env, "mise", "bootstrap", "packages", "status")
}
func TestLocalValuesSurviveApply(t *testing.T) {
	root := repository(t)
	env := isolated(t)
	home := homeOf(env)
	checkout := filepath.Join(home, "checkout")
	if err := os.Mkdir(checkout, 0755); err != nil {
		t.Fatal(err)
	}
	copyTree(t, filepath.Join(root, "home"), filepath.Join(checkout, "home"))
	local := filepath.Join(checkout, "mise.local.toml")
	isolatedMiseConfig(t, filepath.Join(root, "mise.toml"), filepath.Join(checkout, "mise.toml"), home)
	if err := os.WriteFile(local, []byte("# Machine secrets remain local\n[vars]\nemail = \"keep@example.com\"\nwakatime_api_key = \"keep-key\"\nghq_root = \"C:/old\"\n[env]\nKEEP_SETTING = \"yes\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	env = append(env, "MISE_TRUSTED_CONFIG_PATHS="+checkout+string(os.PathListSeparator)+home)
	command(t, checkout, env, "mise", "config", "set", "--file", local, "--type", "string", "vars.ghq_root", "D:/source repos")
	for _, tc := range []struct{ name, value string }{
		{name: "email", value: "keep@example.com"},
		{name: "api_key", value: "keep-key"},
		{name: "environment", value: "KEEP_SETTING = \"yes\""},
		{name: "comment", value: "# Machine secrets remain local"},
	} {
		t.Run("preserves_"+tc.name, func(t *testing.T) {
			contains(t, local, tc.value)
		})
	}
	for _, step := range []string{"first_apply", "repeat_apply"} {
		t.Run(step, func(t *testing.T) {
			command(t, checkout, env, "mise", "dot", "apply")
			if runtime.GOOS != "windows" {
				contains(t, filepath.Join(home, ".config/git/config"), "D:/source repos")
			}
		})
	}
}
func TestBiomeConfig(t *testing.T) {
	if _, err := exec.LookPath("biome"); err != nil {
		t.Skip("biome is not installed")
	}
	env := isolated(t)
	home := homeOf(env)
	project := filepath.Join(home, "project")
	if err := os.Mkdir(project, 0755); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(home, ".config/biome")
	if runtime.GOOS == "darwin" {
		dir = filepath.Join(home, "Library/Application Support/biome")
	}
	if runtime.GOOS == "windows" {
		dir = filepath.Join(home, "AppData/Roaming/biome/config")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "biome.json"), read(t, filepath.Join(repository(t), "home/.config/biome/biome.json")), 0600); err != nil {
		t.Fatal(err)
	}
	sample := "const value = {\n  enabled: true,\n};\n"
	format := func() string {
		cmd := exec.Command("biome", "format", "--stdin-file-path=example.js")
		cmd.Dir = project
		cmd.Env = env
		cmd.Stdin = strings.NewReader(sample)
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		return string(out)
	}
	if format() != sample {
		t.Fatal("shared formatting differs")
	}
	if err := os.WriteFile(filepath.Join(project, "biome.json"), []byte(`{"formatter":{"indentStyle":"tab"}}`), 0600); err != nil {
		t.Fatal(err)
	}
	if format() != strings.ReplaceAll(sample, "  enabled", "\tenabled") {
		t.Fatal("project override ignored")
	}
}

func TestPlatformSources(t *testing.T) {
	root := repository(t)
	// Let mise validate and resolve every platform's declared source mapping.
	config := string(read(t, filepath.Join(root, "mise.toml")))
	for _, platform := range []string{"linux", "macos", "windows"} {
		t.Run("platform/"+platform, func(t *testing.T) {
			if !strings.Contains(config, `os = "`+platform+`"`) {
				t.Fatalf("missing platform %s", platform)
			}
		})
	}
	for _, line := range strings.Split(config, "\n") {
		marker := `source = "`
		index := strings.Index(line, marker)
		if index < 0 {
			continue
		}
		source := strings.SplitN(line[index+len(marker):], `"`, 2)[0]
		t.Run("source/"+filepath.ToSlash(source), func(t *testing.T) {
			info, err := os.Stat(filepath.Join(root, source))
			if err != nil || info.IsDir() {
				t.Fatalf("missing declared source %s", source)
			}
		})
	}
	t.Run("windows_preserves_main_git_configuration", func(t *testing.T) {
		gitconfig := strings.SplitN(strings.SplitN(config, "[dotfiles.\"gitconfig\"]", 2)[1], "[dotfiles.biome]", 2)[0]
		if strings.Contains(gitconfig, `os = "windows"`) {
			t.Fatal("Windows main Git configuration must be preserved")
		}
	})
}
