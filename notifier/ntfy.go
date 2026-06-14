package notifier

import "go-harvester/fetcher"

type Notifier interface {
	Notify(item fetcher.Item) error
}

type NtfyNotifier struct {
	URL   string
	Topic string
}

func (n *NtfyNotifier) Notify(item fetcher.Item) error {
	// TODO: ntfy に HTTP POST する
	return nil
}
