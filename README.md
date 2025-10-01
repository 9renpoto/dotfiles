# dotfiles

Opinionated dotfiles and setup scripts for macOS, Linux/WSL, and Windows terminals.

## Quick Start

1. Install [`chezmoi`](https://www.chezmoi.io/):
   - macOS / Linux: `brew install chezmoi`
   - Windows: `winget install twpayne.chezmoi`
2. Apply the configuration: `chezmoi init 9renpoto/dotfiles --apply`
3. Finish runtime setup: `mise install`

Re-run `chezmoi apply` whenever you pull updates to keep `$HOME` in sync.

### macOS Preparation

Before running `chezmoi apply` on macOS, make sure the base toolchain is ready:

- Install Apple's Command Line Tools: `xcode-select --install`
- Install Homebrew if it's missing (see [brew.sh](https://brew.sh/) for the latest install command)
- Confirm Homebrew works: `brew doctor`

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
- Host-specific tweaks can be added by extending the override maps inside the platform templates or by adding per-host data files that follow chezmoi's naming (for example `chezmoidata.<hostname>.toml`).
- Sensitive or machine-local overrides can stay outside the repo: use `~/.config/chezmoi/chezmoi.toml` and nest keys under `[data]` to match the structure in `chezmoidata`. For the GitHub SSH key, set:

  ```toml
  [data.ssh]
    github_identity_file = "~/.ssh/your_custom_key"
  ```

- Templates such as `private_dot_ssh/config.tmpl` read `data.ssh.github_identity_file`, falling back to `~/.ssh/id_ed25519` when nothing is defined. This keeps the SSH config portable while letting each environment opt into a different key path.

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
- `lefthook run pre-commit --all-files` runs Secretlint locally; enable the hook permanently with `lefthook install`.

## Working Locally

When hacking on the repository from a development clone, `./initialize.sh` ensures `chezmoi` is installed, applies the local working tree with `chezmoi apply`, and optionally runs the Brewfile on macOS/Linux hosts.

## Configuration Layout

Source files follow the `chezmoi` naming conventions: `dot_*` represent files that land directly in `$HOME` (e.g., `dot_zshrc` → `~/.zshrc`) and `dot_config/` mirrors `~/.config/`. Update the source tree with these prefixes when adding new files so `chezmoi` can track them correctly.

## Contributing

Review [AGENTS.md](AGENTS.md) for contributor expectations, coding standards, and pull request requirements.

## License

MIT
