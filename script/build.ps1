# script/build: build binary 

Set-Location "$PSScriptRoot\.."

# this is also set in script/release and is a copy paste
$GIT_SHA=@(git rev-parse --short HEAD)
$VERSION=@(git rev-parse --abbrev-ref HEAD)
$LDFLAGS="-X github.com/rsvihladremio/ssdownloader/cmd.GitSha=$GIT_SHA -X github.com/rsvihladremio/ssdownloader/cmd.Version=$VERSION"
go build -ldflags "$LDFLAGS" -o ./bin/ssdownloader.exe