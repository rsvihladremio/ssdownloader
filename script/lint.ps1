# script/lint: Run gofmt and golangci-lint run

Set-Location "$PSScriptRoot\.."

go fmt ./...

golangci-lint run