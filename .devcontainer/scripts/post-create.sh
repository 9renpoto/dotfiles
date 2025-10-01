#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)

cd "$REPO_ROOT"

bash "$SCRIPT_DIR/run-brew-bundle.sh"

if command -v mise >/dev/null 2>&1; then
  echo "Running mise install to provision runtimes..."
  if ! mise install; then
    echo "mise install failed." >&2
    exit 1
  fi
else
  echo "mise is not available; skipping runtime provisioning." >&2
fi

if command -v brew >/dev/null 2>&1; then
  system_name=$(uname -s)
  if [ "$system_name" = "Linux" ]; then
    tmp_brewfile=$(mktemp)
    trap 'rm -f "$tmp_brewfile"' EXIT
    awk '/^(tap|brew)/ { print }' Brewfile >"$tmp_brewfile"
    echo "Validating Brewfile taps and brews with brew bundle check..."
    if ! brew bundle check --file="$tmp_brewfile"; then
      echo "brew bundle check reported issues." >&2
      exit 1
    fi
    rm -f "$tmp_brewfile"
  else
    echo "Validating Brewfile with brew bundle check..."
    if ! brew bundle check --file=Brewfile; then
      echo "brew bundle check reported issues." >&2
      exit 1
    fi
  fi
else
  echo "Homebrew not available; skipping brew bundle check." >&2
fi

if command -v lefthook >/dev/null 2>&1; then
  if command -v docker >/dev/null 2>&1; then
    echo "Running lefthook pre-commit checks..."
    if ! lefthook run pre-commit --all-files; then
      echo "lefthook pre-commit failed." >&2
      exit 1
    fi
  else
    echo "docker command not available; skipping lefthook pre-commit." >&2
  fi
else
  echo "lefthook not installed; skipping pre-commit hook execution." >&2
fi
