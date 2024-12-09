package sdcheck

// Interface 检查接口
type Interface interface {
	// Check 检查是否满足条件，如果error不为nil，则表示不满足条件
	Check() error
}
