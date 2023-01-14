package mockbldr

import (
	"fmt"
	"go/types"
)

type Parser interface {
	GetInterface(name string) (types.Object, error)
}

type Mock struct {
	Name    string
	Methods []*types.Func
}

func (m *Mock) Dump() {
	for _, method := range m.Methods {
		fmt.Println(method)
		dumpObject(0, method)
	}
}

func Build(parser Parser, ifaceName string) (*Mock, error) {
	o, err := parser.GetInterface(ifaceName)
	if err != nil {
		return nil, err
	}

	iface := o.Type().Underlying().(*types.Interface)

	m := &Mock{
		Name: o.Name(),
	}

	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		m.Methods = append(m.Methods, method)
	}

	return m, nil
}

func printf(lvl int, msg string, args ...any) {
	for i := 0; i < lvl; i++ {
		fmt.Print(" ")
	}
	fmt.Printf(msg, args...)
	fmt.Print("\n")
}

func dumpObject(lvl int, o types.Object) {
	printf(lvl, "Object: %v", o.Name())
	dumpType(lvl+1, o.Type())
}

func dumpType(lvl int, t types.Type) {
	printf(lvl, "Type: %v", t.String())

	u := t.Underlying()

	if u != t {
		printf(lvl, "Underlying:")
		dumpType(lvl+1, t.Underlying())
	}

	switch t := t.(type) {
	case *types.Signature:
		printf(lvl, "Signature:")
		printf(lvl, "Params:")
		dumpTuple(lvl+1, t.Params())
		printf(lvl, "Results:")
		dumpTuple(lvl+1, t.Results())

	case *types.Named:
		pkgPath := "<global>"
		if pkg := t.Obj().Pkg(); pkg != nil {
			pkgPath = pkg.Path()
		}
		printf(lvl, "Named: %v %v", pkgPath, t.Obj().Name())

	case *types.Pointer:
		printf(lvl, "Pointer:")
		dumpType(lvl+1, t.Elem())
	}
}

func dumpTuple(lvl int, t *types.Tuple) {
	for i := 0; i < t.Len(); i++ {
		v := t.At(i)
		printf(lvl, "%v: %v", i, v.Name())
		dumpType(lvl+1, v.Type())
	}
}
