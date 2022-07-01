# script/update: Update application to run for its current checkout.

Set-Location "$PSScriptRoot\.."

Write-Output "==> Running bootstrap"

.\script\bootstrap

Write-Output "==> Cleaning bin folder"

.\script\clean