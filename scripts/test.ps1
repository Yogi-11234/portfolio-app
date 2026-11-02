$ErrorActionPreference = "Stop"

$env:GOPATH = "D:\tools\gopath"
$env:GOCACHE = "D:\tools\gocache"
$env:GOTMPDIR = "D:\tools\tmp"
$env:TEMP = "D:\tools\tmp"
$env:TMP = "D:\tools\tmp"

$goExe = "D:\tools\go\bin\go.exe"
if (-not (Test-Path $goExe)) {
    $goCmd = Get-Command go -ErrorAction SilentlyContinue
    if ($goCmd) { $goExe = $goCmd.Source }
}

Write-Host "=== Running Backend Unit Tests (Go) ===" -ForegroundColor Cyan
Set-Location "$PSScriptRoot\..\apps\api"
& $goExe test -count=1 ./...
if ($LASTEXITCODE -ne 0) { throw "Backend tests failed" }
Write-Host "Backend tests passed!" -ForegroundColor Green

Write-Host "`n=== Running Frontend TypeScript Check (tsc) ===" -ForegroundColor Cyan
Set-Location "$PSScriptRoot\..\apps\web"
& npx tsc --noEmit
if ($LASTEXITCODE -ne 0) { throw "Frontend TypeScript check failed" }
Write-Host "Frontend type checks passed!" -ForegroundColor Green

Write-Host "`nAll tests and validations passed successfully!" -ForegroundColor Green
