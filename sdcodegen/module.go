package sdcodegen

// Module 模块接口，可以在其中在代码生成器中生成一组文件
type Module interface {
	Apply(i Interface)
}

// ModuleFunc 以一个函数形式呈现的模块
type ModuleFunc func(i Interface)

func (f ModuleFunc) Apply(i Interface) {
	if f != nil {
		f(i)
	}
}
