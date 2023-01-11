package hammock

import (
	"fmt"
	"strings"
	"sync"
)

type Mock struct {
	controller *Controller

	mu      sync.Mutex
	methods map[string][]*Expectation
}

func newMock(c *Controller) *Mock {
	m := &Mock{
		controller: c,
		methods:    make(map[string][]*Expectation),
	}
	return m
}

func (m *Mock) Expect(method string, args []any) *Expectation {
	m.mu.Lock()
	defer m.mu.Unlock()

	e := newExpectation(args)
	m.methods[method] = append(m.methods[method], e)
	return e
}

func (m *Mock) Call(method string, args []any) []any {
	m.controller.t.Helper()

	m.mu.Lock()

	expectations := m.methods[method]
	if expectations == nil {
		m.mu.Unlock()
		m.controller.Failf("Method %q has no expectations", method)
		return nil
	}

	for _, e := range expectations {
		if !e.isOpen() {
			continue
		}

		if e.isMatch(args) {
			e.calls++

			// Must release the mutex before processing any callbacks as this can
			// result in reentrancy.
			m.mu.Unlock()

			if e.do != nil {
				return e.do(args)
			}
			return e.results
		}
	}

	// Failed to find a matching expectation.
	// Log and fail.
	var open []*Expectation
	for _, e := range expectations {
		if e.isOpen() {
			open = append(open, e)
		}
	}
	if len(open) > 0 {
		m.controller.Logf("Calling:")
		m.controller.Logf("  %v(%v)", method, formatArgs(args))
		m.controller.Logf("Open expectations:")
		for _, e := range open {
			m.controller.Logf("  %v(%v)", method, formatArgs(e.args))
		}
	}

	// Must release the mutex before Failf because it can panic.
	m.mu.Unlock()
	m.controller.Failf("No matching expectations")

	return nil
}

func (m *Mock) checkExpectations() {
	m.controller.t.Helper()

	m.mu.Lock()

	first := true
	for method, expectations := range m.methods {
		for _, e := range expectations {
			if !e.isSatisfied() {
				if first {
					m.controller.Logf("Unsatisfied expectations:")
					first = false
				}
				m.controller.Logf("  %v(%v)", method, formatArgs(e.args))
			}
		}
	}

	m.mu.Unlock()

	if !first {
		m.controller.Failf("Mock has unsatisfied expectations")
	}
}

func formatArgs(args []any) string {
	var b strings.Builder
	for i, arg := range args {
		if i != 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%+v", arg)
	}
	return b.String()
}
