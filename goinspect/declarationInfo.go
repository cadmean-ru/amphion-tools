package goinspect

import "unicode"

type DeclarationInfo struct {
	Name    string
	Package string
}

func IsNameExported(name string) bool {
	return unicode.IsUpper([]rune(name)[0])
}