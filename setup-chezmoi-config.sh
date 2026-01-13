#!/usr/bin/env bash

# Setup script for chezmoi configuration file
# Creates ~/.config/chezmoi/chezmoi.toml with user-provided values

set -euo pipefail

# --- Color Output ---
info() {
    printf "\033[1;34m%s\033[0m\n" "$1"
}

success() {
    printf "\033[1;32m%s\033[0m\n" "$1"
}

warn() {
    printf "\033[1;33m%s\033[0m\n" "$1"
}

error() {
    printf "\033[1;31m%s\033[0m\n" "$1" >&2
}

# --- Main Logic ---
CHEZMOI_CONFIG_DIR="${HOME}/.config/chezmoi"
CHEZMOI_CONFIG_FILE="${CHEZMOI_CONFIG_DIR}/chezmoi.toml"

main() {
    info "Chezmoi Configuration Setup"
    echo

    # Check if config file already exists
    if [ -f "$CHEZMOI_CONFIG_FILE" ]; then
        warn "Configuration file already exists at: $CHEZMOI_CONFIG_FILE"
        echo
        read -rp "Do you want to update it? (y/N): " response
        case "$response" in
            [yY][eE][sS]|[yY])
                info "Updating existing configuration..."
                ;;
            *)
                info "Keeping existing configuration. Exiting."
                exit 0
                ;;
        esac
    fi

    # Create config directory if it doesn't exist
    mkdir -p "$CHEZMOI_CONFIG_DIR"

    # Interactive prompts
    echo "This script will help you set up your chezmoi configuration."
    echo "Press Enter to skip optional fields."
    echo

    # WakaTime API Key
    read -rp "Enter your WakaTime API key (optional): " wakatime_api_key
    echo

    # Email (optional override)
    read -rp "Enter your email address (optional, overrides default): " email
    echo

    # SSH key path (optional override)
    read -rp "Enter custom SSH key path for GitHub (optional): " github_ssh_key
    echo

    # Machine profile (optional)
    read -rp "Enter machine profile (e.g., dev, work, personal) (optional): " machine_profile
    echo

    # Generate config file
    info "Creating configuration file at: $CHEZMOI_CONFIG_FILE"

    {
        echo "# Chezmoi configuration file"
        echo "# This file should never be committed to version control"
        echo ""
        echo "[data]"

        # Email
        if [ -n "${email:-}" ]; then
            echo ""
            echo "[data.user]"
            echo "  email = \"$email\""
        fi

        # WakaTime
        if [ -n "${wakatime_api_key:-}" ]; then
            echo ""
            echo "[data.wakatime]"
            echo "  api_key = \"$wakatime_api_key\""
        fi

        # SSH
        if [ -n "${github_ssh_key:-}" ]; then
            echo ""
            echo "[data.ssh]"
            echo "  github_identity_file = \"$github_ssh_key\""
        fi

        # Machine profile
        if [ -n "${machine_profile:-}" ]; then
            echo ""
            echo "[data.machine]"
            echo "  profile = \"$machine_profile\""
        fi

        # Add a final newline
        echo ""
    } > "$CHEZMOI_CONFIG_FILE"

    # Set appropriate permissions
    chmod 600 "$CHEZMOI_CONFIG_FILE"

    success "✅ Configuration file created successfully!"
    echo
    info "Configuration saved to: $CHEZMOI_CONFIG_FILE"

    if [ -n "${wakatime_api_key:-}" ]; then
        echo
        info "WakaTime is configured. Run 'chezmoi apply' to generate ~/.wakatime.cfg"
    fi

    echo
    info "You can manually edit the file at any time:"
    echo "  ${CHEZMOI_CONFIG_FILE}"
    echo
    info "To apply your dotfiles with this configuration:"
    echo "  chezmoi apply"
}

main "$@"
