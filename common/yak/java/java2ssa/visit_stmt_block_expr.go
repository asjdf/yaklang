package java2ssa

import (
	"github.com/google/uuid"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
	javaparser "github.com/yaklang/yaklang/common/yak/java/parser"
	"github.com/yaklang/yaklang/common/yak/ssa"
	"github.com/yaklang/yaklang/common/yak/yak2ssa"
)

func (y *builder) VisitBlock(raw javaparser.IBlockContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.BlockContext)
	if i == nil {
		return nil
	}
	block := i.GetText()
	log.Info("block", block)
	if ret := i.BlockStatementList(); ret != nil {
		y.BuildSyntaxBlock(func() {
			y.VisitBlockStatementList(ret)
		})
	}

	return nil
}

func (y *builder) VisitBlockStatement(raw javaparser.IBlockStatementContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.BlockStatementContext)
	if i == nil {
		return nil
	}

	if ret := i.LocalVariableDeclaration(); ret != nil {
		y.VisitLocalVariableDeclaration(ret)
	} else if ret := i.LocalTypeDeclaration(); ret != nil {
		y.VisitLocalTypeDeclaration(ret)
	} else if ret := i.Statement(); ret != nil {
		y.VisitStatement(ret)
	}

	return nil
}

func (y *builder) VisitExpression(raw javaparser.IExpressionContext) ssa.Value {
	if y == nil || raw == nil {
		return nil
	}

	var opcode ssa.BinaryOpcode
	var unaryOpcode ssa.UnaryOpcode
	var variable *ssa.Variable
	var value ssa.Value
	var handlerJumpExpression = func(cond func(string) ssa.Value, trueExpr, falseExpr func() ssa.Value) ssa.Value {
		// 为了聚合产生Phi指令
		id := uuid.NewString()
		variable := y.CreateVariable(id)
		y.AssignVariable(variable, y.EmitConstInstAny())
		// 只需要使用b.WriteValue设置value到此ID，并最后调用b.ReadValue可聚合产生Phi指令，完成语句预期行为
		ifb := y.CreateIfBuilder()
		ifb.AppendItem(
			func() ssa.Value {
				return cond(id)
			},
			func() {
				v := trueExpr()
				variable := y.CreateVariable(id)
				y.AssignVariable(variable, v)
			},
		)
		ifb.SetElse(func() {
			v := falseExpr()
			variable := y.CreateVariable(id)
			y.AssignVariable(variable, v)
		})
		ifb.Build()
		// generator phi instruction
		v := y.ReadValue(id)
		v.SetName(raw.GetText())
		return v
	}

	switch ret := raw.(type) {
	case *javaparser.PrimaryExpressionContext:
		// 处理主要表达式
		if ret.Primary() != nil {
			return y.VisitPrimary(ret.Primary())
		}
		return nil
	case *javaparser.SliceCallExpressionContext:
		// 处理切片调用表达式
		expr := y.VisitExpression(ret.Expression(0))
		key := y.VisitExpression(ret.Expression(1))
		if utils.IsNil(key) {
			y.NewError(ssa.Error, "javaast", yak2ssa.AssignRightSideEmpty())
			return nil
		}
		return y.ReadMemberCallVariable(expr, key)
	case *javaparser.MemberCallExpressionContext:
		// 处理成员调用表达式，如通过点操作符访问成员
		// todo: 成员调用表达式
		value = y.ReadValue(ret.GetText())
		return value
	case *javaparser.FunctionCallExpressionContext:
		// 处理函数调用表达式
		if s := ret.MethodCall(); s != nil {
			return y.VisitMethodCall(s)
		}
		return nil
	case *javaparser.MethodReferenceExpressionContext:
		// 处理方法引用表达式
		// todo: 方法引用表达式
		return nil
	case *javaparser.ConstructorReferenceExpressionContext:
		// 处理构造器引用表达式
		// todo: 构造器引用表达式
		return nil
	case *javaparser.Java17SwitchExpressionContext:
		// 处理 Java 17 的 switch 表达式
		value := y.VisitSwitchExpression(ret.SwitchExpression(), true)
		return value
	case *javaparser.PostfixExpression1Context:
		// 处理后缀表达式，如自增、自减操作
		if s := ret.GetLeftExpression(); s != nil {
			text := s.GetText()
			variable = y.CreateVariable(text)
		}

		if variable == nil {
			y.NewError(ssa.Error, "javaast", yak2ssa.AssignLeftSideEmpty())
			return nil
		}

		if postfix := ret.GetPostfix().GetText(); postfix == "++" {
			value = y.EmitBinOp(ssa.OpAdd, y.ReadValueByVariable(variable), y.EmitConstInst(1))
		} else if postfix == "--" {
			value = y.EmitBinOp(ssa.OpSub, y.ReadValueByVariable(variable), y.EmitConstInst(1))
		}

		y.AssignVariable(variable, value)
		return value

	case *javaparser.PostfixExpression2Context:
		if s := ret.Identifier(); s != nil {
			text := s.GetText()
			variable = y.CreateVariable(text)
		}
		if variable == nil {
			y.NewError(ssa.Error, "javaast", yak2ssa.AssignLeftSideEmpty())
			return nil
		}

		if postfix := ret.GetPostfix().GetText(); postfix == "++" {
			value = y.EmitBinOp(ssa.OpAdd, y.ReadValueByVariable(variable), y.EmitConstInst(1))
		} else if postfix == "--" {
			value = y.EmitBinOp(ssa.OpSub, y.ReadValueByVariable(variable), y.EmitConstInst(1))
		}

		y.AssignVariable(variable, value)
		return value

	case *javaparser.PrefixUnaryExpressionContext:
		// 处理前缀表达式，如正负号、逻辑非等
		if ret.Expression() != nil {
			value = y.VisitExpression(ret.Expression())
		} else {
			y.NewError(ssa.Error, "javaast", yak2ssa.AssignRightSideEmpty())
		}
		switch ret.GetPrefix().GetText() {
		case "+":
			unaryOpcode = ssa.OpPlus
		case "-":
			unaryOpcode = ssa.OpNeg
		case "~":
			unaryOpcode = ssa.OpBitwiseNot
		case "!":
			unaryOpcode = ssa.OpNot
		default:
			y.NewError(ssa.Error, "javaast", yak2ssa.UnaryOperatorNotSupport(ret.GetText()))
		}
		return y.EmitUnOp(unaryOpcode, value)
	case *javaparser.PrefixBinayExpression1Context:
		// 处理前缀表达式中的"--"和"++"
		if s := ret.GetLeftExpression(); s != nil {
			text := s.GetText()
			variable = y.CreateVariable(text)
		}
		if variable == nil {
			y.NewError(ssa.Error, "javaast", yak2ssa.AssignLeftSideEmpty())
			return nil
		}

		value = y.ReadValueByVariable(variable)
		if prefix := ret.GetPrefix().GetText(); prefix == "++" {
			y.AssignVariable(variable, y.EmitBinOp(ssa.OpAdd, value, y.EmitConstInst(1)))
		} else if prefix == "--" {
			y.AssignVariable(variable, y.EmitBinOp(ssa.OpSub, value, y.EmitConstInst(1)))
		}
		return value

	case *javaparser.PrefixBinayExpression2Context:
		if s := ret.Identifier(); s != nil {
			text := s.GetText()
			variable = y.CreateVariable(text)
		}
		if variable == nil {
			y.NewError(ssa.Error, "javaast", yak2ssa.AssignLeftSideEmpty())
			return nil
		}

		value = y.ReadValueByVariable(variable)
		if prefix := ret.GetPrefix().GetText(); prefix == "++" {
			y.AssignVariable(variable, y.EmitBinOp(ssa.OpAdd, value, y.EmitConstInst(1)))
		} else if prefix == "--" {
			y.AssignVariable(variable, y.EmitBinOp(ssa.OpSub, value, y.EmitConstInst(1)))
		}
		return value
	case *javaparser.CastExpressionContext:
		// 处理类型转换表达式
		// todo: 类型转换表达式
		return nil
	case *javaparser.NewCreatorExpressionContext:
		// 处理创建对象的表达式
		// todo: 创建对象的表达式
		y.VisitCreator(ret.Creator())
		return nil
	case *javaparser.MultiplicativeExpressionContext:
		// 处理乘法、除法、模运算表达式
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		switch ret.GetBop().GetText() {
		case "*":
			opcode = ssa.OpMul
		case "/":
			opcode = ssa.OpDiv
		case "%":
			opcode = ssa.OpMod
		default:
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		return y.EmitBinOp(opcode, op1, op2)
	case *javaparser.AdditiveExpressionContext:
		// 处理加法和减法表达式
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		switch ret.GetBop().GetText() {
		case "+":
			opcode = ssa.OpAdd
		case "-":
			opcode = ssa.OpSub
		default:
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		return y.EmitBinOp(opcode, op1, op2)
	case *javaparser.ShiftExpressionContext:
		// 处理位移表达式
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		var ltNum int
		var rtNum int

		for index, _ := range ret.AllLT() {
			_ = index
			ltNum++
		}
		for index, _ := range ret.AllGT() {
			_ = index
			rtNum++
		}

		if ltNum != 0 && rtNum != 0 {
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		if ltNum == 2 {
			opcode = ssa.OpShl
		} else if rtNum == 2 {
			opcode = ssa.OpShr
		} else if rtNum == 3 {
			//todo: 无符号右移运算符
			opcode = ssa.OpShr
		} else {
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		return y.EmitBinOp(opcode, op1, op2)
	case *javaparser.RelationalExpressionContext:
		// 处理关系运算表达式，如大于、小于等
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		switch ret.GetBop().GetText() {
		case "<":
			opcode = ssa.OpLt
		case ">":
			opcode = ssa.OpGt
		case "<=":
			opcode = ssa.OpLtEq
		case ">=":
			opcode = ssa.OpGtEq
		default:
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		return y.EmitBinOp(opcode, op1, op2)

	case *javaparser.InstanceofExpressionContext:
		// 处理 instanceof 表达式
		// todo instanceof 表达式
	case *javaparser.EqualityExpressionContext:
		// 处理等于和不等于表达式
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		switch ret.GetBop().GetText() {
		case "==":
			opcode = ssa.OpEq
		case "!=":
			opcode = ssa.OpNotEq
		default:
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		return y.EmitBinOp(opcode, op1, op2)
	case *javaparser.BitwiseAndExpressionContext:
		// 处理按位与表达式
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		if bop := ret.GetBop().GetText(); bop == "&" {
			opcode = ssa.OpAnd
		} else {
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		return y.EmitBinOp(opcode, op1, op2)
	case *javaparser.BitwiseXORExpressionContext:
		// 处理按位异或表达式
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		if bop := ret.GetBop().GetText(); bop == "^" {
			opcode = ssa.OpXor
		} else {
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		return y.EmitBinOp(opcode, op1, op2)
	case *javaparser.BitwiseORExpressionContext:
		// 处理按位或表达式
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		if bop := ret.GetBop().GetText(); bop == "|" {
			opcode = ssa.OpOr
		} else {
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		return y.EmitBinOp(opcode, op1, op2)
	case *javaparser.LogicANDExpressionContext:
		// 处理逻辑与表达式

		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		return handlerJumpExpression(
			func(id string) ssa.Value {
				return op1
			},
			func() ssa.Value {
				return op2
			},
			func() ssa.Value {
				return op1
			},
		)
	case *javaparser.LogicORExpressionContext:
		// 处理逻辑或表达式
		op1 := y.VisitExpression(ret.Expression(0))
		op2 := y.VisitExpression(ret.Expression(1))
		return handlerJumpExpression(
			func(id string) ssa.Value {
				return op1
			},
			func() ssa.Value {
				return op1
			},
			func() ssa.Value {
				return op2
			},
		)
	case *javaparser.TernaryExpressionContext:
		// 处理三元运算符表达式
		builder := y.CreateIfBuilder()
		allExpr := ret.AllExpression()
		if allExpr != nil {
			builder.AppendItem(func() ssa.Value {
				return y.VisitExpression(ret.Expression(0))
			},
				func() { y.VisitExpression(ret.Expression(1)) })
			builder.SetElse(func() { y.VisitExpression(ret.Expression(2)) })
			builder.Build()
		}
		if y.VisitExpression(ret.Expression(0)) == y.EmitConstInst(true) {
			return y.VisitExpression(ret.Expression(1))
		} else {
			return y.VisitExpression(ret.Expression(2))
		}
	case *javaparser.AssignmentExpression1Context:
		// 处理赋值表达式，包括所有赋值运算符
		variable = y.CreateVariable(ret.GetLeftExpression().GetText())
		if variable == nil {
			y.NewError(ssa.Error, "javaast", yak2ssa.AssignLeftSideEmpty())
			return nil
		}
		v := y.VisitExpression(ret.Expression(1))
		switch ret.GetBop().GetText() {
		case "+=":
			opcode = ssa.OpAdd
		case "-=":
			opcode = ssa.OpSub
		case "*=":
			opcode = ssa.OpMul
		case "/=":
			opcode = ssa.OpDiv
		case "%=":
			opcode = ssa.OpMod
		case "<<=":
			opcode = ssa.OpShl
		case ">>=":
			opcode = ssa.OpShr
		case ">>>=":
			opcode = ssa.OpShr //todo: 无符号右移运算符
		case "&=":
			opcode = ssa.OpAnd
		case "|=":
			opcode = ssa.OpOr
		case "^=":
			opcode = ssa.OpXor
		default:
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		value := y.EmitBinOp(opcode, y.ReadValueByVariable(variable), v)
		y.AssignVariable(variable, value)
		return value

	case *javaparser.AssignmentExpression2Context:
		// 处理赋值表达式，包括所有赋值运算符
		variable = y.CreateVariable(ret.Identifier().GetText())
		if variable == nil {
			y.NewError(ssa.Error, "javaast", yak2ssa.AssignLeftSideEmpty())
			return nil
		}
		v := y.VisitExpression(ret.Expression())
		switch ret.GetBop().GetText() {
		case "+=":
			opcode = ssa.OpAdd
		case "-=":
			opcode = ssa.OpSub
		case "*=":
			opcode = ssa.OpMul
		case "/=":
			opcode = ssa.OpDiv
		case "%=":
			opcode = ssa.OpMod
		case "<<=":
			opcode = ssa.OpShl
		case ">>=":
			opcode = ssa.OpShr
		case ">>>=":
			opcode = ssa.OpShr //todo: 无符号右移运算符
		case "&=":
			opcode = ssa.OpAnd
		case "|=":
			opcode = ssa.OpOr
		case "^=":
			opcode = ssa.OpXor
		default:
			y.NewError(ssa.Error, "javaast", yak2ssa.BinaryOperatorNotSupport(ret.GetText()))
			return nil
		}
		value := y.EmitBinOp(opcode, y.ReadValueByVariable(variable), v)
		y.AssignVariable(variable, value)
		return value

	case *javaparser.AssignmentEqExpression1Context:
		// 处理赋值表达式的等于号
		if s := ret.GetLeftExpression(); s != nil {
			variable = y.CreateVariable(s.GetText())
		}
		if id := ret.Identifier(); id != nil {
			value = y.ReadValue(id.GetText())
		} else if expr := ret.Expression(1); expr != nil {
			value = y.VisitExpression(expr)
		}
		y.AssignVariable(variable, value)
		return nil

	case *javaparser.AssignmentEqExpression2Context:
		// 处理赋值表达式的等于号
		if s := ret.Identifier(0); s != nil {
			variable = y.CreateVariable(s.GetText())
		}
		if id := ret.Identifier(1); id != nil {
			value = y.ReadValue(id.GetText())
		} else if expr := ret.Expression(); expr != nil {
			value = y.VisitExpression(expr)
		}
		y.AssignVariable(variable, value)
		return nil
	case *javaparser.Java8LambdaExpressionContext:
		// 处理 Java 8 的 lambda 表达式
		// todo: Java 8 的 lambda 表达式
		return y.EmitConstInstNil()
	default:
		// 默认情况，可能是不支持的表达式类型
		log.Errorf("unsupported expression type: %T", ret)
		panic("unsupported expression type")
	}

	return y.EmitUndefined("_")
}

func (y *builder) VisitMethodCall(raw javaparser.IMethodCallContext) ssa.Value {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.MethodCallContext)
	if i == nil {
		return nil
	}

	var v ssa.Value
	if ret := i.Identifier(); ret != nil {
		v = y.ReadValue(ret.GetText())
	} else if ret := i.THIS(); ret != nil {
		v = y.ReadValue(ret.GetText())
	} else if ret = i.SUPER(); ret != nil {
		v = y.ReadValue(ret.GetText())
	}
	var args []ssa.Value
	if argument := i.Arguments(); argument != nil {
		args = y.VisitArguments(i.Arguments())
		c := y.NewCall(v, args)
		return y.EmitCall(c)
	}
	return nil
}

func (y *builder) VisitPrimary(raw javaparser.IPrimaryContext) ssa.Value {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.PrimaryContext)
	if i == nil {
		return nil
	}

	if ret := i.Expression(); ret != nil {
		return y.VisitExpression(ret)
	}
	if ret := i.Literal(); ret != nil {
		return y.VisitLiteral(ret)
	}

	if ret := i.Identifier(); ret != nil {
		text := ret.GetText()
		if text == "_" {
			y.NewError(ssa.Warn, "javaast", "cannot use _ as value")
			return nil
		}
		v := y.ReadValue(text)
		return v
	}
	return nil
}

func (y *builder) VisitStatement(raw javaparser.IStatementContext) interface{} {
	if y.IsBlockFinish() {
		return nil
	}
	if y == nil || raw == nil {
		return nil
	}

	switch ret := raw.(type) {
	case *javaparser.BlockLabelStatementContext:
		return y.VisitBlock(ret.Block())
	case *javaparser.AssertStatementContext:
		// 处理断言语句
		getExpr := func(i int) ssa.Value {
			if expr := ret.Expression(i); expr != nil {
				return y.VisitExpression(expr)
			}
			y.NewError(ssa.Error, "javaast", yak2ssa.UnexpectedAssertStmt())
			return nil
		}
		exprs := ret.AllExpression()
		lenExprs := len(exprs)

		var cond, msgV ssa.Value

		cond = getExpr(0)
		if lenExprs > 1 {
			msgV = getExpr(1)
		}
		y.EmitAssert(cond, msgV, exprs[0].GetText())
		return nil
	case *javaparser.IfStatementContext:
		// 处理 if 语句
		y.VisitIfStmt(ret.Ifstmt())
		return nil
	case *javaparser.ForStatementContext:
		//处理 for 语句
		if ret.ForControl() != nil {
			loop := y.VisitForControl(ret.ForControl())
			//设置循环体
			loop.SetBody(func() {
				if block := ret.Block(); block != nil {
					y.VisitBlock(block)
				}
			})
			loop.Finish()
		}
		return nil
	case *javaparser.WhileStatementContext:
		// 处理 while 语句
		loop := y.CreateLoopBuilder()
		parExpr := ret.ParExpression().(*javaparser.ParExpressionContext)

		if parExpr != nil {
			expr := parExpr.Expression()
			loop.SetCondition(func() ssa.Value {
				condition := y.VisitExpression(expr)
				if condition == nil {
					condition = y.EmitConstInst(true)
				} else {
					condition = y.VisitExpression(expr)
				}
				return condition
			})
		}
		loop.SetBody(func() {
			if block := ret.Block(); block != nil {
				y.VisitBlock(block)
			}
		})
		loop.Finish()
		return nil
	case *javaparser.DoWhileStatementContext:
		// 处理 do while 语句
		loop := y.CreateLoopBuilder()

		loop.SetCondition(
			func() ssa.Value {
				return y.EmitConstInst(true)
			})
		parExprs := ret.ParExpressionList().(*javaparser.ParExpressionListContext)
		if parExprs != nil {
			exprs := parExprs.ExpressionList()
			loop.SetThird(func() []ssa.Value {
				return y.VisitExpressionList(exprs)
			})
		}

		loop.SetBody(func() {
			if block := ret.Block(); block != nil {
				y.VisitBlock(block)
			}
		})

		loop.Finish()
		return nil
	case *javaparser.TryStatementContext:
		// 处理 try 语句
		if ret.TRY() != nil {
			tryBuilder := y.BuildTry()

			tryBuilder.BuildTryBlock(func() {
				if ret := ret.Block(); ret != nil {
					y.VisitBlock(ret)
				}
			})
			// todo: 多catch case的情况
			// 处理第一个catch
			if catchClause := ret.CatchClause(0).(*javaparser.CatchClauseContext); catchClause != nil {
				tryBuilder.BuildError(func() string {
					id := catchClause.Identifier().GetText()
					return id
				})
				tryBuilder.BuildCatch(func() {
					y.VisitBlock(catchClause.Block())
				})
				if finallyBlock := ret.FinallyBlock(); finallyBlock != nil {
					tryBuilder.BuildFinally(func() {
						y.VisitBlock(finallyBlock.(*javaparser.FinallyBlockContext).Block())
					})
				}
			} else if finallyBlock := ret.FinallyBlock(); finallyBlock != nil {
				tryBuilder.BuildFinally(func() {
					y.VisitBlock(finallyBlock.(*javaparser.FinallyBlockContext).Block())
				})
			}
			tryBuilder.Finish()

		}
		return nil
	case *javaparser.TryWithResourcesStatementContext:
		// 处理 try with resources 语句
		// todo 处理try with resources语句
		return nil
	case *javaparser.SwitchStatementContext:
		// 处理 switch 语句
		SwitchBuilder := y.BuildSwitch()
		SwitchBuilder.AutoBreak = false
		// 设置switch的参数
		var cond ssa.Value
		parExpr := ret.ParExpression().(*javaparser.ParExpressionContext)
		if expr := parExpr.Expression(); expr != nil {
			SwitchBuilder.BuildCondition(func() ssa.Value {
				cond = y.VisitExpression(expr)
				return cond
			})
		} else {
			y.NewError(ssa.Warn, "javaast", "switch expression is nil")
		}
		// 设置case数目
		allcase := ret.AllCASE()
		SwitchBuilder.BuildCaseSize(len(allcase))
		// 设置case参数
		SwitchBuilder.SetCase(func(i int) []ssa.Value {
			if exprList := ret.ExpressionList(i); exprList != nil {
				return y.VisitExpressionList(exprList)
			}
			return nil
		})
		// 设置case执行体
		SwitchBuilder.BuildBody(func(i int) {
			if stmtList := ret.StatementList(i); stmtList != nil {
				y.VisitStatementList(stmtList)
			}
		})
		//设置defalut
		if ret.DEFAULT() != nil {
			if stmtlist := ret.StatementList(len(allcase)); stmtlist != nil {
				SwitchBuilder.BuildDefault(func() {
					log.Infof("aaa:%v", stmtlist.GetText())
					y.VisitStatementList(stmtlist)
				})
			}
		}
		SwitchBuilder.Finish()
	case *javaparser.SynchronizedStatementContext:
		// 处理 synchronized 语句
		return nil
	case *javaparser.ReturnStatementContext:
		// 处理 return 语句
		if ret.Expression() != nil {
			value := y.VisitExpression(ret.Expression())
			y.EmitReturn([]ssa.Value{value})
		} else {
			y.EmitReturn(nil)
		}
		return nil
	case *javaparser.ThrowStatementContext:
		// 处理 throw 语句
		return nil
	case *javaparser.BreakStatementContext:
		// 处理 break 语句
		// todo break使用标签
		if !y.Break() {
			y.NewError(ssa.Error, "javaast", yak2ssa.UnexpectedBreakStmt())
		}
		return nil
	case *javaparser.ContinueStatementContext:
		// 处理 continue 语句
		// todo continue使用标签
		if !y.Continue() {
			y.NewError(ssa.Error, "javaast", yak2ssa.UnexpectedContinueStmt())
		}
		return nil
	case *javaparser.YieldStatementContext:
		// 处理 yield 语句
		return y.VisitExpression(ret.Expression())
	case *javaparser.ExpressionStatementContext:
		// 处理表达式语句
		return y.VisitExpression(ret.Expression())
	case *javaparser.SwitchArrowExpressionContext:
		// 处理 switch 箭头语句
		_ = y.VisitSwitchExpression(ret.SwitchExpression(), false)
		return nil
	case *javaparser.IdentifierLabelStatementContext:
		// 处理标识符标签语句
		return nil
	default:
		return nil
	}
	return nil
}

func (y *builder) VisitLocalTypeDeclaration(raw javaparser.ILocalTypeDeclarationContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.LocalTypeDeclarationContext)
	if i == nil {
		return nil
	}

	return nil
}

func (y *builder) VisitLocalVariableDeclaration(raw javaparser.ILocalVariableDeclarationContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.LocalVariableDeclarationContext)
	if i == nil {
		return nil
	}

	if ret := i.Identifier(); ret != nil {
		log.Infof("visit local variable declaration: %v", ret.GetText())
		variable := y.CreateLocalVariable(ret.GetText())
		value := y.VisitExpression(i.Expression())
		y.AssignVariable(variable, value)
	} else if ret := i.VariableDeclarators(); ret != nil {
		log.Infof("visit local variable declaration: %v", ret.GetText())
		decls := ret.(*javaparser.VariableDeclaratorsContext)
		for _, decl := range decls.AllVariableDeclarator() {
			y.VisitVariableDeclarator(decl)
		}
	}

	return nil
}

func (y *builder) VisitVariableDeclarator(raw javaparser.IVariableDeclaratorContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.VariableDeclaratorContext)
	if i == nil {
		return nil
	}

	if i.VariableInitializer() != nil {
		text := i.VariableDeclaratorId().(*javaparser.VariableDeclaratorIdContext).Identifier().GetText()
		variable := y.CreateVariable(text)
		value := y.VisitVariableInitializer(i.VariableInitializer())
		y.AssignVariable(variable, value)
	} else {
		text := i.VariableDeclaratorId().(*javaparser.VariableDeclaratorIdContext).Identifier().GetText()
		y.CreateVariable(text)
	}

	return nil

}

func (y *builder) VisitVariableInitializer(raw javaparser.IVariableInitializerContext) ssa.Value {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.VariableInitializerContext)
	if i == nil {
		return nil
	}

	if ret := i.Expression(); ret != nil {
		return y.VisitExpression(ret)
	} else if ret := i.ArrayInitializer(); ret != nil {
		return y.VisitArrayInitializer(ret)
	}
	return nil
}

func (y *builder) VisitArguments(raw javaparser.IArgumentsContext) []ssa.Value {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.ArgumentsContext)
	if i == nil {
		return nil
	}

	var args []ssa.Value
	if ret := i.ExpressionList(); ret != nil {
		exprs := ret.(*javaparser.ExpressionListContext)
		for _, expr := range exprs.AllExpression() {
			args = append(args, y.VisitExpression(expr))
		}
	}
	return args
}

func (y *builder) VisitExpressionList(raw javaparser.IExpressionListContext) []ssa.Value {
	if y == nil || raw == nil {
		return nil
	}
	i, _ := raw.(*javaparser.ExpressionListContext)
	if i == nil {
		return nil
	}
	exprs := i.AllExpression()
	valueLen := len(exprs)
	values := make([]ssa.Value, 0, valueLen)
	for _, expr := range exprs {
		if expr != nil {
			if v := y.VisitExpression(expr); !utils.IsNil(v) {
				values = append(values, v)
			}
		}
	}
	return values
}

func (y *builder) VisitStatementList(raw javaparser.IStatementListContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}
	i, _ := raw.(*javaparser.StatementListContext)
	if i == nil {
		return nil
	}

	for _, stmt := range i.AllStatement() {
		if stmt != nil {
			y.VisitStatement(stmt)
		}
	}
	return nil
}

func (y *builder) VisitForControl(raw javaparser.IForControlContext) *ssa.LoopBuilder {
	if y == nil || raw == nil {
		return nil
	}
	i, _ := raw.(*javaparser.ForControlContext)
	if i == nil {
		return nil
	}

	var cond javaparser.IExpressionContext
	loop := y.CreateLoopBuilder()
	if i.EnhancedForControl() != nil {
		//处理增强for循环形式(for each)
		enhanced := i.EnhancedForControl().(*javaparser.EnhancedForControlContext)
		var value ssa.Value
		loop.SetFirst(func() []ssa.Value {
			value = y.VisitExpression(enhanced.Expression())
			return []ssa.Value{value}
		})

		loop.SetCondition(func() ssa.Value {
			var variable *ssa.Variable
			if decl := enhanced.VariableDeclaratorId(); decl != nil {
				text := decl.(*javaparser.VariableDeclaratorIdContext).Identifier().GetText()
				variable = y.CreateVariable(text)
			}
			_, field, ok := y.EmitNext(value, false)
			y.AssignVariable(variable, field)
			return ok
		})
		return loop
	} else {
		// 处理标准for循环形式
		// 设置第一个参数
		if first := i.ForInit(); first != nil {
			loop.SetFirst(func() []ssa.Value { return y.VisitForInit(first) })
		}
		// 设置第二个参数
		if expr := i.Expression(); expr != nil {
			cond = expr
		}
		// 设置第三个参数
		if third := i.GetForUpdate(); third != nil {
			loop.SetThird(func() []ssa.Value { return y.VisitExpressionList(third) })
		}
	}
	// 设置循环条件
	loop.SetCondition(func() ssa.Value {
		var condition ssa.Value
		if cond == nil {
			condition = y.EmitConstInst(true)
		} else {
			condition = y.VisitExpression(cond)
		}
		return condition
	})
	return loop
}

func (y *builder) VisitForInit(raw javaparser.IForInitContext) []ssa.Value {
	if y == nil || raw == nil {
		return nil
	}
	i, _ := raw.(*javaparser.ForInitContext)
	if i == nil {
		return nil
	}
	// for循环first为局部变量声明
	// for(int a=1;;){}
	var value []ssa.Value
	if ret := i.LocalVariableDeclaration(); ret != nil {
		y.VisitLocalVariableDeclaration(ret)
		// 访问expressionlist获取变量名的Value
		if name := ret.(*javaparser.LocalVariableDeclarationContext).Identifier(); name != nil {
			text := name.GetText()
			value = append(value, y.ReadValue(text))
		} else if name := ret.(*javaparser.LocalVariableDeclarationContext).VariableDeclarators(); name != nil {
			// 访问localVariableDeclaration，定义变量，并获取变量名的value
			y.VisitLocalVariableDeclaration(ret)
			// 获取所有定义变量的变量名
			decls := name.(*javaparser.VariableDeclaratorsContext)
			for _, decl := range decls.AllVariableDeclarator() {
				if decl != nil {
					variableDeclaratorId := decl.(*javaparser.VariableDeclaratorContext).VariableDeclaratorId()
					text := variableDeclaratorId.(*javaparser.VariableDeclaratorIdContext).Identifier().GetText()
					value = append(value, y.ReadValue(text))
				}
			}
		}
	}
	return value
}

func (y *builder) VisitIfStmt(raw javaparser.IIfstmtContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}

	builder := y.CreateIfBuilder()

	var build func(raw javaparser.IIfstmtContext) func()
	build = func(raw javaparser.IIfstmtContext) func() {
		if raw == nil {
			return nil
		}
		i, _ := raw.(*javaparser.IfstmtContext)
		if i == nil {
			return nil
		}

		if parExpr := i.ParExpression(); parExpr != nil {
			expr := parExpr.(*javaparser.ParExpressionContext).Expression()
			if block := i.Block(); block != nil {
				builder.AppendItem(
					func() ssa.Value { return y.VisitExpression(expr) },
					func() {
						y.VisitBlock(i.Block())
					})
			} else {
				// 没有block的情况
				builder.AppendItem(
					func() ssa.Value { return y.VisitExpression(expr) },
					func() {})
			}

		}

		for _, elseIfBlock := range i.AllElseIfBlock() {
			if elseIfBlock != nil {
				elseIfStmt := elseIfBlock.(*javaparser.ElseIfBlockContext)
				block := elseIfStmt.Block()
				parExpr := elseIfStmt.ParExpression()
				expr := parExpr.(*javaparser.ParExpressionContext).Expression()
				builder.AppendItem(
					func() ssa.Value {
						return (y.VisitExpression(expr))
					},
					func() {
						y.VisitBlock(block)
					},
				)
			}
		}
		elseStmt := i.ElseBlock()
		if elseStmt != nil {
			if elseBlock := elseStmt.(*javaparser.ElseBlockContext).Block(); elseBlock != nil {
				return func() { y.VisitBlock(elseBlock) }
			}
		}
		return nil
	}
	elseBlock := build(raw)
	builder.SetElse(elseBlock)
	builder.Build()
	return nil
}

func (y *builder) VisitSwitchExpression(raw javaparser.ISwitchExpressionContext, isExpression bool) ssa.Value {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.SwitchExpressionContext)
	if i == nil {
		return nil
	}

	switchBuilder := y.BuildSwitch()
	switchBuilder.AutoBreak = false

	parExpr := i.ParExpression().(*javaparser.ParExpressionContext)
	expr := parExpr.Expression()
	if expr != nil {
		switchBuilder.BuildCondition(func() ssa.Value {
			return y.VisitExpression(expr)
		})
	} else {
		y.NewError(ssa.Error, "javaast", "switch expression is nil")
	}

	switchLabels := i.AllSwitchLabeledRule()
	caseNum := len(switchLabels)
	//得到case后面参数的value
	getCaseValue := func(i int) []ssa.Value {
		switchStmt := switchLabels[i].(*javaparser.SwitchLabeledRuleContext)
		if switchStmt.ExpressionList() != nil {
			return y.VisitExpressionList(switchStmt.ExpressionList())
		} else if switchStmt.NULL_LITERAL() != nil {
			return []ssa.Value{y.EmitConstInstNil()}
		} else if switchStmt.GuardedPattern() != nil {
			return []ssa.Value{y.EmitConstInstNil()} // todo: 处理guarded pattern
		} else {
			return nil
		}
	}

	switchBuilder.BuildCaseSize(caseNum)
	switchBuilder.SetCase(func(i int) []ssa.Value {
		return getCaseValue(i)
	})

	switchBuilder.BuildBody(func(i int) {
		switchStmt := switchLabels[i].(*javaparser.SwitchLabeledRuleContext)
		if switchRuleOutCome := switchStmt.SwitchRuleOutcome(); switchRuleOutCome != nil {
			s := switchRuleOutCome.(*javaparser.SwitchRuleOutcomeContext)
			if s.Block() != nil {
				y.VisitBlock(s.Block())
			}
			for _, stmt := range s.AllBlockStatement() {
				y.VisitBlockStatement(stmt)
			}
		}
	})

	if i.DefaultLabeledRule() != nil {
		switchBuilder.BuildDefault(func() {
			if defaultStmt := i.DefaultLabeledRule().(*javaparser.DefaultLabeledRuleContext); defaultStmt != nil {
				switchRuleOutCome := defaultStmt.SwitchRuleOutcome()
				s := switchRuleOutCome.(*javaparser.SwitchRuleOutcomeContext)
				if s.Block() != nil {
					y.VisitBlock(s.Block())
				}
				for _, stmt := range s.AllBlockStatement() {
					y.VisitBlockStatement(stmt)
				}
			}
		})
	}

	switchBuilder.Finish()
	// switch 作为expression
	if isExpression {
		// 当switch作为expression的时候需要返回ssa.Value
		// 得到blockStatement的ssa.Value
		// 因为blockStatement并不所有的语句都会返回ssa.Value
		// 而switch作为expression的时候需要返回ssa.Value
		// todo 处理yeild语句
		getBlockValue := func(stmt javaparser.IBlockContext) []ssa.Value {
			if stmt == nil {
				return nil
			}
			block := stmt.(*javaparser.BlockContext)
			if blockStmtList := block.BlockStatementList(); blockStmtList != nil {
				blockStmts := blockStmtList.(*javaparser.BlockStatementListContext)
				for _, blockStmt := range blockStmts.AllBlockStatement() {
					blockStatement := blockStmt.(*javaparser.BlockStatementContext)
					if blockStatement.Statement() != nil {
						statement := blockStatement.Statement()
						switch ret := statement.(type) {
						case *javaparser.YieldStatementContext:
							return []ssa.Value{y.VisitExpression(ret.Expression())}
						}
					}
				}
			}

			return nil
		}

		getBlockStmtValue := func(stmt javaparser.IBlockStatementContext) []ssa.Value {
			if stmt == nil {
				return nil
			}
			blockStatement := stmt.(*javaparser.BlockStatementContext)
			if blockStatement.Statement() != nil {
				statement := blockStatement.Statement()
				switch ret := statement.(type) {
				case *javaparser.ExpressionStatementContext:
					return []ssa.Value{y.VisitExpression(ret.Expression())}
				}
			}
			return nil
		}
		// 遍历case的switchRuleOutcome的block和blockStatement
		getSwitchOutcomeValue := func(i int) ssa.Value {
			var value []ssa.Value
			switchStmt := switchLabels[i].(*javaparser.SwitchLabeledRuleContext)
			if switchRuleOutCome := switchStmt.SwitchRuleOutcome(); switchRuleOutCome != nil {
				s := switchRuleOutCome.(*javaparser.SwitchRuleOutcomeContext)
				if s.Block() != nil {
					block := s.Block().(*javaparser.BlockContext)
					value = append(value, getBlockValue(block)...)

				}
				for _, blockStmt := range s.AllBlockStatement() {
					value = append(value, getBlockStmtValue(blockStmt)...)
				}
			}
			// switch 作为参数的时候只能返回一个value
			if len(value) > 1 {
				y.NewError(ssa.Warn, "javaast", "switch as expression can only return one value")
				return nil
			} else {
				return value[0]
			}
		}
		// 遍历default的switchRuleOutcome的block和blockStatement
		getDefalutOutCome := func() ssa.Value {
			var value []ssa.Value
			defaultStmt := i.DefaultLabeledRule().(*javaparser.DefaultLabeledRuleContext)
			if switchRuleOutCome := defaultStmt.SwitchRuleOutcome(); switchRuleOutCome != nil {
				s := switchRuleOutCome.(*javaparser.SwitchRuleOutcomeContext)
				if s.Block() != nil {
					block := s.Block().(*javaparser.BlockContext)
					value = append(value, getBlockValue(block)...)

				}
				for _, blockStmt := range s.AllBlockStatement() {
					value = append(value, getBlockStmtValue(blockStmt)...)
				}
			}
			if len(value) > 1 {
				y.NewError(ssa.Warn, "javaast", "switch as expression can only return one value")
				return nil
			} else if len(value) == 1 {
				return value[0]
			} else {
				return y.EmitConstInstNil()
			}
		}
		// switch参数的value
		value1 := y.VisitExpression(expr)
		for i := 0; i < caseNum; i++ {
			// case参数的value
			value2 := getCaseValue(i)
			for _, v := range value2 {
				if value1.String() == v.String() {
					return getSwitchOutcomeValue(i)
				}
			}
		}
		if i.DefaultLabeledRule() != nil {
			return getDefalutOutCome()
		} else {
			return nil
		}
	} else {
		return nil
	}

}

func (y *builder) VisitGuardedPattern(raw javaparser.IGuardedPatternContext) []ssa.Value {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.GuardedPatternContext)
	if i == nil {
		return nil
	}
	return nil

}

func (y *builder) VisitBlockStatementList(raw javaparser.IBlockStatementListContext) {
	if y == nil || raw == nil {
		return
	}

	i, _ := raw.(*javaparser.BlockStatementListContext)
	if i == nil {
		return
	}
	for _, stmt := range i.AllBlockStatement() {
		if stmt != nil {
			y.VisitBlockStatement(stmt)
		}
	}
}

func (y *builder) VisitCreator(raw javaparser.ICreatorContext) ssa.Value {
	if y == nil || raw == nil {
		return nil
	}
	i, _ := raw.(*javaparser.CreatorContext)
	if i == nil {
		return nil
	}
	// todo new声明类的方式
	if ret := i.ArrayCreatorRest(); ret != nil {
		return y.VisitArrayCreatorRest(ret)
	}
	return nil
}

func (y *builder) VisitArrayCreatorRest(raw javaparser.IArrayCreatorRestContext) ssa.Value {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*javaparser.ArrayCreatorRestContext)
	if i == nil {
		return nil
	}
	// 数组声明
	if ret := i.ArrayInitializer(); ret != nil {
		return y.VisitArrayInitializer(ret)
	}
	// todo 初始化数组

	return nil
}
