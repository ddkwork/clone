package main

import (
	"os"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

func TestName(t *testing.T) { // 模块代理刷新的不及时，需要禁用代理
	mylog.Check(os.Setenv("GOPROXY", "direct"))
	for s := range strings.Lines(`
     go get -x gioui.org@main
	 go get -x gioui.org/cmd@main
	 go get -x gioui.org/example@main
	 go get -x gioui.org/x@main
	 go get -x github.com/oligo/gvcode@main
	 go get -x github.com/ddkwork/golibrary@master
	 go get -x github.com/ddkwork/ux@master
	 go get -x github.com/ddkwork/ddk@master
	 go get -x github.com/ddkwork/bindgen@master
	 go get -x github.com/ddkwork/crypt@master
	 go get -x github.com/google/go-cmp@master
	 go get -x golang.org/x/arch@master
	 go get -x modernc.org/ccgo@master
	 go get -x golang.org/x/tools/gopls@master
	 go get -x github.com/ddkwork/app@master
	 go get -x github.com/ddkwork/toolbox@master
	 go get -x github.com/ddkwork/unison@master
	 go get -x github.com/ebitengine/purego@main
	 go get -x github.com/saferwall/pe@main
`) {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "::") || strings.HasPrefix(s, "//") || s == "" {
			continue
		}
		stream.RunCommand(s)
	}
	g := stream.NewGeneratedFile()
	for s := range stream.ReadFileToLines("go.mod") {
		if strings.Contains(s, "gvcode") {
			s = `github.com/oligo/gvcode v0.2.1-0.20250424030509-8138ffc92f73`
		}
		g.P(s)
	}
	stream.WriteTruncate("go.mod", g.Bytes())
	mylog.Json("mod", string(mylog.Check2(os.ReadFile("go.mod"))))
	mylog.Json("log", string(mylog.Check2(os.ReadFile("log.log"))))
}
