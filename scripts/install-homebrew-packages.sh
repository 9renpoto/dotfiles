#!/bin/sh
set -eu
cd "$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)"
exec mise exec -- go run ./cmd/dotfiles packages
