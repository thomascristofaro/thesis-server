param ([Parameter(Mandatory)][string]$service, [string]$function = '')

$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
# go env

Set-Location .\$service
if ($function -eq '') {
    $files = Get-ChildItem -Path .\functions -Name    
}
else {
    $files = @($function)
}
foreach ($file in $files) {
    Write-Host $file
    go build -ldflags="-s -w" -o bin/$file functions/$file/main.go
}
Write-Host "END BUILD $service"
