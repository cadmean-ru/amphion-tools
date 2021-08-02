package goinspect

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

const (
	AmphionScope = "amphion"
	ProjectScope = "project"
)

type Inspector struct {
	scopes       []*Scope
	amphionScope *Scope
	projectScope *Scope
}

func (i *Inspector) NewScope(name, path string) (*Scope, error) {
	goMod, err := ParseGoMod(path)
	if err != nil {
		return nil, err
	}

	s := &Scope{
		Name:       name,
		Path:       path,
		Module:     goMod.ModuleName,
		functions:  map[string]*FuncInfo{},
		structs:    map[string]*StructInfo{},
		interfaces: map[string]*InterfaceInfo{},
	}

	i.scopes = append(i.scopes, s)

	if name == AmphionScope {
		i.amphionScope = s
	} else if name == ProjectScope {
		i.projectScope = s
	}

	return s, nil
}

func (i *Inspector) GetScope(name string) *Scope {
	if name == AmphionScope {
		return i.amphionScope
	} else if name == ProjectScope {
		return i.projectScope
	}

	for _, scope := range i.scopes {
		if scope.Name == name {
			return scope
		}
	}

	return nil
}

func (i *Inspector) ForEachScope(action func(scope *Scope) bool) {
	for _, s := range i.scopes {
		if !action(s) {
			return
		}
	}
}

func (i *Inspector) GetExportedComponents(scope *Scope) []*StructInfo {
	comps := make([]*StructInfo, 0, 10)

	componentInterface := i.amphionScope.GetInterface("Component")

	for _, s := range scope.GetStructs() {
		if s.IsExported() && componentInterface.CheckImplements(s) {
			comps = append(comps, s)
		}
	}

	return comps
}

func (i *Inspector) InspectAmphion() (err error) {
	if i.amphionScope == nil {
		err = errors.New("no amphion scope")
		return
	}

	err = i.InspectSemantics(i.amphionScope, "common")
	if err != nil {
		return
	}

	err = i.InspectSemantics(i.amphionScope, "common/a")
	if err != nil {
		return
	}

	err = i.InspectSemantics(i.amphionScope, "rendering")
	if err != nil {
		return
	}

	err = i.InspectSemantics(i.amphionScope, "engine")
	if err != nil {
		return
	}

	err = i.InspectSemantics(i.amphionScope, "engine/builtin")
	if err != nil {
		return
	}

	return
}

func (i *Inspector) InspectComponents(scope *Scope) []string {
	messages := make([]string, 0, 10)
	componentInterface := i.amphionScope.GetInterface("Component")

	for _, cs := range scope.GetStructs() {
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
		Parameters: []*FieldInfo{},
		Returns: []*FieldInfo{{
			Name:        "_",
			TypeName:    "string",
			PassageKind: ByValue,
		}},
	}

	messages := make([]string, 0, 10)

	for _, m := range comp.GetAllMethods() {
		if m.Receiver.PassageKind == ByValue {
			messages = append(messages, fmt.Sprintf("Warning: Avoid using value receivers for component methods: %s.%s", comp.Package, comp.Name))
		}

		if getNameInfo.Matches(m) {
			messages = append(messages, fmt.Sprintf("Warning: Unnecessary method GetName() string: %s.%s", comp.Package, comp.Name))
		}
	}

	return messages
}

func (i *Inspector) InspectSemantics(scope *Scope, scopeRelativePath string) error {
	projectPath := filepath.Clean(scope.Path)
	scopeRelativePath = filepath.Clean(scopeRelativePath)

	path := filepath.Join(projectPath, scopeRelativePath)
	goMod, err := ParseGoMod(projectPath)
	if err != nil {
		return err
	}

	packageName := goMod.ModuleName + "/" + scopeRelativePath
	structs, funcs, interfaces, err := i.findDeclarations(path, packageName)
	if err != nil {
		return err
	}

	i.findMethods(structs, funcs)
	i.findTopLevelFunctions(scope, funcs)
	scope.AddStructs(structs...)
	scope.AddInterfaces(interfaces...)
	i.inspectStructEmbeddings(structs)
	i.inspectInterfaceEmbeddings(interfaces)

	return nil
}

func (i *Inspector) findTopLevelFunctions(scope *Scope, funcs []*FuncInfo) {
	for _, f := range funcs {
		if !f.IsMethod() {
			scope.functions[f.Name] = f
		}
	}
}

func (i *Inspector) inspectStructEmbeddings(structs []*StructInfo) {
	for _, s := range structs {
		for _, e := range s.Embeddings {
			i.ForEachScope(func(scope *Scope) bool {
				for _, ss := range scope.structs {
					if e.TypeName == ss.Name {
						e.StructInfo = ss
						return false
					}
				}
				return true
			})
		}
	}
}

func (i *Inspector) inspectInterfaceEmbeddings(structs []*InterfaceInfo) {
	for _, ii := range structs {
		for _, e := range ii.Embeddings {
			i.ForEachScope(func(scope *Scope) bool {
				for _, iii := range scope.interfaces {
					if e.TypeName == iii.Name {
						e.InterfaceInfo = iii
						return false
					}
				}
				return true
			})
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
		scopes: make([]*Scope, 0, 2),
	}
}
