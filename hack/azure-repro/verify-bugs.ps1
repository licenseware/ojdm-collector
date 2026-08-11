# Plants a matrix of real Java installations on the VM and runs one scan over
# them, so each remaining suspected defect either shows up in the CSV or is
# disproven. Run through run-azure-repro.sh with OJDM_SCRIPT pointed here.
#
# Layouts planted:
#   C:\Program Files\Java\jre1.8.0_481   JRE 8, default location (already there)
#   C:\Program Files\Java\jdk-21         JDK 21, exercises javac / jps / jinfo
#   C:\bin-tools\jdk-21                  root containing "bin" before the real
#                                        bin dir -> processPath truncation
#   D:\Java\jdk-21                       installation on a non-C: drive

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$results = New-Object System.Collections.ArrayList
$sasBase = '__ARTIFACT_SAS__'
$variant = '__VARIANT__'
$work = 'C:\ojdm-repro'
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
function Invoke-Native([scriptblock]$block) {
    $previous = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try { & $block } finally { $ErrorActionPreference = $previous }
}

# ---------------------------------------------------------------- installs
# Self-contained: the VM may have been rebuilt, rebooted or cleaned between
# runs, so never assume an earlier script left an installation behind.
function Install-Temurin([string]$target, [string]$uri, [string]$stageName) {
    if (Test-Path (Join-Path $target 'bin\java.exe')) { return }
    $zip = Join-Path $work "$stageName.zip"
    Invoke-WebRequest -UseBasicParsing -OutFile $zip -Uri $uri
    $stage = Join-Path $work $stageName
    Remove-Item -Recurse -Force $stage -ErrorAction SilentlyContinue
    Expand-Archive -Path $zip -DestinationPath $stage -Force
    New-Item -ItemType Directory -Force -Path (Split-Path $target) | Out-Null
    Remove-Item -Recurse -Force $target -ErrorAction SilentlyContinue
    Move-Item (Get-ChildItem -Directory $stage | Select-Object -First 1).FullName $target
}

$jre8 = 'C:\Program Files\Java\jre1.8.0_481'
Install-Temurin $jre8 'https://api.adoptium.net/v3/binary/latest/8/ga/windows/x64/jre/hotspot/normal/eclipse?project=jdk' 'jre8'

$jdk21 = 'C:\Program Files\Java\jdk-21'
Install-Temurin $jdk21 'https://api.adoptium.net/v3/binary/latest/21/ga/windows/x64/jdk/hotspot/normal/eclipse?project=jdk' 'jdk21'


# ------------------------------------------------- copies in awkward locations
$binTools = 'C:\bin-tools\jdk-21'
if (-not (Test-Path (Join-Path $binTools 'bin\java.exe'))) {
    New-Item -ItemType Directory -Force -Path (Split-Path $binTools) | Out-Null
    Copy-Item -Recurse -Force $jdk21 $binTools
}

$otherDrive = 'D:\Java\jdk-21'
$haveOtherDrive = Test-Path 'D:\'
if ($haveOtherDrive -and -not (Test-Path (Join-Path $otherDrive 'bin\java.exe'))) {
    New-Item -ItemType Directory -Force -Path (Split-Path $otherDrive) | Out-Null
    Copy-Item -Recurse -Force $jdk21 $otherDrive
}
Write-Output "PLANTED jdk21=$(Test-Path $jdk21) binTools=$(Test-Path $binTools) otherDrive=$(Test-Path $otherDrive)"

# ---------------------------------------------------------------- keep a JVM busy
# jps/jinfo only report live JVMs, so compile a trivial class and hold one open
# for the duration of the scan.
$sleeperSrc = Join-Path $work 'Sleeper.java'
Set-Content -Path $sleeperSrc -Encoding ASCII -Value @'
public class Sleeper {
    public static void main(String[] a) throws Exception { Thread.sleep(600000); }
}
'@
Invoke-Native { & (Join-Path $jdk21 'bin\javac.exe') -d $work $sleeperSrc *>&1 | Write-Output }
$busy = Start-Process -FilePath (Join-Path $jdk21 'bin\java.exe') `
    -ArgumentList '-Xmx64m', '-cp', $work, 'Sleeper' -PassThru -WindowStyle Hidden `
    -ErrorAction SilentlyContinue
Start-Sleep -Seconds 5
Write-Output ("BUSY_JVM pid={0} alive={1}" -f $busy.Id, (-not $busy.HasExited))

# ---------------------------------------------------------------- scan
$exe = Join-Path $work "ojdm-$variant.exe"
Get-Artifact "ojdm-$variant.exe" $exe

$csv = Join-Path $work "$variant-matrix.csv"
$log = Join-Path $work "$variant-matrix.log"
Remove-Item $csv, $log -ErrorAction SilentlyContinue

$extra = @('-search-paths', 'C:\bin-tools,D:\Java')
Invoke-Native { & $exe -output-path $csv @extra *>&1 | Out-File -FilePath $log -Encoding utf8 }

if ($busy) { Stop-Process -Id $busy.Id -Force -ErrorAction SilentlyContinue }

Put-Artifact $csv "results/$variant-matrix.csv"
Put-Artifact $log "results/$variant-matrix.log"

# ---------------------------------------------------------------- assertions
$logDir = Join-Path (Split-Path $csv) 'logs'
$debugLog = Join-Path $logDir 'debug.log'
Write-Output ("DEBUGLOG exists={0} entries={1} archives={2}" -f (Test-Path $debugLog),
    @(Get-Content $debugLog -ErrorAction SilentlyContinue).Count,
    @(Get-ChildItem $logDir -Filter 'debug-*.log' -ErrorAction SilentlyContinue).Count)
if (Test-Path $debugLog) {
    Put-Artifact $debugLog "results/$variant-debug.log"
    Get-Content $debugLog | Where-Object { $_ -match '"level":"warn"' } | Select-Object -First 2 |
        ForEach-Object { Write-Output "WARNSAMPLE $_" }
}

$rows = @(Import-Csv $csv)
Write-Output ("ROWS total={0}" -f $rows.Count)
foreach ($r in $rows) {
    $javacOk = if ($r.JavaCBinPath) { Test-Path $r.JavaCBinPath } else { 'n/a' }
    $dllOk = if ($r.DynLibBinPath) { Test-Path $r.DynLibBinPath } else { 'n/a' }
    Write-Output ("ROW home='{0}' isJdk={1} javac='{2}' javacExists={3} dllExists={4} procRunning={5}" -f `
            $r.JavaHome, $r.IsJDK, $r.JavaCBinPath, $javacOk, $dllOk, $r.ProcessRunning)
}
foreach ($expected in @($jdk21, $binTools, $otherDrive)) {
    $needle = ($expected -replace '\\', '/')
    $seen = @($rows | Where-Object { ($_.JavaHome -replace '\\', '/') -like "$needle*" }).Count
    Write-Output ("EXPECT install='{0}' present={1}" -f $expected, ($seen -gt 0))
}
