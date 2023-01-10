package hammock

import "testing"

type Controller struct {
	t     testing.TB
	mocks []*Mock
}

func NewController(t testing.TB) *Controller {
	return &Controller{t: t}
}

func (c *Controller) NewMock() *Mock {
	m := newMock(c)
	c.mocks = append(c.mocks, m)
	return m
}

func (c *Controller) Finish() {
	c.t.Helper()
	for _, m := range c.mocks {
		m.checkExpectations()
	}
}

func (c *Controller) Failf(msg string, args ...any) {
	c.t.Helper()
	c.t.Fatalf(msg, args...)
}

func (c *Controller) Logf(msg string, args ...any) {
	c.t.Helper()
	c.t.Logf(msg, args...)
}
