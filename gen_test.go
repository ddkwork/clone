package main

import (
	"os"
	"strings"
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
		"gioui.org":                    "main",
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
		"github.com/ddkwork/app":       "master",
		"github.com/ddkwork/toolbox":   "master",
		"github.com/ddkwork/unison":    "master",
		"github.com/ebitengine/purego": "main",
		"github.com/saferwall/pe":      "main",
	}
	var all []string
	for k, v := range reps {
		if strings.Contains(k, "gvcode") {
			v = "v0.2.1-0.20250424030509-8138ffc92f73"
		}
		all = append(all, k+"@"+v)
	}
	stream.RunCommand("go", "get", strings.Join(all, " "))
	g := stream.NewGeneratedFile()
	for s := range stream.ReadFileToLines("go.mod") {
		if strings.Contains(s, "gvcode") {
			s = `github.com/oligo/gvcode v0.2.1-0.20250424030509-8138ffc92f73`
		}
		g.P(s)
	}
	stream.WriteTruncate("go.mod", g.Bytes())
	mylog.Json("mod", string(mylog.Check2(os.ReadFile("go.mod"))))
}
