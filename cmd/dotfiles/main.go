package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/9renpoto/dotfiles/internal/skills"
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
	if len(args) != 1 || args[0] != "packages" {
		return fmt.Errorf("usage: dotfiles {bootstrap [args...]|integrate|git-credentials|sync-skills|packages}")
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
