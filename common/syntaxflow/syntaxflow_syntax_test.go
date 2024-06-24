package syntaxflow

import (
	"testing"

	"github.com/yaklang/yaklang/common/syntaxflow/sfvm"
)

func TestSyntaxInOne(t *testing.T) {
	for _, i := range []string{
		"aa",
		// "$",
		"*",
		".*",
		"a.*",
		"exec",    // Ref
		".member", // Field
		".*exec*",
		"*exec",
		"exe*c",
		"/$reexc/",
		"./$reexc/",
		"a[1]",
		"a.b",
		"c.d",
		"$c.d",
		"a[1]",
		"b?{>1}",
		"b?{!have: /abc/}",
		"f()",
		"f() as $a",
		"f(*,)",
		"f(* as $arg1,)",
		"f(,*)",
		"f(*,*)",
		"f(*,*,)",
		"f(,*,*,)",
		"f(,*,*)",
		"runtime.exec",
		"$runtime.exec",
		"runtime.exec()",
		"runtime.exec(,,)",
		"/(?i)runtime/.exec(,,,#>*exec)",
		"exec as $rough",
		"/(?i)runtime/.exec(,,,#>*exec) as $a",
		"/(?i)runtime/.exec(,,,#>*exec as $b) as $a",
		"/(?i)runtime/.exec(,,,#>*exec,,) as $a",
		"/(?i)runtime/.exec(,,,#>*exec as $ccc,,) as $a",
		"/(?i)runtime/.exec(,,,#>*exec,,,) as $a",
		"$a->b",
		"$a#>b",
		"$a-->b",
		"$a#->b",
		"a->b",
		"a#>b",
		"a-->b",
		"a#->b",
		"a->b->c",
		"a#>b#>c",
		"a-->b.exec()-->c",
		"a#->b.exec()#->c",
		"a-{}->b",
		"a#{}->b",
		"system(#>get)",
		"system(#>*)",
		"system(#>a*get)",
		"$system(#>a*get)",
		"a*",
		"*a",
		"*",
		"system(#{}*a)",
		"system(#{}*)",
		"a-{depth:1}->b",
		"a#{depth:1}->b",
		"$a-{depth:1, " + "\nkey:value}->b",
		"a-{depth:1, " + "\nkey:value}->b",
		"a#{depth:1, " + "\nkey:value}->b",
		"a-{depth:1, " + "\nkey:value,}->b",
		"a#{depth:1, " + "\nkey:value,}->b",
		`*.mem as $a
		$a.exec() as $exec
		`,
		`
		a -{
			until:` + " `*->a`" + `
		}-> *`,
		"a -<abc>- b",
		"a -<abc{depth: 1}>- b",
		"a -<abc(depth: 1, asdf: `a -{}-> *`)>- b",
		"a ->",
		"a #>",
		"a #->",
		"a #{}->",
		"a #{depth: 1}->",
		"a -->",
		"a -{depth: 1}->",
		"a -{}->",
		"a -> as $b",
		"a #> as $b",
		"a #-> as $b",
		"a --> as $b",
		`a;a;a;a;;;;;`,
		"check $abc",
		"check $abc;;;;;",
		"check $abc then Finished else BAD",
		"check $abc then 'asdfasdfasd' else \"abc\"",
		"check $abc else BAD",
		"check $abc else 'Bad Boy~~~ asdfasd adsf asdfpouasdf jpoasdf '",
		"check $abc then 'Finished Hello?'",
		"check $abc then Finished else BAD;;;; desc{a: b};; desc{title: SprintChecking}",
		"desc(a: b, c: eee, e,e,e,e)",
		`desc(title: "你好,世界，你可以在这里输入任何内容")`,
		"desc(a: b, c: eee, e,e,e,e);;;;check $abc then GOOD else BAD",
		`"abc";'abc';` + "'abcasdfasdf'",
		`"a\"bc";'a\'bc';` + "'abcasdfasdfaaa'",
		`"\\";'\\'`,
		`"\";'\'`,
		"a?{opcode: const} as $a",
		"$a?{opcode: const} as $a",
		"exec(* #-> ?{opcode: call,phi})",
		"exec(* #-> ?{!(opcode: call,phi,)})",
		"exec(* #-> ?{!(opcode: call,phi,)})",
		"exec(* #-> ?{(!have: 'example')})",
		"exec(* #-> ?{(have: 'example')})",
		"exec(* #-> ?{(not any: 'example')})",
		"exec(* #-> ?{(any: 'example',abc,)})",
		"exec(* #-> ?{((any: 'example',abc,) && (opcode: call,phi))})",
		"exec(* #-> ?{(any: 'example',abc,) && (opcode: call,phi)})",
		"exec(* #-> ?{((any: 'example',abc,) && (opcode: phi))})",
	} {
		vm := sfvm.NewSyntaxFlowVirtualMachine().Debug(true)
		_, err := vm.Compile(i)
		if err != nil {
			t.Fatalf("syntax failed: %#v, reason: %v", i, err)
		}
	}
}
