package damock

import (
	"fmt"
	"reflect"
	"strings"
)

type Matcher struct {
	Match    func(v any) bool
	ToString func() string
}

func (m *Matcher) String() string {
	if m.ToString != nil {
		return m.ToString()
	}
	return "<Matcher>"
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
