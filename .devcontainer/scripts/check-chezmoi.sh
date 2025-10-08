#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)

# Initialize Homebrew environment if available so we can install dependencies.
if [ -x /home/linuxbrew/.linuxbrew/bin/brew ]; then
  eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
fi

# Ensure chezmoi is available for template validation.
if ! command -v chezmoi >/dev/null 2>&1; then
  if command -v brew >/dev/null 2>&1; then
    echo "Installing chezmoi via Homebrew for checks..."
    brew install chezmoi
  else
    echo "chezmoi is not installed and Homebrew is unavailable." >&2
    exit 1
  fi
fi

TEMP_HOME=$(mktemp -d)
trap 'rm -rf "$TEMP_HOME"' EXIT

export HOME="$TEMP_HOME"
export XDG_CONFIG_HOME="$TEMP_HOME/.config"
export XDG_CACHE_HOME="$TEMP_HOME/.cache"
export XDG_DATA_HOME="$TEMP_HOME/.local/share"

mkdir -p "$XDG_CONFIG_HOME" "$XDG_CACHE_HOME" "$XDG_DATA_HOME"

export CHEZMOI_SOURCE_DIR="$REPO_ROOT"
export CHEZMOI_DEST_DIR="$TEMP_HOME"

echo "Running chezmoi doctor..."
chezmoi doctor

echo "Applying dotfiles into isolated destination..."
chezmoi apply --exclude=externals --verbose

echo "Verifying rendered state matches the source..."
chezmoi verify --exclude=externals --verbose
