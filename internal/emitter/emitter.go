package emitter

import (
	"bytes"
	"fmt"
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
	fmt.Fprintf(e.bodyBuf, "type %v struct {\n", "mock"+m.Name)
	fmt.Fprintf(e.bodyBuf, "  *%v.Mock\n", e.packageAlias(hammockPackage))
	fmt.Fprintf(e.bodyBuf, "}")

	for _, method := range m.Methods {
		_ = method
	}
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
