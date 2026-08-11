# Runs on the Azure Windows VM via `az vm run-command invoke` (as SYSTEM).
# Placeholders __ARTIFACT_SAS__ / __VARIANT__ are substituted by run-azure-repro.sh
# before upload, because run-command parameter passing differs across API versions.

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$sasBase = '__ARTIFACT_SAS__'   # https://<acct>.blob.core.windows.net/<container>?<sas>
$variant = '__VARIANT__'        # label for this build, e.g. buggy / fixed

$results = New-Object System.Collections.ArrayList
$work = 'C:\ojdm-repro'
$javaHome = 'C:\Program Files\Java\jre1.8.0_481'
$churnRoot = 'C:\Program Files\AAA_churn'
$brokenRoot = 'C:\Program Files\AAA_broken'
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

# ---------------------------------------------------------------- Java install
# Temurin 8 (no Oracle login wall) laid out under the exact Oracle JRE 8 path
# from the customer screenshot, so the directory layout the collector walks and
# the bin/server/jvm.dll placement match a real 8u481 install.
if (-not (Test-Path (Join-Path $javaHome 'bin\java.exe'))) {
    $zip = Join-Path $work 'jre8.zip'
    $api = 'https://api.adoptium.net/v3/binary/latest/8/ga/windows/x64/jre/hotspot/normal/eclipse?project=jdk'
    Invoke-WebRequest -Uri $api -OutFile $zip -UseBasicParsing
    $stage = Join-Path $work 'jre8'
    Remove-Item -Recurse -Force $stage -ErrorAction SilentlyContinue
    Expand-Archive -Path $zip -DestinationPath $stage -Force
    $inner = (Get-ChildItem -Directory $stage | Select-Object -First 1).FullName
    New-Item -ItemType Directory -Force -Path (Split-Path $javaHome) | Out-Null
    Move-Item $inner $javaHome
}
# java.exe and the collector write to stderr, which PowerShell turns into a
# terminating NativeCommandError while ErrorActionPreference is 'Stop'.
function Invoke-Native([scriptblock]$block) {
    $previous = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try { & $block } finally { $ErrorActionPreference = $previous }
}

Invoke-Native { & (Join-Path $javaHome 'bin\java.exe') -version 2>&1 | Out-String | Write-Output }

# ---------------------------------------------------------------- collector
$exe = Join-Path $work "ojdm-$variant.exe"
Get-Artifact "ojdm-$variant.exe" $exe

# ---------------------------------------------------------------- error injector
# Reproduces the real trigger: entries that vanish between the directory listing
# and the Lstat of each child, which surfaces as ENOENT (NOT a permission error)
# and therefore propagates out of the Walk callback. On customer machines the
# same class of error comes from EDR agents, updaters, broken junctions and
# OneDrive placeholders. AAA_churn sorts before Java\, so the abort happens
# before the walk ever reaches the installation.
#
# Primary injector, deterministic: NTFS accepts a filename with a trailing
# space when it is created through the \\?\ prefix, but Win32 path
# normalization strips that space on every subsequent lookup. The directory
# listing therefore hands the walk a name whose Lstat fails with
# ERROR_FILE_NOT_FOUND -- not a permission error, so the callback propagates
# it. Same error class as an EDR deleting a file mid-scan, minus the race.
function New-BrokenEntry([string]$root) {
    Remove-Item -Recurse -Force $root -ErrorAction SilentlyContinue
    New-Item -ItemType Directory -Force -Path $root | Out-Null
    foreach ($name in @('ghost ', 'ghost.')) {
        try {
            [System.IO.File]::WriteAllText("\\?\$root\$name", 'x')
        } catch {
            Write-Host "could not create '$name': $($_.Exception.Message)"
        }
    }
    $listed = (Get-ChildItem -LiteralPath $root -Force | Measure-Object).Count
    Write-Host "broken entries planted under ${root}: $listed"
    return $listed
}

$churnScript = {
    param($root)
    $ErrorActionPreference = 'SilentlyContinue'
    New-Item -ItemType Directory -Force -Path $root | Out-Null
    $deadline = (Get-Date).AddMinutes(6)
    while ((Get-Date) -lt $deadline) {
        foreach ($i in 1..80) {
            $d = Join-Path $root "sub$i\bin"
            New-Item -ItemType Directory -Force -Path $d | Out-Null
            foreach ($j in 1..40) { Set-Content -Path (Join-Path $d "f$j.tmp") -Value 'x' }
        }
        Remove-Item -Recurse -Force (Join-Path $root '*')
    }
}

function Start-Churn { return Start-Job -ScriptBlock $churnScript -ArgumentList $churnRoot }
function Stop-Churn($job) {
    if ($job) { Stop-Job $job -ErrorAction SilentlyContinue; Remove-Job $job -Force -ErrorAction SilentlyContinue }
    Remove-Item -Recurse -Force $churnRoot -ErrorAction SilentlyContinue
}

function Invoke-Collector([string]$tag, [string[]]$extraArgs, [bool]$withChurn) {
    $csv = Join-Path $work "$variant-$tag.csv"
    $log = Join-Path $work "$variant-$tag.log"
    Remove-Item $csv, $log -ErrorAction SilentlyContinue

    $job = $null
    if ($withChurn) {
        New-BrokenEntry $brokenRoot | Out-Null
        $job = Start-Churn
        Start-Sleep -Seconds 5
    }
    try {
        Invoke-Native { & $exe -output-path $csv @extraArgs *>&1 | Out-File -FilePath $log -Encoding utf8 }
    } finally {
        if ($withChurn) {
            Stop-Churn $job
            Remove-Item -Recurse -Force "\\?\$brokenRoot" -ErrorAction SilentlyContinue
        }
    }

    $rows = 0
    if (Test-Path $csv) { $rows = [Math]::Max(0, (Import-Csv $csv | Measure-Object).Count) }
    $dyn = 0
    if (Test-Path $csv) { $dyn = (Import-Csv $csv | Where-Object { $_.DynLibBinPath } | Measure-Object).Count }
    $found = 0
    if (Test-Path $log) { $found = (Select-String -Path $log -Pattern '^Found ' -ErrorAction SilentlyContinue | Measure-Object).Count }

    Put-Artifact $csv "results/$variant-$tag.csv"
    Put-Artifact $log "results/$variant-$tag.log"

    $null = $results.Add(("RESULT variant={0} run={1} rows={2} foundLines={3} dynLibRows={4}" -f $variant, $tag, $rows, $found, $dyn))
    return [int]$rows
}

# A: control. No injected errors -> the collector must see the JRE.
$rowsClean = Invoke-Collector 'A-clean' @() $false

# B: same scan with errors injected before Java\ in walk order.
#    Buggy build: Walk aborts -> rows drop (expected 0 java rows).
#    Fixed build: errors logged and skipped -> rows match run A.
$rowsChurn = Invoke-Collector 'B-churn' @() $true

# C: discriminator. Same injected errors, but the search root IS the Java dir,
#    so nothing sorts before it. Rows here + no rows in B proves the walk aborted
#    rather than the JVM being undetectable or java.exe being blocked from exec.
$rowsScoped = Invoke-Collector 'C-churn-scoped' @('-search-paths', 'C:\Program Files\Java') $true

$results | ForEach-Object { Write-Output $_ }
Write-Output ("SUMMARY variant={0} clean={1} churn={2} scoped={3}" -f $variant, $rowsClean, $rowsChurn, $rowsScoped)
if ($rowsClean -gt 0 -and $rowsChurn -eq 0 -and $rowsScoped -gt 0) {
    Write-Output 'VERDICT=REPRODUCED walk aborted on a non-permission error'
} elseif ($rowsClean -gt 0 -and $rowsChurn -eq $rowsClean) {
    Write-Output 'VERDICT=NOT_REPRODUCED scan survived the injected errors'
} else {
    Write-Output 'VERDICT=INCONCLUSIVE inspect the uploaded logs'
}
