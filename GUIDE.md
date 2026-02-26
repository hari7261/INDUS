# INDUS CLI - User Guide

## Table of Contents
1. [Getting Started](#getting-started)
2. [Interactive Mode](#interactive-mode)
3. [Command Reference](#command-reference)
4. [Examples](#examples)
5. [Configuration](#configuration)

---

## Getting Started

### Installation

Build from source:
```bash
go build -o indus.exe ./cmd/indus
```

### Running INDUS

**Interactive Mode** (REPL):
```bash
indus
```

**Direct Command Mode**:
```bash
indus <command> [flags]
```

---

## Interactive Mode

When you run `indus` without arguments, you enter an interactive REPL (Read-Eval-Print Loop).

### Features
- Beautiful ASCII art banner with Indian flag colors
- Command history
- Tab completion hints
- Colorized output
- Built-in help system

### REPL Commands
- `help` or `?` - Show help
- `clear` or `cls` - Clear screen
- `exit`, `quit`, or `q` - Exit REPL
- Press `Ctrl+C` - Cancel current operation
- Press `Ctrl+D` - Exit REPL

---

## Command Reference

### 1. init - Initialize Project

Create a new project structure with standard directories.

**Syntax:**
```bash
init --name <project-name> [--dir <directory>]
```

**Flags:**
- `--name` (required) - Project name
- `--dir` (optional) - Target directory (default: current directory)

**Example:**
```bash
> init --name myapp --dir ~/projects
```

**Output:**
```
project_dir=~/projects/myapp
project_name=myapp
```

**Created Structure:**
```
myapp/
├── cmd/
├── internal/
├── pkg/
├── config/
└── README.md
```

---

### 2. run - Execute Workload

Simulate a workload with bounded concurrency using worker pools.

**Syntax:**
```bash
run [--workers <n>] [--tasks <n>]
```

**Flags:**
- `--workers` (optional) - Number of concurrent workers (default: 4)
- `--tasks` (optional) - Total tasks to process (default: 20)

**Example:**
```bash
> run --workers 8 --tasks 50
```

**Output:**
```
Starting run with 8 workers processing 50 tasks...
Progress: 5/50 tasks completed
Progress: 10/50 tasks completed
...
completed=50
failed=0
total=50
```

**Use Cases:**
- Testing concurrent processing
- Benchmarking worker pools
- Demonstrating graceful cancellation (Ctrl+C)

---

### 3. version - Show Version

Display version information.

**Syntax:**
```bash
version
```

**Example:**
```bash
> version
```

**Output:**
```
version=1.0.0
commit=abc123def
build_time=2026-02-26T10:00:00Z
```

---

### 4. http - Make HTTP Requests

Perform HTTP requests with automatic retry and timeout handling.

#### GET Request

**Syntax:**
```bash
http get <url> [--headers 'Key:Value,Key2:Value2']
```

**Examples:**
```bash
> http get https://api.github.com

> http get https://api.example.com/users --headers 'Authorization:Bearer token123'
```

#### POST Request

**Syntax:**
```bash
http post <url> <data> [--headers 'Key:Value']
```

**Examples:**
```bash
> http post https://api.example.com/users '{"name":"John","email":"john@example.com"}'

> http post https://httpbin.org/post '{"test":"data"}' --headers 'Content-Type:application/json'
```

#### PUT Request

**Syntax:**
```bash
http put <url> <data>
```

**Example:**
```bash
> http put https://api.example.com/users/123 '{"name":"Jane"}'
```

#### DELETE Request

**Syntax:**
```bash
http delete <url>
```

**Example:**
```bash
> http delete https://api.example.com/users/123
```

---

## Examples

### Example 1: Quick API Test
```bash
# Start REPL
indus

# Test GitHub API
> http get https://api.github.com/users/octocat

# Check rate limit
> http get https://api.github.com/rate_limit
```

### Example 2: Create and Test Project
```bash
# Initialize project
> init --name my-api-client --dir ~/dev

# Test concurrent processing
> run --workers 10 --tasks 100

# Check version
> version
```

### Example 3: API Integration Testing
```bash
# POST data
> http post https://jsonplaceholder.typicode.com/posts '{"title":"Test","body":"Content","userId":1}'

# GET the created resource
> http get https://jsonplaceholder.typicode.com/posts/1

# DELETE resource
> http delete https://jsonplaceholder.typicode.com/posts/1
```

### Example 4: Direct Command Mode
```bash
# Run commands directly from shell
indus version
indus init --name testapp
indus http get https://api.github.com
indus run --workers 4 --tasks 20
```

---

## Configuration

### Config File Location

Default: `~/.config/indus/config.yaml`

Override with environment variable:
```bash
set INDUS_CONFIG=C:\custom\path\config.yaml
indus
```

### Config File Format

```yaml
api_timeout: 30
max_retries: 3
```

### Environment Variables

- `INDUS_CONFIG` - Custom config file path

---

## Tips & Tricks

### 1. Canceling Operations
Press `Ctrl+C` to gracefully cancel any running operation. The CLI will clean up and return to the prompt.

### 2. Command History
Use arrow keys (↑/↓) to navigate through command history in the REPL.

### 3. Clear Screen
Type `clear` or `cls` to clear the screen and redisplay the banner.

### 4. Quick Exit
Type `q` for quick exit instead of typing `exit`.

### 5. Piping Output
In direct command mode, pipe output to files:
```bash
indus http get https://api.github.com > response.json
indus version > version.txt
```

### 6. Chaining Commands
In shell (not REPL):
```bash
indus init --name myapp && cd myapp && indus version
```

---

## Troubleshooting

### Issue: Command not found
**Solution:** Make sure `indus.exe` is in your PATH or use the full path:
```bash
C:\path\to\indus.exe
```

### Issue: HTTP request timeout
**Solution:** Increase timeout in config file or check network connection.

### Issue: Permission denied when creating project
**Solution:** Run with appropriate permissions or choose a different directory.

---

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+C` | Cancel current operation |
| `Ctrl+D` | Exit REPL |
| `↑` / `↓` | Navigate command history |
| `Tab` | Auto-complete (future) |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Internal error |
| 2 | User error (bad flags, invalid input) |
| 130 | Canceled (Ctrl+C) |

---

## Support

For issues or feature requests, refer to the project repository.

**Happy coding with INDUS! 🇮🇳**
