package resinspect

import (
	"amphion-tools/project"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Inspector struct {

}

func (i *Inspector) FindResources(projPath string, config *project.Config) []*ResInfo {
	resPath := filepath.Join(projPath, "res")

	resources := make([]*ResInfo, 0, 10)

	_ = filepath.Walk(resPath, func(path string, info os.FileInfo, err error) error {
		if info.Name()[0] == '.' || info.IsDir() {
			return nil
		}

		sPath := stripResPath(config.Name, path)
		n := resName(sPath)

		if !validVarName.MatchString(n) {
			return nil
		}

		resources = append(resources, newResInfo(n, sPath))

		return nil
	})

	return resources
}

func NewInspector() *Inspector {
	return &Inspector{}
}

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
