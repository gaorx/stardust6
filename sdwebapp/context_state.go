package sdwebapp

func (c Context) State() any {
	return c.Get(akState)
}

func (c Context) SetState(v any) {
	c.Set(akState, v)
}
