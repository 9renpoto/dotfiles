#!/bin/sh
set -e

# Applies dotfiles to the running devcontainer using `chezmoi docker`.

# --- Helper Functions ---
info() {
    printf "\033[1;34m%s\033[0m\n" "$1"
}

error() {
    printf "\033[1;31m%s\033[0m\n" "$1" >&2
}

get_container_id() {
    # Heuristic: find the container with the mounted workspace directory.
    docker ps --filter "label=devcontainer.local_folder=$(pwd)" --format "{{.ID}}" | head -n 1
}

# --- Main Logic ---
main() {
    if ! command -v chezmoi >/dev/null 2>&1; then
        error "chezmoi is not installed or not in your PATH."
        error "Install it from https://www.chezmoi.io/install/ and try again."
        exit 1
    fi

    info "Looking for the running devcontainer..."
    CONTAINER_ID=$(get_container_id)

    if [ -z "$CONTAINER_ID" ]; then
        error "Could not find a running devcontainer for this project."
        error "Please ensure the container is running and you are in the project root."
        exit 1
    fi

    info "Found devcontainer with ID: $CONTAINER_ID"
    info "Applying dotfiles using 'chezmoi docker exec'..."

    # The source is the current working directory, which is mounted inside the container.
    # The destination inside the container is the vscode user's home directory.
    chezmoi docker exec \
        -- \
        "$CONTAINER_ID" \
        init --apply --source="/workspaces/dotfiles" --destination="/home/vscode"

    info "✅ Dotfiles applied successfully to the devcontainer."
}

main