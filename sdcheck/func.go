package sdcheck

// Func 函数检查
type Func func() error

// Check 实现 Interface.Check
func (f Func) Check() error {
	if f == nil {
		return nil
	}
	return f()
}
