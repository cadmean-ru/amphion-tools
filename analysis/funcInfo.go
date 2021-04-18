package analysis

type FuncInfo struct {
	Name       string
	Receiver   *FuncParameter
	Parameters []*FuncParameter
	ReturnType TypeName
}

type FuncParameter struct {
	Name        string
	TypeName    TypeName
	PassageKind PassageKind
}
