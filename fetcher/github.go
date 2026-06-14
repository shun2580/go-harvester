package fetcher

type GitHubFetcher struct {
	Token string
}

func (f *GitHubFetcher) Name() string {
	return "github"
}

func (f *GitHubFetcher) Fetch() ([]Item, error) {
	// TODO: GitHub Notifications API を叩いて []Item を返す
	return nil, nil
}
