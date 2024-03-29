package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const FilePath = "settings.json"

type Container struct {
	LastProject    *LastProject `json:"lastProject"`
	GoRoot         string       `json:"goroot"`
	AndroidNdkHome string       `json:"androidNdkHome"`
	AndroidHome    string       `json:"androidHome"`
}

var Current *Container

func Load() *Container {
	Current = Default()

	file, err := ioutil.ReadFile(FilePath)
	if err != nil {
		return Default()
	}

	c := &Container{}
	err = json.Unmarshal(file, c)
	if err != nil {
		return Default()
	}

	Current = c
	return c
}

func Save(container *Container) (err error) {
	var data []byte
	data, err = json.Marshal(container)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(FilePath, data, os.FileMode(0777))
	return
}

func Default() *Container {
	return &Container{}
}

func InputGoRoot() (err error) {
	fmt.Print("Please enter a valid GOROOT path:")
	fmt.Scanln(&Current.GoRoot)
	_, err = os.Stat(Current.GoRoot)
	if err != nil {
		return
	}
	err = Save(Current)
	return
}