package generators

import (
	"regexp"
	"strings"
)

var notW = regexp.MustCompile("\\W")
var validVarName = regexp.MustCompile("^[_A-z]+[_A-z$0-9]*$")

func resName(path string) string {
	var s = path
	s = s[:strings.LastIndex(s, ".")]
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, " ", "")
	s = notW.ReplaceAllString(s, "")
	s = strings.Title(s)
	return s
}

func stripResPath(projName, path string) string {
	flag := projName + "/res"
	p := strings.ReplaceAll(path, "\\", "/")
	return p[strings.Index(p, flag)+len(flag)+1:]
}
