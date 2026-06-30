$ErrorActionPreference = 'Continue'
Set-Location 'd:\code-space\coding-soler\permission-system-0d41021b\frontend'

"Starting npm install at $(Get-Date)" | Out-File 'install-result.txt' -Encoding utf8

& npm.cmd install --no-audit --no-fund *>> 'install-result.txt'

"npm install exit code: $LASTEXITCODE" | Out-File 'install-result.txt' -Append -Encoding utf8

if (Test-Path 'node_modules') {
    $count = (Get-ChildItem 'node_modules' -Directory | Measure-Object).Count
    "node_modules FOUND with $count directories" | Out-File 'install-result.txt' -Append -Encoding utf8
} else {
    "node_modules NOT FOUND" | Out-File 'install-result.txt' -Append -Encoding utf8
}

"Completed at $(Get-Date)" | Out-File 'install-result.txt' -Append -Encoding utf8
