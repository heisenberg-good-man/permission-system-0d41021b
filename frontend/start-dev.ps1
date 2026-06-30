$ErrorActionPreference = 'Continue'
Set-Location 'd:\code-space\coding-soler\permission-system-0d41021b\frontend'

$cacheDir = "d:\code-space\coding-soler\permission-system-0d41021b\frontend\.npm-cache"
$env:npm_config_cache = $cacheDir

Write-Host 'Starting Vite dev server on :5173...'
& "D:\dev-tools\node-v22.22.3-win-x64\node.exe" "D:\dev-tools\node-v22.22.3-win-x64\node_modules\npm\bin\npm-cli.js" run dev -- --cache "$cacheDir"
