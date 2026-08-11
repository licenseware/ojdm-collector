# Provisioning only: makes sure a JDK and a JRE exist on this host and prints
# the environment the acceptance suite needs. Everything else, including the
# planted layouts and every assertion, lives in the Go suite.
#
#   .\plant_windows.ps1 | Invoke-Expression
#   go test -tags acceptance .\test\acceptance\ -v
#
# Installs into C:\Program Files\Java, so run it on a throwaway host only.

[CmdletBinding()]
param(
    [string]$JdkPath = 'C:\Program Files\Java\jdk-21',
    [string]$JrePath = 'C:\Program Files\Java\jre1.8.0_481',
    [string]$Staging = "$env:TEMP\ojdm-plant"
)

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

New-Item -ItemType Directory -Force -Path $Staging | Out-Null

function Install-Temurin([string]$target, [string]$uri, [string]$stageName) {
    if (Test-Path (Join-Path $target 'bin\java.exe')) { return }

    $zip = Join-Path $Staging "$stageName.zip"
    Invoke-WebRequest -UseBasicParsing -OutFile $zip -Uri $uri
    $stage = Join-Path $Staging $stageName
    Remove-Item -Recurse -Force $stage -ErrorAction SilentlyContinue
    Expand-Archive -Path $zip -DestinationPath $stage -Force
    New-Item -ItemType Directory -Force -Path (Split-Path $target) | Out-Null
    Remove-Item -Recurse -Force $target -ErrorAction SilentlyContinue
    Move-Item (Get-ChildItem -Directory $stage | Select-Object -First 1).FullName $target
}

# Temurin has no login wall. The JRE is laid out under the Oracle 8 directory
# name so the planted host matches what customers actually run.
Install-Temurin $JdkPath 'https://api.adoptium.net/v3/binary/latest/21/ga/windows/x64/jdk/hotspot/normal/eclipse?project=jdk' 'jdk21'
Install-Temurin $JrePath 'https://api.adoptium.net/v3/binary/latest/8/ga/windows/x64/jre/hotspot/normal/eclipse?project=jdk' 'jre8'

Write-Output "`$env:OJDM_JDK = '$JdkPath'"
Write-Output "`$env:OJDM_JRE = '$JrePath'"

# A second drive exercises installations outside C:, which the default search
# paths never reach. Absent on many hosts, in which case the suite says so.
if (Test-Path 'D:\') {
    New-Item -ItemType Directory -Force -Path 'D:\Java' | Out-Null
    Write-Output "`$env:OJDM_EXTRA_ROOT = 'D:\Java'"
}
