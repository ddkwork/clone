package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	UpdateDependencies()
	TestParseGoMod(t)
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
	mylog.Json("mod", string(mylog.Check2(os.ReadFile("go.mod"))))
}

func TestParseGoMod(t *testing.T) {
	g := stream.NewGeneratedFile()
	m := ParseGoMod()
	for k, v := range m.Range() {
		cmd := "go get -x " + k + "@" + v
		g.P(cmd)
	}
	g.P("go mod tidy")
	stream.WriteTruncate("dep.txt", g.String())
	stream.WriteTruncate(filepath.Join(GetDesktopDir(), "dep.txt"), g.String())
	println(g.String())
}

func ParseGoMod() *safemap.M[string, string] {
	path := "go.mod"
	f := mylog.Check2(modfile.Parse(path, mylog.Check2(os.ReadFile(path)), nil))
	return safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
		for _, req := range f.Require {
			yield(req.Mod.Path, req.Mod.Version)
		}
	})
}

func GetDesktopDir() string {
	// 获取用户主目录
	homeDir := mylog.Check2(os.UserHomeDir())
	// 根据操作系统处理路径
	switch runtime.GOOS {
	case "windows", "darwin":
		// Windows和macOS直接拼接Desktop
		return filepath.Join(homeDir, "Desktop")
	case "linux":
		// Linux优先检查XDG环境变量
		if xdgDir := os.Getenv("XDG_DESKTOP_DIR"); xdgDir != "" {
			return xdgDir
		}
		// 默认使用主目录下的Desktop
		return filepath.Join(homeDir, "Desktop")
	default:
		panic("unsupported platform")
	}
}
