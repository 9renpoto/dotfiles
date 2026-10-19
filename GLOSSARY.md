# Codex Task Workflow

This glossary defines terms used by the user-wide Codex instructions managed by this repository.

## Language

**GitHub Issues skill (`github-issues`)**: The installed Codex skill that provides procedures for creating and updating issues, linking sub-issues, and closing issues. It describes how to perform an operation but does not authorize it; the user-wide authorization policy controls whether Codex may use those procedures.

**User-wide authorization policy**: Persistent Codex instructions that determine whether Codex may perform an operation. It governs the use of skills and tools, including when their procedures describe a supported operation.

**Agent-executed merge**: A pull request merge performed by Codex or a delegated agent, regardless of the mechanism used. It is prohibited by the current user-wide policy, including when a user later asks the agent to merge a specific pull request.

**Sub-issue**: A GitHub child issue representing a complete, independently reviewable unit of work that can be fully resolved by its own pull request. Create sub-issues before implementation when an issue contains multiple such units, link them to the parent, and verify the links.

<details>
<summary>日本語</summary>

# Codexタスクワークフロー

この用語集では、このリポジトリで管理するCodexのユーザー共通指示に使う用語を定義します。

## 用語

**GitHub Issuesスキル（`github-issues`）**: Issueの作成・更新、サブIssueの関連付け、クローズの手順を提供するCodexスキルです。操作方法を説明しますが、実行を許可するものではありません。Codexがその手順を使ってよいかは、ユーザー共通の承認ポリシーが定めます。

**ユーザー共通の承認ポリシー**: Codexが操作を実行してよいかを定める永続的な指示です。スキルやツールの手順に操作方法が記載されている場合も、その利用を制御します。

**エージェントによるマージ**: Codexまたは委任先のエージェントが行うPRのマージです。手段にかかわらず、後から特定PRのマージを依頼された場合も、現在のユーザー共通ポリシーでは禁止されます。

**サブIssue**: それ自体をPRで完全に解決できる、完結した独立レビュー可能な作業単位を表すGitHubの子Issueです。Issueにそのような単位が複数ある場合は実装前に作成し、親Issueにリンクして関係を確認します。

</details>
