# INDUS Terminal 🇮🇳

<div align="center">

![INDUS Terminal](images/image1.png)

**Production-grade interactive terminal for API orchestration, agent runners, and developer tooling.**

Built with zero external dependencies using only the Go standard library.

[![Version](https://img.shields.io/badge/version-1.0.0-blue)](https://github.com/hari7261/indus)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/platform-Windows-0078D6?logo=windows)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/license-Production%20Use-green)](LICENSE)

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Documentation](#-documentation) • [Contributing](#-contributing)

</div>

---

## 🎯 What is INDUS?

INDUS is a modern, feature-rich terminal emulator designed for developers and DevOps engineers. It combines the power of traditional terminals (PowerShell, Git Bash) with built-in developer tools for API testing, project scaffolding, and concurrent workload simulation.

### Why INDUS?

- 🚀 **All-in-One**: Terminal + HTTP Client + Project Generator + Workload Simulator
- 🎨 **Beautiful UI**: Indian flag-themed banner with ANSI colors
- ⚡ **Zero Dependencies**: Pure Go standard library
- 💻 **Full Shell**: Works like PowerShell/Git Bash with `cd`, `pwd`, and all system commands
- 🌐 **Built-in HTTP Client**: Test APIs without curl or Postman
- 🔧 **Developer-Friendly**: Project initialization, concurrent testing, and more

---

## ✨ Features

### 🖥️ Interactive Terminal

![INDUS Terminal Interface](images/image.png)

- **Beautiful ASCII Banner** with Indian flag colors (Saffron, White, Green)
- **Smart Prompt** showing current directory: `INDUS ~/Documents >`
- **Full Shell Capabilities** like PowerShell/Git Bash
- **ANSI Color Support** for better readability
- **Command History** (coming soon)

### 🛠️ Built-in Commands

#### Project Initialization
```bash
init --name myproject --dir ~/projects
```
Creates a complete project structure with standard directories.

#### HTTP Client
```bash
http get https://api.github.com
http post https://api.example.com/users '{"name":"John"}'
http put https://api.example.com/users/123 '{"name":"Jane"}'
http delete https://api.example.com/users/123
```
Make HTTP requests with automatic retry and timeout handling.

#### Concurrent Workload Simulation
```bash
run --workers 8 --tasks 50
```
Test concurrent processing with worker pools and bounded concurrency.

#### Version Information
```bash
version
```
Display version, commit hash, and build time.

### 💻 System Integration

Run **ANY** Windows command:
```bash
ipconfig              # Network configuration
ping google.com       # Test connectivity
git status            # Git commands
docker ps             # Docker commands
npm install           # Node.js commands
python script.py      # Python scripts
cd Documents          # Change directory
pwd                   # Print working directory
```

### 🎨 Shell Commands

- `cd <path>` - Change directory
- `pwd` - Print working directory
- `clear` / `cls` - Clear screen
- `exit` / `quit` - Exit terminal
- `help` - Show all commands

---

## 🚀 Installation

### Quick Install (Recommended)

1. **Download** the latest release from [Releases](https://github.com/hari7261/indus/releases)
2. **Extract** all files to a folder
3. **Right-click** `install.bat` and select **"Run as administrator"**
4. Follow the on-screen instructions

The installer will:
- ✅ Install INDUS to `%LOCALAPPDATA%\INDUS`
- ✅ Add INDUS to your PATH
- ✅ Create Desktop shortcut
- ✅ Create Start Menu entry

### Portable Mode

Just download `indus.exe` and double-click to run. No installation needed!

### Build from Source

```bash
# Prerequisites
go install github.com/akavel/rsrc@latest

# Clone repository
git clone https://github.com/hari7261/indus.git
cd indus

# Build
build.bat

# Or use Make
make build
```

See [README-INSTALL.md](README-INSTALL.md) for detailed installation instructions.

---

## 💡 Usage

### Launch INDUS Terminal

After installation, launch INDUS in multiple ways:

1. **Desktop Shortcut**: Double-click "INDUS Terminal" on your desktop
2. **Start Menu**: Search for "INDUS" and click
3. **Command Line**: Open any terminal and type `indus`
4. **Direct**: Navigate to installation folder and run `indus.exe`

### Quick Start Examples

```bash
# Check version
version

# Test an API
http get https://api.github.com/users/hari7261

# Initialize a new project
init --name my-api-project --dir ~/projects

# Navigate directories
cd Documents
pwd
cd ~

# Run concurrent workload test
run --workers 10 --tasks 100

# Use any system command
ipconfig /all
ping google.com
git log --oneline
docker ps -a

# Clear screen
clear

# Exit
exit
```

---

## 📚 Documentation

- **[Installation Guide](README-INSTALL.md)** - Detailed installation and setup
- **[User Guide](GUIDE.md)** - Complete command reference with examples
- **[Capabilities](CAPABILITIES.md)** - Full feature list and use cases

---

## 🎯 Use Cases

### 1. API Development & Testing
```bash
# Test your API endpoints
http get http://localhost:8080/api/users
http post http://localhost:8080/api/users '{"name":"Test User"}'

# Check external APIs
http get https://api.github.com/repos/golang/go
```

### 2. DevOps & CI/CD
```bash
# Health checks
http get https://myapp.com/health

# System diagnostics
ipconfig
netstat -an
ping myserver.com
```

### 3. Project Setup
```bash
# Initialize multiple projects
init --name backend-api --dir ~/projects
init --name frontend-app --dir ~/projects
cd ~/projects/backend-api
```

### 4. Concurrent Testing
```bash
# Test worker pools
run --workers 50 --tasks 1000

# Simulate load
run --workers 100 --tasks 5000
```

---

## 🏗️ Project Structure

```
indus/
├── cmd/
│   └── indus-terminal/       # Main terminal application
├── internal/
│   ├── cli/                  # CLI framework
│   ├── commands/             # Command implementations
│   ├── config/               # Configuration loader
│   └── httpclient/           # HTTP client with retry
├── build/
│   └── icon.ico              # Application icon
├── images/                   # Screenshots
├── indus.exe                 # Compiled binary
├── install.bat               # Installer script
├── uninstall.bat             # Uninstaller script
├── build.bat                 # Build script
└── README.md                 # This file
```

---

## 🔧 Configuration

INDUS looks for configuration at:
```
~/.config/indus/config.yaml
```

Override with environment variable:
```bash
set INDUS_CONFIG=C:\path\to\config.yaml
```

---

## 🆚 Comparison

| Feature | INDUS | PowerShell | Git Bash | CMD |
|---------|-------|------------|----------|-----|
| Built-in HTTP Client | ✅ | ❌ | ❌ | ❌ |
| Project Scaffolding | ✅ | ❌ | ❌ | ❌ |
| Concurrent Workloads | ✅ | ❌ | ❌ | ❌ |
| System Commands | ✅ | ✅ | ✅ | ✅ |
| ANSI Colors | ✅ | ✅ | ✅ | ⚠️ |
| Custom Prompt | ✅ | ✅ | ✅ | ❌ |
| Portable | ✅ | ❌ | ❌ | ✅ |
| Zero Dependencies | ✅ | ❌ | ❌ | ✅ |

---

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Report Bugs**: Open an issue with details
2. **Suggest Features**: Share your ideas
3. **Submit PRs**: Fork, code, and submit pull requests
4. **Improve Docs**: Help make documentation better
5. **Share**: Star ⭐ the repo and share with others

### Development Setup

```bash
# Clone repository
git clone https://github.com/hari7261/indus.git
cd indus

# Install dependencies
go mod download

# Build
go build -o indus.exe ./cmd/indus-terminal

# Test
go test ./...
```

---

## 📝 Roadmap

- [ ] Command history with arrow keys
- [ ] Tab completion
- [ ] Configuration file support
- [ ] Plugin system
- [ ] Database commands
- [ ] Cloud provider integrations
- [ ] Scripting capabilities
- [ ] Multi-tab support
- [ ] Themes and customization

---

## 🐛 Troubleshooting

### "indus is not recognized"
- Restart your command prompt after installation
- Or add `%LOCALAPPDATA%\INDUS` to PATH manually

### Icon not showing
- Run `refresh-icon.bat` as administrator
- Or restart Windows Explorer

### Colors not working
- INDUS automatically enables ANSI support
- If issues persist, run in Windows Terminal

See [GUIDE.md](GUIDE.md) for more troubleshooting tips.

---

## 📄 License

Production use - no tutorial restrictions.

---

## 🙏 Credits

**Made with ♥ by [hari7261](https://github.com/hari7261)**

Namaste! 🙏

### Special Thanks

- Go Team for the amazing standard library
- All contributors and users
- The open-source community

---

## 📞 Support

- **GitHub**: [https://github.com/hari7261](https://github.com/hari7261)
- **Issues**: [Report bugs](https://github.com/hari7261/indus/issues)
- **Discussions**: [Join the conversation](https://github.com/hari7261/indus/discussions)

---

<div align="center">

**INDUS Terminal** - A production-grade CLI for the modern developer

Made in India 🇮🇳 with ♥

[⬆ Back to Top](#indus-terminal-)

</div>

## 🚀 Quick Start

### Installation

1. Download the latest release
2. Run `install.bat` as administrator
3. Launch from Desktop or Start Menu

Or use portable mode - just double-click `indus.exe`

See [README-INSTALL.md](README-INSTALL.md) for detailed installation instructions.

### First Run

```bash
# Launch INDUS Terminal
indus

# Try some commands
version
http get https://api.github.com
cd Documents
pwd
```

## ✨ Features

### Interactive Terminal
- Beautiful ASCII banner with Indian flag colors
- Smart prompt showing current directory
- Full shell capabilities like PowerShell/Git Bash
- ANSI color support
- Built-in command history

### Built-in Commands
- `init` - Initialize project structures
- `http` - Make HTTP requests (GET, POST, PUT, DELETE)
- `run` - Simulate concurrent workloads
- `version` - Show version information
- `cd`, `pwd`, `clear`, `exit` - Standard shell commands

### System Integration
- Run ANY Windows command
- Works with git, docker, npm, python, etc.
- Proper PATH integration
- Desktop & Start Menu shortcuts

## 📚 Documentation

- [Installation Guide](README-INSTALL.md) - How to install and use
- [User Guide](GUIDE.md) - Detailed command reference
- [Capabilities](CAPABILITIES.md) - Complete feature list

## 🎯 Use Cases

- API testing & integration
- DevOps & CI/CD workflows
- Microservices testing
- System diagnostics
- Project initialization
- Quick scripting & automation

## 🛠️ Building from Source

### Prerequisites
- Go 1.21 or higher
- rsrc tool: `go install github.com/akavel/rsrc@latest`

### Build

```bash
# Windows
build.bat

# Or manually
rsrc -ico build/icon.ico -o cmd/indus-terminal/rsrc.syso
go build -o indus.exe ./cmd/indus-terminal
```

## 📁 Project Structure

```
indus/
├── cmd/
│   └── indus-terminal/    # Main terminal application
├── internal/
│   ├── cli/               # CLI framework
│   ├── commands/          # Command implementations
│   ├── config/            # Configuration loader
│   └── httpclient/        # HTTP client with retry
├── build/
│   └── icon.ico           # Application icon
├── indus.exe              # Compiled binary
├── install.bat            # Installer script
├── uninstall.bat          # Uninstaller script
└── build.bat              # Build script
```

## 🎨 Screenshots

```
████████████████████████████████████████████████████████████████████████████████████████████████████████████████
████████████████████████████████████████████████████████████████████████████████████████████████████████████████
████████████████████████████████████████████████████████████████████████████████████████████████████████████████
██                                                                                                            ██
██  ██╗███╗   ██╗██████╗ ██╗   ██╗███████╗    Terminal v1.0.0                                            ██
██  ██║████╗  ██║██╔══██╗██║   ██║██╔════╝                                                                    ██
██  ██║██╔██╗ ██║██║  ██║██║   ██║███████╗    Production-Grade Interactive Terminal                       ██
██  ██║██║╚██╗██║██║  ██║██║   ██║╚════██║    for Developers & DevOps                                     ██
██  ██║██║ ╚████║██████╔╝╚██████╔╝███████║                                                                    ██
██  ╚═╝╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚══════╝                                                                    ██
██                                                                                                            ██
████████████████████████████████████████████████████████████████████████████████████████████████████████████████
████████████████████████████████████████████████████████████████████████████████████████████████████████████████
████████████████████████████████████████████████████████████████████████████████████████████████████████████████

    🙏  Namaste! Welcome to INDUS Terminal

    Made with ♥ by hari7261
    GitHub: https://github.com/hari7261
```

## 🤝 Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest features
- Submit pull requests

## 📄 License

Production use - no tutorial restrictions.

## 🙏 Credits

Made with ♥ by [hari7261](https://github.com/hari7261)

Namaste! 🙏

---

**INDUS Terminal** - A production-grade CLI for the modern developer

## Features

- Subcommands with per-command flags
- Proper exit codes and error handling
- Clean stdout/stderr separation
- Context cancellation and signal handling
- Bounded concurrency with worker pools
- Configuration loading
- HTTP client with retry and backoff
- Future-ready architecture

## Installation

```bash
go build -o indus ./cmd/indus
```

## Build with Version Information

```bash
go build -ldflags "\
  -X main.version=1.0.0 \
  -X main.commit=$(git rev-parse HEAD) \
  -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o indus ./cmd/indus
```

## Cross-Compilation Examples

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o indus-linux-amd64 ./cmd/indus

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o indus-darwin-arm64 ./cmd/indus

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o indus-windows-amd64.exe ./cmd/indus
```

## Usage

### Interactive Mode (REPL)

Simply run without arguments to enter interactive mode:
```bash
indus
```

You'll see a beautiful banner with Indian flag colors and an interactive prompt where you can type commands.

### Direct Command Mode

```bash
# Show help
indus help

# Initialize a new project
indus init --name myproject --dir /path/to/parent

# Run workload simulation
indus run --workers 8 --tasks 50

# Make HTTP requests
indus http get https://api.github.com
indus http post https://api.example.com/data '{"key":"value"}'

# Show version
indus version
```

## Configuration

Default config location: `~/.config/indus/config.yaml`

Override with environment variable:
```bash
export INDUS_CONFIG=/custom/path/config.yaml
```

## Architecture

### Project Structure

```
cmd/indus/           - Entry point
internal/cli/        - CLI kernel (app, command interface, errors)
internal/commands/   - Command implementations
internal/config/     - Configuration loading
internal/httpclient/ - HTTP client with retry logic
```

### Design Principles

- No global mutable state
- Explicit dependency injection
- Context propagation for cancellation
- Stdout for machine output, stderr for human logs
- Commands never call os.Exit
- Exit codes mapped from error types

### Exit Codes

- 0: Success
- 1: Internal error
- 2: User error (bad flags, invalid input)
- 130: Canceled (SIGINT/SIGTERM)

## Interactive Features

- **ASCII Art Banner** with Indian flag colors (Saffron, White, Green)
- **REPL Commands**: `help`, `clear`, `exit`, `quit`
- **Colorized Output** for better readability
- **Command History** navigation
- **Graceful Cancellation** with Ctrl+C

## Commands

### init

Initialize a new project structure.

Flags:
- `--name`: Project name (required)
- `--dir`: Target directory (default: current directory)

Output: Machine-readable key=value pairs to stdout

### run

Execute a simulated workload demonstrating:
- Worker pool pattern
- Bounded concurrency
- Fan-out/fan-in
- Graceful cancellation

Flags:
- `--workers`: Number of concurrent workers (default: 4)
- `--tasks`: Total tasks to process (default: 20)

Output: Progress to stderr, results to stdout

### version

Print version information.

Output: Machine-readable version details to stdout

### http

Make HTTP requests with automatic retry and timeout handling.

Subcommands:
- `http get <url>` - Make GET request
- `http post <url> <data>` - Make POST request
- `http put <url> <data>` - Make PUT request
- `http delete <url>` - Make DELETE request

Flags:
- `--headers`: Custom headers in format 'Key:Value,Key2:Value2'

Examples:
```bash
indus http get https://api.github.com
indus http post https://httpbin.org/post '{"test":"data"}'
indus http get https://api.example.com --headers 'Authorization:Bearer token'
```

## Signal Handling

Press Ctrl+C to gracefully cancel operations. The CLI will:
1. Catch SIGINT/SIGTERM
2. Cancel the root context
3. Wait for commands to clean up
4. Exit with code 130

## Future Extensions

This foundation is ready for:
- API client commands ✅ (Already implemented: http get/post/put/delete)
- Agent orchestration
- Pipeline execution
- Resource management
- Plugin system
- Advanced configuration formats
- Database operations
- Cloud provider integrations
- Container orchestration commands

## What Can INDUS Do?

See [CAPABILITIES.md](CAPABILITIES.md) for a comprehensive list of all features and use cases.

**Quick Summary:**
- 🚀 Initialize projects with standard structure
- 🌐 Make HTTP requests (GET, POST, PUT, DELETE) with retry logic
- ⚡ Simulate concurrent workloads with worker pools
- 💻 Run ANY Windows system command (ipconfig, ping, git, docker, etc.)
- 🎨 Beautiful interactive REPL with Indian flag colors
- 📊 Real-time progress tracking and statistics
- 🔧 Perfect for API testing, DevOps, and development workflows

## License

Production use - no tutorial restrictions.
