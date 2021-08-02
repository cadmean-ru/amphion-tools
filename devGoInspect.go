package main

import (
	"amphion-tools/goinspect"
	"fmt"
	"github.com/TwinProduction/go-color"
	"os"
)

func devGoInspect() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Not enough arguments")
		return
	}

	path := os.Args[2]

	inspector := goinspect.NewInspector()
	amphionScope, err := inspector.NewScope(goinspect.AmphionScope, path)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = inspector.InspectAmphion()
	if err != nil {
		fmt.Println(err)
		return
	}

	messages := inspector.InspectComponents(amphionScope)

	fmt.Printf("Total messages: %d\n", len(messages))
	fmt.Println()

	for _, m := range messages {
		fmt.Println(color.Ize(color.Yellow, m))
	}

	fmt.Println()
}
