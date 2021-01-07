package project

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

type RunConfig struct {
	Name     string `yaml:"name"`
	Frontend string `yaml:"frontend"`
	Debug    string `yaml:"debug"`
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