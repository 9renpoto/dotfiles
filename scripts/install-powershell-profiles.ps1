$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$documents = [Environment]::GetFolderPath('MyDocuments')
Push-Location $root
try {
  & mise exec -- go run ./cmd/dotfiles windows-profiles --documents $documents
  if ($LASTEXITCODE -ne 0) { throw "Profile setup failed with exit code $LASTEXITCODE." }
} finally {
  Pop-Location
}
