package server

import (
	"amphion-tools/generators"
	"amphion-tools/gotools"
	"fmt"
	"path/filepath"
)

type IosBuildDelegate struct {

}

func (i *IosBuildDelegate) Build(ctx *BuildDelegateContext) (err error) {
	srcPath := filepath.Join(ctx.projPath, ctx.proj.Name)

	//1. Generate code
	err = generators.IosMain(ctx.mainTemplateData, ctx.projPath, ctx.proj, ctx.runConfig)
	if err != nil {
		return
	}
	err = generators.Ios(ctx.mainTemplateData, ctx.projPath, ctx.proj, ctx.runConfig)
	if err != nil {
		return
	}

	//2. Run gomobile bind
	var dstPath = filepath.Join(ctx.buildPath, "Amphion.framework")

	err = gotools.MobileBind(gotools.MobileBindParams{
		Target:             gotools.BindTargetIos,
		SrcPath:            srcPath,
		DstFilePath:        dstPath,
		PackagesToBind:     []string{
			ctx.proj.Name+"/generated/iosCli",
			"github.com/cadmean-ru/amphion/frontend/cli",
			"github.com/cadmean-ru/amphion/common/atext",
			"github.com/cadmean-ru/amphion/common/dispatch",
		},
	})

	if err == nil {
		fmt.Printf("iOS framework was successfully created: %s\n", dstPath)
	}

	return
}

