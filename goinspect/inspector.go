package goinspect

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Inspector struct {
	functions []*FuncInfo
	structs   []*StructInfo
	interfaces []*InterfaceInfo
}

func (i *Inspector) GetFunctions() []*FuncInfo {
	return i.functions
}

func (i *Inspector) GetStructs() []*StructInfo {
	return i.structs
}

func (i *Inspector) GetInterfaces() []*InterfaceInfo {
	return i.interfaces
}

func (i *Inspector) GetInterface(name string) *InterfaceInfo {
	for _, ii := range i.interfaces {
		if ii.Name == name {
			return ii
		}
	}

	return nil
}

func (i *Inspector) InspectSemantics(path string) error {
	structs, funcs, interfaces, err := i.findDeclarations(path)
	if err != nil {
		return err
	}

	i.findMethods(structs, funcs)

	i.findTopLevelFunctions(funcs)

	i.structs = append(i.structs, structs...)
	i.interfaces = append(i.interfaces, interfaces...)

	i.inspectStructEmbeddings(structs)
	i.inspectInterfaceEmbeddings(interfaces)

	return nil
}

func (i *Inspector) findTopLevelFunctions(funcs []*FuncInfo) {
	for _, f := range funcs {
		if !f.IsMethod() {
			i.functions = append(i.functions, f)
		}
	}
}

func (i *Inspector) inspectStructEmbeddings(structs []*StructInfo) {
	for _, s := range structs {
		for _, e := range s.Embeddings {
			for _, ss := range i.structs {
				if e.TypeName == ss.Name {
					e.StructInfo = ss
					break
				}
			}
		}
	}
}

func (i *Inspector) inspectInterfaceEmbeddings(structs []*InterfaceInfo) {
	for _, ii := range structs {
		for _, e := range ii.Embeddings {
			for _, iii := range i.interfaces {
				if e.TypeName == iii.Name {
					e.InterfaceInfo = iii
					break
				}
			}
		}
	}
}

func (i *Inspector) findMethods(structList []*StructInfo, funcList []*FuncInfo) {
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

func (i *Inspector) findDeclarations(codePath string) ([]*StructInfo, []*FuncInfo, []*InterfaceInfo, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, codePath, nil, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	structList := make([]*StructInfo, 0, 10)
	funcList := make([]*FuncInfo, 0, 10)
	interfaceList := make([]*InterfaceInfo, 0, 10)

	for _, pack := range packages {
		//fmt.Printf("Package name: %s\n", pack.Name)

		ast.Inspect(pack, func(node ast.Node) bool {
			if structInfo, ok := tryParseStruct(node); ok {
				structList = append(structList, structInfo)
			} else if interfaceInfo, ok := tryParseInterface(node); ok {
				interfaceList = append(interfaceList, interfaceInfo)
			} else if funcInfo, ok := tryParseFunc(node); ok {
				funcList = append(funcList, funcInfo)
			}
			return true
		})
	}

	return structList, funcList, interfaceList, nil
}

func NewInspector() *Inspector {
	return &Inspector{
		functions: []*FuncInfo{},
		structs: []*StructInfo{},
		interfaces: []*InterfaceInfo{},
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

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return
	}

	embeddings := make([]*StructEmbeddingInfo, 0)
	fields := make([]*FieldInfo, 0)

	if structType.Fields != nil {
		fieldInfos := make([]*FieldInfo, 0, len(structType.Fields.List))
		for _, f := range structType.Fields.List {
			fieldInfos = append(fieldInfos, parseField(f)...)
		}

		for _, f := range fieldInfos {
			//fmt.Println(f)

			if f.Name == "_" {
				embeddings = append(embeddings, &StructEmbeddingInfo{
					FieldInfo:  *f,
					StructInfo: nil,
				})
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
	ok = true

	return
}

func tryParseInterface(node ast.Node) (i *InterfaceInfo, ok bool) {
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

	interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return
	}

	embeddings := make([]*InterfaceEmbeddingInfo, 0)
	methods := make([]*FuncInfo, 0, 2)

	for _, m := range interfaceType.Methods.List {
		if m.Names == nil || len(m.Names) == 0 {
			embeddings = append(embeddings, &InterfaceEmbeddingInfo{
				FieldInfo:     *(parseField(m)[0]),
				InterfaceInfo: nil,
			})
		} else {
			methods = append(methods, &FuncInfo{
				Name:       m.Names[0].Name,
				Receiver:   nil,
				Parameters: parseFuncFieldList(m.Type.(*ast.FuncType).Params),
				Returns:    parseFuncFieldList(m.Type.(*ast.FuncType).Results),
			})
		}
	}

	i = &InterfaceInfo{
		Name:       typeSpec.Name.Name,
		Embeddings: embeddings,
		Methods:    methods,
	}
	ok = true

	return
}

func tryParseFunc(node ast.Node) (f *FuncInfo, ok bool) {
	var decl *ast.FuncDecl
	decl, ok = node.(*ast.FuncDecl)
	if !ok {
		return
	}

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