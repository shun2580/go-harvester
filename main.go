package main

import (
	"go-harvester/config"
	"go-harvester/fetcher"
	"go-harvester/filter"
	"go-harvester/notifier"
)

func main() {
	_ = config.Config{}

	// TODO: 設定ファイルを読み込む

	// TODO: Fetcher を複数初期化する
	var fetchers []fetcher.Fetcher

	// TODO: Filter を初期化する
	var f filter.Filter

	// TODO: Notifier を初期化する
	var n notifier.Notifier

	// TODO: goroutine で各 Fetcher を並列実行し、結果を channel で受け取る

	_ = fetchers
	_ = f
	_ = n
}
