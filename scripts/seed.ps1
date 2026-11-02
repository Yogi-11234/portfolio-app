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

Set-Location "$PSScriptRoot\..\apps\api"
& $goExe run ./cmd/seed
