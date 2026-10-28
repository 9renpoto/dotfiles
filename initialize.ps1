[CmdletBinding()]
param(
  [switch]$SkipPackages,
  [switch]$SkipApply,
  [switch]$ConfigureWakaTime,
  [string]$GhqRoot = (Join-Path $HOME 'src')
)
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
function Write-Info { param([string]$Message) Write-Host "==> $Message" }
function Update-ProcessPath {
  $machinePath = [Environment]::GetEnvironmentVariable('Path', 'Machine')
  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  $separator = [IO.Path]::PathSeparator
  $pathEntries = @($machinePath, $userPath, $env:Path) `
    -split [Regex]::Escape($separator) `
    | Where-Object { $_ } `
    | Select-Object -Unique
  $env:Path = $pathEntries -join $separator
}

function Sync-WslGhqRoot {
  param([Parameter(Mandatory)][string]$Root)

  if (-not (Get-Command wsl.exe -ErrorAction SilentlyContinue)) {
    Write-Warning 'WSL was not found; skipped WSL ghq root configuration.'
    return
  }

  $wslRootOutput = & wsl.exe --exec wslpath -a -u $Root 2>$null
  if ($LASTEXITCODE -ne 0 -or -not $wslRootOutput) {
    Write-Warning 'Could not convert the Windows ghq root to a WSL path.'
    return
  }

  $wslRoot = ($wslRootOutput | Select-Object -First 1).Trim()
  & wsl.exe --exec mkdir -p $wslRoot
  if ($LASTEXITCODE -ne 0) {
    Write-Warning "Could not create the WSL ghq root at $wslRoot."
    return
  }

  & wsl.exe --exec git config --global ghq.root $wslRoot
  if ($LASTEXITCODE -ne 0) {
    Write-Warning 'Git was not available in the default WSL distribution; skipped its ghq configuration.'
    return
  }
  Write-Info "WSL ghq root: $wslRoot"
}


Set-Location $PSScriptRoot
if (-not (Get-Command mise -ErrorAction SilentlyContinue)) {
  & winget install --id jdx.mise --exact --accept-package-agreements --accept-source-agreements
  if ($LASTEXITCODE -ne 0) { throw 'mise installation failed' }
  Update-ProcessPath
}
& mise trust
if ($LASTEXITCODE -ne 0) { throw 'mise trust failed' }
# Persist machine overrides in the config used by ordinary apply commands.
$localConfig = Join-Path $PSScriptRoot 'mise.local.toml'
$normalizedRoot = [IO.Path]::GetFullPath($GhqRoot).Replace('\', '/')
New-Item -ItemType Directory -Path $normalizedRoot -Force | Out-Null
if (-not (Test-Path -LiteralPath $localConfig)) {
  [IO.File]::WriteAllText($localConfig, '', [Text.UTF8Encoding]::new($false))
}
& mise config set --file $localConfig --type string vars.ghq_root $normalizedRoot
if ($LASTEXITCODE -ne 0) { throw 'Failed to persist ghq root' }
if ($ConfigureWakaTime) {
  # Process-scoped only; no secret is printed or passed on the command line.
  $key = Read-Host 'WakaTime API key' -AsSecureString
  $ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($key)
  try { $env:WAKATIME_API_KEY = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr) }
  finally { [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr) }
}
try {
  if (-not $SkipPackages) {
    & "$PSScriptRoot/scripts/install-winget-packages.ps1"
    if ($LASTEXITCODE -ne 0) { throw 'Package installation failed' }
    Update-ProcessPath
  }
  if (-not $SkipApply) {
    & mise bootstrap --skip packages
    if ($LASTEXITCODE -ne 0) { throw 'mise bootstrap failed' }
    Sync-WslGhqRoot -Root $normalizedRoot
  }
} finally {
  if ($ConfigureWakaTime) { Remove-Item Env:WAKATIME_API_KEY -ErrorAction SilentlyContinue }
}
