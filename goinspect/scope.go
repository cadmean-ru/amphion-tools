package goinspect

type Scope struct {
	Name       string
	Path       string
	Module     string
	functions  map[string]*FuncInfo
	structs    map[string]*StructInfo
	interfaces map[string]*InterfaceInfo
}

func (s *Scope) GetFunctions() map[string]*FuncInfo {
	return s.functions
}

func (s *Scope) GetFunction(name string) *FuncInfo {
	for n, f := range s.functions {
		if n == name {
			return f
		}
	}

	return nil
}

func (s *Scope) AddFunctions(funcs ...*FuncInfo) {
	for _, f := range funcs {
		s.functions[f.Name] = f
	}
}

func (s *Scope) GetStructs() map[string]*StructInfo {
	return s.structs
}

func (s *Scope) GetStruct(name string) *StructInfo {
	for n, st := range s.structs {
		if n == name {
			return st
		}
	}

	return nil
}

func (s *Scope) AddStructs(structs ...*StructInfo) {
	for _, str := range structs {
		s.structs[str.Name] = str
	}
}

func (s *Scope) GetInterfaces() map[string]*InterfaceInfo {
	return s.interfaces
}

func (s *Scope) GetInterface(name string) *InterfaceInfo {
	for n, i := range s.interfaces {
		if n == name {
			return i
		}
	}

	return nil
}

func (s *Scope) AddInterfaces(interfaces ...*InterfaceInfo) {
	for _, i := range interfaces {
		s.interfaces[i.Name] = i
	}
}

func (s *Scope) Clear() {
	s.functions = make(map[string]*FuncInfo)
	s.structs = make(map[string]*StructInfo)
	s.interfaces = make(map[string]*InterfaceInfo)
}