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
		"", "'->'", "'-->'", "';'", "'==>'", "'...'", "'%%'", "'..'", "'<='",
		"'>='", "'>>'", "'=>'", "'=='", "'=~'", "'!~'", "'&&'", "'||'", "'!='",
		"'?{'", "'-{'", "'}->'", "'#{'", "'#>'", "'#->'", "'>'", "'.'", "'<'",
		"'='", "'?'", "'('", "','", "')'", "'['", "']'", "'{'", "'}'", "'#'",
		"'$'", "':'", "'%'", "'!'", "'*'", "'-'", "'as'", "", "", "", "", "",
		"", "'str'", "'list'", "'dict'", "", "'bool'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "DeepFilter", "Deep", "Percent", "DeepDot", "LtEq",
		"GtEq", "DoubleGt", "Filter", "EqEq", "RegexpMatch", "NotRegexpMatch",
		"And", "Or", "NotEq", "ConditionStart", "DeepNextStart", "DeepNextEnd",
		"TopDefStart", "DefStart", "TopDef", "Gt", "Dot", "Lt", "Eq", "Question",
		"OpenParen", "Comma", "CloseParen", "ListSelectOpen", "ListSelectClose",
		"MapBuilderOpen", "MapBuilderClose", "ListStart", "DollarOutput", "Colon",
		"Search", "Bang", "Star", "Minus", "As", "WhiteSpace", "Number", "OctalNumber",
		"BinaryNumber", "HexNumber", "StringLiteral", "StringType", "ListType",
		"DictType", "NumberType", "BoolType", "BoolLiteral", "Identifier", "IdentifierChar",
		"RegexpLiteral",
	}
	staticData.ruleNames = []string{
		"T__0", "T__1", "T__2", "DeepFilter", "Deep", "Percent", "DeepDot",
		"LtEq", "GtEq", "DoubleGt", "Filter", "EqEq", "RegexpMatch", "NotRegexpMatch",
		"And", "Or", "NotEq", "ConditionStart", "DeepNextStart", "DeepNextEnd",
		"TopDefStart", "DefStart", "TopDef", "Gt", "Dot", "Lt", "Eq", "Question",
		"OpenParen", "Comma", "CloseParen", "ListSelectOpen", "ListSelectClose",
		"MapBuilderOpen", "MapBuilderClose", "ListStart", "DollarOutput", "Colon",
		"Search", "Bang", "Star", "Minus", "As", "WhiteSpace", "Number", "OctalNumber",
		"BinaryNumber", "HexNumber", "StringLiteral", "StringType", "ListType",
		"DictType", "NumberType", "BoolType", "BoolLiteral", "Identifier", "IdentifierChar",
		"IdentifierCharStart", "HexDigit", "Digit", "OctalDigit", "RegexpLiteral",
		"RegexpLiteralChar",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 58, 356, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
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
		62, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1,
		3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1,
		7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10,
		1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 14, 1,
		14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17,
		1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1,
		21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 22, 1, 23, 1, 23, 1, 24, 1, 24,
		1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 1, 27, 1, 28, 1, 28, 1, 29, 1, 29, 1,
		30, 1, 30, 1, 31, 1, 31, 1, 32, 1, 32, 1, 33, 1, 33, 1, 34, 1, 34, 1, 35,
		1, 35, 1, 36, 1, 36, 1, 37, 1, 37, 1, 38, 1, 38, 1, 39, 1, 39, 1, 40, 1,
		40, 1, 41, 1, 41, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 1, 43, 1, 44,
		4, 44, 247, 8, 44, 11, 44, 12, 44, 248, 1, 45, 1, 45, 1, 45, 1, 45, 4,
		45, 255, 8, 45, 11, 45, 12, 45, 256, 1, 46, 1, 46, 1, 46, 1, 46, 4, 46,
		263, 8, 46, 11, 46, 12, 46, 264, 1, 47, 1, 47, 1, 47, 1, 47, 4, 47, 271,
		8, 47, 11, 47, 12, 47, 272, 1, 48, 1, 48, 5, 48, 277, 8, 48, 10, 48, 12,
		48, 280, 9, 48, 1, 48, 1, 48, 1, 49, 1, 49, 1, 49, 1, 49, 1, 50, 1, 50,
		1, 50, 1, 50, 1, 50, 1, 51, 1, 51, 1, 51, 1, 51, 1, 51, 1, 52, 1, 52, 1,
		52, 1, 52, 1, 52, 1, 52, 1, 52, 1, 52, 3, 52, 306, 8, 52, 1, 53, 1, 53,
		1, 53, 1, 53, 1, 53, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1,
		54, 1, 54, 3, 54, 322, 8, 54, 1, 55, 1, 55, 5, 55, 326, 8, 55, 10, 55,
		12, 55, 329, 9, 55, 1, 56, 1, 56, 3, 56, 333, 8, 56, 1, 57, 3, 57, 336,
		8, 57, 1, 58, 1, 58, 1, 59, 1, 59, 1, 60, 1, 60, 1, 61, 1, 61, 4, 61, 346,
		8, 61, 11, 61, 12, 61, 347, 1, 61, 1, 61, 1, 62, 1, 62, 1, 62, 3, 62, 355,
		8, 62, 0, 0, 63, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17,
		9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17, 35,
		18, 37, 19, 39, 20, 41, 21, 43, 22, 45, 23, 47, 24, 49, 25, 51, 26, 53,
		27, 55, 28, 57, 29, 59, 30, 61, 31, 63, 32, 65, 33, 67, 34, 69, 35, 71,
		36, 73, 37, 75, 38, 77, 39, 79, 40, 81, 41, 83, 42, 85, 43, 87, 44, 89,
		45, 91, 46, 93, 47, 95, 48, 97, 49, 99, 50, 101, 51, 103, 52, 105, 53,
		107, 54, 109, 55, 111, 56, 113, 57, 115, 0, 117, 0, 119, 0, 121, 0, 123,
		58, 125, 0, 1, 0, 7, 3, 0, 10, 10, 13, 13, 32, 32, 1, 0, 96, 96, 1, 0,
		48, 57, 4, 0, 42, 42, 65, 90, 95, 95, 97, 122, 3, 0, 48, 57, 65, 70, 97,
		102, 1, 0, 48, 55, 1, 0, 47, 47, 361, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0,
		0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0,
		0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0,
		0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1,
		0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35,
		1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0,
		43, 1, 0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0, 0, 0,
		0, 51, 1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1, 0, 0,
		0, 0, 59, 1, 0, 0, 0, 0, 61, 1, 0, 0, 0, 0, 63, 1, 0, 0, 0, 0, 65, 1, 0,
		0, 0, 0, 67, 1, 0, 0, 0, 0, 69, 1, 0, 0, 0, 0, 71, 1, 0, 0, 0, 0, 73, 1,
		0, 0, 0, 0, 75, 1, 0, 0, 0, 0, 77, 1, 0, 0, 0, 0, 79, 1, 0, 0, 0, 0, 81,
		1, 0, 0, 0, 0, 83, 1, 0, 0, 0, 0, 85, 1, 0, 0, 0, 0, 87, 1, 0, 0, 0, 0,
		89, 1, 0, 0, 0, 0, 91, 1, 0, 0, 0, 0, 93, 1, 0, 0, 0, 0, 95, 1, 0, 0, 0,
		0, 97, 1, 0, 0, 0, 0, 99, 1, 0, 0, 0, 0, 101, 1, 0, 0, 0, 0, 103, 1, 0,
		0, 0, 0, 105, 1, 0, 0, 0, 0, 107, 1, 0, 0, 0, 0, 109, 1, 0, 0, 0, 0, 111,
		1, 0, 0, 0, 0, 113, 1, 0, 0, 0, 0, 123, 1, 0, 0, 0, 1, 127, 1, 0, 0, 0,
		3, 130, 1, 0, 0, 0, 5, 134, 1, 0, 0, 0, 7, 136, 1, 0, 0, 0, 9, 140, 1,
		0, 0, 0, 11, 144, 1, 0, 0, 0, 13, 147, 1, 0, 0, 0, 15, 150, 1, 0, 0, 0,
		17, 153, 1, 0, 0, 0, 19, 156, 1, 0, 0, 0, 21, 159, 1, 0, 0, 0, 23, 162,
		1, 0, 0, 0, 25, 165, 1, 0, 0, 0, 27, 168, 1, 0, 0, 0, 29, 171, 1, 0, 0,
		0, 31, 174, 1, 0, 0, 0, 33, 177, 1, 0, 0, 0, 35, 180, 1, 0, 0, 0, 37, 183,
		1, 0, 0, 0, 39, 186, 1, 0, 0, 0, 41, 190, 1, 0, 0, 0, 43, 193, 1, 0, 0,
		0, 45, 196, 1, 0, 0, 0, 47, 200, 1, 0, 0, 0, 49, 202, 1, 0, 0, 0, 51, 204,
		1, 0, 0, 0, 53, 206, 1, 0, 0, 0, 55, 208, 1, 0, 0, 0, 57, 210, 1, 0, 0,
		0, 59, 212, 1, 0, 0, 0, 61, 214, 1, 0, 0, 0, 63, 216, 1, 0, 0, 0, 65, 218,
		1, 0, 0, 0, 67, 220, 1, 0, 0, 0, 69, 222, 1, 0, 0, 0, 71, 224, 1, 0, 0,
		0, 73, 226, 1, 0, 0, 0, 75, 228, 1, 0, 0, 0, 77, 230, 1, 0, 0, 0, 79, 232,
		1, 0, 0, 0, 81, 234, 1, 0, 0, 0, 83, 236, 1, 0, 0, 0, 85, 238, 1, 0, 0,
		0, 87, 241, 1, 0, 0, 0, 89, 246, 1, 0, 0, 0, 91, 250, 1, 0, 0, 0, 93, 258,
		1, 0, 0, 0, 95, 266, 1, 0, 0, 0, 97, 274, 1, 0, 0, 0, 99, 283, 1, 0, 0,
		0, 101, 287, 1, 0, 0, 0, 103, 292, 1, 0, 0, 0, 105, 305, 1, 0, 0, 0, 107,
		307, 1, 0, 0, 0, 109, 321, 1, 0, 0, 0, 111, 323, 1, 0, 0, 0, 113, 332,
		1, 0, 0, 0, 115, 335, 1, 0, 0, 0, 117, 337, 1, 0, 0, 0, 119, 339, 1, 0,
		0, 0, 121, 341, 1, 0, 0, 0, 123, 343, 1, 0, 0, 0, 125, 354, 1, 0, 0, 0,
		127, 128, 5, 45, 0, 0, 128, 129, 5, 62, 0, 0, 129, 2, 1, 0, 0, 0, 130,
		131, 5, 45, 0, 0, 131, 132, 5, 45, 0, 0, 132, 133, 5, 62, 0, 0, 133, 4,
		1, 0, 0, 0, 134, 135, 5, 59, 0, 0, 135, 6, 1, 0, 0, 0, 136, 137, 5, 61,
		0, 0, 137, 138, 5, 61, 0, 0, 138, 139, 5, 62, 0, 0, 139, 8, 1, 0, 0, 0,
		140, 141, 5, 46, 0, 0, 141, 142, 5, 46, 0, 0, 142, 143, 5, 46, 0, 0, 143,
		10, 1, 0, 0, 0, 144, 145, 5, 37, 0, 0, 145, 146, 5, 37, 0, 0, 146, 12,
		1, 0, 0, 0, 147, 148, 5, 46, 0, 0, 148, 149, 5, 46, 0, 0, 149, 14, 1, 0,
		0, 0, 150, 151, 5, 60, 0, 0, 151, 152, 5, 61, 0, 0, 152, 16, 1, 0, 0, 0,
		153, 154, 5, 62, 0, 0, 154, 155, 5, 61, 0, 0, 155, 18, 1, 0, 0, 0, 156,
		157, 5, 62, 0, 0, 157, 158, 5, 62, 0, 0, 158, 20, 1, 0, 0, 0, 159, 160,
		5, 61, 0, 0, 160, 161, 5, 62, 0, 0, 161, 22, 1, 0, 0, 0, 162, 163, 5, 61,
		0, 0, 163, 164, 5, 61, 0, 0, 164, 24, 1, 0, 0, 0, 165, 166, 5, 61, 0, 0,
		166, 167, 5, 126, 0, 0, 167, 26, 1, 0, 0, 0, 168, 169, 5, 33, 0, 0, 169,
		170, 5, 126, 0, 0, 170, 28, 1, 0, 0, 0, 171, 172, 5, 38, 0, 0, 172, 173,
		5, 38, 0, 0, 173, 30, 1, 0, 0, 0, 174, 175, 5, 124, 0, 0, 175, 176, 5,
		124, 0, 0, 176, 32, 1, 0, 0, 0, 177, 178, 5, 33, 0, 0, 178, 179, 5, 61,
		0, 0, 179, 34, 1, 0, 0, 0, 180, 181, 5, 63, 0, 0, 181, 182, 5, 123, 0,
		0, 182, 36, 1, 0, 0, 0, 183, 184, 5, 45, 0, 0, 184, 185, 5, 123, 0, 0,
		185, 38, 1, 0, 0, 0, 186, 187, 5, 125, 0, 0, 187, 188, 5, 45, 0, 0, 188,
		189, 5, 62, 0, 0, 189, 40, 1, 0, 0, 0, 190, 191, 5, 35, 0, 0, 191, 192,
		5, 123, 0, 0, 192, 42, 1, 0, 0, 0, 193, 194, 5, 35, 0, 0, 194, 195, 5,
		62, 0, 0, 195, 44, 1, 0, 0, 0, 196, 197, 5, 35, 0, 0, 197, 198, 5, 45,
		0, 0, 198, 199, 5, 62, 0, 0, 199, 46, 1, 0, 0, 0, 200, 201, 5, 62, 0, 0,
		201, 48, 1, 0, 0, 0, 202, 203, 5, 46, 0, 0, 203, 50, 1, 0, 0, 0, 204, 205,
		5, 60, 0, 0, 205, 52, 1, 0, 0, 0, 206, 207, 5, 61, 0, 0, 207, 54, 1, 0,
		0, 0, 208, 209, 5, 63, 0, 0, 209, 56, 1, 0, 0, 0, 210, 211, 5, 40, 0, 0,
		211, 58, 1, 0, 0, 0, 212, 213, 5, 44, 0, 0, 213, 60, 1, 0, 0, 0, 214, 215,
		5, 41, 0, 0, 215, 62, 1, 0, 0, 0, 216, 217, 5, 91, 0, 0, 217, 64, 1, 0,
		0, 0, 218, 219, 5, 93, 0, 0, 219, 66, 1, 0, 0, 0, 220, 221, 5, 123, 0,
		0, 221, 68, 1, 0, 0, 0, 222, 223, 5, 125, 0, 0, 223, 70, 1, 0, 0, 0, 224,
		225, 5, 35, 0, 0, 225, 72, 1, 0, 0, 0, 226, 227, 5, 36, 0, 0, 227, 74,
		1, 0, 0, 0, 228, 229, 5, 58, 0, 0, 229, 76, 1, 0, 0, 0, 230, 231, 5, 37,
		0, 0, 231, 78, 1, 0, 0, 0, 232, 233, 5, 33, 0, 0, 233, 80, 1, 0, 0, 0,
		234, 235, 5, 42, 0, 0, 235, 82, 1, 0, 0, 0, 236, 237, 5, 45, 0, 0, 237,
		84, 1, 0, 0, 0, 238, 239, 5, 97, 0, 0, 239, 240, 5, 115, 0, 0, 240, 86,
		1, 0, 0, 0, 241, 242, 7, 0, 0, 0, 242, 243, 1, 0, 0, 0, 243, 244, 6, 43,
		0, 0, 244, 88, 1, 0, 0, 0, 245, 247, 3, 119, 59, 0, 246, 245, 1, 0, 0,
		0, 247, 248, 1, 0, 0, 0, 248, 246, 1, 0, 0, 0, 248, 249, 1, 0, 0, 0, 249,
		90, 1, 0, 0, 0, 250, 251, 5, 48, 0, 0, 251, 252, 5, 111, 0, 0, 252, 254,
		1, 0, 0, 0, 253, 255, 3, 121, 60, 0, 254, 253, 1, 0, 0, 0, 255, 256, 1,
		0, 0, 0, 256, 254, 1, 0, 0, 0, 256, 257, 1, 0, 0, 0, 257, 92, 1, 0, 0,
		0, 258, 259, 5, 48, 0, 0, 259, 260, 5, 98, 0, 0, 260, 262, 1, 0, 0, 0,
		261, 263, 2, 48, 49, 0, 262, 261, 1, 0, 0, 0, 263, 264, 1, 0, 0, 0, 264,
		262, 1, 0, 0, 0, 264, 265, 1, 0, 0, 0, 265, 94, 1, 0, 0, 0, 266, 267, 5,
		48, 0, 0, 267, 268, 5, 120, 0, 0, 268, 270, 1, 0, 0, 0, 269, 271, 3, 117,
		58, 0, 270, 269, 1, 0, 0, 0, 271, 272, 1, 0, 0, 0, 272, 270, 1, 0, 0, 0,
		272, 273, 1, 0, 0, 0, 273, 96, 1, 0, 0, 0, 274, 278, 5, 96, 0, 0, 275,
		277, 8, 1, 0, 0, 276, 275, 1, 0, 0, 0, 277, 280, 1, 0, 0, 0, 278, 276,
		1, 0, 0, 0, 278, 279, 1, 0, 0, 0, 279, 281, 1, 0, 0, 0, 280, 278, 1, 0,
		0, 0, 281, 282, 5, 96, 0, 0, 282, 98, 1, 0, 0, 0, 283, 284, 5, 115, 0,
		0, 284, 285, 5, 116, 0, 0, 285, 286, 5, 114, 0, 0, 286, 100, 1, 0, 0, 0,
		287, 288, 5, 108, 0, 0, 288, 289, 5, 105, 0, 0, 289, 290, 5, 115, 0, 0,
		290, 291, 5, 116, 0, 0, 291, 102, 1, 0, 0, 0, 292, 293, 5, 100, 0, 0, 293,
		294, 5, 105, 0, 0, 294, 295, 5, 99, 0, 0, 295, 296, 5, 116, 0, 0, 296,
		104, 1, 0, 0, 0, 297, 298, 5, 105, 0, 0, 298, 299, 5, 110, 0, 0, 299, 306,
		5, 116, 0, 0, 300, 301, 5, 102, 0, 0, 301, 302, 5, 108, 0, 0, 302, 303,
		5, 111, 0, 0, 303, 304, 5, 97, 0, 0, 304, 306, 5, 116, 0, 0, 305, 297,
		1, 0, 0, 0, 305, 300, 1, 0, 0, 0, 306, 106, 1, 0, 0, 0, 307, 308, 5, 98,
		0, 0, 308, 309, 5, 111, 0, 0, 309, 310, 5, 111, 0, 0, 310, 311, 5, 108,
		0, 0, 311, 108, 1, 0, 0, 0, 312, 313, 5, 116, 0, 0, 313, 314, 5, 114, 0,
		0, 314, 315, 5, 117, 0, 0, 315, 322, 5, 101, 0, 0, 316, 317, 5, 102, 0,
		0, 317, 318, 5, 97, 0, 0, 318, 319, 5, 108, 0, 0, 319, 320, 5, 115, 0,
		0, 320, 322, 5, 101, 0, 0, 321, 312, 1, 0, 0, 0, 321, 316, 1, 0, 0, 0,
		322, 110, 1, 0, 0, 0, 323, 327, 3, 115, 57, 0, 324, 326, 3, 113, 56, 0,
		325, 324, 1, 0, 0, 0, 326, 329, 1, 0, 0, 0, 327, 325, 1, 0, 0, 0, 327,
		328, 1, 0, 0, 0, 328, 112, 1, 0, 0, 0, 329, 327, 1, 0, 0, 0, 330, 333,
		7, 2, 0, 0, 331, 333, 3, 115, 57, 0, 332, 330, 1, 0, 0, 0, 332, 331, 1,
		0, 0, 0, 333, 114, 1, 0, 0, 0, 334, 336, 7, 3, 0, 0, 335, 334, 1, 0, 0,
		0, 336, 116, 1, 0, 0, 0, 337, 338, 7, 4, 0, 0, 338, 118, 1, 0, 0, 0, 339,
		340, 7, 2, 0, 0, 340, 120, 1, 0, 0, 0, 341, 342, 7, 5, 0, 0, 342, 122,
		1, 0, 0, 0, 343, 345, 5, 47, 0, 0, 344, 346, 3, 125, 62, 0, 345, 344, 1,
		0, 0, 0, 346, 347, 1, 0, 0, 0, 347, 345, 1, 0, 0, 0, 347, 348, 1, 0, 0,
		0, 348, 349, 1, 0, 0, 0, 349, 350, 5, 47, 0, 0, 350, 124, 1, 0, 0, 0, 351,
		352, 5, 92, 0, 0, 352, 355, 5, 47, 0, 0, 353, 355, 8, 6, 0, 0, 354, 351,
		1, 0, 0, 0, 354, 353, 1, 0, 0, 0, 355, 126, 1, 0, 0, 0, 13, 0, 248, 256,
		264, 272, 278, 305, 321, 327, 332, 335, 347, 354, 1, 6, 0, 0,
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
	SyntaxFlowLexerT__0            = 1
	SyntaxFlowLexerT__1            = 2
	SyntaxFlowLexerT__2            = 3
	SyntaxFlowLexerDeepFilter      = 4
	SyntaxFlowLexerDeep            = 5
	SyntaxFlowLexerPercent         = 6
	SyntaxFlowLexerDeepDot         = 7
	SyntaxFlowLexerLtEq            = 8
	SyntaxFlowLexerGtEq            = 9
	SyntaxFlowLexerDoubleGt        = 10
	SyntaxFlowLexerFilter          = 11
	SyntaxFlowLexerEqEq            = 12
	SyntaxFlowLexerRegexpMatch     = 13
	SyntaxFlowLexerNotRegexpMatch  = 14
	SyntaxFlowLexerAnd             = 15
	SyntaxFlowLexerOr              = 16
	SyntaxFlowLexerNotEq           = 17
	SyntaxFlowLexerConditionStart  = 18
	SyntaxFlowLexerDeepNextStart   = 19
	SyntaxFlowLexerDeepNextEnd     = 20
	SyntaxFlowLexerTopDefStart     = 21
	SyntaxFlowLexerDefStart        = 22
	SyntaxFlowLexerTopDef          = 23
	SyntaxFlowLexerGt              = 24
	SyntaxFlowLexerDot             = 25
	SyntaxFlowLexerLt              = 26
	SyntaxFlowLexerEq              = 27
	SyntaxFlowLexerQuestion        = 28
	SyntaxFlowLexerOpenParen       = 29
	SyntaxFlowLexerComma           = 30
	SyntaxFlowLexerCloseParen      = 31
	SyntaxFlowLexerListSelectOpen  = 32
	SyntaxFlowLexerListSelectClose = 33
	SyntaxFlowLexerMapBuilderOpen  = 34
	SyntaxFlowLexerMapBuilderClose = 35
	SyntaxFlowLexerListStart       = 36
	SyntaxFlowLexerDollarOutput    = 37
	SyntaxFlowLexerColon           = 38
	SyntaxFlowLexerSearch          = 39
	SyntaxFlowLexerBang            = 40
	SyntaxFlowLexerStar            = 41
	SyntaxFlowLexerMinus           = 42
	SyntaxFlowLexerAs              = 43
	SyntaxFlowLexerWhiteSpace      = 44
	SyntaxFlowLexerNumber          = 45
	SyntaxFlowLexerOctalNumber     = 46
	SyntaxFlowLexerBinaryNumber    = 47
	SyntaxFlowLexerHexNumber       = 48
	SyntaxFlowLexerStringLiteral   = 49
	SyntaxFlowLexerStringType      = 50
	SyntaxFlowLexerListType        = 51
	SyntaxFlowLexerDictType        = 52
	SyntaxFlowLexerNumberType      = 53
	SyntaxFlowLexerBoolType        = 54
	SyntaxFlowLexerBoolLiteral     = 55
	SyntaxFlowLexerIdentifier      = 56
	SyntaxFlowLexerIdentifierChar  = 57
	SyntaxFlowLexerRegexpLiteral   = 58
)
