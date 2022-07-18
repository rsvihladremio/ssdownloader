# script/install: Script to install from source, eventually when there 
#                 are releases this will switch to latest release

$ARCH=(Get-CimInstance -ClassName win32_operatingsystem).OSArchitecture
$download=""
$download_folder=""
if ($ARCH -like 'ARM*') { 
   echo "ARM ARCH"
   $download=ssdownloader-windows-arm64.zip
   $download_folder=ssdownloader-windows-arm64
} else { 
   echo "INTEL ARCH" 
   $download=ssdownloader-windows-amd64.zip
   $download_folder=ssdownloader-windows-amd64
}

$url="https://github.com/rsvihladremio/ssdownloader/releases/latest/$download"
Invoke-WebRequest  -Uri $url -OutFile $download -ContentType 'application/octet-stream'

Write-Output "Checking if scoop is installed"
Get-Date

if (Get-Command 'scoop' -errorAction SilentlyContinue) {
    "scoop installed"
} else {
    Write-Output "scoop not found installing"
    Get-Date
    Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
    Invoke-RestMethod get.scoop.sh | Invoke-Expression
}

Write-Output "Checking if unzip is installed"
Get-Date
if (Get-Command 'unzip' -errorAction SilentlyContinue) {
    "unzip installed"
} else {
    Write-Output "unzip not found installing"
    Get-Date
    scoop install unzip
}

unzip .\"$download"
cp $download".."\ssdownloader.exe ..\

Set-Location ..
Remove-Item ".\$download"
