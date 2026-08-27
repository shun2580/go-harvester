package notifier

import (
	"go-harvester/fetcher"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotify(t *testing.T) {
	var gotBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		gotBody = string(body)

		if r.Method != "POST" {
			t.Errorf("Method = %q, want %q", r.Method, "POST")
		}

		if r.URL.Path != "/test" {
			t.Errorf("URL.Path = %q, want %q", r.URL.Path, "/test")
		}
	}))
	defer server.Close()

	n := &NtfyNotifier{URL: server.URL, Topic: "test"}

	item := fetcher.Item{Title: "Test01", URL: "https://example.com"}

	err := n.Notify(item)
	if err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	want := "Test01\nhttps://example.com"
	if gotBody != want {
		t.Errorf("body = %q, want %q", gotBody, want)
	}
}
