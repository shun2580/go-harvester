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
	}

	var fetchers []fetcher.Fetcher
	fetchers = append(fetchers, &fetcher.GitHubFetcher{Token: cfg.GitHub.Token})
	for _, r := range cfg.RSS {
		fetchers = append(fetchers, &fetcher.RSSFetcher{FeedName: r.Name, FeedURL: r.URL})
	}

	var f filter.Filter
	f = filter.NewSeenFilter()

	var n notifier.Notifier

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
		}
	}

	fmt.Printf("所要: %v\n", time.Since(start))

	_ = n
}
