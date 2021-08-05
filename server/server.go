package server

import (
	"amphion-tools/analysis"
	"amphion-tools/generators"
	"amphion-tools/goinspect"
	"amphion-tools/project"
	"amphion-tools/resinspect"
	"amphion-tools/settings"
	"amphion-tools/support"
	"amphion-tools/utils"
	"errors"
	"fmt"
	ccolor "github.com/TwinProduction/go-color"
	"gopkg.in/yaml.v2"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func StartDevelopment(projectPath, runConfigName string) (s *DevServer, err error) {
	config, err := project.FindProjectConfig(projectPath)
	if err != nil {
		return
	}

	runConfig := config.GetRunConfig(runConfigName)
	if runConfig == nil {
		return nil, fmt.Errorf("run configuration not found")
	}

	amVersion, err := getAmphionVersion(filepath.Join(projectPath, config.Name))
	if err != nil {
		return nil, err
	}

	settings.Current.LastProject = &settings.LastProject{
		Name: config.Name,
		Path: projectPath,
	}
	_ = settings.Save(settings.Current)

	s = &DevServer{
		runConfig:      runConfig,
		done:          make(chan bool, 1),
		stopped:       false,
		proj:          config,
		projPath:      projectPath,
		buildPath:     filepath.Join(projectPath, "build"),
		amVersion:     amVersion,
		goInspector:   goinspect.NewInspector(),
		resInspector:  resinspect.NewInspector(),
		buildDelegate: NewBuildDelegateForFrontend(runConfig.Frontend),
	}

	err = s.prepareInspector()
	if err != nil {
		return nil, err
	}

	return
}

func getAmphionVersion(codePath string) (string, error) {
	deps, err := analysis.GetProjectDependencies(codePath)
	if err != nil {
		return "", err
	}

	var ok bool
	var version string
	for _, dep := range deps {
		if dep.Name == "github.com/cadmean-ru/amphion" {
			if support.IsAmphionVersionSupported(dep.Version) {
				ok = true
			}
			version = dep.Version
			break
		}
	}

	if !ok {
		return "nil", errors.New("unsupported amphion version")
	}

	return version, nil
}

type DevServer struct {
	runConfig      *project.RunConfig
	done          chan bool
	stopped       bool
	proj          *project.Config
	projPath      string
	buildPath     string
	webDebug      *WebDebugServer
	amVersion     string
	goInspector   *goinspect.Inspector
	resInspector  *resinspect.Inspector
	buildDelegate BuildDelegate
}

func (s *DevServer) Stop() {
	if s.stopped {
		return
	}
	s.stopped = true
	s.done <- true

	if s.webDebug != nil {
		s.webDebug.Stop()
	}
}

func (s *DevServer) Start() {
	_ = os.Mkdir("./run", os.FileMode(0777))
	_ = os.Mkdir(s.buildPath, os.FileMode(0777))

	s.stopped = false

	if s.runConfig.Frontend == "web" {
		go utils.HttpServeDir("./run", s.runConfig.Url, s.done)

		if s.runConfig.Debug {
			s.webDebug = NewWebDebugServer("4200")
			s.webDebug.Start()
			utils.OpenBrowser("http://" + s.runConfig.Url + "?connectDebugger=4200")
		} else {
			utils.OpenBrowser("http://" + s.runConfig.Url)
		}
	}
}

func (s *DevServer) prepareInspector() error {
	if settings.Current.GoRoot == "" {
		fmt.Println("GOROOT not defined")
		err := settings.InputGoRoot()
		if err != nil {
			return err
		}
	}

	amPath := filepath.Join(settings.Current.GoRoot, "pkg", "mod", "github.com", "cadmean-ru", "amphion@"+s.amVersion)
	if !utils.Exists(amPath) {
		return errors.New(fmt.Sprintf("amphion@%s not found in GOROOT (try running 'go get')", s.amVersion))
	}

	codePath := filepath.Join(s.projPath, s.proj.Name)

	_, _ = s.goInspector.NewScope(goinspect.AmphionScope, amPath)
	_, _ = s.goInspector.NewScope(goinspect.ProjectScope, codePath)

	err := s.goInspector.InspectAmphion()
	if err != nil {
		return err
	}

	return nil
}

func (s *DevServer) handleDirectoryUpdate() {
	//_ = s.BuildProject()
}

func (s *DevServer) BuildProject() error {
	err := s.InspectCode()
	if err != nil {
		return err
	}

	fmt.Println("Building project...")

	err = s.buildExecutable()
	if err != nil {
		return err
	}

	err = s.buildApp()
	if err != nil {
		return err
	}

	return err
}

func (s *DevServer) InspectCode() error {
	fmt.Println("Running code inspection...")

	prScope := s.goInspector.GetScope("project")
	prScope.Clear()

	err := filepath.Walk(prScope.Path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		relPath := strings.TrimPrefix(path, prScope.Path)
		relPath = strings.TrimPrefix(relPath, "/")
		err = s.goInspector.InspectSemantics(prScope, relPath)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, msg := range s.goInspector.InspectComponents(prScope) {
		fmt.Println(ccolor.Ize(ccolor.Yellow, msg))
	}

	return nil
}

func (s *DevServer) buildExecutable() (err error) {
	resources := s.resInspector.FindResources(s.projPath, s.proj)

	projScope := s.goInspector.GetScope(goinspect.ProjectScope)
	components := s.goInspector.GetExportedComponents(projScope)
	components = append(components, s.goInspector.GetExportedComponents(s.goInspector.GetScope(goinspect.AmphionScope))...)
	compData := generators.MakeCompFileTemplateData(components, projScope.Module)

	mainData := generators.MakeMainTemplateData(s.runConfig, resources, compData)

	err = s.generateCommon(mainData)
	if err != nil {
		return err
	}

	builtinResPath := filepath.Join("templates", "res")
	err = utils.CopyDir(builtinResPath, filepath.Join(s.projPath, "res", "builtin"))
	if err != nil {
		return err
	}

	err = s.buildDelegate.Build(&BuildDelegateContext{
		projPath:         s.projPath,
		buildPath:        s.buildPath,
		proj:             s.proj,
		runConfig:         s.runConfig,
		mainTemplateData: mainData,
	})

	return
}

func (s *DevServer) generateCommon(mainData *generators.MainTemplateData) (err error) {
	err = generators.Res(mainData, s.projPath, s.proj)
	if err != nil {
		return
	}

	scope := s.goInspector.GetScope(goinspect.ProjectScope)
	components := s.goInspector.GetExportedComponents(scope)
	codePath := s.proj.GetCodePath(s.projPath)

	err = filepath.Walk(codePath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		packageComponents := make([]*goinspect.StructInfo, 0, len(components))
		relPath := strings.TrimPrefix(path, scope.Path)
		relPath = strings.TrimPrefix(relPath, "/")
		packagePath := scope.Module + "/" + relPath
		packagePath = strings.TrimSuffix(packagePath, "/")

		for _, comp := range components {
			if comp.Package != packagePath {
				continue
			}

			packageComponents = append(packageComponents, comp)
		}

		if len(packageComponents) == 0 {
			return nil
		}

		data := generators.MakeCompFileTemplateData(packageComponents, packagePath)
		return generators.Comp(data, path, "comp.gen.go")
	})

	return
}

func (s *DevServer) buildApp() (err error) {
	//1. Copy files from corresponding frontend folder to run folder
	frontendPath := filepath.Join(s.projPath, "frontend", s.runConfig.Frontend)
	err = utils.CopyDir(frontendPath, "run")
	if err != nil {
		return
	}

	//2. Prepare for debugging if necessary.
	if s.runConfig.Frontend == "web" {
		if s.runConfig.Debug {
			err = utils.CopyFile(filepath.Join("templates", "webDebug", "amphiondebug.js"), filepath.Join("run", "amphiondebug.js"))
		} else {
			_ = os.Remove(filepath.Join("run", "amphiondebug.js"))
		}
		if err != nil {
			return
		}
	}

	//3. Copy res folder to run folder
	resPath := filepath.Join(s.projPath, "res")
	err = utils.CopyDir(resPath, filepath.Join("run", "res"))
	if err != nil {
		return
	}

	//4. Copy build folder
	var buildPath string
	switch s.runConfig.Frontend {
	case "pc":
		buildPath = filepath.Join(s.buildPath, runtime.GOOS)
	case "web":
		buildPath = filepath.Join(s.buildPath, "web")
	}

	err = utils.CopyDir(buildPath, "run")
	if err != nil {
		return
	}

	//5. Generate app.yaml and copy to run
	app := project.NewAppFromConfig(s.proj, s.runConfig)
	data, err := yaml.Marshal(app)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filepath.Clean("./run/app.yaml"), data, os.FileMode(0777))

	return
}

func (s *DevServer) RunProject() (err error) {
	if s.runConfig.Frontend != "web" && s.runConfig.Frontend != "pc" {
		return errors.New("cannot run on this frontend")
	}

	fmt.Println("Running project...")

	//If on pc, we can run the app
	if s.runConfig.Frontend == "pc" {
		cmd := exec.Command("./" + executableName(s.proj, s.runConfig))
		cmd.Dir = "run"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//If on web, refresh page
	if s.runConfig.Frontend == "web" && s.webDebug != nil {
		s.webDebug.Refresh()
	}

	return
}

func executableName(proj *project.Config, runConfig *project.RunConfig) string {
	switch runConfig.Frontend {
	case "web":
		return "app.wasm"
	case "pc":
		switch runtime.GOOS {
		case "windows":
			return proj.Name + ".exe"
		default:
			return proj.Name
		}
	}

	return ""
}
