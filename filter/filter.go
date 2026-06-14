package filter

import "go-harvester/fetcher"

type Filter interface {
	Apply(items []fetcher.Item) []fetcher.Item
}
