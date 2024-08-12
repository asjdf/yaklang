package ssa

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/omap"
	"github.com/yaklang/yaklang/common/yak/ssa/ssadb"
	"go.uber.org/atomic"

	syncAtomic "sync/atomic"
)

var cachePool = omap.NewEmptyOrderedMap[string, *Cache]()

func GetCacheFromPool(programName string) *Cache {
	if cache, ok := cachePool.Get(programName); ok {
		return cache
	}
	cache := NewDBCache(programName, true)
	cachePool.Set(programName, cache)
	return cache
}

type instructionIrCode struct {
	inst   Instruction
	irCode *ssadb.IrCode
}

// Cache : a cache in middle layer of database and application.
//
//	application will Get/Set Instruction,
//
// and save the data to database when the data is expired,
// and load the data from database when the data is not in cache.
type Cache struct {
	ProgramName      string // mark which program handled
	DB               *gorm.DB
	id               *atomic.Int64
	InstructionCache *utils.CacheWithKey[int64, instructionIrCode] // instructionID to instruction

	VariableCache   map[string][]Instruction // variable(name:string) to []instruction
	MemberCache     map[string][]Instruction
	Class2InstIndex map[string][]Instruction
}

// NewDBCache : create a new ssa db cache. if ttl is 0, the cache will never expire, and never save to database.
func NewDBCache(programName string, databaseEnable bool, ConfigTTL ...time.Duration) *Cache {
	ttl := time.Duration(0)
	if databaseEnable {
		// enable database
		ttl = time.Second * 8
	}

	cache := &Cache{
		ProgramName:      programName,
		InstructionCache: utils.NewTTLCacheWithKey[int64, instructionIrCode](ttl),
		VariableCache:    make(map[string][]Instruction),
		MemberCache:      make(map[string][]Instruction),
		Class2InstIndex:  make(map[string][]Instruction),
	}

	if databaseEnable {
		cache.DB = ssadb.GetDB().Where("program_name = ?", programName)
		cache.InstructionCache.SetCheckExpirationCallback(func(key int64, inst instructionIrCode) bool {
			return cache.saveInstruction(inst)
		})
	} else {
		cache.id = atomic.NewInt64(0)
	}

	return cache
}

func (c *Cache) HaveDatabaseBackend() bool {
	return c.DB != nil && c.ProgramName != ""
}

// =============================================== Instruction =======================================================

// SetInstruction : set instruction to cache.
func (c *Cache) SetInstruction(inst Instruction) {
	if inst.GetId() != -1 {
		return
	}

	var id int64
	var instIr instructionIrCode
	if c.HaveDatabaseBackend() {
		// use database
		rawID, irCode := ssadb.RequireIrCode(c.DB, c.ProgramName)
		id = int64(rawID)
		instIr = instructionIrCode{
			inst:   inst,
			irCode: irCode,
		}
	} else {
		// not use database
		id = c.id.Inc()
		instIr = instructionIrCode{
			inst:   inst,
			irCode: nil,
		}
	}
	inst.SetId(id)
	c.InstructionCache.Set(id, instIr)
}

func (c *Cache) DeleteInstruction(inst Instruction) {
	c.InstructionCache.Remove(inst.GetId())
}

// GetInstruction : get instruction from cache.
func (c *Cache) GetInstruction(id int64) Instruction {
	ret, ok := c.InstructionCache.Get(id)
	if !ok && c.HaveDatabaseBackend() {
		// if no in cache, get from database
		// if found in database, create a new lazy instruction
		return c.newLazyInstruction(id)
		// all instruction from database will be lazy instruction
	}
	return ret.inst
}

// =============================================== Variable =======================================================

func (c *Cache) AddVariable(name string, inst Instruction) {
	member := ""
	// field
	if strings.HasPrefix(name, "#") { // member-call variable contain #, see common/yak/ssa/member_call.go:checkCanMemberCall
		if _, memberName, ok := strings.Cut(name, "."); ok {
			member = memberName
		}

		if _, memberKey, ok := strings.Cut(name, "["); ok {
			member, _ = strings.CutSuffix(memberKey, "]")
		}
	}

	if member != "" {
		c.MemberCache[member] = append(c.MemberCache[member], inst)
	} else {
		c.VariableCache[name] = append(c.VariableCache[name], inst)
	}
	if c.HaveDatabaseBackend() {
		SaveVariableIndex(inst, name, member)
	}
}

func (c *Cache) RemoveVariable(name string, inst Instruction) {
	member := ""
	// field
	if strings.HasPrefix(name, "#") { // member-call variable contain #, see common/yak/ssa/member_call.go:checkCanMemberCall
		if _, memberName, ok := strings.Cut(name, "."); ok {
			member = memberName
		}

		if _, memberKey, ok := strings.Cut(name, "["); ok {
			member, _ = strings.CutSuffix(memberKey, "]")
		}
	}

	if member != "" {
		c.MemberCache[member] = utils.RemoveSliceItem(c.MemberCache[member], inst)
	} else {
		c.VariableCache[name] = utils.RemoveSliceItem(c.VariableCache[name], inst)
	}
}

func (c *Cache) AddClassInstance(name string, inst Instruction) {
	c.Class2InstIndex[name] = append(c.Class2InstIndex[name], inst)
	if c.HaveDatabaseBackend() {
		SaveClassIndex(inst, name)
	}
}

// =============================================== Database =======================================================
// only LazyInstruction and false marshal will not be saved to database
func (c *Cache) saveInstruction(instIr instructionIrCode) bool {
	start := time.Now()
	if !c.HaveDatabaseBackend() {
		log.Errorf("BUG: saveInstruction called when DB is nil")
		return false
	}

	// all instruction from database will be lazy instruction
	if lz, ok := ToLazyInstruction(instIr.inst); ok {
		// we just check if this lazy-instruction should be saved again?
		if !lz.ShouldSave() {
			return false
		}
	}

	err := Instruction2IrCode(instIr.inst, instIr.irCode)
	if err != nil {
		log.Errorf("FitIRCode error: %s", err)
		return false
	}

	if instIr.irCode.Opcode == 0 {
		log.Errorf("BUG: saveInstruction called with empty opcode: %v", instIr.inst.GetName())
	}
	if err := c.DB.Save(instIr.irCode).Error; err != nil {
		log.Errorf("Save irCode error: %v", err)
	}
	syncAtomic.AddUint64(&_SSASaveIrCodeCost, uint64(time.Now().Sub(start).Nanoseconds()))

	return true
}

func (c *Cache) SaveToDatabase() {
	if !c.HaveDatabaseBackend() {
		return
	}

	start := time.Now()
	defer func() {
		syncAtomic.AddUint64(&_SSACacheToDatabaseCost, uint64(time.Now().Sub(start).Nanoseconds()))
	}()

	insturctionCache := c.InstructionCache.GetAll()
	c.InstructionCache.Close()

	for _, instIR := range insturctionCache {
		c.saveInstruction(instIR)
	}
}

func (c *Cache) IsExistedSourceCodeHash(programName string, hashString string) bool {
	if programName == "" || !c.HaveDatabaseBackend() {
		return false
	}

	var count int
	if ret := c.DB.Model(&ssadb.IrCode{}).Where(
		"source_code_hash = ?", hashString,
	).Where(
		"program_name = ?", programName,
	).Count(&count).Error; ret != nil {
		log.Warnf("IsExistedSourceCodeHash error: %v", ret)
	}
	return count > 0
}
