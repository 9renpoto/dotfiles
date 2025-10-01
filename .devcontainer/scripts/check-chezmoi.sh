#!/usr/bin/env bash
set -euo pipefail

if ! command -v chezmoi >/dev/null 2>&1; then
  echo "chezmoi is not installed; skipping check." >&2
  exit 0
fi

SOURCE_DIR=${CHEZMOI_SOURCE_DIR:-$(pwd)}
if [ ! -d "$SOURCE_DIR" ]; then
  echo "Specified source directory '$SOURCE_DIR' does not exist; skipping check." >&2
  exit 0
fi

if [ ! -f "$SOURCE_DIR/dot_chezmoi.toml.tmpl" ] && [ ! -f "$SOURCE_DIR/.chezmoiignore" ]; then
  echo "chezmoi source files not present in $SOURCE_DIR; skipping check." >&2
  exit 0
fi

scratch_root=$(mktemp -d "${TMPDIR:-/tmp}/chezmoi-test-XXXXXX")
cleanup() {
  rm -rf "$scratch_root"
}
trap cleanup EXIT

home_dir="$scratch_root/home"
mkdir -p "$home_dir"

export HOME="$home_dir"
export XDG_CONFIG_HOME="$home_dir/.config"
export XDG_DATA_HOME="$home_dir/.local/share"
export XDG_CACHE_HOME="$home_dir/.cache"

apply() {
  chezmoi --source="$SOURCE_DIR" --destination="$home_dir" apply --keep-going --progress=false
}

diff() {
  chezmoi --source="$SOURCE_DIR" --destination="$home_dir" diff
}

echo "Applying chezmoi source from $SOURCE_DIR into scratch HOME $home_dir"
apply

diff_output=$(mktemp "$scratch_root/diff.XXXXXX")
if diff >"$diff_output"; then
  if [ -s "$diff_output" ]; then
    cat "$diff_output"
    echo "chezmoi diff detected changes after apply." >&2
    exit 1
  fi
else
  status=$?
  cat "$diff_output"
  echo "chezmoi diff failed with status $status." >&2
  exit $status
fi

rm -f "$diff_output"

echo "chezmoi apply completed without residual differences."
