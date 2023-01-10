package main

import (
	"fmt"

	"github.com/dcaiafa/damock"
)

type Struct struct {
	A int
}

type Foo interface {
	DoStuff(i int) int
}

type Expectation struct {
	Args []any
}

type AnythingType struct{}

var Anything AnythingType

type MockFoo struct {
	*damock.Mock
}

func (m *MockFoo) DoStuff(i *Struct) int {
	panic("not implemented")
}

type MockFooDoStuffResult struct{}

func (r *MockFooDoStuffResult) Return(i int) {
}

func MockFooDoStuff[A int | AnythingType | Matcher[int]](m *MockFoo, i A) {
	fmt.Println(i)
}

func TestStuff() {
	var m MockFoo

	MockFooDoStuff(&m, &Struct{1})
	MockFooDoStuff(&m, Matcher[*Struct](func(s *Struct) bool { return true }))
	MockFooDoStuff(&m, Anything)

	MockFooDoFoo(&m, 123)
	MockFooDoFoo(&m, Matcher[int](func(int) bool { return true }))
	MockFooDoFoo(&m, Anything)
}
