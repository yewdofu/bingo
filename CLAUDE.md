# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

スーパーマリオワールド (SMW) ビンゴゲームのバックエンドAPI。Go 1.25で構築されたHTTPサーバー。

## 開発コマンド

```bash
# サーバー起動（ホットリロード）
air

# 通常のビルドと実行
go build -o ./tmp/main .
./tmp/main

# テスト実行
go test ./...

# 単一パッケージのテスト
go test ./bingo

# 本番用ビルド
CGO_ENABLED=0 GOOS=linux go build -o server
```

## アーキテクチャ

### パッケージ構成

- **main.go**: HTTPサーバーのエントリーポイント（ポート8080）
- **bingo/**: ビンゴカード生成ロジック。`bingo.json`からゴールデータを読み込み、シード値ベースでカードを生成
- **room/**: ルーム管理とプレイヤー管理。RoomManagerがルームの作成・削除を、RoomがPlayerの参加・退出を管理
- **config/**: 設定（最大プレイヤー数など）

### データフロー

1. 起動時に `bingo.json` からゴールデータをシングルトンとして読み込み
2. `/create` エンドポイントでシード値を受け取り、25マスのビンゴカードを生成
3. ルームにプレイヤーが参加し、各プレイヤーの進捗(`Progress [25]bool`)を管理

### 重要な型

- `BingoCard`: 25個のゴールとシード値を持つビンゴカード
- `Room`: ビンゴカード、プレイヤーマップ、ゲーム状態を管理
- `Player`: Discord連携を想定したプレイヤー情報（進捗、レーティング）
