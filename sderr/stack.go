package sderr

import (
	"runtime"
	"slices"
)

var (
	StackTraceMaxDepth = 10
)

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
		file, fn := trimFilename(file), trimFunc(f)
		frames = append(frames, Frame{pc, file, fn, line})
	}
	return &Stack{frames: frames}
}

func (s *Stack) Frames() []Frame {
	return slices.Clone(s.frames)
}

func trimFilename(file string) string {
	// TODO
	return file
}

func trimFunc(f *runtime.Func) string {
	// TODO
	return f.Name()
}
