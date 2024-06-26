package ssa

import (
	"fmt"

	"github.com/yaklang/yaklang/common/log"
)

func (p *Package) NewFunction(name string) *Function {
	return p.NewFunctionWithParent(name, nil)
}

func (p *Package) NewFunctionWithParent(name string, parent *Function) *Function {
	index := len(p.Funcs)
	if index == 0 && name == "" {
		name = "main"
	}
	if name == "" {
		if parent != nil {
			name = fmt.Sprintf("%s$%d", parent.GetName(), index)
		} else {
			name = fmt.Sprintf("AnonymousFunc-%d", index)
		}
	}
	f := &Function{
		anValue:     NewValue(),
		Package:     p,
		Params:      make([]Value, 0),
		hasEllipsis: false,
		Blocks:      make([]Instruction, 0),
		EnterBlock:  nil,
		ExitBlock:   nil,
		ChildFuncs:  make([]Value, 0),
		parent:      nil,
		FreeValues:  make(map[string]Value),
		SideEffects: make([]*FunctionSideEffect, 0),
		builder:     nil,
	}
	f.SetName(name)

	if parent != nil {
		parent.addAnonymous(f)
		// Pos: parent.CurrentPos,
		f.SetRange(parent.builder.CurrentRange)
	} else {
		p.Funcs[name] = f
	}
	p.Prog.SetVirtualRegister(f)
	// function 's Range is essential!
	if f.GetRange() == nil {
		if editor := p.Prog.getCurrentEditor(); editor != nil {
			f.SetRangeInit(editor)
		} else {
			log.Warnf("the program must contains a editor to init function range: %v", p.Prog.Name)
		}
	}

	f.EnterBlock = f.NewBasicBlock("entry")
	return f
}

func (f *Function) GetType() Type {
	if f.Type != nil {
		return f.Type
	} else {
		return GetAnyType()
	}
}

func (f *Function) SetType(t Type) {
	if funTyp, ok := ToFunctionType(t); ok {
		f.Type = funTyp
	} else if t != nil {
		log.Warnf("ssa.Function type cannot covnert to FunctionType: %v", t)
	}
}

func (f *Function) SetGeneric(b bool) {
	f.isGeneric = b
}

func (f *Function) GetProgram() *Program {
	if f.Package == nil {
		return nil
	}
	return f.Package.Prog
}

func (f *Function) GetFunc() *Function {
	return f
}

func (f *Function) addAnonymous(anon *Function) {
	f.ChildFuncs = append(f.ChildFuncs, anon)
	anon.parent = f
}

func (f *FunctionBuilder) NewParam(name string, pos ...CanStartStopToken) *Parameter {
	p := NewParam(name, false, f)
	f.appendParam(p, pos...)
	return p
}

func (f *FunctionBuilder) NewParameterMember(name string, obj *Parameter, key Value) *ParameterMember {
	paraMember := NewParamMember(name, f, obj, key)
	f.ParameterMembers = append(f.ParameterMembers, paraMember)
	paraMember.FormalParameterIndex = len(f.ParameterMembers) - 1
	if f.MarkedThisObject != nil &&
		obj.GetDefault() != nil &&
		f.MarkedThisObject.GetId() == obj.GetDefault().GetId() {
		f.SetMethod(true, obj.GetType())
	}
	variable := f.CreateVariable(name)
	f.AssignVariable(variable, paraMember)
	return paraMember
}

func (f *FunctionBuilder) appendParam(p *Parameter, token ...CanStartStopToken) {
	f.Params = append(f.Params, p)
	p.FormalParameterIndex = len(f.Params) - 1
	p.IsFreeValue = false
	variable := f.CreateVariable(p.GetName(), token...)
	f.AssignVariable(variable, p)
}

func (f *Function) ReturnValue() []Value {
	exitBlock, ok := ToBasicBlock(f.ExitBlock)
	if !ok {
		log.Warnf("function exit block cannot convert to BasicBlock: %v", f.ExitBlock)
		return nil
	}
	ret := exitBlock.LastInst().(*Return)
	return ret.Results
}

func (f *Function) IsMain() bool {
	return f.GetName() == "main"
}

func (f *Function) GetParent() *Function {
	if f.parent == nil {
		return nil
	}

	fu, ok := ToFunction(f.parent)
	if ok {
		return fu
	}
	log.Warnf("function parent cannot convert to Function: %v", f.parent)
	return nil
}

// just create a function define, only function parameter type \ return type \ ellipsis
func NewFunctionWithType(name string, typ *FunctionType) *Function {
	f := &Function{
		anValue: NewValue(),
	}
	f.SetType(typ)
	f.SetName(name)
	return f
}

func (f *Function) IsMethod() bool {
	if f.Type == nil {
		f.Type = NewFunctionType("", nil, nil, false)
	}
	return f.Type.IsMethod
}

func (f *Function) SetMethod(is bool, objType Type) {
	if f.Type == nil {
		f.Type = NewFunctionType("", nil, nil, false)
	}
	f.Type.IsMethod = is
	f.Type.ObjectType = objType
}
