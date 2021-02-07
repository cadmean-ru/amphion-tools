package main

import (
	"fmt"
	"os"
)

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
	case "analyze":
		analyze()
	default:
		fmt.Println("Unknown command")
	}
}