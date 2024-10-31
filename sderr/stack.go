package sderr

import (
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strings"
)

var (
	StackTraceMaxDepth = 10 // 默认的stack最大深度
	packageName        = reflect.TypeOf(packageTag{}).PkgPath()
)

type packageTag struct{}

// Frame 栈上的一帧
type Frame struct {
	PC   uintptr
	File string // 文件名
	Func string // 函数名
	Line int    // 行号
}

// Stack 一个stacktrace
type Stack struct {
	frames []Frame // stack上的所有帧
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
	return &Stack{frames: frames}
}

// Frames 灰灰一个stacktrace上的所有帧
func (s *Stack) Frames() []Frame {
	if s == nil {
		return nil
	}
	return slices.Clone(s.frames)
}

// String 返回一帧的字符串描述
func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.File, f.Line, f.Func)
}

func trimFunc(f *runtime.Func) string {
	longName := f.Name()
	shortName := longName[strings.LastIndex(longName, "/")+1:]
	return shortName
}
