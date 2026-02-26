@echo off
echo Building INDUS Terminal...
echo.

REM Embed icon
rsrc -ico build\icon.ico -o cmd\indus-terminal\rsrc.syso

REM Build with version info
go build -ldflags "-X main.version=1.0.0 -X main.commit=initial -X main.buildTime=%date:~-4%-%date:~3,2%-%date:~0,2%T%time:~0,2%:%time:~3,2%:%time:~6,2%Z" -o indus.exe ./cmd/indus-terminal

echo.
echo Build complete! indus.exe is ready.
echo.
echo To install: Run install.bat as administrator
echo To test: Double-click indus.exe
echo.
pause
