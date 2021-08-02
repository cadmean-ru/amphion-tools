package server

import (
	"amphion-tools/generators"
	"amphion-tools/gotools"
	"os"
	"path/filepath"
)

type WebBuildDelegate struct {

}

func (w *WebBuildDelegate) Build(ctx *BuildDelegateContext) (err error) {
	srcPath := filepath.Join(ctx.projPath, ctx.proj.Name)

	//1. Generate code
	err = generators.Main(ctx.mainTemplateData, ctx.projPath, ctx.proj, ctx.runConfig)
	if err != nil {
		return
	}

	//2. Run go build
	var dstPath = filepath.Join(ctx.buildPath, "web")
	var dstFileName = executableName(ctx.proj, ctx.runConfig)
	var goos = "js"
	var goarch = "wasm"

	_ = os.Mkdir(dstPath, os.FileMode(0777))

	err = gotools.Build(srcPath, dstPath, dstFileName, goos, goarch)

	return
}

