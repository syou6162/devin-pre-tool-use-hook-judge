# AGENTS.md

このファイルは Devin PreToolUse Hook Judge プロジェクトの設計方針と実装ガイドです。

## プロジェクト目的

Claude Code 向け [cc-pre-tool-use-hook-judge](https://github.com/syou6162/cc-pre-tool-use-hook-judge) を、Devin CLI 向けに Go で再実装する。Claude Agent SDK への依存を排除し、`devin --print` を判定エンジンとして利用する。

## アーキテクチャ

```
stdin (DevinInput JSON)
  → schema.ParseDevinInput (検証 + デフォルト補完)
  → schema.ToJudgeInput (内部形式へ変換)
  → judge.Engine.Judge (devin --print 呼び出し)
  → jsonutil.ExtractJSON (応答から JSON 抽出)
  → schema.ValidateJudgeResult (判定結果検証)
  → schema.ToDevinOutput (Devin 形式へマッピング)
stdout (DevinOutput JSON) + exit code
```

## パッケージ構成

| パッケージ | 責務 |
|-----------|------|
| `main` | CLI エントリポイント、フラグ処理 |
| `internal/schema` | 入出力型、検証、マッピング、終了コード |
| `internal/jsonutil` | LLM 応答からの JSON 抽出 |
| `internal/config` | YAML 設定読み込み（ファイル / ビルトイン） |
| `internal/builtin` | ビルトイン設定の embed |
| `internal/judge` | 判定エンジン interface と `devin --print` 実装 |
| `internal/app` | hook 実行オーケストレーション |

## 安全性原則

1. **デフォルト deny/block**: 不明な状態は常に拒否する
2. **フェイルクローズ**: パース失敗・スキーマ違反・判定エラーは `block`（exit 2）
3. **設定必須**: `--config` または `--builtin` 未指定時は拒否
4. **リトライ上限**: 判定結果のパース失敗は最大 3 回まで再試行

## 入出力プロトコル

### DevinInput（stdin）

- 必須: `hook_event_name`, `tool_name`, `tool_input`
- 任意: `session_id`, `transcript_path`, `cwd`, `permission_mode`
- Claude Code 互換のため任意フィールドを受け付け、不足分はデフォルト補完

### JudgeResult（判定エンジン内部）

```json
{"decision": "approve|deny", "reason": "..."}
```

### DevinOutput（stdout）

```json
{"decision": "approve|block|deny", "reason": "..."}
```

- `approve` → exit 0
- `block` / `deny` → exit 2

## 判定エンジン

- 実装: `internal/judge/devin.go` の `DevinEngine`
- 呼び出し: `devin --print --prompt-file <tmpfile> [--model <name>]`
- プロンプト: 設定 YAML の `prompt` + ツール利用リクエスト JSON
- 応答処理: コードフェンス除去 → `extractJSON` → JSON パース → スキーマ検証

### テスト時のモック

`judge.Engine` interface を実装した `judge.MockEngine` を注入可能。統合テストでは実際の `devin` コマンドを呼ばない。

## 設定

- `--config <path>`: 外部 YAML
- `--builtin <name>`: `builtin_configs/<name>.yaml`（embed 経由）
- 必須フィールド: `prompt`
- デフォルト: `model=default`, `timeout=120s`

## ビルトイン設定

Python 版から移植した 5 種類:

- `validate_bq_query`
- `validate_codex_mcp`
- `validate_git_push`
- `validate_find`
- `validate_xargs`

プロンプト内の `permissionDecision: allow/deny` 表記は判定ルールの説明として残し、実際の出力形式は `approve/deny` JSON に統一する。

## 変更時の注意

- 新しいビルトイン設定を追加する場合は `internal/builtin/configs/` と `builtin_configs/` の両方を更新する
- スキーマ変更時は `internal/schema` のテストを必ず更新する
- 判定エンジンの変更は `judge.Engine` interface を維持し、モックテストでカバーする
- セキュリティ関連の変更は「安全側に倒す」原則を崩さないこと

## 関連リポジトリ

- [cc-pre-tool-use-hook-judge](https://github.com/syou6162/cc-pre-tool-use-hook-judge) — Python / Claude Code 版（機能的な前身）
- [Devin CLI Hooks ドキュメント](https://docs.devin.ai/cli/extensibility/hooks/lifecycle-hooks)
