$ErrorActionPreference = 'Continue'
Set-Location 'd:\code-space\coding-soler\permission-system-0d41021b\frontend'

$cacheDir = "d:\code-space\coding-soler\permission-system-0d41021b\frontend\.npm-cache"
$env:npm_config_cache = $cacheDir

$logFile = "d:\code-space\coding-soler\permission-system-0d41021b\frontend\vite-start.log"
$stdoutFile = "d:\code-space\coding-soler\permission-system-0d41021b\frontend\vite-stdout.log"
$stderrFile = "d:\code-space\coding-soler\permission-system-0d41021b\frontend\vite-stderr.log"

"Starting Vite launch at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" | Set-Content $logFile

$nodeExe = "D:\dev-tools\node-v22.22.3-win-x64\node.exe"
$npmCli = "D:\dev-tools\node-v22.22.3-win-x64\node_modules\npm\bin\npm-cli.js"

$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = $nodeExe
$psi.Arguments = "`"$npmCli`" run dev"
$psi.WorkingDirectory = "d:\code-space\coding-soler\permission-system-0d41021b\frontend"
$psi.RedirectStandardOutput = $true
$psi.RedirectStandardError = $true
$psi.UseShellExecute = $false
$psi.EnvironmentVariables["npm_config_cache"] = $cacheDir

$proc = New-Object System.Diagnostics.Process
$proc.StartInfo = $psi
$proc.Start() | Out-Null

"Started Vite process ID: $($proc.Id)" | Add-Content $logFile

$swOut = [System.IO.StreamWriter]::new($stdoutFile, $true)
$swErr = [System.IO.StreamWriter]::new($stderrFile, $true)

$outAction = {
    param($sender, $e)
    if (-not [string]::IsNullOrEmpty($e.Data)) {
        $swOut.WriteLine($e.Data)
        $swOut.Flush()
    }
}

$errAction = {
    param($sender, $e)
    if (-not [string]::IsNullOrEmpty($e.Data)) {
        $swErr.WriteLine($e.Data)
        $swErr.Flush()
    }
}

Register-ObjectEvent -InputObject $proc -EventName OutputDataReceived -Action $outAction | Out-Null
Register-ObjectEvent -InputObject $proc -EventName ErrorDataReceived -Action $errAction | Out-Null

$proc.BeginOutputReadLine()
$proc.BeginErrorReadLine()

"Process started. Waiting 15 seconds for startup..." | Add-Content $logFile
Start-Sleep -Seconds 15

if ($proc.HasExited) {
    "Process exited with code: $($proc.ExitCode)" | Add-Content $logFile
} else {
    "Vite is still running (PID: $($proc.Id))" | Add-Content $logFile
}

"Check completed at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" | Add-Content $logFile
