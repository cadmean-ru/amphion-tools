package goinspect

import (
	"fmt"
	"strings"
)

type StructInfo struct {
	Name       string
	Methods    []*FuncInfo
	Embeddings []*StructEmbeddingInfo
	Fields     []*FieldInfo
}

func (s *StructInfo) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("struct %s\n", s.Name))

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
