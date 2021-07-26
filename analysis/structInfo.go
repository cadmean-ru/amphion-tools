package analysis

type PassageKind byte

const (
	ByPointer PassageKind = iota
	ByValue
)

//type TypeName string

type StructInfo struct {
	Name       string
	Methods    []*FuncInfo
	Embeddings []*FieldInfo
	Fields     []*FieldInfo
}
