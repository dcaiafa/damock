package emitter

import (
	"bytes"
	"fmt"
	"go/types"
	"io"

	"github.com/dcaiafa/hammock/internal/mockbldr"
)

const hammockPackage = "github.com/dcaiafa/hammock"

type Emitter struct {
	pkgPath string
	pkgName string
	bodyBuf *bytes.Buffer
	imports map[string]string
}

func New(pkgPath, pkgName string) *Emitter {
	e := &Emitter{
		pkgPath: pkgPath,
		pkgName: pkgName,
		bodyBuf: bytes.NewBuffer(make([]byte, 0, 10240)),
		imports: make(map[string]string),
	}
	e.packageAlias(hammockPackage)
	return e
}

func (e *Emitter) WriteMock(m *mockbldr.Mock) {
	fmt.Fprintf(e.bodyBuf, "type %v struct {\n", e.mockName(m))
	fmt.Fprintf(e.bodyBuf, "  *%v.Mock\n", e.packageAlias(hammockPackage))
	fmt.Fprintf(e.bodyBuf, "}\n")

	for _, method := range m.Methods {
		e.emitMethod(m, method)
	}
}

func (e *Emitter) mockName(m *mockbldr.Mock) string {
	return "mock" + m.Name
}

func (e *Emitter) emitMethod(m *mockbldr.Mock, method *types.Func) {
	fmt.Fprintf(e.bodyBuf, "\n")
	fmt.Fprintf(e.bodyBuf, "func (m *%v) %v(%v) {\n",
		e.mockName(m),
		method.Name(),
		"",
	)
	fmt.Fprintf(e.bodyBuf, "}\n")
}

func (e *Emitter) Finish(w io.Writer) error {
	headerBuf := bytes.NewBuffer(make([]byte, 0, 1024))
	fmt.Fprintf(headerBuf, "package %v\n\n", e.pkgName)

	for imp, alias := range e.imports {
		fmt.Fprintf(headerBuf, "import %v %q\n", alias, imp)
	}
	fmt.Fprintf(headerBuf, "\n\n")

	_, err := w.Write(headerBuf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(e.bodyBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (e *Emitter) packageAlias(pkgName string) string {
	alias, ok := e.imports[pkgName]
	if !ok {
		alias = fmt.Sprintf("p%d", len(e.imports))
		e.imports[pkgName] = alias
	}
	return alias
}
