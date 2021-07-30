package goinspect

import (
	"fmt"
	"go/ast"
)

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