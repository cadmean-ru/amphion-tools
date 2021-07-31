package goinspect

import "unicode"

type DeclarationInfo struct {
	Name    string
	Package string
}

func (d DeclarationInfo) IsExported() bool {
	return IsNameExported(d.Name)
}

func IsNameExported(name string) bool {
	return unicode.IsUpper([]rune(name)[0])
}