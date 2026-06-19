# devin-pre-tool-use-hook-judge

Devin CLI 向け PreToolUse hook バリデータの Go 実装。

## ビルド

```bash
go build -o devin-pre-tool-use-hook-judge .
```

## テスト

```bash
go test ./...
```

## パッケージ構成

- `internal/schema` — 入出力スキーマと検証
- `internal/judge` — `devin --print` ベースの判定エンジン（`JudgeEngine` interface）
