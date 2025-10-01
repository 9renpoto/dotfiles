#!/usr/bin/env bash
set -euo pipefail

if [ ! -f Brewfile ]; then
  echo "Brewfile not found; skipping brew bundle."
  exit 0
fi

BREW_BIN="/home/linuxbrew/.linuxbrew/bin/brew"
if [ ! -x "$BREW_BIN" ]; then
  echo "Homebrew not installed at $BREW_BIN; skipping brew bundle" >&2
  exit 0
fi

eval "$("$BREW_BIN" shellenv)"
export HOMEBREW_NO_AUTO_UPDATE=1

system_name=$(uname -s)
if [ "$system_name" = "Linux" ]; then
  tmp_brewfile=$(mktemp)
  trap 'rm -f "$tmp_brewfile"' EXIT
  awk '/^(tap|brew)/ { print }' Brewfile > "$tmp_brewfile"
  echo "Running brew bundle for taps and brews (casks skipped on $system_name)..."
  if ! brew bundle --file="$tmp_brewfile"; then
    echo "brew bundle reported errors; continuing startup" >&2
  fi
else
  echo "Running brew bundle using Brewfile for $system_name..."
  if ! brew bundle --file=Brewfile; then
    echo "brew bundle reported errors; continuing startup" >&2
  fi
fi
