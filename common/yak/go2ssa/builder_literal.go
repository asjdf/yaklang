package go2ssa

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	gol "github.com/yaklang/yaklang/common/yak/antlr4go/parser"
	"github.com/yaklang/yaklang/common/yak/ssa"
)

func (b *astbuilder) buildBoolLiteral(name string) ssa.Value {
	boolLit, err := strconv.ParseBool(name)
	if err != nil {
		b.NewError(ssa.Error, TAG, UnhandledBool())
	}
	return b.EmitConstInst(boolLit)
}

func (b *astbuilder) buildLiteral(exp *gol.LiteralContext) ssa.Value {
	recoverRange := b.SetRange(exp.BaseParserRuleContext)
	defer recoverRange()

	if lit := exp.BasicLit(); lit != nil {
		return b.buildBasicLit(lit.(*gol.BasicLitContext))
	}

	if lit := exp.CompositeLit(); lit != nil {
		return b.buildCompositeLit(lit.(*gol.CompositeLitContext))
	}

	if lit := exp.FunctionLit(); lit != nil {
		return b.buildFunctionLit(lit.(*gol.FunctionLitContext))
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return b.EmitConstInst(0)
}

func (b *astbuilder) buildFunctionLit(exp *gol.FunctionLitContext) ssa.Value {
	recoverRange := b.SetRange(exp.BaseParserRuleContext)
	defer recoverRange()

	b.SupportClosure = true
	newFunc := b.NewFunc("")

	hitDefinedFunction := false
	MarkedFunctionType := b.GetMarkedFunction()
	handleFunctionType := func(fun *ssa.Function) {
		fun.ParamLength = len(fun.Params)
		if MarkedFunctionType == nil {
			return
		}
		if len(fun.Params) != len(MarkedFunctionType.Parameter) {
			return
		}

		for i, p := range fun.Params {
			p.SetType(MarkedFunctionType.Parameter[i])
		}
		hitDefinedFunction = true
	}

	{
		recoverRange := b.SetRange(exp.BaseParserRuleContext)
		b.FunctionBuilder = b.PushFunction(newFunc)

		if para, ok := exp.Signature().(*gol.SignatureContext); ok {
			b.buildSignature(para)
		}

		handleFunctionType(b.Function)

		if block, ok := exp.Block().(*gol.BlockContext); ok {
			b.buildBlock(block)
		}

		b.Finish()
		b.FunctionBuilder = b.PopFunction()
		if hitDefinedFunction {
			b.MarkedFunctions = append(b.MarkedFunctions, newFunc)
		}
		recoverRange()
	}

	b.SupportClosure = false
	return newFunc
}

type keyValue struct {
	key   ssa.Value
	value ssa.Value
	kv    []keyValue
}

func (b *astbuilder) buildCompositeLit(exp *gol.CompositeLitContext) ssa.Value {
	recoverRange := b.SetRange(exp.BaseParserRuleContext)
	defer recoverRange()

	var kvs []keyValue

	typ := b.buildLiteralType(exp.LiteralType().(*gol.LiteralTypeContext))
	if typ == nil { // 目前还没法识别golang的标准库,这里暂时识别为StructType
		typ = ssa.NewStructType()
	}
	if value := exp.LiteralValue(); value != nil {
		if s, ok := value.(*gol.LiteralValueContext); ok {
			switch t := typ.(type) {
			case *ssa.ObjectType:
				if t.GetTypeKind() == ssa.StructTypeKind {
					kvs = b.buildLiteralValue(s, true)
				} else if t.GetTypeKind() == ssa.SliceTypeKind {
					kvs = b.buildLiteralValue(s, false)
				} else {
					kvs = b.buildLiteralValue(s, false)
				}
			case *ssa.AliasType: // 处理golang标准库
				// TODO
				typ = typ.(*ssa.AliasType).GetType()
				kvs = b.buildLiteralValue(s, true)
			default:
				typ = ssa.CreateAnyType()
				kvs = b.buildLiteralValue(s, false)
			}
		}
	}

	var typeHander func(ssa.Type, []keyValue) ssa.Value

	typeHander = func(typ ssa.Type, kvs []keyValue) ssa.Value {
		var obj ssa.Value

		switch typ.GetTypeKind() {
		case ssa.SliceTypeKind, ssa.BytesTypeKind:
			objt := typ.(*ssa.ObjectType)
			if len(kvs) == 0 {
				return b.CreateInterfaceWithMap(nil, nil)
			}
			if kvs[0].value != nil {
				return kvs[0].value
			} else {
				obj = b.InterfaceAddFieldBuild(len(kvs),
					func(i int) ssa.Value {
						return b.EmitConstInst(i)
					},
					func(i int) ssa.Value {
						return typeHander(objt.FieldType, kvs[i].kv)
					})
			}
		case ssa.MapTypeKind:
			objt := typ.(*ssa.ObjectType)
			if len(kvs) == 0 {
				return b.CreateInterfaceWithMap(nil, nil)
			}
			if kvs[0].value != nil {
				return kvs[0].value
			} else {
				obj = b.InterfaceAddFieldBuild(len(kvs),
					func(i int) ssa.Value {
						return kvs[i].key
					},
					func(i int) ssa.Value {
						return typeHander(objt.FieldType, kvs[i].kv)
					})
			}
		case ssa.StructTypeKind:
			objt := typ.(*ssa.ObjectType)
			if len(kvs) == 0 {
				return b.CreateInterfaceWithMap(nil, nil)
			}
			if kvs[0].value != nil {
				return kvs[0].value
			} else {
				if kvs[0].key == nil { // 全部初始化
					obj = b.InterfaceAddFieldBuild(len(kvs),
						func(i int) ssa.Value {
							if i < len(objt.Keys) {
								return objt.Keys[i]
							} else {
								return b.EmitConstInst("")
							}
						},
						func(i int) ssa.Value {
							return typeHander(objt.FieldTypes[i], kvs[i].kv)
						})
				} else { // 部分初始化
					obj = b.InterfaceAddFieldBuild(len(objt.Keys),
						func(i int) ssa.Value {
							return objt.Keys[i]
						},
						func(i int) ssa.Value {
							for y, kv := range kvs {
								if objt.Keys[i].String() == kv.key.String() {
									return typeHander(objt.FieldTypes[i], kvs[y].kv)
								}
							}
							return b.GetDefaultValue(objt.FieldTypes[i])
						})
				}
			}
		case ssa.InterfaceTypeKind:
			// TODO
			obj = b.InterfaceAddFieldBuild(0,
				func(i int) ssa.Value {
					return b.EmitConstInst(i)
				},
				func(i int) ssa.Value {
					return b.EmitConstInst(i)
				})
		case ssa.AnyTypeKind:
			// TODO
			return b.EmitUndefined(typ.String())
		case ssa.UndefinedTypeKind:
			obj = b.InterfaceAddFieldBuild(0,
				func(i int) ssa.Value {
					return b.EmitConstInst(i)
				},
				func(i int) ssa.Value {
					return b.EmitConstInst(i)
				})
		case ssa.NumberTypeKind, ssa.StringTypeKind, ssa.BooleanTypeKind:
			return kvs[0].value
		default:
			if kvs[0].value != nil {
				return kvs[0].value
			}
			b.NewError(ssa.Error, TAG, "unhandled type")
			return b.EmitConstInst(0)
		}
		coverType(obj.GetType(), typ)
		return obj
	}

	return typeHander(typ, kvs)
}

func (b *astbuilder) buildLiteralValue(exp *gol.LiteralValueContext, iscreate bool) (ret []keyValue) {
	recoverRange := b.SetRange(exp.BaseParserRuleContext)
	defer recoverRange()

	if list := exp.ElementList(); list != nil {
		for _, e := range list.(*gol.ElementListContext).AllKeyedElement() {
			kv := b.buildKeyedElement(e.(*gol.KeyedElementContext), iscreate)
			ret = append(ret, kv)
		}
	}

	return ret
}

func (b *astbuilder) buildKeyedElement(exp *gol.KeyedElementContext, iscreate bool) (ret keyValue) {
	recoverRange := b.SetRange(exp.BaseParserRuleContext)
	defer recoverRange()
	var keys ssa.Value
	var kvs []keyValue

	if key := exp.Key(); key != nil {
		keys = b.buildKey(key.(*gol.KeyContext), iscreate)
	}

	if elem := exp.Element(); elem != nil {
		if s, ok := elem.(*gol.ElementContext); ok {
			kvs = b.buildElement(s, iscreate)
		}
	}

	return keyValue{
		key:   keys,
		value: nil,
		kv:    kvs,
	}
}

func (b *astbuilder) buildKey(exp *gol.KeyContext, iscreate bool) ssa.Value {
	recoverRange := b.SetRange(exp.BaseParserRuleContext)
	defer recoverRange()

	if e := exp.Expression(); e != nil {
		if iscreate {
			var leftv ssa.Value
			if p := e.(*gol.ExpressionContext).PrimaryExpr(); p != nil {
				if o := p.(*gol.PrimaryExprContext).Operand(); o != nil {
					if n := o.(*gol.OperandContext).OperandName(); n != nil {
						id := n.(*gol.OperandNameContext).IDENTIFIER()
						leftv = b.EmitConstInst(id.GetText())
					}
				}
			}
			return leftv
		} else {
			rightv, _ := b.buildExpression(e.(*gol.ExpressionContext), false)
			return rightv
		}
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return b.EmitConstInst(0)
}

func (b *astbuilder) buildElement(exp *gol.ElementContext, iscreate bool) (ret []keyValue) {
	recoverRange := b.SetRange(exp.BaseParserRuleContext)
	defer recoverRange()

	if e := exp.Expression(); e != nil {
		right, _ := b.buildExpression(e.(*gol.ExpressionContext), false)
		kv := keyValue{
			key:   nil,
			value: right,
			kv:    []keyValue{},
		}
		ret = append(ret, kv)
		return ret
	}

	if e := exp.LiteralValue(); e != nil {
		return b.buildLiteralValue(e.(*gol.LiteralValueContext), iscreate)
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return ret
}

func (b *astbuilder) buildLiteralType(stmt *gol.LiteralTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	if name := stmt.TypeName(); name != nil {
		return b.buildTypeName(name.(*gol.TypeNameContext))
	}

	if stmt.ELLIPSIS() != nil {
		return b.buildSliceTypeELiteral(stmt)
	}

	// slice type literal
	if s, ok := stmt.SliceType().(*gol.SliceTypeContext); ok {
		return b.buildSliceTypeLiteral(s)
	}

	// array type literal
	if s, ok := stmt.ArrayType().(*gol.ArrayTypeContext); ok {
		return b.buildArrayTypeLiteral(s)
	}

	// map type literal
	if s, ok := stmt.MapType().(*gol.MapTypeContext); ok {
		return b.buildMapTypeLiteral(s)
	}

	// struct type literal
	if s, ok := stmt.StructType().(*gol.StructTypeContext); ok {
		return b.buildStructTypeLiteral(s)
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return ssa.CreateAnyType()
}

func (b *astbuilder) buildTypeLit(stmt *gol.TypeLitContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()
	text := stmt.GetText()

	// slice type literal
	if s := stmt.SliceType(); s != nil {
		return b.buildSliceTypeLiteral(s.(*gol.SliceTypeContext))
	}

	// array type literal
	if s := stmt.ArrayType(); s != nil {
		return b.buildArrayTypeLiteral(s.(*gol.ArrayTypeContext))
	}

	// map type literal
	if strings.HasPrefix(text, "map") {
		if s := stmt.MapType(); s != nil {
			return b.buildMapTypeLiteral(s.(*gol.MapTypeContext))
		}
	}

	// struct type literal
	if strings.HasPrefix(text, "struct") {
		if s := stmt.StructType(); s != nil {
			return b.buildStructTypeLiteral(s.(*gol.StructTypeContext))
		}
	}

	// pointer type literal
	if strings.HasPrefix(text, "*") {
		if p := stmt.PointerType(); p != nil {
			if t := p.(*gol.PointerTypeContext).Type_(); t != nil {
				return b.buildType(t.(*gol.Type_Context))
			}
		}
	}

	// function type literal
	if strings.HasPrefix(text, "func") {
		if s := stmt.FunctionType(); s != nil {
			return b.buildFunctionTypeLiteral(s.(*gol.FunctionTypeContext))
		}
	}

	// interface type literal
	if strings.HasPrefix(text, "interface") {
		if s := stmt.InterfaceType(); s != nil {
			return b.buildInterfaceTypeLiteral(s.(*gol.InterfaceTypeContext))
		}
	}

	// channel type literal
	if strings.HasPrefix(text, "chan") ||
		strings.HasPrefix(text, "<-chan") ||
		strings.HasPrefix(text, "chan<-") {
		if s := stmt.ChannelType(); s != nil {
			return b.buildChanTypeLiteral(s.(*gol.ChannelTypeContext))
		}
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return ssa.CreateAnyType()
}

func (b *astbuilder) buildFunctionTypeLiteral(stmt *gol.FunctionTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	if signature := stmt.Signature(); signature != nil {
		paramt, rett := b.buildSignature(signature.(*gol.SignatureContext))
		return ssa.NewFunctionType("", paramt, rett, false)
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return ssa.CreateAnyType()
}

func (b *astbuilder) buildInterfaceTypeLiteral(stmt *gol.InterfaceTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	interfacetyp := ssa.NewInterfaceType("", "")

	for _, t := range stmt.AllTypeElement() {
		ssatyp := b.buildTypeElement(t.(*gol.TypeElementContext))
		switch t := ssatyp.(type) {
		case *ssa.InterfaceType:
			interfacetyp.AddFatherInterfaceType(t)
		case *ssa.ObjectType:
			interfacetyp.AddStructure(t.Name, t)
		}
	}

	for _, f := range stmt.AllMethodSpec() {
		b.buildMethodSpec(f.(*gol.MethodSpecContext), interfacetyp)
	}

	return interfacetyp
}

func (b *astbuilder) buildChanTypeLiteral(stmt *gol.ChannelTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	if etyp := stmt.ElementType(); etyp != nil {
		if typ := etyp.(*gol.ElementTypeContext).Type_(); typ != nil {
			ssatyp := b.buildType(typ.(*gol.Type_Context))
			return ssa.NewChanType(ssatyp)
		}
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return ssa.CreateAnyType()
}

func (b *astbuilder) buildMapTypeLiteral(stmt *gol.MapTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	var keyTyp ssa.Type
	var valueTyp ssa.Type
	if s, ok := stmt.Type_().(*gol.Type_Context); ok {
		keyTyp = b.buildType(s)
	}

	// value
	if s, ok := stmt.ElementType().(*gol.ElementTypeContext); ok {
		valueTyp = b.buildType(s.Type_().(*gol.Type_Context))
	}
	if keyTyp != nil && valueTyp != nil {
		return ssa.NewMapType(keyTyp, valueTyp)
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return ssa.CreateAnyType()
}

func (b *astbuilder) buildSliceTypeLiteral(stmt *gol.SliceTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	var ssatyp ssa.Type
	if stmt.GetText() == "[]byte" || stmt.GetText() == "[]uint8" {
		return ssa.BasicTypes[ssa.BytesTypeKind]
	}
	if s, ok := stmt.ElementType().(*gol.ElementTypeContext); ok {
		if eleTyp := b.buildType(s.Type_().(*gol.Type_Context)); eleTyp != nil {
			ssatyp = ssa.NewSliceType(eleTyp)
		}
	}
	return ssatyp
}

func (b *astbuilder) buildSliceTypeELiteral(stmt *gol.LiteralTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	var ssatyp ssa.Type
	if s, ok := stmt.ElementType().(*gol.ElementTypeContext); ok {
		if eleTyp := b.buildType(s.Type_().(*gol.Type_Context)); eleTyp != nil {
			ssatyp = ssa.NewSliceType(eleTyp)
		}
	}
	return ssatyp
}

func (b *astbuilder) buildArrayTypeLiteral(stmt *gol.ArrayTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	var value ssa.Value
	var ssatyp ssa.Type
	if s, ok := stmt.ArrayLength().(*gol.ArrayLengthContext); ok {
		if e := s.Expression(); e != nil {
			rightv, _ := b.buildExpression(e.(*gol.ExpressionContext), false)
			value = rightv
		}
	}

	if s, ok := stmt.ElementType().(*gol.ElementTypeContext); ok {
		if eleTyp := b.buildType(s.Type_().(*gol.Type_Context)); eleTyp != nil {
			ssatyp = ssa.NewSliceType(eleTyp)
		}
	}
	_ = value
	return ssatyp
}

func (b *astbuilder) buildStructTypeLiteral(stmt *gol.StructTypeContext) ssa.Type {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	structTyp := ssa.NewStructType()
	for _, s := range stmt.AllFieldDecl() {
		b.buildFieldDecl(s.(*gol.FieldDeclContext), structTyp)
	}
	return structTyp
}

func (b *astbuilder) buildFieldDecl(stmt *gol.FieldDeclContext, structTyp *ssa.ObjectType) {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	var ssatyp ssa.Type = nil
	if typ := stmt.Type_(); typ != nil {
		ssatyp = b.buildType(typ.(*gol.Type_Context))
	}

	if idlist := stmt.IdentifierList(); idlist != nil {
		sList := b.buildStructList(idlist.(*gol.IdentifierListContext))
		if ssatyp != nil {
			for _, p := range sList {
				structTyp.AddField(p, ssatyp)
			}
		}
	}

	if em := stmt.EmbeddedField(); em != nil {
		if typ, ok := em.(*gol.EmbeddedFieldContext); ok {
			parent := b.buildTypeName(typ.TypeName().(*gol.TypeNameContext))
			if a := typ.TypeArgs(); a != nil {
				b.tpHander[b.Function.GetName()] = b.buildTypeArgs(a.(*gol.TypeArgsContext))
			}
			if parent == nil {
				name := typ.TypeName().(*gol.TypeNameContext).GetText()
				b.NewError(ssa.Warn, TAG, StructNotFind(name))
				parent = ssa.NewStructType()
				parent.(*ssa.ObjectType).Name = name
			}
			if p, ok := parent.(*ssa.ObjectType); ok {
				structTyp.AddField(b.EmitConstInst(""), p)
				structTyp.AnonymousField = append(structTyp.AnonymousField, p)
			} else if a, ok := parent.(*ssa.AliasType); ok {
				structTyp.AddField(b.EmitConstInst(a.Name), a.GetType())
			}
		}
	}
}

func (b *astbuilder) buildBasicLit(exp *gol.BasicLitContext) ssa.Value {
	recoverRange := b.SetRange(exp.BaseParserRuleContext)
	defer recoverRange()

	if lit := exp.Integer(); lit != nil {
		return b.buildIntegerLiteral(lit.(*gol.IntegerContext))
	}

	if lit := exp.NIL_LIT(); lit != nil {
		return b.EmitConstInstNil()
	}

	if lit := exp.FLOAT_LIT(); lit != nil {
		t := lit.GetText()
		if strings.HasPrefix(t, ".") {
			t = "0" + t
		}
		f, _ := strconv.ParseFloat(t, 64)
		return b.EmitConstInst(f)
	}

	if lit := exp.String_(); lit != nil {
		return b.buildStringLiteral(lit.(*gol.String_Context))
	}

	if lit := exp.Char_(); lit != nil {
		return b.buildCharLiteral(lit.(*gol.Char_Context))
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return b.EmitConstInst(0)
}

func (b *astbuilder) buildStringLiteral(stmt *gol.String_Context) ssa.Value {
	var text = stmt.GetText()
	if text == "" {
		return b.EmitConstInst(text)
	}

	switch text[0] {
	case '"':
		val, err := strconv.Unquote(text)
		if err != nil {
			b.NewError(ssa.Error, TAG, CannotParseString(stmt.GetText(), err.Error()))
		}
		return b.EmitConstInstWithUnary(val, 0)
	case '`':
		val, err := strconv.Unquote(text)
		if err != nil {
			b.NewError(ssa.Error, TAG, CannotParseString(stmt.GetText(), err.Error()))
		}
		return b.EmitConstInstWithUnary(val, 0)
	}

	b.NewError(ssa.Error, TAG, Unreachable())
	return b.EmitConstInst(0)
}

func (b *astbuilder) buildCharLiteral(stmt *gol.Char_Context) ssa.Value {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	lit := stmt.GetText()
	var s string
	var err error
	if lit == "'\\''" {
		s = "'"
	} else {
		lit = strings.ReplaceAll(lit, `"`, `\"`)
		s, err = strconv.Unquote(fmt.Sprintf("\"%s\"", lit[1:len(lit)-1]))
		if err != nil {
			b.NewError(ssa.Error, TAG, fmt.Sprintf("unquote error %s", err))
			return b.EmitConstInst(0)
		}
	}
	runeChar := []rune(s)[0]
	if runeChar < 256 {
		return b.EmitConstInst(byte(runeChar))
	} else {
		return b.EmitConstInst(runeChar)
	}
}

func (b *astbuilder) buildIntegerLiteral(stmt *gol.IntegerContext) ssa.Value {
	recoverRange := b.SetRange(stmt.BaseParserRuleContext)
	defer recoverRange()

	lit := stmt.GetText()

	if find := strings.Contains(lit, "."); find {
		var f, _ = strconv.ParseFloat(lit, 64)
		return b.EmitConstInst(f)
	} else {
		var err error
		var originStr = stmt.GetText()
		var intStr = strings.ToLower(originStr)
		var resultInt64 int64

		if num := stmt.DECIMAL_LIT(); num != nil { // 十进制
			if strings.Contains(stmt.GetText(), "e") {
				var f, _ = strconv.ParseFloat(intStr, 64)
				return b.EmitConstInst(f)
			}
			resultInt64, err = strconv.ParseInt(intStr, 10, 64)
		} else if num := stmt.HEX_LIT(); num != nil { // 十六进制
			resultInt64, err = strconv.ParseInt(intStr[2:], 16, 64)
		} else if num := stmt.BINARY_LIT(); num != nil { // 二进制
			resultInt64, err = strconv.ParseInt(intStr[2:], 2, 64)
		} else if num := stmt.OCTAL_LIT(); num != nil { // 八进制
			resultInt64, err = strconv.ParseInt(intStr[2:], 8, 64)
		} else {
			b.NewError(ssa.Error, TAG, fmt.Sprintf("cannot parse num for literal: %s", stmt.GetText()))
			return b.EmitConstInst(0)
		}

		if err != nil {
			b.NewError(ssa.Error, TAG, fmt.Sprintf("const parse %s as integer literal... is to large for int64: %v", originStr, err))
			return b.EmitConstInst(0)
		}

		if resultInt64 > math.MaxInt {
			return b.EmitConstInst(int64(resultInt64))
		} else {
			return b.EmitConstInst(int64(resultInt64))
		}
	}
}

func coverType(ityp, iwantTyp ssa.Type) {
	typ, ok := ityp.(*ssa.ObjectType)
	if !ok {
		return
	}
	wantTyp, ok := iwantTyp.(*ssa.ObjectType)
	if !ok {
		return
	}

	typ.SetTypeKind(wantTyp.GetTypeKind())
	switch wantTyp.GetTypeKind() {
	case ssa.SliceTypeKind:
		typ.FieldType = wantTyp.FieldType
	case ssa.MapTypeKind:
		typ.FieldType = wantTyp.FieldType
		typ.KeyTyp = wantTyp.KeyTyp
	case ssa.StructTypeKind:
		typ.FieldType = wantTyp.FieldType
		typ.KeyTyp = wantTyp.KeyTyp
		for n, m := range wantTyp.GetMethod() {
			typ.AddMethod(n, m)
		}
	}
	for _, a := range wantTyp.AnonymousField {
		// TODO: 匿名结构体应该是一个指针，修改时应该要连带父类一起修改
		typ.AnonymousField = append(typ.AnonymousField, a)
	}
}

func (b *astbuilder) GetDefaultValue(ityp ssa.Type) ssa.Value {
	switch ityp.GetTypeKind() {
	case ssa.NumberTypeKind:
		return b.EmitConstInst(0)
	case ssa.StringTypeKind:
		return b.EmitConstInst("")
	case ssa.BooleanTypeKind:
		return b.EmitConstInst(false)
	default:
		return b.EmitConstInst("")
	}
}
