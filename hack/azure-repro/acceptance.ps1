# Thin Azure wrapper: fetches the release candidate and the portable Windows
# acceptance suite, runs it, and ships the artifacts back to blob storage.
# The assertions themselves live in hack/acceptance/windows-acceptance.ps1 so
# CI and this VM verify exactly the same contract.

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$sasBase = '__ARTIFACT_SAS__'
$variant = '__VARIANT__'
$work = 'C:\ojdm-acceptance'
$staging = 'C:\ojdm-repro'
New-Item -ItemType Directory -Force -Path $staging | Out-Null

function Blob-Uri([string]$name) {
    $parts = $sasBase -split '\?', 2
    return "$($parts[0])/$name`?$($parts[1])"
}
function Get-Artifact([string]$name, [string]$dest) {
    Invoke-WebRequest -Uri (Blob-Uri $name) -OutFile $dest -UseBasicParsing
}
function Put-Artifact([string]$localPath, [string]$name) {
    if (-not (Test-Path $localPath)) { return }
    Invoke-WebRequest -Uri (Blob-Uri $name) -Method Put -InFile $localPath `
        -Headers @{ 'x-ms-blob-type' = 'BlockBlob' } -UseBasicParsing | Out-Null
}

$exe = Join-Path $staging "ojdm-$variant.exe"
$suite = Join-Path $staging 'windows-acceptance.ps1'
Get-Artifact "ojdm-$variant.exe" $exe
Get-Artifact 'windows-acceptance.ps1' $suite

& $suite -Binary $exe -Work $work
$suiteExit = $LASTEXITCODE

foreach ($artifact in @('report.csv', 'logs\debug.log', 'console1.txt')) {
    Put-Artifact (Join-Path $work $artifact) ("results/$variant-" + ($artifact -replace '\\', '-'))
}

Write-Output "SUITEEXIT $suiteExit"
