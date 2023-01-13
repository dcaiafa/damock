package parser

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Parser struct {
	packages map[string]*packages.Package
}

func NewAnalyzer() *Parser {
	return &Parser{
		packages: make(map[string]*packages.Package),
	}
}

func (a *Parser) GetInterface(name string) (types.Object, error) {
	packageName, ifaceName, err := ParseFQName(name)
	if err != nil {
		return nil, err
	}

	pkg := a.packages[packageName]
	if pkg == nil {
		cfg := &packages.Config{Mode: packages.NeedName | packages.NeedTypes}
		pkgs, err := packages.Load(cfg, packageName)
		if err != nil {
			return nil, fmt.Errorf("failed to load package %q", packageName)
		}
		if len(pkgs) != 1 {
			return nil, fmt.Errorf(
				"there were %d packages matching %q", len(pkgs), packageName)
		}
		pkg = pkgs[0]
		a.packages[packageName] = pkg
	}

	o := pkg.Types.Scope().Lookup(ifaceName)
	if o == nil {
		return nil, fmt.Errorf(
			"package %q does not have interface %q", packageName, ifaceName)
	}

	iface, ok := o.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, fmt.Errorf("%q is not an interface", name)
	}

	if !iface.IsMethodSet() {
		return nil, fmt.Errorf("%q is not a method set interface", name)
	}

	return o, nil
}

func ParseFQName(name string) (string, string, error) {
	period := strings.LastIndex(name, ".")
	if period == -1 {
		return "", "", fmt.Errorf("name is not fully qualified")
	}
	packageName := name[:period]
	ifaceName := name[period+1:]
	return packageName, ifaceName, nil
}
