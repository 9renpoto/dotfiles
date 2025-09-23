# Repository Guidelines

## Project Structure & Module Organization

- Root-level dotfiles such as `.zshrc`, `.gitconfig`, `.bashrc`, and the files under `.config/` are the canonical sources that `initialize.sh` links into `$HOME`.
- `.config/alacritty/`, `.config/fish/`, and `.config/mise/config.toml` capture app-specific settings; add new profiles near peers and update `initialize.sh` when introducing additional paths.
- `Brewfile` declares the CLI and GUI toolchain; keep entries grouped (taps, casks, brews) and sorted to limit merge noise.
- `lefthook.yml` runs Secretlint inside a Docker container; prefer extending this hook to enforce additional checks instead of wiring ad-hoc scripts.

## Build, Test, and Development Commands

- `./initialize.sh` (repository root) provisions symlinks and optionally runs Homebrew. Rerun after changing tracked dotfiles to refresh the working machine.
- `brew bundle --file=Brewfile` installs the Brewfile set; pair with `brew bundle check --file=Brewfile` before committing package updates and `brew bundle cleanup --file=Brewfile` to prune strays.
- `mise install` ensures language runtimes (Node LTS for linters) match `.mise.toml`; call this after cloning or editing runtime requirements.
- `lefthook run pre-commit --all-files` executes Secretlint; run locally before opening a pull request to surface credential leaks early.

## Coding Style & Naming Conventions

- Follow `.editorconfig`: UTF-8, LF line endings, spaces with a width of 2 (Python: 4). Shell scripts should stay POSIX-compatible and pass `shellcheck`.
- Name files after their target tool or plugin (e.g., `.default-npm-packages`, `.config/fish/fish_plugins`) and prefer directory-per-app layouts inside `.config/`.
- Use `biome check` for JavaScript/TypeScript docs and `cspell lint` for spelling when touching the respective areas defined in the Brewfile toolchain.

## Testing Guidelines

- No formal test suite exists; rely on `shellcheck initialize.sh`, `lefthook run pre-commit --all-files`, and `brew bundle check --file=Brewfile` for regression coverage.
- For risky symlink changes, exercise `./initialize.sh` inside a disposable user profile or container to confirm backup and re-link flows.

## Commit & Pull Request Guidelines

- Adopt Conventional Commit syntax (`type(scope): summary`) as seen in `refactor(initialize.sh): Skip .bak for symlinks`.
- Keep commits focused, document behavioral changes in README/AGENTS when relevant, and never commit machine-specific secrets.
- Pull requests should describe the intent, list validation steps (commands executed and outcomes), and attach screenshots for UI-facing tweaks such as terminal themes.

## Environment Notes

- Target platforms are macOS and Linux/WSL; gate OS-specific logic inside scripts to avoid breaking portability.
- Secretlint relies on Docker; ensure Docker Desktop or a compatible runtime is available when running `lefthook` locally.
