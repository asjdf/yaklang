package ssa

import (
	"sort"
	"strings"

	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/sca/dxtypes"
	"github.com/yaklang/yaklang/common/utils"
	"golang.org/x/exp/slices"

	fi "github.com/yaklang/yaklang/common/utils/filesys/filesys_interface"
	"github.com/yaklang/yaklang/common/utils/memedit"
	"github.com/yaklang/yaklang/common/utils/omap"
	"github.com/yaklang/yaklang/common/yak/ssa/ssautil"
)

func NewChildProgram(prog *Program, name string, add bool) *Program {
	program := NewProgram(name, prog.EnableDatabase, ChildAPP, prog.Loader.GetFilesysFileSystem(), prog.Loader.GetBasePath())
	program.LibraryFile = prog.LibraryFile
	program.Language = prog.Language
	program.Application = prog
	program.FileList = prog.FileList
	program.editorStack = prog.editorStack
	program.Loader = prog.Loader
	program.Build = prog.Build
	program.externType = prog.externType
	program.externBuildValueHandler = prog.externBuildValueHandler
	program.ExternInstance = prog.ExternInstance
	program.ExternLib = prog.ExternLib
	program.GlobalScope = prog.GlobalScope
	program.ScopeCallback = prog.ScopeCallback
	program.Cache.SetFetchId(prog.Cache.fetchId)
	if add {
		prog.ChildApplication = append(prog.ChildApplication, program)
	}
	return program
}

func NewProgram(ProgramName string, enableDatabase bool, kind ProgramKind, fs fi.FileSystem, programPath string) *Program {
	prog := &Program{
		ChildApplication:        make([]*Program, 0),
		Name:                    ProgramName,
		ProgramKind:             kind,
		LibraryFile:             make(map[string][]string),
		UpStream:                make(map[string]*Program),
		DownStream:              make(map[string]*Program),
		errors:                  make([]*SSAError, 0),
		Cache:                   NewDBCache(ProgramName, enableDatabase),
		OffsetMap:               make(map[int]*OffsetItem),
		OffsetSortedSlice:       make([]int, 0),
		Funcs:                   make(map[string]*Function),
		ClassBluePrint:          make(map[string]*ClassBluePrint),
		editorStack:             omap.NewOrderedMap(make(map[string]*memedit.MemEditor)),
		editorMap:               omap.NewOrderedMap(make(map[string]*memedit.MemEditor)),
		FileList:                make(map[string]string),
		cacheExternInstance:     make(map[string]Value),
		externType:              make(map[string]Type),
		externBuildValueHandler: make(map[string]func(b *FunctionBuilder, id string, v any) (value Value)),
		ExternInstance:          make(map[string]any),
		ExternLib:               make(map[string]map[string]any),
		ScopeCallback:           func(scope ScopeIF) ScopeIF { return scope },
	}
	if kind == Application {
		prog.Application = prog
	}
	prog.EnableDatabase = enableDatabase
	prog.Loader = ssautil.NewPackageLoader(
		ssautil.WithFileSystem(fs),
		ssautil.WithIncludePath(programPath),
		ssautil.WithBasePath(programPath),
	)
	return prog
}

func (prog *Program) GetLibrary(name string) (*Program, bool) {
	if prog == nil || utils.IsNil(prog) || prog.Application == nil || utils.IsNil(prog.Application) {
		return nil, false
	}
	// get lib from application
	app := prog.Application
	currentEditor := prog.GetCurrentEditor()
	// this program has current file
	hasFile := func(p *Program) bool {
		if hash, ok := p.FileList[currentEditor.GetFilename()]; ok {
			if hash == currentEditor.SourceCodeMd5() {
				return true
			}
		}
		return false
	}

	// contain in memory
	if p, ok := app.UpStream[name]; ok {
		return p, hasFile(p)
	}

	if p, ok := prog.UpStream[name]; ok {
		app.AddUpStream(p)
		return p, hasFile(p)
	}

	if !app.EnableDatabase {
		return nil, false
	}

	version := ""
	if p := app.GetSCAPackageByName(name); p != nil {
		version = p.Version
	} else {
		return nil, false
	}
	// library in  database, load and set relation
	p, err := GetLibrary(name, version)
	if err != nil {
		return nil, false
	}
	app.AddUpStream(p)
	if !slices.Contains(p.irProgram.UpStream, name) {
		// update up-down stream
		prog.AddUpStream(p)
	}
	return p, hasFile(p)
}

func (prog *Program) AddUpStream(p *Program) {
	prog.UpStream[p.Name] = p
	p.DownStream[prog.Name] = prog
}

func (prog *Program) NewLibrary(name string, path []string) *Program {
	// create lib
	// get program Path
	fs := prog.Loader.GetFilesysFileSystem()
	fullPath := prog.GetCurrentEditor().GetFilename()
	endPath := fs.Join(path...)
	programPath, _, _ := strings.Cut(fullPath, endPath)

	lib := NewProgram(name, prog.EnableDatabase, Library, fs, programPath)
	lib.Loader.AddIncludePath(prog.Loader.GetIncludeFiles()...)
	lib.Language = prog.Language

	//todo: 这里需要加一个测试
	lib.GlobalScope = prog.GlobalScope

	// up-down stream and application
	prog.AddUpStream(lib)
	prog.Application.AddUpStream(lib)
	lib.Application = prog.Application
	lib.ExternLib = prog.ExternLib
	lib.ExternInstance = prog.ExternInstance
	return lib
}

func (prog *Program) GetProgramName() string {
	return prog.Name
}

func (prog *Program) GetAndCreateFunction(pkgName string, funcName string) *Function {
	fun := prog.GetFunction(funcName)
	if fun == nil {
		fun = prog.NewFunction(funcName)
	}

	if fun.GetRange() == nil {
		// if editor := prog.getCurrentEditor(); editor != nil {
		// 	fun.SetRangeInit(editor)
		// } else {
		log.Warnf("the program must contains a editor to init function range: %v", prog.Name)
		// }
	}

	return fun
}

func (prog *Program) GetCacheExternInstance(name string) (Value, bool) {
	v, ok := prog.cacheExternInstance[name]
	return v, ok
}

func (prog *Program) SetCacheExternInstance(name string, v Value) {
	prog.cacheExternInstance[name] = v
}

// create or get main function builder
func (prog *Program) GetAndCreateFunctionBuilder(pkgName string, funcName string) *FunctionBuilder {
	fun := prog.GetAndCreateFunction(pkgName, funcName)
	builder := fun.builder
	if builder == nil {
		builder = NewBuilder(prog.GetCurrentEditor(), fun, nil)
	}

	return builder
}

func (p *Program) GetFunction(name string) *Function {
	if f, ok := p.Funcs[name]; ok {
		return f
	}
	return nil
}

func (prog *Program) EachFunction(handler func(*Function)) {
	var handFunc func(*Function)
	handFunc = func(f *Function) {
		handler(f)
		for _, s := range f.ChildFuncs {
			f, ok := ToFunction(s)
			if !ok {
				log.Warnf("function %s is not a ssa.Function", s.GetName())
				continue
			}
			handFunc(f)
		}
	}

	for _, f := range prog.Funcs {
		handFunc(f)
	}
}

func (prog *Program) Finish() {
	finishOnce := func() {
		for _, program := range prog.ChildApplication {
			program.Finish()
		}
		for _, up := range prog.UpStream {
			up.Finish()
		}
		prog.Cache.SaveToDatabase()
		if prog.EnableDatabase {
			updateToDatabase(prog)
		}
	}
	prog.finishOnce.Do(finishOnce)
}

func (prog *Program) SearchIndexAndOffsetByOffset(searchOffset int) (index int, offset int) {
	index = sort.Search(len(prog.OffsetSortedSlice), func(i int) bool {
		return prog.OffsetSortedSlice[i] >= searchOffset
	})
	if index >= len(prog.OffsetSortedSlice) && len(prog.OffsetSortedSlice) > 0 {
		index = len(prog.OffsetSortedSlice) - 1
	}
	if len(prog.OffsetSortedSlice) > 0 {
		offset = prog.OffsetSortedSlice[index]
	}
	return
}

func (prog *Program) GetFrontValueByOffset(searchOffset int) (offset int, value Value) {
	index, offset := prog.SearchIndexAndOffsetByOffset(searchOffset)
	// 如果二分查找的结果是大于目标值的，那么就需要回退一个
	if offset > searchOffset {
		if index > 0 {
			index -= 1
		}
		offset = prog.OffsetSortedSlice[index]
	}
	if item, ok := prog.OffsetMap[offset]; ok {
		value = item.GetValue()
	}
	return offset, value
}

func (p *Program) GetEditor(url string) (*memedit.MemEditor, bool) {
	return p.editorMap.Get(url)
}

func (p *Program) PushEditor(e *memedit.MemEditor) {
	p.PushEditorex(e, true)
}
func (p *Program) PushEditorex(e *memedit.MemEditor, store bool) {
	p.editorStack.Push(e)
	if store {
		p.editorMap.Set(p.GetCurrentEditor().GetFilename(), p.GetCurrentEditor())
	}
}

func (p *Program) GetIncludeFiles() []string {
	return p.editorMap.Keys()
}
func (p *Program) GetIncludeFileNum() int {
	return p.editorMap.Len()
}

func (p *Program) GetCurrentEditor() *memedit.MemEditor {
	if p.editorStack == nil || p.editorStack.Len() <= 0 {
		return nil
	}
	_, v, ok := p.editorStack.Last()
	if !ok {
		return nil
	}
	return v
}

func (p *Program) PopEditor() {
	if p.editorStack == nil || p.editorStack.Len() <= 0 {
		return
	}
	e := p.editorStack.Pop()
	if e != nil {
		p.FileList[e.GetFilename()] = e.SourceCodeMd5()
	}
}

func (p *Program) GetSCAPackageByName(name string) *dxtypes.Package {
	if p == nil {
		return nil
	}
	for _, pkg := range p.SCAPackages {
		if strings.Contains(pkg.Name, name) {
			return pkg
		}
	}
	return nil
}

func (p *Program) GetApplication() *Program {
	if p == nil {
		return nil
	}
	return p.Application
}

func (p *Program) GetType(name string) Type {
	if t, ok := p.externType[name]; ok {
		return t
	}
	return nil
}
