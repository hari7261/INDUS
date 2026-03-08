@echo off
setlocal EnableDelayedExpansion

echo ========================================================
echo INDUS - Production Build
echo ========================================================
echo.

set "VERSION=1.5.0"
for /f %%i in ('git rev-parse --short HEAD 2^>nul') do set "COMMIT=%%i"
if "!COMMIT!"=="" set "COMMIT=none"
for /f %%i in ('powershell -NoProfile -Command "(Get-Date).ToUniversalTime().ToString(\"yyyy-MM-ddTHH:mm:ssZ\")"') do set "BUILD_TIME=%%i"

echo Version   : %VERSION%
echo Commit    : %COMMIT%
echo Build time: %BUILD_TIME%
echo.

if not exist dist mkdir dist

set "RSRC_EXE="
where rsrc.exe >nul 2>&1
if !errorlevel! equ 0 set "RSRC_EXE=rsrc.exe"
if not defined RSRC_EXE if exist "%USERPROFILE%\go\bin\rsrc.exe" set "RSRC_EXE=%USERPROFILE%\go\bin\rsrc.exe"

if not defined RSRC_EXE (
  echo [1/4] Installing rsrc tool...
  go install github.com/akavel/rsrc@latest
  if exist "%USERPROFILE%\go\bin\rsrc.exe" set "RSRC_EXE=%USERPROFILE%\go\bin\rsrc.exe"
  if not defined RSRC_EXE (
    where rsrc.exe >nul 2>&1
    if !errorlevel! equ 0 set "RSRC_EXE=rsrc.exe"
  )
)

if not defined RSRC_EXE (
  echo ERROR: rsrc.exe is required for icon embedding.
  echo Install manually: go install github.com/akavel/rsrc@latest
  exit /b 1
)

if not exist build\icon.ico (
  echo ERROR: build\icon.ico is missing.
  exit /b 1
)

echo [1/4] Embedding icon resource...
"%RSRC_EXE%" -ico build\icon.ico -o cmd\indus\rsrc.syso
if !errorlevel! neq 0 (
  echo ERROR: failed to embed icon resource.
  exit /b 1
)
if not exist cmd\indus\rsrc.syso (
  echo ERROR: cmd\indus\rsrc.syso was not generated.
  exit /b 1
)

echo [2/4] Building dist\ind.exe...
set "LDFLAGS=-s -w -H windowsgui -X main.version=%VERSION% -X main.commit=%COMMIT% -X main.buildTime=%BUILD_TIME%"
go build -ldflags "%LDFLAGS%" -o dist\indus.exe .\cmd\indus
if !errorlevel! neq 0 (
  echo ERROR: go build failed.
  exit /b 1
)
copy /Y dist\ind.exe dist\indus.exe >nul
if !errorlevel! neq 0 (
  echo ERROR: failed to create dist\indus.exe alias.
  exit /b 1
)

echo [3/4] Running unit tests...
go test ./...
if !errorlevel! neq 0 (
  echo ERROR: unit tests failed.
  exit /b 1
)

if /i "%SKIP_INSTALLER%"=="1" (
  echo [4/4] Installer build skipped (SKIP_INSTALLER=1).
) else (
  echo [4/4] Building installer...
  set "ISCC=C:\Program Files (x86)\Inno Setup 6\ISCC.exe"
  if not exist "%ISCC%" (
    echo ERROR: Inno Setup not found at "%ISCC%".
    echo Install from https://jrsoftware.org/isdl.php
    exit /b 1
  )
  "%ISCC%" installer\indus-setup.iss
  if !errorlevel! neq 0 (
    echo ERROR: installer build failed.
    exit /b 1
  )
  if not exist dist\indus-setup.exe (
    echo ERROR: dist\indus-setup.exe was not created.
    exit /b 1
  )
)

if not "%SIGN_PFX%"=="" (
  echo Signing binaries...
  where signtool.exe >nul 2>&1
  if !errorlevel! neq 0 (
    echo ERROR: signtool.exe not found.
    exit /b 1
  )

  signtool sign /fd SHA256 /f "%SIGN_PFX%" /p "%SIGN_PFX_PASSWORD%" /tr http://timestamp.digicert.com /td SHA256 dist\ind.exe
  if !errorlevel! neq 0 exit /b 1
  signtool sign /fd SHA256 /f "%SIGN_PFX%" /p "%SIGN_PFX_PASSWORD%" /tr http://timestamp.digicert.com /td SHA256 dist\indus.exe
  if !errorlevel! neq 0 exit /b 1
  if exist dist\indus-setup.exe (
    signtool sign /fd SHA256 /f "%SIGN_PFX%" /p "%SIGN_PFX_PASSWORD%" /tr http://timestamp.digicert.com /td SHA256 dist\indus-setup.exe
    if !errorlevel! neq 0 exit /b 1
  )
)

echo.
echo ========================================================
echo Build complete.
echo ========================================================
echo dist\ind.exe
echo dist\indus.exe
if exist dist\indus-setup.exe echo dist\indus-setup.exe
echo.
exit /b 0
