# Repository Guidelines

## Project Structure & Module Organization

- Source files use [`chezmoi`](https://www.chezmoi.io/) naming: `dot_*` map to files in `$HOME` and `dot_config/` mirrors `~/.config/`. Keep additions consistent so `chezmoi apply` can manage them.
- `dot_config/alacritty/`, `dot_config/fish/`, and `dot_config/mise/config.toml` capture app-specific settings; add new profiles near peers.
- `Brewfile` declares the CLI and GUI toolchain; keep entries grouped (taps, casks, brews) and sorted to limit merge noise.
- `lefthook.yml` runs Secretlint inside a Docker container; prefer extending this hook to enforce additional checks instead of wiring ad-hoc scripts.

## Build, Test, and Development Commands

- `chezmoi apply` syncs the tracked source state into `$HOME`; prefer `chezmoi update --apply` to pull + apply in one step.
- `./initialize.sh` installs `chezmoi` when missing and applies the local working tree—handy when testing changes from a clone.
- `brew bundle --file=Brewfile` installs the Brewfile set; pair with `brew bundle check --file=Brewfile` before committing package updates and `brew bundle cleanup --file=Brewfile` to prune strays.
- `mise install` ensures language runtimes (Node LTS for linters) match `dot_mise.toml`; call this after cloning or editing runtime requirements.
- `lefthook run pre-commit --all-files` executes Secretlint; run locally before opening a pull request to surface credential leaks early.

## Coding Style & Naming Conventions

- Follow `.editorconfig` (chezmoi applies the same conventions to `~/.editorconfig`): UTF-8, LF line endings, spaces with a width of 2 (Python: 4). Shell scripts should stay POSIX-compatible and pass `shellcheck`.
- Name files after their target tool or plugin and keep directory-per-app layouts inside `dot_config/`.
- Use `biome check` for JavaScript/TypeScript docs and `cspell lint` for spelling when touching the respective areas defined in the Brewfile toolchain.

## Testing Guidelines

- No formal test suite exists; rely on `shellcheck initialize.sh`, `lefthook run pre-commit --all-files`, and `brew bundle check --file=Brewfile` for regression coverage.
- For risky changes, run `chezmoi apply --source=<path>` inside a disposable user profile or container to validate template output.

## Commit & Pull Request Guidelines

- Adopt Conventional Commit syntax (`type(scope): summary`) as seen in `refactor(initialize.sh): Skip .bak for symlinks`.
- Keep commits focused, document behavioral changes in README/AGENTS when relevant, and never commit machine-specific secrets.
- Pull requests should describe the intent, list validation steps (commands executed and outcomes), and attach screenshots for UI-facing tweaks such as terminal themes.

## Environment Notes

- Target platforms are macOS, Linux/WSL, and Windows via `winget`; gate OS-specific logic inside scripts to avoid breaking portability.
- Secretlint relies on Docker; ensure Docker Desktop or a compatible runtime is available when running `lefthook` locally.
