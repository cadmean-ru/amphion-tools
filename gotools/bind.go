package gotools

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const (
	BindTargetAndroid = "android"
	BindTargetIos     = "android"
)

type MobileBindParams struct {
	Target             string
	SrcPath            string
	DstFilePath        string
	JavaPackageName    string
	AndroidHomePath    string
	AndroidNdkHomePath string
	PackagesToBind     []string
}

func MobileBind(params MobileBindParams) (err error) {
	//if target == BindTargetAndroid {
	//	bind = exec.Command("gomobile",
	//		"bind",
	//		"-target="+target,
	//		"-o", dstFilePath,
	//		"-javapkg=ru.cadmean.amphion.android",
	//		s.proj.Name+"/generated/droidCli",
	//		"github.com/cadmean-ru/amphion/frontend/cli",
	//		"github.com/cadmean-ru/amphion/common/atext",
	//		"github.com/cadmean-ru/amphion/common/dispatch",
	//	)
	//} else if target == BindTargetIos {
	//	bind = exec.Command("gomobile",
	//		"bind",
	//		"-target="+target,
	//		"-o", dstFilePath,
	//		s.proj.Name+"/generated/iosCli",
	//		"github.com/cadmean-ru/amphion/frontend/cli",
	//		"github.com/cadmean-ru/amphion/common/atext",
	//		"github.com/cadmean-ru/amphion/common/dispatch",
	//	)
	//}

	var args []string

	if params.Target == BindTargetAndroid {
		args = []string{
			"bind",
			"-target=" + params.Target,
			"-o", params.DstFilePath,
			"-javapkg=" + params.JavaPackageName,
		}
	} else if params.Target == BindTargetIos {
		args = []string{
			"bind",
			"-target=" + params.Target,
			"-o", params.DstFilePath,
		}
	} else {
		err = errors.New(fmt.Sprintf("unknown target %s", params.Target))
		return
	}

	args = append(args, params.PackagesToBind...)

	bind := exec.Command("gomobile", args...)
	bind.Dir = params.SrcPath
	bind.Env = os.Environ()
	if params.Target == BindTargetAndroid {
		bind.Env = append(bind.Env, "ANDROID_NDK_HOME="+params.AndroidNdkHomePath)
		bind.Env = append(bind.Env, "ANDROID_HOME="+params.AndroidHomePath)
	}

	bind.Stdout = os.Stdout
	bind.Stderr = os.Stderr

	err = bind.Run()
	if err != nil {
		return
	}

	err = bind.Wait()
	return
}
