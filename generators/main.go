package generators

import (
	"amphion-tools/project"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const mainFileTemplate = `
package main

// This file is automatically generated
// DO NOT EDIT THIS FILE MANUALLY!!!

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/frontend/{{ .Frontend }}"
	"runtime"
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
	cm.RegisterComponentType(&builtin.ShapeView{})
	cm.RegisterComponentType(&builtin.CircleBoundary{})
	cm.RegisterComponentType(&builtin.OnClickListener{})
	cm.RegisterComponentType(&builtin.TextView{})
	cm.RegisterComponentType(&builtin.RectBoundary{})
	cm.RegisterComponentType(&builtin.TriangleBoundary{})
	cm.RegisterComponentType(&builtin.BezierView{})
	cm.RegisterComponentType(&builtin.DropdownView{})
	cm.RegisterComponentType(&builtin.ImageView{})
	cm.RegisterComponentType(&builtin.MouseMover{})
	cm.RegisterComponentType(&builtin.BuilderComponent{})
	cm.RegisterComponentType(&builtin.GridLayout{})
	cm.RegisterComponentType(&builtin.NativeInputView{})
	cm.RegisterComponentType(&builtin.EventListener{})

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

const resFileTemplate = `
package res

// This file is automatically generated
// DO NOT EDIT THIS FILE MANUALLY!!!

import "github.com/cadmean-ru/amphion/common/a"

{{ range $i, $res := .ResNames }}
const {{ $res }} a.ResId = {{ $i }}
{{ end }}
`

func Main(projPath string, config *project.Config, runConfig *project.RunConfig) (err error) {
	resPath := filepath.Join(projPath, "res")
	codePath := filepath.Join(projPath, config.Name)

	data := struct {
		Frontend  string
		Resources []string
		ResNames  []string
	}{
		Frontend:  runConfig.Frontend,
		Resources: make([]string, 0),
		ResNames:  make([]string, 0),
	}

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

	mainTmpl := template.Must(template.New("main").Parse(mainFileTemplate))

	mainFile, err := os.Create(filepath.Join(codePath, "main.gen.go"))
	if err != nil {
		return
	}
	defer mainFile.Close()

	err = mainTmpl.Execute(mainFile, data)
	if err != nil {
		return
	}

	resTmpl := template.Must(template.New("main").Parse(resFileTemplate))

	_ = os.MkdirAll(filepath.Join(codePath, "generated", "res"), os.FileMode(0777))

	resFile, err := os.Create(filepath.Join(codePath, "generated", "res", "res.gen.go"))
	if err != nil {
		return
	}
	defer resFile.Close()

	err = resTmpl.Execute(resFile, data)

	return
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
