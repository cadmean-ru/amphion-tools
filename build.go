package main

import (
	"amphion-tools/server"
	. "fmt"
	"os"
)

func build() {
	args := os.Args

	if len(args) < 4 {
		Println("Not enough arguments")
		return
	}

	projectPath := args[2]
	configName := args[3]

	srv, err := server.StartDevelopment(projectPath, configName)
	if err != nil {
		Printf("Failed to start development server: %v\n", err)
		os.Exit(1)
	}

	err = srv.BuildProject()

	if err != nil {
		Printf("Build failed: %v\n", err)
		os.Exit(2)
	} else {
		Println("Build finished successfully")
	}
}
