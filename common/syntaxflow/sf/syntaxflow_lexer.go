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
	*antlr.BaseLexer
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
		"DEFAULT_MODE",
	}
	staticData.literalNames = []string{
		"", "';'", "'->'", "'-->'", "'-<'", "'>-'", "'==>'", "'...'", "'%%'",
		"'..'", "'<='", "'>='", "'>>'", "'=>'", "'=='", "'=~'", "'!~'", "'&&'",
		"'||'", "'!='", "'?{'", "'-{'", "'}->'", "'#{'", "'#>'", "'#->'", "'>'",
		"'.'", "'<'", "'='", "'?'", "'('", "','", "')'", "'['", "']'", "'{'",
		"'}'", "'#'", "'$'", "':'", "'%'", "'!'", "'*'", "'-'", "'as'", "'`'",
		"'''", "'\"'", "", "", "", "", "", "'str'", "'list'", "'dict'", "",
		"'bool'", "", "'check'", "'then'", "", "'else'", "'type'", "'in'", "'call'",
		"", "'phi'", "", "", "'opcode'", "'have'", "'any'", "'not'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "", "", "DeepFilter", "Deep", "Percent", "DeepDot",
		"LtEq", "GtEq", "DoubleGt", "Filter", "EqEq", "RegexpMatch", "NotRegexpMatch",
		"And", "Or", "NotEq", "ConditionStart", "DeepNextStart", "DeepNextEnd",
		"TopDefStart", "DefStart", "TopDef", "Gt", "Dot", "Lt", "Eq", "Question",
		"OpenParen", "Comma", "CloseParen", "ListSelectOpen", "ListSelectClose",
		"MapBuilderOpen", "MapBuilderClose", "ListStart", "DollarOutput", "Colon",
		"Search", "Bang", "Star", "Minus", "As", "Backtick", "SingleQuote",
		"DoubleQuote", "WhiteSpace", "Number", "OctalNumber", "BinaryNumber",
		"HexNumber", "StringType", "ListType", "DictType", "NumberType", "BoolType",
		"BoolLiteral", "Check", "Then", "Desc", "Else", "Type", "In", "Call",
		"Constant", "Phi", "FormalParam", "Return", "Opcode", "Have", "HaveAny",
		"Not", "Identifier", "IdentifierChar", "QuotedStringLiteral", "RegexpLiteral",
		"WS",
	}
	staticData.ruleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "DeepFilter", "Deep", "Percent",
		"DeepDot", "LtEq", "GtEq", "DoubleGt", "Filter", "EqEq", "RegexpMatch",
		"NotRegexpMatch", "And", "Or", "NotEq", "ConditionStart", "DeepNextStart",
		"DeepNextEnd", "TopDefStart", "DefStart", "TopDef", "Gt", "Dot", "Lt",
		"Eq", "Question", "OpenParen", "Comma", "CloseParen", "ListSelectOpen",
		"ListSelectClose", "MapBuilderOpen", "MapBuilderClose", "ListStart",
		"DollarOutput", "Colon", "Search", "Bang", "Star", "Minus", "As", "Backtick",
		"SingleQuote", "DoubleQuote", "WhiteSpace", "Number", "OctalNumber",
		"BinaryNumber", "HexNumber", "StringType", "ListType", "DictType", "NumberType",
		"BoolType", "BoolLiteral", "Check", "Then", "Desc", "Else", "Type",
		"In", "Call", "Constant", "Phi", "FormalParam", "Return", "Opcode",
		"Have", "HaveAny", "Not", "Identifier", "IdentifierChar", "QuotedStringLiteral",
		"IdentifierCharStart", "HexDigit", "Digit", "OctalDigit", "RegexpLiteral",
		"RegexpLiteralChar", "WS",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 79, 546, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2,
		31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36,
		7, 36, 2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7,
		41, 2, 42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46,
		2, 47, 7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2,
		52, 7, 52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57,
		7, 57, 2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7,
		62, 2, 63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66, 7, 66, 2, 67, 7, 67,
		2, 68, 7, 68, 2, 69, 7, 69, 2, 70, 7, 70, 2, 71, 7, 71, 2, 72, 7, 72, 2,
		73, 7, 73, 2, 74, 7, 74, 2, 75, 7, 75, 2, 76, 7, 76, 2, 77, 7, 77, 2, 78,
		7, 78, 2, 79, 7, 79, 2, 80, 7, 80, 2, 81, 7, 81, 2, 82, 7, 82, 2, 83, 7,
		83, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1,
		3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1,
		7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10,
		1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 14, 1,
		14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17,
		1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1, 21, 1,
		21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 23, 1, 23, 1, 23, 1, 24, 1, 24,
		1, 24, 1, 24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 1, 27, 1, 28, 1, 28, 1,
		29, 1, 29, 1, 30, 1, 30, 1, 31, 1, 31, 1, 32, 1, 32, 1, 33, 1, 33, 1, 34,
		1, 34, 1, 35, 1, 35, 1, 36, 1, 36, 1, 37, 1, 37, 1, 38, 1, 38, 1, 39, 1,
		39, 1, 40, 1, 40, 1, 41, 1, 41, 1, 42, 1, 42, 1, 43, 1, 43, 1, 44, 1, 44,
		1, 44, 1, 45, 1, 45, 1, 46, 1, 46, 1, 47, 1, 47, 1, 48, 1, 48, 1, 48, 1,
		48, 1, 49, 4, 49, 301, 8, 49, 11, 49, 12, 49, 302, 1, 50, 1, 50, 1, 50,
		1, 50, 4, 50, 309, 8, 50, 11, 50, 12, 50, 310, 1, 51, 1, 51, 1, 51, 1,
		51, 4, 51, 317, 8, 51, 11, 51, 12, 51, 318, 1, 52, 1, 52, 1, 52, 1, 52,
		4, 52, 325, 8, 52, 11, 52, 12, 52, 326, 1, 53, 1, 53, 1, 53, 1, 53, 1,
		54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 55, 1, 55, 1, 55, 1, 55, 1, 55, 1, 56,
		1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 3, 56, 351, 8, 56, 1,
		57, 1, 57, 1, 57, 1, 57, 1, 57, 1, 58, 1, 58, 1, 58, 1, 58, 1, 58, 1, 58,
		1, 58, 1, 58, 1, 58, 3, 58, 367, 8, 58, 1, 59, 1, 59, 1, 59, 1, 59, 1,
		59, 1, 59, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 61, 1, 61, 1, 61, 1, 61,
		1, 61, 1, 61, 1, 61, 1, 61, 3, 61, 388, 8, 61, 1, 62, 1, 62, 1, 62, 1,
		62, 1, 62, 1, 63, 1, 63, 1, 63, 1, 63, 1, 63, 1, 64, 1, 64, 1, 64, 1, 65,
		1, 65, 1, 65, 1, 65, 1, 65, 1, 66, 1, 66, 1, 66, 1, 66, 1, 66, 1, 66, 1,
		66, 1, 66, 1, 66, 1, 66, 1, 66, 1, 66, 1, 66, 3, 66, 421, 8, 66, 1, 67,
		1, 67, 1, 67, 1, 67, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 1,
		68, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 1, 68, 3, 68,
		444, 8, 68, 1, 69, 1, 69, 1, 69, 1, 69, 1, 69, 1, 69, 1, 69, 1, 69, 1,
		69, 3, 69, 455, 8, 69, 1, 70, 1, 70, 1, 70, 1, 70, 1, 70, 1, 70, 1, 70,
		1, 71, 1, 71, 1, 71, 1, 71, 1, 71, 1, 72, 1, 72, 1, 72, 1, 72, 1, 73, 1,
		73, 1, 73, 1, 73, 1, 74, 1, 74, 5, 74, 479, 8, 74, 10, 74, 12, 74, 482,
		9, 74, 1, 75, 1, 75, 3, 75, 486, 8, 75, 1, 76, 1, 76, 1, 76, 1, 76, 1,
		76, 1, 76, 1, 76, 5, 76, 495, 8, 76, 10, 76, 12, 76, 498, 9, 76, 1, 76,
		1, 76, 1, 76, 1, 76, 1, 76, 1, 76, 1, 76, 1, 76, 1, 76, 5, 76, 509, 8,
		76, 10, 76, 12, 76, 512, 9, 76, 1, 76, 1, 76, 3, 76, 516, 8, 76, 1, 77,
		3, 77, 519, 8, 77, 1, 78, 1, 78, 1, 79, 1, 79, 1, 80, 1, 80, 1, 81, 1,
		81, 4, 81, 529, 8, 81, 11, 81, 12, 81, 530, 1, 81, 1, 81, 1, 82, 1, 82,
		1, 82, 3, 82, 538, 8, 82, 1, 83, 4, 83, 541, 8, 83, 11, 83, 12, 83, 542,
		1, 83, 1, 83, 0, 0, 84, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15,
		8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17,
		35, 18, 37, 19, 39, 20, 41, 21, 43, 22, 45, 23, 47, 24, 49, 25, 51, 26,
		53, 27, 55, 28, 57, 29, 59, 30, 61, 31, 63, 32, 65, 33, 67, 34, 69, 35,
		71, 36, 73, 37, 75, 38, 77, 39, 79, 40, 81, 41, 83, 42, 85, 43, 87, 44,
		89, 45, 91, 46, 93, 47, 95, 48, 97, 49, 99, 50, 101, 51, 103, 52, 105,
		53, 107, 54, 109, 55, 111, 56, 113, 57, 115, 58, 117, 59, 119, 60, 121,
		61, 123, 62, 125, 63, 127, 64, 129, 65, 131, 66, 133, 67, 135, 68, 137,
		69, 139, 70, 141, 71, 143, 72, 145, 73, 147, 74, 149, 75, 151, 76, 153,
		77, 155, 0, 157, 0, 159, 0, 161, 0, 163, 78, 165, 0, 167, 79, 1, 0, 9,
		3, 0, 10, 10, 13, 13, 32, 32, 1, 0, 48, 57, 4, 0, 10, 10, 13, 13, 39, 39,
		92, 92, 4, 0, 10, 10, 13, 13, 34, 34, 92, 92, 4, 0, 42, 42, 65, 90, 95,
		95, 97, 122, 3, 0, 48, 57, 65, 70, 97, 102, 1, 0, 48, 55, 1, 0, 47, 47,
		3, 0, 9, 9, 13, 13, 32, 32, 564, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0,
		5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0,
		13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0,
		0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0,
		0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0,
		0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 43, 1,
		0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0, 0, 0, 0, 51,
		1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1, 0, 0, 0, 0,
		59, 1, 0, 0, 0, 0, 61, 1, 0, 0, 0, 0, 63, 1, 0, 0, 0, 0, 65, 1, 0, 0, 0,
		0, 67, 1, 0, 0, 0, 0, 69, 1, 0, 0, 0, 0, 71, 1, 0, 0, 0, 0, 73, 1, 0, 0,
		0, 0, 75, 1, 0, 0, 0, 0, 77, 1, 0, 0, 0, 0, 79, 1, 0, 0, 0, 0, 81, 1, 0,
		0, 0, 0, 83, 1, 0, 0, 0, 0, 85, 1, 0, 0, 0, 0, 87, 1, 0, 0, 0, 0, 89, 1,
		0, 0, 0, 0, 91, 1, 0, 0, 0, 0, 93, 1, 0, 0, 0, 0, 95, 1, 0, 0, 0, 0, 97,
		1, 0, 0, 0, 0, 99, 1, 0, 0, 0, 0, 101, 1, 0, 0, 0, 0, 103, 1, 0, 0, 0,
		0, 105, 1, 0, 0, 0, 0, 107, 1, 0, 0, 0, 0, 109, 1, 0, 0, 0, 0, 111, 1,
		0, 0, 0, 0, 113, 1, 0, 0, 0, 0, 115, 1, 0, 0, 0, 0, 117, 1, 0, 0, 0, 0,
		119, 1, 0, 0, 0, 0, 121, 1, 0, 0, 0, 0, 123, 1, 0, 0, 0, 0, 125, 1, 0,
		0, 0, 0, 127, 1, 0, 0, 0, 0, 129, 1, 0, 0, 0, 0, 131, 1, 0, 0, 0, 0, 133,
		1, 0, 0, 0, 0, 135, 1, 0, 0, 0, 0, 137, 1, 0, 0, 0, 0, 139, 1, 0, 0, 0,
		0, 141, 1, 0, 0, 0, 0, 143, 1, 0, 0, 0, 0, 145, 1, 0, 0, 0, 0, 147, 1,
		0, 0, 0, 0, 149, 1, 0, 0, 0, 0, 151, 1, 0, 0, 0, 0, 153, 1, 0, 0, 0, 0,
		163, 1, 0, 0, 0, 0, 167, 1, 0, 0, 0, 1, 169, 1, 0, 0, 0, 3, 171, 1, 0,
		0, 0, 5, 174, 1, 0, 0, 0, 7, 178, 1, 0, 0, 0, 9, 181, 1, 0, 0, 0, 11, 184,
		1, 0, 0, 0, 13, 188, 1, 0, 0, 0, 15, 192, 1, 0, 0, 0, 17, 195, 1, 0, 0,
		0, 19, 198, 1, 0, 0, 0, 21, 201, 1, 0, 0, 0, 23, 204, 1, 0, 0, 0, 25, 207,
		1, 0, 0, 0, 27, 210, 1, 0, 0, 0, 29, 213, 1, 0, 0, 0, 31, 216, 1, 0, 0,
		0, 33, 219, 1, 0, 0, 0, 35, 222, 1, 0, 0, 0, 37, 225, 1, 0, 0, 0, 39, 228,
		1, 0, 0, 0, 41, 231, 1, 0, 0, 0, 43, 234, 1, 0, 0, 0, 45, 238, 1, 0, 0,
		0, 47, 241, 1, 0, 0, 0, 49, 244, 1, 0, 0, 0, 51, 248, 1, 0, 0, 0, 53, 250,
		1, 0, 0, 0, 55, 252, 1, 0, 0, 0, 57, 254, 1, 0, 0, 0, 59, 256, 1, 0, 0,
		0, 61, 258, 1, 0, 0, 0, 63, 260, 1, 0, 0, 0, 65, 262, 1, 0, 0, 0, 67, 264,
		1, 0, 0, 0, 69, 266, 1, 0, 0, 0, 71, 268, 1, 0, 0, 0, 73, 270, 1, 0, 0,
		0, 75, 272, 1, 0, 0, 0, 77, 274, 1, 0, 0, 0, 79, 276, 1, 0, 0, 0, 81, 278,
		1, 0, 0, 0, 83, 280, 1, 0, 0, 0, 85, 282, 1, 0, 0, 0, 87, 284, 1, 0, 0,
		0, 89, 286, 1, 0, 0, 0, 91, 289, 1, 0, 0, 0, 93, 291, 1, 0, 0, 0, 95, 293,
		1, 0, 0, 0, 97, 295, 1, 0, 0, 0, 99, 300, 1, 0, 0, 0, 101, 304, 1, 0, 0,
		0, 103, 312, 1, 0, 0, 0, 105, 320, 1, 0, 0, 0, 107, 328, 1, 0, 0, 0, 109,
		332, 1, 0, 0, 0, 111, 337, 1, 0, 0, 0, 113, 350, 1, 0, 0, 0, 115, 352,
		1, 0, 0, 0, 117, 366, 1, 0, 0, 0, 119, 368, 1, 0, 0, 0, 121, 374, 1, 0,
		0, 0, 123, 387, 1, 0, 0, 0, 125, 389, 1, 0, 0, 0, 127, 394, 1, 0, 0, 0,
		129, 399, 1, 0, 0, 0, 131, 402, 1, 0, 0, 0, 133, 420, 1, 0, 0, 0, 135,
		422, 1, 0, 0, 0, 137, 443, 1, 0, 0, 0, 139, 454, 1, 0, 0, 0, 141, 456,
		1, 0, 0, 0, 143, 463, 1, 0, 0, 0, 145, 468, 1, 0, 0, 0, 147, 472, 1, 0,
		0, 0, 149, 476, 1, 0, 0, 0, 151, 485, 1, 0, 0, 0, 153, 515, 1, 0, 0, 0,
		155, 518, 1, 0, 0, 0, 157, 520, 1, 0, 0, 0, 159, 522, 1, 0, 0, 0, 161,
		524, 1, 0, 0, 0, 163, 526, 1, 0, 0, 0, 165, 537, 1, 0, 0, 0, 167, 540,
		1, 0, 0, 0, 169, 170, 5, 59, 0, 0, 170, 2, 1, 0, 0, 0, 171, 172, 5, 45,
		0, 0, 172, 173, 5, 62, 0, 0, 173, 4, 1, 0, 0, 0, 174, 175, 5, 45, 0, 0,
		175, 176, 5, 45, 0, 0, 176, 177, 5, 62, 0, 0, 177, 6, 1, 0, 0, 0, 178,
		179, 5, 45, 0, 0, 179, 180, 5, 60, 0, 0, 180, 8, 1, 0, 0, 0, 181, 182,
		5, 62, 0, 0, 182, 183, 5, 45, 0, 0, 183, 10, 1, 0, 0, 0, 184, 185, 5, 61,
		0, 0, 185, 186, 5, 61, 0, 0, 186, 187, 5, 62, 0, 0, 187, 12, 1, 0, 0, 0,
		188, 189, 5, 46, 0, 0, 189, 190, 5, 46, 0, 0, 190, 191, 5, 46, 0, 0, 191,
		14, 1, 0, 0, 0, 192, 193, 5, 37, 0, 0, 193, 194, 5, 37, 0, 0, 194, 16,
		1, 0, 0, 0, 195, 196, 5, 46, 0, 0, 196, 197, 5, 46, 0, 0, 197, 18, 1, 0,
		0, 0, 198, 199, 5, 60, 0, 0, 199, 200, 5, 61, 0, 0, 200, 20, 1, 0, 0, 0,
		201, 202, 5, 62, 0, 0, 202, 203, 5, 61, 0, 0, 203, 22, 1, 0, 0, 0, 204,
		205, 5, 62, 0, 0, 205, 206, 5, 62, 0, 0, 206, 24, 1, 0, 0, 0, 207, 208,
		5, 61, 0, 0, 208, 209, 5, 62, 0, 0, 209, 26, 1, 0, 0, 0, 210, 211, 5, 61,
		0, 0, 211, 212, 5, 61, 0, 0, 212, 28, 1, 0, 0, 0, 213, 214, 5, 61, 0, 0,
		214, 215, 5, 126, 0, 0, 215, 30, 1, 0, 0, 0, 216, 217, 5, 33, 0, 0, 217,
		218, 5, 126, 0, 0, 218, 32, 1, 0, 0, 0, 219, 220, 5, 38, 0, 0, 220, 221,
		5, 38, 0, 0, 221, 34, 1, 0, 0, 0, 222, 223, 5, 124, 0, 0, 223, 224, 5,
		124, 0, 0, 224, 36, 1, 0, 0, 0, 225, 226, 5, 33, 0, 0, 226, 227, 5, 61,
		0, 0, 227, 38, 1, 0, 0, 0, 228, 229, 5, 63, 0, 0, 229, 230, 5, 123, 0,
		0, 230, 40, 1, 0, 0, 0, 231, 232, 5, 45, 0, 0, 232, 233, 5, 123, 0, 0,
		233, 42, 1, 0, 0, 0, 234, 235, 5, 125, 0, 0, 235, 236, 5, 45, 0, 0, 236,
		237, 5, 62, 0, 0, 237, 44, 1, 0, 0, 0, 238, 239, 5, 35, 0, 0, 239, 240,
		5, 123, 0, 0, 240, 46, 1, 0, 0, 0, 241, 242, 5, 35, 0, 0, 242, 243, 5,
		62, 0, 0, 243, 48, 1, 0, 0, 0, 244, 245, 5, 35, 0, 0, 245, 246, 5, 45,
		0, 0, 246, 247, 5, 62, 0, 0, 247, 50, 1, 0, 0, 0, 248, 249, 5, 62, 0, 0,
		249, 52, 1, 0, 0, 0, 250, 251, 5, 46, 0, 0, 251, 54, 1, 0, 0, 0, 252, 253,
		5, 60, 0, 0, 253, 56, 1, 0, 0, 0, 254, 255, 5, 61, 0, 0, 255, 58, 1, 0,
		0, 0, 256, 257, 5, 63, 0, 0, 257, 60, 1, 0, 0, 0, 258, 259, 5, 40, 0, 0,
		259, 62, 1, 0, 0, 0, 260, 261, 5, 44, 0, 0, 261, 64, 1, 0, 0, 0, 262, 263,
		5, 41, 0, 0, 263, 66, 1, 0, 0, 0, 264, 265, 5, 91, 0, 0, 265, 68, 1, 0,
		0, 0, 266, 267, 5, 93, 0, 0, 267, 70, 1, 0, 0, 0, 268, 269, 5, 123, 0,
		0, 269, 72, 1, 0, 0, 0, 270, 271, 5, 125, 0, 0, 271, 74, 1, 0, 0, 0, 272,
		273, 5, 35, 0, 0, 273, 76, 1, 0, 0, 0, 274, 275, 5, 36, 0, 0, 275, 78,
		1, 0, 0, 0, 276, 277, 5, 58, 0, 0, 277, 80, 1, 0, 0, 0, 278, 279, 5, 37,
		0, 0, 279, 82, 1, 0, 0, 0, 280, 281, 5, 33, 0, 0, 281, 84, 1, 0, 0, 0,
		282, 283, 5, 42, 0, 0, 283, 86, 1, 0, 0, 0, 284, 285, 5, 45, 0, 0, 285,
		88, 1, 0, 0, 0, 286, 287, 5, 97, 0, 0, 287, 288, 5, 115, 0, 0, 288, 90,
		1, 0, 0, 0, 289, 290, 5, 96, 0, 0, 290, 92, 1, 0, 0, 0, 291, 292, 5, 39,
		0, 0, 292, 94, 1, 0, 0, 0, 293, 294, 5, 34, 0, 0, 294, 96, 1, 0, 0, 0,
		295, 296, 7, 0, 0, 0, 296, 297, 1, 0, 0, 0, 297, 298, 6, 48, 0, 0, 298,
		98, 1, 0, 0, 0, 299, 301, 3, 159, 79, 0, 300, 299, 1, 0, 0, 0, 301, 302,
		1, 0, 0, 0, 302, 300, 1, 0, 0, 0, 302, 303, 1, 0, 0, 0, 303, 100, 1, 0,
		0, 0, 304, 305, 5, 48, 0, 0, 305, 306, 5, 111, 0, 0, 306, 308, 1, 0, 0,
		0, 307, 309, 3, 161, 80, 0, 308, 307, 1, 0, 0, 0, 309, 310, 1, 0, 0, 0,
		310, 308, 1, 0, 0, 0, 310, 311, 1, 0, 0, 0, 311, 102, 1, 0, 0, 0, 312,
		313, 5, 48, 0, 0, 313, 314, 5, 98, 0, 0, 314, 316, 1, 0, 0, 0, 315, 317,
		2, 48, 49, 0, 316, 315, 1, 0, 0, 0, 317, 318, 1, 0, 0, 0, 318, 316, 1,
		0, 0, 0, 318, 319, 1, 0, 0, 0, 319, 104, 1, 0, 0, 0, 320, 321, 5, 48, 0,
		0, 321, 322, 5, 120, 0, 0, 322, 324, 1, 0, 0, 0, 323, 325, 3, 157, 78,
		0, 324, 323, 1, 0, 0, 0, 325, 326, 1, 0, 0, 0, 326, 324, 1, 0, 0, 0, 326,
		327, 1, 0, 0, 0, 327, 106, 1, 0, 0, 0, 328, 329, 5, 115, 0, 0, 329, 330,
		5, 116, 0, 0, 330, 331, 5, 114, 0, 0, 331, 108, 1, 0, 0, 0, 332, 333, 5,
		108, 0, 0, 333, 334, 5, 105, 0, 0, 334, 335, 5, 115, 0, 0, 335, 336, 5,
		116, 0, 0, 336, 110, 1, 0, 0, 0, 337, 338, 5, 100, 0, 0, 338, 339, 5, 105,
		0, 0, 339, 340, 5, 99, 0, 0, 340, 341, 5, 116, 0, 0, 341, 112, 1, 0, 0,
		0, 342, 343, 5, 105, 0, 0, 343, 344, 5, 110, 0, 0, 344, 351, 5, 116, 0,
		0, 345, 346, 5, 102, 0, 0, 346, 347, 5, 108, 0, 0, 347, 348, 5, 111, 0,
		0, 348, 349, 5, 97, 0, 0, 349, 351, 5, 116, 0, 0, 350, 342, 1, 0, 0, 0,
		350, 345, 1, 0, 0, 0, 351, 114, 1, 0, 0, 0, 352, 353, 5, 98, 0, 0, 353,
		354, 5, 111, 0, 0, 354, 355, 5, 111, 0, 0, 355, 356, 5, 108, 0, 0, 356,
		116, 1, 0, 0, 0, 357, 358, 5, 116, 0, 0, 358, 359, 5, 114, 0, 0, 359, 360,
		5, 117, 0, 0, 360, 367, 5, 101, 0, 0, 361, 362, 5, 102, 0, 0, 362, 363,
		5, 97, 0, 0, 363, 364, 5, 108, 0, 0, 364, 365, 5, 115, 0, 0, 365, 367,
		5, 101, 0, 0, 366, 357, 1, 0, 0, 0, 366, 361, 1, 0, 0, 0, 367, 118, 1,
		0, 0, 0, 368, 369, 5, 99, 0, 0, 369, 370, 5, 104, 0, 0, 370, 371, 5, 101,
		0, 0, 371, 372, 5, 99, 0, 0, 372, 373, 5, 107, 0, 0, 373, 120, 1, 0, 0,
		0, 374, 375, 5, 116, 0, 0, 375, 376, 5, 104, 0, 0, 376, 377, 5, 101, 0,
		0, 377, 378, 5, 110, 0, 0, 378, 122, 1, 0, 0, 0, 379, 380, 5, 100, 0, 0,
		380, 381, 5, 101, 0, 0, 381, 382, 5, 115, 0, 0, 382, 388, 5, 99, 0, 0,
		383, 384, 5, 110, 0, 0, 384, 385, 5, 111, 0, 0, 385, 386, 5, 116, 0, 0,
		386, 388, 5, 101, 0, 0, 387, 379, 1, 0, 0, 0, 387, 383, 1, 0, 0, 0, 388,
		124, 1, 0, 0, 0, 389, 390, 5, 101, 0, 0, 390, 391, 5, 108, 0, 0, 391, 392,
		5, 115, 0, 0, 392, 393, 5, 101, 0, 0, 393, 126, 1, 0, 0, 0, 394, 395, 5,
		116, 0, 0, 395, 396, 5, 121, 0, 0, 396, 397, 5, 112, 0, 0, 397, 398, 5,
		101, 0, 0, 398, 128, 1, 0, 0, 0, 399, 400, 5, 105, 0, 0, 400, 401, 5, 110,
		0, 0, 401, 130, 1, 0, 0, 0, 402, 403, 5, 99, 0, 0, 403, 404, 5, 97, 0,
		0, 404, 405, 5, 108, 0, 0, 405, 406, 5, 108, 0, 0, 406, 132, 1, 0, 0, 0,
		407, 408, 5, 99, 0, 0, 408, 409, 5, 111, 0, 0, 409, 410, 5, 110, 0, 0,
		410, 411, 5, 115, 0, 0, 411, 421, 5, 116, 0, 0, 412, 413, 5, 99, 0, 0,
		413, 414, 5, 111, 0, 0, 414, 415, 5, 110, 0, 0, 415, 416, 5, 115, 0, 0,
		416, 417, 5, 116, 0, 0, 417, 418, 5, 97, 0, 0, 418, 419, 5, 110, 0, 0,
		419, 421, 5, 116, 0, 0, 420, 407, 1, 0, 0, 0, 420, 412, 1, 0, 0, 0, 421,
		134, 1, 0, 0, 0, 422, 423, 5, 112, 0, 0, 423, 424, 5, 104, 0, 0, 424, 425,
		5, 105, 0, 0, 425, 136, 1, 0, 0, 0, 426, 427, 5, 112, 0, 0, 427, 428, 5,
		97, 0, 0, 428, 429, 5, 114, 0, 0, 429, 430, 5, 97, 0, 0, 430, 444, 5, 109,
		0, 0, 431, 432, 5, 102, 0, 0, 432, 433, 5, 111, 0, 0, 433, 434, 5, 114,
		0, 0, 434, 435, 5, 109, 0, 0, 435, 436, 5, 97, 0, 0, 436, 437, 5, 108,
		0, 0, 437, 438, 5, 95, 0, 0, 438, 439, 5, 112, 0, 0, 439, 440, 5, 97, 0,
		0, 440, 441, 5, 114, 0, 0, 441, 442, 5, 97, 0, 0, 442, 444, 5, 109, 0,
		0, 443, 426, 1, 0, 0, 0, 443, 431, 1, 0, 0, 0, 444, 138, 1, 0, 0, 0, 445,
		446, 5, 114, 0, 0, 446, 447, 5, 101, 0, 0, 447, 448, 5, 116, 0, 0, 448,
		449, 5, 117, 0, 0, 449, 450, 5, 114, 0, 0, 450, 455, 5, 110, 0, 0, 451,
		452, 5, 114, 0, 0, 452, 453, 5, 101, 0, 0, 453, 455, 5, 116, 0, 0, 454,
		445, 1, 0, 0, 0, 454, 451, 1, 0, 0, 0, 455, 140, 1, 0, 0, 0, 456, 457,
		5, 111, 0, 0, 457, 458, 5, 112, 0, 0, 458, 459, 5, 99, 0, 0, 459, 460,
		5, 111, 0, 0, 460, 461, 5, 100, 0, 0, 461, 462, 5, 101, 0, 0, 462, 142,
		1, 0, 0, 0, 463, 464, 5, 104, 0, 0, 464, 465, 5, 97, 0, 0, 465, 466, 5,
		118, 0, 0, 466, 467, 5, 101, 0, 0, 467, 144, 1, 0, 0, 0, 468, 469, 5, 97,
		0, 0, 469, 470, 5, 110, 0, 0, 470, 471, 5, 121, 0, 0, 471, 146, 1, 0, 0,
		0, 472, 473, 5, 110, 0, 0, 473, 474, 5, 111, 0, 0, 474, 475, 5, 116, 0,
		0, 475, 148, 1, 0, 0, 0, 476, 480, 3, 155, 77, 0, 477, 479, 3, 151, 75,
		0, 478, 477, 1, 0, 0, 0, 479, 482, 1, 0, 0, 0, 480, 478, 1, 0, 0, 0, 480,
		481, 1, 0, 0, 0, 481, 150, 1, 0, 0, 0, 482, 480, 1, 0, 0, 0, 483, 486,
		7, 1, 0, 0, 484, 486, 3, 155, 77, 0, 485, 483, 1, 0, 0, 0, 485, 484, 1,
		0, 0, 0, 486, 152, 1, 0, 0, 0, 487, 496, 3, 93, 46, 0, 488, 495, 8, 2,
		0, 0, 489, 490, 5, 92, 0, 0, 490, 495, 5, 39, 0, 0, 491, 492, 5, 92, 0,
		0, 492, 495, 5, 92, 0, 0, 493, 495, 5, 92, 0, 0, 494, 488, 1, 0, 0, 0,
		494, 489, 1, 0, 0, 0, 494, 491, 1, 0, 0, 0, 494, 493, 1, 0, 0, 0, 495,
		498, 1, 0, 0, 0, 496, 494, 1, 0, 0, 0, 496, 497, 1, 0, 0, 0, 497, 499,
		1, 0, 0, 0, 498, 496, 1, 0, 0, 0, 499, 500, 3, 93, 46, 0, 500, 516, 1,
		0, 0, 0, 501, 510, 3, 95, 47, 0, 502, 509, 8, 3, 0, 0, 503, 504, 5, 92,
		0, 0, 504, 509, 5, 34, 0, 0, 505, 506, 5, 92, 0, 0, 506, 509, 5, 92, 0,
		0, 507, 509, 5, 92, 0, 0, 508, 502, 1, 0, 0, 0, 508, 503, 1, 0, 0, 0, 508,
		505, 1, 0, 0, 0, 508, 507, 1, 0, 0, 0, 509, 512, 1, 0, 0, 0, 510, 508,
		1, 0, 0, 0, 510, 511, 1, 0, 0, 0, 511, 513, 1, 0, 0, 0, 512, 510, 1, 0,
		0, 0, 513, 514, 3, 95, 47, 0, 514, 516, 1, 0, 0, 0, 515, 487, 1, 0, 0,
		0, 515, 501, 1, 0, 0, 0, 516, 154, 1, 0, 0, 0, 517, 519, 7, 4, 0, 0, 518,
		517, 1, 0, 0, 0, 519, 156, 1, 0, 0, 0, 520, 521, 7, 5, 0, 0, 521, 158,
		1, 0, 0, 0, 522, 523, 7, 1, 0, 0, 523, 160, 1, 0, 0, 0, 524, 525, 7, 6,
		0, 0, 525, 162, 1, 0, 0, 0, 526, 528, 5, 47, 0, 0, 527, 529, 3, 165, 82,
		0, 528, 527, 1, 0, 0, 0, 529, 530, 1, 0, 0, 0, 530, 528, 1, 0, 0, 0, 530,
		531, 1, 0, 0, 0, 531, 532, 1, 0, 0, 0, 532, 533, 5, 47, 0, 0, 533, 164,
		1, 0, 0, 0, 534, 535, 5, 92, 0, 0, 535, 538, 5, 47, 0, 0, 536, 538, 8,
		7, 0, 0, 537, 534, 1, 0, 0, 0, 537, 536, 1, 0, 0, 0, 538, 166, 1, 0, 0,
		0, 539, 541, 7, 8, 0, 0, 540, 539, 1, 0, 0, 0, 541, 542, 1, 0, 0, 0, 542,
		540, 1, 0, 0, 0, 542, 543, 1, 0, 0, 0, 543, 544, 1, 0, 0, 0, 544, 545,
		6, 83, 0, 0, 545, 168, 1, 0, 0, 0, 22, 0, 302, 310, 318, 326, 350, 366,
		387, 420, 443, 454, 480, 485, 494, 496, 508, 510, 515, 518, 530, 537, 542,
		1, 6, 0, 0,
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
	l.GrammarFileName = "SyntaxFlow.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// SyntaxFlowLexer tokens.
const (
	SyntaxFlowLexerT__0                = 1
	SyntaxFlowLexerT__1                = 2
	SyntaxFlowLexerT__2                = 3
	SyntaxFlowLexerT__3                = 4
	SyntaxFlowLexerT__4                = 5
	SyntaxFlowLexerDeepFilter          = 6
	SyntaxFlowLexerDeep                = 7
	SyntaxFlowLexerPercent             = 8
	SyntaxFlowLexerDeepDot             = 9
	SyntaxFlowLexerLtEq                = 10
	SyntaxFlowLexerGtEq                = 11
	SyntaxFlowLexerDoubleGt            = 12
	SyntaxFlowLexerFilter              = 13
	SyntaxFlowLexerEqEq                = 14
	SyntaxFlowLexerRegexpMatch         = 15
	SyntaxFlowLexerNotRegexpMatch      = 16
	SyntaxFlowLexerAnd                 = 17
	SyntaxFlowLexerOr                  = 18
	SyntaxFlowLexerNotEq               = 19
	SyntaxFlowLexerConditionStart      = 20
	SyntaxFlowLexerDeepNextStart       = 21
	SyntaxFlowLexerDeepNextEnd         = 22
	SyntaxFlowLexerTopDefStart         = 23
	SyntaxFlowLexerDefStart            = 24
	SyntaxFlowLexerTopDef              = 25
	SyntaxFlowLexerGt                  = 26
	SyntaxFlowLexerDot                 = 27
	SyntaxFlowLexerLt                  = 28
	SyntaxFlowLexerEq                  = 29
	SyntaxFlowLexerQuestion            = 30
	SyntaxFlowLexerOpenParen           = 31
	SyntaxFlowLexerComma               = 32
	SyntaxFlowLexerCloseParen          = 33
	SyntaxFlowLexerListSelectOpen      = 34
	SyntaxFlowLexerListSelectClose     = 35
	SyntaxFlowLexerMapBuilderOpen      = 36
	SyntaxFlowLexerMapBuilderClose     = 37
	SyntaxFlowLexerListStart           = 38
	SyntaxFlowLexerDollarOutput        = 39
	SyntaxFlowLexerColon               = 40
	SyntaxFlowLexerSearch              = 41
	SyntaxFlowLexerBang                = 42
	SyntaxFlowLexerStar                = 43
	SyntaxFlowLexerMinus               = 44
	SyntaxFlowLexerAs                  = 45
	SyntaxFlowLexerBacktick            = 46
	SyntaxFlowLexerSingleQuote         = 47
	SyntaxFlowLexerDoubleQuote         = 48
	SyntaxFlowLexerWhiteSpace          = 49
	SyntaxFlowLexerNumber              = 50
	SyntaxFlowLexerOctalNumber         = 51
	SyntaxFlowLexerBinaryNumber        = 52
	SyntaxFlowLexerHexNumber           = 53
	SyntaxFlowLexerStringType          = 54
	SyntaxFlowLexerListType            = 55
	SyntaxFlowLexerDictType            = 56
	SyntaxFlowLexerNumberType          = 57
	SyntaxFlowLexerBoolType            = 58
	SyntaxFlowLexerBoolLiteral         = 59
	SyntaxFlowLexerCheck               = 60
	SyntaxFlowLexerThen                = 61
	SyntaxFlowLexerDesc                = 62
	SyntaxFlowLexerElse                = 63
	SyntaxFlowLexerType                = 64
	SyntaxFlowLexerIn                  = 65
	SyntaxFlowLexerCall                = 66
	SyntaxFlowLexerConstant            = 67
	SyntaxFlowLexerPhi                 = 68
	SyntaxFlowLexerFormalParam         = 69
	SyntaxFlowLexerReturn              = 70
	SyntaxFlowLexerOpcode              = 71
	SyntaxFlowLexerHave                = 72
	SyntaxFlowLexerHaveAny             = 73
	SyntaxFlowLexerNot                 = 74
	SyntaxFlowLexerIdentifier          = 75
	SyntaxFlowLexerIdentifierChar      = 76
	SyntaxFlowLexerQuotedStringLiteral = 77
	SyntaxFlowLexerRegexpLiteral       = 78
	SyntaxFlowLexerWS                  = 79
)
