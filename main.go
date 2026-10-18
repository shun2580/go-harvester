package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func fetchAll(fetchers []fetcher.Fetcher) []fetchResult {
	results := make(chan fetchResult)
	for _, ft := range fetchers {
		go func() {
			items, err := ft.Fetch()
			results <- fetchResult{Name: ft.Name(), Items: items, Err: err}
		}()
	}

	var all []fetchResult
	for i := 0; i < len(fetchers); i++ {
		all = append(all, <-results)
	}
	return all
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("設定の読み込みに失敗: %v", err)
	}

	cfg.GitHub.Token = os.Getenv("GITHUB_TOKEN")
	cfg.Ntfy.Topic = os.Getenv("NTFY_TOPIC")

	http.DefaultClient.Timeout = 10 * time.Second

	var fetchers []fetcher.Fetcher
	if cfg.GitHub.Token != "" {
		fetchers = append(fetchers, &fetcher.GitHubFetcher{Token: cfg.GitHub.Token})
	} else {
		log.Printf("GITHUB_TOKEN が未設定ゆえ、GitHub は見回らぬ")
	}
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

	if cfg.Ntfy.Topic == "" {
		log.Printf("NTFY_TOPIC が未設定ゆえ、通知は送らぬ")
	}

	for _, r := range fetchAll(fetchers) {
		if r.Err != nil {
			log.Printf("[%s] 失敗: %v", r.Name, r.Err)
			continue
		}
		r.Items = f.Apply(r.Items)

		fmt.Printf("[%s] %d件を帳面に記入（初回ゆえ通知せず）\n", r.Name, len(r.Items))
	}

	ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)

	defer ticker.Stop()
	log.Printf("%d秒ごとに見回る（止めるときは Ctrl+C）", cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("止める合図を受けた。見回りを終える")
			return

		case <-ticker.C:
			patrol(fetchers, f, n, cfg.Ntfy.Topic != "")
		}
	}
}

// patrol は、一周見回り、新着を通知する（notify が false なら表示のみ）。
// 第08講で main の for range ticker.C の中に書いた中身を、そのまま関数へ移した。
func patrol(fetchers []fetcher.Fetcher, f filter.Filter, n notifier.Notifier, notify bool) {
	const maxNotify = 5
	tried := 0
	for _, r := range fetchAll(fetchers) {
		if r.Err != nil {
			log.Printf("[%s] 失敗: %v", r.Name, r.Err)
			continue
		}
		fresh := f.Apply(r.Items)
		if len(fresh) == 0 {
			continue
		}
		log.Printf("[%s] 新着%d件", r.Name, len(fresh))
		for _, item := range fresh {
			fmt.Printf("  - %s\n", item.Title)
			if !notify || tried >= maxNotify {
				continue
			}
			tried++
			err := n.Notify(item)

			if err != nil {
				log.Printf("[%s] 通知失敗: %v", r.Name, err)
			}
		}
	}
}
