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

	settings.Load()

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
		serve(false)
	case "last":
		serve(true)
	case "build":
		build()
	case "analyze":
		analyze()
	case "dev-serve":
		devServe()
	case "dev-generate":
		devGenerate()
	case "dev-goinspect":
		devGoInspect()
	default:
		fmt.Println("Unknown command")
	}
}