package emitter

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"strings"

	"github.com/dcaiafa/hammock/internal/mockbldr"
)

const hammockPackage = "github.com/dcaiafa/hammock"
const matchPackage = "github.com/dcaiafa/hammock/match"

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
	hammockPackageAlias := e.packageAlias(hammockPackage)

	fmt.Fprintf(e.bodyBuf, "\n")
	fmt.Fprintf(e.bodyBuf, "func (m *%v) %v(",
		e.mockName(m),
		method.Name(),
	)
	sig := method.Type().(*types.Signature)
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		if i != 0 {
			fmt.Fprintf(e.bodyBuf, ", ")
		}
		fmt.Fprintf(e.bodyBuf, "%v %v", e.argStr(nil, i), e.typeStr(p.Type()))
	}
	fmt.Fprintf(e.bodyBuf, ") ")

	results := sig.Results()
	if results.Len() > 0 {
		fmt.Fprintf(e.bodyBuf, "(")
		for i := 0; i < results.Len(); i++ {
			r := results.At(i)
			if i != 0 {
				fmt.Fprintf(e.bodyBuf, ", ")
			}
			fmt.Fprintf(e.bodyBuf, "%v %v", e.resultStr(nil, i), e.typeStr(r.Type()))
		}
		fmt.Fprintf(e.bodyBuf, ") ")
	}

	paramNames := func() []string {
		names := make([]string, params.Len())
		for i := 0; i < params.Len(); i++ {
			names[i] = e.argStr(nil, i)
		}
		return names
	}

	fmt.Fprintf(e.bodyBuf, "{\n")
	fmt.Fprintf(e.bodyBuf, "  ")
	if results.Len() > 0 {
		fmt.Fprintf(e.bodyBuf, "res := ")
	}
	fmt.Fprintf(
		e.bodyBuf, "m.Call(%q, []any{%v})\n",
		method.Name(),
		strings.Join(paramNames(), ", "))

	for i := 0; i < results.Len(); i++ {
		r := results.At(i)
		fmt.Fprintf(
			e.bodyBuf, "  %v = %v.Get[%v](res, %d)\n",
			e.resultStr(nil, i), hammockPackageAlias, e.typeStr(r.Type()), i)
	}

	if results.Len() > 0 {
		fmt.Fprintf(e.bodyBuf, "  return\n")
	}
	fmt.Fprintf(e.bodyBuf, "}\n")

	fmt.Fprintf(e.bodyBuf, "\n")
	e.emitExpectation(m, method)
}

func (e *Emitter) emitExpectation(m *mockbldr.Mock, method *types.Func) {
	hammockPackageAlias := e.packageAlias(hammockPackage)
	matchPackageAlias := e.packageAlias(matchPackage)

	sig := method.Type().(*types.Signature)
	params := sig.Params()
	results := sig.Results()

	// Expectation struct:
	expectationName := e.expectationName(m, method)
	fmt.Fprintf(e.bodyBuf, "type %v struct {\n", expectationName)
	fmt.Fprintf(e.bodyBuf, "  e *%v.Expectation\n", e.packageAlias(hammockPackage))
	fmt.Fprintf(e.bodyBuf, "}\n")
	fmt.Fprintf(e.bodyBuf, "\n")

	// Return:
	fmt.Fprintf(e.bodyBuf, "func (e *%v) Return(%v) *%v {\n",
		expectationName, e.varNameAndTypePairs(results, e.resultStr), expectationName)
	fmt.Fprintf(e.bodyBuf, "  e.e.Return([]any{%v})\n", e.varNames(results, e.resultStr))
	fmt.Fprintf(e.bodyBuf, "  return e\n")
	fmt.Fprintf(e.bodyBuf, "}\n")
	fmt.Fprintf(e.bodyBuf, "\n")

	// Do:
	fmt.Fprintf(e.bodyBuf, "func (e *%v) Do(f func(%v) (%v)) *%v {\n",
		expectationName, e.varNameAndTypePairs(params, e.argStr),
		e.varTypes(results), expectationName)
	fmt.Fprintf(e.bodyBuf, "  e.e.Do(func(args []any) []any {\n")
	rets := e.varNames(results, e.resultStr)
	fmt.Fprintf(e.bodyBuf, "    ")
	if rets != "" {
		fmt.Fprintf(e.bodyBuf, "%v := ", rets)
	}
	fmt.Fprintf(e.bodyBuf, "f(\n")
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		fmt.Fprintf(e.bodyBuf, "      %v.Get[%v](args, %d),\n",
			hammockPackageAlias, e.typeStr(p.Type()), i)
	}
	fmt.Fprintf(e.bodyBuf, "    )\n")
	fmt.Fprintf(e.bodyBuf, "    return []any{%v}\n", rets)
	fmt.Fprintf(e.bodyBuf, "  })\n")
	fmt.Fprintf(e.bodyBuf, "  return e\n")
	fmt.Fprintf(e.bodyBuf, "}\n")
	fmt.Fprintf(e.bodyBuf, "\n")

	// Times:
	fmt.Fprintf(e.bodyBuf, "func (e *%v) Times(n int) *%v {\n", expectationName, expectationName)
	fmt.Fprintf(e.bodyBuf, "  e.e.Times(n)\n")
	fmt.Fprintf(e.bodyBuf, "  return e\n")
	fmt.Fprintf(e.bodyBuf, "}\n")
	fmt.Fprintf(e.bodyBuf, "\n")

	// Expect:
	fmt.Fprintf(e.bodyBuf, "func Expect_%v_%v[\n", m.Name, method.Name())
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		fmt.Println(p)
		if _, ok := p.Type().(*types.Interface); ok {
			fmt.Fprintf(e.bodyBuf, "  A%d any,\n", i)
		} else {
			fmt.Fprintf(e.bodyBuf, "  A%d %v | %v.BasicMatchers | %v.CustomType[%v],\n",
				i, e.typeStr(p.Type()), matchPackageAlias, matchPackageAlias, e.typeStr(p.Type()))
		}
	}
	fmt.Fprintf(e.bodyBuf, "](\n")
	fmt.Fprintf(e.bodyBuf, "  m *%v,", e.mockName(m))
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		fmt.Fprintf(e.bodyBuf, "  %v A%d,\n", e.argStr(p, i), i)
	}
	fmt.Fprintf(e.bodyBuf, ") *%v {\n", expectationName)
	fmt.Fprintf(e.bodyBuf, "  return nil\n")
	fmt.Fprintf(e.bodyBuf, "}\n")
}

func (e *Emitter) varNames(t *types.Tuple, nameFunc func(v *types.Var, n int) string) string {
	res := make([]string, t.Len())
	for i := 0; i < t.Len(); i++ {
		res[i] = nameFunc(t.At(i), i)
	}
	return strings.Join(res, ", ")
}

func (e *Emitter) varTypes(t *types.Tuple) string {
	typs := make([]string, t.Len())
	for i := 0; i < t.Len(); i++ {
		typs[i] = e.typeStr(t.At(i).Type())
	}
	return strings.Join(typs, ", ")
}

func (e *Emitter) varNameAndTypePairs(t *types.Tuple, nameFunc func(v *types.Var, n int) string) string {
	nameAndTyps := make([]string, t.Len())
	for i := 0; i < t.Len(); i++ {
		v := t.At(i)
		nameAndTyps[i] = fmt.Sprintf("%v %v", nameFunc(v, i), e.typeStr(v.Type()))
	}
	return strings.Join(nameAndTyps, ", ")
}

func (e *Emitter) expectationName(m *mockbldr.Mock, method *types.Func) string {
	return e.mockName(m) + method.Name() + "Expectation"
}

func (e *Emitter) argStr(v *types.Var, n int) string {
	return fmt.Sprintf("a%d", n)
}

func (e *Emitter) resultStr(v *types.Var, n int) string {
	return fmt.Sprintf("r%d", n)
}

func (e *Emitter) signatureStr(sig *types.Signature) string {
	var str strings.Builder

	params := sig.Params()
	str.WriteString("(")
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		if i != 0 {
			str.WriteString(", ")
		}
		str.WriteString(e.typeStr(p.Type()))
	}
	str.WriteString(") ")

	results := sig.Results()
	if results.Len() > 1 {
		str.WriteString("(")
	}

	for i := 0; i < results.Len(); i++ {
		r := results.At(i)
		if i != 0 {
			str.WriteString(", ")
		}
		str.WriteString(e.typeStr(r.Type()))
	}

	if results.Len() > 1 {
		str.WriteString(")")
	}

	return str.String()
}

func (e *Emitter) typeStr(t types.Type) string {
	switch t := t.(type) {
	case *types.Basic:
		return t.String()
	case *types.Pointer:
		return "*" + e.typeStr(t.Elem())
	case *types.Named:
		name := t.Obj().Name()
		pkgPath := ""
		if pkg := t.Obj().Pkg(); pkg != nil {
			pkgPath = pkg.Path()
		}
		prefix := e.packageAlias(pkgPath)
		if prefix != "" {
			prefix += "."
		}
		return prefix + name
	case *types.Interface:
		var str strings.Builder
		str.WriteString("interface {")
		for i := 0; i < t.NumMethods(); i++ {
			method := t.Method(i)
			if i != 0 {
				str.WriteString("; ")
			}
			str.WriteString(method.Name())
			str.WriteString(e.signatureStr(method.Type().(*types.Signature)))
		}
		str.WriteString("}")
		return str.String()

	default:
		panic(fmt.Sprintf("type not supported: %v", t))
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

func (e *Emitter) packageAlias(pkgPath string) string {
	if pkgPath == "" || pkgPath == e.pkgPath {
		return ""
	}

	alias, ok := e.imports[pkgPath]
	if !ok {
		alias = fmt.Sprintf("p%d", len(e.imports))
		e.imports[pkgPath] = alias
	}
	return alias
}

func isGenericUnionCompatible(t types.Type) bool {
	switch t := t.(type) {
	case *types.Named:
		return isGenericUnionCompatible(t.Underlying())
	case *types.Interface:
		return false
	default:
		return true
	}
}
