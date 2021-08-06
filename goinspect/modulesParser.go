package goinspect

import (
	"amphion-tools/utils"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type GoModInfo struct {
	ModuleName string
}

func ParseGoMod(projectPath string) (*GoModInfo, error) {
	path := filepath.Join(projectPath, "go.mod")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	modStr := string(data)

	info := GoModInfo{}

	for _, line := range strings.Split(modStr, utils.NewLineString()) {
		if strings.HasPrefix(line, "module") {
			info.ModuleName = strings.Split(line, " ")[1]
			break
		}
	}

	return &info, nil
}
