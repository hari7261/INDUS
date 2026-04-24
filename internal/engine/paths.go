package engine

import (
	"os"
	"path/filepath"
)

func discoverPaths() (Paths, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Paths{}, err
	}

	root := os.Getenv("INDUS_ROOT_DIR")
	if root == "" {
		root = findRoot(wd)
	}
	if !hasRegistry(root) {
		if exe, exeErr := os.Executable(); exeErr == nil {
			exeRoot := findRoot(filepath.Dir(exe))
			if hasRegistry(exeRoot) {
				root = exeRoot
			} else {
				root = filepath.Dir(exe)
			}
		}
	}

	stateDir := os.Getenv("INDUS_STATE_DIR")
	cacheDir := os.Getenv("INDUS_CACHE_DIR")

	if stateDir == "" {
		configRoot, err := os.UserConfigDir()
		if err != nil {
			configRoot = wd
		}
		stateDir = filepath.Join(configRoot, "indus")
	}
	if cacheDir == "" {
		cacheRoot, err := os.UserCacheDir()
		if err != nil {
			cacheRoot = wd
		}
		cacheDir = filepath.Join(cacheRoot, "indus")
	}
	reportsDir := filepath.Join(stateDir, "reports")
	updatesDir := filepath.Join(stateDir, "updates")

	paths := Paths{
		RootDir:      root,
		StateDir:     stateDir,
		StateFile:    filepath.Join(stateDir, "state.json"),
		CacheDir:     cacheDir,
		ReportsDir:   reportsDir,
		UpdatesDir:   updatesDir,
		RegistryPath: filepath.Join(root, "core", "commands", "registry.json"),
		DocsDir:      filepath.Join(root, "docs"),
	}

	if err := os.MkdirAll(paths.StateDir, 0o755); err != nil {
		return Paths{}, err
	}
	if err := os.MkdirAll(paths.CacheDir, 0o755); err != nil {
		return Paths{}, err
	}
	if err := os.MkdirAll(paths.ReportsDir, 0o755); err != nil {
		return Paths{}, err
	}
	if err := os.MkdirAll(paths.UpdatesDir, 0o755); err != nil {
		return Paths{}, err
	}

	return paths, nil
}

func findRoot(start string) string {
	current := start
	for {
		if hasRegistry(current) {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return start
		}
		current = parent
	}
}

func hasRegistry(root string) bool {
	_, err := os.Stat(filepath.Join(root, "core", "commands", "registry.json"))
	return err == nil
}
