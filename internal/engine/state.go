package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type StateStore struct {
	mu    sync.RWMutex
	path  string
	state PersistentState
}

func NewStateStore(path string) (*StateStore, error) {
	store := &StateStore{path: path}
	if err := store.Load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *StateStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = defaultState()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return s.saveLocked()
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &s.state); err != nil {
		return err
	}

	if s.state.ManagedEnv == nil {
		s.state.ManagedEnv = map[string]string{}
	}
	if s.state.Packages == nil {
		s.state.Packages = map[string]PackageRecord{}
	}
	if s.state.Tasks == nil {
		s.state.Tasks = map[string]TaskRecord{}
	}
	s.state.Profile = normalizeTerminalProfile(s.state.Profile)
	return nil
}

func (s *StateStore) Snapshot() PersistentState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clone := s.state
	clone.ManagedEnv = cloneStringMap(s.state.ManagedEnv)
	clone.Packages = clonePackageMap(s.state.Packages)
	clone.Tasks = cloneTaskMap(s.state.Tasks)
	clone.Workspaces = append([]WorkspaceRecord(nil), s.state.Workspaces...)
	return clone
}

func (s *StateStore) Update(fn func(*PersistentState)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fn(&s.state)
	if s.state.ManagedEnv == nil {
		s.state.ManagedEnv = map[string]string{}
	}
	if s.state.Packages == nil {
		s.state.Packages = map[string]PackageRecord{}
	}
	if s.state.Tasks == nil {
		s.state.Tasks = map[string]TaskRecord{}
	}
	s.state.Profile = normalizeTerminalProfile(s.state.Profile)

	return s.saveLocked()
}

func (s *StateStore) saveLocked() error {
	data, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func defaultState() PersistentState {
	return PersistentState{
		Theme:      "saffron",
		ManagedEnv: map[string]string{},
		Packages:   map[string]PackageRecord{},
		Tasks:      map[string]TaskRecord{},
		Profile:    defaultTerminalProfile(),
	}
}

func cloneStringMap(source map[string]string) map[string]string {
	result := map[string]string{}
	for key, value := range source {
		result[key] = value
	}
	return result
}

func clonePackageMap(source map[string]PackageRecord) map[string]PackageRecord {
	result := map[string]PackageRecord{}
	for key, value := range source {
		result[key] = value
	}
	return result
}

func cloneTaskMap(source map[string]TaskRecord) map[string]TaskRecord {
	result := map[string]TaskRecord{}
	for key, value := range source {
		value.Commands = append([]string(nil), value.Commands...)
		result[key] = value
	}
	return result
}

func defaultTerminalProfile() TerminalProfile {
	return TerminalProfile{
		ShowBanner:       true,
		BannerAnimation:  "mascot-wave",
		BannerDurationMS: 5000,
		CompactMode:      false,
		PromptLabel:      "INDUS",
	}
}

func normalizeTerminalProfile(profile TerminalProfile) TerminalProfile {
	defaults := defaultTerminalProfile()
	isZero := profile == (TerminalProfile{})

	if profile.BannerAnimation == "" {
		profile.BannerAnimation = defaults.BannerAnimation
	}
	if profile.BannerDurationMS <= 0 {
		profile.BannerDurationMS = defaults.BannerDurationMS
	}
	if profile.PromptLabel == "" {
		profile.PromptLabel = defaults.PromptLabel
	}
	if isZero {
		profile.ShowBanner = defaults.ShowBanner
	}
	return profile
}
