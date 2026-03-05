# Copilot Instructions

[chezmoi](https://www.chezmoi.io/)-managed dotfiles for macOS, Linux/WSL, and Windows. Source files are the *source state*; `chezmoi apply` renders them into `$HOME`.

## File Naming

| Prefix | Destination |
|---|---|
| `dot_foo` | `~/.foo` |
| `dot_config/` | `~/.config/` |
| `private_dot_foo` | `~/.foo` (restricted permissions) |
| `foo.tmpl` | rendered via Go templates |

`.chezmoiexternal.toml` pulls external repos (currently: Neovim config). `.chezmoiignore` is itself a template.

## Template Data

Defaults and platform values (Homebrew prefix, package manager) are defined/computed in `dot_chezmoi.toml.tmpl` from `chezmoi.os`/`chezmoi.arch`. Machine-local overrides in `~/.config/chezmoi/chezmoi.toml` (never commit this file).

Platform branching: `{{ if eq .chezmoi.os "darwin" }} ... {{ else if eq .chezmoi.os "linux" }} ... {{ end }}`

## Commands

```sh
chezmoi apply --source=. --destination="$HOME"  # apply from working tree
chezmoi update --apply                           # pull + apply in one step
./initialize.sh                                  # bootstrap (installs chezmoi + Brewfile)
bash .devcontainer/scripts/check-chezmoi.sh     # validate templates (same as CI)

bun install                                      # install JS toolchain (after cloning)
lefthook install                                 # enable pre-commit hook
lefthook run pre-commit --all-files              # run Secretlint manually

brew bundle --file=Brewfile                      # install packages
brew bundle check --file=Brewfile               # check drift
brew bundle cleanup --file=Brewfile             # prune unlisted

shellcheck initialize.sh                         # lint shell scripts
mise install                                     # sync language runtimes
```

## Conventions

- **Commits**: Conventional Commit syntax — `type(scope): summary`
- **Language**: Chat responses should be in Japanese; commit messages and code comments should be in English
- **Brewfile**: grouped (taps → casks → brews), sorted within groups
- **New dotfiles**: `dot_config/<appname>/` for `~/.config/`, `dot_<name>` for `~/.<name>`
- **Secrets**: `private_dot_` prefix + template variables from `~/.config/chezmoi/chezmoi.toml`; never hardcode
- **Shell scripts**: POSIX-compatible, pass `shellcheck`; gate OS-specific logic to avoid cross-platform breakage
- **Indent**: 2 spaces (4 for Python), UTF-8, LF
- **JS/TS**: `biome check`; spelling: `cspell lint`
- **PRs**: describe intent, list validation steps (`commands run + outcomes`); screenshots for UI changes
