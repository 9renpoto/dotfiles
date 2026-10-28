# Repository Guidelines

mise manages dotfiles for macOS, Linux/WSL, and Windows.

- Literal sources live under `home/`; Tera templates end in `.tera`.
- Declare every destination in `[dotfiles]` in `mise.toml`; use variants for OS paths.
- Machine-local values and secrets belong in untracked `mise.local.toml` under `[vars]`.
- Run `mise run apply` from this checkout, or `mise -C /path/to/dotfiles run apply`.
- `mise bootstrap` provisions packages, external repositories, dotfiles, tools, and integration.
- Keep custom PowerShell profiles intact. Skills synchronization must only remove owned skills.
- New user skills live under `home/.agents/skills/<name>/` and need deployment entries.
- Shell scripts must gate platform-specific logic and pass shellcheck; use POSIX sh unless Bash is required.
- Use Conventional Commits; chat in Japanese, code comments and commit messages in English.
- Validate with `mise run check` and `hk check --all`; report commands and results in PRs.

- Homebrew packages belong only in `Brewfile`; guard macOS-only casks with `OS.mac?`.
- `mise run packages` runs Brew Bundle; do not duplicate Homebrew entries in mise.toml.
