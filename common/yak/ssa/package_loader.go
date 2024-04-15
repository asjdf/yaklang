package ssa

import (
	"github.com/yaklang/yaklang/common/utils"
	"path"
	"strings"
)

func (b *FunctionBuilder) AddIncludePath(path string) {
	p := b.GetProgram()
	p.loader.AddIncludePath(path)
}

func (b *FunctionBuilder) BuildFilePackage(filename string, once bool) error {
	p := b.GetProgram()
	file, data, err := p.loader.LoadFilePackage(filename, once)
	if err != nil {
		return err
	}
	_file := b.CurrentFile
	_path := p.loader.GetCurrentPath()

	b.CurrentFile = file
	p.loader.SetCurrentPath(path.Dir(file))

	defer func() {
		b.CurrentFile = _file
		p.loader.SetCurrentPath(_path)
	}()
	err = p.Build(utils.UnsafeBytesToString(data), b)
	return err
}

func (b *FunctionBuilder) BuildDirectoryPackage(name string, once bool) error {
	p := b.GetProgram()
	ch, err := p.loader.LoadDirectoryPackage(name, once)
	if err != nil {
		return err
	}
	for v := range ch {

		_file := b.CurrentFile
		b.CurrentFile = v.PathName

		_path := p.loader.GetCurrentPath()
		p.loader.SetCurrentPath(path.Dir(_file))
		err := p.Build(utils.UnsafeBytesToString(v.Data), b)
		b.CurrentFile = _file
		p.loader.SetCurrentPath(_path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *FunctionBuilder) AddCurrentPackagePath(path []string) *FunctionBuilder {
	p := b.GetProgram()
	p.loader.AddPackagePath(path)

	pkgName := strings.Join(path, ".")
	if pkgName != "" {
		if !p.IsPackagePathInList(pkgName) {
			p.AddPackage(NewPackage(pkgName))
			p.packagePathList = append(p.packagePathList, path)
		}
		return p.GetAndCreateFunctionBuilder(pkgName, "init")
	}
	return nil

}

func (b *FunctionBuilder) GetCurrentPackagePath() []string {
	p := b.GetProgram()
	pkgPath := p.loader.GetPackagePath()
	return pkgPath
}
