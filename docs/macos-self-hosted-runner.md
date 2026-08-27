# macOS self-hosted runner

This setup sends macOS host metrics to Grafana Cloud and registers a GitHub
Actions runner for one private repository. It intentionally does not collect
system logs by default, which reduces sensitive-data exposure and Grafana Cloud
free-tier usage.

## Prerequisites

- Use a dedicated macOS user account for the runner when possible.
- Create a Grafana Cloud Free stack.
- In the stack portal, open **Prometheus details** and create an access policy
  token with `metrics:write` permission.
- Authenticate GitHub CLI with an account that administers the private target
  repository: `gh auth login`.

## Configure Grafana Cloud credentials

Add the following values to the local
`~/.config/chezmoi/chezmoi.toml`. Never commit this file.

```toml
[data.grafana_cloud]
  prometheus_url = "https://prometheus-REGION.grafana.net/api/prom/push"
  prometheus_username = "YOUR_METRICS_INSTANCE_ID"
  api_key = "glc_YOUR_ACCESS_POLICY_TOKEN"
```

Apply the dotfiles and install Homebrew dependencies:

```sh
chezmoi apply
brew bundle --file="$(chezmoi source-path)/Brewfile"
```

Configure and start Alloy from a development clone of this repository:

```sh
./scripts/setup-macos-observability.sh
```

The script validates the configuration, links it into Homebrew's Alloy
configuration directory, and starts the launchd service. The Alloy debugging UI
remains bound to `127.0.0.1:12345`.

Verify the local service:

```sh
brew services info grafana/grafana/alloy
curl --fail http://127.0.0.1:12345/-/ready
```

In Grafana Cloud, install the macOS integration or import its dashboards after
the first metrics arrive. The configuration collects host and Alloy health
metrics at low cardinality. Logs are excluded initially.

## Register the GitHub Actions runner

Run the setup against the private repository. The optional second and third
arguments set the runner name and custom labels.

```sh
./scripts/setup-github-actions-runner.sh OWNER/PRIVATE_REPOSITORY
./scripts/setup-github-actions-runner.sh OWNER/PRIVATE_REPOSITORY mac-builder macos,ai-development
```

The script:

1. Refuses public repositories and root execution.
2. Downloads the latest runner for the Mac architecture.
3. Verifies the release asset SHA-256 digest from the GitHub API.
4. Requests a short-lived registration token without saving it.
5. Installs and starts the official launchd service.

Copy `examples/self-hosted-runner-smoke.yml` into the private repository as
`.github/workflows/self-hosted-runner-smoke.yml`. It only runs through manual
dispatch.

## Add GitHub Actions data to Grafana

Install Grafana's GitHub data source in the Cloud stack. Prefer a GitHub App
restricted to the target private repository. Use the `Workflow runs` and
`Workflow usage` query types for execution status, duration, and usage. This is
separate from Alloy host metrics so GitHub metadata does not increase host
metric cardinality.

## Operations

Check the runner service from its installation directory:

```sh
cd "$HOME/actions-runner"
./svc.sh status
```

Stop or remove the runner service before deleting its directory:

```sh
cd "$HOME/actions-runner"
./svc.sh stop
./svc.sh uninstall
./config.sh remove
```

Do not run workflows from untrusted forks on this runner. Keep repository
Actions permissions and secrets at the minimum required level.

## Known macOS limitation

Alloy's embedded process exporter reads Linux `/proc` and cannot provide the
same per-process grouping on macOS. The initial setup therefore correlates host
resource graphs with GitHub workflow timestamps. Add a dedicated macOS process
collector later only if host-level correlation is insufficient.
