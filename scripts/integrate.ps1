$ErrorActionPreference = 'Stop'
& "$PSScriptRoot/install-powershell-profiles.ps1"
# Preserve host settings and gh's credential helpers in the main config.
& git config --global --replace-all include.path '~/.config/git/config'
if ($LASTEXITCODE -ne 0) { throw 'Failed to configure Git include' }
