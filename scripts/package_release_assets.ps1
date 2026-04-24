param(
  [string]$Version = "1.5.5"
)

$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$assetsDir = Join-Path $repoRoot "dist\release-assets"

if (Test-Path $assetsDir) {
  Remove-Item $assetsDir -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $assetsDir | Out-Null

$windowsBinary = Join-Path $repoRoot "dist\indus-v$Version-windows-amd64.exe"
$linuxAmd64Binary = Join-Path $repoRoot "dist\ind-v$Version-linux-amd64"
$linuxArm64Binary = Join-Path $repoRoot "dist\ind-v$Version-linux-arm64"
$darwinAmd64Binary = Join-Path $repoRoot "dist\ind-v$Version-darwin-amd64"
$darwinArm64Binary = Join-Path $repoRoot "dist\ind-v$Version-darwin-arm64"

foreach ($file in @($windowsBinary, $linuxAmd64Binary, $linuxArm64Binary, $darwinAmd64Binary, $darwinArm64Binary)) {
  if (-not (Test-Path $file)) {
    throw "Expected binary missing: $file"
  }
}

$tmpWin = Join-Path $assetsDir "tmp-win"
New-Item -ItemType Directory -Force -Path $tmpWin | Out-Null
Copy-Item $windowsBinary (Join-Path $tmpWin "indus.exe") -Force
Compress-Archive -Path (Join-Path $tmpWin "*") -DestinationPath (Join-Path $assetsDir "indus-v$Version-windows-amd64.zip") -CompressionLevel Optimal -Force
Remove-Item $tmpWin -Recurse -Force

function New-TarAsset {
  param(
    [Parameter(Mandatory = $true)][string]$SourceBinary,
    [Parameter(Mandatory = $true)][string]$AssetName
  )

  $tempDir = Join-Path $assetsDir ("tmp-" + $AssetName)
  New-Item -ItemType Directory -Force -Path $tempDir | Out-Null
  Copy-Item $SourceBinary (Join-Path $tempDir "ind") -Force
  tar -czf (Join-Path $assetsDir $AssetName) -C $tempDir ind
  Remove-Item $tempDir -Recurse -Force
}

New-TarAsset -SourceBinary $linuxAmd64Binary -AssetName "ind-v$Version-linux-amd64.tar.gz"
New-TarAsset -SourceBinary $linuxArm64Binary -AssetName "ind-v$Version-linux-arm64.tar.gz"
New-TarAsset -SourceBinary $darwinAmd64Binary -AssetName "ind-v$Version-darwin-amd64.tar.gz"
New-TarAsset -SourceBinary $darwinArm64Binary -AssetName "ind-v$Version-darwin-arm64.tar.gz"

Copy-Item $windowsBinary (Join-Path $assetsDir "indus-v$Version-windows-amd64.exe") -Force

$latestSummary = Get-ChildItem (Join-Path $repoRoot "dist") -Filter "smoke-summary-*.txt" | Sort-Object LastWriteTime -Descending | Select-Object -First 1
$latestReport = Get-ChildItem (Join-Path $repoRoot "dist") -Filter "smoke-report-*.json" | Sort-Object LastWriteTime -Descending | Select-Object -First 1
if ($latestSummary) {
  Copy-Item $latestSummary.FullName (Join-Path $assetsDir $latestSummary.Name) -Force
}
if ($latestReport) {
  Copy-Item $latestReport.FullName (Join-Path $assetsDir $latestReport.Name) -Force
}

$hashLines = Get-ChildItem $assetsDir -File | Sort-Object Name | ForEach-Object {
  $hash = (Get-FileHash -Algorithm SHA256 $_.FullName).Hash.ToLower()
  "$hash  $($_.Name)"
}
Set-Content -Path (Join-Path $assetsDir "SHA256SUMS.txt") -Value ($hashLines -join "`n")

Get-ChildItem $assetsDir -File | Sort-Object Name | Select-Object Name, Length, LastWriteTime | Format-Table -AutoSize
