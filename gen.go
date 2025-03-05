package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
	"os"
	"strings"
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

func UpdateDependencies() { //模块代理刷新的不及时，需要禁用代理
	mylog.Check(os.Setenv("GOPROXY", "direct"))
	for s := range strings.Lines(`
     go get -x gioui.org@main
	 go get -x gioui.org/cmd@main
	 go get -x gioui.org/example@main
	 go get -x gioui.org/x@main
	 go get -x github.com/oligo/gvcode@main
	 go get -x github.com/ddkwork/golibrary@master
	 go get -x github.com/ddkwork/ux@master
	 go get -x github.com/google/go-cmp@master
	 go get -x github.com/ddkwork/app@master
	 go get -x github.com/ddkwork/toolbox@master
	 go get -x github.com/ddkwork/unison@master
	 go get -x github.com/ebitengine/purego@main
	 go get -x github.com/saferwall/pe@main
	 ::go get -u -x all
	 //go mod tidy

	//go install mvdan.cc/gofumpt@latest
	//gofumpt -l -w .
	//go install honnef.co/go/tools/cmd/staticcheck@latest
	//staticcheck ./...
	//go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

`) {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "::") || strings.HasPrefix(s, "//") || s == "" {
			continue
		}
		stream.RunCommand(s)
	}
	for s := range stream.ReadFileToLines("go.mod") {
		println(s)
	}
}
