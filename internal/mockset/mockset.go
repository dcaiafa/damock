package mockset

import (
	"go/types"
)

type Parser interface {
	GetInterface(name string) (types.Object, error)
}

type Mock struct {
	IfaceName string
	IfaceObjs map[string]types.Object
	Methods   []*types.Func
}

type MockSet struct {
	parser Parser
	mocks  map[string]*Mock // key is fq interface name
}

func NewMockSet(parser Parser) *MockSet {
	return &MockSet{
		parser: parser,
		mocks:  make(map[string]*Mock),
	}
}

func (s *MockSet) AddMock(ifaceName string) error {
	ifaceObj, err := s.parser.GetInterface(ifaceName)
	if err != nil {
		return err
	}

	m := &Mock{
		IfaceName: ifaceName,
		IfaceObjs: make(map[string]types.Object),
	}

	s.addIface(m, ifaceObj)

	return nil
}

func (s *MockSet) addIface(m *Mock, ifaceObj types.Object) {
	if m.IfaceObjs[ifaceObj.Name()] != nil {
		return
	}

	iface := ifaceObj.Type().Underlying().(*types.Interface)
	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
	}
}
