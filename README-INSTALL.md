# INDUS Terminal - Installation Guide

## 🚀 Quick Install

### Method 1: Automatic Installation (Recommended)

1. **Download** the latest release
2. **Extract** all files to a folder
3. **Right-click** `install.bat` and select **"Run as administrator"**
4. Follow the on-screen instructions

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
