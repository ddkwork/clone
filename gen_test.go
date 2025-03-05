package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"os"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	UpdateDependencies()
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
		println(stream.RunCommand(s).Output.String())
	}
	mylog.Success("mod", mylog.Check2(os.ReadFile("go.mod")))
}
