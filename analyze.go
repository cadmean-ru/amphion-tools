package main

import (
	"amphion-tools/analysis"
	"fmt"
	"os"
)

func analyze() {
	var what, path string

	if len(os.Args) < 4 {
		fmt.Print("Enter what to analyze (dependencies, lines, components):")
		fmt.Scanln(&what)

		fmt.Print("Enter path:")
		fmt.Scanln(&path)
	} else {
		what = os.Args[2]
		path = os.Args[3]
	}

	switch what {
	case "dependencies":
		analyzeDependencies(path)
	case "lines":
		analyzeLines(path)
	case "components":
		analyzeComponents(path)
	default:
		fmt.Println("Dont know what to analyze")
	}
}

func analyzeDependencies(path string) {
	deps, err := analysis.GetProjectDependencies(path)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("List of dependencies of the project %s\n\n", path)
		for _, dep := range deps {
			fmt.Println(dep.ToString())
			for _, usedBy := range dep.UsedBy {
				fmt.Printf("\t-%s\n", usedBy.ToString())
			}
			fmt.Println()
		}
	}
}

func analyzeLines(path string) {
	count := analysis.CountLines(path)
	fmt.Printf("Counted lines in directory: %s\n", path)
	fmt.Printf("Total: %d\n", count.Total)
	fmt.Printf("Not empty: %d\n", count.NotEmpty)
	fmt.Printf("Code: %d\n", count.Code)
	fmt.Printf("Comments: %d\n", count.Comments)
}

func analyzeComponents(path string) {
	//inspector := goinspect.NewInspector()
	//
	//err := inspector.InspectSemantics("/Users/alex/Projects/AmphionEngine/amphion", "common")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//err = inspector.InspectSemantics("/Users/alex/Projects/AmphionEngine/amphion", "common/a")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//err = inspector.InspectSemantics("/Users/alex/Projects/AmphionEngine/amphion", "rendering")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//err = inspector.InspectSemantics("/Users/alex/Projects/AmphionEngine/amphion", "engine")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//err = inspector.InspectSemantics("/Users/alex/Projects/AmphionEngine/amphion", "engine/builtin")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//err = inspector.InspectSemantics(path, "components")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println("Found top level functions:")
	//for _, f := range inspector.GetFunctions() {
	//	fmt.Println(f)
	//}
	//
	//fmt.Println()
	//fmt.Println("Found interfaces:")
	//for _, i := range inspector.GetInterfaces() {
	//	fmt.Println(i)
	//}
	//
	//fmt.Println()
	//fmt.Println("Found structs:")
	//for _, s := range inspector.GetStructs() {
	//	fmt.Println(s)
	//}
	//
	//componentInterface := inspector.GetInterface("Component")
	//
	//fmt.Println()
	//fmt.Println("Found components:")
	//for _, s := range inspector.GetStructs() {
	//	if componentInterface.CheckImplements(s) {
	//		fmt.Println(s.Name)
	//	}
	//}

	fmt.Println("deprecated")
}