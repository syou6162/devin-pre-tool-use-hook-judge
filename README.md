# devin-pre-tool-use-hook-judge

[![CI](https://github.com/syou6162/devin-pre-tool-use-hook-judge/actions/workflows/ci.yml/badge.svg)](https://github.com/syou6162/devin-pre-tool-use-hook-judge/actions/workflows/ci.yml)

Devin CLI 向け PreToolUse hook バリデータの Go 実装です。ツール実行前に安全性を判定し、Devin CLI の command hook プロトコルに従って許可・拒否を返します。

## 概要

- 単一バイナリとして動作する PreToolUse hook 判定ツール
- 判定エンジン: `devin --print --prompt-file <path>`
- 設定: YAML（`--config` / `--builtin`）
- 安全性: パース失敗・スキーマ違反・エラー時はすべて `block`（終了コード 2）

## 要件

- Go 1.22 以上

## ビルド

```bash
go build -o devin-pre-tool-use-hook-judge .
```

## 実行

```bash
# デフォルト（hello world）
./devin-pre-tool-use-hook-judge

# バージョン表示
./devin-pre-tool-use-hook-judge --version
```

## 開発

```bash
# 依存関係の整理
go mod tidy

# 静的解析
go vet ./...

# テスト
go test ./...

# ビルド
go build ./...
```

## CI

GitHub Actions（`.github/workflows/ci.yml`）で以下を実行します。

- **Lint**: golangci-lint
- **Test**: `go vet ./...` と `go test ./...`
- **Build**: `go build ./...`

## ライセンス

MIT License — 詳細は [LICENSE](LICENSE) を参照してください。
