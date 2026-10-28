#!/bin/sh
set -eu
container=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" --format '{{.ID}}' | head -n 1)
[ -n "$container" ] || { echo 'No running devcontainer found' >&2; exit 1; }
docker exec -u vscode -w /workspaces/dotfiles "$container" sh -c 'export PATH="$HOME/.local/bin:$PATH"; exec sh ./initialize.sh --only dotfiles,task'
