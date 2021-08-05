package main

import (
	"amphion-tools/project"
	"amphion-tools/server"
	"amphion-tools/settings"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func serve(lastProject bool) {
	var projectPath, runConfig string

	set := settings.Current
	lastProjectPath := set.LastProject.Path
	if !lastProject && set.LastProject != nil {
		fmt.Printf("Last project: %s - %s\n", set.LastProject.Name, set.LastProject.Path)
		fmt.Println("Enter \"last\" project path to serve it.")
	}

	if len(os.Args) < 4 {
		scanner := bufio.NewScanner(os.Stdin)

		if lastProject {
			projectPath = lastProjectPath
		} else {
			fmt.Print("Enter project path: ")
			scanner.Scan()
			projectPath = scanner.Text()

			if projectPath == "last" && lastProjectPath != "" {
				projectPath = lastProjectPath
			}
		}

		p, err := project.FindProjectConfig(projectPath)
		if err != nil {
			panic("failed to find project config file")
		}

		fmt.Println("Select run config:")

		for i, conf := range p.Configurations {
			fmt.Printf("%d - %s (%s)\n", i, conf.Name, conf.Frontend)
		}

		scanner.Scan()
		numStr := scanner.Text()
		var num int
		num, err = strconv.Atoi(numStr)
		if err != nil || num < 0 || num > len(p.Configurations) {
			num = 0
		}

		runConfig = p.Configurations[num].Name

		fmt.Printf("Selected config: %s\n", runConfig)
	} else {
		projectPath = os.Args[2]
		runConfig = os.Args[3]
	}

	s, err := server.StartDevelopment(projectPath, runConfig)
	if err != nil {
		fmt.Printf("Failed to start development server: %s\n", err)
		return
	}

	fmt.Println("Development server started")

	s.Start()
	err = s.BuildProject()
	if err != nil {
		fmt.Println(err)
	} else {
		err = s.RunProject()
		if err != nil {
			fmt.Println(err)
		}
	}

	for {
		fmt.Print("Enter command: ")
		var input string
		fmt.Scanln(&input)

		if input == "exit" {
			break
		}

		switch input {
		case "build", "b":
			err = s.BuildProject()
			if err != nil {
				fmt.Printf("Build failed: %s\n", err)
			}
		case "run", "r":
			err = s.RunProject()
			if err != nil {
				fmt.Printf("Failed to run project: %s\n", err)
			}
		case "br":
			err = s.BuildProject()
			if err != nil {
				fmt.Printf("Build failed: %s\n", err)
				continue
			}
			err = s.RunProject()
			if err != nil {
				fmt.Printf("Failed to run project: %s\n", err)
			}
		}
	}

	fmt.Println("Exiting...")
	s.Stop()
	fmt.Println("Bye")
}
