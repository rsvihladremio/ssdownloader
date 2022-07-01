# script/clean: Remove bin folder

Set-Location "$PSScriptRoot\.."

Remove-Item .\bin -Recurse -Force -ErrorAction SilentlyContinue