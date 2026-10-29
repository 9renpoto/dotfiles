#!/bin/sh
set -eu
ROOT="$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT"
exec mise exec -- go run ./cmd/dotfiles macos-runner "$@"
