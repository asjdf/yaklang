package ssaapi

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yak/ssa"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
)

// GetTopDefs desc all of 'Defs' is not used by any other value
func (i *Value) GetTopDefs(opt ...OperationOption) Values {
	actx := NewAnalyzeContext(opt...)
	actx.Self = i
	ret := i.getTopDefs(actx, opt...)
	ret = lo.UniqBy(ret, func(item *Value) int64 {
		return item.GetId()
	})
	return ret
}

func (v Values) GetTopDefs(opts ...OperationOption) Values {
	ret := make(Values, 0)
	for _, sub := range v {
		ret = append(ret, sub.GetTopDefs(opts...)...)
	}
	return ret
}

func (i *Value) visitedDefsDefault(actx *AnalyzeContext, opt ...OperationOption) Values {
	var vals Values
	if i.node == nil {
		return vals
	}
	if !actx.TheDefaultShouldBeVisited(i) {
		return vals
	}
	for _, def := range i.node.GetValues() {
		if ret := i.NewValue(def).AppendEffectOn(i).getTopDefs(actx, opt...); len(ret) > 0 {
			vals = append(vals, ret...)
		}
	}

	if len(vals) <= 0 {
		vals = append(vals, i)
	}
	if maskable, ok := i.node.(ssa.Maskable); ok {
		for _, def := range maskable.GetMask() {
			if ret := i.NewValue(def).AppendEffectOn(i).getTopDefs(actx, opt...); len(ret) > 0 {
				vals = append(vals, ret...)
			}
		}
	}
	return vals
}

func (i *Value) getTopDefs(actx *AnalyzeContext, opt ...OperationOption) Values {
	if i == nil {
		return nil
	}
	if actx == nil {
		actx = NewAnalyzeContext(opt...)
	}

	actx.EnterRecursive()
	// 1w recursive call check
	if !utils.InGithubActions() {
		if actx.GetRecursiveCounter() > 10000 {
			log.Warnf("recursive call is over 10000, stop it")
			return Values{i}
		}
	}
	if actx.IsReachedDepthLimited() {
		log.Warnf("reached depth limit,stop it")
		return Values{i}
	}

	actx.depth--
	defer func() {
		actx.depth++
	}()
	i.SetDepth(actx.depth)
	if actx.depth > 0 && actx.config.MaxDepth > 0 && actx.depth > actx.config.MaxDepth {
		actx.ReachDepthLimited()
		return Values{i}
	}
	if actx.depth < 0 && actx.config.MinDepth < 0 && actx.depth < actx.config.MinDepth {
		actx.ReachDepthLimited()
		return Values{i}
	}
	if actx.depth > 0 && i.GetDepth() > 0 && actx.depth > i.GetDepth() {
		actx.ReachDepthLimited()
		return Values{i}
	}
	if actx.depth < 0 && i.GetDepth() < 0 && actx.depth < i.GetDepth() {
		actx.ReachDepthLimited()
		return Values{i}
	}

	// hook everynode
	if len(actx.config.HookEveryNode) > 0 {
		for _, hook := range actx.config.HookEveryNode {
			if err := hook(i); err != nil {
				if err.Error() != "abort" {
					log.Errorf("hook-every-node error: %v", err)
				}
				return Values{}
			}
		}
	}

	{
		obj, key, member := actx.GetCurrentObject()
		_ = obj
		_ = key
		_ = member
		if obj != nil && i.IsObject() && i != obj {
			if m := i.GetMember(key); m != nil && m != member {
				actx.PopObject()
				return m.getTopDefs(actx, opt...)
			}
		}
	}

	getMemberCall := func(apiValue *Value, value ssa.Value, actx *AnalyzeContext) Values {
		if value.HasValues() {
			return i.visitedDefsDefault(actx, opt...)
		}
		if value.IsMember() {
			obj := i.NewValue(value.GetObject())
			key := i.NewValue(value.GetKey())
			if err := actx.PushObject(obj, key, i); err != nil {
				log.Errorf("%v", err)
				return i.visitedDefsDefault(actx, opt...)
			}
			obj.AppendDependOn(apiValue)
			ret := obj.getTopDefs(actx, opt...)
			if !ValueCompare(i, actx.Self) {
				ret = append(ret, i)
			}
			return ret
		}
		return i.visitedDefsDefault(actx, opt...)
	}

	switch inst := i.node.(type) {
	case *ssa.LazyInstruction:
		var ok bool
		i.node, ok = inst.Self().(ssa.Value)
		if !ok {
			log.Errorf("BUG: %T is not ssa.Value", inst.Self())
			return Values{}
		}
		return i.getTopDefs(actx, opt...)
	case *ssa.Undefined:
		// ret[n]
		return getMemberCall(i, inst, actx)
	case *ssa.ConstInst:
		return i.visitedDefsDefault(actx, opt...)
	case *ssa.Phi:
		if !actx.ThePhiShouldBeVisited(i) {
			return Values{}
		}

		conds := inst.GetControlFlowConditions()
		result := getMemberCall(i, inst, actx)
		for _, cond := range conds {
			v := i.NewValue(cond)
			ret := v.AppendEffectOn(i).getTopDefs(actx, opt...)
			result = append(result, v)
			result = append(result, ret...)
		}
		return result
	case *ssa.Call:
		calleeInst := inst.Method
		if calleeInst == nil {
			return Values{i} // return self
		}

		fun, isFunc := ssa.ToFunction(calleeInst)
		// callee := i.NewValue(fun)
		if !isFunc && calleeInst.GetReference() != nil {
			fun, isFunc = ssa.ToFunction(calleeInst.GetReference())
			// callee = i.NewValue(fun)
		}

		switch {
		case isFunc && !fun.IsExtern():
			callee := i.NewValue(fun)
			callee.SetContextValue(ANALYZE_RUNTIME_CTX_TOPDEF_CALL_ENTRY, i)
			callee.AppendEffectOn(i)
			needRecover := actx.CrossProcess(i, callee)
			if needRecover {
				defer actx.RecoverCrossProcess()
			}
			// inherit return index
			val, ok := i.GetContextValue(ANALYZE_RUNTIME_CTX_TOPDEF_CALL_ENTRY_TRACE_INDEX)
			if ok {
				callee.SetContextValue(ANALYZE_RUNTIME_CTX_TOPDEF_CALL_ENTRY_TRACE_INDEX, val)
			}
			return callee.getTopDefs(actx, opt...).AppendEffectOn(callee)
		default:
			callee := i.NewValue(calleeInst)
			i.AppendDependOn(callee)
			nodes := Values{callee}
			for _, val := range inst.Args {
				arg := i.NewValue(val)
				i.AppendDependOn(arg)
				nodes = append(nodes, arg)
			}
			for _, value := range inst.Binding {
				arg := i.NewValue(value)
				i.AppendDependOn(arg)
				nodes = append(nodes, arg)
			}
			var results Values
			for _, subNode := range nodes {
				if subNode == nil {
					continue
				}
				// getTopDefs(nil,opt...)第一个参数指定为nil
				// 提供一个新的上下文，避免指向新的actx.self导致多余的结果
				vals := subNode.GetTopDefs(opt...).AppendEffectOn(subNode)
				//vals := subNode.getTopDefs(nil, opt...).AppendEffectOn(subNode)
				results = append(results, vals...)
			}
			return results
		}
	case *ssa.Function:
		var vals Values
		// handle return
		returnIndex, traceIndexedReturn := i.GetContextValue(ANALYZE_RUNTIME_CTX_TOPDEF_CALL_ENTRY_TRACE_INDEX)
		if traceIndexedReturn {
			retIndexRaw := returnIndex.GetConstValue()
			retIndexRawStr := fmt.Sprint(retIndexRaw)
			if utils.IsValidInteger(retIndexRawStr) {
				targetIdx := codec.Atoi(retIndexRawStr)
				var traceRets Values
				for _, retInsRaw := range inst.Return {
					retIns, ok := ssa.ToReturn(retInsRaw)
					if !ok {
						log.Warnf("BUG: %T is not *Return", retInsRaw)
						continue
					}
					for idx, traceVal := range retIns.Results {
						if idx == targetIdx {
							traceRets = append(traceRets, i.NewValue(traceVal).AppendEffectOn(i))
						}
					}
				}
				return lo.FlatMap(traceRets, func(item *Value, index int) []*Value {
					return item.getTopDefs(actx, opt...)
				})
			} else {
				// string literal member
				var traceRets Values
				for _, retInsRaw := range inst.Return {
					retIns, ok := ssa.ToReturn(retInsRaw)
					if !ok {
						log.Warnf("BUG: %T is not *Return", retInsRaw)
						continue
					}
					for _, traceVal := range retIns.Results {
						val, ok := traceVal.GetStringMember(retIndexRawStr)
						if ok {
							traceRets = append(traceRets, i.NewValue(val).AppendEffectOn(i))
							// trace mask ?
							if len(inst.Blocks) > 0 {
								name, ok := ssa.CombineMemberCallVariableName(traceVal, ssa.NewConst(retIndexRawStr))
								if ok {
									lastBlockRaw, _ := lo.Last(inst.Blocks)
									lastBlock, ok := ssa.ToBasicBlock(lastBlockRaw)
									if ok {
										variableInstance := lastBlock.ScopeTable.ReadVariable(name)
										_ = variableInstance.String()
									}
								}
							}
						}
					}
				}
				return lo.FlatMap(traceRets, func(item *Value, index int) []*Value {
					return item.getTopDefs(actx, opt...)
				})
			}
		}

		handlerReturn := func(value *Value) {
			fun, ok := ssa.ToFunction(value.node)
			if !ok {
				return
			}
			for _, r := range fun.Return {
				for _, subVal := range r.GetValues() {
					if ret := value.NewValue(subVal).AppendEffectOn(value).getTopDefs(actx, opt...); len(ret) > 0 {
						vals = append(vals, ret...)
					}
				}
			}
		}

		handlerReturn(i)
		if len(vals) == 0 {
			vals = append(vals, i)
		}
		// handler child-class function
		for _, child := range inst.GetPointer() {
			handlerReturn(i.NewValue(child))
		}
		return vals.AppendEffectOn(i)
	case *ssa.ParameterMember:
		var vals Values
		getParameter := func() Values {
			log.Infof("formal parameter index: %d is out of range", inst.FormalParameterIndex)
			fun := i.GetFunction()
			switch inst.MemberCallKind {
			case ssa.ParameterMemberCall:
				if para := fun.GetParameter(inst.MemberCallObjectIndex); para != nil {
					return para.getTopDefs(actx, opt...)
				}
			case ssa.FreeValueMemberCall:
				if fv := fun.GetFreeValue(inst.MemberCallObjectName); fv != nil {
					return fv.getTopDefs(actx, opt...)
				}
			}
			return Values{i}
		}
		getCalledByValue := func(called *Value) Values {
			if called == nil {
				return nil
			}
			calledInstance, ok := ssa.ToCall(called.node)
			if !ok {
				log.Infof("BUG: Parameter getCalledByValue called is not callInstruction %s", called.GetOpcode())
				return Values{}
			}
			var actualParam ssa.Value
			if inst.FormalParameterIndex >= len(calledInstance.ArgMember) {
				return getParameter()
			}
			actualParam = calledInstance.ArgMember[inst.FormalParameterIndex]
			traced := i.NewValue(actualParam).AppendEffectOn(called)
			ret := traced.getTopDefs(actx, opt...)
			if len(ret) > 0 {
				return ret
			} else {
				return Values{traced}
			}
		}
		called := actx.GetCallFromLastCrossProcess()
		if called != nil {
			hash, needRecover := actx.ReverseProcessWithDirection(i, called)
			calledByValue := getCalledByValue(called)
			if needRecover {
				actx.RecoverReverseProcess(hash)
			}
			vals = append(vals, calledByValue...)
		}
		if actx.config.AllowIgnoreCallStack {
			if fun := i.GetFunction(); fun != nil {
				fun.GetCalledBy().ForEach(func(value *Value) {
					hash, needRecover := actx.ReverseProcessWithDirection(i, called)
					val := getCalledByValue(called)
					if needRecover {
						actx.RecoverReverseProcess(hash)
					}
					vals = append(vals, val...)
				})
			}
		}
		if len(vals) == 0 {
			return getParameter()
		}
		return vals.AppendEffectOn(i)
	case *ssa.Parameter:
		// 查找被调用函数的TopDef
		getCalledByValue := func(called *Value) Values {
			if called == nil {
				return nil
			}
			calledInstance, ok := ssa.ToCall(called.node)
			if !ok {
				log.Infof("BUG: Parameter getCalledByValue called is not callInstruction %s", called.GetOpcode())
				return Values{}
			}

			thisFunc := i.GetFunction()
			if !ValueCompare(i.NewValue(calledInstance.Method), thisFunc) {
				log.Errorf("call stack function %s(%d) not same with Parameter function %s(%d)",
					calledInstance.Method.GetName(), calledInstance.Method.GetId(),
					thisFunc.GetName(), thisFunc.GetId(),
				)
				return Values{}
			}

			var actualParam ssa.Value
			if inst.IsFreeValue {
				// free value
				if tmp, ok := calledInstance.Binding[inst.GetName()]; ok {
					actualParam = tmp
				} else {
					log.Errorf("free value: %v is not found in binding", inst.GetName())
					return getMemberCall(i, i.node, actx)
				}
			} else {
				// parameter
				if inst.FormalParameterIndex >= len(calledInstance.Args) {
					log.Infof("formal parameter index: %d is out of range", inst.FormalParameterIndex)
					return getMemberCall(i, i.node, actx)
				}
				actualParam = calledInstance.Args[inst.FormalParameterIndex]
			}

			traced := i.NewValue(actualParam).AppendEffectOn(called)
			if !actx.TheParameterShouldBeVisited(traced) {
				return Values{traced}
			}
			ret := traced.getTopDefs(actx, opt...)

			if len(ret) > 0 {
				return ret
			} else {
				return Values{traced}
			}
		}

		if inst.GetDefault() != nil {
			return i.NewValue(inst.GetDefault()).getTopDefs(actx, opt...)
		}
		var vals Values
		called := actx.GetCallFromLastCrossProcess()
		if called != nil {
			if !called.IsCall() {
				log.Infof("parent function is not called by any other function, skip (%T)", called)
				return Values{i}
			}
			hash, needRecover := actx.ReverseProcess()
			calledByValue := getCalledByValue(called)
			if needRecover {
				actx.RecoverReverseProcess(hash)
			}
			vals = append(vals, calledByValue...)
		}

		// if not found in call stack, then find in called-by
		if actx.config.AllowIgnoreCallStack && len(vals) == 0 {
			fun := i.GetFunction()
			if fun != nil {
				call2fun := fun.GetCalledBy()
				call2fun.ForEach(func(call *Value) {
					hash, needRecover := actx.ReverseProcessWithDirection(i, call)
					val := getCalledByValue(call)
					if needRecover {
						actx.RecoverReverseProcess(hash)
					}
					vals = append(vals, val...)
				})
			}
		}

		if len(vals) == 0 {
			vals = append(vals, i)
		}
		return vals.AppendEffectOn(i)
	case *ssa.SideEffect:
		callIns := inst.CallSite
		if callIns != nil {
			call := i.NewValue(callIns).AppendEffectOn(i)
			needRecover := actx.CrossProcess(i, call)
			if needRecover {
				defer actx.RecoverCrossProcess()
			}
			v := i.NewValue(inst.Value).AppendEffectOn(i)
			return v.getTopDefs(actx)
		} else {
			log.Errorf("side effect: %v is not created from call instruction", i.String())
		}
	case *ssa.Make:
		var values Values
		values = append(values, i)
		for key, member := range inst.GetAllMember() {
			// if key.String() == "__ref__" {
			// 	continue
			// }
			value := i.NewValue(member)
			if err := actx.PushObject(i, i.NewValue(key), value); err != nil {
				log.Errorf("push object failed: %v", err)
				continue
			}
			values = append(values, value.getTopDefs(actx, opt...)...)
			actx.PopObject()
		}
		return values
	}
	return getMemberCall(i, i.node, actx)
}
