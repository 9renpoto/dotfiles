package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/9renpoto/dotfiles/internal/skills"
	dotwindows "github.com/9renpoto/dotfiles/internal/windows"
)

type commandRunner func(context.Context, string, ...string) error

func main() {
	root, err := os.Getwd()
	if err == nil {
		err = run(context.Background(), os.Args[1:], root, os.Stdout, os.Stderr, execCommand)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execCommand(ctx context.Context, name string, args ...string) error {
	command := exec.CommandContext(ctx, name, args...)
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	return command.Run()
}

func run(ctx context.Context, args []string, root string, stdout, stderr io.Writer, command commandRunner) error {
	if len(args) > 0 && args[0] == "bootstrap" {
		return command(ctx, "mise", append([]string{"bootstrap"}, args[1:]...)...)
	}
	if len(args) == 1 && args[0] == "git-credentials" {
		return gitCredentials(ctx, stdout, command, exec.LookPath)
	}
	if len(args) == 1 && args[0] == "sync-skills" {
		return syncSkills(ctx, root, stdout)
	}
	if len(args) == 1 && args[0] == "integrate" {
		if err := gitCredentials(ctx, stdout, command, exec.LookPath); err != nil {
			return err
		}
		return syncSkills(ctx, root, stdout)
	}
	if len(args) == 3 && args[0] == "windows-profiles" && args[1] == "--documents" {
		return installWindowsProfiles(args[2], stdout)
	}
	if len(args) == 3 && args[0] == "windows-integrate" && args[1] == "--documents" {
		if err := installWindowsProfiles(args[2], stdout); err != nil {
			return err
		}
		return command(ctx, "git", "config", "--global", "--replace-all", "include.path", "~/.config/git/config")
	}
	if len(args) == 3 && args[0] == "windows-wsl-ghq" && args[1] == "--root" {
		return syncWslGhqRoot(ctx, args[2], stdout)
	}
	if len(args) > 0 && args[0] == "windows-init" {
		return runWindowsInitialize(ctx, root, args[1:], stdout)
	}
	if len(args) != 1 || args[0] != "packages" {
		return fmt.Errorf("usage: dotfiles {bootstrap [args...]|integrate|git-credentials|sync-skills|windows-integrate|windows-profiles --documents PATH|windows-wsl-ghq --root PATH|packages}")
	}
	if runtime.GOOS == "windows" {
		manifest := filepath.Join(root, "winget-packages.json")
		if _, err := os.Stat(manifest); err != nil {
			return fmt.Errorf("winget manifest %s: %w", manifest, err)
		}
		if err := command(ctx, "winget.exe", dotwindows.WingetImportArgs(manifest)...); err != nil {
			return fmt.Errorf("winget import failed: %w", err)
		}
		return nil
	}
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return nil
	}
	manifest := filepath.Join(root, "Brewfile")
	if _, err := os.Stat(manifest); err != nil {
		return fmt.Errorf("Homebrew manifest %s: %w", manifest, err)
	}
	fmt.Fprintln(stdout, "Installing Homebrew packages from Brewfile.")
	if err := command(ctx, "brew", "bundle", "--file="+manifest); err != nil {
		fmt.Fprintln(stderr, "Homebrew package installation failed.")
		return err
	}
	return nil
}

func runWindowsInitialize(ctx context.Context, root string, args []string, output io.Writer) error {
	options := dotwindows.InitializeOptions{Checkout: root}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--skip-packages":
			options.SkipPackages = true
		case "--skip-apply":
			options.SkipApply = true
		case "--ghq-root":
			i++
			if i == len(args) {
				return fmt.Errorf("--ghq-root requires a path")
			}
			options.GhqRoot = args[i]
		default:
			return fmt.Errorf("unknown windows-init argument %q", args[i])
		}
	}
	run := func(ctx context.Context, name string, commandArgs ...string) ([]byte, error) {
		cmd := exec.CommandContext(ctx, name, commandArgs...)
		data, err := cmd.CombinedOutput()
		if len(data) > 0 {
			secret := os.Getenv("WAKATIME_API_KEY")
			if secret != "" {
				data = []byte(strings.ReplaceAll(string(data), secret, "[REDACTED]"))
			}
			_, _ = output.Write(data)
		}
		return data, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return dotwindows.RunInitialize(ctx, options, home, func(name string) error { _, err := exec.LookPath(name); return err }, run, func(ctx context.Context) error { return refreshWindowsPath(ctx, run) }, output)
}

func refreshWindowsPath(ctx context.Context, run dotwindows.Runner) error {
	const script = `$machine=[Environment]::GetEnvironmentVariable('Path','Machine'); $user=[Environment]::GetEnvironmentVariable('Path','User'); [pscustomobject]@{Machine=$machine;User=$user}|ConvertTo-Json -Compress`
	data, err := run(ctx, "powershell.exe", "-NoProfile", "-Command", script)
	if err != nil {
		return fmt.Errorf("could not refresh process PATH: %w", err)
	}
	var paths struct {
		Machine string
		User    string
	}
	if err := json.Unmarshal(data, &paths); err != nil {
		return fmt.Errorf("could not read Windows PATH: %w", err)
	}
	merged := dotwindows.MergePathEntries(paths.Machine, paths.User, os.Getenv("PATH"), string(os.PathListSeparator))
	return os.Setenv("PATH", strings.Join(merged, string(os.PathListSeparator)))
}

func syncWslGhqRoot(ctx context.Context, root string, output io.Writer) error {
	return dotwindows.SyncWslGhqRoot(ctx, root, func(name string) error { _, err := exec.LookPath(name); return err }, func(ctx context.Context, name string, args ...string) ([]byte, error) {
		cmd := exec.CommandContext(ctx, name, args...)
		return cmd.CombinedOutput()
	}, output)
}

func installWindowsProfiles(documents string, output io.Writer) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	preserved, err := dotwindows.InstallProfiles(documents, home)
	if err != nil {
		return err
	}
	for _, path := range preserved {
		fmt.Fprintf(output, "Preserving customized PowerShell profile: %s\n", path)
	}
	fmt.Fprintf(output, "PowerShell profiles load %s\n", filepath.Join(home, ".config", "powershell", "profile.ps1"))
	return nil
}

func syncSkills(ctx context.Context, root string, output io.Writer) error {
	file, err := os.Open(filepath.Join(root, "Skillsfile"))
	if err != nil {
		return err
	}
	defer file.Close()
	desired, err := skills.ParseDesired(file)
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	state := os.Getenv("XDG_STATE_HOME")
	if state == "" {
		state = filepath.Join(home, ".local", "state")
	}
	return skills.Sync(ctx, home, filepath.Join(state, "mise-dotfiles", "installed-skills.tsv"), filepath.Join(home, ".agents", ".skill-lock.json"), desired, func(ctx context.Context, name string, args ...string) ([]byte, error) {
		cmd := exec.CommandContext(ctx, name, args...)
		return cmd.CombinedOutput()
	}, output)
}

func gitCredentials(ctx context.Context, stdout io.Writer, command commandRunner, lookPath func(string) (string, error)) error {
	if _, err := lookPath("gh"); err != nil {
		fmt.Fprintln(stdout, "gh not found; skipping 'gh auth setup-git'.")
		return nil
	}
	if err := command(ctx, "gh", "auth", "status"); err != nil {
		fmt.Fprintln(stdout, "gh is not authenticated; run 'gh auth login' and then 'mise run apply' again.")
		return nil
	}
	fmt.Fprintln(stdout, "Configuring the git credential helper with 'gh auth setup-git'...")
	return command(ctx, "gh", "auth", "setup-git")
}
