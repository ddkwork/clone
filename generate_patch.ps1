param(
    [Parameter(Mandatory=$false)]
    [string]$Commit = "HEAD~1",
    
    [Parameter(Mandatory=$false)]
    [string]$OutputPath = "driver_compilation_fix.patch"
)

$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path

$ErrorActionPreference = "Stop"

Write-Host "Generating patch from $Commit to HEAD..." -ForegroundColor Green

$repoRoot = git rev-parse --show-toplevel
if (-not $repoRoot) {
    Write-Error "Not in a git repository"
    exit 1
}

Write-Host "Repository root: $repoRoot" -ForegroundColor Cyan

$patchPath = Join-Path $SCRIPT_DIR $OutputPath
Write-Host "Output path: $patchPath" -ForegroundColor Cyan

$excludeFiles = @(
    "hwdbg/sim/hwdbg/DebuggerModuleTestingBRAM/test.sh"
)

$excludeArgs = $excludeFiles | ForEach-Object { "':!$_'" }
$excludeArgsString = $excludeArgs -join " "

$command = "git diff $Commit HEAD -- . $excludeArgsString"
Write-Host "Executing: $command" -ForegroundColor Yellow

git diff $Commit HEAD -- . $excludeArgsString | Out-File -FilePath $patchPath -Encoding utf8 -NoNewline

if (Test-Path $patchPath) {
    $size = (Get-Item $patchPath).Length
    Write-Host "Patch generated successfully: $patchPath ($size bytes)" -ForegroundColor Green
} else {
    Write-Error "Failed to generate patch file"
    exit 1
}