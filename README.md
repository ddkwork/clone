# HyperDbg 统一工程

这是一个合并了所有HyperDbg模块的统一工程，便于编译和维护。

## 背景

由于原始HyperDbg项目将驱动分为多个独立的模块（hyperkd、hyperhv、hyperlog等），在调试蓝屏问题时，PDB符号文件无法正确加载，因为模块已经被卸载。为了解决这个问题，我们创建了这个统一工程，将所有驱动模块合并为一个。

## 模块列表

此工程包含以下模块：

- **hyperkd**: 内核调试器驱动
- **hyperhv**: 虚拟机监控器
- **hyperlog**: 日志系统
- **kdserial**: 串口通信
- **symbol-parser**: 符号解析器
- **script-engine**: 脚本引擎
- **libhyperdbg**: 调试器库
- **hyperdbg-cli**: 命令行界面
- **dependencies/zydis**: 反汇编器
- **dependencies/pdbex**: PDB符号提取器

## 编译方法

### 前置要求

- Visual Studio 2019 或更高版本
- Windows Driver Kit (WDK) 10 或更高版本
- CMake 3.28 或更高版本

### 编译步骤

1. **生成项目**：
   ```bash
   mkdir build
   cd build
   cmake ../HyperDbgUnified -G "Visual Studio 17 2022" -A x64 -DCMAKE_SYSTEM_NAME=Windows -DCMAKE_SYSTEM_VERSION=10.0.22621.0
   ```

2. **编译项目**：
   ```bash
   cmake --build . --config Release
   ```

或者直接运行构建脚本：
```bash
build.bat
```

## 目录结构

```
HyperDbgUnified/
├── CMakeLists.txt          # 主CMake文件
├── build.bat              # 构建脚本
├── README.md              # 本文档
├── .gitignore            # Git忽略文件
├── HyperDbg/             # HyperDbg源码
│   ├── include/           # 公共头文件
│   ├── hyperdbg/
│   │   ├── hyperkd/      # 内核调试器驱动
│   │   ├── hyperhv/      # 虚拟机监控器
│   │   ├── hyperlog/     # 日志系统
│   │   ├── kdserial/     # 串口通信
│   │   ├── symbol-parser/# 符号解析器
│   │   ├── script-engine/ # 脚本引擎
│   │   ├── libhyperdbg/  # 调试器库
│   │   ├── hyperdbg-cli/ # 命令行界面
│   │   └── dependencies/ # 外部依赖
│   │       ├── zydis/     # 反汇编器
│   │       ├── pdbex/     # PDB符号提取器
│   │       ├── keystone/   # 汇编器
│   │       └── ia32-doc/   # IA32文档
│   └── FindWdk.cmake     # WDK查找脚本
└── build/                # 构建输出目录
    └── Release/
        ├── hyperkd_unified.sys  # 统一驱动
        └── hyperkd_unified.pdb  # PDB符号文件
```

## 调试蓝屏问题

当遇到蓝屏问题时，可以按照以下步骤进行调试：

1. **启用小内存转储**：
   ```powershell
   # 运行 enable_crash_dump.ps1 脚本
   .\enable_crash_dump.ps1
   ```

2. **触发蓝屏**：
   ```bash
   # 运行调试器两次
   .\hyperdbg.exe -driver-only
   # 按Enter退出
   .\hyperdbg.exe -driver-only
   ```

3. **分析转储文件**：
   ```bash
   # 使用WinDbg分析转储文件
   windbg -y C:\Users\Admin\Downloads\release -z C:\Windows\Minidump\*.dmp
   ```

4. **加载PDB符号**：
   ```
   .sympath C:\Users\Admin\Downloads\release
   .symopt+ 0x402
   .symopt+ 0x800
   .reload /f
   !unloaded
   ```

## 注意事项

1. **首次编译**：首次编译需要下载并安装所有依赖项，可能需要较长时间。

2. **编译环境**：确保使用正确的Visual Studio版本和WDK版本。

3. **符号文件**：编译完成后，将生成的PDB符号文件复制到指定目录以便调试。

4. **测试签名**：在测试环境中，需要启用测试签名：
   ```bash
   bcdedit /set testsigning on
   ```

5. **EWDK路径**：构建脚本默认使用E:\作为EWDK路径，如需修改请编辑build.bat。

## 许可证

本项目遵循HyperDbg的GPL v3许可证。

## 联系方式

- HyperDbg项目: https://github.com/HyperDbg/HyperDbg
- 问题反馈: https://github.com/HyperDbg/HyperDbg/issues

## 致谢

感谢HyperDbg项目的所有贡献者。
