package main

import (
	"github.com/dcaiafa/hammock"
	"github.com/dcaiafa/hammock/match"
)

type MockFoo struct {
	*hammock.Mock
}

func NewMockFoo(t hammock.Test) *MockFoo {
	return &MockFoo{Mock: hammock.NewMock(t)}
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

func (r *mockFooDoStuff) Do(f func(i int, s *Struct) (int, error)) {
	r.e.Do(func(args []any) []any {
		r1, r2 := f(
			hammock.Get[int](args, 0),
			hammock.Get[*Struct](args, 1),
		)
		return []any{r1, r2}
	})
}

func Expect_Foo_DoStuff[
	A1 int | match.BasicMatchers | match.CustomType[int],
	A2 *Struct | match.BasicMatchers | match.CustomType[*Struct],
](m *MockFoo, a1 A1, a2 A2) *mockFooDoStuff {
	args := []any{a1, a2}
	e := m.Expect("DoStuff", args)
	return &mockFooDoStuff{e}
}
