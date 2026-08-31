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

toml_string_value() {
    section=$1
    key=$2

    [ -f "$CHEZMOI_CONFIG_FILE" ] || return 0

    awk -v section="[$section]" -v key="$key" '
        $0 == section { in_section = 1; next }
        in_section && /^\[/ { exit }
        in_section && $0 ~ "^[[:space:]]*" key "[[:space:]]*=" {
            sub("^[[:space:]]*" key "[[:space:]]*=[[:space:]]*", "")
            sub(/[[:space:]]*#.*$/, "")
            sub(/^"/, "")
            sub(/"[[:space:]]*$/, "")
            print
            exit
        }
    ' "$CHEZMOI_CONFIG_FILE"
}

toml_escape() {
    value=$1
    value=${value//\\/\\\\}
    value=${value//\"/\\\"}
    printf '%s' "$value"
}

configure_grafana_cloud() {
    existing_url=$(toml_string_value "data.grafana_cloud" "prometheus_url")
    existing_username=$(toml_string_value "data.grafana_cloud" "prometheus_username")
    existing_api_key=$(toml_string_value "data.grafana_cloud" "api_key")

    read -rp "Configure Grafana Cloud metrics? (y/N): " response
    case "$response" in
        [yY][eE][sS]|[yY]) ;;
        *) return 1 ;;
    esac

    if [ -n "$existing_url" ]; then
        read -rp "Prometheus remote write URL [keep existing]: " grafana_prometheus_url
        grafana_prometheus_url=${grafana_prometheus_url:-$existing_url}
    else
        read -rp "Prometheus remote write URL: " grafana_prometheus_url
    fi

    if [ -n "$existing_username" ]; then
        read -rp "Metrics instance ID [keep existing]: " grafana_prometheus_username
        grafana_prometheus_username=${grafana_prometheus_username:-$existing_username}
    else
        read -rp "Metrics instance ID: " grafana_prometheus_username
    fi

    if [ -n "$existing_api_key" ]; then
        read -rsp "Access policy token [keep existing]: " grafana_api_key
        grafana_api_key=${grafana_api_key:-$existing_api_key}
    else
        read -rsp "Access policy token (metrics:write): " grafana_api_key
    fi
    echo

    if [ -z "$grafana_prometheus_url" ] ||
       [ -z "$grafana_prometheus_username" ] ||
       [ -z "$grafana_api_key" ]; then
        error "All Grafana Cloud values are required when monitoring is enabled."
        exit 1
    fi

    return 0
}

upsert_grafana_cloud() {
    temporary_file=$(mktemp "${CHEZMOI_CONFIG_FILE}.tmp.XXXXXX")
    trap 'rm -f "${temporary_file:-}"' EXIT

    awk '
        $0 == "[data.grafana_cloud]" { skip = 1; next }
        skip && /^\[/ { skip = 0 }
        !skip { print }
    ' "$CHEZMOI_CONFIG_FILE" > "$temporary_file"

    {
        printf '\n[data.grafana_cloud]\n'
        printf '  prometheus_url = "%s"\n' "$(toml_escape "$grafana_prometheus_url")"
        printf '  prometheus_username = "%s"\n' "$(toml_escape "$grafana_prometheus_username")"
        printf '  api_key = "%s"\n' "$(toml_escape "$grafana_api_key")"
    } >> "$temporary_file"

    chmod 600 "$temporary_file"
    mv "$temporary_file" "$CHEZMOI_CONFIG_FILE"
    trap - EXIT
}

main() {
    info "Chezmoi Configuration Setup"
    echo

    umask 077

    # Existing configurations are updated in place to preserve unknown sections.
    if [ -f "$CHEZMOI_CONFIG_FILE" ]; then
        warn "Configuration file already exists at: $CHEZMOI_CONFIG_FILE"
        echo
        if configure_grafana_cloud; then
            upsert_grafana_cloud
            success "✅ Grafana Cloud configuration updated."
        else
            info "Keeping existing configuration. Exiting."
        fi
        exit 0
    fi

    # Create config directory if it doesn't exist
    mkdir -p "$CHEZMOI_CONFIG_DIR"

    # Interactive prompts
    echo "This script will help you set up your chezmoi configuration."
    echo "Press Enter to skip optional fields."
    echo

    # WakaTime API Key
    read -rsp "Enter your WakaTime API key (optional): " wakatime_api_key
    echo

    # Email (optional override)
    read -rp "Enter your email address [9renpoto@gmail.com]: " email
    email=${email:-9renpoto@gmail.com}
    echo

    # Machine profile (optional)
    read -rp "Enter machine profile (e.g., dev, work, personal) [dev]: " machine_profile
    machine_profile=${machine_profile:-dev}
    echo

    grafana_prometheus_url=""
    grafana_prometheus_username=""
    grafana_api_key=""
    grafana_enabled=false
    if [ "$(uname -s)" = "Darwin" ] && configure_grafana_cloud; then
        grafana_enabled=true
    fi
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

        # Machine profile
        if [ -n "${machine_profile:-}" ]; then
            echo ""
            echo "[data.machine]"
            echo "  profile = \"$machine_profile\""
        fi

        # Grafana Cloud
        if [ "$grafana_enabled" = true ]; then
            echo ""
            echo "[data.grafana_cloud]"
            printf '  prometheus_url = "%s"\n' "$(toml_escape "$grafana_prometheus_url")"
            printf '  prometheus_username = "%s"\n' "$(toml_escape "$grafana_prometheus_username")"
            printf '  api_key = "%s"\n' "$(toml_escape "$grafana_api_key")"
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

    if [ "$grafana_enabled" = true ]; then
        echo
        info "Grafana Cloud is configured. Run the macOS observability setup after chezmoi apply."
    fi

    echo
    info "You can manually edit the file at any time:"
    echo "  ${CHEZMOI_CONFIG_FILE}"
    echo
    info "To apply your dotfiles with this configuration:"
    echo "  chezmoi apply"
}

main "$@"
