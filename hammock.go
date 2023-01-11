package hammock

type Test interface {
	Logf(msg string, args ...any)
	Fatalf(msg string, args ...any)
	Cleanup(f func())
}

func Get[T any](r []any, n int) T {
	var ret T
	if r[n] != nil {
		ret = r[n].(T)
	}
	return ret
}
