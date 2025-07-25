package main

import (
	"os"
	"sync"
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

func TestUpdateAllLocalRep(t *testing.T) {
	stream.UpdateAllLocalRep()
}
func TestName(t *testing.T) { // 模块代理刷新的不及时，需要禁用代理
	mylog.Check(os.Setenv("GOPROXY", "direct"))
	reps := map[string]string{
		//"gioui.org":                    "main",
		//"github.com/ddkwork/gio":       "main",
		"gioui.org/cmd":                "main",
		"gioui.org/example":            "main",
		"gioui.org/x":                  "main",
		"github.com/oligo/gvcode":      "main",
		"github.com/ddkwork/golibrary": "master",
		"github.com/ddkwork/ux":        "master",
		"github.com/ddkwork/ddk":       "master",
		"github.com/ddkwork/bindgen":   "master",
		"github.com/ddkwork/crypt":     "master",
		"github.com/google/go-cmp":     "master",
		"golang.org/x/arch":            "master",
		"modernc.org/ccgo":             "master",
		"golang.org/x/tools/gopls":     "master",
		"github.com/ebitengine/purego": "main",
		"github.com/saferwall/pe":      "main",
	}
	fack := `
replace gioui.org => github.com/ddkwork/gio latest
`
	stream.NewBuffer("go.mod").Append(stream.NewBuffer(fack)).ReWriteSelf()

	w := sync.WaitGroup{}
	for k, v := range reps {
		w.Go(func() {
			stream.RunCommandSafe("go", "get", k+"@"+v)
		})
	}
	w.Wait()
	g := stream.NewGeneratedFile()
	for s := range stream.ReadFileToLines("go.mod") {
		g.P(s)
	}
	stream.WriteTruncate("go.mod", g.Bytes())
	mylog.Json("mod", string(mylog.Check2(os.ReadFile("go.mod"))))
}
