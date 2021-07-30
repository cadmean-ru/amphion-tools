package goinspect

import (
	"fmt"
	"go/ast"
	"strings"
)

const MissingMethodsName = "MissingMethods"

type FuncInfo struct {
	DeclarationInfo
	Receiver   *FieldInfo
	Parameters []*FieldInfo
	Returns    []*FieldInfo
}

func (f FuncInfo) String() string {
	recv := ""
	pack := f.Package + "."
	if f.Receiver != nil {
		recv = "(" + f.Receiver.String() + ") "
		pack = ""
	}

	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}

	returns := make([]string, len(f.Returns))
	for i, r := range f.Returns {
		returns[i] = r.String()
	}

	return fmt.Sprintf("func %s%s%s(%s) (%s)", recv, pack, f.Name, strings.Join(params, ", "), strings.Join(returns, ", "))
}

func (f FuncInfo) Matches(other *FuncInfo) bool {
	if f.Name != other.Name {
		return false
	}

	if len(f.Parameters) != len(other.Parameters) {
		return false
	}

	for i, p := range f.Parameters {
		if !p.Matches(other.Parameters[i]) {
			return false
		}
	}

	if len(f.Returns) != len(other.Returns) {
		return false
	}

	for i, r := range f.Returns {
		if !r.Matches(other.Returns[i]) {
			return false
		}
	}

	return true
}

func (f FuncInfo) IsMethod() bool {
	return f.Receiver != nil
}

func NewFuncInfo(name string) *FuncInfo {
	return &FuncInfo{
		DeclarationInfo: DeclarationInfo{
			Name: name,
		},
		Parameters: []*FieldInfo{},
		Returns:    []*FieldInfo{},
	}
}

func NewMethodInfo(name string, receiver *FieldInfo) *FuncInfo {
	return &FuncInfo{
		DeclarationInfo: DeclarationInfo{
			Name: name,
		},
		Parameters: []*FieldInfo{},
		Returns:    []*FieldInfo{},
		Receiver:   receiver,
	}
}

func tryParseFunc(node ast.Node, packageName string) (f *FuncInfo, ok bool) {
	var decl *ast.FuncDecl
	decl, ok = node.(*ast.FuncDecl)
	if !ok {
		return
	}

	var receiver *FieldInfo
	if decl.Recv != nil {
		receiver = parseFuncFieldList(decl.Recv)[0]
	}

	paramInfos := parseFuncFieldList(decl.Type.Params)
	returnInfos := parseFuncFieldList(decl.Type.Results)

	f = &FuncInfo{
		DeclarationInfo: DeclarationInfo{
			Name:    decl.Name.Name,
			Package: packageName,
		},
		Receiver:   receiver,
		Parameters: paramInfos,
		Returns:    returnInfos,
	}

	return
}

func parseFuncFieldList(list *ast.FieldList) []*FieldInfo {
	if list == nil || list.List == nil {
		return []*FieldInfo{}
	}

	params := make([]*FieldInfo, 0, len(list.List))
	for _, p := range list.List {
		params = append(params, parseField(p)...)
	}
	return params
}
