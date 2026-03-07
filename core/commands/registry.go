package commandcatalog

import (
	"embed"
	"encoding/json"
	"os"
)

//go:embed registry.json
var embeddedRegistry embed.FS

type Entry struct {
	Description string   `json:"description"`
	Module      string   `json:"module"`
	Category    string   `json:"category"`
	Usage       string   `json:"usage"`
	Example     string   `json:"example"`
	Since       string   `json:"since"`
	Aliases     []string `json:"aliases,omitempty"`
	CacheTTLMS  int      `json:"cache_ttl_ms,omitempty"`
}

type File struct {
	Version  string           `json:"version"`
	Commands map[string]Entry `json:"commands"`
}

func Load() (*File, error) {
	data, err := embeddedRegistry.ReadFile("registry.json")
	if err != nil {
		return nil, err
	}
	return decode(data)
}

func LoadFromPath(path string) (*File, error) {
	if path == "" {
		return Load()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Load()
	}
	return decode(data)
}

func decode(data []byte) (*File, error) {
	var file File
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	if file.Commands == nil {
		file.Commands = map[string]Entry{}
	}
	return &file, nil
}
