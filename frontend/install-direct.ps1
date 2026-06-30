$ErrorActionPreference = 'Continue'
Set-Location 'd:\code-space\coding-soler\permission-system-0d41021b\frontend'

"Starting npm install at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" | Set-Content 'install-status2.txt'

$nodeExe = "D:\dev-tools\node-v22.22.3-win-x64\node.exe"
$npmCli = "D:\dev-tools\node-v22.22.3-win-x64\node_modules\npm\bin\npm-cli.js"

$psi = New-Object System.Diagnostics.ProcessStartInfo
$cacheDir = "d:\code-space\coding-soler\permission-system-0d41021b\frontend\.npm-cache"
if (-not (Test-Path $cacheDir)) { New-Item -ItemType Directory -Path $cacheDir -Force | Out-Null }
$psi.FileName = $nodeExe
$psi.Arguments = "`"$npmCli`" install --no-audit --no-fund --cache `"$cacheDir`""
$psi.EnvironmentVariables["npm_config_cache"] = $cacheDir
$psi.WorkingDirectory = "d:\code-space\coding-soler\permission-system-0d41021b\frontend"
$psi.RedirectStandardOutput = $true
$psi.RedirectStandardError = $true
$psi.UseShellExecute = $false

$proc = New-Object System.Diagnostics.Process
$proc.StartInfo = $psi
$proc.Start() | Out-Null

$stdout = $proc.StandardOutput.ReadToEnd()
$stderr = $proc.StandardError.ReadToEnd()
$proc.WaitForExit()

"STDOUT: $stdout" | Add-Content 'install-status2.txt'
"STDERR: $stderr" | Add-Content 'install-status2.txt'
"ExitCode: $($proc.ExitCode)" | Add-Content 'install-status2.txt'

if (Test-Path 'node_modules') {
    $count = (Get-ChildItem 'node_modules' -Directory | Measure-Object).Count
    "SUCCESS: node_modules with $count dirs" | Add-Content 'install-status2.txt'
} else {
    "FAILED: no node_modules" | Add-Content 'install-status2.txt'
}

"Done at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" | Add-Content 'install-status2.txt'
