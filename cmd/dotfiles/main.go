package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	if len(args) != 1 || args[0] != "packages" {
		return fmt.Errorf("usage: dotfiles {bootstrap [args...]|packages}")
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
