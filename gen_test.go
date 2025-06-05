package main

import (
	"os"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) { // 模块代理刷新的不及时，需要禁用代理
	mylog.Check(os.Setenv("GOPROXY", "direct"))
	for s := range strings.Lines(`
     go get -x gioui.org@main
	 go get -x gioui.org/cmd@main
	 go get -x gioui.org/example@main
	 go get -x gioui.org/x@main
	 //go get -x github.com/oligo/gvcode@v0.2.1-0.20250424030509-8138ffc92f73
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

	s := `go get -x github.com/oligo/gvcode@v0.2.1-0.20250424030509-8138ffc92f73`
	b := stream.NewBuffer("go.mod")
	b.NewLine().WriteString(s)
	b.ReWriteSelf()
	mylog.Json("mod", string(mylog.Check2(os.ReadFile("go.mod"))))
}
