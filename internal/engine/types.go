package engine

import (
	"context"
	"time"

	commandcatalog "indus/core/commands"
)

type Mode int

const (
	ModeExecutable Mode = iota
	ModeInteractive
)

type Theme struct {
	Name   string
	Prompt string
}

type Session struct {
	cwd   string
	theme Theme
}

type Effects struct {
	NextDir     string
	ClearScreen bool
	Exit        bool
	Theme       Theme
}

type Response struct {
	Output   string
	Warning  string
	Err      *IndError
	Cached   bool
	Duration time.Duration
	Effects  Effects
}

type Invocation struct {
	Path      string
	Args      []string
	Parsed    ParsedArgs
	Meta      commandcatalog.Entry
	Session   *Session
	Command   string
	RootToken string
	Mode      Mode
}

type Module interface {
	Execute(context.Context, Invocation) Response
}

type CommandMeta struct {
	Path string
	commandcatalog.Entry
}

type WorkspaceRecord struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Pinned   bool   `json:"pinned"`
	LastUsed string `json:"last_used,omitempty"`
}

type PackageRecord struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	InstalledAt string `json:"installed_at"`
}

type TaskRecord struct {
	Name      string   `json:"name"`
	Commands  []string `json:"commands"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	LastRunAt string   `json:"last_run_at,omitempty"`
}

type TerminalProfile struct {
	ShowBanner       bool   `json:"show_banner"`
	BannerAnimation  string `json:"banner_animation"`
	BannerDurationMS int    `json:"banner_duration_ms"`
	CompactMode      bool   `json:"compact_mode"`
	PromptLabel      string `json:"prompt_label"`
}

type PersistentState struct {
	Theme           string                   `json:"theme"`
	ManagedEnv      map[string]string        `json:"managed_env"`
	Workspaces      []WorkspaceRecord        `json:"workspaces"`
	ActiveWorkspace string                   `json:"active_workspace"`
	Packages        map[string]PackageRecord `json:"packages"`
	Tasks           map[string]TaskRecord    `json:"tasks"`
	Profile         TerminalProfile          `json:"profile"`
}

type Paths struct {
	RootDir      string
	DocsDir      string
	RegistryPath string
	StateDir     string
	StateFile    string
	CacheDir     string
	ReportsDir   string
	UpdatesDir   string
}

type Metric struct {
	Command  string
	Duration time.Duration
	Cached   bool
	At       time.Time
}

type Options struct {
	Version   string
	Commit    string
	BuildTime string
}

func (s *Session) CWD() string {
	return s.cwd
}

func (s *Session) SetCWD(path string) {
	s.cwd = path
}

func (s *Session) Theme() Theme {
	return s.theme
}

func (s *Session) SetTheme(theme Theme) {
	s.theme = theme
}
