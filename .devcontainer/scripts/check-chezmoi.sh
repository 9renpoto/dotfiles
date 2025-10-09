#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/../.." && pwd)

TEMP_HOME=$(mktemp -d)
trap 'rm -rf "$TEMP_HOME"' EXIT

export HOME="$TEMP_HOME"
export XDG_CONFIG_HOME="$TEMP_HOME/.config"
export XDG_CACHE_HOME="$TEMP_HOME/.cache"
export XDG_DATA_HOME="$TEMP_HOME/.local/share"

mkdir -p "$XDG_CONFIG_HOME" "$XDG_CACHE_HOME" "$XDG_DATA_HOME" "$XDG_DATA_HOME/chezmoi" "$TEMP_HOME/bin"

export CHEZMOI_SOURCE_DIR="$REPO_ROOT"
export CHEZMOI_DEST_DIR="$TEMP_HOME"
export PATH="$TEMP_HOME/bin:$PATH"
export HOMEBREW_NO_AUTO_UPDATE=1

# Initialize Homebrew environment if available so we can install dependencies.
if [ -x /home/linuxbrew/.linuxbrew/bin/brew ]; then
  eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
fi

# Ensure chezmoi is available for template validation.
if ! command -v chezmoi >/dev/null 2>&1; then
  if command -v brew >/dev/null 2>&1; then
    echo "Installing chezmoi via Homebrew for checks..."
    if ! brew install chezmoi; then
      echo "Homebrew install failed; falling back to script installer." >&2
    fi
  fi
fi

if ! command -v chezmoi >/dev/null 2>&1; then
  if command -v curl >/dev/null 2>&1; then
    echo "Installing chezmoi to temporary bin via upstream installer..."
    curl -fsSL https://chezmoi.io/get | sh -s -- -b "$TEMP_HOME/bin"
  else
    echo "curl is required to install chezmoi but is not available." >&2
    exit 1
  fi
fi

if ! command -v chezmoi >/dev/null 2>&1; then
  echo "chezmoi is still unavailable after installation attempts." >&2
  exit 1
fi

echo "Running chezmoi doctor..."
chezmoi doctor --source="$CHEZMOI_SOURCE_DIR" --destination="$CHEZMOI_DEST_DIR"

echo "Applying dotfiles into isolated destination..."
chezmoi apply --exclude=externals --verbose

echo "Verifying rendered state matches the source..."
chezmoi verify --exclude=externals --verbose
