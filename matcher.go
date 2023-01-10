package hammock

import (
	"fmt"
	"reflect"
	"strings"
)

type Matcher[T any] func(v T) bool

func (m Matcher[T]) ToMatcherAny() Matcher[any] {
	return Matcher[any](
		func(v any) bool {
			t := v.(T)
			return m(t)
		},
	)
}

type AnyType struct{}

func (a AnyType) String() string {
	return "<any>"
}

var Any AnyType

func isMatch(expected, actual any) bool {
	switch expected := expected.(type) {
	case Matcher[any]:
		return expected(actual)
	case AnyType:
		return true
	default:
		return reflect.DeepEqual(expected, actual)
	}
}

func formatArgs(args []any) string {
	var b strings.Builder
	for i, arg := range args {
		if i != 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%+v", arg)
	}
	return b.String()
}
