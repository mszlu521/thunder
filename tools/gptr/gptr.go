package gptr

func Of[T any](v T) *T {
	return &v
}
