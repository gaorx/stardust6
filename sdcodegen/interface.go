package sdcodegen

import (
	"path/filepath"
)

// Interface 代码生成接口
type Interface interface {
	// Sub 通过一个dirname生成新的Interface，在其下的所有代码生成都会在这个dirname下
	Sub(dirname string) Interface

	// Add 添加一个生成器
	Add(name string, handler Handler, middlewares ...Middleware) Interface

	// AddFile 加入一个文件
	AddFile(f *File) Interface

	// AddFiles 加入多个文件
	AddFiles(files []*File) Interface

	// Also 在当前生成器上执行一个函数，并返回当前生成器
	Also(f func(i Interface)) Interface

	// AddModule 添加一个模块
	AddModule(m Module) Interface
}

type subGenerator struct {
	parent Interface
	sub    string
}

func newSubGenerator(parent Interface, sub string) Interface {
	if sub == "" {
		return parent
	}
	return subGenerator{
		parent: parent,
		sub:    sub,
	}
}

func (s subGenerator) Sub(dirname string) Interface {
	return newSubGenerator(s, dirname)
}

func (s subGenerator) Add(name string, handler Handler, middlewares ...Middleware) Interface {
	s.parent.Add(filepath.Join(s.sub, name), handler, middlewares...)
	return s
}

func (s subGenerator) AddFile(f *File) Interface {
	if f == nil {
		return s
	}
	f1 := f.Clone()
	if s.sub != "" {
		f1.Name = filepath.Join(s.sub, f1.Name)
	}
	s.parent.AddFile(f1)
	return s
}

func (s subGenerator) AddFiles(files []*File) Interface {
	for _, f := range files {
		s.AddFile(f)
	}
	return s
}

func (s subGenerator) Also(f func(i Interface)) Interface {
	if f != nil {
		f(s)
	}
	return s
}

func (s subGenerator) AddModule(m Module) Interface {
	if m != nil {
		m.Apply(s)
	}
	return s
}
