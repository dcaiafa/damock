package main

import (
	"bytes"
	"io"
)

// Struct is very special.
type Struct struct {
	A int
	B int
}

type Foo interface {
	DoStuff(i int, s *Struct, w io.Writer) (int, error)
	Interface(i interface {
		Foo() string
		Bar(*bytes.Buffer) error
		Baz()
	})
}
