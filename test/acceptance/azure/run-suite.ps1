# Runs on the Azure VM through `az vm run-command invoke`. Fetches the
# collector, the compiled acceptance suite and the planting script, provisions
# the host, then runs the suite. Placeholders are substituted by run.sh.
#
# The suite is shipped pre-compiled, so the VM never needs a Go toolchain.

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$sasBase = '__ARTIFACT_SAS__'
$variant = '__VARIANT__'
$work = 'C:\ojdm-acceptance'
New-Item -ItemType Directory -Force -Path $work | Out-Null

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

$collector = Join-Path $work 'ojdm-collector.exe'
$suite = Join-Path $work 'acceptance.test.exe'
$plant = Join-Path $work 'plant_windows.ps1'

Get-Artifact "ojdm-collector-$variant.exe" $collector
Get-Artifact "acceptance-$variant.test.exe" $suite
Get-Artifact 'plant_windows.ps1' $plant

# The planting script emits the environment the suite expects.
& $plant | Invoke-Expression

$env:OJDM_BINARY = $collector

$output = Join-Path $work "$variant-suite.txt"
$previous = $ErrorActionPreference
$ErrorActionPreference = 'Continue'
# The arguments must be quoted: unquoted, PowerShell splits -test.v at the dot
# and the suite rejects an unknown -test flag.
& $suite '-test.v' '-test.timeout=20m' *>&1 | Tee-Object -FilePath $output
$suiteExit = $LASTEXITCODE
$ErrorActionPreference = $previous

Put-Artifact $output "results/$variant-suite.txt"

Write-Output ("SUITEEXIT {0}" -f $suiteExit)
