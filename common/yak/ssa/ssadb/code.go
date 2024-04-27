package ssadb

import (
	"encoding/json"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
	"sync"

	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"

	"github.com/jinzhu/gorm"
)

type IrSourceCodeResource struct {
	gorm.Model
	SourceCodeHash     string `json:"source_code_hash" gorm:"index"`
	SourceCodeFilename string `json:"source_code_file"`
	QuotedSourceCode   string `json:"quoted_source_code"`
}

type IrCode struct {
	gorm.Model

	ProgramName string `json:"program_name" gorm:"index"`
	PackageName string `json:"package_name" gorm:"index"`

	// source code
	SourceCodeStartOffset int64  `json:"source_code_start_offset"`
	SourceCodeEndOffset   int64  `json:"source_code_end_offset"`
	SourceCodeHash        string `json:"source_code_hash"`

	// opcode
	Opcode     int64  `json:"opcode"`
	OpcodeName string `json:"opcode_name"`

	// just for binary and unary operator
	OpcodeOperator string `json:"opcode_operator"`

	// basic info
	Name             string `json:"name"`
	VerboseName      string `json:"verbose_name"`
	ShortVerboseName string `json:"short_verbose_name"`

	// any IrCode in one block inner one  function
	CurrentBlock    int64 `json:"current_block"`
	CurrentFunction int64 `json:"current_function"`

	// FunctionDefs
	IsFunction       bool       `json:"is_function"`
	FormalArgs       Int64Slice `json:"formal_args" gorm:"type:text"`
	IsVariadic       bool       `json:"is_variadic"`
	ReturnCodes      Int64Slice `json:"return_codes" gorm:"type:text"`
	SideEffects      Int64Slice `json:"side_effects" gorm:"type:text"`
	IsExternal       bool       `json:"is_external"`
	CodeBlocks       Int64Slice `json:"code_blocks" gorm:"type:text"`
	EnterBlock       int64      `json:"enter_block"`
	ExitBlock        int64      `json:"exit_block"`
	DeferBlock       int64      `json:"defer_block"`
	ChildrenFunction Int64Slice `json:"children_function" gorm:"type:text"`
	ParentFunction   int64      `json:"parent_function"`

	// block
	IsBlock   bool       `json:"is_block"`
	PredBlock Int64Slice `json:"pred_block" gorm:"type:text"`
	SuccBlock Int64Slice `json:"succ_block" gorm:"type:text"`
	Phis      Int64Slice `json:"phis_in_block" gorm:"type:text"`

	// Use-Def Chains Relation
	Defs  Int64Slice `json:"defs" gorm:"type:text"`
	Users Int64Slice `json:"users" gorm:"type:text"`
	// this is user is call and method is this IR self
	CalledBy Int64Slice `json:"is_called_by" gorm:"type:text"`

	// OOP Supporting
	IsObject       bool
	ObjectMembers  Int64Map `json:"object_members" gorm:"type:text"`
	IsObjectMember bool
	ObjectParent   int64 `json:"object_parent"`
	ObjectKey      int64 `json:"object_key"`

	// Maskable
	MaskedCodes Int64Slice `json:"masked_codes" gorm:"type:text"`
	IsMasked    bool       `json:"is_masked"`

	// Call instruction
	IsCalled       bool       `json:"is_called"`
	Method         int64      `json:"method"`
	ActualArgs     Int64Slice `json:"actual_args" gorm:"type:text"`
	BingDings      Int64Map   `json:"bing_dings" gorm:"type:text"`
	ArgumentMember Int64Map   `json:"argument_member" gorm:"type:text"`

	// Variable
	Variable StringSlice `json:"variable" gorm:"type:text"`

	// Constant
	ConstantValue string `json:"constant_value" gorm:"type:text"`

	// compile hash means: hash[ (file-content)+(program-name)+(package-name)+(program-index) ]
	ProgramCompileHash string `json:"program_compile_hash" gorm:"index"`

	// not important information
	ExtraInformation string `json:"extra_information"`
}

func emptyIrCode() *IrCode {
	return &IrCode{}
}

var verifyExisted = new(sync.Once)

func RequireIrCode(db *gorm.DB, program string) (uint, *IrCode) {
	verifyExisted.Do(func() {
		db.AutoMigrate(&IrCode{})
	})
	db = db.Model(&IrCode{})
	// save new ircode
	ircode := emptyIrCode()
	ircode.ProgramName = program
	db.Create(ircode)
	return ircode.ID, ircode
}

func NewLazyInstruction(i any) (*LazyInstruction, error) {
	codec.Atoi(codec.AnyToString(i))
}

func GetIrCodeById(db *gorm.DB, id int64) *IrCode {
	if id == -1 {
		return nil
	}
	return db.Model(&IrCode{}).Where("id = ?", id).First(&IrCode{}).Value.(*IrCode)
}

func GetIrByVariable(db *gorm.DB, program, name string) []*IrCode {
	var irVariable IrVariable
	if err := db.Model(&IrVariable{}).Where("variable_name = ? AND program_name = ?", name, program).First(&irVariable).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
	}

	ret := make([]*IrCode, 0, len(irVariable.InstructionID))
	for _, id := range irVariable.InstructionID {
		r := GetIrCodeById(db, id)
		ret = append(ret, r)
	}

	return ret
}

func DeleteProgram(db *gorm.DB, program string) {
	db.Model(&IrCode{}).Where("program_name = ?", program).Unscoped().Delete(&IrCode{})
}

func (r *IrCode) IsEmptySourceCodeHash() bool {
	if r == nil {
		return true
	}
	if r.SourceCodeHash == "" {
		return true
	}
	for _, hash := range []func(...any) string{
		utils.CalcMd5,
		utils.CalcSha1,
		utils.CalcSha256,
	} {
		if r.SourceCodeHash == hash("") {
			return true
		}
	}
	return false
}

func (r *IrCode) GetIdInt64() int64 {
	return int64(r.ID)
}

func (r *IrCode) GetIdInt() int {
	return int(r.ID)
}

func (r *IrCode) SetExtraInfo(params map[string]any) {
	results, err := json.Marshal(params)
	if err != nil {
		log.Warnf("BUG: cannot fetch instruction's extra info: %v", err)
	}
	r.ExtraInformation = string(results)
}

func (r *IrCode) GetExtraInfo() map[string]any {
	results := make(map[string]any)
	if err := json.Unmarshal([]byte(r.ExtraInformation), &results); err != nil {
		log.Warnf("BUG: cannot fetch instruction's extra info: %v", err)
	}
	return results
}
