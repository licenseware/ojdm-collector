# Windows acceptance suite for a release candidate. Plants the installation
# layouts and failure modes every fix targets, then asserts the delivered
# contract: a complete report, a debug log beside it, and rotation of the
# previous log. Run through run-azure-repro.sh with OJDM_SCRIPT set to this.

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$sasBase = '__ARTIFACT_SAS__'
$variant = '__VARIANT__'
$work = 'C:\ojdm-acceptance'
$results = New-Object System.Collections.ArrayList
$script:pass = 0
$script:fail = 0

Remove-Item -Recurse -Force $work -ErrorAction SilentlyContinue
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
function Check([string]$description, [bool]$ok) {
    if ($ok) { $script:pass++ } else { $script:fail++ }
    $null = $results.Add((("PASS ", "FAIL ")[!$ok]) + $description)
}

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

# ------------------------------------------------------------------ installs
$jre8 = 'C:\Program Files\Java\jre1.8.0_481'
$jdk21 = 'C:\Program Files\Java\jdk-21'
$binTools = 'C:\bin-tools\jdk-21'
$otherDrive = 'D:\Java\jdk-21'

Install-Temurin $jre8 'https://api.adoptium.net/v3/binary/latest/8/ga/windows/x64/jre/hotspot/normal/eclipse?project=jdk' 'jre8'
Install-Temurin $jdk21 'https://api.adoptium.net/v3/binary/latest/21/ga/windows/x64/jdk/hotspot/normal/eclipse?project=jdk' 'jdk21'
foreach ($copy in @($binTools, $otherDrive)) {
    if (-not (Test-Path (Join-Path $copy 'bin\java.exe'))) {
        New-Item -ItemType Directory -Force -Path (Split-Path $copy) | Out-Null
        Copy-Item -Recurse -Force $jdk21 $copy
    }
}
$plantedAll = @($jre8, $jdk21, $binTools, $otherDrive |
        ForEach-Object { Test-Path (Join-Path $_ 'bin\java.exe') }) -notcontains $false
Check "all four installations planted" $plantedAll

# ------------------------------------------- non-permission failure injection
# NTFS accepts a trailing space via \\?\, Win32 lookup strips it: the listing
# yields a name whose Lstat fails with ERROR_FILE_NOT_FOUND. Sorts before Java\.
$broken = 'C:\Program Files\AAA_broken'
Remove-Item -Recurse -Force "\\?\$broken" -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $broken | Out-Null
foreach ($name in @('ghost ', 'ghost.')) { [System.IO.File]::WriteAllText("\\?\$broken\$name", 'x') }

# ------------------------------------------------------------- live jvm
$sleeperSrc = Join-Path $work 'Sleeper.java'
Set-Content -Path $sleeperSrc -Encoding ASCII -Value @'
public class Sleeper {
    public static void main(String[] a) throws Exception { Thread.sleep(600000); }
}
'@
Invoke-Native { & (Join-Path $jdk21 'bin\javac.exe') -d $work $sleeperSrc *>&1 | Out-Null }
$busy = Start-Process -FilePath (Join-Path $jdk21 'bin\java.exe') `
    -ArgumentList '-Xmx64m', '-cp', $work, 'Sleeper' -PassThru -WindowStyle Hidden
Start-Sleep -Seconds 5

# ------------------------------------------------------------------ run 1
$exe = Join-Path $work "ojdm-$variant.exe"
Get-Artifact "ojdm-$variant.exe" $exe

$report = Join-Path $work 'report.csv'
$log = Join-Path $work 'logs\debug.log'
Invoke-Native { & $exe -search-paths 'C:\bin-tools,D:\Java' -output-path $report *>&1 |
        Out-File -FilePath (Join-Path $work 'console1.txt') -Encoding utf8 }

if ($busy) { Stop-Process -Id $busy.Id -Force -ErrorAction SilentlyContinue }

Check "csv report created" (Test-Path $report)
Check "debug log created beside the report" (Test-Path $log)

$rows = @(Import-Csv $report)
$entries = @(Get-Content $log | Where-Object { $_.Trim() } | ForEach-Object { $_ | ConvertFrom-Json })

Check "all four installations reported" ($rows.Count -eq 4)
Check "no installation reported twice" (@($rows | Group-Object JavaHome | Where-Object { $_.Count -gt 1 }).Count -eq 0)
Check "every row resolves its vm shared library" (@($rows | Where-Object { -not $_.DynLibBinPath -or -not (Test-Path $_.DynLibBinPath) }).Count -eq 0)
Check "every jdk row resolves javac.exe" (@($rows | Where-Object { $_.IsJDK -eq 'true' -and -not (Test-Path $_.JavaCBinPath) }).Count -eq 0)
Check "running jvm identified" (@($rows | Where-Object { $_.ProcessRunning -eq 'true' }).Count -ge 1)
Check "logical processors recorded" (@($rows | Where-Object { [int]$_.HostLogicalProcessors -lt 1 }).Count -eq 0)

Check "debug log parses as json throughout" ($entries.Count -gt 0)
Check "log records the search paths" (@($entries | Where-Object { $_.message -eq 'scanning for java installations' }).Count -eq 1)
Check "unreadable path reported and scan survived" (@($entries | Where-Object { $_.message -eq 'skipping unreadable path' -and $_.path -like '*AAA_broken*' }).Count -ge 1)
Check "debug detail present in the file" (@($entries | Where-Object { $_.level -eq 'debug' }).Count -gt 0)
Check "final entry points the operator at the log" ($entries[-1].message -like 'done*')
Check "console stays quiet at info level" (-not (Select-String -Path (Join-Path $work 'console1.txt') -Pattern 'found java file' -Quiet))

# ------------------------------------------------------------------ run 2
Start-Sleep -Seconds 1
Invoke-Native { & $exe -output-path $report *>&1 | Out-File -FilePath (Join-Path $work 'console2.txt') -Encoding utf8 }
$archives = @(Get-ChildItem (Split-Path $log) -Filter 'debug-*.log')
Check "previous log archived on rerun" ($archives.Count -eq 1)
Check "archive kept the earlier run's content" ($archives.Count -eq 1 -and (Select-String -Path $archives[0].FullName -Pattern 'AAA_broken' -Quiet))

# ------------------------------------------------------------------ run 3
$customLog = Join-Path $work 'custom\my.log'
Invoke-Native { & $exe -search-paths 'C:\bin-tools' -output-path (Join-Path $work 'r3.csv') `
        -log-path $customLog -log-level debug *>&1 | Out-File -FilePath (Join-Path $work 'console3.txt') -Encoding utf8 }
Check "-log-path honoured" (Test-Path $customLog)
Check "-log-level=debug surfaces detail on the console" (Select-String -Path (Join-Path $work 'console3.txt') -Pattern 'found java file' -Quiet)

# ------------------------------------------------------------------ report
Put-Artifact $report "results/$variant-report.csv"
Put-Artifact $log "results/$variant-debug.log"

$results | ForEach-Object { Write-Output $_ }
Write-Output ("ROWS total={0} logEntries={1}" -f $rows.Count, $entries.Count)
Write-Output ("SUITE pass={0} fail={1}" -f $script:pass, $script:fail)
