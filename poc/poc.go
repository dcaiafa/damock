package main

type Struct struct {
	A int
	B int
}

type Foo interface {
	DoStuff(i int, s *Struct) (int, error)
}
