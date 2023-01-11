package hammock

import (
	"sync"
)

type Logger interface {
	Logf(msg string, args ...any)
	Fatalf(msg string, args ...any)
}

type Controller struct {
	logger Logger

	mu    sync.Mutex
	mocks []*Mock
}

func NewController(logger Logger) *Controller {
	return &Controller{logger: logger}
}

func (c *Controller) NewMock() *Mock {
	m := newMock(c.logger)

	c.mu.Lock()
	c.mocks = append(c.mocks, m)
	c.mu.Unlock()

	return m
}

func (c *Controller) Finish() {
	c.mu.Lock()
	mocks := c.mocks
	c.mocks = nil
	c.mu.Unlock()

	for _, m := range mocks {
		m.checkExpectations()
	}
}
