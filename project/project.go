package project

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Name           string       `yaml:"name"`
	Author         string       `yaml:"author"`
	CompanyDomain  string       `yaml:"companyDomain"`
	GitUrl         string       `yaml:"gitUrl"`
	PublicUrl      string       `yaml:"publicUrl"`
	Configurations []*RunConfig `yaml:"configurations"`
	MainScene      string       `yaml:"mainScene"`
}

func (c *Config) GetRunConfig(name string) *RunConfig {
	for _, rc := range c.Configurations {
		if rc.Name == name {
			return rc
		}
	}

	return nil
}

func (c *Config) GetCodePath(path string) string {
	return filepath.Join(path, c.Name)
}

type RunConfig struct {
	Name     string `yaml:"name"`
	Frontend string `yaml:"frontend"`
	Debug    bool   `yaml:"debug"`
	Url      string `yaml:"url"`
}

func ParseConfig(configFilePath string) (config *Config, err error) {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}

	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return
	}

	return &c, nil
}

func FindProjectConfig(projectPath string) (config *Config, err error) {
	_, name := filepath.Split(projectPath)

	configFilePath := filepath.Join(projectPath, name+".config.yaml")

	config, err = ParseConfig(configFilePath)
	return
}
