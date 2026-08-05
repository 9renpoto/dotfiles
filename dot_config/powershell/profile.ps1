if ($global:DotfilesPowerShellProfileLoaded) {
  return
}
$global:DotfilesPowerShellProfileLoaded = $true

if (Get-Command starship -ErrorAction SilentlyContinue) {
  Invoke-Expression (& starship init powershell)
}

if (Get-Command mise -ErrorAction SilentlyContinue) {
  (& mise activate pwsh) | Out-String | Invoke-Expression
}
