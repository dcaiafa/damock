package main

import (
	"fmt"
	"testing"

	"github.com/dcaiafa/hammock"
)

func TestStuff(t *testing.T) {
	c := hammock.NewController(t)
	defer c.Finish()

	m := NewMockFoo(c)
	Expect_Foo_DoStuff(m, 1, hammock.Any).
		Times(1).
		Return(1, nil)

	Expect_Foo_DoStuff(m, 1, hammock.MatcherT[*Struct](func(s *Struct) bool { return (s.A+s.B)%2 == 0 })).
		Times(1).
		Do(func(i int, s *Struct) (int, error) {
			return i + s.A + s.B, nil
		})

	i, err := m.DoStuff(1, &Struct{3, 4})
	fmt.Println(i, err)

	i, err = m.DoStuff(1, &Struct{2, 6})
	fmt.Println(i, err)
}
