package goinspect

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type StructInfo struct {
	DeclarationInfo
	Methods    []*FuncInfo
	Embeddings []*StructEmbeddingInfo
	Fields     []*FieldInfo
}

func (s *StructInfo) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("struct %s.%s\n", s.Package, s.Name))

	for _, e := range s.Embeddings {
		sb.WriteString(fmt.Sprintf("\tembeds %s\n", e.TypeName))
	}

	for _, f := range s.Fields {
		sb.WriteString(fmt.Sprintf("\tfield %v\n", f))
	}

	for _, m := range s.GetAllMethods() {
		sb.WriteString(fmt.Sprintf("\t%v\n", m))
	}

	return sb.String()
}

func (s *StructInfo) GetAllMethods() []*FuncInfo {
	methodsMap := make(map[string]*FuncInfo)

	for _, e := range s.Embeddings {
		if e.StructInfo == nil {
			methodsMap[MissingMethodsName] = NewMethodInfo(MissingMethodsName, nil)
			continue
		}

		for _, m := range e.StructInfo.GetAllMethods() {
			methodsMap[m.Name] = m
		}
	}

	var knownReceiver *FieldInfo
	for _, m := range s.Methods {
		methodsMap[m.Name] = m

		if knownReceiver == nil {
			knownReceiver = m.Receiver
		}
	}

	if m, ok := methodsMap[MissingMethodsName]; ok {
		m.Receiver = knownReceiver
	}

	methods := make([]*FuncInfo, len(methodsMap))
	i := 0
	for _, m := range methodsMap {
		methods[i] = m
		i++
	}
	return methods
}

type StructEmbeddingInfo struct {
	FieldInfo
	StructInfo *StructInfo
}

func tryParseStruct(node ast.Node, packageName string) (s *StructInfo, ok bool) {
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
		DeclarationInfo: DeclarationInfo{
			Name:    typeSpec.Name.String(),
			Package: packageName,
		},
		Methods:    []*FuncInfo{},
		Embeddings: embeddings,
		Fields:     fields,
	}
	ok = true

	return
}