package damock

type Mock struct {
	controller *Controller
	methods    map[string][]*Expectation
}

func newMock(c *Controller) *Mock {
	m := &Mock{
		controller: c,
		methods:    make(map[string][]*Expectation),
	}
	return m
}

func (m *Mock) On(method string, args []any) {
	expectations := m.methods[method]
	if expectations == nil {
		m.controller.Failf("There is no expectation for method %q", method)
		return
	}

	foundMatch := false
	for _, e := range expectations {
		if e.isOpen() {
			continue
		}
		if e.isMatch(args) {
			e.calls++
			foundMatch = true
			break
		}
	}

	if !foundMatch {
		var open []*Expectation
		for _, e := range expectations {
			if e.isOpen() {
				open = append(open, e)
			}
		}
		if len(open) > 0 {
			m.controller.Logf("On call:")
			m.controller.Logf("  %v(%v)", method, formatArgs(args))
			m.controller.Logf("Open expectations:")
			for _, e := range open {
				m.controller.Logf("  %v(%v)", method, formatArgs(e.args))
			}
		}
	}
}

func (m *Mock) checkExpectations() {
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
	if !first {
		m.controller.Failf("Mock has unsatisfied expectations")
	}
}
