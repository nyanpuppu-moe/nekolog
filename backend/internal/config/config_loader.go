package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load() (*Config, error) {
	var cfg Config

	configPath := "configs/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist at: %s", configPath)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	if cfg.Server.SessionStore.PrivateKey == "" {
		return nil, fmt.Errorf("session store private key is empty")
	}

	return &cfg, nil
}
