package main

import (
	"fmt"
	"testing"

	"github.com/dcaiafa/hammock"
)

type Struct struct {
	A int
}

type Foo interface {
	DoStuff(i int, s *Struct) (int, error)
}

type MockFoo struct {
	*hammock.Mock
}

func NewMockFoo(c *hammock.Controller) *MockFoo {
	return &MockFoo{Mock: c.NewMock()}
}

func (m *MockFoo) DoStuff(i int, s *Struct) (int, error) {
	res := m.Call("DoStuff", []any{i, s})
	r1 := res[0].(int)
	r2 := res[1].(error)
	return r1, r2
}

type mockFooDoStuff struct {
	e *hammock.Expectation
}

func (r *mockFooDoStuff) Return(r1 int, r2 error) {
	r.e.Return([]any{r1, r2})
}

func ExpectFooDoStuff[
	A1 int | hammock.AnyType | hammock.Matcher[int],
	A2 *Struct | hammock.AnyType | hammock.Matcher[*Struct],
](m *MockFoo, a1 A1, a2 A2) *mockFooDoStuff {
	args := [2]any{a1, a2}
	if m, ok := args[0].(hammock.Matcher[int]); ok {
		args[0] = m.ToMatcherAny()
	}
	if m, ok := args[1].(hammock.Matcher[*Struct]); ok {
		args[1] = m.ToMatcherAny()
	}
	e := m.Expect("DoStuff", args[:])
	return &mockFooDoStuff{e}
}

func TestStuff(t *testing.T) {
	c := hammock.NewController(t)
	defer c.Finish()

	m := NewMockFoo(c)
	ExpectFooDoStuff(m, 1, hammock.Any).Return(1, nil)

	i, err := m.DoStuff(1, &Struct{3})
	fmt.Println(i, err)
}
