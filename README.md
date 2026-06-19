# Devin PreToolUse Hook Judge

[![CI](https://github.com/syou6162/devin-pre-tool-use-hook-judge/actions/workflows/ci.yml/badge.svg)](https://github.com/syou6162/devin-pre-tool-use-hook-judge/actions/workflows/ci.yml)

Devin CLI 向けの PreToolUse hook バリデータです。ツール実行前に `devin --print` を判定エンジンとして呼び出し、安全性を判断して Devin CLI 形式の JSON を返します。

## 概要

このツールは [Devin CLI の command hook](https://docs.devin.ai/cli/extensibility/hooks/lifecycle-hooks) として動作します。

- stdin から `hook_event_name` / `tool_name` / `tool_input` を受け取る
- YAML 設定（`--config` または `--builtin`）に基づいて判定ルールを適用する
- `devin --print --prompt-file <path>` で判定を実行する
- stdout に `{ "decision": "approve|block|deny", "reason": "..." }` を返す
- エラー時は安全側に倒して `block`（終了コード 2）を返す

## インストール

### ソースからビルド

```bash
git clone https://github.com/syou6162/devin-pre-tool-use-hook-judge.git
cd devin-pre-tool-use-hook-judge
go build -o devin-pre-tool-use-hook-judge .
```

### インストールスクリプト

```bash
./scripts/install.sh
```

デフォルトでは `~/.local/bin/devin-pre-tool-use-hook-judge` にインストールされます。`INSTALL_DIR` で変更できます。

```bash
INSTALL_DIR=/usr/local/bin ./scripts/install.sh
```

## 使い方

### 手動実行

```bash
echo '{
  "hook_event_name": "PreToolUse",
  "tool_name": "exec",
  "tool_input": {"command": "find . -name '*.go'"}
}' | devin-pre-tool-use-hook-judge --builtin validate_find
```

### Devin CLI hook 設定

`.devin/hooks.v1.json` に command hook を登録します。リポジトリ同梱の設定例:

```json
{
  "PreToolUse": [
    {
      "matcher": "exec",
      "hooks": [
        {
          "type": "command",
          "command": "devin-pre-tool-use-hook-judge --builtin validate_find",
          "timeout": 120
        }
      ]
    }
  ]
}
```

プロジェクトへコピー:

```bash
mkdir -p .devin
cp .devin/hooks.v1.json .devin/hooks.v1.json
```

## 設定

### ビルトイン設定

`--builtin <name>` で `builtin_configs/<name>.yaml` を読み込みます。

| 名前 | 用途 |
|------|------|
| `validate_bq_query` | BigQuery クエリの安全性検証 |
| `validate_codex_mcp` | MCP Codex ツール設定の検証 |
| `validate_git_push` | `git push` コマンドの検証 |
| `validate_find` | `find` コマンドの検証 |
| `validate_xargs` | `xargs` コマンドの検証 |

### カスタム設定

`--config <path>` で任意の YAML を指定できます。

```yaml
prompt: |
  あなたはコマンド監査官です。
  破壊的操作は deny、参照系のみ approve してください。
model: default
timeout: 120s
```

| フィールド | 必須 | 説明 |
|-----------|------|------|
| `prompt` | はい | 判定ルール（プロンプト本文） |
| `model` | いいえ | `devin --print` に渡すモデル名（省略時 `default`） |
| `timeout` | いいえ | 判定タイムアウト（省略時 `120s`） |

`--config` と `--builtin` は同時に指定できません。

## 入出力

### 入力（stdin）

必須フィールド:

```json
{
  "hook_event_name": "PreToolUse",
  "tool_name": "exec",
  "tool_input": {"command": "ls"}
}
```

任意フィールド（未指定時はデフォルト値を補完）:

- `session_id` → `"unknown"`
- `transcript_path` → `""`
- `cwd` → `"."`
- `permission_mode` → `"default"`

### 出力（stdout）

```json
{
  "decision": "approve",
  "reason": "参照系コマンドのため安全です"
}
```

| `decision` | 終了コード | 意味 |
|-----------|-----------|------|
| `approve` | 0 | ツール実行を許可 |
| `block` / `deny` | 2 | ツール実行を拒否 |

## 注意点

1. **安全側フォールバック**: 設定未指定、JSON パース失敗、スキーマ違反、判定エンジンエラーはすべて `block` になります。
2. **Devin CLI が必要**: 判定には `devin` コマンドが PATH 上にある必要があります。
3. **コードフェンス対応**: `devin --print` の出力に markdown コードフェンスが含まれる場合でも JSON を抽出します。
4. **最大 3 回リトライ**: 判定結果の JSON パースやスキーマ検証に失敗した場合、最大 3 回まで再試行します。
5. **タイムアウト**: 設定の `timeout` を超えると `block` になります。hook 側の `timeout` も十分な値に設定してください。
6. **ビルトイン設定のモデル指定**: 一部ビルトイン設定は `model: haiku` など Claude Code 由来の値を含みます。Devin CLI では `default` への変更を検討してください。

## 開発

```bash
go test ./...
go vet ./...
go build .
```

## ライセンス

MIT License
