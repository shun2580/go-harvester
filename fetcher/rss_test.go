package fetcher

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testRSS = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>test feed</title>
    <item>
      <title>first</title>
      <link>https://example.com/1</link>
    </item>
    <item>
      <title>second</title>
      <link>https://example.com/2</link>
    </item>
  </channel>
</rss>`

func TestRSSFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, testRSS)
	}))
	defer server.Close()

	f := &RSSFetcher{FeedName: "test", FeedURL: server.URL}
	items, err := f.Fetch()
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len = %d, want 2", len(items))
	}
	if items[0].Title != "first" || items[0].URL != "https://example.com/1" {
		t.Errorf("items[0] = %+v", items[0])
	}
	if items[0].Source != "test" {
		t.Errorf("items[0].Source = %q, want %q", items[0].Source, "test")
	}

}

func TestRSSFetchBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, testRSS)
	}))
	defer server.Close()

	f := &RSSFetcher{FeedName: "test", FeedURL: server.URL}
	_, err := f.Fetch()
	if err == nil {
		t.Errorf("Fetch did not return error on bad status")
	}

}
