; ============================================================
; INDUS Terminal - Inno Setup Installer Script
; Build command: iscc installer\indus-setup.iss
; ============================================================

#define AppName      "INDUS Terminal"
#define AppVersion   GetVersionNumbersString(".\..\dist\ind.exe")
#define AppPublisher "Hariom Kumar Pandit"
#define AppURL       "https://github.com/hari7261/INDUS"
#define AppExeName   "ind.exe"
#define AppId        "{{B4F2C8D1-3A7E-4F9B-82C6-1D5E3A7F9B2C}"

[Setup]
AppId={#AppId}
AppName={#AppName}
AppVersion={#AppVersion}
AppVerName={#AppName} {#AppVersion}
AppPublisher={#AppPublisher}
AppPublisherURL={#AppURL}
AppSupportURL={#AppURL}/issues
AppUpdatesURL={#AppURL}/releases
DefaultDirName={localappdata}\Programs\INDUS Terminal
DefaultGroupName={#AppName}
DisableProgramGroupPage=no
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ChangesEnvironment=yes
OutputDir=..\dist
OutputBaseFilename=indus-setup
Compression=lzma2/ultra64
SolidCompression=yes
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
WizardStyle=modern
WizardSizePercent=120
SetupIconFile=..\icon.ico
UninstallDisplayIcon={app}\{#AppExeName}
UninstallDisplayName={#AppName}
VersionInfoVersion={#AppVersion}
VersionInfoCompany={#AppPublisher}
VersionInfoDescription={#AppName} Installer
VersionInfoProductName={#AppName}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "Create a &desktop shortcut"; GroupDescription: "Additional shortcuts:"
Name: "startmenuicon"; Description: "Create a Start &Menu shortcut"; GroupDescription: "Additional shortcuts:"
Name: "contextmenu"; Description: "Add ""Open INDUS Terminal here"" to right-click context menu"; GroupDescription: "Shell integration:"
Name: "addtopath"; Description: "Add INDUS to &PATH (use 'ind' from any terminal)"; GroupDescription: "System integration:"

[Files]
Source: "..\dist\ind.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\dist\indus.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\icon.ico"; DestDir: "{app}"; DestName: "indus.ico"; Flags: ignoreversion
Source: "..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist
Source: "..\README.md"; DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist

[Icons]
Name: "{group}\{#AppName}"; Filename: "{app}\{#AppExeName}"; WorkingDir: "{app}"; IconFilename: "{app}\indus.ico"; Tasks: startmenuicon
Name: "{group}\Uninstall {#AppName}"; Filename: "{uninstallexe}"; Tasks: startmenuicon
Name: "{userdesktop}\{#AppName}"; Filename: "{app}\{#AppExeName}"; WorkingDir: "{app}"; IconFilename: "{app}\indus.ico"; Tasks: desktopicon

[Registry]
Root: HKCU; Subkey: "Software\Classes\Directory\shell\INDUS Terminal"; ValueType: string; ValueName: ""; ValueData: "Open INDUS Terminal here"; Flags: uninsdeletekey; Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\shell\INDUS Terminal"; ValueType: string; ValueName: "Icon"; ValueData: "{app}\{#AppExeName},0"; Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\shell\INDUS Terminal\command"; ValueType: string; ValueName: ""; ValueData: """{app}\{#AppExeName}"" --cwd ""%1"""; Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\Background\shell\INDUS Terminal"; ValueType: string; ValueName: ""; ValueData: "Open INDUS Terminal here"; Flags: uninsdeletekey; Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\Background\shell\INDUS Terminal"; ValueType: string; ValueName: "Icon"; ValueData: "{app}\{#AppExeName},0"; Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\Background\shell\INDUS Terminal\command"; ValueType: string; ValueName: ""; ValueData: """{app}\{#AppExeName}"" --cwd ""%V"""; Tasks: contextmenu

[Run]
Filename: "{app}\{#AppExeName}"; Description: "Launch INDUS Terminal"; Flags: nowait postinstall skipifsilent

[Code]
const
  UserEnvironmentKey = 'Environment';
  HWND_BROADCAST = $ffff;
  WM_SETTINGCHANGE = $001A;
  SMTO_ABORTIFHUNG = $0002;

function SendMessageTimeout(hWnd: HWND; Msg: UINT; wParam: WPARAM; lParam: string;
  fuFlags, uTimeout: UINT; var lpdwResult: DWORD): LRESULT;
  external 'SendMessageTimeoutW@user32.dll stdcall';

procedure AddToUserPath(AppDir: string);
var
  OldPath, NewPath: string;
begin
  if not RegQueryStringValue(HKCU, UserEnvironmentKey, 'Path', OldPath) then
    OldPath := '';
  if Pos(LowerCase(AppDir), LowerCase(OldPath)) = 0 then
  begin
    if (Length(OldPath) > 0) and (OldPath[Length(OldPath)] <> ';') then
      NewPath := OldPath + ';' + AppDir
    else
      NewPath := OldPath + AppDir;
    RegWriteStringValue(HKCU, UserEnvironmentKey, 'Path', NewPath);
  end;
end;

procedure RemoveFromUserPath(AppDir: string);
var
  OldPath, NewPath, Segment: string;
  Parts: TStringList;
  I: Integer;
begin
  if not RegQueryStringValue(HKCU, UserEnvironmentKey, 'Path', OldPath) then
    Exit;
  Parts := TStringList.Create;
  try
    Parts.Delimiter := ';';
    Parts.DelimitedText := OldPath;
    NewPath := '';
    for I := 0 to Parts.Count - 1 do
    begin
      Segment := Trim(Parts[I]);
      if (Segment <> '') and (CompareText(Segment, AppDir) <> 0) then
      begin
        if NewPath <> '' then
          NewPath := NewPath + ';';
        NewPath := NewPath + Segment;
      end;
    end;
    RegWriteStringValue(HKCU, UserEnvironmentKey, 'Path', NewPath);
  finally
    Parts.Free;
  end;
end;

procedure BroadcastEnvironmentChange();
var
  ResultCode: DWORD;
begin
  SendMessageTimeout(HWND_BROADCAST, WM_SETTINGCHANGE, 0, 'Environment', SMTO_ABORTIFHUNG, 5000, ResultCode);
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    if IsTaskSelected('addtopath') then
    begin
      AddToUserPath(ExpandConstant('{app}'));
      BroadcastEnvironmentChange();
    end;
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
  begin
    RemoveFromUserPath(ExpandConstant('{app}'));
    BroadcastEnvironmentChange();
  end;
end;
