#!/bin/sh

# Bootstrap script for the dotfiles repository.
# Ensures chezmoi is available, applies the source state, and optionally runs Homebrew.

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

seed_default_data_template() {
    template_src="$DOTFILES_ROOT/dot_config/chezmoi/data.toml.sample"
    template_dest="$HOME/.config/chezmoi/data.toml"

    if [ ! -f "$template_src" ]; then
        return
    fi

    if [ -e "$template_dest" ]; then
        info "Existing chezmoi data file detected at $template_dest. Skipping template seeding."
        return
    fi

    info "Seeding default chezmoi data template to $template_dest"
    mkdir -p "$(dirname "$template_dest")"
    cp "$template_src" "$template_dest"
}

ensure_chezmoi() {
    if command -v chezmoi >/dev/null 2>&1; then
        return 0
    fi

    info "chezmoi not found on PATH. Installing..."

    if command -v brew >/dev/null 2>&1; then
        brew install chezmoi
        return 0
    fi

    cat <<'EOM'
chezmoi is required to apply the dotfiles.
Install it using one of the supported methods:
  • macOS/Linux (Homebrew): brew install chezmoi
  • Windows (winget): winget install twpayne.chezmoi
  • Manual options: https://www.chezmoi.io/install/
Re-run initialize.sh after installing chezmoi.
EOM
    exit 1
}

apply_with_chezmoi() {
    info "Applying dotfiles using chezmoi..."
    chezmoi apply --source="$DOTFILES_ROOT" --destination="$HOME"
    info "chezmoi apply completed."
}

setup_homebrew() {
    info "Running Homebrew setup..."
    if ! command -v brew >/dev/null 2>&1; then
        info "Homebrew not detected. Skipping package installation."
        return
    fi

    brewfile="$DOTFILES_ROOT/Brewfile"
    if [ -f "$brewfile" ]; then
        info "Installing packages listed in Brewfile..."
        brew bundle --file="$brewfile"
    else
        info "Brewfile not found in repository root."
    fi
}

main() {
    seed_default_data_template
    ensure_chezmoi
    apply_with_chezmoi

    case "$(uname -s)" in
        Darwin|Linux)
            setup_homebrew
            ;;
        *)
            info "Homebrew automation skipped on $(uname -s)."
            ;;
    esac

    info "✅ Dotfiles initialization complete!"
}

main
