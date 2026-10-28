# dotfiles

mise-managed workstation configuration for macOS, Linux/WSL, and Windows.

Deploy configuration files, install platform packages and runtime tools, and
integrate workstation services from one checkout. Configuration sources live
under `home/`; machine-specific values stay in untracked local configuration.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Install

Requires **mise 2026.10.19 or newer**, Git, and the platform build toolchain
(Rust/Cargo and C/C++ tools for hk). Install Homebrew on macOS/Linux first.

Clone this repository, review `mise.toml`, and configure local values and newly
issued secrets as described under [local configuration and secrets](#local-configuration-and-secrets)
before running:

```sh
mise trust
./initialize.sh
```

On Windows, install mise (`winget install --id jdx.mise --exact`) and run
`./initialize.ps1` from PowerShell. The wrapper refreshes PATH after package
installation, installs the PowerShell 5.1/7 shared-profile loaders without
replacing custom profiles, and synchronizes the ghq root with the default WSL
distribution. The ghq root is saved in `mise.local.toml`, preserving other local
values, so ordinary `mise run apply` keeps using it. Options: `-SkipPackages`, `-SkipApply`, `-GhqRoot D:\src`.
`-ConfigureWakaTime` reads a hidden key for this invocation only. Store it in
local configuration to persist it.

`mise bootstrap` installs Homebrew packages from `Brewfile` through a pre-packages
hook (or imports `winget-packages.json` on Windows) and natively manages
Neovim/TPM Git checkouts,
dotfiles and runtime tools, then runs the integration task. Existing external
repositories are preserved; `mise bootstrap --update` explicitly updates them.
Review the package lists before bootstrap. Homebrew itself is a prerequisite.
`Brewfile` is the single source for Homebrew taps, formulae and casks.
Homebrew determines platform compatibility; this repository adds no OS guards
to package declarations. The Firefox language option is preserved.
`winget-packages.json` is the single source for Windows packages. Both manifests
are invoked by `mise run packages` and the bootstrap pre-packages hook.

<details>
<summary>日本語</summary>

macOS、Linux/WSL、Windows の設定を mise で管理します。mise 2026.10.19 以降、
Git、hk のビルド環境が必要です。macOS/Linux では Homebrew も先に導入します。
設定とパッケージ一覧を確認し、ローカル設定を用意してから `mise trust` と
`./initialize.sh` を実行します。Windows は PowerShell から `./initialize.ps1` を
実行します。既存の PowerShell プロファイルは保持されます。
Homebrew パッケージは `Brewfile`、Windows パッケージは
`winget-packages.json` が管理元です。外部リポジトリの更新は
`mise bootstrap --update` で明示的に行います。

</details>

## Usage

### Local configuration and secrets

Rotate WakaTime and Grafana credentials with their providers and configure the
new values; setup does not import existing credentials or local settings.
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

### Daily use

Run commands inside the checkout (or use `mise -C /path/to/dotfiles ...`):

```sh
mise dot diff
mise run apply
mise bootstrap --dry-run
mise run packages
brew bundle check --file=Brewfile
# Windows package manifest:
# winget import -i winget-packages.json --ignore-unavailable --ignore-versions
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

`Skillsfile` is synchronized by the integration task on Unix. Only entries
recorded by this setup are eligible for removal. Personal Codex skills live in
`home/.agents/skills/`. Windows does not synchronize `gh skill` entries.

<details>
<summary>日本語</summary>

端末固有の値と新しく発行した認証情報は `mise.local.toml` の `[vars]` に設定します。
このファイルはコミットせず、Unix では権限を 0600 にします。
`WAKATIME_API_KEY` と `GRAFANA_CLOUD_API_KEY` はローカルのキーを上書きできます。
設定元の `home/` を編集し、チェックアウト内で `mise run apply` を実行します。
配置先の直接編集は次回の適用で上書きされる可能性があります。
競合は確認してから必要に応じて `mise dot apply --force` で解消します。
スキル同期で削除対象になるのは、この仕組みが管理したスキルだけです。

</details>

## Contributing

Edit configuration sources under `home/` and declare destinations in
`mise.toml`. Use Conventional Commits and keep each change independently
reviewable. Run both checks before submitting a pull request:

```sh
mise run check
hk check --all
```

`mise run check` deploys into a temporary home, verifies idempotency, secret
permissions and Git/Biome output without installing packages or executing
integration scripts. CI runs this same check. Devcontainer setup installs mise;
run `./devcontainer-setup.sh` from the host to opt into dotfiles deployment.

The shared hk settings live in `home/.config/hk/config.pkl`; install its global
launcher with `mise exec cargo:hk -- hk install --global --mise` (Git 2.54+).
macOS runner and monitoring setup is documented in
[the runner guide](docs/macos-self-hosted-runner.md).

The decision to migrate workstation automation to Go and publish Octocov
coverage badges on GitHub Pages is recorded in
[ADR 0001](docs/adr/0001-go-automation-and-coverage.md).
`mise run check` runs the Go test suite, including isolated deployment and
Biome configuration checks. CI records `coverage.out` as the `go-coverage`
artifact and reports it in the Octocov Job Summary. Public badge hosting
remains tracked in the coverage migration issue.

<details>
<summary>日本語</summary>

変更は独立してレビューできる単位に分け、Conventional Commits を使います。
提出前に `mise run check` と `hk check --all` を実行してください。
現在のテストは一時ホームへの設定配置を検証し、パッケージ導入や統合処理は実行しません。
Go への移行と Octocov・GitHub Pages の採用方針は ADR 0001 に記録しています。
Go テストへの移行後、CI はカバレッジを artifact と Octocov の Job Summary に記録します。
badge 公開はカバレッジ移行 Issue で追跡します。

</details>

## License

No license has been specified in this repository.

<details>
<summary>日本語</summary>

このリポジトリにはライセンスが指定されていません。

</details>
