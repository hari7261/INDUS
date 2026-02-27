# INDUS Terminal — Installation Guide

## 🚀 Install from a Release (like Git)

1. Go to the **[Releases](https://github.com/hari7261/indus-terminal/releases)** page
2. Download **`indus-setup-vX.Y.Z-windows-amd64.exe`**
3. Double-click it — a setup wizard opens
4. Click **Next** through the wizard (choose shortcuts, PATH etc.)
5. Click **Finish** — INDUS is installed

Open any new terminal and type:
```
indus
```

Or use the **INDUS Terminal** shortcut on your Desktop / Start Menu.
Or **right-click any folder** → **"Open INDUS Terminal here"**.

---

## 📦 What the Installer Does

| Step | What happens |
|------|--------------|
| Files | Copies `indus.exe` to `%LOCALAPPDATA%\INDUS\` |
| PATH | Adds that folder to your user PATH (no admin needed) |
| Desktop | Creates **INDUS Terminal** shortcut (optional) |
| Start Menu | Creates **INDUS Terminal** entry (optional) |
| Context menu | Adds **"Open INDUS Terminal here"** to right-click on folders (optional) |
| Uninstaller | Registered in **Apps & Features** for clean removal |

---

## 🗂️ Portable Install (no wizard)

Just download `indus.exe` from Releases, put it anywhere,  
and double-click — no installation needed.

---

## 🔧 Build & Release Yourself

### Build the binary + installer locally

```bat
:: 1. Install Inno Setup (free) from https://jrsoftware.org/isdl.php
:: 2. Run:
build.bat
:: Outputs:  dist\indus.exe  +  dist\indus-setup.exe
```

### Publish a release on GitHub

```bash
git tag v1.2.0
git push --tags
```

GitHub Actions (`.github/workflows/release.yml`) automatically:
1. Builds `indus.exe` (Windows 64-bit)
2. Compiles `indus-setup-v1.2.0-windows-amd64.exe` via Inno Setup
3. Publishes a **GitHub Release** with both files + `checksums.txt`

No manual uploads needed.

---

## 🖥️ Windows Terminal Integration (optional)

Add INDUS as a profile in Windows Terminal:
1. Open Windows Terminal → Settings → **Add a new profile**
2. Set **Command line** to:  `%LOCALAPPDATA%\INDUS\indus.exe`
3. Set **Starting directory** to: `%USERPROFILE%`
4. Give it the name **INDUS Terminal** and save

---

## 🔧 Configuration

INDUS reads config from:
```
~\.config\indus\config.yaml
```

Supported keys:
```yaml
api_timeout = 30    # HTTP timeout in seconds
max_retries = 3     # HTTP retry attempts (GET/HEAD/PUT only)
```

Override location with:
```bat
set INDUS_CONFIG=C:\path\to\config.yaml
```

---

## 🗑️ Uninstall

**Method 1 — Apps & Features (recommended)**
1. Windows Settings → Apps → search **INDUS Terminal** → Uninstall

**Method 2 — manual**
1. Delete `%LOCALAPPDATA%\INDUS\`
2. Remove `%LOCALAPPDATA%\INDUS` from your user PATH
3. Delete Desktop / Start Menu shortcuts

---

## 🐛 Troubleshooting

| Problem | Fix |
|---------|-----|
| `indus` not found after install | Open a **new** terminal window (PATH only loads at start) |
| Colors look broken | Run inside **Windows Terminal** for full ANSI support |
| Installer blocked by Windows | Click **More info → Run anyway** (SmartScreen warning on unsigned exe) |

---

## 🙏 Credits

Made with ♥ by [hari7261](https://github.com/hari7261)


The installer will:
- Install INDUS to `%LOCALAPPDATA%\INDUS`
- Add INDUS to your PATH
- Create Desktop shortcut
- Create Start Menu entry

### Method 2: Portable (No Installation)

1. **Download** `indus.exe`
2. **Place** it anywhere you like
3. **Double-click** to run

## 🎯 How to Use

### Launch INDUS Terminal

After installation, you can launch INDUS in multiple ways:

1. **Desktop Shortcut**: Double-click "INDUS Terminal" on your desktop
2. **Start Menu**: Search for "INDUS" and click
3. **Command Line**: Open any terminal and type `indus`
4. **Direct**: Navigate to installation folder and run `indus.exe`

### First Time Usage

When you open INDUS Terminal, you'll see:
- Beautiful Indian flag banner
- Version information
- Quick start guide

### Basic Commands

```bash
# Navigate directories (like PowerShell/Bash)
cd Documents
cd ~
pwd

# INDUS commands
version
http get https://api.github.com
init --name myproject
run --workers 4 --tasks 20

# System commands (all Windows commands work)
ipconfig
ping google.com
dir
git status
docker ps
npm install

# Terminal commands
help          # Show all commands
clear         # Clear screen
exit          # Exit terminal
```

## 📦 What Gets Installed

```
%LOCALAPPDATA%\INDUS\
├── indus.exe           # Main terminal executable
└── indus-core.exe      # Core CLI (optional)

Desktop\
└── INDUS Terminal.lnk  # Desktop shortcut

Start Menu\Programs\INDUS\
└── INDUS Terminal.lnk  # Start menu shortcut
```

## 🔧 Configuration

INDUS looks for configuration at:
```
~/.config/indus/config.yaml
```

You can override with environment variable:
```bash
set INDUS_CONFIG=C:\path\to\config.yaml
```

## 🗑️ Uninstallation

1. **Right-click** `uninstall.bat` and select **"Run as administrator"**
2. Follow the on-screen instructions

Or manually:
1. Delete `%LOCALAPPDATA%\INDUS` folder
2. Remove desktop shortcut
3. Remove from Start Menu
4. Remove from PATH (optional)

## 🆚 INDUS vs Other Terminals

| Feature | INDUS | PowerShell | Git Bash | CMD |
|---------|-------|------------|----------|-----|
| Built-in HTTP Client | ✅ | ❌ | ❌ | ❌ |
| Project Scaffolding | ✅ | ❌ | ❌ | ❌ |
| Concurrent Workloads | ✅ | ❌ | ❌ | ❌ |
| System Commands | ✅ | ✅ | ✅ | ✅ |
| ANSI Colors | ✅ | ✅ | ✅ | ⚠️ |
| Custom Prompt | ✅ | ✅ | ✅ | ❌ |
| Portable | ✅ | ❌ | ❌ | ✅ |
| Indian Flag Banner | ✅ | ❌ | ❌ | ❌ |

## 💡 Pro Tips

1. **Add to Windows Terminal**: 
   - Open Windows Terminal settings
   - Add new profile pointing to `indus.exe`

2. **Set as Default Terminal**:
   - Right-click any folder
   - "Open INDUS Terminal here" (coming soon)

3. **Use with VS Code**:
   - Set INDUS as integrated terminal
   - Settings → Terminal → External: Windows Exec

4. **Keyboard Shortcuts**:
   - `Ctrl+C` - Cancel current command
   - `Ctrl+D` - Exit terminal
   - `↑/↓` - Command history (coming soon)

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

## 🔄 Updates

To update INDUS:
1. Download the latest version
2. Run `install.bat` again (overwrites old version)

## 📞 Support

- GitHub: https://github.com/hari7261
- Issues: Report bugs on GitHub
- Docs: See GUIDE.md and CAPABILITIES.md

## 🙏 Credits

Made with ♥ by hari7261

Namaste! 🙏
