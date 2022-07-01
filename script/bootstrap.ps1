# script/bootstrap: Resolve all dependencies that the application requires to
#                   run.

Write-Output "Checking if goveralls is installed"
Get-Date 

if (Get-Command 'goveralls' -errorAction SilentlyContinue) {
    "goveralls installed"
} else {
    Write-Output "goveralls not found installing"
    Get-Date
    go install github.com/mattn/goveralls@latest
}

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

Write-Output "Checking if golaningci-lint is installed"
Get-Date

if (Get-Command 'golangci-lint' -errorAction SilentlyContinue) {
    "golangci-lint installed"
} else {
    Write-Output "golangci-lint not found installing"
    Get-Date
    scoop install golangci-lint
}

Write-Output "Checking if gosec is installed"
Get-Date

if (Get-Command 'gosec' -errorAction SilentlyContinue) {
    "gosec installed"
} else {
    Write-Output "gosec not found installing"
    Get-Date
    go install github.com/securego/gosec/v2/cmd/gosec@latest
}

Write-Output "Checking if gh is installed"
Get-Date

if (Get-Command 'gh' -errorAction SilentlyContinue) {
    "gh installed"
} else {
    Write-Output "gh not found installing"
    Get-Date
    scoop install gh
}


Write-Output "Checking if zip is installed"
Get-Date
if (Get-Command 'zip' -errorAction SilentlyContinue) {
    "zip installed"
} else {
    Write-Output "zip not found installing"
    Get-Date
    scoop install zip
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