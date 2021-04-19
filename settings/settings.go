package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const FilePath = "settings.json"

type Container struct {
	LastProject *LastProject `json:"lastProject"`
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