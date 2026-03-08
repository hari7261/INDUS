package engine

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// executeSystemCommand executes a command through the system shell (cmd.exe on Windows, sh on Unix)
func (e *Engine) executeSystemCommand(ctx context.Context, session *Session, tokens []string) Response {
	if len(tokens) == 0 {
		return Response{Err: missingCommandError()}
	}

	start := time.Now()
	
	var cmd *exec.Cmd
	commandLine := strings.Join(tokens, " ")
	
	// Use the appropriate shell for the platform
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd.exe", "/C", commandLine)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", commandLine)
	}
	
	// Set the working directory to the session's current directory
	cmd.Dir = session.CWD()
	
	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Execute the command
	err := cmd.Run()
	
	duration := time.Since(start)
	
	// Prepare response
	output := stdout.String()
	errorOutput := stderr.String()
	
	if err != nil {
		// Command failed - return error with stderr
		if errorOutput != "" {
			return Response{
				Duration: duration,
				Err: &IndError{
					Code:       "IND_ERR_SYSTEM",
					Command:    commandLine,
					Message:    errorOutput,
					Suggestion: "verify the command syntax and try again",
				},
			}
		}
		return Response{
			Duration: duration,
			Err: &IndError{
				Code:       "IND_ERR_SYSTEM",
				Command:    commandLine,
				Message:    err.Error(),
				Suggestion: "verify the command exists and is available in your PATH",
			},
		}
	}
	
	// Success - return output
	// If both stdout and stderr have content, combine them
	finalOutput := output
	if errorOutput != "" {
		if finalOutput != "" {
			finalOutput += "\n" + errorOutput
		} else {
			finalOutput = errorOutput
		}
	}
	
	return Response{
		Output:   strings.TrimSpace(finalOutput),
		Duration: duration,
	}
}

// isLikelySystemCommand checks if a command is likely a system command
// This helps provide better error messages
func isLikelySystemCommand(cmd string) bool {
	// Common system commands and development tools
	commonCommands := map[string]bool{
		// Windows commands
		"dir": true, "cd": true, "copy": true, "move": true, "del": true, "mkdir": true,
		"rmdir": true, "type": true, "echo": true, "cls": true, "ipconfig": true,
		"ping": true, "netstat": true, "tasklist": true, "taskkill": true,
		
		// Unix commands
		"ls": true, "pwd": true, "cat": true, "grep": true, "find": true,
		"chmod": true, "chown": true, "ps": true, "kill": true, "wget": true, "curl": true,
		
		// Version control
		"git": true, "svn": true, "hg": true,
		
		// Programming languages & tools
		"python": true, "python3": true, "pip": true, "pip3": true,
		"node": true, "npm": true, "npx": true, "yarn": true, "pnpm": true,
		"go": true, "cargo": true, "rustc": true,
		"java": true, "javac": true, "mvn": true, "gradle": true,
		"gcc": true, "g++": true, "clang": true, "make": true, "cmake": true,
		"dotnet": true, "msbuild": true,
		"ruby": true, "gem": true,
		"php": true, "composer": true,
		
		// Container & orchestration
		"docker": true, "docker-compose": true, "podman": true,
		"kubectl": true, "helm": true, "minikube": true,
		
		// Editors & IDEs (command line)
		"code": true, "vim": true, "vi": true, "nano": true, "emacs": true,
		
		// Build tools
		"ant": true, "bazel": true, "buck": true,
		
		// Package managers
		"apt": true, "apt-get": true, "yum": true, "dnf": true, "brew": true, "choco": true,
		
		// Database clients
		"mysql": true, "psql": true, "mongo": true, "redis-cli": true, "sqlite3": true,
		
		// Cloud CLIs
		"aws": true, "az": true, "gcloud": true, "heroku": true,
		
		// Other common tools
		"ssh": true, "scp": true, "ftp": true, "telnet": true,
		"tar": true, "gzip": true, "zip": true, "unzip": true,
		"systemctl": true, "service": true,
	}
	
	return commonCommands[strings.ToLower(cmd)]
}
