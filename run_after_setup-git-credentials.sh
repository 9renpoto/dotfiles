#!/bin/sh

# Configure GitHub authentication for git through the GitHub CLI.
# This runs after every `chezmoi apply` on purpose: applying rewrites ~/.gitconfig
# from dot_gitconfig.tmpl, which is also where `gh auth setup-git` stores its
# credential helper entries. `gh auth setup-git` is idempotent.

set -e

if ! command -v gh >/dev/null 2>&1; then
    echo "gh not found; skipping 'gh auth setup-git'."
    exit 0
fi

if ! gh auth status >/dev/null 2>&1; then
    echo "gh is not authenticated; run 'gh auth login' and then 'chezmoi apply' again."
    exit 0
fi

echo "Configuring the git credential helper with 'gh auth setup-git'..."
gh auth setup-git
