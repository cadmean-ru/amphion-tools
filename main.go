package main

import (
	"amphion-tools/server"
	"amphion-tools/utils"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// /Users/alex/Projects/AmphionEngine/projects

func main() {
	args := os.Args

	var command string

	if len(args) < 2 {
		fmt.Print("Enter command: ")
		fmt.Scanln(&command)
	} else {
		command = args[1]
	}

	switch command {
	case "create":
		createProjectInteractive()
	case "serve":
		serve()
	default:
		fmt.Println("Unknown command")
	}
}

func createProjectInteractive() {
	var name string
	fmt.Print("Enter project name: ")
	fmt.Scanln(&name)

	for !checkProjectName(name) {
		fmt.Print("Invalid project name. Try again: ")
		fmt.Scanln(&name)
	}

	var path string
	fmt.Print("Enter project path (where the project directory will be created): ")
	fmt.Scanln(&path)

	var author string
	fmt.Print("Enter author name, e.g. your name or company name: ")
	fmt.Scanln(&author)

	for len(author) < 2 {
		fmt.Print("Invalid author name. Try again: ")
		fmt.Scanln(&author)
	}

	var domain string
	fmt.Print("Enter company domain (e.g. cadmean.ru) [optional]: ")
	fmt.Scanln(&domain)

	var gitUrl string
	fmt.Print("Enter git repository url [optional]: ")
	fmt.Scanln(&gitUrl)

	err := createProject(path, name, author, domain, gitUrl)
	if err != nil {
		fmt.Println("Failed to create a project:")
		fmt.Println(err)
	} else {
		fmt.Println("Project was successfully created")
	}
}

var projectNameRegex = regexp.MustCompile("^[A-z]+[A-z0-9-]{2,}$")
var invalidProjectNames = []string { "amphion", "build", "res", "frontend", "pc", "web" }

func checkProjectName(name string) bool {
	for _, n := range invalidProjectNames {
		if n == name {
			return false
		}
	}

	return projectNameRegex.MatchString(name)
}

func createProject(path, name, author, companyDomain, gitUrl string) (err error) {
	fullProjectPath := filepath.Join(path, name)

	_, err = os.Stat(fullProjectPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.FileMode(0777)); err != nil {
			return
		}
	} else if err != nil {
		return
	} else {
		return fmt.Errorf("project already exists")
	}

	templateDirPath := "./templates/basicProject"
	err = utils.CopyDir(templateDirPath, fullProjectPath)
	if err != nil {
		return
	}

	filesToDelete := make([]string, 0)
	filesToRename := make([]string, 0)
	err = filepath.Walk(fullProjectPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), "__PROJECT_NAME__") {
			filesToRename = append(filesToRename, path)
		} else if info.Name() == "deleteme.txt" {
			filesToDelete = append(filesToDelete, path)
		}
		return nil
	})

	for _, f := range filesToRename {
		d, n := filepath.Split(f)
		newFileName := strings.ReplaceAll(n, "__PROJECT_NAME__", name)
		newPath := filepath.Join(d, newFileName)
		_ = os.Rename(f, newPath)
	}

	for _, f := range filesToDelete {
		_ = os.Remove(f)
	}

	config := struct {
		ProjectName   string
		Author        string
		CompanyDomain string
		GitUrl        string
	}{
		ProjectName:   name,
		Author:        author,
		CompanyDomain: companyDomain,
		GitUrl:        gitUrl,
	}

	configFilePath := filepath.Join(fullProjectPath, name + ".config.yaml")
	tmpl := template.Must(template.ParseFiles(configFilePath))

	configFile, err := os.Create(configFilePath)
	if err != nil {
		return
	}
	defer configFile.Close()

	err = tmpl.Execute(configFile, config)
	if err != nil {
		return
	}

	cmd := exec.Command("git", "init", fullProjectPath)
	err = cmd.Run()
	if err != nil {
		return
	}

	codeDirPath := filepath.Join(fullProjectPath, name)

	cmd = exec.Command("go", "mod", "init", name)
	cmd.Dir = codeDirPath
	err = cmd.Run()
	if err != nil {
		return
	}

	cmd = exec.Command("go", "get", "-u", "github.com/cadmean-ru/amphion@v0.1.0rc4")
	cmd.Dir = codeDirPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	fmt.Printf("%s\n", output)

	return
}

func serve() {
	var projectPath, runConfig string

	if len(os.Args) < 4 {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Enter project path: ")
		scanner.Scan()
		projectPath = scanner.Text()

		fmt.Print("Enter run config: ")
		scanner.Scan()
		runConfig = scanner.Text()
	} else {
		projectPath = os.Args[2]
		runConfig = os.Args[3]
	}

	s, err := server.StartDevelopment(projectPath, runConfig)
	if err != nil {
		fmt.Println("Failed to start development server")
		return
	}

	s.Start()
	err = s.BuildProject()
	if err != nil {
		fmt.Println(err)
	}
	err = s.RunProject()
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Print("Enter command: ")
		var input string
		fmt.Scanln(&input)

		if input == "exit" {
			break
		}

		switch input {
		case "build":
			err = s.BuildProject()
			if err != nil {
				fmt.Println(err)
			}
		case "run":
			err = s.RunProject()
			if err != nil {
				fmt.Println(err)
			}
		case "br":
			err = s.BuildProject()
			if err != nil {
				fmt.Println(err)
			}
			err = s.RunProject()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	s.Stop()
}