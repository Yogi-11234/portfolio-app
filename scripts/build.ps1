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

Write-Host "=== 1/2 Building Backend (Go API) ===" -ForegroundColor Cyan
Set-Location "$PSScriptRoot\..\apps\api"
& $goExe build -o ./bin/server.exe ./cmd/server
if ($LASTEXITCODE -ne 0) { throw "Backend build failed" }
Write-Host "Backend build successful -> apps/api/bin/server.exe" -ForegroundColor Green

Write-Host "`n=== 2/2 Building Frontend (React Vite) ===" -ForegroundColor Cyan
Set-Location "$PSScriptRoot\..\apps\web"
npm run build
if ($LASTEXITCODE -ne 0) { throw "Frontend build failed" }
Write-Host "Frontend build successful -> apps/web/dist" -ForegroundColor Green

Write-Host "`nAll builds completed successfully!" -ForegroundColor Green
