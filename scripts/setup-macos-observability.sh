#!/bin/sh

set -eu

info() {
  printf '%s\n' "$1"
}

fail() {
  printf 'Error: %s\n' "$1" >&2
  exit 1
}

link_config() {
  source_file=$1
  target_file=$2

  if [ -L "$target_file" ]; then
    current_target=$(readlink "$target_file")
    if [ "$current_target" = "$source_file" ]; then
      return
    fi
  elif [ -e "$target_file" ]; then
    fail "$target_file already exists. Move it aside before running this script."
  fi

  ln -sfn "$source_file" "$target_file"
}

[ "$(uname -s)" = "Darwin" ] || fail "This setup is only supported on macOS."
command -v brew >/dev/null 2>&1 || fail "Homebrew is required."
command -v alloy >/dev/null 2>&1 || fail "Grafana Alloy is not installed. Run brew bundle first."

config_home=${XDG_CONFIG_HOME:-"$HOME/.config"}
source_dir="$config_home/alloy"
source_config="$source_dir/config.alloy"
source_env="$source_dir/config.env"
source_args="$source_dir/extra-args.txt"

[ -f "$source_config" ] || fail "$source_config was not found. Run chezmoi apply first."
[ -f "$source_env" ] || fail "$source_env was not found. Run chezmoi apply first."
[ -f "$source_args" ] || fail "$source_args was not found. Run chezmoi apply first."

# Load credentials only into this process; do not print them.
# shellcheck disable=SC1090
. "$source_env"
: "${GRAFANA_CLOUD_PROMETHEUS_URL:?Configure grafana_cloud.prometheus_url in chezmoi.toml}"
: "${GRAFANA_CLOUD_PROMETHEUS_USERNAME:?Configure grafana_cloud.prometheus_username in chezmoi.toml}"
: "${GRAFANA_CLOUD_API_KEY:?Configure grafana_cloud.api_key in chezmoi.toml}"

alloy validate "$source_config"

brew_prefix=$(brew --prefix)
target_dir="$brew_prefix/etc/alloy"
mkdir -p "$target_dir"

link_config "$source_config" "$target_dir/config.alloy"
link_config "$source_env" "$target_dir/config.env"
link_config "$source_args" "$target_dir/extra-args.txt"

chmod 600 "$source_env"
brew services restart grafana/grafana/alloy

info "Grafana Alloy is configured and running."
info "Check readiness at http://127.0.0.1:12345/-/ready"
