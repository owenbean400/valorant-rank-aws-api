# PowerShell Script to build Golang ZIP for AWS Lambda
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o bootstrap
if (Test-Path ".\function.zip") {
    Remove-Item ".\function.zip" -Force
}
Compress-Archive -Path .\bootstrap -DestinationPath function.zip