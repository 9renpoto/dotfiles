# dotfiles

Opinionated dotfiles and setup scripts for macOS, Linux/WSL, and Windows terminals.

## Quick Start

1. Install [`chezmoi`](https://www.chezmoi.io/):
   - macOS / Linux: `brew install chezmoi`
   - Windows: `winget install twpayne.chezmoi`
2. Apply the configuration: `chezmoi init 9renpoto/dotfiles --apply`
3. Finish runtime setup: `mise install`

Re-run `chezmoi apply` whenever you pull updates to keep `$HOME` in sync.

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
