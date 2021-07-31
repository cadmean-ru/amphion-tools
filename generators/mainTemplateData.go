package generators

import (
	"amphion-tools/project"
	"amphion-tools/resinspect"
)

type MainTemplateData struct {
	Frontend  string
	Resources []string
	ResNames  []string
}

func MakeMainTemplateData(runConfig *project.RunConfig, resources []*resinspect.ResInfo) *MainTemplateData {
	data := &MainTemplateData{
		Frontend:  runConfig.Frontend,
	}

	data.Resources = make([]string, len(resources))
	data.ResNames = make([]string, len(resources))

	for i, ri := range resources {
		data.Resources[i] = ri.Path
		data.ResNames[i] = ri.Name
	}

	return data
}