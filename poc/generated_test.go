package main

import "github.com/dcaiafa/hammock"

type MockFoo struct {
	*hammock.Mock
}

func NewMockFoo(c *hammock.Controller) *MockFoo {
	return &MockFoo{Mock: c.NewMock()}
}

func (m *MockFoo) DoStuff(i int, s *Struct) (int, error) {
	res := m.Call("DoStuff", []any{i, s})
	r1 := hammock.Get[int](res, 0)
	r2 := hammock.Get[error](res, 1)
	return r1, r2
}

type mockFooDoStuff struct {
	e *hammock.Expectation
}

func (r *mockFooDoStuff) Return(r1 int, r2 error) *mockFooDoStuff {
	r.e.Return([]any{r1, r2})
	return r
}

func (r *mockFooDoStuff) Times(n int) *mockFooDoStuff {
	r.e.Times(n)
	return r
}

func Expect_Foo_DoStuff[
	A1 int | hammock.BasicMatchers | hammock.MatcherT[int],
	A2 *Struct | hammock.BasicMatchers | hammock.MatcherT[*Struct],
](m *MockFoo, a1 A1, a2 A2) *mockFooDoStuff {
	args := []any{a1, a2}
	e := m.Expect("DoStuff", args)
	return &mockFooDoStuff{e}
}
