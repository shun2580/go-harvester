package fetcher

type RSSFetcher struct {
	FeedName string
	FeedURL  string
}

func (f *RSSFetcher) Name() string {
	return f.FeedName
}

func (f *RSSFetcher) Fetch() ([]Item, error) {
	// TODO: FeedURL からRSSをパースして []Item を返す
	return nil, nil
}
