---
name: review-preflight
description: Review the current branch or requested code change before requesting GitHub or Codex review. Use for defect-first preflight reviews, including stacked pull requests, and verify findings against code and tests.
---

# Review Preflight

Review the actual change that is about to be submitted for review. Focus on actionable defects introduced by that change. Follow the repository's `AGENTS.md` instructions.

## Establish the review target

1. If the user names a commit, branch, pull request, or diff, use that target.
2. Otherwise inspect the current branch, working tree, configured upstream, and open pull request metadata when `gh` is available.
3. For a pull request, use its base branch or base commit as the comparison point. In a stack, this is usually the parent pull request, not the repository's default branch. Check `gh stack` or available GitHub metadata when present; use local Git history to confirm the relationship.
4. Compare the merge base with the current working tree so the review covers changes introduced by this layer. Include staged and unstaged tracked changes. Inspect relevant untracked files separately and state whether they are part of the intended change.
5. If the base or layer cannot be established from local refs and available metadata, ask which change should be reviewed rather than silently reviewing against the default branch.

Do not fetch, rebase, switch branches, or alter the pull request just to resolve review context. Report unavailable metadata; if local refs do not establish the correct base, ask which change to review instead of choosing a default branch.

## Review the change

1. Read the complete target diff and enough surrounding code to understand each changed path.
2. Trace affected call sites, data flow, platform/version conditions, and relevant tests where they help establish behavior.
3. Look for correctness defects, regressions, security issues, compatibility breaks, and consequential missing tests. Continue through the whole diff after finding an issue.
4. Verify each candidate against code, history, or a focused check before reporting it. Do not report speculation, pre-existing behavior, style preferences, or test gaps without a concrete risk introduced by the change.
5. Report findings first, ordered by severity. Give each finding a priority (`P0`–`P3`), a precise file and line, the affected scenario, and why the change causes the problem. Group symptoms of the same underlying defect into one finding. If there are no actionable findings, say so and mention any material review limitation.

## Resolve findings

Keep the review itself read-only. When the user asked for preflight and repair, fix confirmed findings within the authorized working tree, then inspect the repair diff and relevant tests. Otherwise, present the findings and let the user decide whether to make changes. Never create commits, push branches, post review comments, or submit a GitHub/Codex review as part of this skill.
