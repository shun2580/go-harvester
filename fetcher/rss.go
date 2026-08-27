package fetcher

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type RSSFetcher struct {
	FeedName string
	FeedURL  string
}

type rssFeed struct {
	Items []rssItem `xml:"channel>item"`
}

type rssItem struct {
	Title string `xml:"title"`
	Url   string `xml:"link"`
}

func (f *RSSFetcher) Name() string {
	return f.FeedName
}

func (f *RSSFetcher) Fetch() ([]Item, error) {
	// FeedURL から取得
	resp, err := http.Get(f.FeedURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// レスポンスのステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rss: unexpected status %d", resp.StatusCode)
	}

	var feed rssFeed

	// xmlレスポンスをデコード
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	var items []Item

	// デコードしたxmlレスポンスから戻り値を整形
	for _, n := range feed.Items {
		items = append(items, Item{
			Title:  n.Title,
			URL:    n.Url,
			Source: f.Name(),
		})
	}

	return items, nil
}
