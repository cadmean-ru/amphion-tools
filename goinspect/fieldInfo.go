package goinspect

import "fmt"

type PassageKind byte

const (
	ByPointer PassageKind = iota
	ByValue
)

type FieldInfo struct {
	Name        string
	TypeName    string
	PassageKind PassageKind
}

func (f FieldInfo) Matches(other *FieldInfo) bool {
	return f.Name == other.Name && f.TypeName == other.TypeName && f.PassageKind == other.PassageKind
}

func (f FieldInfo) String() string {
	star := ""
	if f.PassageKind == ByPointer {
		star = "*"
	}
	return fmt.Sprintf("%s %s%s", f.Name, star, f.TypeName)
}
