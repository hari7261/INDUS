# INDUS v1.5.0 - System Passthrough & Toolchain Intelligence

**Release Date:** March 8, 2026  
**Tag:** v1.5.0  
**Target:** master (commit 46fa5a7)

---

## 🚀 Major Features

### **System Command Passthrough**
Unknown commands now pass through to the system shell automatically. INDUS intelligently detects when a command isn't in its registry and executes it via Windows PowerShell/CMD, making it a true **universal terminal**.

```bash
ind npm install express    # Passes through to npm
ind docker ps              # Passes through to Docker
ind python script.py       # Passes through to Python
```

### **30+ Toolchain Detection**
INDUS automatically detects installed developer toolchains and integrates them seamlessly:
- **Languages:** Node.js, Python, Rust, Go, Java, PHP, Ruby, .NET
- **Containers:** Docker, Podman
- **Version Control:** Git, Mercurial, SVN
- **Build Tools:** Make, CMake, Gradle, Maven, Cargo
- **Package Managers:** npm, pip, cargo, composer, gem
- **Cloud CLI:** AWS CLI, Azure CLI, Google Cloud SDK

### **103 Total Commands**
Expanded from 101 to **103 commands** across **11 modules**:
- Core (6) • System (8) • Project (7) • Environment (9)
- Filesystem (5) • Network (6) • Developer (4) • Package (7)
- Terminal (3) • Workspace (3) • **Toolchain (45)** ← NEW!

### **New Toolchain Module**
Intelligent toolchain management with detection, validation, and passthrough:
- `ind toolchain list` - Show all detected toolchains
- `ind toolchain check <tool>` - Verify toolchain availability
- `ind toolchain detect` - Scan system for developer tools
- Plus 42 passthrough commands for common tools

---

## 📚 Premium Documentation Redesign

### **Visual Overhaul**
- **Pure charcoal black** (#0a0a0a) background
- **Saffron accents** (#FF9933) for Indian national identity
- **icon.ico** displayed as logo (navigation + footer)
- **Fully responsive** design (mobile/tablet/desktop)
- **Google Search Console** verified

### **7 Unique Mermaid Architecture Diagrams**
Each page includes interactive architecture visualizations:
1. **Homepage** - 11-module system architecture
2. **Commands** - Command module flow diagram
3. **Versions** - Evolution timeline (v1.3.0→v1.5.0)
4. **v1.5.0** - System passthrough architecture with toolchain detection
5. **v1.4.1** - Registry with production hardening
6. **v1.4.0** - Registry system introduction
7. **v1.3.0** - Pre-registry legacy architecture

### **Pages Updated**
- ✅ **index.html** (22 KB) - Homepage with full feature overview
- ✅ **commands.html** (21 KB) - All 103 commands organized by module
- ✅ **versions.html** (21 KB) - Complete version history with diagrams
- ✅ **v1.5.0.html** (18 KB) - Current release notes
- ✅ **v1.4.1.html** (14 KB) - Production patch notes
- ✅ **v1.4.0.html** (20 KB) - Platform release notes
- ✅ **v1.3.0.html** (18 KB) - Legacy baseline notes

### **Assets Cleanup**
- ✅ Deleted old CSS (style.css, commands.css, index.css)
- ✅ Deleted old JS (script.js, commands.js, index.js)
- ✅ New **main.css** (13 KB) - Premium stylesheet
- ✅ Cleaned all backup/old files

---

## 🔧 Technical Improvements

### **Architecture**
- New `system_exec.go` for passthrough execution
- New `module_toolchain.go` for toolchain detection
- Enhanced registry system with 103 commands
- Improved module architecture with dynamic loading

### **Registry Updates**
- Added 2 new commands (toolchain list, toolchain check)
- Updated command metadata for better help output
- Registry version still 1.4.0 (backward compatible)

### **Build System**
- Added **BUILD-GUIDE.md** with comprehensive build instructions
- Updated `build.bat` with v1.5.0 version
- Included pre-built binaries (indus-terminal.exe, indus-gui.exe)

### **Documentation**
- Updated **CAPABILITIES.md** with v1.5.0 features
- Updated **README.md** with system passthrough examples
- All documentation credits **Hariom Kumar Pandit**

---

## 📦 Installation

### **Windows Installer** (Recommended)
Download and run `indus-setup-v1.5.0.exe` from the release assets below.

### **Manual Installation**
1. Download `indus-terminal.exe`
2. Add to system PATH
3. Run `ind version` to verify

### **From Source**
```bash
git clone https://github.com/hari7261/INDUS.git
cd INDUS
.\build.bat
```

---

## 🎯 Migration from v1.4.x

**100% Backward Compatible** - All v1.4.x commands work identically.

### **What's New for You:**
- Unknown commands now work automatically (no more "command not found")
- Run any system command through INDUS without prefix
- Toolchain detection shows which dev tools you have installed
- Premium documentation with visual architecture diagrams

### **Upgrade Steps:**
```bash
# Download v1.5.0 installer and run
# Or download indus-terminal.exe and replace existing binary

# Verify upgrade
ind version
# Output: INDUS v1.5.0 (103 commands, 11 modules)

# Try system passthrough
ind node --version
ind docker ps
ind git status
```

---

## 🐛 Bug Fixes

- Fixed command parsing for passthrough commands with arguments
- Improved error handling for non-existent toolchains
- Enhanced help output formatting
- Fixed registry loading for new toolchain module

---

## 🙏 Credits

**Built by Hariom Kumar Pandit**  
- GitHub: [@hari7261](https://github.com/hari7261)
- LinkedIn: [Hariom Kumar Pandit](https://linkedin.com/in/hariom-kumar-pandit-2k3)
- Website: [dreamsbuilder.tech](https://www.dreamsbuilder.tech)

---

## 📊 Statistics

- **Lines Changed:** 4,328 insertions, 4,368 deletions
- **Files Modified:** 34 files
- **Documentation Size:** 134 KB (all HTML + CSS)
- **Binary Size:** ~2.5 MB (indus-terminal.exe)
- **Supported Commands:** 103
- **Supported Modules:** 11
- **Toolchain Detected:** 30+

---

## 🔗 Resources

- **Documentation:** https://hari7261.github.io/INDUS/
- **Repository:** https://github.com/hari7261/INDUS
- **Issue Tracker:** https://github.com/hari7261/INDUS/issues
- **Releases:** https://github.com/hari7261/INDUS/releases

---

## 📝 License

MIT License - See [LICENSE](LICENSE) file for details

---

**Full Changelog:** https://github.com/hari7261/INDUS/compare/v1.4.5...v1.5.0
