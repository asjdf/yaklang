package java

import "testing"

func TestSimpleSearchType(t *testing.T) {

	code := `
	class C extends D {}

	class B  extends C {
		void methodB(int i) {
			println(i);
		}
	}

	class A {
		int a; 
		void main() {
			int B = 1;
			B b1  = new B();
			b1.methodB(1);

			B b2  = new B();
			b2.methodB(2);
		}
	}
	`

	// just append Class Instance to cache and database, we can pass this test.
	t.Run("get class instance and variable", func(t *testing.T) {
		test(t, &TestCase{
			Code:    code,
			SF:      `B as $target`,
			Contain: false,
			Expect: map[string][]string{
				"target": {
					"1",
					"make(B)",
					"make(B)",
				},
			},
		})
	})

	t.Run("get extern class instance", func(t *testing.T) {
		test(t, &TestCase{
			Code:    code,
			SF:      `C as $target`,
			Contain: false,
			Expect: map[string][]string{
				"target": {
					"make(B)",
					"make(B)",
				},
			},
		})
	})

	t.Run("get multiple level extern class instance", func(t *testing.T) {
		test(t, &TestCase{
			Code:    code,
			SF:      `D as $target`,
			Contain: false,
			Expect: map[string][]string{
				"target": {
					"make(B)",
					"make(B)",
				},
			},
		})
	})

	t.Run("get class method", func(t *testing.T) {
		test(t, &TestCase{
			Code:    code,
			SF:      `B.methodB as $target`,
			Contain: false,
			Expect: map[string][]string{
				"target": {
					"Undefined-b1.methodB(valid)",
					"Undefined-b2.methodB(valid)",
				},
			},
		})
	})

	t.Run("get class instance method", func(t *testing.T) {
		test(t, &TestCase{
			Code:    code,
			SF:      `b1.methodB() as $target`,
			Contain: false,
			Expect: map[string][]string{
				"target": {
					"Undefined-b1.methodB(valid)(make(B),1)",
				},
			},
		})
	})

	t.Run("get class method call", func(t *testing.T) {
		test(t, &TestCase{
			Code:    code,
			SF:      `B.methodB() as $target`,
			Contain: false,
			Expect: map[string][]string{
				"target": {
					"Undefined-b1.methodB(valid)(make(B),1)",
					"Undefined-b2.methodB(valid)(make(B),2)",
				},
			},
		})
	})

	t.Run("method function should has called", func(t *testing.T) {
		test(t, &TestCase{
			Code:    code,
			SF:      `println(* #-> * as $target)`,
			Contain: false,
			Expect: map[string][]string{
				"target": {"1", "2"},
			},
		})
	})

}
