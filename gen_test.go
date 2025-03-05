package main

import (
	"bufio"
	"iter"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	UpdateDependencies()
}

func UpdateDependencies() { //模块代理刷新的不及时，需要禁用代理
	Check(os.Setenv("GOPROXY", "direct"))
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
		RunCommand(s)
	}
	for s := range ReadFileToLines("go.mod") {
		println(s)
	}
}

func RunCommand(command string) string { // std error not support
	fnInitCmd := func() *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", command)
		}
		return exec.Command("bash", "-c", command)
	}
	return string(Check2(fnInitCmd().CombinedOutput()))
}

func ReadFileToLines(path string) iter.Seq[string] {
	return func(yield func(string) bool) {
		f := Check2(os.Open(path))
		defer func() { Check(f.Close()) }()
		scanner := bufio.NewScanner(f)
		// scanner.Split(bufio.ScanLines)
		// scanner.Buffer(nil, 1024*1024)
		lineNumber := 1
		for scanner.Scan() {
			yield(scanner.Text())
			lineNumber++
		}
		Check(scanner.Err())
	}
}

func Check[T any](result T) {
	switch err := any(result).(type) {
	case error:
		if err != nil {
			panic(err)
		}
	default:
	}
}

func Check2[T any](ret T, err error) T {
	if err != nil {
		panic(err)
	}
	return ret
}
