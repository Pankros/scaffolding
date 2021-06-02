package main

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	"strings"
)

func main() {
	// 1. Handle arguments to command
	if len(os.Args) != 3 {
		failErr(fmt.Errorf("expected exactly one argument: <source type>"))
	}

	sourceType := os.Args[1]
	tableName := os.Args[2]
	sourceTypePackage, sourceTypeName := splitSourceType(sourceType)

	// 2. Inspect package and use type checker to infer imported types
	pkg := loadPackage(sourceTypePackage)

	// 3. Lookup the given source type name in the package declarations
	obj := pkg.Types.Scope().Lookup(sourceTypeName)
	if obj == nil {
		failErr(fmt.Errorf("%s not found in declared types of %s",
			sourceTypeName, pkg))
	}

	// 4. We check if it is a declared type
	if _, ok := obj.(*types.TypeName); !ok {
		failErr(fmt.Errorf("%v is not a named type", obj))
	}
	// 5. We expect the underlying type to be a struct
	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		failErr(fmt.Errorf("type %v is not a struct", obj))
	}

	err := generate(sourceTypeName, tableName, structType)
	if err != nil {
		failErr(err)
	}
	// 6. Now we can iterate through fields and access tags
	//printStruct(structType)
}

func printStruct(structType *types.Struct) {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tagValue := structType.Tag(i)
		if st, ok := field.Type().Underlying().(*types.Struct); ok {
			printStruct(st)
		}
		fmt.Println(field.Name(), tagValue, field.Type(), field.IsField())
	}
}

// github.com/mercadolibre/fury_payment-methods-write-v2/src/api/internal/model/entities.PaymentMethodType
// github.com/Pankros/scaffolding/generate.TestStruct
func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func splitSourceType(sourceType string) (string, string) {
	idx := strings.LastIndexByte(sourceType, '.')
	if idx == -1 {
		failErr(fmt.Errorf(`expected qualified type as "pkg/path.MyType"`))
	}
	sourceTypePackage := sourceType[0:idx]
	sourceTypeName := sourceType[idx+1:]
	return sourceTypePackage, sourceTypeName
}

func failErr(err error) {
	if err != nil {
		fmt.Printf("%s", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
