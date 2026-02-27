package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

	// Parse simple "key = value" format (lines starting with # are comments).
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "api_timeout":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.APITimeout = n
			}
		case "max_retries":
			if n, err := strconv.Atoi(val); err == nil && n >= 0 {
				cfg.MaxRetries = n
			}
		}
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
