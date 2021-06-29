package generators

import "amphion-tools/project"

type MainTemplateData struct {
	Frontend  string
	Resources []string
	ResNames  []string
}

func MakeMainTemplateData(projPath string, config *project.Config, runConfig *project.RunConfig) *MainTemplateData {
	data := &MainTemplateData{
		Frontend:  runConfig.Frontend,
	}
	findResources(data, projPath, config)
	return data
}