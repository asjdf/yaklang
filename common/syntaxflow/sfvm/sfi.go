package sfvm

import (
	"fmt"
	"strconv"
)

type SFVMOpCode int

const (
	OpPass SFVMOpCode = iota

	// OpPushNumber and OpPushString and OpPushBool can push literal into stack
	OpPushNumber
	OpPushString
	OpPushBool
	OpPop
	OpWithdraw

	// OpPushSearchExact can push data from origin
	OpPushSearchExact
	OpPushSearchGlob
	OpPushSearchRegexp

	// handle function call
	OpGetCallArgs
	OpGetMembers

	OpTopDefs
	OpBottomUse

	// ListOperation
	OpListIndex

	// => variable
	OpNewRef
	OpUpdateRef

	/*
		Binary Operator
		Fetch TWO in STACK, calc result, push result into stack
	*/
	OpEq
	OpNotEq
	OpGt
	OpGtEq
	OpLt
	OpLtEq
	OpLogicAnd
	OpLogicOr
	OpLogicBang

	/*
		Unary Operator: Fetch ONE in STACK, calc result, push result into stack
	*/
	OpReMatch
	OpGlobMatch
	OpNot

	OpCheckStackTop
)

type SFI struct {
	OpCode   SFVMOpCode
	UnaryInt int
	UnaryStr string
	Desc     string
	Values   []string
}

const verboseLen = "%-12s"

func (s *SFI) String() string {
	switch s.OpCode {
	case OpPass:
		return "- pass -"
	case OpPushBool:
		v := "false"
		if s.UnaryInt > 0 {
			v = "true"
		}
		return fmt.Sprintf(verboseLen+" %v", "push", v)
	case OpPushString:
		return fmt.Sprintf(verboseLen+" (len:%v) %v", "push", len(s.UnaryStr), strconv.Quote(s.UnaryStr))
	case OpPushNumber:
		return fmt.Sprintf(verboseLen+" %v", "push", s.UnaryInt)
	case OpPushSearchGlob:
		return fmt.Sprintf(verboseLen+" %v", "push$glob", s.UnaryStr)
	case OpPushSearchExact:
		return fmt.Sprintf(verboseLen+" %v", "push$exact", s.UnaryStr)
	case OpPushSearchRegexp:
		return fmt.Sprintf(verboseLen+" %v", "push$regexp", s.UnaryStr)
	case OpNewRef:
		return fmt.Sprintf(verboseLen+" %v", "new$ref", s.UnaryStr)
	case OpUpdateRef:
		return fmt.Sprintf(verboseLen+" %v", "update$ref", s.UnaryStr)
	case OpEq:
		return fmt.Sprintf(verboseLen+" %v", "(operator) ==", s.UnaryStr)
	case OpNotEq:
		return fmt.Sprintf(verboseLen+" %v", "(operator) !=", s.UnaryStr)
	case OpGt:
		return fmt.Sprintf(verboseLen+" %v", "(operator) >", s.UnaryStr)
	case OpGtEq:
		return fmt.Sprintf(verboseLen+" %v", "(operator) >=", s.UnaryStr)
	case OpLt:
		return fmt.Sprintf(verboseLen+" %v", "(operator) <", s.UnaryStr)
	case OpLtEq:
		return fmt.Sprintf(verboseLen+" %v", "(operator) <=", s.UnaryStr)
	case OpReMatch:
		return fmt.Sprintf(verboseLen+" %v", "(operator) ~=", s.UnaryStr)
	case OpGlobMatch:
		return fmt.Sprintf(verboseLen+" %v", "(operator) *~", s.UnaryStr)
	case OpNot:
		return fmt.Sprintf(verboseLen+" %v", "(operator) !", s.UnaryStr)
	case OpLogicAnd:
		return fmt.Sprintf(verboseLen+" %v", "(operator) &&", s.UnaryStr)
	case OpLogicOr:
		return fmt.Sprintf(verboseLen+" %v", "(operator) ||", s.UnaryStr)
	case OpPop:
		return fmt.Sprintf(verboseLen+" %v", "pop", s.UnaryStr)
	case OpWithdraw:
		return fmt.Sprint("withdraw last stack value")
	case OpCheckStackTop:
		return fmt.Sprint("check stack top")
	default:
		panic("unhandled default case")
	}
	return ""
}
