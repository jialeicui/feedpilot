package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the configuration for the server.
type Config struct {
	// Addr is the address to listen on. e.g. ":8080"
	Addr string `json:"addr" yaml:"addr"`

	// Debug is the flag to enable debug mode.
	Debug bool `json:"debug" yaml:"debug"`

	Platforms []*PlatformConfig `json:"platforms" yaml:"platforms"`
}

type PlatformConfig struct {
	// Name is the name of the platform.
	Name string `json:"name" yaml:"name"`

	Weibo *WeiboConfig `json:"weibo,omitempty" yaml:"weibo,omitempty"`
}

type WeiboConfig struct {
	RandomSleep bool `json:"random_sleep" yaml:"random_sleep"`
}

func Load(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := new(Config)
	if err := yaml.Unmarshal(content, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
