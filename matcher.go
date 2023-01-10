package hammock

import (
	"fmt"
	"reflect"
	"strings"
)

type Matcher interface {
	Match(v any) bool
}

type MatcherT[T any] func(v T) bool

func (m MatcherT[T]) Match(v any) bool {
	return m(v.(T))
}

func (m MatcherT[T]) String() string {
	return "<custom-matcher>"
}

type AnyType struct{}

func (a AnyType) Match(v any) bool {
	return true
}

func (a AnyType) String() string {
	return "any"
}

var Any AnyType

type NotNilType struct{}

func (t NotNilType) Match(v any) bool {
	return v != nil
}

func (t NotNilType) String() string {
	return "!=nil"
}

var NotNil NotNilType

type BasicMatchers interface {
	AnyType | NotNilType
}

func isMatch(expected, actual any) bool {
	if m, ok := expected.(Matcher); ok {
		return m.Match(actual)
	}
	return reflect.DeepEqual(expected, actual)
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
