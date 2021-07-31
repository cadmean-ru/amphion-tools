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
		runConfig:    runConfig,
		done:         make(chan bool, 1),
		stopped:      false,
		proj:         config,
		projPath:     projectPath,
		buildPath:    filepath.Join(projectPath, "build"),
		amVersion:    amVersion,
		goInspector:  goinspect.NewInspector(),
		resInspector: resinspect.NewInspector(),
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
	runConfig    *project.RunConfig
	done         chan bool
	stopped      bool
	proj         *project.Config
	projPath     string
	buildPath    string
	webDebug     *WebDebugServer
	amVersion    string
	goInspector  *goinspect.Inspector
	resInspector *resinspect.Inspector
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
	amPath := filepath.Join(settings.Current.GoRoot, "pkg", "mod", "github.com", "cadmean-ru", "amphion@"+s.amVersion)
	if !utils.Exists(amPath) {
		return errors.New("amphion not found in GOPATH (try running 'go get')")
	}

	codePath := filepath.Join(s.projPath, s.proj.Name)

	amScope, _ := s.goInspector.NewScope("amphion", amPath)
	_, _ = s.goInspector.NewScope("project", codePath)

	err := s.goInspector.InspectSemantics(amScope, "common")
	if err != nil {
		return err
	}

	err = s.goInspector.InspectSemantics(amScope, "common/a")
	if err != nil {
		return err
	}

	err = s.goInspector.InspectSemantics(amScope, "rendering")
	if err != nil {
		return err
	}

	err = s.goInspector.InspectSemantics(amScope, "engine")
	if err != nil {
		return err
	}

	err = s.goInspector.InspectSemantics(amScope, "engine/builtin")
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

	switch s.runConfig.Frontend {
	case "android":
		return s.androidBuild(mainData)
	case "ios":
		return s.iosBuild(mainData)
	default:
		return s.defaultBuild(mainData)
	}
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
		return generators.Comp(data, path)
	})

	return
}

func (s *DevServer) defaultBuild(mainData *generators.MainTemplateData) (err error) {
	srcPath := filepath.Join(s.projPath, s.proj.Name)

	//0. Copy builtin res folder to run folder
	builtinResPath := filepath.Join("templates", "res")
	err = utils.CopyDir(builtinResPath, filepath.Join(s.projPath, "res", "builtin"))
	if err != nil {
		return err
	}

	//1. Generate code
	err = generators.Main(mainData, s.projPath, s.proj, s.runConfig)
	if err != nil {
		return
	}

	//2. Run go build
	var dstPath string
	var dstFileName string
	var goos, goarch string

	switch s.runConfig.Frontend {
	case "pc":
		dstPath = filepath.Join(s.buildPath, runtime.GOOS)
		dstFileName = executableName(s.proj, s.runConfig)
		goos = os.Getenv("GOOS")
		goarch = os.Getenv("GOARCH")
	case "web":
		dstPath = filepath.Join(s.buildPath, "web")
		dstFileName = executableName(s.proj, s.runConfig)
		goos = "js"
		goarch = "wasm"
	default:
		return fmt.Errorf("unknown platform")
	}

	_ = os.Mkdir(dstPath, os.FileMode(0777))

	err = goBuild(srcPath, dstPath, dstFileName, goos, goarch)

	return
}

func (s *DevServer) androidBuild(mainData *generators.MainTemplateData) (err error) {
	srcPath := filepath.Join(s.projPath, s.proj.Name)

	//1. Generate code
	err = generators.AndroidMain(mainData, s.projPath, s.proj, s.runConfig)
	if err != nil {
		return
	}
	err = generators.Android(mainData, s.projPath, s.proj, s.runConfig)
	if err != nil {
		return
	}

	//2. Run gomobile bind
	var dstPath = filepath.Join(s.buildPath, s.proj.Name+".android.aar")

	err = s.goMobileBind("android", srcPath, dstPath)

	if err == nil {
		fmt.Printf("Android library file was successfully created: %s\n", dstPath)
	}

	return
}

func (s *DevServer) iosBuild(mainData *generators.MainTemplateData) (err error) {
	srcPath := filepath.Join(s.projPath, s.proj.Name)

	//1. Generate code
	err = generators.IosMain(mainData, s.projPath, s.proj, s.runConfig)
	if err != nil {
		return
	}
	err = generators.Ios(mainData, s.projPath, s.proj, s.runConfig)
	if err != nil {
		return
	}

	//2. Run gomobile bind
	var dstPath = filepath.Join(s.buildPath, "Amphion.framework")

	err = s.goMobileBind("ios", srcPath, dstPath)

	if err == nil {
		fmt.Printf("iOS framework was successfully created: %s\n", dstPath)
	}

	return
}

func (s *DevServer) InspectCode() error {
	if settings.Current.GoRoot == "" {
		s.inputGoRoot()
	}

	fmt.Println("Running code inspection...")

	prScope := s.goInspector.GetScope("project")
	prScope.Clear()

	_ = filepath.Walk(prScope.Path, func(path string, info fs.FileInfo, err error) error {
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

	for _, msg := range s.goInspector.InspectComponents(prScope) {
		fmt.Println(ccolor.Ize(ccolor.Yellow, msg))
	}

	return nil
}

func (s *DevServer) inputGoRoot() {
	fmt.Print("GOROOT no defined. Please enter a valid GOROOT path:")
	fmt.Scanln(&settings.Current.GoRoot)
	_, err := os.Stat(settings.Current.GoRoot)
	if err != nil {
		panic(err)
	}
	_ = settings.Save(settings.Current)
}

func (s *DevServer) RunProject() (err error) {
	if s.runConfig.Frontend != "web" && s.runConfig.Frontend != "pc" {
		return errors.New("cannot run on this frontend")
	}

	fmt.Println("Running project...")

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

	//6. If on pc, we can run the app
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

func goBuild(srcPath, dstPath, dstFileName, goos, goarch string) (err error) {
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

func (s *DevServer) goMobileBind(target, srcPath, dstFilePath string) (err error) {

	var bind *exec.Cmd

	if target == "android" {
		bind = exec.Command("gomobile",
			"bind",
			"-target="+target,
			"-o", dstFilePath,
			"-javapkg=ru.cadmean.amphion.android",
			s.proj.Name+"/generated/droidCli",
			"github.com/cadmean-ru/amphion/frontend/cli",
			"github.com/cadmean-ru/amphion/common/atext",
			"github.com/cadmean-ru/amphion/common/dispatch",
		)
	} else {
		bind = exec.Command("gomobile",
			"bind",
			"-target="+target,
			"-o", dstFilePath,
			s.proj.Name+"/generated/iosCli",
			"github.com/cadmean-ru/amphion/frontend/cli",
			"github.com/cadmean-ru/amphion/common/atext",
			"github.com/cadmean-ru/amphion/common/dispatch",
		)
	}

	bind.Dir = srcPath
	bind.Env = os.Environ()
	if target == "android" {
		bind.Env = append(bind.Env, "ANDROID_NDK_HOME=/Users/alex/Library/Android/sdk/ndk/23.0.7196353")
		bind.Env = append(bind.Env, "ANDROID_HOME=/Users/alex/Library/Android/sdk")
	}

	output, err := bind.CombinedOutput()
	fmt.Printf("%s\n", output)

	if err != nil {
		return
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
