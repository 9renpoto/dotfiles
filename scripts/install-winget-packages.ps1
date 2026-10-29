$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
Push-Location $root
try {
  & mise exec -- go run ./cmd/dotfiles packages
  if ($LASTEXITCODE -ne 0) { throw "Package installation failed with exit code $LASTEXITCODE." }
} finally {
  Pop-Location
}
