package analysis

type InterfaceInfo struct {
	Name    string
	Methods []*FuncInfo
}

//func (i *InterfaceInfo) CheckImplements(structInfo *StructInfo) bool {
//	implementedCount := 0
//
//	for _, interfaceMethod := range i.Methods {
//		for _, structMethod := range structInfo.Methods {
//
//		}
//	}
//}
