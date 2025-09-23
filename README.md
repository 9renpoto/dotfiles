# dotfiles

Opinionated dotfiles and setup scripts for macOS and Linux/WSL shells.

## Quick Start

1. Clone the repository into `~/src/github.com/9renpoto/dotfiles` (or similar).
2. Run `mise install` to install the Node LTS runtime defined in `.mise.toml` for tooling such as Secretlint.
3. Execute `./initialize.sh` from the repository root to create symlinks and optionally install Homebrew packages.

## Maintenance Commands

- `brew bundle --file=Brewfile` keeps CLI/GUI packages in sync with the tracked Brewfile.
- `brew bundle check --file=Brewfile` inspects for drift; pair with `brew bundle cleanup --file=Brewfile` to prune unused packages.
- `lefthook run pre-commit --all-files` runs Secretlint locally; enable the hook permanently with `lefthook install`.

## Configuration Layout

Tracked dotfiles live at the repository root (`.zshrc`, `.gitconfig`, `.bashrc`) with application-specific settings under `.config/` (Alacritty, Fish, Mise). Update `initialize.sh` whenever you add new paths so the bootstrap step keeps machines consistent.

## Contributing

Review [AGENTS.md](AGENTS.md) for contributor expectations, coding standards, and pull request requirements.

## License

MIT
