# INDUS Terminal Build Guide

## Two Build Options

INDUS Terminal can be built in two modes depending on your use case:

### 1. Console Subsystem (Default) - `indus.exe`
**Best for:** Running from existing terminals (CMD, PowerShell, Windows Terminal)

```bash
go build -ldflags="-s -w" -o indus.exe ./cmd/indus
```

**Behavior:**
- Runs inside the terminal you launch it from
- Works perfectly when called from CMD/PowerShell
- Standard console application behavior

**Use when:**
- Using INDUS as a shell within Windows Terminal
- Running from command line
- Scripting and automation

---

### 2. GUI Subsystem (Independent) - `indus-gui.exe`
**Best for:** Standalone desktop application with its own window

```bash
go build -ldflags="-s -w -H windowsgui" -o indus-gui.exe ./cmd/indus-terminal
```

**Behavior:**
- Opens in its OWN independent console window
- No "Command Prompt" or "PowerShell" in title bar
- Completely standalone when double-clicked
- Creates its own console window programmatically

**Use when:**
- Desktop shortcut to launch INDUS directly
- Start Menu launcher
- Standalone terminal experience
- Want INDUS to appear like an independent application

---

## Build Instructions

### Quick Build (Console Version)
```bash
cd "C:\Users\Hariom kumar\Desktop\INDUS"
go build -o indus.exe ./cmd/indus
```

### Quick Build (GUI Version - Standalone)
```bash
cd "C:\Users\Hariom kumar\Desktop\INDUS"
go build -ldflags="-H windowsgui" -o indus-gui.exe ./cmd/indus
```

### Production Build (Both Versions)
```bash
# Build console version
go build -ldflags="-s -w -X main.version=1.5.0 -X main.commit=$(git rev-parse --short HEAD)" -o build/indus.exe ./cmd/indus

# Build GUI version  
go build -ldflags="-s -w -H windowsgui -X main.version=1.5.0 -X main.commit=$(git rev-parse --short HEAD)" -o build/indus-gui.exe ./cmd/indus
```

---

## What's the Difference?

### Console Subsystem (`indus.exe`)
```
User double-clicks indus.exe
  ↓
Windows creates CMD window
  ↓
INDUS runs inside that CMD window
  ↓
User sees CMD with INDUS running
```

### GUI Subsystem (`indus-gui.exe`)
```
User double-clicks indus-gui.exe
  ↓
Windows starts as GUI app (no automatic console)
  ↓
INDUS creates its own console window
  ↓
User sees standalone INDUS Terminal window
```

---

## Recommended Setup

### For Desktop Shortcut:
Use `indus-gui.exe` - gives you a clean, independent window

### For Windows Terminal Profile:
Use `indus.exe` - integrates smoothly with Windows Terminal

### For PATH:
Add both versions to PATH with different names:
```
C:\Users\Hariom kumar\AppData\Local\INDUS\
  ├── indus.exe         (console version)
  └── indus-gui.exe     (standalone version)
```

Then:
- Type `indus` in any terminal → uses console version
- Double-click desktop shortcut → launches GUI version

---

## Build Script

The included `build.bat` script builds production-ready `indus.exe` with version info and embedded icon.

To build:
```batch
.\build.bat
```

This creates `dist\indus.exe` which opens in its own window when double-clicked.

---

## Technical Details

### `-H windowsgui` Flag
- Changes Windows subsystem from "console" to "GUI"
- Prevents automatic console window creation by Windows
- Application must create its own console using `AllocConsole()`
- This is how standalone terminal emulators work (Git Bash, Cmder, etc.)

### Console Management
The `console_windows.go` file handles:
1. Detecting if a console already exists
2. Attaching to parent console (for command-line usage)
3. Creating new console if launched standalone
4. Enabling ANSI colors and proper terminal features

---

## Testing

### Test Console Version:
```bash
# From PowerShell
.\indus.exe version
.\indus.exe tools scan
```

### Test GUI Version:
```bash
# Launch standalone
Start-Process .\indus-gui.exe

# Or double-click the file in Explorer
```

---

## FAQ

**Q: Which version should I use?**  
A: Use GUI version (`indus-gui.exe`) for desktop shortcuts and standalone use. Use console version (`indus.exe`) for command-line and Windows Terminal integration.

**Q: Can I have both versions?**  
A: Yes! Keep both. They're the same code, just compiled differently.

**Q: The GUI version creates a second window when I run it from CMD?**  
A: That's expected. GUI version always creates its own window. Use console version for terminal integration.

**Q: How do I make the installer use GUI version?**  
A: Update `installer/indus-setup.iss` to install `indus-gui.exe` as the main executable for desktop shortcuts.

---

## Installer Configuration

For Inno Setup installer to use standalone GUI version:

```iss
[Files]
Source: "dist\ind.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "dist\indus-gui.exe"; DestDir: "{app}"; Flags: ignoreversion; DestName: "indus-standalone.exe"

[Icons]
Name: "{autoprograms}\INDUS Terminal"; Filename: "{app}\indus-standalone.exe"
Name: "{autodesktop}\INDUS Terminal"; Filename: "{app}\indus-standalone.exe"
```

This way:
- `ind.exe` / `indus.exe` - for command line (in PATH)
- `indus-standalone.exe` - for desktop/Start Menu (GUI version)
