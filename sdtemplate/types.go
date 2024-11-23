package sdtemplate

// Executor 模板执行器，通过模板和数据生成字符串
type Executor interface {
	// Exec 通过模板生成字符串
	Exec(template string, data any) (string, error)
	// ExecOr 通过模板生成字符串，如果出错则返回默认值
	ExecOr(template string, data any, def string) string
	// MustExec 通过模板生成字符串，如果出错则 panic
	MustExec(template string, data any) string
}
