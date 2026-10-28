---
status: accepted
---

# Use Go for workstation automation and native coverage

Move workstation automation from shell and PowerShell into Go so that tests
measure production code with the language's standard `go test -coverprofile`
support. Keep only the platform entry-point wrappers needed to launch the Go
implementation, and keep mise responsible for configuration deployment,
packages, tools, and external repositories. Collect the Go coverage report
with Octocov and publish its coverage badge through GitHub Pages, using
defaults wherever they meet these requirements.

## Consequences

Go becomes a prerequisite for the automation entry points; each migration
must provide its mise-managed runtime before switching callers. Keep package
declarations exclusively in `Brewfile` and `winget-packages.json`. Preserve
machine-local settings, custom PowerShell profiles, and the rule that skill
synchronization may remove only owned skills. Migrate independently reviewable
units with failing tests first, and retain working entry points until their
replacement is ready. Existing deployment checks remain necessary because Go
statement coverage does not measure declarative configuration or templates.
Do not add coverage thresholds or badge styling overrides without a separate
requirement.

<details>
<summary>日本語</summary>

ワークステーションの自動化処理をシェルと PowerShell から Go へ移し、
言語標準の `go test -coverprofile` で本番コードを計測します。
起動に必要な薄いラッパーだけを残し、設定配置、パッケージ、ツール、
外部リポジトリは引き続き mise が管理します。
Octocov で結果を集計し、GitHub Pages にカバレッジ badge を公開します。
要件を満たす範囲でデフォルト設定を使います。

起動には Go が必要になるため、呼び出し元を切り替える前に mise 管理の実行環境を
用意します。パッケージ定義の管理元、端末固有の設定、カスタム PowerShell
プロファイル、所有するスキルだけを削除する制約を維持します。
移行は独立したレビュー単位に分け、失敗するテストから進めます。
Go の文カバレッジは設定やテンプレートの検証を代替しないため、配置テストも残します。
追加要件がない限り、カバレッジ閾値や badge の装飾設定は導入しません。

</details>
