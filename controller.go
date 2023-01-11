package hammock

import (
	"sync"
	"testing"
)

type Controller struct {
	t testing.TB

	mu    sync.Mutex
	mocks []*Mock
}

func NewController(t testing.TB) *Controller {
	return &Controller{t: t}
}

func (c *Controller) NewMock() *Mock {
	m := newMock(c)

	c.mu.Lock()
	c.mocks = append(c.mocks, m)
	c.mu.Unlock()

	return m
}

func (c *Controller) Finish() {
	c.t.Helper()

	c.mu.Lock()
	mocks := c.mocks
	c.mocks = nil
	c.mu.Unlock()

	for _, m := range mocks {
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
