param(
  [string]$BinaryPath = "",
  [string]$ReportDir = ""
)

$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
if (-not $BinaryPath) {
  $BinaryPath = Join-Path $repoRoot "dist\ind.exe"
}

if (-not (Test-Path $BinaryPath)) {
  throw "Binary not found at $BinaryPath. Build first with: go build -o dist/ind.exe ./cmd/indus-terminal"
}

if (-not $ReportDir) {
  $ReportDir = Join-Path $repoRoot "dist"
}

New-Item -ItemType Directory -Force -Path $ReportDir | Out-Null

$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$smokeRoot = Join-Path $env:TEMP "indus-smoke-$timestamp"
New-Item -ItemType Directory -Force -Path $smokeRoot | Out-Null
New-Item -ItemType Directory -Force -Path (Join-Path $smokeRoot "docs") | Out-Null

Set-Content -Path (Join-Path $smokeRoot "README.md") -Value "# INDUS Smoke Workspace`n"
Set-Content -Path (Join-Path $smokeRoot "docs\index.html") -Value "<!doctype html><title>Smoke Docs</title>"
Set-Content -Path (Join-Path $smokeRoot "seed.json") -Value "{`"hello`":`"world`"}"

$appData = Join-Path $smokeRoot "AppData\Roaming"
$localAppData = Join-Path $smokeRoot "AppData\Local"
New-Item -ItemType Directory -Force -Path $appData | Out-Null
New-Item -ItemType Directory -Force -Path $localAppData | Out-Null

$env:APPDATA = $appData
$env:LOCALAPPDATA = $localAppData

Set-Location $smokeRoot

$commands = @(
  @{ name = "ind about"; args = @("ind", "about") },
  @{ name = "ind doctor"; args = @("ind", "doctor") },
  @{ name = "ind docs"; args = @("ind", "docs") },
  @{ name = "ind scan"; args = @("ind", "scan") },
  @{ name = "ind status"; args = @("ind", "status") },
  @{ name = "ind version"; args = @("ind", "version") },

  @{ name = "ind sys stats"; args = @("ind", "sys", "stats") },
  @{ name = "ind sys info"; args = @("ind", "sys", "info") },
  @{ name = "ind sys clean"; args = @("ind", "sys", "clean") },
  @{ name = "ind sys doctor"; args = @("ind", "sys", "doctor") },
  @{ name = "ind sys watch --interval 200ms --count 2"; args = @("ind", "sys", "watch", "--interval", "200ms", "--count", "2") },

  @{ name = "ind proj create orbit-app --dir ."; args = @("ind", "proj", "create", "orbit-app", "--dir", ".") },
  @{ name = "ind proj init --name smoke-root"; args = @("ind", "proj", "init", "--name", "smoke-root") },
  @{ name = "ind proj list --path ."; args = @("ind", "proj", "list", "--path", ".") },
  @{ name = "ind proj clean ."; args = @("ind", "proj", "clean", ".") },
  @{ name = "ind proj build ."; args = @("ind", "proj", "build", ".") },
  @{ name = "ind proj run ."; args = @("ind", "proj", "run", ".") },

  @{ name = "ind env list"; args = @("ind", "env", "list") },
  @{ name = "ind env set INDUS_PROFILE production"; args = @("ind", "env", "set", "INDUS_PROFILE", "production") },
  @{ name = "ind env export --file indus-env.json"; args = @("ind", "env", "export", "--file", "indus-env.json") },
  @{ name = "ind env unset INDUS_PROFILE"; args = @("ind", "env", "unset", "INDUS_PROFILE") },
  @{ name = "ind env import --file indus-env.json"; args = @("ind", "env", "import", "--file", "indus-env.json") },

  @{ name = "ind fs tree . --depth 2"; args = @("ind", "fs", "tree", ".", "--depth", "2") },
  @{ name = "ind fs find docs --path ."; args = @("ind", "fs", "find", "docs", "--path", ".") },
  @{ name = "ind fs inspect README.md"; args = @("ind", "fs", "inspect", "README.md") },
  @{ name = "ind fs size ."; args = @("ind", "fs", "size", ".") },
  @{ name = "ind fs sync docs docs-copy"; args = @("ind", "fs", "sync", "docs", "docs-copy") },
  @{ name = "ind fs digest docs/index.html"; args = @("ind", "fs", "digest", "docs/index.html") },

  @{ name = "ind net scan example.com"; args = @("ind", "net", "scan", "example.com") },
  @{ name = "ind net pingx example.com --port 443"; args = @("ind", "net", "pingx", "example.com", "--port", "443") },
  @{ name = "ind net trace example.com"; args = @("ind", "net", "trace", "example.com") },
  @{ name = "ind net ports --from 8080 --to 8090"; args = @("ind", "net", "ports", "--from", "8080", "--to", "8090") },
  @{ name = "ind net status --url https://example.com"; args = @("ind", "net", "status", "--url", "https://example.com") },
  @{ name = "ind net fetch https://example.com --method GET"; args = @("ind", "net", "fetch", "https://example.com", "--method", "GET") },

  @{ name = "ind dev bench --command 'ind sys stats' --runs 3"; args = @("ind", "dev", "bench", "--command", "ind sys stats", "--runs", "3") },
  @{ name = "ind dev watch --path . --seconds 1"; args = @("ind", "dev", "watch", "--path", ".", "--seconds", "1") },
  @{ name = "ind dev cache"; args = @("ind", "dev", "cache") },
  @{ name = "ind dev reload"; args = @("ind", "dev", "reload") },
  @{ name = "ind dev debug"; args = @("ind", "dev", "debug") },
  @{ name = "ind dev report --output indus-report.json"; args = @("ind", "dev", "report", "--output", "indus-report.json") },

  @{ name = "ind pkg search aurora"; args = @("ind", "pkg", "search", "aurora") },
  @{ name = "ind pkg install aurora-kit"; args = @("ind", "pkg", "install", "aurora-kit") },
  @{ name = "ind pkg update aurora-kit"; args = @("ind", "pkg", "update", "aurora-kit") },
  @{ name = "ind pkg list"; args = @("ind", "pkg", "list") },
  @{ name = "ind pkg audit"; args = @("ind", "pkg", "audit") },
  @{ name = "ind pkg remove aurora-kit"; args = @("ind", "pkg", "remove", "aurora-kit") },

  @{ name = "ind term clearx"; args = @("ind", "term", "clearx") },
  @{ name = "ind term theme saffron"; args = @("ind", "term", "theme", "saffron") },
  @{ name = "ind term history --limit 10"; args = @("ind", "term", "history", "--limit", "10") },
  @{ name = "ind term speed"; args = @("ind", "term", "speed") },
  @{ name = "ind term reset"; args = @("ind", "term", "reset") },
  @{ name = "ind term doctor"; args = @("ind", "term", "doctor") },

  @{ name = "ind work init orbit-space"; args = @("ind", "work", "init", "orbit-space") },
  @{ name = "ind work list"; args = @("ind", "work", "list") },
  @{ name = "ind work switch orbit-space"; args = @("ind", "work", "switch", "orbit-space") },
  @{ name = "ind work clean ."; args = @("ind", "work", "clean", ".") },
  @{ name = "ind work archive ."; args = @("ind", "work", "archive", ".") },
  @{ name = "ind work pin orbit-space"; args = @("ind", "work", "pin", "orbit-space") }
)

$results = [System.Collections.Generic.List[object]]::new()
$passed = 0
$failed = 0

foreach ($cmd in $commands) {
  $sw = [System.Diagnostics.Stopwatch]::StartNew()
  $rawOutput = (& $BinaryPath @($cmd.args) 2>&1 | Out-String)
  $exitCode = $LASTEXITCODE
  $sw.Stop()

  $ok = $exitCode -eq 0
  if ($ok) { $passed++ } else { $failed++ }

  $results.Add([pscustomobject]@{
    name        = $cmd.name
    args        = $cmd.args
    command     = ($cmd.args -join " ")
    passed      = $ok
    exit_code   = $exitCode
    duration_ms = [int][Math]::Round($sw.Elapsed.TotalMilliseconds)
    output      = $rawOutput.Trim()
  })

  $state = if ($ok) { "PASS" } else { "FAIL" }
  Write-Host ("[{0}] {1} ({2} ms)" -f $state, $cmd.name, [int][Math]::Round($sw.Elapsed.TotalMilliseconds))
}

$summary = [pscustomobject]@{
  generated_at     = (Get-Date).ToString("o")
  binary           = $BinaryPath
  smoke_root       = $smokeRoot
  total_commands   = $commands.Count
  passed           = $passed
  failed           = $failed
  pass_rate        = if ($commands.Count -gt 0) { [Math]::Round(($passed / $commands.Count) * 100, 2) } else { 0 }
}

$summaryPath = Join-Path $ReportDir "smoke-summary-$timestamp.txt"
$jsonPath = Join-Path $ReportDir "smoke-report-$timestamp.json"

$summaryText = @(
  "INDUS release smoke test",
  "generated_at: $($summary.generated_at)",
  "binary: $($summary.binary)",
  "smoke_root: $($summary.smoke_root)",
  "total_commands: $($summary.total_commands)",
  "passed: $($summary.passed)",
  "failed: $($summary.failed)",
  "pass_rate: $($summary.pass_rate)%"
) -join [Environment]::NewLine

Set-Content -Path $summaryPath -Value $summaryText
[pscustomobject]@{
  summary = $summary
  results = $results
} | ConvertTo-Json -Depth 6 | Set-Content -Path $jsonPath

Write-Host ""
Write-Host "Summary file: $summaryPath"
Write-Host "JSON report : $jsonPath"

if ($failed -gt 0) {
  Write-Error "$failed command(s) failed in smoke run."
  exit 1
}

exit 0
