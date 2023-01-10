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
	return newMock(c)
}

func (c *Controller) Finish() {
	for _, m := range c.mocks {
		m.checkExpectations()
	}
}

func (c *Controller) Failf(msg string, args ...any) {
	c.t.Fatalf(msg, args...)
}

func (c *Controller) Logf(msg string, args ...any) {
	c.t.Logf(msg, args...)
}
