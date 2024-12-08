package sdlo

// EnsureSlice 返回一个非nil的slice
func EnsureSlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

// EnsureMap 返回一个非nil的map
func EnsureMap[M ~map[K]V, K comparable, V any](m M) M {
	if m == nil {
		return M{}
	}
	return m
}
