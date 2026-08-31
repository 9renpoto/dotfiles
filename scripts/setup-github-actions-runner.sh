#!/bin/sh

set -eu

info() {
  printf '%s\n' "$1"
}

fail() {
  printf 'Error: %s\n' "$1" >&2
  exit 1
}

usage() {
  printf 'Usage: %s OWNER/PRIVATE_REPOSITORY [RUNNER_NAME] [LABELS]\n' "$0" >&2
  exit 2
}

[ "$#" -ge 1 ] && [ "$#" -le 3 ] || usage
[ "$(uname -s)" = "Darwin" ] || fail "This setup is only supported on macOS."
[ "$(id -u)" -ne 0 ] || fail "Do not configure a GitHub Actions runner as root."

command -v gh >/dev/null 2>&1 || fail "GitHub CLI is required."
command -v jq >/dev/null 2>&1 || fail "jq is required."
command -v curl >/dev/null 2>&1 || fail "curl is required."
gh auth status >/dev/null 2>&1 || fail "Authenticate GitHub CLI with gh auth login first."

repository=$1
runner_name=${2:-"$(scutil --get LocalHostName 2>/dev/null || hostname -s)"}
labels=${3:-"macos,ai-development"}

case "$repository" in
  */*) ;;
  *) usage ;;
esac

visibility=$(gh repo view "$repository" --json visibility --jq '.visibility')
[ "$visibility" = "PRIVATE" ] || fail "$repository is not a private repository."

case "$(uname -m)" in
  x86_64) runner_arch=x64 ;;
  arm64) runner_arch=arm64 ;;
  *) fail "Unsupported architecture: $(uname -m)" ;;
esac

runner_dir=${ACTIONS_RUNNER_DIR:-"$HOME/actions-runner"}
if [ -f "$runner_dir/.runner" ]; then
  info "A runner is already configured in $runner_dir."
  "$runner_dir/svc.sh" status
  exit 0
fi

if [ -d "$runner_dir" ] && [ -n "$(find "$runner_dir" -mindepth 1 -maxdepth 1 -print -quit)" ]; then
  fail "$runner_dir is not empty. Choose another ACTIONS_RUNNER_DIR."
fi

temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT HUP INT TERM

release_json="$temp_dir/release.json"
archive="$temp_dir/actions-runner.tar.gz"
gh api repos/actions/runner/releases/latest > "$release_json"

version=$(jq -r '.tag_name | ltrimstr("v")' "$release_json")
asset_name="actions-runner-osx-$runner_arch-$version.tar.gz"
asset_url=$(jq -r --arg name "$asset_name" '.assets[] | select(.name == $name) | .browser_download_url' "$release_json")
asset_digest=$(jq -r --arg name "$asset_name" '.assets[] | select(.name == $name) | .digest // empty' "$release_json")

[ -n "$asset_url" ] || fail "Could not find release asset $asset_name."
[ -n "$asset_digest" ] || fail "GitHub did not provide a digest for $asset_name."

info "Downloading GitHub Actions runner $version for $runner_arch."
curl -fL --retry 3 --output "$archive" "$asset_url"

expected_digest=${asset_digest#sha256:}
actual_digest=$(shasum -a 256 "$archive" | awk '{print $1}')
[ "$actual_digest" = "$expected_digest" ] || fail "Runner archive checksum verification failed."

mkdir -p "$runner_dir"
tar -xzf "$archive" -C "$runner_dir"

registration_token=$(gh api --method POST "repos/$repository/actions/runners/registration-token" --jq '.token')

(
  cd "$runner_dir"
  ./config.sh \
    --unattended \
    --url "https://github.com/$repository" \
    --token "$registration_token" \
    --name "$runner_name" \
    --labels "$labels" \
    --work _work
  ./svc.sh install
  ./svc.sh start
  ./svc.sh status
)

info "GitHub Actions runner $runner_name is registered for $repository."
