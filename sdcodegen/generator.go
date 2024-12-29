package sdcodegen

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"io/fs"
)

// Generator 代码生成器，可以生成一个或多个文件
type Generator struct {
	middlewares []Middleware
	entries     []entry
	files       []*File
	logger      Logger
}

type entry struct {
	name        string
	handler     Handler
	middlewares []Middleware
}

var _ Interface = (*Generator)(nil)

// New 创建一个代码生成器
func New() *Generator {
	return &Generator{}
}

// SetLogger 设置日志
func (g *Generator) SetLogger(logger Logger) *Generator {
	if logger == nil {
		logger = NoLog()
	}
	g.logger = logger
	return g
}

// Use 添加全局中间件
func (g *Generator) Use(middlewares ...Middleware) *Generator {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

// Sub 通过一个dirname生成新的Generator，在其下的所有代码生成都会在这个dirname下
func (g *Generator) Sub(dirname string) Interface {
	return newSubGenerator(g, dirname)
}

// Add 添加一个生成器
func (g *Generator) Add(name string, handler Handler, middlewares ...Middleware) Interface {
	e := entry{
		name:        name,
		handler:     handler,
		middlewares: middlewares,
	}
	g.entries = append(g.entries, e)
	return g
}

// AddFile 加入一个文件
func (g *Generator) AddFile(f *File) Interface {
	if f == nil {
		return g
	}
	g.files = append(g.files, f.Clone())
	return g
}

// AddFiles 加入多个文件
func (g *Generator) AddFiles(files []*File) Interface {
	for _, f := range files {
		g.AddFile(f)
	}
	return g
}

// Also 在当前生成器上执行一个函数，并返回当前生成器
func (g *Generator) Also(f func(i Interface)) Interface {
	if f != nil {
		f(g)
	}
	return g
}

// AddModule 添加一个模块
func (g *Generator) AddModule(m Module) Interface {
	if m != nil {
		m.Apply(g)
	}
	return g
}

// TryOne 尝试生成一个文件，但实际不写入
func (g *Generator) TryOne(name string, current *File) (*File, error) {
	e, ok := lo.Find(g.entries, func(e entry) bool {
		return e.name == name
	})
	if !ok {
		for _, f := range g.files {
			if f.Name == name {
				return f.Clone(), nil
			}
		}
		return nil, sderr.With("file", name).Newf("not found generator")
	}
	middlewares := mergeMiddlewares([]Middleware{SetLogger(g.logger)}, g.middlewares, e.middlewares)
	return Generate(e.name, current, e.handler, middlewares...)
}

// TryAll 尝试生成所有文件，但实际不写入
func (g *Generator) TryAll(currents Files) (Files, error) {
	var files Files
	for _, e := range g.entries {
		current := currents.Get(e.name)
		middlewares := mergeMiddlewares([]Middleware{SetLogger(g.logger)}, g.middlewares, e.middlewares)
		f, err := Generate(e.name, current, e.handler, middlewares...)
		if err != nil {
			return nil, sderr.Wrap(err)
		}
		files = append(files, f)
	}
	for _, f := range g.files {
		files = append(files, f.Clone())
	}
	return files, nil
}

// TryFS 尝试从指定的fs中读取文件，然后生成所有文件，但实际不写入
func (g *Generator) TryFS(fsys fs.FS, nameFilters ...StringFilter) (Files, error) {
	currents, err := ReadDirFS(fsys, nameFilters...)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return g.TryAll(currents)
}

// TryDir 尝试从指定的目录中读取文件，然后生成所有文件，但实际不写入
func (g *Generator) TryDir(dirname string, nameFilters ...StringFilter) (Files, error) {
	currents, err := ReadDir(dirname, nameFilters...)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return g.TryAll(currents)
}

// GenerateOne 生成一个文件
// 尝试先读取指定目录中的文件，然后生成新内容，最后写入到这个文件中
func (g *Generator) GenerateOne(dirname string, name string) error {
	absDirname, ok := toAbs(dirname)
	if !ok {
		return sderr.With("dir", dirname).Newf("illegall dirname")
	}
	current, err := ReadFile(dirname, name)
	if err != nil && !IsNotExistErr(err) {
		return sderr.Wrap(err)
	}
	f, err := g.TryOne(name, current)
	if err != nil {
		return sderr.Wrap(err)
	}
	err = f.Write(absDirname, g.logger)
	return sderr.Wrap(err)
}

// Generate 生成所有文件，并写入到目录中
// 先尝试读取文件，之后再生成新内容，最后写入到这些文件中
func (g *Generator) Generate(dirname string, nameFilters ...StringFilter) error {
	absDirname, ok := toAbs(dirname)
	if !ok {
		return sderr.With("dir", dirname).Newf("illegall dirname")
	}
	files, err := g.TryDir(absDirname, nameFilters...)
	if err != nil {
		return sderr.Wrap(err)
	}
	err = files.Write(absDirname, g.logger)
	return sderr.Wrap(err)
}
