package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Load は、path の JSON ファイルを読んで Config を作る。
// 秘匿値（GitHub トークン・ntfy の topic）はファイルに書かず、呼んだ側で環境変数から入れる。
func Load(path string) (Config, error) {
	var cfg Config

	if file, err := os.Open(path); err != nil {
		return cfg, err
	} else {
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&cfg); err != nil {
			return cfg, err
		}
	}

	// 間隔が0以下だと、time.NewTicker が落ちる。読み込みの段で弾いておく
	if cfg.Interval <= 0 {
		return cfg, fmt.Errorf("config: interval must be positive, got %d", cfg.Interval)
	}

	return cfg, nil
}
