@echo off
echo ========================================
echo INDUS Terminal Uninstaller
echo ========================================
echo.

set "INSTALL_DIR=%LOCALAPPDATA%\INDUS"

echo Removing INDUS Terminal...
echo.

REM Remove desktop shortcut
if exist "%USERPROFILE%\Desktop\INDUS Terminal.lnk" (
    echo Removing desktop shortcut...
    del "%USERPROFILE%\Desktop\INDUS Terminal.lnk"
)

REM Remove Start Menu shortcut
if exist "%APPDATA%\Microsoft\Windows\Start Menu\Programs\INDUS" (
    echo Removing Start Menu shortcuts...
    rmdir /S /Q "%APPDATA%\Microsoft\Windows\Start Menu\Programs\INDUS"
)

REM Remove installation directory
if exist "%INSTALL_DIR%" (
    echo Removing installation files...
    rmdir /S /Q "%INSTALL_DIR%"
)

echo.
echo ========================================
echo Uninstallation Complete!
echo ========================================
echo.
echo INDUS Terminal has been removed from your system.
echo.
echo Note: You may need to manually remove INDUS from your PATH:
echo   1. Search for "Environment Variables" in Windows
echo   2. Edit PATH variable
echo   3. Remove: %INSTALL_DIR%
echo.
pause
