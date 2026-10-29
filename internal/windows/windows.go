package windows

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Runner func(context.Context, string, ...string) ([]byte, error)

func SyncWslGhqRoot(ctx context.Context, root string, lookPath func(string) error, run Runner, output io.Writer) error {
	if err := lookPath("wsl.exe"); err != nil {
		fmt.Fprintln(output, "WSL was not found; skipped WSL ghq root configuration.")
		return nil
	}
	converted, err := run(ctx, "wsl.exe", "--exec", "wslpath", "-a", "-u", root)
	if err != nil || strings.TrimSpace(string(converted)) == "" {
		fmt.Fprintln(output, "Could not convert the Windows ghq root to a WSL path.")
		return nil
	}
	wslRoot := strings.TrimSpace(strings.SplitN(string(converted), "\n", 2)[0])
	if _, err := run(ctx, "wsl.exe", "--exec", "mkdir", "-p", wslRoot); err != nil {
		fmt.Fprintf(output, "Could not create the WSL ghq root at %s.\n", wslRoot)
		return nil
	}
	if _, err := run(ctx, "wsl.exe", "--exec", "git", "config", "--global", "ghq.root", wslRoot); err != nil {
		fmt.Fprintln(output, "Git was not available in the default WSL distribution; skipped its ghq configuration.")
		return nil
	}
	fmt.Fprintf(output, "WSL ghq root: %s\n", wslRoot)
	return nil
}

type InitializeOptions struct {
	Checkout     string
	GhqRoot      string
	SkipPackages bool
	SkipApply    bool
}

func RunInitialize(ctx context.Context, options InitializeOptions, home string, lookPath func(string) error, run Runner, refreshPath func(context.Context) error, output io.Writer) error {
	if err := lookPath("mise"); err != nil {
		if _, err := run(ctx, "winget.exe", "install", "--id", "jdx.mise", "--exact", "--accept-package-agreements", "--accept-source-agreements"); err != nil {
			return fmt.Errorf("mise installation failed: %w", err)
		}
		if err := refreshPath(ctx); err != nil {
			return err
		}
	}
	if _, err := run(ctx, "mise", "trust"); err != nil {
		return fmt.Errorf("mise trust failed: %w", err)
	}
	ghqRoot := options.GhqRoot
	if ghqRoot == "" {
		ghqRoot = filepath.Join(home, "src")
	}
	ghqRoot, err := filepath.Abs(ghqRoot)
	if err != nil {
		return err
	}
	ghqRoot = strings.ReplaceAll(ghqRoot, `\`, "/")
	if err := os.MkdirAll(ghqRoot, 0o755); err != nil {
		return err
	}
	localConfig := filepath.Join(options.Checkout, "mise.local.toml")
	if _, err := os.Stat(localConfig); os.IsNotExist(err) {
		if err := os.WriteFile(localConfig, nil, 0o600); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if _, err := run(ctx, "mise", "config", "set", "--file", localConfig, "--type", "string", "vars.ghq_root", ghqRoot); err != nil {
		return fmt.Errorf("failed to persist ghq root: %w", err)
	}
	if !options.SkipPackages {
		manifest := filepath.Join(options.Checkout, "winget-packages.json")
		if _, err := os.Stat(manifest); err != nil {
			return fmt.Errorf("winget manifest %s: %w", manifest, err)
		}
		if _, err := run(ctx, "winget.exe", WingetImportArgs(manifest)...); err != nil {
			return fmt.Errorf("package installation failed: %w", err)
		}
		if err := refreshPath(ctx); err != nil {
			return err
		}
	}
	if !options.SkipApply {
		if _, err := run(ctx, "mise", "bootstrap", "--skip", "packages"); err != nil {
			return fmt.Errorf("mise bootstrap failed: %w", err)
		}
		if err := SyncWslGhqRoot(ctx, ghqRoot, lookPath, run, output); err != nil {
			return err
		}
	}
	return nil
}

func WingetImportArgs(manifest string) []string {
	return []string{"import", "--import-file", manifest, "--ignore-unavailable", "--ignore-versions", "--accept-package-agreements", "--accept-source-agreements"}
}

func MergePathEntries(machine, user, process, separator string) []string {
	seen := make(map[string]struct{})
	var merged []string
	for _, value := range []string{machine, user, process} {
		for _, entry := range strings.Split(value, separator) {
			entry = strings.TrimSpace(entry)
			key := strings.ToLower(entry)
			if entry == "" {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			merged = append(merged, entry)
		}
	}
	return merged
}

func ProfileLoader(sharedProfilePath string) string {
	sharedProfilePath = strings.ReplaceAll(sharedProfilePath, "'", "''")
	return "$sharedProfilePath = '" + sharedProfilePath + "'\n" +
		"if (Test-Path -LiteralPath $sharedProfilePath) {\n" +
		"  . $sharedProfilePath\n" +
		"}"
}

func MergeProfile(existing []byte, exists bool, loader string) ([]byte, bool) {
	content := strings.TrimSpace(strings.TrimPrefix(string(existing), "\ufeff"))
	if !exists || content == "" || content == strings.TrimSpace(loader) {
		return []byte("\ufeff" + loader + "\n"), true
	}
	return existing, false
}

func InstallProfiles(documents, home string) ([]string, error) {
	shared := filepath.Join(home, ".config", "powershell", "profile.ps1")
	loader := ProfileLoader(shared)
	var preserved []string
	for _, name := range []string{"WindowsPowerShell", "PowerShell"} {
		dir := filepath.Join(documents, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return preserved, err
		}
		path := filepath.Join(dir, "profile.ps1")
		content, err := os.ReadFile(path)
		exists := err == nil
		if err != nil && !os.IsNotExist(err) {
			return preserved, err
		}
		updated, write := MergeProfile(content, exists, loader)
		if write {
			if err := os.WriteFile(path, updated, 0o644); err != nil {
				return preserved, err
			}
		} else {
			preserved = append(preserved, path)
		}
	}
	return preserved, nil
}
