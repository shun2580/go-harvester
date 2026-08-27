package notifier

import (
	"go-harvester/fetcher"
	"net/http"
	"strings"
)

type Notifier interface {
	Notify(item fetcher.Item) error
}

type NtfyNotifier struct {
	URL   string
	Topic string
}

func (n *NtfyNotifier) Notify(item fetcher.Item) error {
	url := n.URL + "/" + n.Topic
	contentType := "text/plain"
	body := strings.NewReader(item.Title + "\n" + item.URL)

	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
