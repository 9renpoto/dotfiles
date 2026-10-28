$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

if (-not (Get-Command winget.exe -ErrorAction SilentlyContinue)) {
  throw 'winget.exe was not found. Install App Installer from Microsoft Store and retry.'
}
$manifest = Join-Path $PSScriptRoot '../winget-packages.json'
& winget.exe import `
  --import-file $manifest `
  --ignore-unavailable `
  --ignore-versions `
  --accept-package-agreements `
  --accept-source-agreements
if ($LASTEXITCODE -ne 0) {
  throw "winget import failed with exit code $LASTEXITCODE."
}
