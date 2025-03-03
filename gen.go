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
	g.P(`
name: clone仓库

on:
  push:
    branches:
      - master  # 或者你想要触发工作流的分支名
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: 检出代码
        uses: actions/checkout@v4
`)
	g.P("      - name: 克隆", name, "仓库")
	g.P("        run: git clone --recursive ", url)
	g.P("      - name: 打包", name, "项目")
	g.P("        run: tar -czvf "+name, ".tar.gz ", name)
	g.P()
	g.P("      - name: 上传打包文件")
	g.P("        uses: actions/upload-artifact@v4")
	g.P("        with:")
	g.P("          name: ", name, "-tar-gz")
	g.P("          path: ", name, ".tar.gz")
	stream.WriteTruncate(".github/workflows/clone.yml", g.String())
}
