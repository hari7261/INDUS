# INDUS CLI - Complete Capabilities (v1.5.0)

## 🎯 What Can INDUS Do?

INDUS is a production-grade CLI that combines built-in developer tools with full system command access. **Version 1.5.0** brings enhanced system command execution, development toolchain detection, and production-ready stability.

---

## 🚀 Built-in INDUS Commands (103 Commands)

### 1. Project Initialization
```bash
init --name myproject [--dir /path]
```
**What it does:**
- Creates a complete project structure
- Sets up standard directories (cmd/, internal/, pkg/, config/)
- Generates README.md
- Perfect for starting new Go projects or any structured project

**Use cases:**
- Bootstrap new applications
- Create microservices scaffolding
- Set up API projects
- Initialize CLI tools

---

### 2. HTTP Client (API Testing & Integration)
```bash
# GET requests
http get https://api.github.com
http get https://api.example.com/users --headers 'Authorization:Bearer token123'

# POST requests
http post https://api.example.com/users '{"name":"John","email":"john@example.com"}'
http post https://httpbin.org/post '{"test":"data"}' --headers 'Content-Type:application/json'

# PUT requests
http put https://api.example.com/users/123 '{"name":"Jane"}'

# DELETE requests
http delete https://api.example.com/users/123
```

**What it does:**
- Makes HTTP requests with automatic retry
- Handles timeouts gracefully
- Supports custom headers
- Shows response status and body
- Built-in exponential backoff

**Use cases:**
- Test REST APIs
- Debug API endpoints
- Integration testing
- Quick API exploration
- CI/CD health checks
- Webhook testing
- Microservice communication testing

---

### 3. Concurrent Workload Simulation
```bash
run --workers 8 --tasks 50
```

**What it does:**
- Simulates concurrent processing with worker pools
- Demonstrates bounded concurrency
- Shows real-time progress
- Handles graceful cancellation (Ctrl+C)
- Reports completion statistics

**Use cases:**
- Test concurrent processing patterns
- Benchmark worker pool implementations
- Demonstrate fan-out/fan-in patterns
- Load testing preparation
- Learning concurrency concepts
- Stress testing systems

---

### 4. Version Information
```bash
version
```

**What it does:**
- Shows CLI version
- Displays Git commit hash
- Shows build timestamp
- Useful for debugging and support

---

### 5. Development Toolchain Detection (NEW in v1.5.0)
```bash
# Scan all installed development tools
ind tools scan

# Check if a specific tool is available
ind tools check python
ind tools check docker
ind tools check git
```

**What it does:**
- Detects 30+ development tools and languages
- Shows version information for each tool
- Identifies missing toolchains
- Supports: Python, Node.js, Go, Rust, Java, C++, Docker, Kubernetes, Git, and more
- Helps verify development environment setup

**Use cases:**
- Environment setup validation
- CI/CD environment verification
- Onboarding new developers
- Troubleshooting missing dependencies
- System auditing

**Example Output:**
```
platform=windows/amd64
installed=12
missing=18

INSTALLED:
  Python         Python 3.11.0
  Node.js        v18.16.0
  Git            git version 2.40.0
  Docker         Docker version 24.0.2
  Go             go version go1.21.0
  ...

NOT FOUND:
  Rust
  Java
  Kubernetes
  ...
```

---

## 💻 System Command Passthrough (FULLY FUNCTIONAL in v1.5.0)

INDUS can run **ANY Windows command** directly:

### Network Commands
```bash
ipconfig                    # Network configuration
ipconfig /all              # Detailed network info
ping google.com            # Test connectivity
tracert google.com         # Trace route
netstat -an                # Network statistics
nslookup google.com        # DNS lookup
curl https://api.github.com # HTTP requests
```

### File System Commands
```bash
dir                        # List directory
cd folder                  # Change directory
mkdir newfolder            # Create directory
copy file1.txt file2.txt   # Copy files
move file1.txt folder/     # Move files
del file.txt               # Delete file
type file.txt              # View file content
tree                       # Directory tree
```

### System Information
```bash
systeminfo                 # System details
tasklist                   # Running processes
taskkill /PID 1234         # Kill process
hostname                   # Computer name
whoami                     # Current user
date                       # Current date
time                       # Current time
```

### Development Tools
```bash
git status                 # Git status
git log                    # Git history
go version                 # Go version
node --version             # Node version
npm install                # Install packages
python script.py           # Run Python
docker ps                  # Docker containers
kubectl get pods           # Kubernetes pods
```

### Text Processing
```bash
findstr "pattern" file.txt # Search in files
sort file.txt              # Sort lines
more file.txt              # Page through file
```

---

## 🎨 Interactive REPL Features

### REPL Commands
```bash
help or ?                  # Show all commands
clear or cls               # Clear screen and show banner
exit, quit, or q           # Exit REPL
```

### Features
- Beautiful ASCII banner with Indian flag colors
- Real-time command execution
- Command history navigation
- Colorized output
- Graceful error handling
- Context cancellation support

---

## 🔧 Advanced Use Cases

### 1. API Development & Testing
```bash
# Test your API endpoints
http get http://localhost:8080/api/users
http post http://localhost:8080/api/users '{"name":"Test"}'

# Check external APIs
http get https://api.github.com/users/hari7261
```

### 2. DevOps & CI/CD
```bash
# Health checks
http get https://myapp.com/health

# Deploy verification
http get https://api.production.com/version

# System monitoring
ping myserver.com
netstat -an | findstr "8080"
```

### 3. Microservices Testing
```bash
# Test service endpoints
http get http://service1:8080/health
http get http://service2:8081/metrics
http post http://service3:8082/process '{"data":"test"}'
```

### 4. Learning & Experimentation
```bash
# Test concurrency
run --workers 10 --tasks 100

# Explore APIs
http get https://jsonplaceholder.typicode.com/posts
http get https://api.github.com/repos/golang/go
```

### 5. Quick Scripting & Automation
```bash
# Check multiple endpoints
http get https://api1.com/health
http get https://api2.com/health
http get https://api3.com/health

# Network diagnostics
ping google.com
ping 8.8.8.8
tracert google.com
```

### 6. Project Setup
```bash
# Initialize multiple projects
init --name backend-api --dir ~/projects
init --name frontend-app --dir ~/projects
init --name shared-lib --dir ~/projects
```

---

## 🎯 Real-World Scenarios

### Scenario 1: API Integration Testing
```bash
> http post https://api.example.com/auth '{"username":"test","password":"pass"}'
> http get https://api.example.com/users --headers 'Authorization:Bearer token123'
> http delete https://api.example.com/users/123
```

### Scenario 2: System Diagnostics
```bash
> ipconfig /all
> ping 8.8.8.8
> netstat -an
> systeminfo
```

### Scenario 3: Development Workflow
```bash
> init --name my-new-api
> cd my-new-api
> git init
> http get https://api.github.com/repos/golang/go
```

### Scenario 4: Load Testing Preparation
```bash
> run --workers 50 --tasks 1000
> http get http://localhost:8080/api/test
```

---

## 🌟 Key Advantages

1. **All-in-One Tool**: No need to switch between multiple terminals
2. **Built-in HTTP Client**: No need for curl, Postman, or other tools
3. **System Integration**: Run any Windows command
4. **Developer-Friendly**: Beautiful UI with Indian flag colors
5. **Production-Ready**: Proper error handling, retries, timeouts
6. **Extensible**: Easy to add new commands
7. **Zero Dependencies**: Pure Go standard library
8. **Cross-Platform**: Works on Windows, Linux, macOS

---

## 📊 Performance Features

- **Automatic Retry**: HTTP requests retry with exponential backoff
- **Timeout Handling**: Configurable timeouts (default 30s)
- **Concurrent Processing**: Worker pool pattern with bounded concurrency
- **Graceful Cancellation**: Ctrl+C properly cleans up resources
- **Context Propagation**: All operations respect context cancellation

---

## 🔐 Security Features

- **No Credential Storage**: Doesn't store sensitive data
- **HTTPS Support**: Full TLS/SSL support
- **Header Control**: Custom headers for authentication
- **Safe Defaults**: Reasonable timeouts and retry limits

---

## 🎓 Learning Tool

INDUS is also great for learning:
- HTTP/REST API concepts
- Concurrent programming patterns
- CLI design principles
- Go programming best practices
- System administration commands

---

## 🚀 Future Extensibility

The architecture supports adding:
- Database commands (SQL queries)
- Cloud provider integrations (AWS, Azure, GCP)
- Container orchestration (Docker, Kubernetes)
- File processing commands
- Data transformation tools
- Custom plugins
- Scripting capabilities
- Configuration management

---

## 💡 Pro Tips

1. **Use in CI/CD**: Add INDUS to your pipeline for API testing
2. **Quick Debugging**: Test APIs without leaving terminal
3. **System Admin**: Combine system commands with HTTP checks
4. **Learning**: Experiment with concurrency patterns
5. **Automation**: Chain commands for complex workflows

---

## 📝 Summary

INDUS can:
✅ Initialize projects
✅ Make HTTP requests (GET, POST, PUT, DELETE)
✅ Test APIs with custom headers
✅ Simulate concurrent workloads
✅ Run any Windows system command
✅ Provide network diagnostics
✅ Manage files and directories
✅ Check system information
✅ Support development workflows
✅ Enable quick scripting
✅ Facilitate learning and experimentation

**Made with ♥ by hari7261**
**GitHub: https://github.com/hari7261**
