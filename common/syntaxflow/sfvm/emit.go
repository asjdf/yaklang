package sfvm

import (
	"fmt"
	"github.com/yaklang/yaklang/common/utils/omap"
)

func (y *SyntaxFlowVisitor) EmitNewRef(i string) {
	y.codes = append(y.codes, &SFI{
		OpCode:   OpNewRef,
		UnaryStr: i,
	})
}

func (y *SyntaxFlowVisitor) EmitUpdate(i string) {
	y.codes = append(y.codes, &SFI{
		OpCode:   OpUpdateRef,
		UnaryStr: i,
	})
}

func (y *SyntaxFlowVisitor) EmitWithdraw() {
	y.codes = append(y.codes, &SFI{
		OpCode: OpWithdraw,
	})
}

func (y *SyntaxFlowVisitor) EmitOperator(i string) {
	switch i {
	case ">":
		y.codes = append(y.codes, &SFI{OpCode: OpGt})
	case ">=":
		y.codes = append(y.codes, &SFI{OpCode: OpGtEq})
	case "<":
		y.codes = append(y.codes, &SFI{OpCode: OpLt})
	case "<=":
		y.codes = append(y.codes, &SFI{OpCode: OpLtEq})
	case "==", "=":
		y.codes = append(y.codes, &SFI{OpCode: OpEq})
	case "!=":
		y.codes = append(y.codes, &SFI{OpCode: OpNotEq})
	case "&&":
		y.codes = append(y.codes, &SFI{OpCode: OpLogicAnd})
	case "||":
		y.codes = append(y.codes, &SFI{OpCode: OpLogicOr})
	case "!":
		y.codes = append(y.codes, &SFI{OpCode: OpLogicBang})
	default:
		panic(fmt.Sprintf("unknown operator: %s", i))
	}
}

func (v *SyntaxFlowVisitor) EmitPushGlob(i string) {
	v.codes = append(v.codes, &SFI{
		OpCode:   OpGlobMatch,
		UnaryStr: i,
	})
}

func (y *SyntaxFlowVisitor) EmitRegexpMatch(i string) {
	y.codes = append(y.codes, &SFI{
		OpCode:   OpReMatch,
		UnaryStr: i,
	})
}

func (v *SyntaxFlowVisitor) EmitPushLiteral(i any) {
	switch ret := i.(type) {
	case string:
		v.codes = append(v.codes, &SFI{
			OpCode:   OpPushString,
			UnaryStr: ret,
		})
	case int:
		v.codes = append(v.codes, &SFI{
			OpCode:   OpPushNumber,
			UnaryInt: ret,
		})
	case bool:
		if ret {
			v.codes = append(v.codes, &SFI{
				OpCode:   OpPushBool,
				UnaryInt: 1,
			})
		} else {
			v.codes = append(v.codes, &SFI{
				OpCode:   OpPushString,
				UnaryInt: 0,
			})
		}
	default:
		panic(fmt.Sprintf("unknown type: %T", ret))
	}

}

func (v *SyntaxFlowVisitor) EmitEqual(i any) {
	switch i.(type) {
	case string:
	case int:
	}
}

func (v *SyntaxFlowVisitor) EmitSearchExact(i string) {
	v.codes = append(v.codes, &SFI{
		OpCode:   OpPushSearchExact,
		UnaryStr: i,
	})
}

func (v *SyntaxFlowVisitor) EmitSearchGlob(i string) {
	v.codes = append(v.codes, &SFI{
		OpCode:   OpPushSearchGlob,
		UnaryStr: i,
	})
}

func (v *SyntaxFlowVisitor) EmitSearchRegexp(i string) {
	v.codes = append(v.codes, &SFI{
		OpCode:   OpPushSearchRegexp,
		UnaryStr: i,
	})
}

func (v *SyntaxFlowVisitor) EmitGetMembers() {
	v.codes = append(v.codes, &SFI{
		OpCode: OpGetMembers,
	})
}

func (v *SyntaxFlowVisitor) EmitPushCallArgs(i int) {
	v.codes = append(v.codes, &SFI{OpCode: OpGetCallArgs, UnaryInt: i})
}

func (v *SyntaxFlowVisitor) Show() {
	for _, c := range v.codes {
		fmt.Println(c.String())
	}
}

func (v *SyntaxFlowVisitor) CreateFrame(vars *omap.OrderedMap[string, any]) *SFFrame {
	return NewSFFrame(vars, v.text, v.codes)
}

func (y *SyntaxFlowVisitor) EmitPop() {
	y.codes = append(y.codes, &SFI{
		OpCode: OpPop,
	})
}

func (y *SyntaxFlowVisitor) EmitCheckStackTop() {
	y.codes = append(y.codes, &SFI{
		OpCode: OpCheckStackTop,
	})
}

func (y *SyntaxFlowVisitor) EmitGetTopDef() {
	y.codes = append(y.codes, &SFI{
		OpCode: OpTopDefs,
	})
}

func (y *SyntaxFlowVisitor) EmitListIndex(i int) {
	y.codes = append(y.codes, &SFI{OpCode: OpListIndex, UnaryInt: i})
}

func (v *SyntaxFlowVisitor) EmitPass() {
	v.codes = append(v.codes, &SFI{
		OpCode: OpPass,
	})
}
