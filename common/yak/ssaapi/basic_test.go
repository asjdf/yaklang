package ssaapi

import (
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yak/ssa"
)

func TestYaklangBasic_variable(t *testing.T) {
	prog, err := Parse(`
	a = 1
	f = (a) => {
	}
	`)
	if err != nil {
		t.Fatal("prog parse error", err)
	}
	prog.Ref("a").ShowWithSource()
}

func TestYaklangBasic_Const(t *testing.T) {
	code := `
	a = 1
	b = a + 1 
	println(b)
	`

	prog, err := Parse(code)
	if err != nil {
		t.Fatal("prog parse error", err)
	}
	// prog.Ref("a").Show()
	prog.Ref("b").Show().ForEach(func(v *Value) {
		if len(v.GetOperands().Show()) != 2 {
			t.Fatalf("const 2 should have 2 operands")
		}
	})
}

func TestYaklangBasic_RecursivePhi_1(t *testing.T) {
	const code = `
count = 100

b = 1
a = (ffff) => {
	b ++
	if b > 100 {
		return
	}
	for i = 0; i < b; i ++ {
		dump(b)
	}
	c = () => { d = a(b); sink(d) }
	c()
}
e = a(1)
`
	prog, err := Parse(code)
	if err != nil {
		t.Fatal("prog parse error", err)
	}
	prog.Show()
	prog.Ref("a")
}

func TestYaklangBasic_RecursivePhi_2(t *testing.T) {
	const code = `
count = 100

a = (b) => {
	b ++
	if b > 100 {
		return
	}
	for i = 0; i < b; i ++ {
		b := virtual(i, b)
		dump(b)
	}
	a(b)
}
e = a(1)          // e
f = a(v2(e))      // f
`
	prog, err := Parse(code)
	if err != nil {
		t.Fatal("prog parse error", err)
	}
	prog.Show()
	prog.Ref("a")
}

func TestYaklangBasic_DoublePhi(t *testing.T) {
	const code = `var a = 1; for i:=0; i<n; i ++ { a += i }; println(a)`
	prog, err := Parse(code)
	if err != nil {
		t.Fatal("prog parse error", err)
	}
	prog.Show()
	prog.Ref("a")
}

func TestYaklangBasic_Used(t *testing.T) {
	token := utils.RandStringBytes(10)
	prog, err := Parse(`var a, b
` + token + `(a)
`)
	if err != nil {
		t.Fatal("prog parse error", err)
	}
	traceFinished := false
	prog.Ref("a").ForEach(func(value *Value) {
		value.GetUsers().ForEach(func(value *Value) {
			log.Infof("a's uses include: %v", value.String())
			if strings.Contains(value.String(), token+"(") {
				traceFinished = true
			}
		})
	})
	if !traceFinished {
		t.Error("trace failed: var cannot trace to call actual arguments")
	}
}

func TestYaklangBasic_if_phi(t *testing.T) {
	prog, err := Parse(`var a, b

dump(a)

if cond {
	a = a + b
} else {
	c := 1 + b 
}
println(a)
`)
	if err != nil {
		t.Fatal("prog parse error", err)
	}
	var traceToCall_via_if bool
	prog.Ref("a").ForEach(func(value *Value) {
		if _, ok := value.node.(*ssa.Phi); ok {
			value.GetUsers().ForEach(func(value *Value) {
				if _, ok := value.node.(*ssa.Call); ok {
					traceToCall_via_if = true
					log.Infof("a's deep uses include: %v", value.String())
				}
			})
		}
	})
	if !traceToCall_via_if {
		t.Error("trace failed: var cannot trace to call actual arguments")
	}
}

func MustParse(code string, t *testing.T) *Program {
	prog, err := Parse(code)
	if err != nil {
		t.Fatal("prog parse error", err)
	}
	return prog
}

func TestYaklangBasic_Foreach(t *testing.T) {
	t.Run("for each with chan", func(t *testing.T) {
		test := assert.New(t)
		prog := MustParse(`
		ch = make(chan int)

		for i in ch { 
			_ = i 
		}
		`, t)
		prog.Show()

		vs := prog.Ref("i")
		test.Equal(1, len(vs))

		v := vs[0]
		test.NotNil(v)

		kind := v.GetTypeKind()
		log.Info("type kind", kind)
		test.Equal(kind, ssa.Number)
	})
}

func TestYaklangParameter(t *testing.T) {
	t.Run("test parameter used", func(t *testing.T) {
		test := assert.New(t)
		prog := MustParse(`
		f = (a) => {
			return a
		}
		`, t)
		as := prog.Ref("a").ShowWithSource()
		test.Equal(1, len(as))
		test.Equal("a", *as[0].GetRange().SourceCode)
	})

	t.Run("test parameter not used", func(t *testing.T) {
		test := assert.New(t)
		prog := MustParse(`
		f = (a) => {
			return 1
		}
		`, t)
		as := prog.Ref("a").ShowWithSource()
		test.Equal(1, len(as))
		test.NotNil(as[0].GetRange())
		test.Equal("a", *as[0].GetRange().SourceCode)
	})

	t.Run("test free value used", func(t *testing.T) {
		test := assert.New(t)
		prog := MustParse(`
		f = () => {
			return a
		}
		`, t)
		as := prog.Ref("a").ShowWithSource()
		test.Equal(1, len(as))
		test.Equal("a", *as[0].GetRange().SourceCode)
	})
}
func TestExternLibInClosure(t *testing.T) {
	test := assert.New(t)
	prog, err := Parse(`
	a = () => {
		lib.method()
	}
	`,
		WithExternLib("lib", map[string]any{
			"method": func() {},
		}),
	)
	test.Nil(err)
	libVariables := prog.Ref("lib").ShowWithSource()
	// TODO: handler this
	// test.Equal(1, len(libVariables))
	test.NotEqual(0, len(libVariables))
	libVariable := libVariables[0]

	test.False(libVariable.IsParameter())
	test.True(libVariable.IsExternLib())
}

func check(code string, want []string, t *testing.T) *Program {
	test := assert.New(t)

	prog, err := Parse(code)

	test.Nil(err)

	prog.Show()

	println := prog.Ref("println").ShowWithSource()
	// test.Equal(1, len(println), "println should only 1")
	got := lo.Map(
		println.GetUsers().ShowWithSource().Flat(func(v *Value) Values {
			return Values{v.GetOperand(1)}
		}),
		func(v *Value, _ int) string {
			return v.String()
		},
	)
	// sort.Strings(got)
	log.Info("got :", got)
	// sort.Strings(want)
	log.Info("want :", want)
	test.Equal(want, got)

	return prog
}

func TestYaklangBasic_Variable_InBlock(t *testing.T) {
	t.Run("test simple assign", func(t *testing.T) {
		check(`
		a = 1
		println(a)
		a = 2
		println(a)
		`, []string{
			"1",
			"2",
		}, t)
	})

	t.Run("test sub-scope capture parent scope in basic block", func(t *testing.T) {
		check(`
		a = 1
		println(a)
		{
			a = 2
			println(a)
		}
		println(a)
		`, []string{
			"1",
			"2",
			"2",
		}, t)
	})

	t.Run("test sub-scope local variable in basic block", func(t *testing.T) {
		check(`
		a = 1
		println(a)
		{
			a := 2
			println(a)
		}
		println(a)
		`, []string{
			"1",
			"2",
			"1",
		}, t)
	})

	t.Run("test sub-scope and return", func(t *testing.T) {
		check(`
		a = 1
		println(a) // 1
		{
			a  = 2 
			println(a) // 2
			return 
		}
		println(a) // 1
		`,
			[]string{
				"1", "2",
			}, t)
	})

	t.Run("test variable, but not cover", func(t *testing.T) {
		check(`
		{
			a = 2
			println(a) // 2
		}
		println(a) // undefine-a
		`, []string{
			"2",
			"Undefined-a",
		}, t)
	})

	t.Run("test ++ expression", func(t *testing.T) {
		check(`
		a = 1
		{
			a ++
			println(a) // 2
		}
		`,
			[]string{
				"2",
			},
			t)
	})
}

func TestYaklangBasic_Variable_InIf(t *testing.T) {
	t.Run("test simple if", func(t *testing.T) {
		check(`
		a = 1
		println(a)
		if c {
			a = 2
			println(a)
		}
		println(a)
		`, []string{
			"1",
			"2",
			"phi(a)[2,1]",
		}, t)
	})
	t.Run("test simple if with local vairable", func(t *testing.T) {
		check(`
		a = 1
		println(a)
		if c {
			a := 2
			println(a)
		}
		println(a) // 1
		`, []string{
			"1",
			"2",
			"1",
		}, t)
	})

	t.Run("test multiple phi if", func(t *testing.T) {
		prog := check(`
		a = 1
		if c {
			a = 2
		}
		println(a)
		println(a)
		println(a)
		`, []string{
			"phi(a)[2,1]",
			"phi(a)[2,1]",
			"phi(a)[2,1]",
		}, t)

		phi := prog.Ref("a").Filter(func(v *Value) bool {
			return v.IsPhi()
		}).Show()

		if len(phi) != 1 {
			t.Fatalf("got %v, want %v", phi, 1)
		}
	})

	t.Run("test simple if else", func(t *testing.T) {
		check(`
		a = 1
		println(a)
		if c {
			a = 2
			println(a)
		} else {
			a = 3
			println(a)
		}
		println(a)
		`, []string{
			"1",
			"2",
			"3",
			"phi(a)[2,3]",
		}, t)
	})

	t.Run("test simple if else with origin branch", func(t *testing.T) {
		check(`
		a = 1
		println(a)
		if c {
			// a = 1
		} else {
			a = 3
			println(a)
		}
		println(a) // phi(a)[1, 3]
		`, []string{
			"1",
			"3",
			"phi(a)[1,3]",
		}, t)
	})

	t.Run("test if-elseif", func(t *testing.T) {
		check(`
		a = 1
		println(a)
		if c {
			a = 2
			println(a)
		}else if c == 2{
			a = 3 
			println(a)
		}
		println(a)
		`,
			[]string{
				"1",
				"2",
				"3",
				"phi(a)[2,3,1]",
			}, t)
	})

	t.Run("test with return, no DoneBlock", func(t *testing.T) {
		check(`
		a = 1
		println(a) // 1
		if c {
			if b {
				a = 2
				println(a) // 2
				return 
			}else {
				a = 3
				println(a) // 3
				return 
			}
			println(a) // unreachable // phi[2, 3]
		}
		println(a) // 1
		`, []string{
			"1",
			"2",
			"3",
			"1",
		}, t)
	})
}

func TestYaklangBasic_Variable_Loop(t *testing.T) {
	t.Run("simple loop not change", func(t *testing.T) {
		check(`
		a = 1
		for i=0; i < 10 ; i++ {
			println(a) // 1
		}
		println(a) //1 
		`,
			[]string{
				"1",
				"1",
			},
			t)
	})

	t.Run("simple loop only condition", func(t *testing.T) {
		check(`
		i = 1
		for i < 10 { 
			println(i) // phi
			i = 2 
			println(i) // 2
		}
		println(i) // phi
		`, []string{
			"phi(i)[1,2]",
			"2",
			"phi(i)[1,2]",
		}, t)
	})

	t.Run("simple loop", func(t *testing.T) {
		check(`
		i=0
		for i=0; i<10; i++ {
			println(i) // phi[0, i+1]
		}
		println(i)
		`,
			[]string{
				"phi(i)[0,add(i, 1)]",
				"phi(i)[0,add(i, 1)]",
			}, t)
	})

	t.Run("loop with spin, signal phi", func(t *testing.T) {
		check(`
		a = 1
		for i := 0; i < 10; i ++ { // i=0; i=phi[0,1]; i=0+1=1
			println(a) // phi[0, $+1]
			a = 0
			println(a) // 0 
		}
		println(a)  // phi[0, 1]
		`,
			[]string{
				"phi(a)[1,0]",
				"0",
				"phi(a)[1,0]",
			},
			t)
	})

	t.Run("loop with spin, double phi", func(t *testing.T) {
		check(`
		a = 1
		for i := 0; i < 10; i ++ {
			a += 1
			println(a) // add(phi, 1)
		}
		println(a)  // phi[1, add(phi, 1)]
		`,
			[]string{
				"add(phi(a)[1,add(a, 1)], 1)",
				"phi(a)[1,add(a, 1)]",
			},
			t)
	})
}
