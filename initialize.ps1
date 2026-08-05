[CmdletBinding()]
param(
  [switch]$SkipPackages,
  [switch]$SkipApply,
  [switch]$ConfigureWakaTime,
  [string]$GhqRoot = (Join-Path $HOME 'src')
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Write-Info {
  param([Parameter(Mandatory)][string]$Message)

  Write-Host "==> $Message" -ForegroundColor Cyan
}

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

function Initialize-ChezmoiConfig {
  $configDirectory = Join-Path $HOME '.config\chezmoi'
  $configPath = Join-Path $configDirectory 'chezmoi.toml'

  if (Test-Path $configPath) {
    return
  }

  Write-Info "Creating default chezmoi config at $configPath"
  New-Item -ItemType Directory -Path $configDirectory -Force | Out-Null
  $config = @'
[data.ssh]
  github_identity_file = "~/.ssh/id_ed25519"
'@
  [IO.File]::WriteAllText(
    $configPath,
    $config + [Environment]::NewLine,
    [Text.UTF8Encoding]::new($false)
  )
}

function Add-WakaTimeConfig {
  param([Parameter(Mandatory)][string]$ConfigPath)

  $configContent = [IO.File]::ReadAllText($ConfigPath)
  if ($configContent -match '(?m)^\s*\[data\.wakatime\]\s*$') {
    Write-Info 'Existing WakaTime configuration was preserved.'
    return
  }

  $secureApiKey = Read-Host 'Enter your WakaTime API key' -AsSecureString
  if ($secureApiKey.Length -eq 0) {
    Write-Warning 'No WakaTime API key was entered; configuration was skipped.'
    return
  }

  $apiKeyPointer = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secureApiKey)
  try {
    $apiKey = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($apiKeyPointer)
    $escapedApiKey = $apiKey.Replace('\', '\\').Replace('"', '\"')
    $newLine = [Environment]::NewLine
    $prefix = if ($configContent.Length -eq 0) {
      ''
    } elseif ($configContent.EndsWith("\n")) {
      $newLine
    } else {
      $newLine + $newLine
    }
    $wakatimeConfig = @"
[data.wakatime]
  api_key = "$escapedApiKey"
"@
    [IO.File]::AppendAllText(
      $ConfigPath,
      $prefix + $wakatimeConfig.Trim() + $newLine,
      [Text.UTF8Encoding]::new($false)
    )
  } finally {
    [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($apiKeyPointer)
  }

  Write-Info 'WakaTime API key was added to the local chezmoi config.'
}

function Set-GhqRoot {
  param([Parameter(Mandatory)][string]$Root)

  $git = Get-Command git.exe -ErrorAction SilentlyContinue
  if (-not $git) {
    throw 'git.exe was not found after package installation.'
  }

  $resolvedRoot = [IO.Path]::GetFullPath($Root)
  New-Item -ItemType Directory -Path $resolvedRoot -Force | Out-Null

  $gitPath = $git.Source
  & $gitPath config --global ghq.root ($resolvedRoot -replace '\\', '/')
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to configure the Windows ghq root at $resolvedRoot."
  }
  Write-Info "Windows ghq root: $resolvedRoot"

  if (-not (Get-Command wsl.exe -ErrorAction SilentlyContinue)) {
    Write-Warning 'WSL was not found; skipped WSL ghq root configuration.'
    return
  }

  $wslRootOutput = & wsl.exe --exec wslpath -a -u $resolvedRoot 2>$null
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

$repositoryRoot = $PSScriptRoot
$packageManifest = Join-Path $repositoryRoot 'winget-packages.json'

if (-not $SkipPackages) {
  if (-not (Get-Command winget.exe -ErrorAction SilentlyContinue)) {
    throw 'winget.exe was not found. Install App Installer from Microsoft Store and retry.'
  }

  Write-Info 'Installing Windows packages with winget'
  & winget.exe import `
    --import-file $packageManifest `
    --ignore-unavailable `
    --ignore-versions `
    --accept-package-agreements `
    --accept-source-agreements

  if ($LASTEXITCODE -ne 0) {
    throw "winget import failed with exit code $LASTEXITCODE."
  }

  Update-ProcessPath
}

if (-not $SkipApply) {
  $chezmoi = Get-Command chezmoi.exe -ErrorAction SilentlyContinue
  if (-not $chezmoi) {
    throw 'chezmoi.exe was not found. Open a new PowerShell session and run .\initialize.ps1 -SkipPackages.'
  }

  Initialize-ChezmoiConfig
  if ($ConfigureWakaTime) {
    $chezmoiConfigPath = Join-Path $HOME '.config\chezmoi\chezmoi.toml'
    Add-WakaTimeConfig -ConfigPath $chezmoiConfigPath
  }
  Write-Info 'Applying dotfiles with chezmoi'
  $chezmoiPath = $chezmoi.Source
  & $chezmoiPath apply --source $repositoryRoot --destination $HOME

  if ($LASTEXITCODE -ne 0) {
    throw "chezmoi apply failed with exit code $LASTEXITCODE."
  }

  Set-GhqRoot -Root $GhqRoot
}

Write-Info 'Windows dotfiles initialization complete'
