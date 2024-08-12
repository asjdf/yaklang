package php2ssa

import (
	phpparser "github.com/yaklang/yaklang/common/yak/php/parser"
)

func (y *builder) VisitPhpBlock(raw phpparser.IPhpBlockContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}
	recoverRange := y.SetRange(raw)
	defer recoverRange()

	i, _ := raw.(*phpparser.PhpBlockContext)
	if i == nil {
		return nil
	}

	// import? what the fuck?
	if len(i.AllImportStatement()) > 0 {
		// handle ImportStmt
	}
	if y.MoreParse && len(i.AllNamespaceDeclaration()) == 0 {
		return nil
	}
	for _, namespace := range i.AllNamespaceDeclaration() {
		y.VisitNamespaceDeclaration(namespace)
	}
	for _, usedecl := range i.AllUseDeclaration() {
		y.VisitUseDeclaration(usedecl)
	}
	for _, global := range i.AllGlobalConstantDeclaration() {
		y.VisitGlobalConstantDeclaration(global)
	}
	for _, functiondecl := range i.AllFunctionDeclaration() {
		y.VisitFunctionDeclaration(functiondecl)
	}
	for _, classdecl := range i.AllClassDeclaration() {
		y.VisitClassDeclaration(classdecl)
	}
	for _, stmt := range i.AllStatement() {
		y.VisitStatement(stmt)
	}
	for _, enum := range i.AllEnumDeclaration() {
		y.VisitEnumDeclaration(enum)
	}
	return nil
}
