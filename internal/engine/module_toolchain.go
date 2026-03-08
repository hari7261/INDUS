package engine

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// toolchainModule provides commands for detecting and working with development toolchains
type toolchainModule struct {
	engine *Engine
}

func (m *toolchainModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "tools scan":
		return m.scan(inv)
	case "tools check":
		return m.check(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

// toolchainInfo represents information about a detected toolchain
type toolchainInfo struct {
	Name      string
	Command   string
	Installed bool
	Version   string
	Path      string
}

// scan detects all available development toolchains
func (m *toolchainModule) scan(inv Invocation) Response {
	toolchains := []toolchainInfo{
		// Programming languages
		{Name: "Python", Command: "python --version"},
		{Name: "Python3", Command: "python3 --version"},
		{Name: "Node.js", Command: "node --version"},
		{Name: "Go", Command: "go version"},
		{Name: "Java", Command: "java -version"},
		{Name: "Rust", Command: "rustc --version"},
		{Name: "GCC", Command: "gcc --version"},
		{Name: "G++", Command: "g++ --version"},
		{Name: "Clang", Command: "clang --version"},
		{Name: "Ruby", Command: "ruby --version"},
		{Name: "PHP", Command: "php --version"},
		{Name: ".NET", Command: "dotnet --version"},

		// Package managers
		{Name: "npm", Command: "npm --version"},
		{Name: "pip", Command: "pip --version"},
		{Name: "pip3", Command: "pip3 --version"},
		{Name: "yarn", Command: "yarn --version"},
		{Name: "pnpm", Command: "pnpm --version"},
		{Name: "Cargo", Command: "cargo --version"},
		{Name: "Maven", Command: "mvn --version"},
		{Name: "Gradle", Command: "gradle --version"},
		{Name: "Composer", Command: "composer --version"},
		{Name: "gem", Command: "gem --version"},

		// Version control
		{Name: "Git", Command: "git --version"},

		// Container & orchestration
		{Name: "Docker", Command: "docker --version"},
		{Name: "Kubernetes", Command: "kubectl version --client"},

		// Build tools
		{Name: "Make", Command: "make --version"},
		{Name: "CMake", Command: "cmake --version"},
	}

	// Detect toolchains concurrently
	var wg sync.WaitGroup
	results := make(chan toolchainInfo, len(toolchains))

	for _, tc := range toolchains {
		wg.Add(1)
		go func(tool toolchainInfo) {
			defer wg.Done()
			info := detectToolchain(tool.Name, tool.Command)
			results <- info
		}(tc)
	}

	wg.Wait()
	close(results)

	// Collect results
	var installed, missing []toolchainInfo
	for info := range results {
		if info.Installed {
			installed = append(installed, info)
		} else {
			missing = append(missing, info)
		}
	}

	// Format output
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "platform=%s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(buffer, "installed=%d\nmissing=%d\n\n", len(installed), len(missing))

	if len(installed) > 0 {
		fmt.Fprintln(buffer, "INSTALLED:")
		for _, info := range installed {
			if info.Version != "" {
				fmt.Fprintf(buffer, "  %-15s %s\n", info.Name, info.Version)
			} else {
				fmt.Fprintf(buffer, "  %-15s available\n", info.Name)
			}
		}
	}

	if len(missing) > 0 {
		fmt.Fprintln(buffer, "\nNOT FOUND:")
		for _, info := range missing {
			fmt.Fprintf(buffer, "  %s\n", info.Name)
		}
	}

	return Response{Output: strings.TrimSpace(buffer.String())}
}

// check verifies if a specific tool is available
func (m *toolchainModule) check(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind tools check <tool>")}
	}

	toolName := inv.Parsed.Positionals[0]

	// Try common version commands
	versionCommands := []string{
		fmt.Sprintf("%s --version", toolName),
		fmt.Sprintf("%s -version", toolName),
		fmt.Sprintf("%s version", toolName),
		fmt.Sprintf("%s -v", toolName),
	}

	var info toolchainInfo
	for _, cmd := range versionCommands {
		info = detectToolchain(toolName, cmd)
		if info.Installed {
			break
		}
	}

	if !info.Installed {
		return Response{Output: fmt.Sprintf("tool=%s\nstatus=not_found", toolName)}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "tool=%s\nstatus=installed\n", toolName)
	if info.Version != "" {
		fmt.Fprintf(buffer, "version=%s\n", info.Version)
	}
	if info.Path != "" {
		fmt.Fprintf(buffer, "path=%s\n", info.Path)
	}

	return Response{Output: strings.TrimSpace(buffer.String())}
}

// detectToolchain checks if a toolchain is installed and retrieves version info
func detectToolchain(name, versionCmd string) toolchainInfo {
	parts := strings.Fields(versionCmd)
	if len(parts) == 0 {
		return toolchainInfo{Name: name, Installed: false}
	}

	command := parts[0]
	args := parts[1:]

	// Find the executable path
	path, err := exec.LookPath(command)
	if err != nil {
		return toolchainInfo{Name: name, Command: command, Installed: false}
	}

	// Get version information
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	version := ""
	if err == nil && len(output) > 0 {
		// Extract first line as version info
		lines := strings.Split(string(output), "\n")
		if len(lines) > 0 {
			version = strings.TrimSpace(lines[0])
			// Limit version string length
			if len(version) > 60 {
				version = version[:60] + "..."
			}
		}
	}

	return toolchainInfo{
		Name:      name,
		Command:   command,
		Installed: true,
		Version:   version,
		Path:      path,
	}
}
