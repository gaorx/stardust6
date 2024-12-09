package sdcheck

func funcOf(c Interface) Func {
	if f, ok := c.(Func); ok {
		return f
	} else {
		return func() error {
			return c.Check()
		}
	}
}
