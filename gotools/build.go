package gotools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Build(srcPath, dstPath, dstFileName, goos, goarch string) (err error) {
	outFilePath := filepath.Join(dstPath, dstFileName)

	build := exec.Command("go", "build", "-o", outFilePath)
	build.Dir = srcPath
	build.Env = os.Environ()
	build.Env = append(build.Env, "GOOS="+goos)
	build.Env = append(build.Env, "GOARCH="+goarch)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	err = build.Run()
	if err != nil {
		fmt.Println(err)
	}

	return
}
