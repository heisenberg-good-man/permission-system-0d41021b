$ErrorActionPreference = 'Stop'
Set-Location 'd:\code-space\coding-soler\permission-system-0d41021b\backend'
Write-Host 'Building Go backend...'
go build -mod=mod -o server.exe .
Write-Host 'Build complete, starting server on :8080...'
.\server.exe
