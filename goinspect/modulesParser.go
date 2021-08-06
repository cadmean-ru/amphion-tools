package goinspect

import (
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

	for _, line := range strings.Split(modStr, "\n") {
		if strings.HasPrefix(line, "module") {
			info.ModuleName = strings.ReplaceAll(strings.Split(line, " ")[1], "\r", "")
			break
		}
	}

	return &info, nil
}
