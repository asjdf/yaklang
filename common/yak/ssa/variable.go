package ssa

import "github.com/yaklang/yaklang/common/utils"

// --------------- for assign
type LeftValue interface {
	Assign(Value, *FunctionBuilder)
	GetPosition() *Position
	GetValue(*FunctionBuilder) Value
}

// --------------- only point variable to value with `f.currentDef`
// --------------- is SSA value
type IdentifierLV struct {
	name         string
	pos          *Position
	isSideEffect bool
}

func (i *IdentifierLV) Assign(v Value, f *FunctionBuilder) {
	v.AddLeftPositions(i.GetPosition())
	f.WriteVariable(i.name, v)
	if i.isSideEffect {
		f.AddSideEffect(i.name, v)
	}
}

func (i *IdentifierLV) GetValue(f *FunctionBuilder) Value {
	v := f.ReadVariable(i.name, true)
	return v
}

func (i *IdentifierLV) GetPosition() *Position {
	return i.pos
}
func (i *IdentifierLV) SetIsSideEffect(b bool) {
	i.isSideEffect = b
}

func NewIdentifierLV(variable string, pos *Position) *IdentifierLV {
	return &IdentifierLV{
		name: variable,
		pos:  pos,
	}
}

var _ LeftValue = (*IdentifierLV)(nil)

// --------------- point variable to value `f.symbol[variable]value`
// --------------- it's memory address, not SSA value
func (field *Field) Assign(v Value, f *FunctionBuilder) {
	f.EmitUpdate(field, v)
}

func (f *Field) GetValue(_ *FunctionBuilder) Value {
	return f
}

var _ LeftValue = (*Field)(nil)

// --------------- `f.currentDef` handler, read && write
func (f *FunctionBuilder) WriteVariable(variable string, value Value) {
	variable = f.GetScopeLocalVariable(variable)
	f.writeVariableByBlock(variable, value, f.CurrentBlock)
}

func (f *Function) ReplaceVariable(variable string, v, to Value) {
	for _, block := range f.Blocks {
		if vs, ok := block.symbolTable[variable]; ok {
			vs = utils.ReplaceSliceItem(vs, v, to)
			block.symbolTable[variable] = vs
		}
	}
}

func (b *Function) writeVariableByBlock(variable string, value Value, block *BasicBlock) {
	vs, ok := block.symbolTable[variable]
	if !ok {
		vs = make(Values, 0)
	}
	vs = append(vs, value)
	block.symbolTable[variable] = vs
}


// get value by variable and block
//
//	return : undefined \ value \ phi

// * first check builder.currentDef
//
// * if block sealed; just create a phi
// * if len(block.preds) == 0: undefined
// * if len(block.preds) == 1: just recursive
// * if len(block.preds) >  1: create phi and builder
func (b *FunctionBuilder) ReadVariable(variable string, create bool) Value {
	var ret Value
	b.ReadVariableEx(variable, create, func(vs []Value) {
		if len(vs) > 0 {
			ret = vs[len(vs)-1]
		} else {
			ret = nil
		}
	})
	return ret
}

func (b *FunctionBuilder) ReadVariableBefore(variable string, create bool, before Instruction) Value {
	var ret Value
	b.ReadVariableEx(variable, create, func(vs []Value) {
		for i := len(vs) - 1; i >= 0; i-- {
			vpos := vs[i].GetPosition()
			bpos := before.GetPosition()
			if vpos.StartLine <= bpos.StartLine {
				ret = vs[i]
				return
			}
		}
	})
	return ret
}

func (b *FunctionBuilder) ReadVariableEx(variable string, create bool, fun func([]Value)) {
	variable = b.GetScopeLocalVariable(variable)
	var ret []Value
	block := b.CurrentBlock
	if block == nil {
		block = b.ExitBlock
	}
	ret = b.readVariableByBlockEx(variable, block, create)
	fun(ret)
}

func (b *FunctionBuilder) deleteVariableByBlock(variable string, block *BasicBlock) {
	delete(block.symbolTable, variable)
}

func (b *FunctionBuilder) readVariableByBlock(variable string, block *BasicBlock, create bool) Value {
	ret := b.readVariableByBlockEx(variable, block, create)
	if len(ret) > 0 {
		return ret[len(ret)-1]
	} else {
		return nil
	}
}

func (b *FunctionBuilder) readVariableByBlockEx(variable string, block *BasicBlock, create bool) []Value {
	if vs, ok := block.symbolTable[variable]; ok && len(vs) > 0 {
		return vs
	}

	if block.Skip {
		return nil
	}

	var v Value
	// if block in sealedBlock
	if !block.isSealed {
		if create {
			phi := NewPhi(block, variable, create)
			phi.SetPosition(b.CurrentPos)
			block.inCompletePhi = append(block.inCompletePhi, phi)
			v = phi
		}
	} else if len(block.Preds) == 0 {
		// v = nil
		if create && b.CanBuildFreeValue(variable) {
			v = b.BuildFreeValue(variable)
		} else if i := b.TryBuildExternValue(variable); i != nil {
			v = i
		} else if create {
			un := NewUndefined(variable)
			// b.emitInstructionBefore(un, block.LastInst())
			b.EmitToBlock(un, block)
			v = un
		} else {
			v = nil
		}
	} else if len(block.Preds) == 1 {
		vs := b.readVariableByBlockEx(variable, block.Preds[0], create)
		if len(vs) > 0 {
			v = vs[len(vs)-1]
		} else {
			v = nil
		}
	} else {
		phi := NewPhi(block, variable, create)
		phi.SetPosition(b.CurrentPos)
		v = phi.Build()
	}
	if v != nil {
		b.writeVariableByBlock(variable, v, block)
		return []Value{v}
	} else {
		return nil
	}
}

// --------------- `f.freeValue`

func (b *FunctionBuilder) BuildFreeValue(variable string) Value {
	freeValue := NewParam(variable, true, b.Function)
	b.FreeValues = append(b.FreeValues, freeValue)
	b.WriteVariable(variable, freeValue)
	return freeValue
}

func (b *FunctionBuilder) CanBuildFreeValue(variable string) bool {
	parent := b.parentBuilder
	scope := b.parentScope
	block := b.parentCurrentBlock
	for parent != nil {
		variable = scope.GetLocalVariable(variable)
		v := parent.readVariableByBlock(variable, block, false)
		if v != nil && !v.IsExtern() {
			return true
		}

		// parent symbol and block
		scope = parent.parentScope
		block = parent.parentCurrentBlock
		// next parent
		parent = parent.parentBuilder
	}
	return false
}
