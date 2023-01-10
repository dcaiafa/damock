package damock

type Expectation struct {
	args     []any
	results  []any
	minCalls int
	maxCalls int
	on       func(args []any) []any
	calls    int
}

func NewExpectation(args []any) *Expectation {
	return &Expectation{
		args:     args,
		maxCalls: 1,
	}
}

func (e *Expectation) Return(results []any) {
	e.results = results
}

func (e *Expectation) Times(n int) {
	e.minCalls = n
	e.maxCalls = n
}

func (e *Expectation) On(on func(args []any) []any) {
	e.on = on
}

func (e *Expectation) isOpen() bool {
	return e.calls < e.maxCalls
}

func (e *Expectation) isSatisfied() bool {
	return e.calls > e.minCalls
}

func (e *Expectation) isMatch(args []any) bool {
	if len(args) != len(e.args) {
		return false
	}
	for i := range args {
		if !isMatch(e.args[i], args[i]) {
			return false
		}
	}
	return true
}
