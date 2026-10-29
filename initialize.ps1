[CmdletBinding()]
param(
  [switch]$SkipPackages,
  [switch]$SkipApply,
  [switch]$ConfigureWakaTime,
  [string]$GhqRoot = (Join-Path $HOME 'src')
)
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
Set-Location $PSScriptRoot

function Update-ProcessPath {
  $machinePath = [Environment]::GetEnvironmentVariable('Path', 'Machine')
  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  $entries = @($machinePath, $userPath, $env:Path) -split [Regex]::Escape([IO.Path]::PathSeparator) | Where-Object { $_ } | Select-Object -Unique
  $env:Path = $entries -join [IO.Path]::PathSeparator
}

# mise must exist before its managed Go runtime can launch the migrated setup.
if (-not (Get-Command mise -ErrorAction SilentlyContinue)) {
  & winget install --id jdx.mise --exact --accept-package-agreements --accept-source-agreements
  if ($LASTEXITCODE -ne 0) { throw 'mise installation failed' }
  Update-ProcessPath
}
$minimumVersionMatch = [regex]::Match(
  (Get-Content (Join-Path $PSScriptRoot 'mise.toml') -Raw),
  '(?m)^min_version\s*=\s*"([^"]+)"'
)
if (-not $minimumVersionMatch.Success) { throw 'mise minimum version was not found' }
$minimumVersion = $minimumVersionMatch.Groups[1].Value
$installedVersion = ((& mise --version) -split '\s+')[0]
if ($LASTEXITCODE -ne 0) { throw 'Failed to read the installed mise version' }
if ([version]::Parse($installedVersion) -lt [version]::Parse($minimumVersion)) {
  & mise self-update --yes --force --no-plugins $minimumVersion
  if ($LASTEXITCODE -ne 0) { throw "Failed to update mise to $minimumVersion" }
}
& mise trust
if ($LASTEXITCODE -ne 0) { throw 'mise trust failed' }

if ($ConfigureWakaTime) {
  # Keep the secret out of command arguments and logs; the Go child inherits only its environment.
  $key = Read-Host 'WakaTime API key' -AsSecureString
  $ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($key)
  try { $env:WAKATIME_API_KEY = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr) }
  finally { [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr) }
}
try {
  $arguments = @('exec', '--', 'go', 'run', './cmd/dotfiles', 'windows-init', '--ghq-root', $GhqRoot)
  if ($SkipPackages) { $arguments += '--skip-packages' }
  if ($SkipApply) { $arguments += '--skip-apply' }
  & mise @arguments
  if ($LASTEXITCODE -ne 0) { throw "Initialization failed with exit code $LASTEXITCODE." }
} finally {
  if ($ConfigureWakaTime) { Remove-Item Env:WAKATIME_API_KEY -ErrorAction SilentlyContinue }
}
