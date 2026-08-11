# Windows acceptance suite for a release candidate. Plants the installation
# layouts and failure modes every fix targets, then asserts the delivered
# contract: a complete report, a debug log beside it, and rotation of the
# previous log.
#
# Runs anywhere Windows runs: a CI runner, a throwaway VM, or a developer box.
# Exits non-zero when any assertion fails, which is what gates CI.
#
#   .\windows-acceptance.ps1 -Binary C:\path\ojdm-collector.exe

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)][string]$Binary,
    [string]$Work = 'C:\ojdm-acceptance',
    # A second drive is not present on every runner; the checks that need one
    # are skipped rather than failed when it is missing.
    [string]$SecondDrive = 'D:'
)

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$script:pass = 0
$script:fail = 0

function Check([string]$description, [bool]$ok) {
    if ($ok) { $script:pass++ } else { $script:fail++ }
    Write-Output ((("PASS ", "FAIL ")[!$ok]) + $description)
}

function Invoke-Native([scriptblock]$block) {
    $previous = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try { & $block } finally { $ErrorActionPreference = $previous }
}

function Install-Temurin([string]$target, [string]$uri, [string]$stageName) {
    if (Test-Path (Join-Path $target 'bin\java.exe')) { return }
    $zip = Join-Path $Work "$stageName.zip"
    Invoke-WebRequest -UseBasicParsing -OutFile $zip -Uri $uri
    $stage = Join-Path $Work $stageName
    Remove-Item -Recurse -Force $stage -ErrorAction SilentlyContinue
    Expand-Archive -Path $zip -DestinationPath $stage -Force
    New-Item -ItemType Directory -Force -Path (Split-Path $target) | Out-Null
    Remove-Item -Recurse -Force $target -ErrorAction SilentlyContinue
    Move-Item (Get-ChildItem -Directory $stage | Select-Object -First 1).FullName $target
}

Remove-Item -Recurse -Force $Work -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $Work | Out-Null

# ------------------------------------------------------------------ installs
$jre8 = 'C:\Program Files\Java\jre1.8.0_481'
$jdk21 = 'C:\Program Files\Java\jdk-21'
$binTools = 'C:\bin-tools\jdk-21'
$otherDrive = Join-Path $SecondDrive 'Java\jdk-21'

Install-Temurin $jre8 'https://api.adoptium.net/v3/binary/latest/8/ga/windows/x64/jre/hotspot/normal/eclipse?project=jdk' 'jre8'
Install-Temurin $jdk21 'https://api.adoptium.net/v3/binary/latest/21/ga/windows/x64/jdk/hotspot/normal/eclipse?project=jdk' 'jdk21'

# A root whose own name contains "bin", which used to truncate the java home.
if (-not (Test-Path (Join-Path $binTools 'bin\java.exe'))) {
    New-Item -ItemType Directory -Force -Path (Split-Path $binTools) | Out-Null
    Copy-Item -Recurse -Force $jdk21 $binTools
}

$haveSecondDrive = Test-Path $SecondDrive
if ($haveSecondDrive -and -not (Test-Path (Join-Path $otherDrive 'bin\java.exe'))) {
    New-Item -ItemType Directory -Force -Path (Split-Path $otherDrive) | Out-Null
    Copy-Item -Recurse -Force $jdk21 $otherDrive
}

$expected = @($jre8, $jdk21, $binTools)
$searchPaths = 'C:\bin-tools'
if ($haveSecondDrive) {
    $expected += $otherDrive
    $searchPaths += ",$(Join-Path $SecondDrive 'Java')"
} else {
    Write-Output "SKIP second drive $SecondDrive not present on this host"
}

$planted = @($expected | ForEach-Object { Test-Path (Join-Path $_ 'bin\java.exe') }) -notcontains $false
Check "every installation planted" $planted

# ------------------------------------------- non-permission failure injection
# NTFS accepts a trailing space via \\?\, Win32 lookup strips it: the listing
# yields a name whose Lstat fails with ERROR_FILE_NOT_FOUND. Sorts before Java\.
$broken = 'C:\Program Files\AAA_broken'
Remove-Item -Recurse -Force "\\?\$broken" -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $broken | Out-Null
foreach ($name in @('ghost ', 'ghost.')) { [System.IO.File]::WriteAllText("\\?\$broken\$name", 'x') }

# ------------------------------------------------------------------ live jvm
$sleeperSrc = Join-Path $Work 'Sleeper.java'
Set-Content -Path $sleeperSrc -Encoding ASCII -Value @'
public class Sleeper {
    public static void main(String[] a) throws Exception { Thread.sleep(600000); }
}
'@
Invoke-Native { & (Join-Path $jdk21 'bin\javac.exe') -d $Work $sleeperSrc *>&1 | Out-Null }
$busy = Start-Process -FilePath (Join-Path $jdk21 'bin\java.exe') `
    -ArgumentList '-Xmx64m', '-cp', $Work, 'Sleeper' -PassThru -WindowStyle Hidden
Start-Sleep -Seconds 5

# ------------------------------------------------------------------ run 1
$report = Join-Path $Work 'report.csv'
$log = Join-Path $Work 'logs\debug.log'
$console1 = Join-Path $Work 'console1.txt'
Invoke-Native { & $Binary -search-paths $searchPaths -output-path $report *>&1 |
        Out-File -FilePath $console1 -Encoding utf8 }

if ($busy) { Stop-Process -Id $busy.Id -Force -ErrorAction SilentlyContinue }

Check "csv report created" (Test-Path $report)
Check "debug log created beside the report" (Test-Path $log)

$rows = @(Import-Csv $report)
$entries = @(Get-Content $log | Where-Object { $_.Trim() } | ForEach-Object { $_ | ConvertFrom-Json })

Check "every installation reported" ($rows.Count -eq $expected.Count)
Check "no installation reported twice" (@($rows | Group-Object JavaHome | Where-Object { $_.Count -gt 1 }).Count -eq 0)
Check "every row resolves its vm shared library" (@($rows | Where-Object { -not $_.DynLibBinPath -or -not (Test-Path $_.DynLibBinPath) }).Count -eq 0)
Check "every jdk row resolves javac.exe" (@($rows | Where-Object { $_.IsJDK -eq 'true' -and -not (Test-Path $_.JavaCBinPath) }).Count -eq 0)
Check "installation under a bin-named root reported" (@($rows | Where-Object { $_.JavaHome -like '*bin-tools*' }).Count -eq 1)
Check "running jvm identified" (@($rows | Where-Object { $_.ProcessRunning -eq 'true' }).Count -ge 1)
Check "logical processors recorded" (@($rows | Where-Object { [int]$_.HostLogicalProcessors -lt 1 }).Count -eq 0)

Check "debug log parses as json throughout" ($entries.Count -gt 0)
Check "log records the search paths" (@($entries | Where-Object { $_.message -eq 'scanning for java installations' }).Count -eq 1)
Check "unreadable path reported and scan survived" (@($entries | Where-Object { $_.message -eq 'skipping unreadable path' -and $_.path -like '*AAA_broken*' }).Count -ge 1)
Check "debug detail present in the file" (@($entries | Where-Object { $_.level -eq 'debug' }).Count -gt 0)
Check "final entry points the operator at the log" ($entries[-1].message -like 'done*')
Check "console stays quiet at info level" (-not (Select-String -Path $console1 -Pattern 'found java file' -Quiet))

# ------------------------------------------------------------------ run 2
Start-Sleep -Seconds 1
Invoke-Native { & $Binary -search-paths $searchPaths -output-path $report *>&1 |
        Out-File -FilePath (Join-Path $Work 'console2.txt') -Encoding utf8 }
$archives = @(Get-ChildItem (Split-Path $log) -Filter 'debug-*.log')
Check "previous log archived on rerun" ($archives.Count -eq 1)
Check "archive kept the earlier run's content" ($archives.Count -eq 1 -and (Select-String -Path $archives[0].FullName -Pattern 'AAA_broken' -Quiet))

# ------------------------------------------------------------------ run 3
$customLog = Join-Path $Work 'custom\my.log'
$console3 = Join-Path $Work 'console3.txt'
Invoke-Native { & $Binary -search-paths 'C:\bin-tools' -output-path (Join-Path $Work 'r3.csv') `
        -log-path $customLog -log-level debug *>&1 | Out-File -FilePath $console3 -Encoding utf8 }
Check "-log-path honoured" (Test-Path $customLog)
Check "-log-level=debug surfaces detail on the console" (Select-String -Path $console3 -Pattern 'found java file' -Quiet)

Write-Output ("ROWS total={0} logEntries={1}" -f $rows.Count, $entries.Count)
Write-Output ("SUITE pass={0} fail={1}" -f $script:pass, $script:fail)

if ($script:fail -gt 0) { exit 1 }
