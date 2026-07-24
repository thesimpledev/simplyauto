# Builds the SimplyAuto MSI from bin\simplyauto.exe using WiX 3.
# Usage: .\installer\build-msi.ps1 -Version 1.0.0
param(
    [Parameter(Mandatory = $true)][string]$Version
)

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$wixBin = 'C:\Program Files (x86)\WiX Toolset v3.14\bin'
$exe = Join-Path $root 'bin\simplyauto.exe'
$obj = Join-Path $root 'bin\simplyauto.wixobj'
$msi = Join-Path $root "bin\SimplyAuto-$Version.msi"

if (-not (Test-Path $exe)) { throw "Missing $exe - run the build first" }

& "$wixBin\candle.exe" -nologo -arch x64 "-dVersion=$Version" "-dExePath=$exe" `
    -out $obj (Join-Path $PSScriptRoot 'simplyauto.wxs')
if ($LASTEXITCODE -ne 0) { throw 'candle failed' }

& "$wixBin\light.exe" -nologo -out $msi $obj
if ($LASTEXITCODE -ne 0) { throw 'light failed' }

Remove-Item $obj -ErrorAction SilentlyContinue
Remove-Item ($msi -replace '\.msi$', '.wixpdb') -ErrorAction SilentlyContinue
Write-Host "Built: $msi"
Write-Host "Silent install test: msiexec /i `"$msi`" /qn"
