param(
  [string]$BinaryPath = "",
  [string]$ReportDir = ""
)

$ErrorActionPreference = "Stop"
$PSNativeCommandUseErrorActionPreference = $false

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
if (-not $BinaryPath) {
  $BinaryPath = Join-Path $repoRoot "dist\indus.exe"
}

if (-not (Test-Path $BinaryPath)) {
  throw "Binary not found at $BinaryPath. Build first with: go build -ldflags='-H windowsgui' -o dist/indus.exe ./cmd/indus-terminal"
}
$BinaryPath = (Resolve-Path $BinaryPath).Path

if (-not $ReportDir) {
  $ReportDir = Join-Path $repoRoot "dist"
}
elseif (-not [System.IO.Path]::IsPathRooted($ReportDir)) {
  $ReportDir = Join-Path $repoRoot $ReportDir
}
$ReportDir = [System.IO.Path]::GetFullPath($ReportDir)

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

$smokePort = 19081
$listenerJob = Start-Job -ArgumentList $smokePort -ScriptBlock {
  param($Port)
  $listener = New-Object System.Net.Sockets.TcpListener([System.Net.IPAddress]::Parse("127.0.0.1"), [int]$Port)
  $listener.Start()
  try {
    while ($true) {
      $client = $listener.AcceptTcpClient()
      try {
        $stream = $client.GetStream()
        $buffer = New-Object byte[] 2048
        if ($stream.CanRead) {
          $null = $stream.Read($buffer, 0, $buffer.Length)
        }
        $body = "indus-smoke-ok"
        $response = "HTTP/1.1 200 OK`r`nContent-Type: text/plain`r`nContent-Length: $($body.Length)`r`nConnection: close`r`n`r`n$body"
        $payload = [System.Text.Encoding]::ASCII.GetBytes($response)
        $stream.Write($payload, 0, $payload.Length)
        $stream.Flush()
      } finally {
        if ($stream) { $stream.Close() }
        $client.Close()
      }
    }
  } finally {
    $listener.Stop()
  }
}
Start-Sleep -Milliseconds 300
$listenerReady = $false
for ($i = 0; $i -lt 20; $i++) {
  try {
    $probe = New-Object System.Net.Sockets.TcpClient
    $probe.Connect("127.0.0.1", $smokePort)
    $probe.Close()
    $listenerReady = $true
    break
  } catch {
    Start-Sleep -Milliseconds 200
  }
}
if (-not $listenerReady) {
  throw "Failed to start local smoke listener on port $smokePort"
}

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

  @{ name = "ind net scan"; args = @("ind", "net", "scan") },
  @{ name = "ind net pingx 127.0.0.1 --port $smokePort"; args = @("ind", "net", "pingx", "127.0.0.1", "--port", "$smokePort"); retries = 3 },
  @{ name = "ind net trace localhost"; args = @("ind", "net", "trace", "localhost") },
  @{ name = "ind net ports --from $smokePort --to $smokePort"; args = @("ind", "net", "ports", "--from", "$smokePort", "--to", "$smokePort") },
  @{ name = "ind net status --url http://127.0.0.1:$smokePort"; args = @("ind", "net", "status", "--url", "http://127.0.0.1:$smokePort") },
  @{ name = "ind net fetch http://127.0.0.1:$smokePort --method GET"; args = @("ind", "net", "fetch", "http://127.0.0.1:$smokePort", "--method", "GET") },
  @{ name = "ind http get http://127.0.0.1:$smokePort"; args = @("ind", "http", "get", "http://127.0.0.1:$smokePort") },
  @{ name = "ind net http post http://127.0.0.1:$smokePort --data '{""ok"":true}'"; args = @("ind", "net", "http", "post", "http://127.0.0.1:$smokePort", "--data", '{"ok":true}') },

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
  @{ name = "ind term profile"; args = @("ind", "term", "profile") },
  @{ name = "ind term profile set prompt INDUS-PRO"; args = @("ind", "term", "profile", "set", "prompt", "INDUS-PRO") },
  @{ name = "ind term profile set compact on"; args = @("ind", "term", "profile", "set", "compact", "on") },
  @{ name = "ind term history --limit 10"; args = @("ind", "term", "history", "--limit", "10") },
  @{ name = "ind term speed"; args = @("ind", "term", "speed") },
  @{ name = "ind term reset"; args = @("ind", "term", "reset") },
  @{ name = "ind term doctor"; args = @("ind", "term", "doctor") },

  @{ name = "ind task create smoke-flow --commands 'ind env set INDUS_SMOKE ok && ind env list'"; args = @("ind", "task", "create", "smoke-flow", "--commands", "ind env set INDUS_SMOKE ok && ind env list") },
  @{ name = "ind task list"; args = @("ind", "task", "list") },
  @{ name = "ind task show smoke-flow"; args = @("ind", "task", "show", "smoke-flow") },
  @{ name = "ind task run smoke-flow"; args = @("ind", "task", "run", "smoke-flow") },
  @{ name = "ind task remove smoke-flow"; args = @("ind", "task", "remove", "smoke-flow") },

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

try {
  foreach ($cmd in $commands) {
    $maxAttempts = 1
    if ($cmd.ContainsKey("retries")) {
      $parsedRetries = 1
      if ([int]::TryParse([string]$cmd.retries, [ref]$parsedRetries) -and $parsedRetries -gt 1) {
        $maxAttempts = $parsedRetries
      }
    }

    $sw = [System.Diagnostics.Stopwatch]::StartNew()
    $attempt = 0
    $exitCode = 1
    $rawOutput = ""

    while ($attempt -lt $maxAttempts) {
      $attempt++
      $previousPreference = $ErrorActionPreference
      $ErrorActionPreference = "Continue"
      $rawOutput = (& $BinaryPath @($cmd.args) 2>&1 | Out-String)
      $exitCode = $LASTEXITCODE
      $ErrorActionPreference = $previousPreference

      if ($exitCode -eq 0) {
        break
      }
      if ($attempt -lt $maxAttempts) {
        Start-Sleep -Milliseconds 250
      }
    }
    $sw.Stop()

    $ok = $exitCode -eq 0
    if ($ok) { $passed++ } else { $failed++ }

    $results.Add([pscustomobject]@{
      name        = $cmd.name
      args        = $cmd.args
      command     = ($cmd.args -join " ")
      attempts    = $attempt
      passed      = $ok
      exit_code   = $exitCode
      duration_ms = [int][Math]::Round($sw.Elapsed.TotalMilliseconds)
      output      = $rawOutput.Trim()
    })

    $state = if ($ok) { "PASS" } else { "FAIL" }
    Write-Host ("[{0}] {1} ({2} ms)" -f $state, $cmd.name, [int][Math]::Round($sw.Elapsed.TotalMilliseconds))
  }
} finally {
  if ($listenerJob) {
    Stop-Job -Job $listenerJob -ErrorAction SilentlyContinue | Out-Null
    Remove-Job -Job $listenerJob -Force -ErrorAction SilentlyContinue | Out-Null
  }
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
$latestSummaryPath = Join-Path $ReportDir "smoke-summary.txt"
$latestReportPath = Join-Path $ReportDir "smoke-report.json"

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

Copy-Item -Path $summaryPath -Destination $latestSummaryPath -Force
Copy-Item -Path $jsonPath -Destination $latestReportPath -Force

Write-Host ""
Write-Host "Summary file: $summaryPath"
Write-Host "JSON report : $jsonPath"
Write-Host "Latest      : $latestSummaryPath"
Write-Host "Latest JSON : $latestReportPath"

if ($failed -gt 0) {
  Write-Error "$failed command(s) failed in smoke run."
  exit 1
}

exit 0
