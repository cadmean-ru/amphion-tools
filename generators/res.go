package generators

import (
	"amphion-tools/project"
	"os"
	"path/filepath"
	"text/template"
)

const resFileTemplate = `
package res

// This file is automatically generated
// DO NOT EDIT THIS FILE MANUALLY!!!

import "github.com/cadmean-ru/amphion/common/a"

{{ range $i, $res := .ResNames }}
const {{ $res }} a.ResId = {{ $i }}
{{ end }}
`

func Res(data *MainTemplateData, projPath string, config *project.Config) (err error) {
	codePath := filepath.Join(projPath, config.Name)

	resTmpl := template.Must(template.New("main").Parse(resFileTemplate))

	_ = os.MkdirAll(filepath.Join(codePath, "generated", "res"), os.FileMode(0777))

	resFile, err := os.Create(filepath.Join(codePath, "generated", "res", "res.gen.go"))
	if err != nil {
		return err
	}
	defer resFile.Close()

	err = resTmpl.Execute(resFile, *data)
	return
}

func findResources(data *MainTemplateData, projPath string, config *project.Config) {
	resPath := filepath.Join(projPath, "res")

	_ = filepath.Walk(resPath, func(path string, info os.FileInfo, err error) error {
		if info.Name()[0] == '.' || info.IsDir() {
			return nil
		}

		sPath := stripResPath(config.Name, path)
		n := resName(sPath)

		if !validVarName.MatchString(n) {
			return nil
		}

		data.Resources = append(data.Resources, sPath)
		data.ResNames = append(data.ResNames, n)

		return nil
	})
}