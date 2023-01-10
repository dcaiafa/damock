package mockgears

func Get[T any](r []any, n int) T {
	var ret T
	if r[n] != nil {
		ret = r[n].(T)
	}
	return ret
}
