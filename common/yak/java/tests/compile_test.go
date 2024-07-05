package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/syntaxflow/sfvm"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yak/ssa"
	"github.com/yaklang/yaklang/common/yak/ssa/ssadb"
	"github.com/yaklang/yaklang/common/yak/ssaapi"
	"golang.org/x/exp/slices"
)

func comparePackage(prog *ssadb.IrProgram) error {
	// compare package-library with database
	var out ssadb.IrProgram
	err := ssadb.GetDB().Model(&ssadb.IrProgram{}).Where(
		"program_name = ?", prog.ProgramName,
	).First(&out).Error
	if err != nil {
		return utils.Errorf("get program %s error : %v", prog.ProgramName, err)
	}
	log.Infof("out pkgName: %v", out)
	if out.ProgramName != prog.ProgramName {
		return fmt.Errorf("program name not match want %v, got %v", prog.ProgramName, out.ProgramName)
	}
	if prog.ProgramKind != out.ProgramKind {
		return fmt.Errorf("program kind not match want %v, got %v", prog.ProgramKind, out.ProgramKind)
	}
	if slices.Compare(prog.UpStream, out.UpStream) != 0 {
		return fmt.Errorf("upstream not match want %v, got %v", prog.UpStream, out.UpStream)
	}
	if slices.Compare(prog.DownStream, out.DownStream) != 0 {
		return fmt.Errorf("downstream not match want %v, got %v", prog.DownStream, out.DownStream)
	}
	return nil
}

func TestCompileProgram_OnlySource(t *testing.T) {
	pkgName := "a" + strings.ReplaceAll(uuid.NewString(), "-", "")
	code := fmt.Sprintf(`
package %s; 
public class A {
	int a;
} 
	`, pkgName)
	ssadb.DeleteProgram(ssadb.GetDB(), pkgName)

	// compile only source
	prog, err := ssaapi.Parse(code, ssaapi.WithLanguage(ssaapi.JAVA))
	assert.NoError(t, err)
	prog.Show()
	assert.Equal(t, 1, len(prog.Program.UpStream))
	assert.Equal(t, 0, len(prog.Program.DownStream))
	if slices.Contains(ssadb.AllPrograms(ssadb.GetDB()), pkgName) {
		t.Fatalf("package %s should not be in the database", pkgName)
	}
}
func TestCompileProgram_WithDatabase(t *testing.T) {
	pkgName := "a" + strings.ReplaceAll(uuid.NewString(), "-", "")
	code := fmt.Sprintf(`
package %s; 
public class A {
	int a;
} 
	`, pkgName)
	ssadb.DeleteProgram(ssadb.GetDB(), pkgName)
	// compile with database
	programId := uuid.NewString()
	ssadb.DeleteProgram(ssadb.GetDB(), programId)

	prog, err := ssaapi.Parse(code,
		ssaapi.WithLanguage(ssaapi.JAVA),
		ssaapi.WithDatabaseProgramName(programId),
	)
	assert.NoError(t, err)

	defer func() {
		ssadb.DeleteProgram(ssadb.GetDB(), programId)
		ssadb.DeleteProgram(ssadb.GetDB(), pkgName)
	}()

	prog.Show()
	assert.Equal(t, 1, len(prog.Program.UpStream))
	assert.Equal(t, 0, len(prog.Program.DownStream))
	{
		// check in ir-code table
		var programsInIrCode []string
		ssadb.GetDB().Model(&ssadb.IrCode{}).Select("DISTINCT(program_name)").Pluck("program_name", &programsInIrCode)
		if !slices.Contains(programsInIrCode, pkgName) {
			t.Fatalf("package %s should be in the ir-code table", pkgName)
		}

		// check in ir-program table
		var programsInIrProgram []string
		ssadb.GetDB().Model(&ssadb.IrProgram{}).Select("DISTINCT(program_name)").Pluck("program_name", &programsInIrProgram)
		if !slices.Contains(programsInIrProgram, pkgName) {
			t.Fatalf("package %s should be in the ir-program table", pkgName)
		}
	}

	// compare program-application with database
	if err := comparePackage(&ssadb.IrProgram{
		ProgramName: programId,
		ProgramKind: string(ssa.Application),
		UpStream:    []string{pkgName},
		DownStream:  []string{},
	}); err != nil {
		t.Fatalf("compare package failed: %v", err)
	}
	// compare package-library with database
	if err := comparePackage(&ssadb.IrProgram{
		ProgramName: pkgName,
		ProgramKind: string(ssa.Library),
		UpStream:    []string{},
		DownStream:  []string{programId},
	}); err != nil {
		t.Fatalf("compare package failed: %v", err)
	}
}

func TestCompileProgram_OnlyDatabase(t *testing.T) {
	pkgName := "a" + strings.ReplaceAll(uuid.NewString(), "-", "")
	programId := uuid.NewString()
	{
		code := fmt.Sprintf(`
	package %s; 
	public class A {
		int a;
	} 
		`, pkgName)
		ssadb.DeleteProgram(ssadb.GetDB(), pkgName)
		// compile with database
		ssadb.DeleteProgram(ssadb.GetDB(), programId)

		_, err := ssaapi.Parse(code,
			ssaapi.WithLanguage(ssaapi.JAVA),
			ssaapi.WithDatabaseProgramName(programId),
		)
		assert.NoError(t, err)

		defer func() {
			ssadb.DeleteProgram(ssadb.GetDB(), programId)
			ssadb.DeleteProgram(ssadb.GetDB(), pkgName)
		}()
	}

	prog, err := ssaapi.FromDatabase(programId)
	assert.NoError(t, err)
	prog.Show()
	// assert.Equal(t, 1, len(prog.Program.UpStream))
	// assert.Equal(t, 0, len(prog.Program.DownStream))
	// for name := range prog.Program.UpStream {
	// 	if name != pkgName {
	// 		t.Fatalf("upstream should be %s, but got %s", pkgName, name)
	// 	}
	// }
}

func TestCompileProgram_Delete(t *testing.T) {

	pkgName := "a" + strings.ReplaceAll(uuid.NewString(), "-", "")
	programId := uuid.NewString()
	// compile
	{
		code := fmt.Sprintf(`
	package %s; 
	public class A {
		int a;
	} 
		`, pkgName)
		ssadb.DeleteProgram(ssadb.GetDB(), pkgName)
		// compile with database
		ssadb.DeleteProgram(ssadb.GetDB(), programId)

		_, err := ssaapi.Parse(code,
			ssaapi.WithLanguage(ssaapi.JAVA),
			ssaapi.WithDatabaseProgramName(programId),
		)
		assert.NoError(t, err)
	}

	// delete
	ssadb.DeleteProgram(ssadb.GetDB(), programId)

	// check in ir-code table
	var programsInIrCode []string
	ssadb.GetDB().Model(&ssadb.IrCode{}).Select("DISTINCT(program_name)").Pluck("program_name", &programsInIrCode)
	if slices.Contains(programsInIrCode, pkgName) {
		t.Fatalf("package %s should not be in the ir-code table", pkgName)
	}

	// check in ir-program table
	var programsInIrProgram []string
	ssadb.GetDB().Model(&ssadb.IrProgram{}).Select("DISTINCT(program_name)").Pluck("program_name", &programsInIrProgram)
	if slices.Contains(programsInIrProgram, pkgName) {
		t.Fatalf("package %s should not be in the ir-program table", pkgName)
	}
}
func TestCompileProgram_ReUseLibrary(t *testing.T) {

	pkgName := "a" + strings.ReplaceAll(uuid.NewString(), "-", "")
	code := fmt.Sprintf(`
	package %s; 
	public class A {
		public static void main(String[] args) {
			int a = 1;
		}
	} 
		`, pkgName)

	// compile with database
	// compile
	programID1 := uuid.NewString()
	{
		ssadb.DeleteProgram(ssadb.GetDB(), programID1)
		_, err := ssaapi.Parse(code,
			ssaapi.WithLanguage(ssaapi.JAVA),
			ssaapi.WithDatabaseProgramName(programID1),
		)
		defer ssadb.DeleteProgram(ssadb.GetDB(), programID1)
		assert.NoError(t, err)
	}
	if err := comparePackage(&ssadb.IrProgram{
		ProgramName: pkgName,
		ProgramKind: string(ssa.Library),
		UpStream:    []string{},
		DownStream:  []string{programID1},
	}); err != nil {
		t.Fatalf("compare package failed: %v", err)
	}
	programID2 := uuid.NewString()
	{
		ssadb.DeleteProgram(ssadb.GetDB(), programID2)
		_, err := ssaapi.Parse(code,
			ssaapi.WithLanguage(ssaapi.JAVA),
			ssaapi.WithDatabaseProgramName(programID2),
		)
		defer ssadb.DeleteProgram(ssadb.GetDB(), programID2)
		assert.NoError(t, err)
	}

	// check ir-package re-use
	if err := comparePackage(&ssadb.IrProgram{
		ProgramName: pkgName,
		ProgramKind: string(ssa.Library),
		UpStream:    []string{},
		DownStream:  []string{programID1, programID2},
	}); err != nil {
		t.Fatalf("compare package failed: %v", err)
	}

	// query
	prog, err := ssaapi.FromDatabase(programID2)
	assert.NoError(t, err)

	res := prog.SyntaxFlow(`a as $a`, sfvm.WithEnableDebug())
	assert.Equal(t,
		[]string{"1"},
		lo.Map(
			res.GetValues("a"),
			func(v *ssaapi.Value, _ int) string { return v.String() },
		),
	)

	// delete
	// delete program 1, will not affect program 2
	ssadb.DeleteProgram(ssadb.GetDB(), programID1)
	if err := comparePackage(&ssadb.IrProgram{
		ProgramName: pkgName,
		ProgramKind: string(ssa.Library),
		UpStream:    []string{},
		DownStream:  []string{programID2},
	}); err != nil {
		t.Fatalf("compare package failed: %v", err)
	}

	// delete program 2, pkg should be deleted
	ssadb.DeleteProgram(ssadb.GetDB(), programID2)
	var programsInIrProgram []string
	ssadb.GetDB().Model(&ssadb.IrProgram{}).Select("DISTINCT(program_name)").Pluck("program_name", &programsInIrProgram)
	if slices.Contains(programsInIrProgram, pkgName) {
		t.Fatalf("package %s should not be in the ir-program table", pkgName)
	}

}
