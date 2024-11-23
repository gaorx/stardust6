package sdtemplate

var (
	// Text 模板执行器，不对HTML中的特殊字符进行转义
	Text Executor = textExecutor{}
	// Html 模板执行器，对HTML中的特殊字符进行转义
	Html Executor = htmlExecutor{}
)
