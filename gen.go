package main

import (
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
)

func main() {
	m := safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
		yield("miqt", "https://github.com/ddkwork/miqt.git")
		yield("HyperDbg", "https://github.com/HyperDbg/HyperDbg.git")
		yield("sointu", "https://github.com/vsariola/sointu.git")
		yield("unison", "https://github.com/richardwilkes/unison.git")
	})
	name := "unison"
	url := m.GetMust(name)
	g := stream.NewGeneratedFile()
	g.PKeepSpace(`
name: clone仓库

on:
  push:
    branches:
      - master  # 或者你想要触发工作流的分支名
  workflow_dispatch:

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ ubuntu-latest ]

    steps:
      - name: 检出代码
        uses: actions/checkout@v4
`)

	g.P("      - name: 克隆", name, "仓库")
	g.P("        run: git clone --recursive ", url)
	g.P("      - name: 打包", name, "项目")
	g.P("        run: tar -czvf "+name, ".tar.gz ", name)
	g.P()

	g.PKeepSpace(`      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24.0'

      # manjaro linux
      # gioui提示vulkan构建约束被排除的原因是没有安装gcc，termux可以试试是不是这个原因
      # sudo pacman -S pkg-config vulkan-headers gcc clang cmake
      # yay -S android-sdk-build-tools  android-sdk-platform-tools
      # yay -S android-platform
      # go env -w GOPROXY=https://goproxy.cn,direct
      # go run 需要sudo

      #- name: Install C library dependencies on Ubuntu
      #  if: matrix.os == 'ubuntu-latest'
      #  run: |
      #    sudo apt-get update
      #    sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

      - name: Run tests
        run: go test -v .

`)

	g.P("      - name: 上传打包文件")
	g.P("        uses: actions/upload-artifact@v4")
	g.P("        with:")
	g.P("          name: ", name, "-tar-gz")
	g.P("          path: ", name, ".tar.gz")
	g.P()

	g.P(`
      # 新增步骤：上传依赖文件
      - name: 上传依赖文件
        uses: actions/upload-artifact@v4
        with:
          name: dep-txt
          path: dep.txt

      - name: 上传mod文件
        uses: actions/upload-artifact@v4
        with:
          name: mod
          path: go.mod
`)
	stream.WriteTruncate(".github/workflows/clone.yml", g.String())
}
