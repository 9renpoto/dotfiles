#!/bin/sh
set -eu
case "$(uname -s)" in
  Darwin|Linux) ;;
  *) exit 0 ;;
esac
command -v brew >/dev/null 2>&1 || {
  echo 'Install Homebrew first: https://brew.sh' >&2
  exit 1
}
repository_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
exec brew bundle --file="$repository_root/Brewfile"
