$ErrorActionPreference = 'Continue'
Set-Location 'd:\code-space\coding-soler\permission-system-0d41021b\frontend'

$cacheDir = 'd:\code-space\coding-soler\permission-system-0d41021b\frontend\.npm-cache'
$env:npm_config_cache = $cacheDir

$stdoutLog = 'd:\code-space\coding-soler\permission-system-0d41021b\frontend\vite-stdout2.log'
$stderrLog = 'd:\code-space\coding-soler\permission-system-0d41021b\frontend\vite-stderr2.log'

"[$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')] Starting Vite dev server..." | Out-File $stdoutLog -Encoding utf8

$nodeExe = 'D:\dev-tools\node-v22.22.3-win-x64\node.exe'
$npmCli = 'D:\dev-tools\node-v22.22.3-win-x64\node_modules\npm\bin\npm-cli.js'

$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = $nodeExe
$psi.Arguments = "`"$npmCli`" run dev"
$psi.WorkingDirectory = 'd:\code-space\coding-soler\permission-system-0d41021b\frontend'
$psi.RedirectStandardOutput = $true
$psi.RedirectStandardError = $true
$psi.UseShellExecute = $false
$psi.EnvironmentVariables['npm_config_cache'] = $cacheDir

$proc = New-Object System.Diagnostics.Process
$proc.StartInfo = $psi
$proc.Start() | Out-Null

"[$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')] Vite dev server started. PID: $($proc.Id)" | Add-Content $stdoutLog

$swOut = [System.IO.StreamWriter]::new($stdoutLog, $true)
$swErr = [System.IO.StreamWriter]::new($stderrLog, $true)

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

Write-Host "Vite dev server PID: $($proc.Id)"
$proc.Id | Out-File 'd:\code-space\coding-soler\permission-system-0d41021b\frontend\vite-pid.txt'

Start-Sleep -Seconds 12
if ($proc.HasExited) {
    Write-Host "ERROR: Vite dev server exited with code: $($proc.ExitCode)"
    exit 1
} else {
    Write-Host "Vite dev server is running on port 5173 (PID: $($proc.Id))"
    exit 0
}
