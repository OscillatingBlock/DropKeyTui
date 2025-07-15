package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	UserID     string `json:"user_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Could not get user config directory")
		return "", fmt.Errorf("Could not get user config directory: %w", err)
	}

	appConfigDir := filepath.Join(configDir, "Drop-Key-TUI")
	if err := os.MkdirAll(appConfigDir, 0o755); err != nil {
		slog.Error("could not create app config directory")
		return "", fmt.Errorf("could not create app config directory: %w", err)
	}

	return filepath.Join(appConfigDir, "config.json"), nil
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	rawConfigData, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("config file not found, please register first")
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(rawConfigData) == 0 {
		slog.Error("config file is empty")
		return nil, fmt.Errorf("config file is empty")
	}

	var config Config
	if err := json.Unmarshal(rawConfigData, &config); err != nil {
		slog.Error("error while decoding config JSONL %w", "error", err)
		return nil, fmt.Errorf("failed to decode config JSON: %w", err)
	}

	return &config, nil
}

func Save(config *Config) error {
	if config.UserID == "" || config.PublicKey == "" || config.PrivateKey == "" {
		return errors.New("cannot save incomplete config")
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	jsonBody, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	if err := os.WriteFile(configPath, jsonBody, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
