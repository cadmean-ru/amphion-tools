package main

import (
	"amphion-tools/support"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Amphion tools v%s\n", support.ToolsVersion)
	fmt.Println("©Cadmean 2021")

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
	case "analyze":
		analyze()
	case "dev-serve":
		devServe()
	default:
		fmt.Println("Unknown command")
	}
}