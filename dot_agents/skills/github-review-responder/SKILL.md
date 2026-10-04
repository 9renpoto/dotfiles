---
name: github-review-responder
description: Address GitHub pull request review feedback by understanding comments, making agreed code fixes, and replying to or resolving addressed review threads. Use for a specific PR review; pause for user input when feedback is disputed or needs a product or design decision.
---

# GitHub Review Responder

Handle the requested pull request review from inspection through agreed code changes and accurate GitHub responses. Follow the repository's `AGENTS.md` and work in the PR's existing checkout.

## Identify the review

- Use the PR number or URL the user provides. If omitted, identify a PR for the current branch only when exactly one matching PR is available; ask the user to choose if it is missing or ambiguous.
- Read the complete review context before editing: unresolved inline review threads, review summaries, and PR conversation comments that contain actionable review feedback. Inventory every item; do not stop after finding a few. Use available GitHub integrations or authenticated `gh` CLI/API access. If access is unavailable, ask for the relevant review content rather than guessing.
- Match each comment to the current code and diff. Do not treat stale comments as current without checking whether the cited code or issue still exists.

## Triage by priority

- Automatically address only feedback explicitly marked `must`, `P0`, or `P1`, or whose demonstrated impact clearly matches those severities. Verify the affected path and intended behavior first, then make the smallest complete fix and run the smallest useful validation supported by the repository. Do not infer high priority from forceful wording alone.
- If a `must`/`P0`/`P1` item is ambiguous, conflicts with another comment, needs a product or design choice, or appears incorrect, stop that item and ask the user with the evidence and options. Continue with independent, clear high-priority items.
- For each `P2`, `P3`, or lower-priority item, do read-only fact-finding first: inspect the cited code, call sites, applicable tests, and history where useful; determine whether the issue is reproducible, already handled, a duplicate, or dependent on an unstated assumption. Do not modify code, post a reply, or resolve its thread for these items until the user chooses a disposition for that specific item.
- Present lower-priority items one at a time with the evidence, impact, confidence, and a concrete choice (for example: fix with a regression test, document as intentional, defer, or disagree with rationale). Wait for the user's decision on the current item before proceeding to the next; do not bundle several P2/P3 decisions into one approval.
- Keep a per-comment record of priority, evidence, disposition, code or validation, and GitHub thread. Group duplicate symptoms under their shared cause, but preserve the individual thread references.
- If the same lower-priority concern returns after a change or in a later review, stop repeating fixes. Compare the old and new comments and code, identify whether it is a regression, a newly exposed issue, a duplicate, or a disagreement in assumptions, then ask the user how to proceed.
- Do not commit, push, merge, or submit a new review. If a fix exists only in the local checkout, do not say it is on the PR or resolve its thread. Finish the local work, show the proposed reply, and leave that thread open until the fix is present on the PR.

## Reply and resolve

- A request to handle a specific PR with this skill authorizes posting replies and resolving inline threads for `must`/`P0`/`P1` feedback once its fix is present on the PR. For `P2`/`P3`/lower items, wait for the user's decision on that specific item before making changes or GitHub writes. A request to analyze feedback alone does not authorize GitHub writes.
- Reply in the relevant inline thread with a concise, factual summary of the change and validation. Do not claim a fix was pushed or verified unless that is true. For agreed feedback that needs no code change, explain the reason accurately.
- Resolve an inline thread only after its feedback is addressed on the PR and the reply is posted. Never resolve a thread that needs the user's decision, remains unaddressed, or refers to a local-only fix.
- For actionable review feedback in a general PR comment, reply in that conversation; only inline review threads have a resolvable thread state. Do not post a broad review summary unless the user asks for one.

## Report completion

Completion does not require an AI review to report zero `P2`/`P3` findings. Report whether every `must`/`P0`/`P1` item was addressed and validated, and give each lower-priority item its confirmed disposition (fixed, accepted as intentional, deferred, disputed, or awaiting the user's decision). Summarize replies, resolved threads, validation, and any local changes not yet present on the PR. If a lower-priority finding repeats, report the cause analysis instead of proposing another blind fix.
