package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type HTTP struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	Seed    int64    `yaml:"seed"`
	Symbols []string `yaml:"symbols"`
	Hz      float64  `yaml:"hz"` // ticks per second per symbol
	HTTP    HTTP     `yaml:"http"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig(), nil
	} // fallback
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func defaultConfig() *Config {
	return &Config{
		Seed:    42,
		Symbols: []string{"AAPL", "MSFT", "TSLA"},
		Hz:      1.0,
		HTTP:    HTTP{Addr: ":8080"},
	}
}

// Utility
func NowMillis() int64 { return time.Now().UnixNano() / int64(time.Millisecond) }
