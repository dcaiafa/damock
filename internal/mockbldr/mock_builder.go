package mockbldr

import (
	"fmt"
	"go/types"
)

type Parser interface {
	GetInterface(name string) (types.Object, error)
}

type Mock struct {
	Methods []*types.Func
}

func (m *Mock) Dump() {
	for _, method := range m.Methods {
		fmt.Println(method)
		dumpTypes(method)
	}
}

func Build(parser Parser, ifaceName string) (*Mock, error) {
	o, err := parser.GetInterface(ifaceName)
	if err != nil {
		return nil, err
	}

	iface := o.Type().Underlying().(*types.Interface)

	m := &Mock{}

	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		m.Methods = append(m.Methods, method)
	}

	return m, nil
}
