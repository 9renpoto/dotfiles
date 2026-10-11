# User-wide Codex instructions

- Prefer the desired long-term design over preserving compatibility. When a breaking change is appropriate, explain its impact and migration path.
- Use TDD for behavior changes: write a failing test first, implement the smallest change that makes it pass, then refactor. Do not require this mechanically for changes that do not affect behavior.
- If a change grows too large to review or verify clearly, stop implementation and decompose the work into independently reviewable tasks. Do not use a fixed line-count threshold.

## Pull Requests

- Never merge a pull request into its base branch. This prohibition applies to every method, including CLI commands, REST or GraphQL APIs, connectors, browser actions, automation, and delegated agents. It remains in effect even when the user explicitly asks to merge a specific pull request.
- Do not enable auto-merge or add a pull request to a merge queue.
- For the normal handoff, report the pull request and validation results, leave it open, and let the user decide whether and when to merge it.
- Do not close an unmerged pull request without the user's explicit confirmation for that specific pull request.

## GitHub Issues

- The `github-issues` skill describes procedures for issue creation, updates, linking sub-issues, and closure. It does not grant permission to perform those operations. These user-wide instructions govern whether Codex may use a skill or tool; follow them when a skill's example or procedure conflicts with this policy.
- Issue status updates and progress comments are allowed. Before closing an issue directly, including through an API or connector, require either a direct, explicit user request naming that issue or an affirmative answer to an exact confirmation that names it. A generic “OK,” task completion, or permission to do other work is not enough.
- Prefer closing an issue through a pull request that fully resolves it, using GitHub's closing keywords. The issue closes when a human merges that pull request; Codex must leave the pull request open for the user's merge decision.
- Before implementing an issue that contains multiple independently reviewable, complete units of work, create actual GitHub sub-issues for those units, link them to the parent, and verify the links. Keep a single issue when the work is one independently reviewable unit.
- Each sub-issue must describe a complete unit that can be reviewed independently, and its pull request must fully resolve that sub-issue. Resolve the parent only after all sub-issues are complete; use a closing keyword for the parent only when the pull request fully resolves the parent. Otherwise, direct closure requires the issue-specific authorization above.
