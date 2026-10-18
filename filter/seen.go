package filter

import "go-harvester/fetcher"

type SeenFilter struct {
	seen map[string]bool
}

func NewSeenFilter() *SeenFilter {
	makeSeen := make(map[string]bool)
	return &SeenFilter{
		seen: makeSeen,
	}
}

func (s *SeenFilter) Apply(items []fetcher.Item) []fetcher.Item {
	var fresh []fetcher.Item
	for _, item := range items {
		if s.seen[item.URL] {
			continue
		}
		s.seen[item.URL] = true
		fresh = append(fresh, item)
	}
	return fresh
}
