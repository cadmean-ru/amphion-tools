package server

import (
	"amphion-tools/analysis"
	"amphion-tools/generators"
	"amphion-tools/project"
	"amphion-tools/support"
	"amphion-tools/utils"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	deps, err := analysis.GetProjectDependencies(filepath.Join(projectPath, config.Name))
	if err != nil {
		return nil, err
	}

	var ok bool
	for _, dep := range deps {
		if dep.Name == "github.com/cadmean-ru/amphion" {
			if support.IsAmphionVersionSupported(dep.Version) {
				ok = true
			}
			break
		}
	}

	if !ok {
		return nil, errors.New("unsupported amphion version")
	}

	s = &DevServer{
		runConfig: runConfig,
		done:      make(chan bool),
		stopped:   false,
		proj:      config,
		projPath:  projectPath,
		buildPath: filepath.Join(projectPath, "build"),
	}

	return
}


type DevServer struct {
	runConfig *project.RunConfig
	done      chan bool
	stopped   bool
	proj      *project.Config
	projPath  string
	buildPath string
	webDebug  *WebDebugServer
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

	//dirs := []string { filepath.Join(s.projPath, "res"), filepath.Join(s.projPath, s.proj.Name) }
	//go utils.WatchDirs(dirs, s.done, s.handleDirectoryUpdate)

	if s.runConfig.Frontend == "web" {
		go utils.HttpServeDir("./run", s.runConfig.Url, s.done)

		if s.runConfig.Debug {
			s.webDebug = NewWebDebugServer("4242")
			s.webDebug.Start()
			utils.OpenBrowser("http://" + s.runConfig.Url + "?connectDebugger=4242")
		} else {
			utils.OpenBrowser("http://" + s.runConfig.Url)
		}
	}
}

func (s *DevServer) handleDirectoryUpdate() {
	//_ = s.BuildProject()
}

func (s *DevServer) BuildProject() (err error) {
	srcPath := filepath.Join(s.projPath, s.proj.Name)

	//0. Copy builtin res folder to run folder
	builtinResPath := filepath.Join("templates", "res")
	err = utils.CopyDir(builtinResPath, filepath.Join(s.projPath, "res", "builtin"))
	if err != nil {
		return err
	}

	//1. Generate code
	err = generators.Main(s.projPath, s.proj, s.runConfig)
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
		dstFileName = s.proj.Name
		goos = os.Getenv("GOOS")
		goarch = os.Getenv("GOARCH")
	case "web":
		dstPath = filepath.Join(s.buildPath, "web")
		dstFileName = "app.wasm"
		goos = "js"
		goarch = "wasm"
	default:
		return fmt.Errorf("unknown platform")
	}

	_ = os.Mkdir(dstPath, os.FileMode(0777))

	err = goBuild(srcPath, dstPath, dstFileName, goos, goarch)

	return
}

func (s *DevServer) RunProject() (err error) {
	//1. Copy files from corresponding frontend folder to run folder
	frontendPath := filepath.Join(s.projPath, "frontend", s.runConfig.Frontend)
	err = utils.CopyDir(frontendPath, "./run")
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

	err = utils.CopyDir(buildPath, "./run")
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
		cmd := exec.Command("./" + s.proj.Name)
		cmd.Dir = "run"
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", output)
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
	build.Env = append(build.Env, "GOOS=" + goos)
	build.Env = append(build.Env, "GOARCH=" + goarch)

	output, err := build.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", output)
		return
	}

	fmt.Printf("%s\n", output)

	return
}