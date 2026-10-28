# dotfiles

mise-managed workstation configuration for macOS, Linux/WSL, and Windows.
Requires **mise 2026.10.19 or newer**, Git, and the platform build toolchain
(Rust/Cargo and C/C++ tools for hk). Install Homebrew on macOS/Linux first.

## Setup

Clone this repository, review `mise.toml`, then run:

```sh
mise trust
./initialize.sh
```

On Windows, install mise (`winget install --id jdx.mise --exact`) and run
`./initialize.ps1` from PowerShell. The wrapper refreshes PATH after package
installation, installs the PowerShell 5.1/7 shared-profile loaders without
replacing custom profiles, and synchronizes the ghq root with the default WSL
distribution. Options: `-SkipPackages`, `-SkipApply`, `-GhqRoot D:\src`.
`-ConfigureWakaTime` reads a hidden key for this invocation only. Store it in
local configuration to persist it.

`mise bootstrap` installs Homebrew packages from `Brewfile` through a pre-packages
hook and natively manages Windows packages, Neovim/TPM Git checkouts,
dotfiles and runtime tools, then runs the integration task. Existing external
repositories are preserved; `mise bootstrap --update` explicitly updates them.
Review the package lists before bootstrap. Homebrew itself is a prerequisite.
`Brewfile` is the single source for Homebrew taps, formulae and casks; casks are
guarded by `OS.mac?` so the same file works on Linux. The Firefox language
option is preserved. Windows packages remain in `mise.toml`.

## Existing chezmoi users

The initialization wrappers copy known personal values from
`~/.config/chezmoi/chezmoi.toml` into `mise.local.toml` once, using a managed
Python runtime. They never modify or remove the old file. An existing
`mise.local.toml` is preserved; merge old values manually if it already exists.
Back up local edits and compare `mise dot diff` before the first application.
Do not run the old manager after migration: both would write the same files.
After checking all configurations, uninstall chezmoi using your package manager.
The old configuration and ownership record are only read for migration.

## Local configuration and secrets

Copy `mise.local.toml.example` to `mise.local.toml`, then `chmod 600 mise.local.toml`
on Unix. Use your Windows account's private checkout and restrict access with
Windows ACLs. This file is ignored by Git and must not be committed or copied
into a shared container. WakaTime and Grafana output files use mode 0600 on
Unix; mise does not enforce POSIX permissions on Windows.

Environment variables `WAKATIME_API_KEY` and `GRAFANA_CLOUD_API_KEY` can override
the local keys. Missing optional keys produce empty configuration values.
Other values, including Git email, ghq root and Grafana endpoint, live under
`[vars]`. Templates use Tera, with `os()` for branching. Target variants select
Biome's standard path on each OS. Use `mise.work.toml` and `mise -E work ...`
for profiles; `mise.local.toml` overrides shared defaults.

## Daily use

Run commands inside the checkout (or use `mise -C /path/to/dotfiles ...`):

```sh
mise dot diff
mise run apply
mise bootstrap --dry-run
mise run packages
brew bundle check --file=Brewfile
# Windows packages only:
mise bootstrap packages status
mise bootstrap repos update
mise run check
```

Deployment uses copies, so applications receive ordinary files on Windows too.
Edit sources under `home/`, then apply; direct destination edits can be overwritten.
Conflicts are reported by mise: inspect them and back up files before explicitly
using `mise dot apply --force`. Nothing in the wrappers forces replacement.
Global runtime versions are deployed separately to `~/.config/mise/config.toml`.
The workstation tasks and local values belong to this checkout.

Git's common settings live in `~/.config/git/config`. On Windows the existing
main `.gitconfig` is preserved and includes that file. On Unix the main config
is rendered and `gh auth setup-git` is reapplied afterward; authenticate using
`gh auth login` and run `mise run apply` again.

`Skillsfile` is synchronized by the integration task on Unix. It adopts the old
ownership record and removes only previously owned entries. Personal Codex
skills live in `home/.agents/skills/`. Windows retains the previous behavior of
not synchronizing `gh skill` entries.

## Development

`mise run check` deploys into a temporary home, verifies idempotency, secret
permissions and Git/Biome output without installing packages or executing
integration scripts. CI runs this same check. Devcontainer setup installs mise;
run `./devcontainer-setup.sh` from the host to opt into dotfiles deployment.

The shared hk settings live in `home/.config/hk/config.pkl`; install its global
launcher with `mise exec cargo:hk -- hk install --global --mise` (Git 2.54+).
macOS runner and monitoring setup is documented in
[the runner guide](docs/macos-self-hosted-runner.md).
