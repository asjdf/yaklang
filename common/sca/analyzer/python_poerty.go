package analyzer

import (
	"path"
	"strings"

	"github.com/aquasecurity/go-dep-parser/pkg/python/poetry"
	"github.com/aquasecurity/go-dep-parser/pkg/python/pyproject"
	"github.com/yaklang/yaklang/common/sca/types"
)

const (
	TypPythonPoetry TypAnalyzer = "python-poetry-lang"

	PoetryLockFile = "poetry.lock"
	PyProjectFile  = "pyproject.toml"

	statusPoetry int = 1
)

func init() {
	RegisterAnalyzer(TypPythonPoetry, NewPythonPoetryAnalyzer())
}

type pythonPoetryAnalyzer struct{}

func NewPythonPoetryAnalyzer() *pythonPoetryAnalyzer {
	return &pythonPoetryAnalyzer{}
}

func (a pythonPoetryAnalyzer) Match(info MatchInfo) int {
	if strings.HasSuffix(info.path, PoetryLockFile) || strings.HasSuffix(info.path, PyProjectFile) {
		return statusPIP
	}
	return 0
}

func (a pythonPoetryAnalyzer) Analyze(afi AnalyzeFileInfo) ([]types.Package, error) {
	fi := afi.self

	switch fi.matchStatus {
	case statusPoetry:
		pkgs, err := ParseLanguageConfiguration(fi, poetry.NewParser())
		if err != nil {
			return nil, err
		}
		if pkgs == nil {
			return nil, nil
		}

		// Parse pyproject.toml to identify the direct dependencies
		pyprojectPath := path.Join(path.Dir(fi.path), PyProjectFile)
		if pyprojectFi, ok := afi.matchedFileInfos[pyprojectPath]; ok {
			pyProjectParser := pyproject.NewParser()
			parsed, err := pyProjectParser.Parse(pyprojectFi.f)
			if err != nil {
				return nil, err
			}
			for i, pkg := range pkgs {
				// Identify the direct/transitive dependencies
				if _, ok := parsed[pkg.Name]; !ok {
					pkgs[i].Indirect = true
				}
			}
		}

		return pkgs, nil
	}

	return nil, nil
}
