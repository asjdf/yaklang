package ssadb

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils/bizhelper"
)

func YieldIrCodesProgramName(db *gorm.DB, ctx context.Context, program string) chan *IrCode {
	db = db.Model(&IrCode{}).Where("program_name = ?", program)
	return yieldIrCodes(db, ctx)
}

func yieldIrCodes(db *gorm.DB, ctx context.Context) chan *IrCode {
	db = db.Model(&IrCode{})
	outC := make(chan *IrCode)
	go func() {
		defer close(outC)

		var page = 1
		for {
			var items []*IrCode
			if _, b := bizhelper.Paging(db, page, 100, &items); b.Error != nil {
				log.Errorf("paging failed: %s", b.Error)
				return
			}

			page++
			for _, d := range items {
				select {
				case <-ctx.Done():
					return
				case outC <- d:
				}
			}

			if len(items) < 100 {
				return
			}
		}
	}()
	return outC
}

func yieldIrIndex(db *gorm.DB, ctx context.Context) chan int64 {
	db = db.Model(&IrIndex{})
	outC := make(chan int64)
	go func() {
		defer close(outC)

		filter := make(map[int64]struct{})

		var page = 1
		for {
			var items []*IrIndex
			if _, b := bizhelper.Paging(db, page, 100, &items); b.Error != nil {
				log.Errorf("paging failed: %s", b.Error)
				return
			}

			page++
			for _, d := range items {
				id := d.ValueID
				if _, ok := filter[id]; ok {
					continue
				}
				filter[id] = struct{}{}

				select {
				case <-ctx.Done():
					return
				case outC <- id:
				}
			}

			if len(items) < 100 {
				return
			}
		}
	}()
	return outC
}

// type MatchMode int
const (
	NameMatch int = 1
	KeyMatch      = 1 << 1
	BothMatch     = NameMatch | KeyMatch
)

const (
	ExactCompare int = iota
	GlobCompare
	RegexpCompare
)

func SearchVariable(db *gorm.DB, compareMode, matchMod int, value string) chan int64 {
	switch compareMode {
	case ExactCompare:
		return ExactSearchVariable(db, matchMod, value)
	case GlobCompare:
		return GlobSearchVariable(db, matchMod, value)
	case RegexpCompare:
		return RegexpSearchVariable(db, matchMod, value)
	}
	return nil
}

func ExactSearchVariable(DB *gorm.DB, mod int, value string) chan int64 {
	db := DB.Model(&IrIndex{})
	switch mod {
	case NameMatch:
		db = db.Where("variable_name = ? OR class_name = ?", value, value)
	case KeyMatch:
		db = db.Where("field_name = ?", value, value)
	case BothMatch:
		db = db.Where("variable_name = ? OR class_name = ? OR field_name = ?", value, value, value)
	}
	return yieldIrIndex(db, context.Background())
}

func GlobSearchVariable(DB *gorm.DB, mod int, value string) chan int64 {
	db := DB.Model(&IrIndex{})
	switch mod {
	case NameMatch:
		db = db.Where("variable_name GLOB ?", value)
	case KeyMatch:
		db = db.Where("field_name GLOB ?", value, value)
	case BothMatch:
		db = db.Where("variable_name GLOB ? OR class_name GLOB ? OR field_name GLOB ?", value, value, value)
	}
	return yieldIrIndex(db, context.Background())
}
func RegexpSearchVariable(DB *gorm.DB, mod int, value string) chan int64 {
	db := DB.Model(&IrIndex{})
	switch mod {
	case NameMatch:
		db = db.Where("variable_name REGEXP ? OR class_name REGEXP ?", value, value)
	case KeyMatch:
		db = db.Where("field_name REGEXP ?", value)
	case BothMatch:
		db = db.Where("variable_name REGEXP ? OR class_name REGEXP ? OR field_name REGEXP ?", value, value, value)
	}
	return yieldIrIndex(db, context.Background())
}

func GetVariableByValue(valueID int64) []*IrIndex {
	db := GetDB()
	var ir []*IrIndex
	if err := db.Model(&IrIndex{}).Where("value_id = ? and variable_name != ?", valueID, "").Find(&ir).Error; err != nil {
		return nil
	}
	return ir
}
