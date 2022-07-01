# script/test: Run test suite for application.

Set-Location "$PSScriptRoot\.."

go test -covermode atomic -coverprofile=covprofile ./...
go tool cover -func=covprofile