package sderr

func ensurePtr[T any](p *T) *T {
	if p == nil {
		return new(T)
	}
	return p
}
