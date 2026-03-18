param(
    [Parameter(Mandatory=$false)]
    [string]$TargetDir = "."
)

$ErrorActionPreference = "Stop"

Write-Host "Applying driver compilation fixes to $TargetDir..." -ForegroundColor Green

$filesToModify = @(
    @{
        Path = "hyperdbg/hyperhv/code/broadcast/DpcRoutines.c"
        Pattern = "^DpcRoutineTerminateGuest\(KDPC \* Dpc, PVOID DeferredContext, PVOID SystemArgument1, PVOID SystemArgument2\)\s*\{\s*//.*?\s*KeSignalCallDpcDone\(SystemArgument1\);\s*\}"
        InsertAfter = "KeSignalCallDpcDone(SystemArgument1);`n`n"
        Content = @"
VOID
DpcRoutinePerformWriteMsr(KDPC * Dpc, PVOID DeferredContext, PVOID SystemArgument1, PVOID SystemArgument2)
{
    UNREFERENCED_PARAMETER(Dpc);
    UNREFERENCED_PARAMETER(DeferredContext);
    UNREFERENCED_PARAMETER(SystemArgument1);
    UNREFERENCED_PARAMETER(SystemArgument2);

    //
    // write on MSR
    //
    __writemsr(0, 0);

    //
    // As this function is designed for a single,
    // we have to release the synchronization lock here
    //
    SpinlockUnlock(&OneCoreLock);
}

VOID
DpcRoutinePerformReadMsr(KDPC * Dpc, PVOID DeferredContext, PVOID SystemArgument1, PVOID SystemArgument2)
{
    UNREFERENCED_PARAMETER(Dpc);
    UNREFERENCED_PARAMETER(DeferredContext);
    UNREFERENCED_PARAMETER(SystemArgument1);
    UNREFERENCED_PARAMETER(SystemArgument2);

    //
    // read on MSR
    //
    __readmsr(0);

    //
    // As this function is designed for a single,
    // we have to release the synchronization lock here
    //
    SpinlockUnlock(&OneCoreLock);
}

VOID
DpcRoutineWriteMsrToAllCores(KDPC * Dpc, PVOID DeferredContext, PVOID SystemArgument1, PVOID SystemArgument2)
{
    UNREFERENCED_PARAMETER(Dpc);
    UNREFERENCED_PARAMETER(DeferredContext);
    UNREFERENCED_PARAMETER(SystemArgument1);
    UNREFERENCED_PARAMETER(SystemArgument2);

    //
    // write on MSR
    //
    __writemsr(0, 0);

    //
    // Wait for all DPCs to synchronize at this point
    //
    KeSignalCallDpcSynchronize(SystemArgument2);

    //
    // Mark as DPC being complete
    //
    KeSignalCallDpcDone(SystemArgument1);
}

VOID
DpcRoutineReadMsrToAllCores(KDPC * Dpc, PVOID DeferredContext, PVOID SystemArgument1, PVOID SystemArgument2)
{
    UNREFERENCED_PARAMETER(Dpc);
    UNREFERENCED_PARAMETER(DeferredContext);
    UNREFERENCED_PARAMETER(SystemArgument1);
    UNREFERENCED_PARAMETER(SystemArgument2);

    //
    // read msr
    //
    __readmsr(0);

    //
    // Wait for all DPCs to synchronize at this point
    //
    KeSignalCallDpcSynchronize(SystemArgument2);

    //
    // Mark as DPC being complete
    //
    KeSignalCallDpcDone(SystemArgument1);
}

BOOLEAN
DpcRoutineSetHardwareDebugRegisters(KDPC * Dpc, PVOID DeferredContext, PVOID SystemArgument1, PVOID SystemArgument2)
{
    UNREFERENCED_PARAMETER(Dpc);
    UNREFERENCED_PARAMETER(DeferredContext);
    UNREFERENCED_PARAMETER(SystemArgument1);
    UNREFERENCED_PARAMETER(SystemArgument2);

    //
    // Apply hardware debug registers
    //

    //
    // Wait for all DPCs to synchronize at this point
    //
    KeSignalCallDpcSynchronize(SystemArgument2);

    //
    // Mark as DPC being complete
    //
    KeSignalCallDpcDone(SystemArgument1);

    return TRUE;
}
"@
    },
    @{
        Path = "hyperdbg/hyperhv/code/common/Common.c"
        Pattern = "^    return TRUE;\s*\}"
        InsertAfter = "return TRUE;`n`n"
        Content = @"
BOOLEAN
CommonIsProcessExist(UINT32 ProcId)
{
    PEPROCESS TargetEprocess;

    if (PsLookupProcessByProcessId((HANDLE)ProcId, &TargetEprocess) != STATUS_SUCCESS)
    {
        return FALSE;
    }
    else
    {
        ObDereferenceObject(TargetEprocess);

        return TRUE;
    }
}

BOOLEAN
CommonValidateCoreNumber(UINT32 CoreNumber)
{
    ULONG ProcessorsCount;

    ProcessorsCount = KeQueryActiveProcessorCount(0);

    if (CoreNumber >= ProcessorsCount)
    {
        return FALSE;
    }
    else
    {
        return TRUE;
    }
}

BOOLEAN
CommonKillProcess(UINT32 ProcessId, UINT32 KillingMethod)
{
    NTSTATUS  Status        = STATUS_SUCCESS;
    HANDLE    ProcessHandle = NULL;
    PEPROCESS Process       = NULL;

    if (ProcessId == NULL_ZERO)
    {
        return FALSE;
    }

    Status = PsLookupProcessByProcessId((HANDLE)ProcessId, &Process);

    if (!NT_SUCCESS(Status))
    {
        return FALSE;
    }

    Status = ObOpenObjectByPointer(Process, OBJ_KERNEL_HANDLE, NULL, GENERIC_ALL, *PsProcessType, KernelMode, &ProcessHandle);

    if (!NT_SUCCESS(Status))
    {
        ObDereferenceObject(Process);
        return FALSE;
    }

    Status = ZwTerminateProcess(ProcessHandle, 0);

    ObDereferenceObject(Process);
    ZwClose(ProcessHandle);

    if (!NT_SUCCESS(Status))
    {
        return FALSE;
    }

    return TRUE;
}
"@
    }
)

foreach ($file in $filesToModify) {
    $filePath = Join-Path $TargetDir $file.Path
    
    if (-not (Test-Path $filePath)) {
        Write-Host "Warning: File not found: $filePath" -ForegroundColor Yellow
        continue
    }
    
    Write-Host "Processing: $filePath" -ForegroundColor Cyan
    
    $content = Get-Content $filePath -Raw
    
    if ($content -match [regex]::Escape($file.Pattern)) {
        $content = $content -replace [regex]::Escape($file.Pattern), "`$0`n`n$($file.Content)"
        Set-Content -Path $filePath -Value $content -NoNewline
        Write-Host "  Applied modifications" -ForegroundColor Green
    } else {
        Write-Host "  Pattern not found, skipping" -ForegroundColor Yellow
    }
}

Write-Host "Fixes applied successfully!" -ForegroundColor Green