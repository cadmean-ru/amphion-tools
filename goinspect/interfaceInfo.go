package goinspect

import (
	"fmt"
	"strings"
)

type InterfaceInfo struct {
	Name       string
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

	sb.WriteString(fmt.Sprintf("interface %s\n", i.Name))

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