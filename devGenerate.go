package main

import (
	"amphion-tools/generators"
	"amphion-tools/goinspect"
	"fmt"
	"os"
)

func devGenerate() {
	args := os.Args
	if len(args) < 6 {
		fmt.Println("not enough arguments")
		return
	}

	path := os.Args[2]
	packPath := os.Args[3]
	pack := os.Args[4]
	file := os.Args[5]

	inspector := goinspect.NewInspector()
	scope, err := inspector.NewScope(goinspect.AmphionScope, path)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = inspector.InspectSemantics(scope, "common")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = inspector.InspectSemantics(scope, "common/a")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = inspector.InspectSemantics(scope, "rendering")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = inspector.InspectSemantics(scope, "engine")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = inspector.InspectSemantics(scope, "engine/builtin")
	if err != nil {
		fmt.Println(err)
		return
	}

	components := make([]*goinspect.StructInfo, 0)
	for _, comp := range inspector.GetExportedComponents(scope) {
		if comp.Package == pack {
			components = append(components, comp)
		}
	}

	if len(components) == 0 {
		fmt.Println("no components found")
		return
	}

	data := generators.MakeCompFileTemplateData(components, pack)
	err = generators.Comp(data, packPath, file)
	fmt.Println(err)
}
