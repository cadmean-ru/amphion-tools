package project

import "runtime"

type App struct {
	Name          string `yaml:"name"`
	Author        string `yaml:"author"`
	CompanyDomain string `yaml:"companyDomain"`
	PublicUrl     string `yaml:"publicUrl"`
	Frontend      string `yaml:"frontend"`
	Debug         bool   `yaml:"debug"`
	MainScene     string `yaml:"mainScene"`
	HostOS        string `yaml:"hostOS"`
}

func NewAppFromConfig(config *Config, runConfig *RunConfig) *App {
	return &App{
		Name:          config.Name,
		Author:        config.Author,
		CompanyDomain: config.CompanyDomain,
		PublicUrl:     runConfig.Url,
		Frontend:      runConfig.Frontend,
		Debug:         runConfig.Debug,
		MainScene:     config.MainScene,
		HostOS:        runtime.GOOS,
	}
}