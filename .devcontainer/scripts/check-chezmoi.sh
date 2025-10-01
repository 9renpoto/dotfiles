#!/bin/bash
set -euo pipefail

# chezmoi apply を実行し、差分があればエラーを出力して終了する
if ! chezmoi apply --dry-run --verbose; then
  echo "chezmoi apply found differences."
  exit 1
fi

echo "chezmoi configuration is up to date."
