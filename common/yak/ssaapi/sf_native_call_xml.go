package ssaapi

import (
	"encoding/xml"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/syntaxflow/sfvm"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/xml2"
	"github.com/yaklang/yaklang/common/yak/ssa"
	"github.com/yaklang/yaklang/common/yak/ssa/ssadb"
	"regexp"
	"strings"
)

var (
	mybatisVarExtractor = regexp.MustCompile(`\$\{\s*([^}]+)\s*}`)
)

type mybatisXMLMapper struct {
	FullClassName string
	ClassName     string
	Namespace     string
	entityStack   *utils.Stack[*mybatisXMLQuery]
}

type mybatisXMLQuery struct {
	mapper         *mybatisXMLMapper
	Id             string
	CheckParamName []string
}

func (m *mybatisXMLQuery) SyntaxFlowFirst() string {
	if m.mapper == nil {
		return ""
	}
	var builder strings.Builder
	builder.WriteString(m.mapper.ClassName)
	builder.WriteString(".")
	builder.WriteString(m.Id)
	builder.WriteString("(*?{!have: this && opcode: param && any: " + strings.Join(m.CheckParamName, ",") + " } as $__output)")
	return builder.String()
}

func (m *mybatisXMLQuery) SyntaxFlowFinal() string {
	if m.mapper == nil {
		return ""
	}
	var builder strings.Builder
	builder.WriteString(m.mapper.ClassName)
	builder.WriteString(".")
	builder.WriteString(m.Id)
	builder.WriteString("(*?{!have: this && opcode: param } as $__output)")

	return builder.String()
}

var nativeCallMybatixXML = func(v sfvm.ValueOperator, frame *sfvm.SFFrame, params *sfvm.NativeCallActualParams) (bool, sfvm.ValueOperator, error) {
	prog, err := fetchProgram(v)
	if err != nil {
		return false, nil, err
	}

	var vals []sfvm.ValueOperator

	for name, value := range prog.Program.ExtraFile {
		log.Infof("start to handling: %v len: %v", name, len(value))

		if len(value) <= 128 {
			editor, _ := ssadb.GetIrSourceFromHash(value)
			if editor != nil {
				value = editor.GetSourceCode()
			}
		}

		mapperStack := utils.NewStack[*mybatisXMLMapper]()
		mapperStack.Push(&mybatisXMLMapper{
			entityStack: utils.NewStack[*mybatisXMLQuery](),
		})
		xml2.Handle(value, xml2.WithDirectiveHandler(func(directive xml.Directive) bool {
			if utils.MatchAnyOfSubString(string(directive), "dtd/mybatis-", "mybatis.org") {
				return true
			}
			return false
		}), xml2.WithStartElementHandler(func(element xml.StartElement) {
			if element.Name.Local == "mapper" {
				mapperStack.Push(&mybatisXMLMapper{
					entityStack: utils.NewStack[*mybatisXMLQuery](),
				})
				mapper := mapperStack.Peek()
				for _, attr := range element.Attr {
					if attr.Name.Local == "namespace" {
						mapper.Namespace = attr.Value
						mapper.FullClassName = attr.Value
						idx := strings.LastIndex(attr.Value, ".")
						if idx > 0 {
							mapper.ClassName = attr.Value[idx+1:]
						} else {
							mapper.ClassName = attr.Value
						}
					}
				}
				return
			}

			if element.Name.Local == "resultMap" {
				return
			}

			i := mapperStack.Peek()
			if utils.IsNil(i) {
				return
			}
			i.entityStack.Push(&mybatisXMLQuery{mapper: mapperStack.Peek()})
			query := i.entityStack.Peek()
			for _, attr := range element.Attr {
				if attr.Name.Local == "id" {
					query.Id = attr.Value
				}
			}
		}), xml2.WithEndElementHandler(func(element xml.EndElement) {
			if element.Name.Local == "mapper" {
				mapper := mapperStack.Pop()
				if mapper != nil {
					log.Infof("mapper: %v", mapper)
				}
			}

			if element.Name.Local == "resultMap" {
				return
			}

			i := mapperStack.Peek()
			if utils.IsNil(i) {
				return
			}

			i.entityStack.Pop()
		}), xml2.WithCharDataHandler(func(data xml.CharData) {
			i := mapperStack.Peek()
			if utils.IsNil(i) {
				return
			}
			query := i.entityStack.Peek()
			if query != nil {
				for _, groups := range mybatisVarExtractor.FindAllStringSubmatch(string(data), -1) {
					variableName := groups[1]
					query.CheckParamName = append(query.CheckParamName, variableName)
				}
				if len(query.CheckParamName) > 0 {
					for _, sf := range []string{
						query.SyntaxFlowFirst(), query.SyntaxFlowFinal(),
					} {
						if sf == "" {
							continue
						}
						val := prog.NewValue(ssa.NewConst(sf))
						_ = val.AppendPredecessor(v, frame.WithPredecessorContext("mybatis-${...}"))
						log.Infof("mybatis-${...}: fetch query: %v", sf)
						_, _, err := nativeCallEval(val, frame, nil)
						if err != nil {
							log.Warnf("mybatis-${...}: fetch query: %v, error: %v", sf, err)
						}
						results, ok := frame.GetSymbolTable().Get("__output")
						haveResult := false
						if ok {
							_ = results.Recursive(func(operator sfvm.ValueOperator) error {
								haveResult = true
								vals = append(vals, operator)
								return nil
							})
						}
						if haveResult {
							break
						}
					}
				}
			}
		}))
	}

	return true, sfvm.NewValues(vals), nil
}
