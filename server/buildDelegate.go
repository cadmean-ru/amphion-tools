package server

import (
	"amphion-tools/generators"
	"amphion-tools/project"
)

type BuildDelegate interface {
	Build(ctx *BuildDelegateContext) error
}

type BuildDelegateContext struct {
	projPath         string
	buildPath        string
	proj             *project.Config
	runConfig         *project.RunConfig
	mainTemplateData *generators.MainTemplateData
}

func NewBuildDelegateForFrontend(frontend string) BuildDelegate {
	switch frontend {
	case "android":
		return &AndroidBuildDelegate{}
	case "ios":
		return &IosBuildDelegate{}
	case "pc":
		return &PcBuildDelegate{}
	case "web":
		return &WebBuildDelegate{}
	default:
		return nil
	}
}