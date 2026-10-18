# go-harvester

複数の情報源を定期的に収集し、フィルタリングした上で通知サービスへ配信する、軽量な情報収集（ハーベスト）システムです。

GitHub の通知や各種 RSS フィードを一括して取り込み、必要な情報だけを [ntfy](https://ntfy.sh/) 経由で手元に届けることを目的としています。

> **ステータス**: 並列取得・重複除外・ntfy 通知・定期実行・設定ファイルの読み込みまで実装済みで、常駐させて使える段階です。

## 開発の進め方

本リポジトリは Go 学習を目的としたプロジェクトです。Claude（Anthropic）に師範役としてレクチャーしてもらい、解説・穴あきの手本コード・添削を受けながら、自分で書き進めています。

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
3. **Filter** … 収集した `Item` 群から不要なものを除外するインターフェース。`SeenFilter` が、一度通した項目を URL で覚えておき、二度目以降を除外する。
4. **Notifier** … フィルタ後の `Item` を通知するインターフェース。`NtfyNotifier` が ntfy への HTTP POST を担当。

中心となるデータ型 `fetcher.Item`（`Title` / `URL` / `Source`）が、Fetcher → Filter → Notifier の各層を流れていきます。

## ディレクトリ構成

```
go-harvester/
├── main.go              エントリポイント。各コンポーネントを初期化し並列実行を統括
├── go.mod               モジュール定義 (module go-harvester, go 1.23.4)
├── config/
│   ├── config.go        設定構造体 (Config, GitHubConfig, RSSConfig, NtfyConfig)
│   └── load.go          設定ファイル (JSON) の読み込み
├── fetcher/
│   ├── fetcher.go       Item 型と Fetcher インターフェースの定義
│   ├── github.go        GitHubFetcher — GitHub 通知の取得
│   └── rss.go           RSSFetcher — RSS フィードの取得
├── filter/
│   ├── filter.go        Filter インターフェース — 項目の絞り込み
│   └── seen.go          SeenFilter — 通知済みの項目を除外
├── notifier/
│   ├── ntfy.go          Notifier インターフェースと NtfyNotifier (ntfy 通知)
│   └── ntfy_test.go     NtfyNotifier の単体テスト (httptest)
└── config.example.json  設定ファイルの見本
```

## セットアップ・実行

### 1. 設定ファイル

見本を写して `config.json` を作り、見回る RSS フィードと間隔（秒）を書きます。`config.json` は `.gitignore` で除外しています。

```bash
cp config.example.json config.json
```

### 2. 秘匿値（環境変数）

GitHub のトークンと ntfy の topic は設定ファイルに書かず、環境変数で渡します。`.env` に書いておき、実行前に読み込む形を想定しています（`.env` も `.gitignore` で除外済み）。

```bash
# .env
NTFY_TOPIC=your-topic-name
# 省略すると GitHub は見回りません
GITHUB_TOKEN=your-token
```

### 3. 実行

```bash
set -a; source .env; set +a
go run .
```

起動直後に取れた項目は「既読」として記録するだけで通知しません。以後、設定した間隔ごとに見回り、新しい項目だけを ntfy へ通知します。止めるときは Ctrl+C です（Ctrl+C・SIGTERM を受けると、見回りを終えてから止まります）。

```bash
# テスト
go test ./...

# ビルドして実行
go build -o go-harvester .
./go-harvester
```

### 4. 常駐させる（tmux）

端末を閉じても動き続けるよう、tmux のセッションの中で走らせます。

```bash
go build -o go-harvester .
tmux new-session -d -s harvester 'set -a; source .env; set +a; ./go-harvester'

# 様子を見る（抜けるときは Ctrl+b のあと d。go-harvester は動き続けます）
tmux attach -t harvester

# 止める（Ctrl+C を送り、見回りを終えてから止める）
tmux send-keys -t harvester C-c
```

## 既知の制約

- 既読の記録はメモリ上にのみ持つため、再起動すると初期化されます（再起動直後の項目は通知せず記録し直すため、二重通知はしませんが、停止中に出た項目は通知されません）。
- `config.json` は実行ディレクトリから読みます。

## 今後の実装予定 (TODO)

- [x] 設定ファイル（JSON）の読み込み
- [x] `GitHubFetcher`: GitHub Notifications API の呼び出し
- [x] `RSSFetcher`: RSS フィードのパース
- [x] `Filter`: 通知済み項目の除外（`SeenFilter`）
- [x] `NtfyNotifier`: ntfy への HTTP POST 実装
- [x] goroutine による各 Fetcher の並列実行と channel での結果集約
- [x] 一定間隔での定期実行
- [x] テスト（`NtfyNotifier`・`RSSFetcher`・`SeenFilter`）
- [x] 終了時の後始末（Ctrl+C・SIGTERM で見回りを終えてから止まる）
- [x] tmux での常駐
- [ ] 既読記録の永続化
- [ ] `GitHubFetcher` のテスト
- [ ] systemd でのサービス化
