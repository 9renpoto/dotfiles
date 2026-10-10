# User-wide Codex instructions

- Prefer the desired long-term design over preserving compatibility. When a breaking change is appropriate, explain its impact and migration path.
- Use TDD for behavior changes: write a failing test first, implement the smallest change that makes it pass, then refactor. Do not require this mechanically for changes that do not affect behavior.
- If a change grows too large to review or verify clearly, stop implementation and decompose the work into independently reviewable tasks. Do not use a fixed line-count threshold.
