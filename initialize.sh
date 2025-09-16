#!/bin/sh

# This script initializes the dotfiles repository.
# It sets up symlinks and installs packages using Homebrew.

set -e

# --- Helper Functions ---
info() {
    printf "\033[1;34m%s\033[0m\n" "$1"
}

get_dotfiles_root() {
    # Resolves the absolute path of the script's directory.
    cd "$(dirname "$0")" && pwd
}

# --- Main Logic ---
DOTFILES_ROOT=$(get_dotfiles_root)
info "Dotfiles repository located at: $DOTFILES_ROOT"

# --- Symlinking Function ---
link_dotfiles() {
    info "Creating symlinks..."

    # List of files/directories to symlink in the format "source:destination".
    links="
.agignore:.agignore
.bashrc:.bashrc
.config:.config
.default-cargo-crates:.default-cargo-crates
.default-npm-packages:.default-npm-packages
.editorconfig:.editorconfig
.gitconfig:.gitconfig
.ignore:.gitignore
.tmux.conf:.tmux.conf
.zshrc:.zshrc
"

    for item in $links; do
        # Split the item into source and destination
        source_in_repo=$(echo "$item" | cut -d: -f1)
        dest_in_home=$(echo "$item" | cut -d: -f2)

        source_path="$DOTFILES_ROOT/$source_in_repo"
        dest_path="$HOME/$dest_in_home"

        if [ ! -e "$source_path" ]; then
            echo "Warning: Source file not found, skipping: $source_path"
            continue
        fi

        # Backup existing file or symlink before creating a new one
        if [ -e "$dest_path" ] || [ -L "$dest_path" ]; then
            echo "Backing up existing file at $dest_path to $dest_path.bak"
            mv "$dest_path" "$dest_path.bak"
        fi

        echo "Linking $source_path to $dest_path"
        ln -s "$source_path" "$dest_path"
    done

    info "Symlinking complete."
}

# --- Homebrew Setup Function ---
setup_homebrew() {
    info "Running Homebrew setup..."
    if ! command -v brew >/dev/null; then
        info "Homebrew is not installed. Please install it first to manage packages."
        info "See: https://brew.sh/"
        return
    fi

    local brewfile="$DOTFILES_ROOT/Brewfile"
    if [ -f "$brewfile" ]; then
        info "Found Brewfile. Installing packages..."
        brew bundle --file="$brewfile"
    else
        info "Brewfile not found in root directory. Skipping package installation."
    fi
}

# --- Main Execution ---
main() {
    link_dotfiles

    # Run Homebrew setup on supported platforms
    case "$(uname -s)" in
        Darwin|Linux)
            setup_homebrew
            ;;
        *)
            info "Unsupported OS: $(uname -s). Skipping Homebrew setup."
            ;;
    esac

    info "✅ Dotfiles initialization complete!"
}

# Execute the script
main
