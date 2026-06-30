$ErrorActionPreference = 'Continue'
Set-Location 'd:\code-space\coding-soler\permission-system-0d41021b\backend'

$stdoutLog = 'd:\code-space\coding-soler\permission-system-0d41021b\backend\server-stdout.log'
$stderrLog = 'd:\code-space\coding-soler\permission-system-0d41021b\backend\server-stderr.log'

"[$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')] Starting backend server..." | Out-File $stdoutLog -Encoding utf8

$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = 'd:\code-space\coding-soler\permission-system-0d41021b\backend\server.exe'
$psi.WorkingDirectory = 'd:\code-space\coding-soler\permission-system-0d41021b\backend'
$psi.RedirectStandardOutput = $true
$psi.RedirectStandardError = $true
$psi.UseShellExecute = $false

$proc = New-Object System.Diagnostics.Process
$proc.StartInfo = $psi
$proc.Start() | Out-Null

"[$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')] Backend server started. PID: $($proc.Id)" | Add-Content $stdoutLog

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

Write-Host "Backend server PID: $($proc.Id)"
$proc.Id | Out-File 'd:\code-space\coding-soler\permission-system-0d41021b\backend\server-pid.txt'

Start-Sleep -Seconds 3
if ($proc.HasExited) {
    Write-Host "ERROR: Backend server exited with code: $($proc.ExitCode)"
    exit 1
} else {
    Write-Host "Backend server is running on port 8080 (PID: $($proc.Id))"
    exit 0
}
