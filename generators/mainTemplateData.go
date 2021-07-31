package generators

import (
	"amphion-tools/project"
	"amphion-tools/resinspect"
)

type MainTemplateData struct {
	Frontend  string
	Resources []string
	ResNames  []string
	CompData  *CompFileTemplateData
}

func MakeMainTemplateData(runConfig *project.RunConfig, resources []*resinspect.ResInfo, components *CompFileTemplateData) *MainTemplateData {
	data := &MainTemplateData{
		Frontend:  runConfig.Frontend,
		Resources: make([]string, len(resources)),
		ResNames:  make([]string, len(resources)),
		CompData:  components,
	}

	for i, ri := range resources {
		data.Resources[i] = ri.Path
		data.ResNames[i] = ri.Name
	}

	return data
}
