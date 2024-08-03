package ssa

import (
	"github.com/yaklang/yaklang/common/consts"
	fi "github.com/yaklang/yaklang/common/utils/filesys/filesys_interface"
)

type ExtraFileAnalyzer interface {
	EnableExtraFileAnalyzer() bool
	ProgramHandler(fi.FileSystem, *FunctionBuilder, string) error
}

type Builder interface {
	Build(string, bool, *FunctionBuilder) error
	FilterFile(string) bool
	GetLanguage() consts.Language
	ExtraFileAnalyzer
}

var _ ExtraFileAnalyzer = &DummyExtraFileAnalyzer{}

type DummyExtraFileAnalyzer struct{}

func (d *DummyExtraFileAnalyzer) EnableExtraFileAnalyzer() bool {
	return false
}

func (d *DummyExtraFileAnalyzer) ProgramHandler(fi.FileSystem, *FunctionBuilder, string) error {
	return nil
}
