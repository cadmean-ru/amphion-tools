package server

import "amphion-tools/project"

type BuildDelegate interface {
	Build(projPath string, proj *project.Config, runConfig *project.RunConfig) error
}
