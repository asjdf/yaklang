package go2ssa

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/yaklang/yaklang/common/consts"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
	fi "github.com/yaklang/yaklang/common/utils/filesys/filesys_interface"
	"github.com/yaklang/yaklang/common/yak/antlr4util"
	"github.com/yaklang/yaklang/common/yak/ssa"

	gol "github.com/yaklang/yaklang/common/yak/antlr4go/parser"
)

type SSABuilder struct {
	ssa.DummyPreHandler
}

var Builder = &SSABuilder{}

func (s *SSABuilder) PreHandlerProject(fileSystem fi.FileSystem, functionBuilder *ssa.FunctionBuilder, path string) error {
	prog := functionBuilder.GetProgram()
	if prog == nil {
		log.Errorf("program is nil")
		return nil
	}
	if prog.ExtraFile == nil {
		prog.ExtraFile = make(map[string]string)
	}

	dirname, filename := fileSystem.PathSplit(path)
	_ = dirname
	_ = filename

	// go.mod
	if strings.TrimLeft(filename, string(fileSystem.GetSeparators())) == "go.mod" {
		raw, err := fileSystem.ReadFile(path)
		if err != nil {
			log.Warnf("read go.mod error: %v", err)
			return nil
		}
		text := string(raw)
		pattern := `module(.*?)\n`
		re, err := regexp.Compile(pattern)
		if err != nil {
			log.Warnf("compile regexp error: %v", err)
			return nil
		}
		matches := re.FindAllString(text, -1)
		matche := strings.Split(matches[0], " ")
		if len(matches) > 0 {
			path := matche[1]
			prog.ExtraFile["go.mod"] = path[:len(path)-1]
		}
	}

	return nil
}

func (*SSABuilder) Build(src string, force bool, builder *ssa.FunctionBuilder) error {
	ast, err := Frontend(src, force)
	if err != nil {
		return err
	}
	builder.SupportClosure = true
	astBuilder := &astbuilder{
		FunctionBuilder: builder,
		cmap:            []map[string]struct{}{},
		globalv:         map[string]ssa.Value{},
		structTypes:     map[string]*ssa.ObjectType{},
		aliasTypes:      map[string]*ssa.AliasType{},
		result:          []string{},
		extendFuncs:     map[string]map[string]*ssa.Function{},
		tpHander:        map[string]func(){},
	}
	log.Infof("ast: %s", ast.ToStringTree(ast.GetParser().GetRuleNames(), ast.GetParser()))
	astBuilder.build(ast)
	return nil
}

func (*SSABuilder) FilterFile(path string) bool {
	return filepath.Ext(path) == ".go"
}

type astbuilder struct {
	*ssa.FunctionBuilder
	cmap        []map[string]struct{}
	globalv     map[string]ssa.Value
	structTypes map[string]*ssa.ObjectType
	aliasTypes  map[string]*ssa.AliasType
	result      []string
	extendFuncs map[string]map[string]*ssa.Function
	tpHander    map[string]func()
}

func Frontend(src string, must bool) (*gol.SourceFileContext, error) {
	errListener := antlr4util.NewErrorListener()
	lexer := gol.NewGoLexer(antlr.NewInputStream(src))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errListener)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := gol.NewGoParser(tokenStream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errListener)
	parser.SetErrorHandler(antlr.NewDefaultErrorStrategy())
	ast := parser.SourceFile().(*gol.SourceFileContext)
	if must || len(errListener.GetErrors()) == 0 {
		return ast, nil
	}
	return nil, utils.Errorf("parse AST FrontEnd error : %v", errListener.GetErrorString())
}

func (b *astbuilder) AddToCmap(key string) {
	b.cmap[len(b.cmap)-1][key] = struct{}{}
}

func (b *astbuilder) GetFromCmap(key string) bool {
	for _, m := range b.cmap {
		if _, ok := m[key]; ok {
			return true
		}
	}
	return false
}

func (b *astbuilder) InCmapLevel() {
	b.cmap = append(b.cmap, make(map[string]struct{}))
}

func (b *astbuilder) OutCmapLevel() {
	b.cmap = b.cmap[:len(b.cmap)-1]
}

func (*SSABuilder) GetLanguage() consts.Language {
	return consts.GO
}

func (b *astbuilder) AddGlobalVariable(name string, v ssa.Value) {
	b.globalv[name] = v
}

func (b *astbuilder) GetGlobalVariable(name string) ssa.Value {
	if b.globalv[name] == nil {
		return nil
	}
	return b.globalv[name]
}

func (b *astbuilder) AddResultDefault(name string) {
	b.result = append(b.result, name)
}

func (b *astbuilder) GetResultDefault() []string {
	return b.result
}

func (b *astbuilder) CleanResultDefault() {
	b.result = []string{}
}

func (b *astbuilder) AddExtendFuncs(name string, funcs map[string]*ssa.Function) {
	b.extendFuncs[name] = funcs
}

func (b *astbuilder) GetExtendFuncs(name string) map[string]*ssa.Function {
	if b.extendFuncs[name] == nil {
		return nil
	}
	return b.extendFuncs[name]
}

// ====================== Object type
func (b *astbuilder) AddStruct(name string, t *ssa.ObjectType) {
	b.structTypes[name] = t
}

func (b *astbuilder) GetStructByStr(name string) *ssa.ObjectType {
	if b.structTypes[name] == nil {
		return nil
	}
	return b.structTypes[name]
}

func (b *astbuilder) GetStructAll() map[string]*ssa.ObjectType {
	return b.structTypes
}

// ====================== Alias type
func (b *astbuilder) AddAlias(name string, t *ssa.AliasType) {
	b.aliasTypes[name] = t
}

func (b *astbuilder) DelAliasByStr(name string) {
	delete(b.aliasTypes, name)
}

func (b *astbuilder) GetAliasByStr(name string) ssa.Type {
	if b.aliasTypes[name] == nil {
		return nil
	}
	return b.aliasTypes[name]
}
