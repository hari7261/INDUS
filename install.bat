@echo off
echo ========================================
echo INDUS Terminal Installer
echo ========================================
echo.

REM Get installation directory
set "INSTALL_DIR=%LOCALAPPDATA%\INDUS"

echo Installing INDUS Terminal to: %INSTALL_DIR%
echo.

REM Create installation directory
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy files
echo Copying files...
copy /Y dist\ind.exe "%INSTALL_DIR%\ind.exe" >nul
copy /Y dist\ind.exe "%INSTALL_DIR%\indus.exe" >nul

REM Add to PATH
echo Adding INDUS to PATH...
setx PATH "%PATH%;%INSTALL_DIR%" >nul

REM Create desktop shortcut
echo Creating desktop shortcut...
powershell -Command "$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%USERPROFILE%\Desktop\INDUS Terminal.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\ind.exe'; $Shortcut.WorkingDirectory = '%USERPROFILE%'; $Shortcut.Save()"

REM Create Start Menu shortcut
echo Creating Start Menu shortcut...
if not exist "%APPDATA%\Microsoft\Windows\Start Menu\Programs\INDUS" mkdir "%APPDATA%\Microsoft\Windows\Start Menu\Programs\INDUS"
powershell -Command "$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%APPDATA%\Microsoft\Windows\Start Menu\Programs\INDUS\INDUS Terminal.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\ind.exe'; $Shortcut.WorkingDirectory = '%USERPROFILE%'; $Shortcut.Save()"

echo.
echo ========================================
echo Installation Complete!
echo ========================================
echo.
echo INDUS Terminal has been installed to:
echo %INSTALL_DIR%
echo.
echo You can now:
echo   1. Double-click "INDUS Terminal" on your Desktop
echo   2. Search for "INDUS" in Start Menu
echo   3. Open any terminal and type: ind
echo.
echo Please restart your command prompt for PATH changes to take effect.
echo.
pause
