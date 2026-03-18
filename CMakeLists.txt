cmake_minimum_required(VERSION 3.28)
project(HyperDbgUnified C ASM_MASM)
set_property(GLOBAL PROPERTY USE_FOLDERS ON)

set(CMAKE_MODULE_PATH "${CMAKE_CURRENT_LIST_DIR}")
find_package(WDK REQUIRED)

set(WDK_NTDDI_VERSION "0x0A000008")

include_directories(
    "."
    "HyperDbg/hyperdbg/include"
    "HyperDbg/hyperdbg/dependencies/zydis/include"
    "HyperDbg/hyperdbg/dependencies/zydis/src"
    "HyperDbg/hyperdbg/dependencies/zydis/dependencies/zycore/include"
    "HyperDbg/hyperdbg/dependencies"
)

add_definitions(
    -DZYAN_NO_LIBC
    -DZYDIS_NO_LIBC
    -DZYCORE_STATIC_BUILD
    -DZYDIS_STATIC_BUILD
    -DNTDDI_VERSION=0x0A000008
    -DHYPERDBG_STATIC_BUILD
)

file(GLOB ZydisSourceFiles "HyperDbg/hyperdbg/dependencies/zydis/src/*.c")
file(GLOB ZycoreSourceFiles "HyperDbg/hyperdbg/dependencies/zydis/dependencies/zycore/*.c")
set(ZydisSourceFiles ${ZydisSourceFiles} ${ZycoreSourceFiles})

file(GLOB_RECURSE HyperkdSourceFiles "HyperDbg/hyperdbg/hyperkd/code/*.c")
list(FILTER HyperkdSourceFiles EXCLUDE REGEX ".*code/common/Common\\.c$")
list(FILTER HyperkdSourceFiles EXCLUDE REGEX ".*code/debugger/broadcast/DpcRoutines\\.c$")
file(GLOB HyperkdAsmFiles "HyperDbg/hyperdbg/hyperkd/code/assembly/*.asm")

file(GLOB_RECURSE HyperhvSourceFiles "HyperDbg/hyperdbg/hyperhv/code/*.c")
list(FILTER HyperhvSourceFiles EXCLUDE REGEX ".*code/common/UnloadDll\\.c$")
file(GLOB HyperhvAsmFiles "HyperDbg/hyperdbg/hyperhv/code/assembly/*.asm")

file(GLOB HyperlogSourceFiles "HyperDbg/hyperdbg/hyperlog/code/*.c")

file(GLOB KdserialSourceFiles "HyperDbg/hyperdbg/kdserial/*.c")

file(GLOB ScriptEvalSourceFiles "HyperDbg/hyperdbg/script-eval/code/*.c")

file(GLOB HyperevadeSourceFiles "HyperDbg/hyperdbg/hyperevade/code/*.c")

file(GLOB HypertraceSourceFiles "HyperDbg/hyperdbg/hypertrace/code/*.c")

file(GLOB IncludeOptimizationsSourceFiles "HyperDbg/hyperdbg/include/components/optimizations/code/*.c")
file(GLOB IncludeSpinlockSourceFiles "HyperDbg/hyperdbg/include/components/spinlock/code/*.c")
file(GLOB IncludePlatformSourceFiles "HyperDbg/hyperdbg/include/platform/kernel/code/*.c")

wdk_add_library(ZydisLib STATIC ${ZydisSourceFiles})

wdk_add_library(HyperhvLib STATIC ${HyperhvSourceFiles} ${HyperhvAsmFiles} ${IncludeOptimizationsSourceFiles} ${IncludeSpinlockSourceFiles})
target_include_directories(HyperhvLib PRIVATE
    "HyperDbg/hyperdbg/hyperhv"
    "HyperDbg/hyperdbg/hyperhv/header"
    "HyperDbg/hyperdbg/hyperhv/code"
    "HyperDbg/hyperdbg/include/components/optimizations/code"
    "HyperDbg/hyperdbg/include/components/spinlock/code"
)
target_compile_definitions(HyperhvLib PRIVATE DHYPERDBG_VMM)
target_precompile_headers(HyperhvLib PRIVATE "HyperDbg/hyperdbg/hyperhv/pch.h")

wdk_add_library(HyperlogLib STATIC ${HyperlogSourceFiles})
target_include_directories(HyperlogLib PRIVATE
    "HyperDbg/hyperdbg/hyperlog"
    "HyperDbg/hyperdbg/hyperlog/header"
    "HyperDbg/hyperdbg/hyperlog/code"
)
target_compile_definitions(HyperlogLib PRIVATE DHYPERDBG_HYPERLOG IMPORT_EXPORT_HYPERLOG=)
target_precompile_headers(HyperlogLib PRIVATE "HyperDbg/hyperdbg/hyperlog/header/pch.h")

wdk_add_library(KdserialLib STATIC ${KdserialSourceFiles})
target_include_directories(KdserialLib PRIVATE
    "HyperDbg/hyperdbg/kdserial"
)

wdk_add_library(HyperevadeLib STATIC ${HyperevadeSourceFiles})
target_include_directories(HyperevadeLib PRIVATE
    "HyperDbg/hyperdbg/hyperevade"
    "HyperDbg/hyperdbg/hyperevade/header"
    "HyperDbg/hyperdbg/hyperevade/code"
)
target_compile_definitions(HyperevadeLib PRIVATE DHYPERDBG_HYPEREVADE IMPORT_EXPORT_HYPEREVADE=)
target_precompile_headers(HyperevadeLib PRIVATE "HyperDbg/hyperdbg/hyperevade/header/pch.h")

wdk_add_library(HypertraceLib STATIC ${HypertraceSourceFiles})
target_include_directories(HypertraceLib PRIVATE
    "HyperDbg/hyperdbg/hypertrace"
    "HyperDbg/hyperdbg/hypertrace/header"
    "HyperDbg/hyperdbg/hypertrace/code"
)
target_compile_definitions(HypertraceLib PRIVATE DHYPERDBG_HYPERTRACE IMPORT_EXPORT_HYPERTRACE=)
target_precompile_headers(HypertraceLib PRIVATE "HyperDbg/hyperdbg/hypertrace/header/pch.h")

wdk_add_library(PlatformLib STATIC ${IncludePlatformSourceFiles})
target_include_directories(PlatformLib PRIVATE
    "HyperDbg/hyperdbg/include/platform/kernel/code"
)
target_precompile_headers(PlatformLib PRIVATE "HyperDbg/hyperdbg/include/platform/kernel/code/pch.h")

wdk_add_driver(hyperkd
    KMDF 1.15
    ${HyperkdSourceFiles}
    ${HyperkdAsmFiles}
    ${IncludeOptimizationsSourceFiles}
    ${IncludeSpinlockSourceFiles}
    ${ScriptEvalSourceFiles}
)
target_include_directories(hyperkd PRIVATE
    "HyperDbg/hyperdbg/hyperkd"
    "HyperDbg/hyperdbg/hyperkd/header"
    "HyperDbg/hyperdbg/hyperkd/code"
    "HyperDbg/hyperdbg/include/components/optimizations/code"
    "HyperDbg/hyperdbg/include/components/spinlock/code"
    "HyperDbg/hyperdbg/script-eval"
)
target_compile_definitions(hyperkd PRIVATE DHYPERDBG_DEBUGGER IMPORT_EXPORT_HYPERLOG= IMPORT_EXPORT_HYPERTRACE= IMPORT_EXPORT_HYPEREVADE=)
target_precompile_headers(hyperkd PRIVATE "HyperDbg/hyperdbg/hyperkd/header/pch.h")

target_compile_options(hyperkd PRIVATE /Zi)
target_link_options(hyperkd PRIVATE /DEBUG)

target_link_libraries(hyperkd
    ZydisLib
    HyperhvLib
    HyperlogLib
    KdserialLib
    HyperevadeLib
    HypertraceLib
    PlatformLib
)
