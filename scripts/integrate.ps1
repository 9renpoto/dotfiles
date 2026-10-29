$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$documents = [Environment]::GetFolderPath('MyDocuments')
Push-Location $root
try {
  & mise exec -- go run ./cmd/dotfiles windows-integrate --documents $documents
  if ($LASTEXITCODE -ne 0) { throw "Windows integration failed with exit code $LASTEXITCODE." }
} finally {
  Pop-Location
}
