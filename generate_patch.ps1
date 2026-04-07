param(
    [Parameter(Mandatory=$false)]
    [string]$OutputPath = "driver_compilation_fix.patch"
)

$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$HYPERDBG_DIR = Join-Path $SCRIPT_DIR "HyperDbg"

$ErrorActionPreference = "Stop"

Write-Host "=== HyperDbg Patch Generator ===" -ForegroundColor Cyan

if (-not (Test-Path $HYPERDBG_DIR)) {
    Write-Error "HyperDbg directory not found: $HYPERDBG_DIR"
    exit 1
}

Push-Location $HYPERDBG_DIR
try {
    $repoRoot = git rev-parse --show-toplevel 2>$null
    if (-not $repoRoot) {
        Write-Error "Not a git repository: $HYPERDBG_DIR"
        exit 1
    }

    Write-Host "Repository root: $repoRoot" -ForegroundColor Green

    $patchPath = Join-Path $SCRIPT_DIR $OutputPath
    Write-Host "Output path: $patchPath" -ForegroundColor Green

    Write-Host "Generating patch from working directory changes..." -ForegroundColor Yellow

    $excludeArgs = @(
        ":(exclude)hwdbg/sim/hwdbg/DebuggerModuleTestingBRAM/test.sh"
        ":(exclude)hyperdbg/dependencies/zydis"
    )

    git diff --output="$patchPath" -- . @excludeArgs

    if (Test-Path $patchPath) {
        $size = (Get-Item $patchPath).Length
        if ($size -gt 0) {
            Write-Host "Patch generated successfully: $patchPath ($size bytes)" -ForegroundColor Green
        } else {
            Remove-Item $patchPath
            Write-Host "No changes to export" -ForegroundColor Yellow
        }
    } else {
        Write-Host "No changes to export" -ForegroundColor Yellow
    }
} finally {
    Pop-Location
}
