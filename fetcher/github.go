package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GitHubFetcher struct {
	Token string
}

type githubNotification struct {
	Subject struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"subject"`
}

func (f *GitHubFetcher) Name() string {
	return "github"
}

func (f *GitHubFetcher) Fetch() ([]Item, error) {
	url := "https://api.github.com/notifications"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+f.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github: unexpected status %d", resp.StatusCode)
	}

	var notifications []githubNotification

	if err := json.NewDecoder(resp.Body).Decode(&notifications); err != nil {
		return nil, err
	}

	var items []Item

	for _, n := range notifications {
		items = append(items, Item{
			Title:  n.Subject.Title,
			URL:    n.Subject.URL,
			Source: f.Name(),
		})
	}

	return items, nil
}
