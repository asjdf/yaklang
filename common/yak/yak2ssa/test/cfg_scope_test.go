package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaklang/yaklang/common/yak/ssaapi"
)

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

	t.Run("simple test", func(t *testing.T) {
		check(`
		println(a)
		`, []string{"Undefined-a"}, t)
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
		println(a) // 1
		{
			a := 2
			println(a) // 2
		}
		println(a) // 1
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
		println(a) // unreachable
		`,
			[]string{
				"1", "2",
			}, t)
	})

	t.Run("undefine variable in sub-scope", func(t *testing.T) {
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

		phi := prog.Ref("a").Filter(func(v *ssaapi.Value) bool {
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

	t.Run("in if sub-scope", func(t *testing.T) {
		check(`
		if c {
			a = 2
		}
		println(a)
		`, []string{"Undefined-a"}, t)
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
func TestYaklangParameter(t *testing.T) {
	check := func(code string, t *testing.T) {
		test := assert.New(t)
		prog, err := ssaapi.Parse(`
		f = (a) => {
			return a
		}
		`)
		test.Nil(err)
		as := prog.Ref("a").ShowWithSource()
		test.Equal(1, len(as))
		test.Equal("a", *as[0].GetRange().SourceCode)

	}
	t.Run("test parameter used", func(t *testing.T) {
		check(
			`
		f = (a) => {
			return a
		}
		`, t)
	})

	t.Run("test parameter not used", func(t *testing.T) {
		check(`
		f = (a) => {
			return 1
		}
		`, t)
	})

	t.Run("test free value used", func(t *testing.T) {
		check(`
		f = () => {
			return a
		}
		`, t)
	})
}
