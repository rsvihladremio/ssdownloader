# script/install: Script to install from source, eventually when there 
#                 are releases this will switch to latest release

#commenting out until this script makes it into a tag
#$latestTag=@(git describe --tags @(git rev-list --tags --max-count=1))
$latestTag="main"

Invoke-WebRequest -outfile bootstrap.ps1 "https://raw.githubusercontent.com/rsvihladremio/ssdownloader/$latestTag/script/bootstrap.ps1"
.\bootstrap.ps1

#$url="https://github.com/rsvihladremio/ssdownloader/archive/refs/tags/$latestTag.zip"
$url="https://github.com/rsvihladremio/ssdownloader/archive/refs/heads/main.zip"
$fileName="$latestTag.zip"
Invoke-WebRequest  -Uri $url -OutFile $fileName -ContentType 'application/octet-stream'


unzip .\"$latestTag.zip"
#for some reason tag loses v portion
$version=$latestTag.Trim("v"," ")
Set-Location ssdownloader-$version
go build -o ..\ssdownloader

Set-Location ..
Remove-Item .\ssdownloader-$version -Force -Recurse 
Remove-Item .\bootstrap.ps1
