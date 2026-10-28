#!/bin/sh
set -eu
cd "$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"
command -v mise >/dev/null 2>&1 || {
  echo 'Install mise 2026.10.19+ first: https://mise.jdx.dev/installing-mise.html' >&2
  exit 1
}
mise trust
exec mise exec -- go run ./cmd/dotfiles bootstrap "$@"
