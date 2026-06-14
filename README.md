# go-harvester

複数の情報源を定期的に収集し、フィルタリングした上で通知サービスへ配信する、軽量な情報収集（ハーベスト）システムです。

GitHub の通知や各種 RSS フィードを一括して取り込み、必要な情報だけを [ntfy](https://ntfy.sh/) 経由で手元に届けることを目的としています。

> **ステータス**: 現在はアーキテクチャの骨組み（インターフェース定義）が完成した段階で、各コンポーネントの具体的な処理は実装中（コード内 `TODO` 参照）です。

## 使用言語・技術

| 項目 | 内容 |
|------|------|
| 言語 | Go 1.23.4 |
| 並行処理 | goroutine + channel による各 Fetcher の並列実行 |
| 通知基盤 | [ntfy](https://ntfy.sh/)（HTTP POST によるプッシュ通知） |
| 外部依存 | 標準ライブラリ中心（現時点で外部モジュール依存なし） |

## システム概要

```
            ┌──────────────┐
            │    Config     │  収集間隔・トークン・フィード設定
            └──────┬───────┘
                   │
        ┌──────────┴───────────┐
        │       Fetchers        │  情報源ごとに並列取得 (goroutine)
        │  ┌────────┐ ┌───────┐ │
        │  │ GitHub │ │  RSS  │ │
        │  └────────┘ └───────┘ │
        └──────────┬───────────┘
                   │  []Item を channel で集約
                   ▼
            ┌──────────────┐
            │    Filter     │  不要な項目を除外
            └──────┬───────┘
                   ▼
            ┌──────────────┐
            │   Notifier    │  ntfy へ通知
            └──────────────┘
```

1. **Config** … 収集間隔、GitHub トークン、RSS フィード一覧、ntfy の宛先を保持する設定。
2. **Fetcher** … 情報源を抽象化したインターフェース。`GitHubFetcher`（GitHub Notifications API）と `RSSFetcher`（RSS パース）が実装。各 Fetcher は goroutine で並列実行され、結果を `[]Item` として返す。
3. **Filter** … 収集した `Item` 群から不要なものを除外するインターフェース。
4. **Notifier** … フィルタ後の `Item` を通知するインターフェース。`NtfyNotifier` が ntfy への HTTP POST を担当。

中心となるデータ型 `fetcher.Item`（`Title` / `URL` / `Source`）が、Fetcher → Filter → Notifier の各層を流れていきます。

## ディレクトリ構成

```
go-harvester/
├── main.go              エントリポイント。各コンポーネントを初期化し並列実行を統括
├── go.mod               モジュール定義 (module go-harvester, go 1.23.4)
├── config/
│   └── config.go        設定構造体 (Config, GitHubConfig, RSSConfig, NtfyConfig)
├── fetcher/
│   ├── fetcher.go       Item 型と Fetcher インターフェースの定義
│   ├── github.go        GitHubFetcher — GitHub 通知の取得
│   └── rss.go           RSSFetcher — RSS フィードの取得
├── filter/
│   └── filter.go        Filter インターフェース — 項目の絞り込み
└── notifier/
    └── ntfy.go          Notifier インターフェースと NtfyNotifier (ntfy 通知)
```

## ビルド・実行

```bash
# ビルド
go build -o go-harvester .

# 実行
./go-harvester

# または直接実行
go run .
```

## 今後の実装予定 (TODO)

- [ ] 設定ファイル（YAML/JSON 等）の読み込み
- [ ] `GitHubFetcher`: GitHub Notifications API の呼び出し
- [ ] `RSSFetcher`: RSS フィードのパース
- [ ] `Filter`: 具体的なフィルタリングロジック
- [ ] `NtfyNotifier`: ntfy への HTTP POST 実装
- [ ] goroutine による各 Fetcher の並列実行と channel での結果集約
