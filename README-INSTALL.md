# INDUS Terminal Installation Guide

## 1. Recommended install (Windows installer)

1. Open the latest release on GitHub.
2. Download `indus-setup-vX.Y.Z-windows-amd64.exe`.
3. Run the installer wizard.
4. Keep `Add INDUS to PATH` checked.
5. Finish setup and launch **INDUS Terminal**.

After install, you can run:

```powershell
ind version
```

## 2. What the installer configures

- Installs binaries into `%LOCALAPPDATA%\Programs\INDUS Terminal`
- Creates desktop and start-menu shortcuts (optional)
- Adds context menu: **Open INDUS Terminal here** (optional)
- Adds install path to user PATH (optional)
- Registers uninstall entry in Windows Apps list

## 3. Portable install (no wizard)

1. Download `ind.exe`.
2. Place it in any folder.
3. Run directly, or add that folder to PATH manually.

## 4. Local production build

```bat
build.bat
```

This now enforces:

- icon embedding (`rsrc`)
- Go unit tests
- installer build validation (unless `SKIP_INSTALLER=1`)

## 5. Signed release configuration

If these GitHub secrets are configured, CI signs Windows artifacts:

- `WINDOWS_SIGN_PFX_BASE64`
- `WINDOWS_SIGN_PFX_PASSWORD`

If either secret is missing, CI still releases unsigned artifacts and logs a warning.

## 6. SmartScreen note

Windows SmartScreen can show warnings for unsigned or low-reputation binaries.
Use a trusted code-signing certificate for production distribution.
