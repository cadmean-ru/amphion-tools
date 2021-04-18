package analysis

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func GetStructList(codePath string) ([]*StructInfo, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, codePath, nil, 0)
	if err != nil {
		return nil, err
	}

	structList := make([]*StructInfo, 0, 10)

	for _, pack := range packages{
		//fmt.Printf("Package name: %s\n", name)
		ast.Inspect(pack, func(node ast.Node) bool {
			if structInfo, ok := tryParseStruct(node); ok {
				structList = append(structList, structInfo)
			}
			//else if decl, ok := node.(*ast.FuncDecl); ok {
			//	fmt.Printf("Func: %s\n", decl.Name.Name)
			//	fmt.Println()
			//}
			return true
		})
	}

	return structList, nil
}

func tryParseStruct(node ast.Node) (s *StructInfo, ok bool) {
	var decl *ast.GenDecl
	decl, ok = node.(*ast.GenDecl);
	if !ok {
		return
	}

	if decl.Tok != token.TYPE || len(decl.Specs) != 1 {
		ok = false
		return
	}

	typeSpec := decl.Specs[0].(*ast.TypeSpec)

	_, ok = typeSpec.Type.(*ast.StructType)
	if !ok {
		return
	}

	s = &StructInfo{
		Name:    typeSpec.Name.String(),
		Methods: []*FuncInfo{},
	}

	return
}