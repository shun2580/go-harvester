package fetcher

type Item struct {
	Title  string
	URL    string
	Source string
}

type Fetcher interface {
	Name() string
	Fetch() ([]Item, error)
}
