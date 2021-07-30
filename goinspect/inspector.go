package goinspect

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type Inspector struct {
	functions  []*FuncInfo
	structs    []*StructInfo
	interfaces []*InterfaceInfo
}

func (i *Inspector) GetFunctions() []*FuncInfo {
	return i.functions
}

func (i *Inspector) GetFunction(name string) *FuncInfo {
	for _, f := range i.functions {
		if f.Name == name {
			return f
		}
	}

	return nil
}

func (i *Inspector) GetStructs() []*StructInfo {
	return i.structs
}

func (i *Inspector) GetStruct(name string) *StructInfo {
	for _, s := range i.structs {
		if s.Name == name {
			return s
		}
	}

	return nil
}

func (i *Inspector) GetInterfaces() []*InterfaceInfo {
	return i.interfaces
}

func (i *Inspector) GetInterface(name string) *InterfaceInfo {
	for _, ii := range i.interfaces {
		if ii.Name == name {
			return ii
		}
	}

	return nil
}

func (i *Inspector) InspectComponents() []string {
	componentInterface := i.GetInterface("Component")
	messages := make([]string, 0, 10)

	for _, cs := range i.GetStructs() {
		if componentInterface.CheckImplements(cs) {
			messages = append(messages, i.InspectComponent(cs)...)
		}
	}

	return messages
}

func (i *Inspector) InspectComponent(comp *StructInfo) []string {
	getNameInfo := &FuncInfo{
		DeclarationInfo: DeclarationInfo{
			Name:    "GetName",
			Package: "",
		},
		Receiver:   nil,
		Parameters: []*FieldInfo {},
		Returns:    []*FieldInfo {{
			Name:        "_",
			TypeName:    "string",
			PassageKind: ByValue,
		}},
	}

	messages := make([]string, 0, 10)

	for _, m := range comp.GetAllMethods() {
		if m.Receiver.PassageKind == ByValue {
			messages = append(messages, fmt.Sprintf("Warning: Avoid using value receivers for component methods: %s", comp.Name))
		}

		if getNameInfo.Matches(m) {
			messages = append(messages, fmt.Sprintf("Warning: Unnecessary method GetName() string: %s", comp.Name))
		}
	}

	return messages
}

func (i *Inspector) InspectSemantics(projectPath, projectRelativePath string) error {
	projectPath = filepath.Clean(projectPath)
	projectRelativePath = filepath.Clean(projectRelativePath)

	path := filepath.Join(projectPath, projectRelativePath)
	goMod, err := ParseGoMod(projectPath)
	if err != nil {
		return err
	}

	packageName := goMod.ModuleName + "/" + projectRelativePath
	structs, funcs, interfaces, err := i.findDeclarations(path, packageName)
	if err != nil {
		return err
	}

	i.findMethods(structs, funcs)
	i.findTopLevelFunctions(funcs)
	i.structs = append(i.structs, structs...)
	i.interfaces = append(i.interfaces, interfaces...)
	i.inspectStructEmbeddings(structs)
	i.inspectInterfaceEmbeddings(interfaces)

	return nil
}

func (i *Inspector) findTopLevelFunctions(funcs []*FuncInfo) {
	for _, f := range funcs {
		if !f.IsMethod() {
			i.functions = append(i.functions, f)
		}
	}
}

func (i *Inspector) inspectStructEmbeddings(structs []*StructInfo) {
	for _, s := range structs {
		for _, e := range s.Embeddings {
			for _, ss := range i.structs {
				if e.TypeName == ss.Name {
					e.StructInfo = ss
					break
				}
			}
		}
	}
}

func (i *Inspector) inspectInterfaceEmbeddings(structs []*InterfaceInfo) {
	for _, ii := range structs {
		for _, e := range ii.Embeddings {
			for _, iii := range i.interfaces {
				if e.TypeName == iii.Name {
					e.InterfaceInfo = iii
					break
				}
			}
		}
	}
}

func (i *Inspector) findMethods(structList []*StructInfo, funcList []*FuncInfo) {
	for _, f := range funcList {
		if f.Receiver == nil {
			continue
		}

		for _, s := range structList {
			if f.Receiver.TypeName == s.Name {
				s.Methods = append(s.Methods, f)
			}
		}
	}
}

func (i *Inspector) findDeclarations(codePath, fullPackageName string) ([]*StructInfo, []*FuncInfo, []*InterfaceInfo, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, codePath, nil, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	structList := make([]*StructInfo, 0, 10)
	funcList := make([]*FuncInfo, 0, 10)
	interfaceList := make([]*InterfaceInfo, 0, 10)

	for _, pack := range packages {
		ast.Inspect(pack, func(node ast.Node) bool {
			if structInfo, ok := tryParseStruct(node, fullPackageName); ok {
				structList = append(structList, structInfo)
			} else if interfaceInfo, ok := tryParseInterface(node, fullPackageName); ok {
				interfaceList = append(interfaceList, interfaceInfo)
			} else if funcInfo, ok := tryParseFunc(node, fullPackageName); ok {
				funcList = append(funcList, funcInfo)
			}
			return true
		})
	}

	return structList, funcList, interfaceList, nil
}

func NewInspector() *Inspector {
	return &Inspector{
		functions:  []*FuncInfo{},
		structs:    []*StructInfo{},
		interfaces: []*InterfaceInfo{},
	}
}
