package ssa

func NewPhi(block *BasicBlock, variable string, create bool) *Phi {
	p := &Phi{
		anInstruction: NewInstruction(),
		anValue:       NewValue(),
		Edge:          make([]Value, 0, len(block.Preds)),
		create:        create,
	}
	p.SetVariable(variable)
	p.SetBlock(block)
	p.SetFunc(block.GetFunc())
	return p
}

func (b *BasicBlock) Sealed() {
	for _, p := range b.inCompletePhi {
		v := p.Build()
		v.SetPosition(p.GetPosition())
	}
	b.inCompletePhi = nil
	b.isSealed = true
}

func (p *Phi) AddEdge(v Value) {
	p.Edge = append(p.Edge, v)
}

func (phi *Phi) Name() string { return phi.GetVariable() }

func (phi *Phi) Build() Value {
	phi.GetBlock().Skip = true
	for _, predBlock := range phi.GetBlock().Preds {
		v := phi.GetFunc().builder.readVariableByBlock(phi.GetVariable(), predBlock, phi.create)
		phi.Edge = append(phi.Edge, v)
	}
	phi.GetBlock().Skip = false
	// phi.SetPosition(phi.GetBlock().GetPosition())
	v := phi.tryRemoveTrivialPhi()
	if v == phi {
		block := phi.GetBlock()
		block.Phis = append(block.Phis, phi)
		phi.GetFunc().SetReg(phi)
	}
	if v != nil {
		fixupUseChain(v)
	}
	return v
}

func (phi *Phi) tryRemoveTrivialPhi() Value {
	w1, w2 := phi.wit1, phi.wit2
	getValue := func(pass Value) Value {
		for _, v := range phi.Edge {
			if v == phi || v == pass {
				continue
			}
			return v
		}
		return nil
	}
	if w1 == nil || w2 == nil {
		// init w1 w2
		w1 = getValue(nil)
		w2 = getValue(w1)
	} else {
		if w1 == phi || w1 == w2 {
			w1 = getValue(w2)
		}
		if w2 == phi || w2 == w1 {
			w2 = getValue(w1)
		}
	}

	var ret Value
	ret = phi
	if w1 == nil {
		if w2 == nil {
			ret = nil
		} else {
			ret = w2
		}
	} else if w2 == nil {
		ret = w1
	}
	if ret != nil && ret != phi {
		phi.Replace(ret)
	}
	return ret
}

func (phi *Phi) Replace(to Value) {
	ReplaceValue(phi, to)
	for _, user := range phi.GetUsers() {
		switch p := user.(type) {
		case *Phi:
			p.tryRemoveTrivialPhi()
		}
	}
}
