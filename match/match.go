package match

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

type matcher interface {
	Match(v any) bool
}

type CustomType[T any] func(v T) bool

func (m CustomType[T]) Match(v any) bool {
	return m(v.(T))
}

func (m CustomType[T]) String() string {
	return "<matcher>"
}

func Custom[T any](f func(v T) bool) CustomType[T] {
	return CustomType[T](f)
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

func GT[T constraints.Ordered](v T) CustomType[T] {
	return CustomType[T](func(x T) bool { return x > v })
}

func GE[T constraints.Ordered](v T) CustomType[T] {
	return CustomType[T](func(x T) bool { return x >= v })
}

func LT[T constraints.Ordered](v T) CustomType[T] {
	return CustomType[T](func(x T) bool { return x < v })
}

func LE[T constraints.Ordered](v T) CustomType[T] {
	return CustomType[T](func(x T) bool { return x <= v })
}

func IsMatch(expected, actual any) bool {
	if m, ok := expected.(matcher); ok {
		return m.Match(actual)
	}
	return reflect.DeepEqual(expected, actual)
}
