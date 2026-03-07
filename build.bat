@echo off
setlocal EnableDelayedExpansion

echo ========================================================
echo  INDUS Terminal - Build Script
echo ========================================================
echo.

:: ── version info ──────────────────────────────────────────
set VERSION=1.4.0
for /f %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT=%%i
if "!COMMIT!"=="" set COMMIT=none
set BUILD_TIME=%date:~-4%-%date:~3,2%-%date:~0,2%T%time:~0,2%:%time:~3,2%:%time:~6,2%Z

echo Version   : %VERSION%
echo Commit    : %COMMIT%
echo Build time: %BUILD_TIME%
echo.

:: ── create dist folder ────────────────────────────────────
if not exist dist mkdir dist

:: ── embed icon (optional - skip if rsrc not installed) ────
where rsrc >nul 2>&1
if %errorlevel%==0 (
    echo [1/3] Embedding icon...
    rsrc -ico build\icon.ico -o cmd\indus-terminal\rsrc.syso
) else (
    echo [1/3] rsrc not found - skipping icon embed
    echo       Run: go install github.com/akavel/rsrc@latest
)

:: ── build binary ───────────────────────────────────────────
echo [2/3] Building ind.exe...
set LDFLAGS=-s -w -X main.version=%VERSION% -X main.commit=%COMMIT% -X "main.buildTime=%BUILD_TIME%"
go build -ldflags "%LDFLAGS%" -o dist\ind.exe .\cmd\indus-terminal
if %errorlevel% neq 0 (
    echo ERROR: Go build failed!
    pause & exit /b 1
)
copy /Y dist\ind.exe dist\indus.exe >nul
echo       dist\ind.exe    OK
echo       dist\indus.exe  compatibility alias OK

:: ── build installer (requires Inno Setup) ─────────────────
echo [3/3] Building installer...
set ISCC="C:\Program Files (x86)\Inno Setup 6\ISCC.exe"
if exist %ISCC% (
    %ISCC% installer\indus-setup.iss
    if %errorlevel% neq 0 (
        echo WARNING: Installer build failed - binary still OK
    ) else (
        echo       dist\indus-setup.exe  OK
    )
) else (
    echo       Inno Setup not found - skipping installer
    echo       Download: https://jrsoftware.org/isdl.php
)

echo.
echo ========================================================
echo  Build complete!
echo ========================================================
echo.
echo  dist\ind.exe          - Portable binary
echo  dist\indus-setup.exe  - Windows installer wizard
echo.
echo  To release: git tag v%VERSION% ^&^& git push --tags
echo  GitHub Actions will auto-build ^& publish the release.
echo.
pause
endlocal
