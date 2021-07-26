package analysis

import (
	"fmt"
	"strings"
)

type FuncInfo struct {
	Name       string
	Receiver   *FieldInfo
	Parameters []*FieldInfo
	Returns    []*FieldInfo
}

func (f FuncInfo) String() string {
	recv := ""
	if f.Receiver != nil {
		recv = "(" + f.Receiver.String() + ") "
	}

	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}

	returns := make([]string, len(f.Returns))
	for i, r := range f.Returns {
		returns[i] = r.String()
	}

	return fmt.Sprintf("func %s%s(%s) (%s)", recv, f.Name, strings.Join(params, ", "), strings.Join(returns, ", "))
}

//func (f FuncInfo) Matches(other FuncInfo) bool {
//	return other.Name == f.Name
//}

type FieldInfo struct {
	Name        string
	TypeName    string
	PassageKind PassageKind
}

func (p FieldInfo) String() string {
	star := ""
	if p.PassageKind == ByPointer {
		star = "*"
	}
	return fmt.Sprintf("%s %s%s", p.Name, star, p.TypeName)
}
