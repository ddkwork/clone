param(
    [Parameter(Mandatory=$false)]
    [string]$Commit = "HEAD~1",
    
    [Parameter(Mandatory=$false)]
    [string]$OutputPath = "../driver_compilation_fix.patch"
)

$ErrorActionPreference = "Stop"

Write-Host "Generating patch from $Commit to HEAD..." -ForegroundColor Green

$repoRoot = git rev-parse --show-toplevel
if (-not $repoRoot) {
    Write-Error "Not in a git repository"
    exit 1
}

Write-Host "Repository root: $repoRoot" -ForegroundColor Cyan

$patchPath = Join-Path $repoRoot $OutputPath
Write-Host "Output path: $patchPath" -ForegroundColor Cyan

$excludeFiles = @(
    "hwdbg/sim/hwdbg/DebuggerModuleTestingBRAM/test.sh"
)

$excludeArgs = $excludeFiles | ForEach-Object { "':!$_'" }
$excludeArgsString = $excludeArgs -join " "

$command = "git diff $Commit HEAD -- . $excludeArgsString > `"$patchPath`""
Write-Host "Executing: $command" -ForegroundColor Yellow

Invoke-Expression $command

if (Test-Path $patchPath) {
    $size = (Get-Item $patchPath).Length
    Write-Host "Patch generated successfully: $patchPath ($size bytes)" -ForegroundColor Green
} else {
    Write-Error "Failed to generate patch file"
    exit 1
}