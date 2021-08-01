package main

import (
	"amphion-tools/settings"
	"amphion-tools/support"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Amphion tools v%s\n", support.ToolsVersion)
	fmt.Println("©Cadmean 2021")

	args := os.Args

	s := settings.Load()
	lastProjectPath := ""
	if s.LastProject != nil {
		fmt.Printf("Last project: %s - %s\n", s.LastProject.Name, s.LastProject.Path)
		fmt.Println("Use command \"last\" to serve it.")
		lastProjectPath = s.LastProject.Path
	}

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
		serve("")
	case "last":
		serve(lastProjectPath)
	case "analyze":
		analyze()
	case "dev-serve":
		devServe()
	case "dev-generate":
		devGenerate()
	default:
		fmt.Println("Unknown command")
	}
}