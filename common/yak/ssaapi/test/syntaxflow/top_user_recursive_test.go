package syntaxflow

import (
	"testing"

	"github.com/yaklang/yaklang/common/yak/ssaapi/test/ssatest"
)

func TestSideEffectAndMask_Recursive(t *testing.T) {
	t.Run("mask", func(t *testing.T) {
		code := `
var a = 3
b = () => {
	a ++
}
if c {
	b()
}
e = a
		`
		ssatest.CheckSyntaxFlow(t, code,
			`e #-> * as $target`,
			map[string][]string{
				"target": {"1", "3", "Undefined-c"},
			})
	})
}

func TestObject_Recursive(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		code := `
		a = {}
		b = () => {
			a.b = 333333
		}
		b()
		c = a.b;
		`
		ssatest.CheckSyntaxFlowContain(t, code,
			`c #-> * as $target`,
			map[string][]string{
				"target": {"333333"},
			},
		)
	})

	t.Run("object mark self", func(t *testing.T) {
		code := `
			a = {}; 
			a.b = 1; 
			a.c = 3; 
			d = a.c + a.b
		`
		ssatest.CheckSyntaxFlow(t, code,
			`d #-> * as $target`,
			map[string][]string{
				"target": {"1", "3"},
			},
		)
	})

	t.Run("object mark self, bottom user", func(t *testing.T) {
		// func Test_Yaklang_BottomUser(t *testing.T) {
		code := `
		f = () =>{
			a = 11
			return a
		}
		f2 = (i) => {
			println(i)
		}
		t = f()
		f2(t)
		`
		ssatest.CheckSyntaxFlow(t, code,
			`a --> * as $target`,
			map[string][]string{
				"target": {
					"FreeValue-println(Parameter-i)",
				},
			},
		)
	})
}

func TestParameter_TopDef_REcursive(t *testing.T) {
	t.Run("parameter top def recursive  1", func(t *testing.T) {
		ssatest.CheckSyntaxFlow(t, `
		f1 = (a1) => {
			return a1
		}
		f2 = (a2)  => {
			target = f1(a2)
			f2(target)
		}
		`, `
		target #-> * as $target
		`, map[string][]string{
			// step 1 :找到f1(a2),f1(a2)是个Call，因此PushCall
			// step 2 :找到返回值a1
			// step 3 :找到调用a1的点，即f1(a2)，进行PopCall并对其形参a2进行topdef
			// step 4 :找到调用a2的点，还是f1(a2)，但是此时CallStack是空的，启用NegativeCallStack，进行topdef
			// step 5 :NegativeCallStack判断查找了两次参数a2，停止递归
			"target": {"FreeValue-f1(Parameter-a2)"},
		})
	})
	t.Run("parameter top def recursive 2 ", func(t *testing.T) {
		ssatest.CheckSyntaxFlow(t, `
		f2 = (a)  => {
			f2(a)
		}
		`, `
		a?{opcode: param} #-> * as $target`, map[string][]string{
			"target": {"Parameter-a"},
		})
	})
}
