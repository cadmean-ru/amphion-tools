package analysis

import (
	"fmt"
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

	for _, pack := range packages {
		//fmt.Printf("Package name: %s\n", name)
		ast.Inspect(pack, func(node ast.Node) bool {
			if structInfo, ok := tryParseStruct(node); ok {
				structList = append(structList, structInfo)
			}
			return true
		})
	}

	return structList, nil
}

func GetFunctionsList(codePath string) ([]*FuncInfo, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, codePath, nil, 0)
	if err != nil {
		return nil, err
	}

	funcList := make([]*FuncInfo, 0, 10)

	for _, pack := range packages {
		//fmt.Printf("Package name: %s\n", name)
		ast.Inspect(pack, func(node ast.Node) bool {
			if funcInfo, ok := tryParseFunc(node); ok {
				funcList = append(funcList, funcInfo)
			}
			return true
		})
	}

	return funcList, nil
}

func GetMethods(structList []*StructInfo, funcList []*FuncInfo) {
	for _, f := range funcList {
		if f.Receiver == nil {
			continue
		}

		for _, s := range structList {
			if f.Receiver.TypeName == s.Name {
				s.Methods = append(s.Methods, f)
			}
		}
	}
}

func tryParseStruct(node ast.Node) (s *StructInfo, ok bool) {
	var decl *ast.GenDecl
	decl, ok = node.(*ast.GenDecl)
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

	structType := typeSpec.Type.(*ast.StructType)

	embeddings := make([]*FieldInfo, 0)
	fields := make([]*FieldInfo, 0)

	if structType.Fields != nil {
		fieldInfos := make([]*FieldInfo, 0, len(structType.Fields.List))
		for _, f := range structType.Fields.List {
			fieldInfos = append(fieldInfos, parseField(f)...)
		}

		for _, f := range fieldInfos {
			fmt.Println(f)

			if f.Name == "_" {
				embeddings = append(embeddings, f)
			} else {
				fields = append(fields, f)
			}
		}
	}

	s = &StructInfo{
		Name:       typeSpec.Name.String(),
		Methods:    []*FuncInfo{},
		Embeddings: embeddings,
		Fields:     fields,
	}

	return
}

func tryParseFunc(node ast.Node) (f *FuncInfo, ok bool) {
	var decl *ast.FuncDecl
	decl, ok = node.(*ast.FuncDecl)
	if !ok {
		return
	}

	//fmt.Printf("Found function: %+v\n", decl)
	//if decl.Recv != nil {
	//	fmt.Printf("Function has receiver %+v\n", decl.Recv.List[0].Type)
	//}

	//if decl.Name.Name == "IsSelected" {
	//	fmt.Println("HERE")
	//}

	var receiver *FieldInfo
	if decl.Recv != nil {
		receiver = parseFuncFieldList(decl.Recv)[0]
	}

	paramInfos := parseFuncFieldList(decl.Type.Params)
	returnInfos := parseFuncFieldList(decl.Type.Results)

	f = &FuncInfo{
		Name:       decl.Name.Name,
		Receiver:   receiver,
		Parameters: paramInfos,
		Returns:    returnInfos,
	}

	return
}

func parseFuncFieldList(list *ast.FieldList) []*FieldInfo {
	if list == nil || list.List == nil {
		return []*FieldInfo{}
	}

	params := make([]*FieldInfo, 0, len(list.List))
	for _, p := range list.List {
		params = append(params, parseField(p)...)
	}
	return params
}

func parseField(field *ast.Field) []*FieldInfo {
	fieldTypeName, kind := parseFieldType(field.Type)

	if field.Names == nil || len(field.Names) == 0 {
		return []*FieldInfo{
			{
				Name:        "_",
				TypeName:    fieldTypeName,
				PassageKind: kind,
			},
		}
	} else {
		infos := make([]*FieldInfo, len(field.Names))
		for i, n := range field.Names {
			infos[i] = &FieldInfo{
				Name:        n.Name,
				TypeName:    fieldTypeName,
				PassageKind: kind,
			}
		}
		return infos
	}
}

func parseFieldType(param ast.Expr) (name string, pointer PassageKind) {
	//fmt.Printf("%T %+v\n", param, param)

	switch param.(type) {
	case *ast.Ident:
		return param.(*ast.Ident).Name, ByValue
	case *ast.SelectorExpr:
		return param.(*ast.SelectorExpr).Sel.Name, ByValue
	case *ast.StarExpr:
		parameterType, _ := parseFieldType(param.(*ast.StarExpr).X)
		return parameterType, ByPointer
	case *ast.MapType:
		return "map", ByPointer
	case *ast.ArrayType:
		return "array", ByValue
	case *ast.SliceExpr:
		return "slice", ByPointer
	default:
		return "", ByValue
	}
}
