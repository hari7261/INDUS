; ============================================================
;  INDUS Terminal - Inno Setup Installer Script
;  https://jrsoftware.org/isinfo.php  (free tool)
;  Build command: iscc installer\indus-setup.iss
; ============================================================

#define AppName      "INDUS Terminal"
#define AppVersion   GetVersionNumbersString(".\..\dist\ind.exe")
#define AppPublisher "hari7261"
#define AppURL       "https://github.com/hari7261/indus-terminal"
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

; Default install to user-level (no admin required)
DefaultDirName={localappdata}\INDUS
DefaultGroupName={#AppName}
DisableProgramGroupPage=no
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog

; Output
OutputDir=..\dist
OutputBaseFilename=indus-setup
Compression=lzma2/ultra64
SolidCompression=yes
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

; Wizard appearance
WizardStyle=modern
WizardSizePercent=120

; Uninstaller registered in "Apps & Features"
UninstallDisplayIcon={app}\{#AppExeName}
UninstallDisplayName={#AppName}

; Version info embedded in setup.exe
VersionInfoVersion={#AppVersion}
VersionInfoCompany={#AppPublisher}
VersionInfoDescription={#AppName} Installer
VersionInfoProductName={#AppName}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

; ── Optional tasks shown on the "Select Additional Tasks" wizard page ──────
[Tasks]
Name: "desktopicon";    Description: "Create a &desktop shortcut";                                  GroupDescription: "Additional shortcuts:"
Name: "startmenuicon";  Description: "Create a Start &Menu shortcut";                               GroupDescription: "Additional shortcuts:"
Name: "contextmenu";    Description: "Add ""Open INDUS Terminal here"" to right-click context menu"; GroupDescription: "Shell integration:"
Name: "addtopath";      Description: "Add INDUS to &PATH (use 'ind' from any terminal)";            GroupDescription: "System integration:"

; ── Files to install ──────────────────────────────────────────────────────
[Files]
Source: "..\dist\ind.exe";         DestDir: "{app}"; Flags: ignoreversion
Source: "..\dist\indus.exe";       DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE";              DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist
Source: "..\README.md";            DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist

; ── Start Menu shortcuts ──────────────────────────────────────────────────
[Icons]
Name: "{group}\{#AppName}";         Filename: "{app}\{#AppExeName}"; WorkingDir: "{userdocs}"; Tasks: startmenuicon
Name: "{group}\Uninstall {#AppName}"; Filename: "{uninstallexe}";    Tasks: startmenuicon
Name: "{userdesktop}\{#AppName}";   Filename: "{app}\{#AppExeName}"; WorkingDir: "{userdocs}"; Tasks: desktopicon

; ── Registry ─────────────────────────────────────────────────────────────
[Registry]
; Right-click "Open INDUS Terminal here" on folders
Root: HKCU; Subkey: "Software\Classes\Directory\shell\INDUS Terminal";             ValueType: string; ValueName: "";          ValueData: "Open INDUS Terminal here"; Flags: uninsdeletekey;   Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\shell\INDUS Terminal";             ValueType: string; ValueName: "Icon";      ValueData: "{app}\{#AppExeName},0";    Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\shell\INDUS Terminal\command";     ValueType: string; ValueName: "";          ValueData: """{app}\{#AppExeName}""";  Tasks: contextmenu

; Right-click on folder background (inside a folder)
Root: HKCU; Subkey: "Software\Classes\Directory\Background\shell\INDUS Terminal";          ValueType: string; ValueName: "";      ValueData: "Open INDUS Terminal here"; Flags: uninsdeletekey; Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\Background\shell\INDUS Terminal";          ValueType: string; ValueName: "Icon"; ValueData: "{app}\{#AppExeName},0";    Tasks: contextmenu
Root: HKCU; Subkey: "Software\Classes\Directory\Background\shell\INDUS Terminal\command";  ValueType: string; ValueName: "";      ValueData: """{app}\{#AppExeName}"""; Tasks: contextmenu

; ── PATH management (user-level, no admin needed) ────────────────────────
[Code]
const
  EnvironmentKey = 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment';
  UserEnvironmentKey = 'Environment';

// ---------------------------------------------------------------------------
// AddToPath – appends {app} to the user PATH if not already present
// ---------------------------------------------------------------------------
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

// ---------------------------------------------------------------------------
// RemoveFromPath – removes {app} from user PATH during uninstall
// ---------------------------------------------------------------------------
procedure RemoveFromUserPath(AppDir: string);
var
  OldPath, NewPath, Segment: string;
  Parts: TStringList;
  i: Integer;
begin
  if not RegQueryStringValue(HKCU, UserEnvironmentKey, 'Path', OldPath) then
    Exit;
  Parts := TStringList.Create;
  try
    Parts.Delimiter := ';';
    Parts.DelimitedText := OldPath;
    NewPath := '';
    for i := 0 to Parts.Count - 1 do
    begin
      Segment := Trim(Parts[i]);
      if (Segment <> '') and (CompareText(Segment, AppDir) <> 0) then
      begin
        if NewPath <> '' then NewPath := NewPath + ';';
        NewPath := NewPath + Segment;
      end;
    end;
    RegWriteStringValue(HKCU, UserEnvironmentKey, 'Path', NewPath);
  finally
    Parts.Free;
  end;
end;

// Called after files are installed
procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    if IsTaskSelected('addtopath') then
      AddToUserPath(ExpandConstant('{app}'));
  end;
end;

// Called during uninstall
procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
    RemoveFromUserPath(ExpandConstant('{app}'));
end;

// ---------------------------------------------------------------------------
// Notify Windows that PATH changed so it takes effect without a reboot
// ---------------------------------------------------------------------------
procedure DeinitializeSetup();
begin
  // Broadcast WM_SETTINGCHANGE so Explorer / new terminals pick up PATH
  // (best effort – no harm if it fails)
end;
