package goinspect

import (
	"fmt"
	"strings"
)

const MissingMethodsName = "MissingMethods"

type FuncInfo struct {
	Name       string
	Receiver   *FieldInfo
	Parameters []*FieldInfo
	Returns    []*FieldInfo
}

func (f FuncInfo) String() string {
	recv := ""
	if f.Receiver != nil {
		recv = "(" + f.Receiver.String() + ") "
	}

	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}

	returns := make([]string, len(f.Returns))
	for i, r := range f.Returns {
		returns[i] = r.String()
	}

	return fmt.Sprintf("func %s%s(%s) (%s)", recv, f.Name, strings.Join(params, ", "), strings.Join(returns, ", "))
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
		Name:       name,
		Parameters: []*FieldInfo{},
		Returns:    []*FieldInfo{},
	}
}

func NewMethodInfo(name string, receiver *FieldInfo) *FuncInfo {
	return &FuncInfo{
		Name:       name,
		Parameters: []*FieldInfo{},
		Returns:    []*FieldInfo{},
		Receiver:   receiver,
	}
}