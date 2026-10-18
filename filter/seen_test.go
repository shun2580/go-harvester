package filter

import (
	"testing"

	"go-harvester/fetcher"
)

func TestSeenFilterSkipsSeenItems(t *testing.T) {
	f := NewSeenFilter()
	a := fetcher.Item{Title: "A", URL: "https://example.com/a"}
	b := fetcher.Item{Title: "B", URL: "https://example.com/b"}
	c := fetcher.Item{Title: "C", URL: "https://example.com/c"}

	got := f.Apply([]fetcher.Item{a, b})

	if len(got) != 2 {
		t.Fatalf("1回目: len = %d, want 2", len(got))
	}
	if got[0].URL != a.URL {
		t.Errorf("1回目: URL = %q, want %q", got[0].URL, a.URL)
	}
	if got[1].URL != b.URL {
		t.Errorf("1回目: URL = %q, want %q", got[1].URL, b.URL)
	}

	got = f.Apply([]fetcher.Item{a, b, c})

	if len(got) != 1 {
		t.Fatalf("2回目: len = %d, want 1", len(got))
	}
	if got[0].Title != c.Title {
		t.Errorf("2回目: Title = %q, want %q", got[0].Title, c.Title)
	}

	if got[0].URL != c.URL {
		t.Errorf("2回目: URL = %q, want %q", got[0].URL, c.URL)
	}
}
