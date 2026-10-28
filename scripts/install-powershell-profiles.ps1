$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$documentsDirectory = [Environment]::GetFolderPath('MyDocuments')
$sharedProfilePath = Join-Path $HOME '.config\powershell\profile.ps1'
$loader = @"
`$sharedProfilePath = Join-Path `$HOME '.config\powershell\profile.ps1'
if (Test-Path -LiteralPath `$sharedProfilePath) {
  . `$sharedProfilePath
}
"@
$loader = $loader.Trim() + [Environment]::NewLine
$utf8WithoutBom = [Text.UTF8Encoding]::new($false)

function Install-ProfileLoader {
  param(
    [Parameter(Mandatory)][string]$Path,
    [switch]$CreateIfMissing
  )

  if (-not (Test-Path -LiteralPath $Path)) {
    if ($CreateIfMissing) {
      [IO.File]::WriteAllText($Path, $loader, $utf8WithoutBom)
      return $true
    }
    return $false
  }

  $currentContent = [IO.File]::ReadAllText($Path).Trim()
  if (
    $currentContent.Length -eq 0 -or
    $currentContent -eq $loader.Trim()
  ) {
    [IO.File]::WriteAllText($Path, $loader, $utf8WithoutBom)
    return $true
  } else {
    Write-Warning "Preserving customized PowerShell profile: $Path"
    return $false
  }
}

foreach ($profileDirectoryName in @('WindowsPowerShell', 'PowerShell')) {
  $profileDirectory = Join-Path $documentsDirectory $profileDirectoryName
  New-Item -ItemType Directory -Path $profileDirectory -Force | Out-Null

  $allHostsProfile = Join-Path $profileDirectory 'profile.ps1'
  Install-ProfileLoader -Path $allHostsProfile -CreateIfMissing | Out-Null
}

Write-Host "PowerShell profiles load $sharedProfilePath"
