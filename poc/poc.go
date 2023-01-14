package main

// Struct is very special.
type Struct struct {
	A int
	B int
}

type Foo interface {
	DoStuff(i int, s *Struct) (int, error)
}
