package hammock

import "github.com/dcaiafa/hammock/match"

type Expectation struct {
	args     []any
	results  []any
	minCalls int
	maxCalls int
	do       func(args []any) []any
	calls    int
}

func newExpectation(args []any) *Expectation {
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

func (e *Expectation) Do(do func(args []any) []any) {
	e.do = do
}

func (e *Expectation) isOpen() bool {
	return e.calls < e.maxCalls
}

func (e *Expectation) isSatisfied() bool {
	return e.calls >= e.minCalls
}

func (e *Expectation) isMatch(args []any) bool {
	if len(args) != len(e.args) {
		return false
	}
	for i := range args {
		if !match.IsMatch(e.args[i], args[i]) {
			return false
		}
	}
	return true
}
