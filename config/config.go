package config

type Config struct {
	Interval int
	GitHub   GitHubConfig
	RSS      []RSSConfig
	Ntfy     NtfyConfig
}

type GitHubConfig struct {
	Token string
}

type RSSConfig struct {
	Name string
	URL  string
}

type NtfyConfig struct {
	URL   string
	Topic string
}
