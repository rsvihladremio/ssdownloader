
# script/cibuild: Setup environment for CI to run tests. This is primarily
#                 designed to run on the continuous integration server.

Set-Location "$PSScriptRoot\.."

Write-Output "Validating if all dependencies are fullfilled"
Get-Date
.\script\bootstrap.ps1

Write-Output "Tests started at"
Get-Date

.\script\test.ps1


Write-Output "Linting started at"
Get-Date

.\script\lint.ps1

Write-Output "Audit started at"
Get-Date

.\script\audit.ps1