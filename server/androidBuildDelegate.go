package server

import (
	"amphion-tools/generators"
	"amphion-tools/gotools"
	"amphion-tools/settings"
	"fmt"
	"path/filepath"
)

type AndroidBuildDelegate struct {

}

func (a *AndroidBuildDelegate) Build(ctx *BuildDelegateContext) (err error) {
	srcPath := filepath.Join(ctx.projPath, ctx.proj.Name)

	//1. Generate code
	err = generators.AndroidMain(ctx.mainTemplateData, ctx.projPath, ctx.proj, ctx.runConfig)
	if err != nil {
		return
	}
	err = generators.Android(ctx.mainTemplateData, ctx.projPath, ctx.proj, ctx.runConfig)
	if err != nil {
		return
	}

	//2. Run gomobile bind
	var dstPath = filepath.Join(ctx.buildPath, ctx.proj.Name+".android.aar")

	err = gotools.MobileBind(gotools.MobileBindParams{
		Target:             gotools.BindTargetAndroid,
		SrcPath:            srcPath,
		DstFilePath:        dstPath,
		JavaPackageName:    "ru.cadmean.amphion.android",
		AndroidHomePath:    settings.Current.AndroidHome,
		AndroidNdkHomePath: settings.Current.AndroidNdkHome,
		PackagesToBind:     []string{
			ctx.proj.Name+"/generated/droidCli",
			"github.com/cadmean-ru/amphion/frontend/cli",
			"github.com/cadmean-ru/amphion/common/atext",
			"github.com/cadmean-ru/amphion/common/dispatch",
		},
	})

	if err == nil {
		fmt.Printf("Android library file was successfully created: %s\n", dstPath)
	}

	return
}

