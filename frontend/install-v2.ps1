$ErrorActionPreference = 'Continue'
Set-Location 'd:\code-space\coding-soler\permission-system-0d41021b\frontend'

$cacheDir = "d:\code-space\coding-soler\permission-system-0d41021b\frontend\.npm-cache"
if (-not (Test-Path $cacheDir)) { New-Item -ItemType Directory -Path $cacheDir -Force | Out-Null }

$env:npm_config_cache = $cacheDir

"Starting npm install at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" | Set-Content 'result-install.txt'

& "D:\dev-tools\node-v22.22.3-win-x64\node.exe" "D:\dev-tools\node-v22.22.3-win-x64\node_modules\npm\bin\npm-cli.js" install --no-audit --no-fund --cache $cacheDir 1> 'stdout.txt' 2> 'stderr.txt'

"ExitCode: $LASTEXITCODE" | Add-Content 'result-install.txt'

if (Test-Path 'node_modules') {
    $count = (Get-ChildItem 'node_modules' -Directory | Measure-Object).Count
    "SUCCESS: node_modules with $count dirs" | Add-Content 'result-install.txt'
} else {
    "FAILED: no node_modules" | Add-Content 'result-install.txt'
}

"Done at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" | Add-Content 'result-install.txt'
