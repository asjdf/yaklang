// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package sf

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type SyntaxFlowLexer struct {
	SyntaxFlowBaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var syntaxflowlexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func syntaxflowlexerLexerInit() {
	staticData := &syntaxflowlexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE", "HereDocIdentifier", "CRLFHereDoc", "LFHereDoc",
	}
	staticData.literalNames = []string{
		"", "'==>'", "'...'", "'%%'", "'..'", "'<='", "'>='", "'>>'", "'=>'",
		"'=='", "'=~'", "'!~'", "'&&'", "'||'", "'!='", "'${'", "';'", "'?{'",
		"'-{'", "'->'", "'}->'", "'-->'", "'#{'", "'#>'", "'#->'", "'>'", "'.'",
		"'<<<'", "'<'", "'='", "'+'", "'&'", "'?'", "'('", "','", "')'", "'['",
		"']'", "'{'", "'}'", "'#'", "'$'", "':'", "'%'", "'!'", "'*'", "'-'",
		"'as'", "'`'", "'''", "'\"'", "", "'\\n'", "", "", "", "", "", "'str'",
		"'list'", "'dict'", "", "'bool'", "", "'alert'", "'check'", "'then'",
		"", "'else'", "'type'", "'in'", "'call'", "", "", "'phi'", "", "", "'opcode'",
		"'have'", "'any'", "'not'", "'for'", "'r'", "'g'", "'e'",
	}
	staticData.symbolicNames = []string{
		"", "DeepFilter", "Deep", "Percent", "DeepDot", "LtEq", "GtEq", "DoubleGt",
		"Filter", "EqEq", "RegexpMatch", "NotRegexpMatch", "And", "Or", "NotEq",
		"DollarBraceOpen", "Semicolon", "ConditionStart", "DeepNextStart", "UseStart",
		"DeepNextEnd", "DeepNext", "TopDefStart", "DefStart", "TopDef", "Gt",
		"Dot", "StartNowDoc", "Lt", "Eq", "Add", "Amp", "Question", "OpenParen",
		"Comma", "CloseParen", "ListSelectOpen", "ListSelectClose", "MapBuilderOpen",
		"MapBuilderClose", "ListStart", "DollarOutput", "Colon", "Search", "Bang",
		"Star", "Minus", "As", "Backtick", "SingleQuote", "DoubleQuote", "LineComment",
		"BreakLine", "WhiteSpace", "Number", "OctalNumber", "BinaryNumber",
		"HexNumber", "StringType", "ListType", "DictType", "NumberType", "BoolType",
		"BoolLiteral", "Alert", "Check", "Then", "Desc", "Else", "Type", "In",
		"Call", "Function", "Constant", "Phi", "FormalParam", "Return", "Opcode",
		"Have", "HaveAny", "Not", "For", "ConstSearchModePrefixRegexp", "ConstSearchModePrefixGlob",
		"ConstSearchModePrefixExact", "Identifier", "IdentifierChar", "QuotedStringLiteral",
		"RegexpLiteral", "WS", "HereDocIdentifierName", "CRLFHereDocIdentifierBreak",
		"LFHereDocIdentifierBreak", "CRLFEndDoc", "CRLFHereDocText", "LFEndDoc",
		"LFHereDocText",
	}
	staticData.ruleNames = []string{
		"DeepFilter", "Deep", "Percent", "DeepDot", "LtEq", "GtEq", "DoubleGt",
		"Filter", "EqEq", "RegexpMatch", "NotRegexpMatch", "And", "Or", "NotEq",
		"DollarBraceOpen", "Semicolon", "ConditionStart", "DeepNextStart", "UseStart",
		"DeepNextEnd", "DeepNext", "TopDefStart", "DefStart", "TopDef", "Gt",
		"Dot", "StartNowDoc", "Lt", "Eq", "Add", "Amp", "Question", "OpenParen",
		"Comma", "CloseParen", "ListSelectOpen", "ListSelectClose", "MapBuilderOpen",
		"MapBuilderClose", "ListStart", "DollarOutput", "Colon", "Search", "Bang",
		"Star", "Minus", "As", "Backtick", "SingleQuote", "DoubleQuote", "LineComment",
		"BreakLine", "WhiteSpace", "Number", "OctalNumber", "BinaryNumber",
		"HexNumber", "StringType", "ListType", "DictType", "NumberType", "BoolType",
		"BoolLiteral", "Alert", "Check", "Then", "Desc", "Else", "Type", "In",
		"Call", "Function", "Constant", "Phi", "FormalParam", "Return", "Opcode",
		"Have", "HaveAny", "Not", "For", "ConstSearchModePrefixRegexp", "ConstSearchModePrefixGlob",
		"ConstSearchModePrefixExact", "Identifier", "IdentifierChar", "QuotedStringLiteral",
		"RegexpLiteral", "WS", "HereDocIdentifierName", "CRLFHereDocIdentifierBreak",
		"LFHereDocIdentifierBreak", "CRLFEndDoc", "CRLFHereDocText", "LFEndDoc",
		"LFHereDocText", "NameString", "IdentifierCharStart", "HexDigit", "Digit",
		"OctalDigit", "RegexpLiteralChar",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 96, 679, 6, -1, 6, -1, 6, -1, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2,
		7, 2, 2, 3, 7, 3, 2, 4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8,
		7, 8, 2, 9, 7, 9, 2, 10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13,
		2, 14, 7, 14, 2, 15, 7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2,
		19, 7, 19, 2, 20, 7, 20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24,
		7, 24, 2, 25, 7, 25, 2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7,
		29, 2, 30, 7, 30, 2, 31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34,
		2, 35, 7, 35, 2, 36, 7, 36, 2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2,
		40, 7, 40, 2, 41, 7, 41, 2, 42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45,
		7, 45, 2, 46, 7, 46, 2, 47, 7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7,
		50, 2, 51, 7, 51, 2, 52, 7, 52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55,
		2, 56, 7, 56, 2, 57, 7, 57, 2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2,
		61, 7, 61, 2, 62, 7, 62, 2, 63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66,
		7, 66, 2, 67, 7, 67, 2, 68, 7, 68, 2, 69, 7, 69, 2, 70, 7, 70, 2, 71, 7,
		71, 2, 72, 7, 72, 2, 73, 7, 73, 2, 74, 7, 74, 2, 75, 7, 75, 2, 76, 7, 76,
		2, 77, 7, 77, 2, 78, 7, 78, 2, 79, 7, 79, 2, 80, 7, 80, 2, 81, 7, 81, 2,
		82, 7, 82, 2, 83, 7, 83, 2, 84, 7, 84, 2, 85, 7, 85, 2, 86, 7, 86, 2, 87,
		7, 87, 2, 88, 7, 88, 2, 89, 7, 89, 2, 90, 7, 90, 2, 91, 7, 91, 2, 92, 7,
		92, 2, 93, 7, 93, 2, 94, 7, 94, 2, 95, 7, 95, 2, 96, 7, 96, 2, 97, 7, 97,
		2, 98, 7, 98, 2, 99, 7, 99, 2, 100, 7, 100, 2, 101, 7, 101, 1, 0, 1, 0,
		1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3,
		1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7,
		1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11,
		1, 11, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1,
		15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18,
		1, 19, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1, 20, 1, 21, 1, 21, 1,
		21, 1, 22, 1, 22, 1, 22, 1, 23, 1, 23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25,
		1, 25, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 27, 1, 27, 1, 28, 1,
		28, 1, 29, 1, 29, 1, 30, 1, 30, 1, 31, 1, 31, 1, 32, 1, 32, 1, 33, 1, 33,
		1, 34, 1, 34, 1, 35, 1, 35, 1, 36, 1, 36, 1, 37, 1, 37, 1, 38, 1, 38, 1,
		39, 1, 39, 1, 40, 1, 40, 1, 41, 1, 41, 1, 42, 1, 42, 1, 43, 1, 43, 1, 44,
		1, 44, 1, 45, 1, 45, 1, 46, 1, 46, 1, 46, 1, 47, 1, 47, 1, 48, 1, 48, 1,
		49, 1, 49, 1, 50, 1, 50, 1, 50, 1, 50, 5, 50, 346, 8, 50, 10, 50, 12, 50,
		349, 9, 50, 1, 51, 1, 51, 1, 52, 1, 52, 1, 52, 1, 52, 1, 53, 4, 53, 358,
		8, 53, 11, 53, 12, 53, 359, 1, 54, 1, 54, 1, 54, 1, 54, 4, 54, 366, 8,
		54, 11, 54, 12, 54, 367, 1, 55, 1, 55, 1, 55, 1, 55, 4, 55, 374, 8, 55,
		11, 55, 12, 55, 375, 1, 56, 1, 56, 1, 56, 1, 56, 4, 56, 382, 8, 56, 11,
		56, 12, 56, 383, 1, 57, 1, 57, 1, 57, 1, 57, 1, 58, 1, 58, 1, 58, 1, 58,
		1, 58, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 60, 1, 60, 1, 60, 1, 60, 1,
		60, 1, 60, 1, 60, 1, 60, 3, 60, 408, 8, 60, 1, 61, 1, 61, 1, 61, 1, 61,
		1, 61, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 3,
		62, 424, 8, 62, 1, 63, 1, 63, 1, 63, 1, 63, 1, 63, 1, 63, 1, 64, 1, 64,
		1, 64, 1, 64, 1, 64, 1, 64, 1, 65, 1, 65, 1, 65, 1, 65, 1, 65, 1, 66, 1,
		66, 1, 66, 1, 66, 1, 66, 1, 66, 1, 66, 1, 66, 3, 66, 451, 8, 66, 1, 67,
		1, 67, 1, 67, 1, 67, 1, 67, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 1, 69, 1,
		69, 1, 69, 1, 70, 1, 70, 1, 70, 1, 70, 1, 70, 1, 71, 1, 71, 1, 71, 1, 71,
		1, 71, 1, 71, 1, 71, 1, 71, 1, 71, 1, 71, 1, 71, 1, 71, 3, 71, 483, 8,
		71, 1, 72, 1, 72, 1, 72, 1, 72, 1, 72, 1, 72, 1, 72, 1, 72, 1, 72, 1, 72,
		1, 72, 1, 72, 1, 72, 3, 72, 498, 8, 72, 1, 73, 1, 73, 1, 73, 1, 73, 1,
		74, 1, 74, 1, 74, 1, 74, 1, 74, 1, 74, 1, 74, 1, 74, 1, 74, 1, 74, 1, 74,
		1, 74, 1, 74, 1, 74, 1, 74, 1, 74, 1, 74, 3, 74, 521, 8, 74, 1, 75, 1,
		75, 1, 75, 1, 75, 1, 75, 1, 75, 1, 75, 1, 75, 1, 75, 3, 75, 532, 8, 75,
		1, 76, 1, 76, 1, 76, 1, 76, 1, 76, 1, 76, 1, 76, 1, 77, 1, 77, 1, 77, 1,
		77, 1, 77, 1, 78, 1, 78, 1, 78, 1, 78, 1, 79, 1, 79, 1, 79, 1, 79, 1, 80,
		1, 80, 1, 80, 1, 80, 1, 81, 1, 81, 1, 82, 1, 82, 1, 83, 1, 83, 1, 84, 1,
		84, 5, 84, 566, 8, 84, 10, 84, 12, 84, 569, 9, 84, 1, 85, 1, 85, 3, 85,
		573, 8, 85, 1, 86, 1, 86, 1, 86, 1, 86, 1, 86, 1, 86, 1, 86, 5, 86, 582,
		8, 86, 10, 86, 12, 86, 585, 9, 86, 1, 86, 1, 86, 1, 86, 1, 86, 1, 86, 1,
		86, 1, 86, 1, 86, 1, 86, 5, 86, 596, 8, 86, 10, 86, 12, 86, 599, 9, 86,
		1, 86, 1, 86, 3, 86, 603, 8, 86, 1, 87, 1, 87, 4, 87, 607, 8, 87, 11, 87,
		12, 87, 608, 1, 87, 1, 87, 1, 88, 4, 88, 614, 8, 88, 11, 88, 12, 88, 615,
		1, 88, 1, 88, 1, 89, 1, 89, 1, 89, 1, 89, 1, 89, 1, 89, 1, 89, 1, 89, 1,
		89, 3, 89, 629, 8, 89, 1, 90, 1, 90, 1, 90, 1, 90, 1, 90, 1, 90, 1, 90,
		1, 90, 1, 91, 1, 91, 1, 91, 1, 91, 1, 91, 1, 91, 1, 92, 1, 92, 1, 92, 1,
		92, 1, 92, 1, 92, 1, 93, 1, 93, 1, 94, 1, 94, 1, 94, 1, 94, 1, 95, 1, 95,
		1, 96, 1, 96, 5, 96, 661, 8, 96, 10, 96, 12, 96, 664, 9, 96, 1, 97, 3,
		97, 667, 8, 97, 1, 98, 1, 98, 1, 99, 1, 99, 1, 100, 1, 100, 1, 101, 1,
		101, 1, 101, 3, 101, 678, 8, 101, 0, 0, 102, 4, 1, 6, 2, 8, 3, 10, 4, 12,
		5, 14, 6, 16, 7, 18, 8, 20, 9, 22, 10, 24, 11, 26, 12, 28, 13, 30, 14,
		32, 15, 34, 16, 36, 17, 38, 18, 40, 19, 42, 20, 44, 21, 46, 22, 48, 23,
		50, 24, 52, 25, 54, 26, 56, 27, 58, 28, 60, 29, 62, 30, 64, 31, 66, 32,
		68, 33, 70, 34, 72, 35, 74, 36, 76, 37, 78, 38, 80, 39, 82, 40, 84, 41,
		86, 42, 88, 43, 90, 44, 92, 45, 94, 46, 96, 47, 98, 48, 100, 49, 102, 50,
		104, 51, 106, 52, 108, 53, 110, 54, 112, 55, 114, 56, 116, 57, 118, 58,
		120, 59, 122, 60, 124, 61, 126, 62, 128, 63, 130, 64, 132, 65, 134, 66,
		136, 67, 138, 68, 140, 69, 142, 70, 144, 71, 146, 72, 148, 73, 150, 74,
		152, 75, 154, 76, 156, 77, 158, 78, 160, 79, 162, 80, 164, 81, 166, 82,
		168, 83, 170, 84, 172, 85, 174, 86, 176, 87, 178, 88, 180, 89, 182, 90,
		184, 91, 186, 92, 188, 93, 190, 94, 192, 95, 194, 96, 196, 0, 198, 0, 200,
		0, 202, 0, 204, 0, 206, 0, 4, 0, 1, 2, 3, 12, 2, 0, 10, 10, 13, 13, 2,
		0, 13, 13, 32, 32, 1, 0, 48, 57, 4, 0, 10, 10, 13, 13, 39, 39, 92, 92,
		4, 0, 10, 10, 13, 13, 34, 34, 92, 92, 3, 0, 9, 9, 13, 13, 32, 32, 4, 0,
		65, 90, 95, 95, 97, 122, 128, 65534, 5, 0, 48, 57, 65, 90, 95, 95, 97,
		122, 128, 65534, 4, 0, 42, 42, 65, 90, 95, 95, 97, 122, 3, 0, 48, 57, 65,
		70, 97, 102, 1, 0, 48, 55, 1, 0, 47, 47, 697, 0, 4, 1, 0, 0, 0, 0, 6, 1,
		0, 0, 0, 0, 8, 1, 0, 0, 0, 0, 10, 1, 0, 0, 0, 0, 12, 1, 0, 0, 0, 0, 14,
		1, 0, 0, 0, 0, 16, 1, 0, 0, 0, 0, 18, 1, 0, 0, 0, 0, 20, 1, 0, 0, 0, 0,
		22, 1, 0, 0, 0, 0, 24, 1, 0, 0, 0, 0, 26, 1, 0, 0, 0, 0, 28, 1, 0, 0, 0,
		0, 30, 1, 0, 0, 0, 0, 32, 1, 0, 0, 0, 0, 34, 1, 0, 0, 0, 0, 36, 1, 0, 0,
		0, 0, 38, 1, 0, 0, 0, 0, 40, 1, 0, 0, 0, 0, 42, 1, 0, 0, 0, 0, 44, 1, 0,
		0, 0, 0, 46, 1, 0, 0, 0, 0, 48, 1, 0, 0, 0, 0, 50, 1, 0, 0, 0, 0, 52, 1,
		0, 0, 0, 0, 54, 1, 0, 0, 0, 0, 56, 1, 0, 0, 0, 0, 58, 1, 0, 0, 0, 0, 60,
		1, 0, 0, 0, 0, 62, 1, 0, 0, 0, 0, 64, 1, 0, 0, 0, 0, 66, 1, 0, 0, 0, 0,
		68, 1, 0, 0, 0, 0, 70, 1, 0, 0, 0, 0, 72, 1, 0, 0, 0, 0, 74, 1, 0, 0, 0,
		0, 76, 1, 0, 0, 0, 0, 78, 1, 0, 0, 0, 0, 80, 1, 0, 0, 0, 0, 82, 1, 0, 0,
		0, 0, 84, 1, 0, 0, 0, 0, 86, 1, 0, 0, 0, 0, 88, 1, 0, 0, 0, 0, 90, 1, 0,
		0, 0, 0, 92, 1, 0, 0, 0, 0, 94, 1, 0, 0, 0, 0, 96, 1, 0, 0, 0, 0, 98, 1,
		0, 0, 0, 0, 100, 1, 0, 0, 0, 0, 102, 1, 0, 0, 0, 0, 104, 1, 0, 0, 0, 0,
		106, 1, 0, 0, 0, 0, 108, 1, 0, 0, 0, 0, 110, 1, 0, 0, 0, 0, 112, 1, 0,
		0, 0, 0, 114, 1, 0, 0, 0, 0, 116, 1, 0, 0, 0, 0, 118, 1, 0, 0, 0, 0, 120,
		1, 0, 0, 0, 0, 122, 1, 0, 0, 0, 0, 124, 1, 0, 0, 0, 0, 126, 1, 0, 0, 0,
		0, 128, 1, 0, 0, 0, 0, 130, 1, 0, 0, 0, 0, 132, 1, 0, 0, 0, 0, 134, 1,
		0, 0, 0, 0, 136, 1, 0, 0, 0, 0, 138, 1, 0, 0, 0, 0, 140, 1, 0, 0, 0, 0,
		142, 1, 0, 0, 0, 0, 144, 1, 0, 0, 0, 0, 146, 1, 0, 0, 0, 0, 148, 1, 0,
		0, 0, 0, 150, 1, 0, 0, 0, 0, 152, 1, 0, 0, 0, 0, 154, 1, 0, 0, 0, 0, 156,
		1, 0, 0, 0, 0, 158, 1, 0, 0, 0, 0, 160, 1, 0, 0, 0, 0, 162, 1, 0, 0, 0,
		0, 164, 1, 0, 0, 0, 0, 166, 1, 0, 0, 0, 0, 168, 1, 0, 0, 0, 0, 170, 1,
		0, 0, 0, 0, 172, 1, 0, 0, 0, 0, 174, 1, 0, 0, 0, 0, 176, 1, 0, 0, 0, 0,
		178, 1, 0, 0, 0, 0, 180, 1, 0, 0, 0, 1, 182, 1, 0, 0, 0, 1, 184, 1, 0,
		0, 0, 1, 186, 1, 0, 0, 0, 2, 188, 1, 0, 0, 0, 2, 190, 1, 0, 0, 0, 3, 192,
		1, 0, 0, 0, 3, 194, 1, 0, 0, 0, 4, 208, 1, 0, 0, 0, 6, 212, 1, 0, 0, 0,
		8, 216, 1, 0, 0, 0, 10, 219, 1, 0, 0, 0, 12, 222, 1, 0, 0, 0, 14, 225,
		1, 0, 0, 0, 16, 228, 1, 0, 0, 0, 18, 231, 1, 0, 0, 0, 20, 234, 1, 0, 0,
		0, 22, 237, 1, 0, 0, 0, 24, 240, 1, 0, 0, 0, 26, 243, 1, 0, 0, 0, 28, 246,
		1, 0, 0, 0, 30, 249, 1, 0, 0, 0, 32, 252, 1, 0, 0, 0, 34, 255, 1, 0, 0,
		0, 36, 257, 1, 0, 0, 0, 38, 260, 1, 0, 0, 0, 40, 263, 1, 0, 0, 0, 42, 266,
		1, 0, 0, 0, 44, 270, 1, 0, 0, 0, 46, 274, 1, 0, 0, 0, 48, 277, 1, 0, 0,
		0, 50, 280, 1, 0, 0, 0, 52, 284, 1, 0, 0, 0, 54, 286, 1, 0, 0, 0, 56, 288,
		1, 0, 0, 0, 58, 294, 1, 0, 0, 0, 60, 296, 1, 0, 0, 0, 62, 298, 1, 0, 0,
		0, 64, 300, 1, 0, 0, 0, 66, 302, 1, 0, 0, 0, 68, 304, 1, 0, 0, 0, 70, 306,
		1, 0, 0, 0, 72, 308, 1, 0, 0, 0, 74, 310, 1, 0, 0, 0, 76, 312, 1, 0, 0,
		0, 78, 314, 1, 0, 0, 0, 80, 316, 1, 0, 0, 0, 82, 318, 1, 0, 0, 0, 84, 320,
		1, 0, 0, 0, 86, 322, 1, 0, 0, 0, 88, 324, 1, 0, 0, 0, 90, 326, 1, 0, 0,
		0, 92, 328, 1, 0, 0, 0, 94, 330, 1, 0, 0, 0, 96, 332, 1, 0, 0, 0, 98, 335,
		1, 0, 0, 0, 100, 337, 1, 0, 0, 0, 102, 339, 1, 0, 0, 0, 104, 341, 1, 0,
		0, 0, 106, 350, 1, 0, 0, 0, 108, 352, 1, 0, 0, 0, 110, 357, 1, 0, 0, 0,
		112, 361, 1, 0, 0, 0, 114, 369, 1, 0, 0, 0, 116, 377, 1, 0, 0, 0, 118,
		385, 1, 0, 0, 0, 120, 389, 1, 0, 0, 0, 122, 394, 1, 0, 0, 0, 124, 407,
		1, 0, 0, 0, 126, 409, 1, 0, 0, 0, 128, 423, 1, 0, 0, 0, 130, 425, 1, 0,
		0, 0, 132, 431, 1, 0, 0, 0, 134, 437, 1, 0, 0, 0, 136, 450, 1, 0, 0, 0,
		138, 452, 1, 0, 0, 0, 140, 457, 1, 0, 0, 0, 142, 462, 1, 0, 0, 0, 144,
		465, 1, 0, 0, 0, 146, 482, 1, 0, 0, 0, 148, 497, 1, 0, 0, 0, 150, 499,
		1, 0, 0, 0, 152, 520, 1, 0, 0, 0, 154, 531, 1, 0, 0, 0, 156, 533, 1, 0,
		0, 0, 158, 540, 1, 0, 0, 0, 160, 545, 1, 0, 0, 0, 162, 549, 1, 0, 0, 0,
		164, 553, 1, 0, 0, 0, 166, 557, 1, 0, 0, 0, 168, 559, 1, 0, 0, 0, 170,
		561, 1, 0, 0, 0, 172, 563, 1, 0, 0, 0, 174, 572, 1, 0, 0, 0, 176, 602,
		1, 0, 0, 0, 178, 604, 1, 0, 0, 0, 180, 613, 1, 0, 0, 0, 182, 628, 1, 0,
		0, 0, 184, 630, 1, 0, 0, 0, 186, 638, 1, 0, 0, 0, 188, 644, 1, 0, 0, 0,
		190, 650, 1, 0, 0, 0, 192, 652, 1, 0, 0, 0, 194, 656, 1, 0, 0, 0, 196,
		658, 1, 0, 0, 0, 198, 666, 1, 0, 0, 0, 200, 668, 1, 0, 0, 0, 202, 670,
		1, 0, 0, 0, 204, 672, 1, 0, 0, 0, 206, 677, 1, 0, 0, 0, 208, 209, 5, 61,
		0, 0, 209, 210, 5, 61, 0, 0, 210, 211, 5, 62, 0, 0, 211, 5, 1, 0, 0, 0,
		212, 213, 5, 46, 0, 0, 213, 214, 5, 46, 0, 0, 214, 215, 5, 46, 0, 0, 215,
		7, 1, 0, 0, 0, 216, 217, 5, 37, 0, 0, 217, 218, 5, 37, 0, 0, 218, 9, 1,
		0, 0, 0, 219, 220, 5, 46, 0, 0, 220, 221, 5, 46, 0, 0, 221, 11, 1, 0, 0,
		0, 222, 223, 5, 60, 0, 0, 223, 224, 5, 61, 0, 0, 224, 13, 1, 0, 0, 0, 225,
		226, 5, 62, 0, 0, 226, 227, 5, 61, 0, 0, 227, 15, 1, 0, 0, 0, 228, 229,
		5, 62, 0, 0, 229, 230, 5, 62, 0, 0, 230, 17, 1, 0, 0, 0, 231, 232, 5, 61,
		0, 0, 232, 233, 5, 62, 0, 0, 233, 19, 1, 0, 0, 0, 234, 235, 5, 61, 0, 0,
		235, 236, 5, 61, 0, 0, 236, 21, 1, 0, 0, 0, 237, 238, 5, 61, 0, 0, 238,
		239, 5, 126, 0, 0, 239, 23, 1, 0, 0, 0, 240, 241, 5, 33, 0, 0, 241, 242,
		5, 126, 0, 0, 242, 25, 1, 0, 0, 0, 243, 244, 5, 38, 0, 0, 244, 245, 5,
		38, 0, 0, 245, 27, 1, 0, 0, 0, 246, 247, 5, 124, 0, 0, 247, 248, 5, 124,
		0, 0, 248, 29, 1, 0, 0, 0, 249, 250, 5, 33, 0, 0, 250, 251, 5, 61, 0, 0,
		251, 31, 1, 0, 0, 0, 252, 253, 5, 36, 0, 0, 253, 254, 5, 123, 0, 0, 254,
		33, 1, 0, 0, 0, 255, 256, 5, 59, 0, 0, 256, 35, 1, 0, 0, 0, 257, 258, 5,
		63, 0, 0, 258, 259, 5, 123, 0, 0, 259, 37, 1, 0, 0, 0, 260, 261, 5, 45,
		0, 0, 261, 262, 5, 123, 0, 0, 262, 39, 1, 0, 0, 0, 263, 264, 5, 45, 0,
		0, 264, 265, 5, 62, 0, 0, 265, 41, 1, 0, 0, 0, 266, 267, 5, 125, 0, 0,
		267, 268, 5, 45, 0, 0, 268, 269, 5, 62, 0, 0, 269, 43, 1, 0, 0, 0, 270,
		271, 5, 45, 0, 0, 271, 272, 5, 45, 0, 0, 272, 273, 5, 62, 0, 0, 273, 45,
		1, 0, 0, 0, 274, 275, 5, 35, 0, 0, 275, 276, 5, 123, 0, 0, 276, 47, 1,
		0, 0, 0, 277, 278, 5, 35, 0, 0, 278, 279, 5, 62, 0, 0, 279, 49, 1, 0, 0,
		0, 280, 281, 5, 35, 0, 0, 281, 282, 5, 45, 0, 0, 282, 283, 5, 62, 0, 0,
		283, 51, 1, 0, 0, 0, 284, 285, 5, 62, 0, 0, 285, 53, 1, 0, 0, 0, 286, 287,
		5, 46, 0, 0, 287, 55, 1, 0, 0, 0, 288, 289, 5, 60, 0, 0, 289, 290, 5, 60,
		0, 0, 290, 291, 5, 60, 0, 0, 291, 292, 1, 0, 0, 0, 292, 293, 6, 26, 0,
		0, 293, 57, 1, 0, 0, 0, 294, 295, 5, 60, 0, 0, 295, 59, 1, 0, 0, 0, 296,
		297, 5, 61, 0, 0, 297, 61, 1, 0, 0, 0, 298, 299, 5, 43, 0, 0, 299, 63,
		1, 0, 0, 0, 300, 301, 5, 38, 0, 0, 301, 65, 1, 0, 0, 0, 302, 303, 5, 63,
		0, 0, 303, 67, 1, 0, 0, 0, 304, 305, 5, 40, 0, 0, 305, 69, 1, 0, 0, 0,
		306, 307, 5, 44, 0, 0, 307, 71, 1, 0, 0, 0, 308, 309, 5, 41, 0, 0, 309,
		73, 1, 0, 0, 0, 310, 311, 5, 91, 0, 0, 311, 75, 1, 0, 0, 0, 312, 313, 5,
		93, 0, 0, 313, 77, 1, 0, 0, 0, 314, 315, 5, 123, 0, 0, 315, 79, 1, 0, 0,
		0, 316, 317, 5, 125, 0, 0, 317, 81, 1, 0, 0, 0, 318, 319, 5, 35, 0, 0,
		319, 83, 1, 0, 0, 0, 320, 321, 5, 36, 0, 0, 321, 85, 1, 0, 0, 0, 322, 323,
		5, 58, 0, 0, 323, 87, 1, 0, 0, 0, 324, 325, 5, 37, 0, 0, 325, 89, 1, 0,
		0, 0, 326, 327, 5, 33, 0, 0, 327, 91, 1, 0, 0, 0, 328, 329, 5, 42, 0, 0,
		329, 93, 1, 0, 0, 0, 330, 331, 5, 45, 0, 0, 331, 95, 1, 0, 0, 0, 332, 333,
		5, 97, 0, 0, 333, 334, 5, 115, 0, 0, 334, 97, 1, 0, 0, 0, 335, 336, 5,
		96, 0, 0, 336, 99, 1, 0, 0, 0, 337, 338, 5, 39, 0, 0, 338, 101, 1, 0, 0,
		0, 339, 340, 5, 34, 0, 0, 340, 103, 1, 0, 0, 0, 341, 342, 5, 47, 0, 0,
		342, 343, 5, 47, 0, 0, 343, 347, 1, 0, 0, 0, 344, 346, 8, 0, 0, 0, 345,
		344, 1, 0, 0, 0, 346, 349, 1, 0, 0, 0, 347, 345, 1, 0, 0, 0, 347, 348,
		1, 0, 0, 0, 348, 105, 1, 0, 0, 0, 349, 347, 1, 0, 0, 0, 350, 351, 5, 10,
		0, 0, 351, 107, 1, 0, 0, 0, 352, 353, 7, 1, 0, 0, 353, 354, 1, 0, 0, 0,
		354, 355, 6, 52, 1, 0, 355, 109, 1, 0, 0, 0, 356, 358, 3, 202, 99, 0, 357,
		356, 1, 0, 0, 0, 358, 359, 1, 0, 0, 0, 359, 357, 1, 0, 0, 0, 359, 360,
		1, 0, 0, 0, 360, 111, 1, 0, 0, 0, 361, 362, 5, 48, 0, 0, 362, 363, 5, 111,
		0, 0, 363, 365, 1, 0, 0, 0, 364, 366, 3, 204, 100, 0, 365, 364, 1, 0, 0,
		0, 366, 367, 1, 0, 0, 0, 367, 365, 1, 0, 0, 0, 367, 368, 1, 0, 0, 0, 368,
		113, 1, 0, 0, 0, 369, 370, 5, 48, 0, 0, 370, 371, 5, 98, 0, 0, 371, 373,
		1, 0, 0, 0, 372, 374, 2, 48, 49, 0, 373, 372, 1, 0, 0, 0, 374, 375, 1,
		0, 0, 0, 375, 373, 1, 0, 0, 0, 375, 376, 1, 0, 0, 0, 376, 115, 1, 0, 0,
		0, 377, 378, 5, 48, 0, 0, 378, 379, 5, 120, 0, 0, 379, 381, 1, 0, 0, 0,
		380, 382, 3, 200, 98, 0, 381, 380, 1, 0, 0, 0, 382, 383, 1, 0, 0, 0, 383,
		381, 1, 0, 0, 0, 383, 384, 1, 0, 0, 0, 384, 117, 1, 0, 0, 0, 385, 386,
		5, 115, 0, 0, 386, 387, 5, 116, 0, 0, 387, 388, 5, 114, 0, 0, 388, 119,
		1, 0, 0, 0, 389, 390, 5, 108, 0, 0, 390, 391, 5, 105, 0, 0, 391, 392, 5,
		115, 0, 0, 392, 393, 5, 116, 0, 0, 393, 121, 1, 0, 0, 0, 394, 395, 5, 100,
		0, 0, 395, 396, 5, 105, 0, 0, 396, 397, 5, 99, 0, 0, 397, 398, 5, 116,
		0, 0, 398, 123, 1, 0, 0, 0, 399, 400, 5, 105, 0, 0, 400, 401, 5, 110, 0,
		0, 401, 408, 5, 116, 0, 0, 402, 403, 5, 102, 0, 0, 403, 404, 5, 108, 0,
		0, 404, 405, 5, 111, 0, 0, 405, 406, 5, 97, 0, 0, 406, 408, 5, 116, 0,
		0, 407, 399, 1, 0, 0, 0, 407, 402, 1, 0, 0, 0, 408, 125, 1, 0, 0, 0, 409,
		410, 5, 98, 0, 0, 410, 411, 5, 111, 0, 0, 411, 412, 5, 111, 0, 0, 412,
		413, 5, 108, 0, 0, 413, 127, 1, 0, 0, 0, 414, 415, 5, 116, 0, 0, 415, 416,
		5, 114, 0, 0, 416, 417, 5, 117, 0, 0, 417, 424, 5, 101, 0, 0, 418, 419,
		5, 102, 0, 0, 419, 420, 5, 97, 0, 0, 420, 421, 5, 108, 0, 0, 421, 422,
		5, 115, 0, 0, 422, 424, 5, 101, 0, 0, 423, 414, 1, 0, 0, 0, 423, 418, 1,
		0, 0, 0, 424, 129, 1, 0, 0, 0, 425, 426, 5, 97, 0, 0, 426, 427, 5, 108,
		0, 0, 427, 428, 5, 101, 0, 0, 428, 429, 5, 114, 0, 0, 429, 430, 5, 116,
		0, 0, 430, 131, 1, 0, 0, 0, 431, 432, 5, 99, 0, 0, 432, 433, 5, 104, 0,
		0, 433, 434, 5, 101, 0, 0, 434, 435, 5, 99, 0, 0, 435, 436, 5, 107, 0,
		0, 436, 133, 1, 0, 0, 0, 437, 438, 5, 116, 0, 0, 438, 439, 5, 104, 0, 0,
		439, 440, 5, 101, 0, 0, 440, 441, 5, 110, 0, 0, 441, 135, 1, 0, 0, 0, 442,
		443, 5, 100, 0, 0, 443, 444, 5, 101, 0, 0, 444, 445, 5, 115, 0, 0, 445,
		451, 5, 99, 0, 0, 446, 447, 5, 110, 0, 0, 447, 448, 5, 111, 0, 0, 448,
		449, 5, 116, 0, 0, 449, 451, 5, 101, 0, 0, 450, 442, 1, 0, 0, 0, 450, 446,
		1, 0, 0, 0, 451, 137, 1, 0, 0, 0, 452, 453, 5, 101, 0, 0, 453, 454, 5,
		108, 0, 0, 454, 455, 5, 115, 0, 0, 455, 456, 5, 101, 0, 0, 456, 139, 1,
		0, 0, 0, 457, 458, 5, 116, 0, 0, 458, 459, 5, 121, 0, 0, 459, 460, 5, 112,
		0, 0, 460, 461, 5, 101, 0, 0, 461, 141, 1, 0, 0, 0, 462, 463, 5, 105, 0,
		0, 463, 464, 5, 110, 0, 0, 464, 143, 1, 0, 0, 0, 465, 466, 5, 99, 0, 0,
		466, 467, 5, 97, 0, 0, 467, 468, 5, 108, 0, 0, 468, 469, 5, 108, 0, 0,
		469, 145, 1, 0, 0, 0, 470, 471, 5, 102, 0, 0, 471, 472, 5, 117, 0, 0, 472,
		473, 5, 110, 0, 0, 473, 474, 5, 99, 0, 0, 474, 475, 5, 116, 0, 0, 475,
		476, 5, 105, 0, 0, 476, 477, 5, 111, 0, 0, 477, 483, 5, 110, 0, 0, 478,
		479, 5, 102, 0, 0, 479, 480, 5, 117, 0, 0, 480, 481, 5, 110, 0, 0, 481,
		483, 5, 99, 0, 0, 482, 470, 1, 0, 0, 0, 482, 478, 1, 0, 0, 0, 483, 147,
		1, 0, 0, 0, 484, 485, 5, 99, 0, 0, 485, 486, 5, 111, 0, 0, 486, 487, 5,
		110, 0, 0, 487, 488, 5, 115, 0, 0, 488, 498, 5, 116, 0, 0, 489, 490, 5,
		99, 0, 0, 490, 491, 5, 111, 0, 0, 491, 492, 5, 110, 0, 0, 492, 493, 5,
		115, 0, 0, 493, 494, 5, 116, 0, 0, 494, 495, 5, 97, 0, 0, 495, 496, 5,
		110, 0, 0, 496, 498, 5, 116, 0, 0, 497, 484, 1, 0, 0, 0, 497, 489, 1, 0,
		0, 0, 498, 149, 1, 0, 0, 0, 499, 500, 5, 112, 0, 0, 500, 501, 5, 104, 0,
		0, 501, 502, 5, 105, 0, 0, 502, 151, 1, 0, 0, 0, 503, 504, 5, 112, 0, 0,
		504, 505, 5, 97, 0, 0, 505, 506, 5, 114, 0, 0, 506, 507, 5, 97, 0, 0, 507,
		521, 5, 109, 0, 0, 508, 509, 5, 102, 0, 0, 509, 510, 5, 111, 0, 0, 510,
		511, 5, 114, 0, 0, 511, 512, 5, 109, 0, 0, 512, 513, 5, 97, 0, 0, 513,
		514, 5, 108, 0, 0, 514, 515, 5, 95, 0, 0, 515, 516, 5, 112, 0, 0, 516,
		517, 5, 97, 0, 0, 517, 518, 5, 114, 0, 0, 518, 519, 5, 97, 0, 0, 519, 521,
		5, 109, 0, 0, 520, 503, 1, 0, 0, 0, 520, 508, 1, 0, 0, 0, 521, 153, 1,
		0, 0, 0, 522, 523, 5, 114, 0, 0, 523, 524, 5, 101, 0, 0, 524, 525, 5, 116,
		0, 0, 525, 526, 5, 117, 0, 0, 526, 527, 5, 114, 0, 0, 527, 532, 5, 110,
		0, 0, 528, 529, 5, 114, 0, 0, 529, 530, 5, 101, 0, 0, 530, 532, 5, 116,
		0, 0, 531, 522, 1, 0, 0, 0, 531, 528, 1, 0, 0, 0, 532, 155, 1, 0, 0, 0,
		533, 534, 5, 111, 0, 0, 534, 535, 5, 112, 0, 0, 535, 536, 5, 99, 0, 0,
		536, 537, 5, 111, 0, 0, 537, 538, 5, 100, 0, 0, 538, 539, 5, 101, 0, 0,
		539, 157, 1, 0, 0, 0, 540, 541, 5, 104, 0, 0, 541, 542, 5, 97, 0, 0, 542,
		543, 5, 118, 0, 0, 543, 544, 5, 101, 0, 0, 544, 159, 1, 0, 0, 0, 545, 546,
		5, 97, 0, 0, 546, 547, 5, 110, 0, 0, 547, 548, 5, 121, 0, 0, 548, 161,
		1, 0, 0, 0, 549, 550, 5, 110, 0, 0, 550, 551, 5, 111, 0, 0, 551, 552, 5,
		116, 0, 0, 552, 163, 1, 0, 0, 0, 553, 554, 5, 102, 0, 0, 554, 555, 5, 111,
		0, 0, 555, 556, 5, 114, 0, 0, 556, 165, 1, 0, 0, 0, 557, 558, 5, 114, 0,
		0, 558, 167, 1, 0, 0, 0, 559, 560, 5, 103, 0, 0, 560, 169, 1, 0, 0, 0,
		561, 562, 5, 101, 0, 0, 562, 171, 1, 0, 0, 0, 563, 567, 3, 198, 97, 0,
		564, 566, 3, 174, 85, 0, 565, 564, 1, 0, 0, 0, 566, 569, 1, 0, 0, 0, 567,
		565, 1, 0, 0, 0, 567, 568, 1, 0, 0, 0, 568, 173, 1, 0, 0, 0, 569, 567,
		1, 0, 0, 0, 570, 573, 7, 2, 0, 0, 571, 573, 3, 198, 97, 0, 572, 570, 1,
		0, 0, 0, 572, 571, 1, 0, 0, 0, 573, 175, 1, 0, 0, 0, 574, 583, 3, 100,
		48, 0, 575, 582, 8, 3, 0, 0, 576, 577, 5, 92, 0, 0, 577, 582, 5, 39, 0,
		0, 578, 579, 5, 92, 0, 0, 579, 582, 5, 92, 0, 0, 580, 582, 5, 92, 0, 0,
		581, 575, 1, 0, 0, 0, 581, 576, 1, 0, 0, 0, 581, 578, 1, 0, 0, 0, 581,
		580, 1, 0, 0, 0, 582, 585, 1, 0, 0, 0, 583, 581, 1, 0, 0, 0, 583, 584,
		1, 0, 0, 0, 584, 586, 1, 0, 0, 0, 585, 583, 1, 0, 0, 0, 586, 587, 3, 100,
		48, 0, 587, 603, 1, 0, 0, 0, 588, 597, 3, 102, 49, 0, 589, 596, 8, 4, 0,
		0, 590, 591, 5, 92, 0, 0, 591, 596, 5, 34, 0, 0, 592, 593, 5, 92, 0, 0,
		593, 596, 5, 92, 0, 0, 594, 596, 5, 92, 0, 0, 595, 589, 1, 0, 0, 0, 595,
		590, 1, 0, 0, 0, 595, 592, 1, 0, 0, 0, 595, 594, 1, 0, 0, 0, 596, 599,
		1, 0, 0, 0, 597, 595, 1, 0, 0, 0, 597, 598, 1, 0, 0, 0, 598, 600, 1, 0,
		0, 0, 599, 597, 1, 0, 0, 0, 600, 601, 3, 102, 49, 0, 601, 603, 1, 0, 0,
		0, 602, 574, 1, 0, 0, 0, 602, 588, 1, 0, 0, 0, 603, 177, 1, 0, 0, 0, 604,
		606, 5, 47, 0, 0, 605, 607, 3, 206, 101, 0, 606, 605, 1, 0, 0, 0, 607,
		608, 1, 0, 0, 0, 608, 606, 1, 0, 0, 0, 608, 609, 1, 0, 0, 0, 609, 610,
		1, 0, 0, 0, 610, 611, 5, 47, 0, 0, 611, 179, 1, 0, 0, 0, 612, 614, 7, 5,
		0, 0, 613, 612, 1, 0, 0, 0, 614, 615, 1, 0, 0, 0, 615, 613, 1, 0, 0, 0,
		615, 616, 1, 0, 0, 0, 616, 617, 1, 0, 0, 0, 617, 618, 6, 88, 1, 0, 618,
		181, 1, 0, 0, 0, 619, 620, 3, 196, 96, 0, 620, 621, 6, 89, 2, 0, 621, 629,
		1, 0, 0, 0, 622, 623, 5, 39, 0, 0, 623, 624, 3, 196, 96, 0, 624, 625, 6,
		89, 3, 0, 625, 626, 1, 0, 0, 0, 626, 627, 5, 39, 0, 0, 627, 629, 1, 0,
		0, 0, 628, 619, 1, 0, 0, 0, 628, 622, 1, 0, 0, 0, 629, 183, 1, 0, 0, 0,
		630, 631, 5, 13, 0, 0, 631, 632, 5, 10, 0, 0, 632, 633, 1, 0, 0, 0, 633,
		634, 6, 90, 4, 0, 634, 635, 1, 0, 0, 0, 635, 636, 6, 90, 5, 0, 636, 637,
		6, 90, 6, 0, 637, 185, 1, 0, 0, 0, 638, 639, 5, 10, 0, 0, 639, 640, 6,
		91, 7, 0, 640, 641, 1, 0, 0, 0, 641, 642, 6, 91, 5, 0, 642, 643, 6, 91,
		8, 0, 643, 187, 1, 0, 0, 0, 644, 645, 5, 13, 0, 0, 645, 646, 5, 10, 0,
		0, 646, 647, 1, 0, 0, 0, 647, 648, 3, 196, 96, 0, 648, 649, 6, 92, 9, 0,
		649, 189, 1, 0, 0, 0, 650, 651, 9, 0, 0, 0, 651, 191, 1, 0, 0, 0, 652,
		653, 5, 10, 0, 0, 653, 654, 3, 196, 96, 0, 654, 655, 6, 94, 10, 0, 655,
		193, 1, 0, 0, 0, 656, 657, 9, 0, 0, 0, 657, 195, 1, 0, 0, 0, 658, 662,
		7, 6, 0, 0, 659, 661, 7, 7, 0, 0, 660, 659, 1, 0, 0, 0, 661, 664, 1, 0,
		0, 0, 662, 660, 1, 0, 0, 0, 662, 663, 1, 0, 0, 0, 663, 197, 1, 0, 0, 0,
		664, 662, 1, 0, 0, 0, 665, 667, 7, 8, 0, 0, 666, 665, 1, 0, 0, 0, 667,
		199, 1, 0, 0, 0, 668, 669, 7, 9, 0, 0, 669, 201, 1, 0, 0, 0, 670, 671,
		7, 2, 0, 0, 671, 203, 1, 0, 0, 0, 672, 673, 7, 10, 0, 0, 673, 205, 1, 0,
		0, 0, 674, 675, 5, 92, 0, 0, 675, 678, 5, 47, 0, 0, 676, 678, 8, 11, 0,
		0, 677, 674, 1, 0, 0, 0, 677, 676, 1, 0, 0, 0, 678, 207, 1, 0, 0, 0, 29,
		0, 1, 2, 3, 347, 359, 367, 375, 383, 407, 423, 450, 482, 497, 520, 531,
		567, 572, 581, 583, 595, 597, 602, 608, 615, 628, 662, 666, 677, 11, 5,
		1, 0, 6, 0, 0, 1, 89, 0, 1, 89, 1, 1, 90, 2, 4, 0, 0, 5, 2, 0, 1, 91, 3,
		5, 3, 0, 1, 92, 4, 1, 94, 5,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// SyntaxFlowLexerInit initializes any static state used to implement SyntaxFlowLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewSyntaxFlowLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func SyntaxFlowLexerInit() {
	staticData := &syntaxflowlexerLexerStaticData
	staticData.once.Do(syntaxflowlexerLexerInit)
}

// NewSyntaxFlowLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewSyntaxFlowLexer(input antlr.CharStream) *SyntaxFlowLexer {
	SyntaxFlowLexerInit()
	l := new(SyntaxFlowLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &syntaxflowlexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "SyntaxFlowLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// SyntaxFlowLexer tokens.
const (
	SyntaxFlowLexerDeepFilter                  = 1
	SyntaxFlowLexerDeep                        = 2
	SyntaxFlowLexerPercent                     = 3
	SyntaxFlowLexerDeepDot                     = 4
	SyntaxFlowLexerLtEq                        = 5
	SyntaxFlowLexerGtEq                        = 6
	SyntaxFlowLexerDoubleGt                    = 7
	SyntaxFlowLexerFilter                      = 8
	SyntaxFlowLexerEqEq                        = 9
	SyntaxFlowLexerRegexpMatch                 = 10
	SyntaxFlowLexerNotRegexpMatch              = 11
	SyntaxFlowLexerAnd                         = 12
	SyntaxFlowLexerOr                          = 13
	SyntaxFlowLexerNotEq                       = 14
	SyntaxFlowLexerDollarBraceOpen             = 15
	SyntaxFlowLexerSemicolon                   = 16
	SyntaxFlowLexerConditionStart              = 17
	SyntaxFlowLexerDeepNextStart               = 18
	SyntaxFlowLexerUseStart                    = 19
	SyntaxFlowLexerDeepNextEnd                 = 20
	SyntaxFlowLexerDeepNext                    = 21
	SyntaxFlowLexerTopDefStart                 = 22
	SyntaxFlowLexerDefStart                    = 23
	SyntaxFlowLexerTopDef                      = 24
	SyntaxFlowLexerGt                          = 25
	SyntaxFlowLexerDot                         = 26
	SyntaxFlowLexerStartNowDoc                 = 27
	SyntaxFlowLexerLt                          = 28
	SyntaxFlowLexerEq                          = 29
	SyntaxFlowLexerAdd                         = 30
	SyntaxFlowLexerAmp                         = 31
	SyntaxFlowLexerQuestion                    = 32
	SyntaxFlowLexerOpenParen                   = 33
	SyntaxFlowLexerComma                       = 34
	SyntaxFlowLexerCloseParen                  = 35
	SyntaxFlowLexerListSelectOpen              = 36
	SyntaxFlowLexerListSelectClose             = 37
	SyntaxFlowLexerMapBuilderOpen              = 38
	SyntaxFlowLexerMapBuilderClose             = 39
	SyntaxFlowLexerListStart                   = 40
	SyntaxFlowLexerDollarOutput                = 41
	SyntaxFlowLexerColon                       = 42
	SyntaxFlowLexerSearch                      = 43
	SyntaxFlowLexerBang                        = 44
	SyntaxFlowLexerStar                        = 45
	SyntaxFlowLexerMinus                       = 46
	SyntaxFlowLexerAs                          = 47
	SyntaxFlowLexerBacktick                    = 48
	SyntaxFlowLexerSingleQuote                 = 49
	SyntaxFlowLexerDoubleQuote                 = 50
	SyntaxFlowLexerLineComment                 = 51
	SyntaxFlowLexerBreakLine                   = 52
	SyntaxFlowLexerWhiteSpace                  = 53
	SyntaxFlowLexerNumber                      = 54
	SyntaxFlowLexerOctalNumber                 = 55
	SyntaxFlowLexerBinaryNumber                = 56
	SyntaxFlowLexerHexNumber                   = 57
	SyntaxFlowLexerStringType                  = 58
	SyntaxFlowLexerListType                    = 59
	SyntaxFlowLexerDictType                    = 60
	SyntaxFlowLexerNumberType                  = 61
	SyntaxFlowLexerBoolType                    = 62
	SyntaxFlowLexerBoolLiteral                 = 63
	SyntaxFlowLexerAlert                       = 64
	SyntaxFlowLexerCheck                       = 65
	SyntaxFlowLexerThen                        = 66
	SyntaxFlowLexerDesc                        = 67
	SyntaxFlowLexerElse                        = 68
	SyntaxFlowLexerType                        = 69
	SyntaxFlowLexerIn                          = 70
	SyntaxFlowLexerCall                        = 71
	SyntaxFlowLexerFunction                    = 72
	SyntaxFlowLexerConstant                    = 73
	SyntaxFlowLexerPhi                         = 74
	SyntaxFlowLexerFormalParam                 = 75
	SyntaxFlowLexerReturn                      = 76
	SyntaxFlowLexerOpcode                      = 77
	SyntaxFlowLexerHave                        = 78
	SyntaxFlowLexerHaveAny                     = 79
	SyntaxFlowLexerNot                         = 80
	SyntaxFlowLexerFor                         = 81
	SyntaxFlowLexerConstSearchModePrefixRegexp = 82
	SyntaxFlowLexerConstSearchModePrefixGlob   = 83
	SyntaxFlowLexerConstSearchModePrefixExact  = 84
	SyntaxFlowLexerIdentifier                  = 85
	SyntaxFlowLexerIdentifierChar              = 86
	SyntaxFlowLexerQuotedStringLiteral         = 87
	SyntaxFlowLexerRegexpLiteral               = 88
	SyntaxFlowLexerWS                          = 89
	SyntaxFlowLexerHereDocIdentifierName       = 90
	SyntaxFlowLexerCRLFHereDocIdentifierBreak  = 91
	SyntaxFlowLexerLFHereDocIdentifierBreak    = 92
	SyntaxFlowLexerCRLFEndDoc                  = 93
	SyntaxFlowLexerCRLFHereDocText             = 94
	SyntaxFlowLexerLFEndDoc                    = 95
	SyntaxFlowLexerLFHereDocText               = 96
)

// SyntaxFlowLexer modes.
const (
	SyntaxFlowLexerHereDocIdentifier = iota + 1
	SyntaxFlowLexerCRLFHereDoc
	SyntaxFlowLexerLFHereDoc
)

func (l *SyntaxFlowLexer) Action(localctx antlr.RuleContext, ruleIndex, actionIndex int) {
	switch ruleIndex {
	case 89:
		l.HereDocIdentifierName_Action(localctx, actionIndex)

	case 90:
		l.CRLFHereDocIdentifierBreak_Action(localctx, actionIndex)

	case 91:
		l.LFHereDocIdentifierBreak_Action(localctx, actionIndex)

	case 92:
		l.CRLFEndDoc_Action(localctx, actionIndex)

	case 94:
		l.LFEndDoc_Action(localctx, actionIndex)

	default:
		panic("No registered action for: " + fmt.Sprint(ruleIndex))
	}
}

func (l *SyntaxFlowLexer) HereDocIdentifierName_Action(localctx antlr.RuleContext, actionIndex int) {
	this := l
	_ = this

	switch actionIndex {
	case 0:
		this.recordHereDocLabel()

	case 1:
		this.recordHereDocLabel()

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *SyntaxFlowLexer) CRLFHereDocIdentifierBreak_Action(localctx antlr.RuleContext, actionIndex int) {
	this := l
	_ = this

	switch actionIndex {
	case 2:
		this.recordHereDocLF()

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *SyntaxFlowLexer) LFHereDocIdentifierBreak_Action(localctx antlr.RuleContext, actionIndex int) {
	this := l
	_ = this

	switch actionIndex {
	case 3:
		this.recordHereDocLF()

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *SyntaxFlowLexer) CRLFEndDoc_Action(localctx antlr.RuleContext, actionIndex int) {
	this := l
	_ = this

	switch actionIndex {
	case 4:
		this.DocEndDistribute()

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *SyntaxFlowLexer) LFEndDoc_Action(localctx antlr.RuleContext, actionIndex int) {
	this := l
	_ = this

	switch actionIndex {
	case 5:
		this.DocEndDistribute()

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
