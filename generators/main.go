package generators

import (
	"amphion-tools/project"
	"os"
	"path/filepath"
	"text/template"
)

const mainFileTemplate = `
package main

// This file is automatically generated
// DO NOT EDIT THIS FILE MANUALLY!!!

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend/{{ .Frontend }}"
{{ if (eq .Frontend "pc") }}
	"runtime"
{{ end }}
{{ range $imp := .CompData.Imports }}
	"{{ $imp }}"
{{ end }}
)

{{ if (eq .Frontend "pc") }}
func init() {
	runtime.LockOSThread()
}
{{ end }}


func runApp() {
	front := {{ .Frontend }}.NewFrontend()
	front.Init()

	e := engine.Initialize(front)

	cm := e.GetComponentsManager()

{{ range $comp := .CompData.Components }}
	cm.RegisterComponentType(&{{ $comp.LastPackage }}.{{ $comp.Name }}{})
{{ end }}

	registerComponents(cm)
	
	{{ if (gt (len .Resources) 0) }}
	rm := e.GetResourceManager()
		{{ range $res := .Resources }}
	rm.RegisterResource("{{ $res }}")
		{{ end }}
	{{ end }}

	go func() {
		e.Start()
		e.LoadApp()

		{{ if (ne .Frontend "pc") }}
		e.WaitForStop()
		{{ end }}
	}()

	front.Run()
}
`

func Main(data *MainTemplateData, projPath string, config *project.Config, runConfig *project.RunConfig) error {
	return generateMainFile(data, mainFileTemplate, projPath, config, runConfig)
}

func generateMainFile(data *MainTemplateData, tmpl string, projPath string, config *project.Config, runConfig *project.RunConfig) (err error) {
	codePath := filepath.Join(projPath, config.Name)
	mainTmpl := template.Must(template.New("main").Parse(tmpl))

	mainFile, err := os.Create(filepath.Join(codePath, "main.gen.go"))
	if err != nil {
		return
	}
	defer mainFile.Close()

	err = mainTmpl.Execute(mainFile, *data)
	return
}
