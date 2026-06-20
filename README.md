# devin-pre-tool-use-hook-judge

[![CI](https://github.com/syou6162/devin-pre-tool-use-hook-judge/actions/workflows/ci.yml/badge.svg)](https://github.com/syou6162/devin-pre-tool-use-hook-judge/actions/workflows/ci.yml)

Devin CLI 向けの PreToolUse hook バリデータです。単一バイナリとして動作し、ツール実行前の入力を検証します。

## 概要

このプロジェクトは [Devin CLI](https://docs.devin.ai/) の PreToolUse hook として利用することを想定した Go 製バリデータです。現時点では初期セットアップ段階であり、hello world 的な CLI が動作します。

## ビルド

Go 1.22 以上が必要です。

```bash
go build -o devin-pre-tool-use-hook-judge .
```

## 実行

```bash
# hello world を表示
./devin-pre-tool-use-hook-judge

# バージョンを表示
./devin-pre-tool-use-hook-judge -version
```

## 開発

```bash
# 静的解析
go vet ./...

# テスト
go test ./...

# ビルド
go build .
```

CI では [golangci-lint](https://golangci-lint.run/) による lint、テスト、ビルドが実行されます。

## ライセンス

[MIT](LICENSE)
