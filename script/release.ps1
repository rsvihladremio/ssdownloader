
# script/release: build binaries in all supported platforms and upload them with the gh client

param(
     $VERSION
)

Set-Location "$PSScriptRoot\.."

# this is also set in script/build and is a copy paste
$GIT_SHA=@(git rev-parse --short HEAD)
$LDFLAGS="-X github.com/rsvihladremio/ssdownloader/cmd.GitSha=$GIT_SHA -X github.com/rsvihladremio/ssdownloader/cmd.Version=$VERSION"

git tag $VERSION
git push origin $VERSION

Write-Output "Cleaning bin folder"
Get-Date
.\script\clean

Write-Output "Building linux-amd64"
Get-Date
$Env:GOOS='linux' 
$Env:GOARCH='amd64' 
go build -ldflags "$LDFLAGS" -o ./bin/ssdownloader
zip -j .\bin\ssdownloader-linux-amd64.zip .\bin\ssdownloader
Write-Output "Building linux-arm64"
Get-Date
$Env:GOOS='linux' 
$Env:GOARCH='arm64'
go build -ldflags "$LDFLAGS" -o ./bin/ssdownloader
zip -j .\bin\ssdownloader-linux-arm64.zip .\bin\ssdownloader
Write-Output "Building darwin-os-x-amd64"
Get-Date
$Env:GOOS='darwin' 
$Env:GOARCH='amd64'
go build -ldflags "$LDFLAGS" -o ./bin/ssdownloader
zip -j .\bin\ssdownloader-darwin-amd64.zip .\bin\ssdownloader
Write-Output "Building darwin-os-x-arm64"
Get-Date
$Env:GOOS='darwin' 
$Env:GOARCH='arm64'
go build -ldflags "$LDFLAGS" -o ./bin/ssdownloader
zip -j .\bin\ssdownloader-darwin-arm64.zip .\bin\ssdownloader
Write-Output "Building windows-amd64"
Get-Date
$Env:GOOS='windows' 
$Env:GOARCH='amd64'
go build -ldflags "$LDFLAGS" -o ./bin/ssdownloader.exe
zip -j .\bin\ssdownloader-windows-amd64.zip .\bin\ssdownloader.exe
Write-Output "Building windows-arm64"
Get-Date
$Env:GOOS='windows' 
$Env:GOARCH='arm64'
go build -ldflags "$LDFLAGS" -o ./bin/ssdownloader.exe
zip -j .\bin\ssdownloader-windows-arm64.zip .\bin\ssdownloader.exe

Remove-Item -Path Env:\GOOS
Remove-Item -Path Env:\GOARCH 
gh release create $VERSION --title $VERSION --generate-notes .\bin\ssdownloader-windows-arm64.zip .\bin\ssdownloader-windows-amd64.zip .\bin\ssdownloader-darwin-arm64.zip .\bin\ssdownloader-darwin-amd64.zip .\bin\ssdownloader-linux-arm64.zip .\bin\ssdownloader-linux-amd64.zip 
 
