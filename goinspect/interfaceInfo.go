package goinspect

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type InterfaceInfo struct {
	DeclarationInfo
	Embeddings []*InterfaceEmbeddingInfo
	Methods    []*FuncInfo
}

func (i *InterfaceInfo) CheckImplements(structInfo *StructInfo) bool {
	implementedCount := 0
	allMethods := i.GetAllMethods()

	for _, interfaceMethod := range allMethods {
		if interfaceMethod.Name == MissingMethodsName {
			continue
		}

		for _, structMethod := range structInfo.GetAllMethods() {
			if interfaceMethod.Matches(structMethod) {
				implementedCount++
				break
			}
		}
	}

	return implementedCount == len(allMethods)
}

func (i *InterfaceInfo) GetAllMethods() []*FuncInfo {
	methodsMap := make(map[string]*FuncInfo)

	for _, e := range i.Embeddings {
		if e.InterfaceInfo == nil {
			methodsMap[MissingMethodsName] = NewMethodInfo(MissingMethodsName, nil)
			continue
		}

		for _, m := range e.InterfaceInfo.GetAllMethods() {
			methodsMap[m.Name] = m
		}
	}

	for _, m := range i.Methods {
		methodsMap[m.Name] = m
	}

	methods := make([]*FuncInfo, len(methodsMap))
	j := 0
	for _, m := range methodsMap {
		methods[j] = m
		j++
	}
	return methods
}

func (i *InterfaceInfo) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("interface %s.%s\n", i.Package, i.Name))

	for _, e := range i.Embeddings {
		sb.WriteString(fmt.Sprintf("\tembeds %s\n", e.TypeName))
	}

	for _, m := range i.GetAllMethods() {
		var embedded = ""
		for _, e := range i.Embeddings {
			if e.InterfaceInfo == nil {
				continue
			}

			for _, em := range e.InterfaceInfo.GetAllMethods() {
				if m.Name == em.Name {
					embedded = e.TypeName
					break
				}
			}
		}

		if embedded == "" {
			sb.WriteString(fmt.Sprintf("\t%v\n", m))
		} else {
			sb.WriteString(fmt.Sprintf("\t%v : %s\n", m, embedded))
		}
	}

	return sb.String()
}

type InterfaceEmbeddingInfo struct {
	FieldInfo
	InterfaceInfo *InterfaceInfo
}

func tryParseInterface(node ast.Node, packageName string) (i *InterfaceInfo, ok bool) {
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
				DeclarationInfo: DeclarationInfo{
					Name:    m.Names[0].Name,
					Package: packageName,
				},
				Receiver:   nil,
				Parameters: parseFuncFieldList(m.Type.(*ast.FuncType).Params),
				Returns:    parseFuncFieldList(m.Type.(*ast.FuncType).Results),
			})
		}
	}

	i = &InterfaceInfo{
		DeclarationInfo: DeclarationInfo{
			Name:    typeSpec.Name.Name,
			Package: packageName,
		},
		Embeddings: embeddings,
		Methods:    methods,
	}
	ok = true

	return
}