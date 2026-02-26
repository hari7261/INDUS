package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	ConfigPath string
	APITimeout int
	MaxRetries int
}

func Load() (*Config, error) {
	path := configPath()
	
	cfg := Default()
	cfg.ConfigPath = path
	
	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}
	
	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	
	// Simple parsing - in production this would parse YAML
	// For now, just validate it exists and is readable
	if len(data) > 0 {
		// Config loaded successfully
	}
	
	return cfg, nil
}

func Default() *Config {
	return &Config{
		APITimeout: 30,
		MaxRetries: 3,
	}
}

func configPath() string {
	// Check environment variable first
	if path := os.Getenv("INDUS_CONFIG"); path != "" {
		return path
	}
	
	// Use default path
	home := homeDir()
	return filepath.Join(home, ".config", "indus", "config.yaml")
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		if home == "" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
		return home
	}
	
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	
	return "."
}
