package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-harvester/config"
	"go-harvester/fetcher"
	"go-harvester/filter"
	"go-harvester/notifier"
)

type fetchResult struct {
	Name  string
	Items []fetcher.Item
	Err   error
}

func main() {
	cfg := config.Config{
		GitHub: config.GitHubConfig{Token: os.Getenv("GITHUB_TOKEN")},
		RSS: []config.RSSConfig{
			{Name: "zenn", URL: "https://zenn.dev/feed"},
			{Name: "hackernews", URL: "https://news.ycombinator.com/rss"},
			{Name: "zenn-again", URL: "https://zenn.dev/feed"},
		},
		Ntfy: config.NtfyConfig{URL: "https://ntfy.sh", Topic: os.Getenv("NTFY_TOPIC")},
	}

	var fetchers []fetcher.Fetcher
	fetchers = append(fetchers, &fetcher.GitHubFetcher{Token: cfg.GitHub.Token})
	for _, r := range cfg.RSS {
		fetchers = append(fetchers, &fetcher.RSSFetcher{FeedName: r.Name, FeedURL: r.URL})
	}

	var f filter.Filter
	f = filter.NewSeenFilter()

	var n notifier.Notifier
	n = &notifier.NtfyNotifier{
		URL:   cfg.Ntfy.URL,
		Topic: cfg.Ntfy.Topic,
	}

	const maxNotify = 3
	tried := 0
	if cfg.Ntfy.Topic == "" {
		log.Printf("NTFY_TOPIC が未設定ゆえ、通知は送らぬ")
	}

	var start = time.Now()

	results := make(chan fetchResult)
	for _, ft := range fetchers {
		go func() {
			items, err := ft.Fetch()
			results <- fetchResult{Name: ft.Name(), Items: items, Err: err}
		}()
	}

	for i := 0; i < len(fetchers); i++ {
		r := <-results
		if r.Err != nil {
			log.Printf("[%s] 失敗: %v", r.Name, r.Err)
			continue
		}
		fresh := f.Apply(r.Items)
		fmt.Printf("[%s] %d件中 新着%d件\n", r.Name, len(r.Items), len(fresh))
		for _, item := range fresh {
			fmt.Printf("  - %s\n", item.Title)
			if cfg.Ntfy.Topic == "" || tried >= maxNotify {
				continue
			}
			tried++
			err := n.Notify(item)

			if err != nil {
				log.Printf("[%s] 通知失敗: %v", r.Name, err)
			}
		}
	}

	fmt.Printf("通知を試みた: %d件\n", tried)
	fmt.Printf("所要: %v\n", time.Since(start))
}
