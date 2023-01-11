package main

import (
	"flag"
	"fmt"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

func main() {
	flag.Parse()

	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedTypes}
	pkgs, err := packages.Load(cfg, flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		dump(pkg)
	}
}

func dump(pkg *packages.Package) {
	fmt.Println(pkg.PkgPath)
	pkgType := pkg.Types

	for _, n := range pkgType.Scope().Names() {
		o := pkgType.Scope().Lookup(n)
		if !o.Exported() {
			continue
		}
		iface, ok := o.Type().Underlying().(*types.Interface)
		if !ok {
			continue
		}
		if !iface.IsMethodSet() {
			continue
		}
		fmt.Println("  ", n)
		for i := 0; i < iface.NumMethods(); i++ {
			m := iface.Method(i)
			sig := m.Type().(*types.Signature)
			fmt.Println("    ", m.FullName(), sig)
			params := sig.Params()
			for j := 0; j < params.Len(); j++ {
				t := params.At(j)
				fmt.Println("      ", t.Name(), t.Type().String())
			}
		}
	}
}
