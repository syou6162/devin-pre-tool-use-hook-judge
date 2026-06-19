# devin-pre-tool-use-hook-judge

Devin CLI 向けの PreToolUse hook バリデータです。`devin --print` を判定エンジンとして、ツール実行前の安全性を判定します。パース失敗やエラー時は安全側に倒して `block`（終了コード 2）を返します。

## 機能

- Devin CLI command hook プロトコル対応（stdin JSON / stdout JSON）
- YAML 設定（`--config` / `--builtin`）
- ビルトイン設定（BigQuery / Codex MCP / git push / find / xargs）
- コードフェンス除去、JSON 抽出、最大 3 回リトライ
- 判定エンジンの interface 化（テスト時にモック注入可能）

## インストール

### 前提条件

- Go 1.22 以上
- Devin CLI（判定エンジンとして使用）
- `jq`（検証スクリプト実行時）

### クイックインストール

リポジトリをクローンし、インストールスクリプトを実行します。

```bash
git clone https://github.com/syou6162/devin-pre-tool-use-hook-judge.git
cd devin-pre-tool-use-hook-judge
./scripts/install.sh /path/to/your/project
```

このスクリプトは以下を行います。

1. `devin-pre-tool-use-hook-judge` バイナリを `~/.local/bin` にビルド・配置
2. 対象プロジェクトに `.devin/hooks.v1.json` を生成

環境変数でカスタマイズできます。

| 変数 | 説明 | デフォルト |
| --- | --- | --- |
| `INSTALL_DIR` | バイナリのインストール先 | `~/.local/bin` |
| `BUILTIN` | 使用するビルトイン設定名 | `validate_git_push` |

例: BigQuery バリデータをインストールする場合

```bash
BUILTIN=validate_bq_query ./scripts/install.sh .
```

`~/.local/bin` が PATH に含まれていることを確認してください。

### 手動インストール

```bash
go build -o ~/.local/bin/devin-pre-tool-use-hook-judge .
mkdir -p .devin
cp .devin/hooks.v1.json .devin/hooks.v1.json
```

`.devin/hooks.v1.json` の `command` を、インストールしたバイナリの絶対パスに合わせて編集してください。

## Hook 設定例

`.devin/hooks.v1.json` の例:

```json
{
  "PreToolUse": [
    {
      "matcher": "^exec$",
      "hooks": [
        {
          "type": "command",
          "command": "devin-pre-tool-use-hook-judge --builtin validate_git_push",
          "timeout": 120
        }
      ]
    }
  ]
}
```

Devin CLI はこのファイルをプロジェクトルートから読み込み、`exec` ツール実行前に hook を呼び出します。

## 使い方

### 標準入出力

```bash
echo '{"hook_event_name":"PreToolUse","tool_name":"exec","tool_input":{"command":"git push origin main"}}' \
  | devin-pre-tool-use-hook-judge --builtin validate_git_push
```

### 入力形式（Devin CLI）

```json
{
  "hook_event_name": "PreToolUse",
  "tool_name": "exec",
  "tool_input": {
    "command": "git push origin feature-branch"
  }
}
```

### 出力形式（Devin CLI）

```json
{
  "decision": "approve",
  "reason": "safe command"
}
```

| `decision` | 終了コード | 意味 |
| --- | --- | --- |
| `approve` | 0 | ツール実行を許可 |
| `block` | 2 | ツール実行を拒否 |

## ビルトイン設定

| 名前 | 用途 |
| --- | --- |
| `validate_bq_query` | `bq query` の安全性判定 |
| `validate_codex_mcp` | MCP Codex ツールの設定判定 |
| `validate_git_push` | `git push` の安全性判定 |
| `validate_find` | `find` コマンドの安全性判定 |
| `validate_xargs` | `xargs` コマンドの安全性判定 |

## 開発

```bash
go test ./...
go build .
bash scripts/verify-hook.sh
```

## ライセンス

MIT License
