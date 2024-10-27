package sderr

import (
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strings"
)

var (
	StackTraceMaxDepth = 10
	packageName        = reflect.TypeOf(packageTag{}).PkgPath()
)

type packageTag struct{}

type Frame struct {
	PC   uintptr
	File string
	Func string
	Line int
}

type Stack struct {
	frames []Frame
}

func newStacktrace() *Stack {
	goRoot := runtime.GOROOT()
	var frames []Frame
	for i := 0; i < StackTraceMaxDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		isGoPkg := goRoot != "" && strings.Contains(file, goRoot)
		isSderrPkg := strings.Contains(f.Name(), packageName)
		isTestPkg := strings.HasSuffix(file, "_test.go")
		if !isGoPkg && (!isSderrPkg || isTestPkg) {
			frames = append(frames, Frame{pc, file, trimFunc(f), line})
		}
	}
	slices.Reverse(frames)
	return &Stack{frames: frames}
}

func (s *Stack) Frames() []Frame {
	if s == nil {
		return nil
	}
	return slices.Clone(s.frames)
}

func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.File, f.Line, f.Func)
}

func trimFunc(f *runtime.Func) string {
	longName := f.Name()
	shortName := longName[strings.LastIndex(longName, "/")+1:]
	return shortName
}
