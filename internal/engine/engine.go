package engine

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	commandcatalog "indus/core/commands"
	"indus/internal/config"
)

type moduleFactory func(*Engine) Module

type Engine struct {
	version         string
	commit          string
	buildTime       string
	registry        map[string]CommandMeta
	registryVer     string
	moduleFactories map[string]moduleFactory
	modules         map[string]Module
	moduleMu        sync.Mutex
	cache           *responseCache
	state           *StateStore
	paths           Paths
	cfg             *config.Config
	bufferPool      sync.Pool
	metricsMu       sync.RWMutex
	metrics         []Metric
}

func New(opts Options) (*Engine, error) {
	paths, err := discoverPaths()
	if err != nil {
		return nil, err
	}

	registryFile, err := commandcatalog.LoadFromPath(paths.RegistryPath)
	if err != nil {
		return nil, err
	}

	state, err := NewStateStore(paths.StateFile)
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}

	engine := &Engine{
		version:         opts.Version,
		commit:          opts.Commit,
		buildTime:       opts.BuildTime,
		registry:        map[string]CommandMeta{},
		registryVer:     registryFile.Version,
		moduleFactories: defaultModuleFactories(),
		modules:         map[string]Module{},
		cache:           newResponseCache(),
		state:           state,
		paths:           paths,
		cfg:             cfg,
		bufferPool: sync.Pool{
			New: func() any {
				return &bytes.Buffer{}
			},
		},
	}

	for path, entry := range registryFile.Commands {
		engine.registry[path] = CommandMeta{Path: path, Entry: entry}
	}

	return engine, nil
}

func (e *Engine) NewSession(cwd string) *Session {
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	state := e.state.Snapshot()
	theme, ok := themeByName(state.Theme)
	if !ok {
		theme = defaultTheme()
	}
	return &Session{
		cwd:   cwd,
		theme: theme,
	}
}

func (e *Engine) ExecuteLine(ctx context.Context, session *Session, line string) Response {
	tokens := ParseCommandLine(strings.TrimSpace(line))
	return e.ExecuteTokens(ctx, session, tokens, ModeInteractive)
}

func (e *Engine) ExecuteTokens(ctx context.Context, session *Session, tokens []string, mode Mode) Response {
	if session == nil {
		session = e.NewSession("")
	}

	if len(tokens) == 0 {
		return Response{Err: missingCommandError()}
	}

	if len(tokens) == 1 && (tokens[0] == "help" || tokens[0] == "?") {
		return Response{Output: e.HelpText()}
	}

	normalized, warning := e.normalizeTokens(tokens, mode)
	if len(normalized) == 0 {
		return Response{Err: missingCommandError()}
	}

	if len(normalized) == 1 && normalized[0] == "help" {
		return Response{Warning: warning, Output: e.HelpText()}
	}

	commandString := "ind " + strings.Join(normalized, " ")
	meta, args, ok := e.resolve(normalized)
	if !ok {
		return Response{Warning: warning, Err: unknownCommandError(commandString)}
	}

	invocation := Invocation{
		Path:      meta.Path,
		Args:      args,
		Parsed:    ParseArgs(args),
		Meta:      meta.Entry,
		Session:   session,
		Command:   commandString,
		RootToken: "ind",
		Mode:      mode,
	}

	if cached, ok := e.cachedResponse(invocation); ok {
		cached.Warning = warning
		e.recordMetric(invocation.Command, cached.Duration, true)
		return cached
	}

	module, err := e.module(meta.Module)
	if err != nil {
		return Response{Warning: warning, Err: registryError(err)}
	}

	start := time.Now()
	resultCh := make(chan Response, 1)

	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				resultCh <- Response{Err: panicError(invocation.Command, recovered)}
			}
		}()
		resultCh <- module.Execute(ctx, invocation)
	}()

	response := <-resultCh
	response.Warning = mergeWarnings(warning, response.Warning)
	response.Duration = time.Since(start)

	if response.Err == nil && invocation.Meta.CacheTTLMS > 0 && !hasEffects(response) {
		e.cache.Set(e.cacheKey(invocation), response, time.Duration(invocation.Meta.CacheTTLMS)*time.Millisecond)
	}

	e.recordMetric(invocation.Command, response.Duration, response.Cached)
	return response
}

func (e *Engine) HelpText() string {
	grouped := map[string][]CommandMeta{}
	for _, meta := range e.registry {
		grouped[meta.Category] = append(grouped[meta.Category], meta)
	}

	categories := make([]string, 0, len(grouped))
	for category := range grouped {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	buffer := e.getBuffer()
	defer e.putBuffer(buffer)

	fmt.Fprintln(buffer, "INDUS Terminal v"+e.version)
	fmt.Fprintln(buffer, "")
	fmt.Fprintln(buffer, "Usage:")
	fmt.Fprintln(buffer, "  ind <command> [options]")
	fmt.Fprintln(buffer, "")

	for _, category := range categories {
		entries := grouped[category]
		sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
		fmt.Fprintf(buffer, "%s:\n", strings.ToUpper(category))
		for _, entry := range entries {
			fmt.Fprintf(buffer, "  %-18s %s\n", entry.Path, entry.Description)
		}
		fmt.Fprintln(buffer)
	}

	fmt.Fprintln(buffer, "Versioned documentation:")
	fmt.Fprintf(buffer, "  %s\n", filepath.Join(e.paths.DocsDir, "index.html"))
	fmt.Fprintf(buffer, "  %s\n", filepath.Join(e.paths.DocsDir, "commands.html"))
	fmt.Fprintf(buffer, "  %s\n", filepath.Join(e.paths.DocsDir, "versions.html"))
	return strings.TrimSpace(buffer.String())
}

func (e *Engine) normalizeTokens(tokens []string, mode Mode) ([]string, string) {
	if len(tokens) == 0 {
		return nil, ""
	}

	first := strings.ToLower(tokens[0])
	switch first {
	case "ind":
		return append([]string(nil), tokens[1:]...), ""
	case "indus":
		return append([]string(nil), tokens[1:]...), "Deprecated command detected.\nUse:\nind " + strings.Join(tokens[1:], " ")
	}

	if alias, ok := legacyAliases()[first]; ok {
		normalized := append([]string{}, alias...)
		normalized = append(normalized, tokens[1:]...)
		return normalized, "Deprecated command detected.\nUse:\nind " + strings.Join(normalized, " ")
	}

	if mode == ModeInteractive {
		if _, _, ok := e.resolve(tokens); ok {
			return append([]string(nil), tokens...), "Deprecated command detected.\nUse:\nind " + strings.Join(tokens, " ")
		}
	}

	return append([]string(nil), tokens...), ""
}

func (e *Engine) resolve(tokens []string) (CommandMeta, []string, bool) {
	if len(tokens) == 0 {
		return CommandMeta{}, nil, false
	}

	for size := min(2, len(tokens)); size >= 1; size-- {
		path := strings.ToLower(strings.Join(tokens[:size], " "))
		meta, ok := e.registry[path]
		if ok {
			return meta, tokens[size:], true
		}
	}

	return CommandMeta{}, nil, false
}

func (e *Engine) module(name string) (Module, error) {
	e.moduleMu.Lock()
	defer e.moduleMu.Unlock()

	if module, ok := e.modules[name]; ok {
		return module, nil
	}

	factory, ok := e.moduleFactories[name]
	if !ok {
		return nil, fmt.Errorf("module %q is not registered", name)
	}

	module := factory(e)
	e.modules[name] = module
	return module, nil
}

func (e *Engine) cachedResponse(inv Invocation) (Response, bool) {
	if inv.Meta.CacheTTLMS <= 0 {
		return Response{}, false
	}
	return e.cache.Get(e.cacheKey(inv))
}

func (e *Engine) cacheKey(inv Invocation) string {
	return inv.Path + "|" + inv.Session.cwd + "|" + strings.Join(inv.Args, "\x00")
}

func (e *Engine) recordMetric(command string, duration time.Duration, cached bool) {
	e.metricsMu.Lock()
	defer e.metricsMu.Unlock()

	e.metrics = append(e.metrics, Metric{
		Command:  command,
		Duration: duration,
		Cached:   cached,
		At:       time.Now(),
	})
	if len(e.metrics) > 40 {
		e.metrics = append([]Metric(nil), e.metrics[len(e.metrics)-40:]...)
	}
}

func (e *Engine) Metrics() []Metric {
	e.metricsMu.RLock()
	defer e.metricsMu.RUnlock()
	return append([]Metric(nil), e.metrics...)
}

func (e *Engine) Reload() error {
	registryFile, err := commandcatalog.LoadFromPath(e.paths.RegistryPath)
	if err != nil {
		return err
	}

	updated := map[string]CommandMeta{}
	for path, entry := range registryFile.Commands {
		updated[path] = CommandMeta{Path: path, Entry: entry}
	}
	e.registry = updated
	e.registryVer = registryFile.Version
	e.cache.Clear()
	return e.state.Load()
}

func (e *Engine) getBuffer() *bytes.Buffer {
	buffer := e.bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	return buffer
}

func (e *Engine) putBuffer(buffer *bytes.Buffer) {
	buffer.Reset()
	e.bufferPool.Put(buffer)
}

func (e *Engine) StateSnapshot() PersistentState {
	return e.state.Snapshot()
}

func (e *Engine) Paths() Paths {
	return e.paths
}

func (e *Engine) Version() string {
	return e.version
}

func (e *Engine) Commit() string {
	return e.commit
}

func (e *Engine) BuildTime() string {
	return e.buildTime
}

func (e *Engine) RegistryVersion() string {
	return e.registryVer
}

func (e *Engine) Config() *config.Config {
	return e.cfg
}

func hasEffects(response Response) bool {
	return response.Effects != (Effects{})
}

func mergeWarnings(first, second string) string {
	switch {
	case first == "":
		return second
	case second == "":
		return first
	default:
		return first + "\n" + second
	}
}

func lowerSlice(values []string) []string {
	normalized := make([]string, len(values))
	for i, value := range values {
		normalized[i] = strings.ToLower(value)
	}
	return normalized
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func legacyAliases() map[string][]string {
	return map[string][]string{
		"cd":      {"work", "switch"},
		"clear":   {"term", "clearx"},
		"cls":     {"term", "clearx"},
		"color":   {"term", "theme"},
		"envlist": {"env", "list"},
		"http":    {"net", "fetch"},
		"init":    {"proj", "init"},
		"pwd":     {"status"},
		"run":     {"dev", "bench"},
	}
}
