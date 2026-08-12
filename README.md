# dotfiles

Opinionated dotfiles and setup scripts for macOS, Linux/WSL, and Windows terminals.

## Quick Start

1. Install [`chezmoi`](https://www.chezmoi.io/):
   - macOS / Linux: `brew install chezmoi`
   - Windows: `winget install twpayne.chezmoi`
2. Apply the configuration: `chezmoi init 9renpoto/dotfiles --apply`
   - During `init`, you will be prompted for your email and other information. These values are used to personalize your configuration.
3. Finish runtime setup:
   - Run `mise install` to install language runtimes.
   - Switch to **zsh** (highly recommended): `chsh -s $(which zsh)`

Re-run `chezmoi apply` whenever you pull updates to keep `$HOME` in sync.

### Shell Recommendation

While this repository supports multiple shells, **zsh** is the recommended default. Our zsh configuration includes:
- [antidote](https://getantidote.github.io/) for fast plugin management.
- Syntax highlighting and autosuggestions.
- [starship](https://starship.rs/) prompt integration.

### macOS Preparation

Before running `chezmoi apply` on macOS, make sure the base toolchain is ready:

- Install Apple's Command Line Tools: `xcode-select --install`
- Install Homebrew if it's missing (see [brew.sh](https://brew.sh/) for the latest install command)
- Confirm Homebrew works: `brew doctor`

### Windows Preparation

Run the Windows bootstrap from PowerShell in a development clone:

```powershell
Set-ExecutionPolicy -Scope Process Bypass
.\initialize.ps1
```

The script imports `winget-packages.json`, refreshes `PATH` for the current
process, and applies the local source state with `chezmoi`. Use
`.\initialize.ps1 -SkipPackages` to reapply only the dotfiles, or
`.\initialize.ps1 -SkipApply` to install only the packages. It also configures
the ghq root as `%USERPROFILE%\src` on Windows and the corresponding
`/mnt/<drive>/Users/<user>/src` path in the default WSL distribution. Override
the location with `.\initialize.ps1 -GhqRoot D:\src` when needed.

## Configuration

Some configurations can be customized using variables. Create a `~/.config/chezmoi/chezmoi.toml` file to override default values.

For example, to change the Git email address:

```toml
[data]
  email = "your-email@example.com"
```

`chezmoi apply` will then use this email address in your `~/.gitconfig`.

### Environment-Specific Data

- Global defaults live in `chezmoidata.toml`; platform-specific files such as `chezmoidata.darwin.toml.tmpl` and `chezmoidata.linux.toml` layer on automatic values (e.g., the Homebrew prefix).
- Sensitive or machine-local overrides can stay outside the repo: use `~/.config/chezmoi/chezmoi.toml` and nest keys under `[data]` to match the structure in `chezmoidata`.

### GitHub Authentication

`~/.ssh/config` is **not** managed by this repository. GitHub authentication is delegated to the GitHub CLI:

1. Authenticate once per machine: `gh auth login`
2. `run_after_setup-git-credentials.sh` runs `gh auth setup-git` after every `chezmoi apply`, which registers `gh` as the git credential helper for GitHub over HTTPS.

The script re-runs on each apply because applying rewrites `~/.gitconfig` from `dot_gitconfig.tmpl`, which is where `gh auth setup-git` stores its entries. It exits quietly when `gh` is missing or not authenticated, so first-time setup order does not matter — just re-run `chezmoi apply` after `gh auth login`.

Run `gh config set git_protocol https` if you want git remotes to go through that credential helper. With the `ssh` protocol, `gh` keeps handing out SSH URLs and git authenticates with your default key (`~/.ssh/id_ed25519`) instead.

### Secret Management

Secrets are managed locally through `~/.config/chezmoi/chezmoi.toml`.

#### Quick Setup

Run the interactive setup script to create your configuration file:

```sh
./setup-chezmoi-config.sh
```

This script will prompt you for:
- WakaTime API key (optional)
- Email address override (optional)
- Machine profile (optional)

#### Manual Setup

Alternatively, manually create or edit `~/.config/chezmoi/chezmoi.toml`:

```toml
[data.wakatime]
  api_key = "waka_YOUR_API_KEY_HERE"

[data.user]
  email = "your-email@example.com"

[data.machine]
  profile = "dev"
```

After configuration, run `chezmoi apply` to generate your dotfiles with the provided values.

**Important**: Never commit `~/.config/chezmoi/chezmoi.toml` to version control. This file should remain local to each machine.

### Devcontainer

This repository includes a devcontainer configuration that allows you to work in a consistent, pre-configured environment.

1.  **Open in Devcontainer**: Open this project in a VS Code devcontainer. When the container starts, it will install necessary tools but will not apply dotfiles automatically.
2.  **Run Setup from Host**: From your **host machine's terminal**, run the following command in the project root:

    ```sh
    ./devcontainer-setup.sh
    ```

    This script uses `chezmoi docker` to apply the dotfiles from your source tree directly into the running container's home directory.

#### macOS Template Override Example

Use `chezmoidata.darwin.toml.tmpl` as a starting point when you need to pin macOS-specific values locally:

1. Render the current template to your local config: `chezmoi execute-template < "$(chezmoi source-path)/chezmoidata.darwin.toml.tmpl" > ~/.config/chezmoi/chezmoidata.darwin.toml`
2. Edit the generated file and override only the keys you need, for example:

   ```toml
   [paths]
   homebrew_prefix = "/opt/homebrew"

   [machine]
   profile = "work"
   ```

3. Re-run `chezmoi apply` (or `chezmoi update --apply`) to pick up the overrides.

## Maintenance Commands

- `chezmoi update --apply` pulls the latest repository changes and reapplies them.
- `brew bundle --file=Brewfile` keeps CLI/GUI packages in sync with the tracked Brewfile.
- `brew bundle check --file=Brewfile` inspects for drift; pair with `brew bundle cleanup --file=Brewfile` to prune unused packages.
- `winget import --import-file winget-packages.json --ignore-unavailable --ignore-versions` installs the Windows toolset.
- `lefthook run pre-commit --all-files` runs Secretlint locally; enable the hook permanently with `lefthook install`.

## Working Locally

When hacking on the repository from a development clone, `./initialize.sh` ensures `chezmoi` is installed, runs an interactive configuration setup if needed, applies the local working tree with `chezmoi apply`, and optionally runs the Brewfile on macOS/Linux hosts.

1. Clone the repository.
2. Run `./initialize.sh`.
3. Follow the prompts to configure your email and other settings.
   - Run `gh auth login` so `gh auth setup-git` can wire up git credentials (see [GitHub Authentication](#github-authentication)).
4. Switch to **zsh** if you haven't already: `chsh -s $(which zsh)`.

## Configuration Layout

Source files follow the `chezmoi` naming conventions: `dot_*` represent files that land directly in `$HOME` (e.g., `dot_zshrc` → `~/.zshrc`) and `dot_config/` mirrors `~/.config/`. Update the source tree with these prefixes when adding new files so `chezmoi` can track them correctly.

## Contributing

Review [AGENTS.md](AGENTS.md) for contributor expectations, coding standards, and pull request requirements.

## License

MIT
