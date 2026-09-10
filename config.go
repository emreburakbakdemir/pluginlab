package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadConfig(path string) (*appConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil { return nil, fmt.Errorf("read config: %w", err) }

	var cfg appConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	return &cfg, nil
}
