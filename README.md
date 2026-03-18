# HyperDbg Unified Driver

这是一个自动构建HyperDbg统一驱动的项目，用于解决蓝屏调试时PDB符号文件无法正确加载的问题。

## 背景

由于原始HyperDbg项目将驱动分为多个独立的模块（hyperkd、hyperhv、hyperlog等），在调试蓝屏问题时，PDB符号文件无法正确加载，因为模块已经被卸载。为了解决这个问题，我们创建了这个统一工程，将所有驱动模块合并为一个。

## 功能特性

- **自动克隆**：自动从GitHub克隆HyperDbg最新代码
- **自动打补丁**：应用编译修复补丁
- **自动构建**：使用CMake + WDK构建驱动
- **自动发布**：构建成功后自动发布到GitHub Releases
- **PDB符号**：生成完整的调试符号文件用于BSOD调试

## 编译方法

### 本地编译（EWDK环境）

#### 前置要求

- EWDK（Enterprise WDK）环境
- PowerShell
- CMake 3.28 或更高版本

#### 编译步骤

1. **配置EWDK路径**：
   编辑 `build.ps1`，设置正确的WDK和EWDK路径：
   ```powershell
   $WDK_PATH = "E:\Program Files\Windows Kits\10"
   $EWDK_BUILD_ENV = "E:\BuildEnv\SetupBuildEnv.cmd"
   ```

2. **运行构建脚本**：
   ```powershell
   .\build.ps1
   ```

3. **构建产物**：
   ```
   build/Release/
   ├── hyperkd.sys    # 统一驱动文件
   └── hyperkd.pdb    # PDB符号文件
   ```

### GitHub Actions自动构建

#### 触发方式

- **自动触发**：推送到 `master` 或 `dev` 分支
- **定时触发**：每6小时自动运行一次
- **手动触发**：通过GitHub Actions界面

#### 获取构建产物

访问GitHub Releases页面，下载 `latest` tag下的文件：
- `hyperkd.sys` - 驱动文件
- `hyperkd.pdb` - 调试符号文件
- `HyperDbg.tar.zst` - 源码压缩包

## 项目结构

```
HyperDbgUnified/
├── .github/
│   └── workflows/
│       └── build.yml           # GitHub Actions工作流
├── CMakeLists.txt              # CMake构建配置
├── FindWdk.cmake              # WDK查找模块
├── build.ps1                  # 本地构建脚本
├── driver_compilation_fix.patch # 编译修复补丁
├── README.md                  # 本文档
├── .gitignore                # Git忽略规则
├── HyperDbg/                 # 克隆的HyperDbg源码（不提交）
└── build/                    # 构建输出目录（不提交）
    └── Release/
        ├── hyperkd.sys
        └── hyperkd.pdb
```

## 编译修复补丁

`driver_compilation_fix.patch` 包含以下修复：

1. **移除DLL依赖**：将DLL导出改为静态链接
2. **添加缺失函数**：补全DPC例程和公共函数
3. **修复符号冲突**：解决模块间的符号冲突
4. **启用PDB生成**：确保生成调试符号文件

## 调试蓝屏问题

当遇到蓝屏问题时，可以按照以下步骤进行调试：

1. **启用小内存转储**：
   ```powershell
   bcdedit /set minidump on
   ```

2. **加载驱动**：
   ```powershell
   sc create hyperkd type= kernel binPath= "C:\path\to\hyperkd.sys"
   sc start hyperkd
   ```

3. **触发蓝屏**（如果需要）：
   ```powershell
   # 使用HyperDbg命令触发
   ```

4. **分析转储文件**：
   ```powershell
   # 使用WinDbg分析转储文件
   windbg -z C:\Windows\Minidump\*.dmp
   ```

5. **加载PDB符号**：
   ```
   .sympath C:\path\to\pdb
   .symopt+ 0x402
   .symopt+ 0x800
   .reload /f
   !analyze -v
   ```

## 注意事项

1. **测试签名**：在测试环境中，需要启用测试签名：
   ```powershell
   bcdedit /set testsigning on
   ```

2. **EWDK环境**：本地编译需要配置EWDK环境变量

3. **网络连接**：GitHub Actions需要访问GitHub和WDK下载地址

4. **补丁应用**：每次构建都会自动应用补丁，确保代码兼容性

5. **符号文件**：PDB文件与sys文件必须匹配，否则无法正确调试

## 许可证

本项目遵循HyperDbg的GPL v3许可证。

## 相关链接

- HyperDbg项目: https://github.com/HyperDbg/HyperDbg
- GitHub Releases: https://github.com/ddkwork/clone/releases
- GitHub Actions: https://github.com/ddkwork/clone/actions

## 致谢

感谢HyperDbg项目的所有贡献者。
