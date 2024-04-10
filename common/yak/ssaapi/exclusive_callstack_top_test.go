package ssaapi

import (
	"testing"
)

func Test_CallStack_Normal_Parameter(t *testing.T) {
	t.Run("test level1", func(t *testing.T) {
		Check(t, `
		f = (i) => {
			return i
		}
		a = f(333333)
		`,
			CheckTopDef_Equal("a", []string{"333333"}),
		)
	})

	t.Run("test level1 object", func(t *testing.T) {
		Check(t, `
		f = (i) => {
			return {
				"i": i,
			}
		}
		obj = f(333333)
		a = obj.i
		`,
			CheckTopDef_Equal("a", []string{"333333"}),
		)
	})

	t.Run("test level1 class", func(t *testing.T) {
		Check(t, `
		f = (i) => {
			this =  {
				"i": i,
				// this will be set a class-blue-print
				"set": (i)=>{this.i = i}, 
			}
			return this
		}
		obj = f(333333)
		a = obj.i
		`,
			CheckTopDef_Equal("a", []string{"333333"}),
		)
	})

	t.Run("test level1 php class", func(t *testing.T) {
		Check(t, `
		<?php
		class A {
			public $i;
			public function __construct($i) {
				$this->i = $i;
			}
		}
		$obj  = new A(333333);
		$a = $obj->i;
		`,
			CheckTopDef_Equal("$a", []string{"333333"}),
			PHP,
		)
	})

	t.Run("test level2", func(t *testing.T) {
		Check(t, `
		f = (i) => {
			return (j) => {
				return j + i
			} 
		}
		f1 = f(333333)
		a = f1(444444)
		`,
			CheckTopDef_Equal("a", []string{"333333", "444444"}),
		)
	})

	t.Run("test level2 object", func(t *testing.T) {
		Check(t, `
		f = (i) => {
			return (j) => {
				return {
					"i": j + i,
				}
			} 
		}
		f1 = f(333333)
		obj = f1(444444)
		a = obj.i
		`,
			CheckTopDef_Equal("a", []string{"333333", "444444"}),
		)
	})

	t.Run("test level3", func(t *testing.T) {
		Check(t, `
		f = (i) => {
			return (j) => {
				return (k) => {
					return k + j + i
				}
			} 
		}
		f1 = f(333333)
		f2 = f1(444444)
		a = f2(444444)
		`, CheckTopDef_Equal("a", []string{"333333", "444444", "444444"}),
		)
	})

	t.Run("test level3 object", func(t *testing.T) {
		Check(t, `
		f = (i) => {
			return (j) => {
				return (k) => {
					return {
						"i": k + j + i
					}
				}
			} 
		}
		f1 = f(333333)
		f2 = f1(444444)
		obj = f2(555555)
		a = obj.i
		`, CheckTopDef_Equal("a", []string{"333333", "444444", "555555"}),
		)
	})
}

func Test_CallStack_FreeValue_WithDefault(t *testing.T) {
	t.Run("test level1", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			return i
		}
		a = f()
		`, CheckTopDef_Equal("a", []string{"333333"}),
		)
	})

	t.Run("test level1, object", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			return {
				"i": i,
			}
		}
		obj = f()
		a = obj.i
		`,
			CheckTopDef_Equal("a", []string{"333333"}),
		)
	})

	t.Run("test level1 class", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			this =  {
				"i": i,
				// this will be set a class-blue-print
				"set": (i)=>{this.i = i},
			}
			return this
		}
		obj = f()
		a = obj.i
		`,
			CheckTopDef_Equal("a", []string{"333333"}),
		)
	})

	t.Run("test level2", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			j = 444444
			return () => {
				return i + j
			}
		}
		f1 = f()
		a = f1()
		`, CheckTopDef_Equal("a", []string{"333333", "444444"}),
		)
	})

	t.Run("test level2 object", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			j = 444444
			return () => {
				return {
					"i": j + i, 
				}
			}
		}
		f1 = f()
		obj = f1()
		a = obj.i
		println(a)
		`, CheckTopDef_Equal("a", []string{"333333", "444444"}),
		)
	})

	t.Run("test level3", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			j = 444444
			return () => {
				k = 555555
				return () => {
					return i + j + k
				}
			}
		}
		f1 = f()
		f2 = f1() 
		a = f2()
		`, CheckTopDef_Equal("a", []string{"333333", "444444", "555555"}),
		)
	})

	t.Run("test level3 object", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			j = 444444
			return () => {
				k = 555555
				return () => {
					return {
						"i": i + j + k
					}
				}
			}
		}
		f1 = f()
		f2 = f1() 
		obj = f2()
		a = obj.i
		`, CheckTopDef_Equal("a", []string{"333333", "444444", "555555"}),
		)
	})
}

func Test_CallStack_FreeValue_WithoutDefault(t *testing.T) {

	t.Run("test level1", func(t *testing.T) {
		Check(t, `
		f = () => {
			return i
		}
		i = 333333
		a = f()
		`, CheckTopDef_Equal("a", []string{"333333"}),
		)
	})

	t.Run("test level1 object", func(t *testing.T) {
		Check(t, `
		f = () => {
			return {
				"i": i,
			}
		}
		i = 333333
		obj = f()
		a = obj.i
		`, CheckTopDef_Equal("a", []string{"333333"}))
	})

	t.Run("test level2", func(t *testing.T) {
		Check(t, `
		f = () => {
			return () => {
				return i + j
			}
		}
		i = 333333
		j = 444444
		f1 = f()
		a = f1()
		`, CheckTopDef_Equal("a", []string{"333333", "444444"}))
	})

	t.Run("test level2 object", func(t *testing.T) {
		Check(t, `
		f = () => {
			return () => {
				return {
					"i": i + j,
				}
			}
		}
		i = 333333
		j = 444444
		f1 = f()
		obj = f1()
		a = obj.i
		`, CheckTopDef_Equal("a", []string{"333333", "444444"}))
	})

	t.Run("test level3", func(t *testing.T) {
		Check(t, `
		f = () => {
			return () => {
				return () => {
					return i + j + k
				}
			}
		}
		i = 333333
		j = 444444
		k = 555555
		f1 = f()
		f2 = f1()
		a = f2()
		`, CheckTopDef_Equal("a", []string{"333333", "444444", "555555"}))
	})

	t.Run("test level3 object", func(t *testing.T) {
		Check(t, `
		f = () => {
			return () => {
				return () => {
					return {
						"i": i + j + k
					}
				}
			}
		}
		i = 333333
		j = 444444
		k = 555555
		f1 = f()
		f2 = f1()
		obj = f2()
		a = obj.i
		`, CheckTopDef_Equal("a", []string{"333333", "444444", "555555"}))
	})

}

func Test_CallStack_FreeVale_Parameter(t *testing.T) {
	// TODO : implement this testcase,
	// function and call-site should be alignment.
	// t.Run("test level2 with parameter", func(t *testing.T) {
	// 	Check(t, `
	// 	f = (i) => {
	// 		return () => {
	// 			return i + j
	// 		}
	// 	}
	// 	j = 444444
	// 	f1 = f(333333)
	// 	a = f1()
	// 	`, CheckTopDef_Equal("a", []string{"333333", "444444"}))
	// })
}

func Test_CallStack_Normal_SideEffect(t *testing.T) {
	t.Run("test level1", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			i = 444444
		}
		f()
		a = i
		`, CheckTopDef_Equal("a", []string{"444444"}),
		)
	})

	t.Run("test level1, object member", func(t *testing.T) {
		code := `
		a = {}
		b = () => {
			a.b = 333333
		}
		b()
		c = a.b;
		`

		Check(t, code,
			CheckTopDef_Contain("c", []string{"Function-b(", "333333"}),
		)
	})

	t.Run("test level2", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			j = 444444
			return () => {
				i = j
			}
		}
		f1 = f()
		f1()
		a = i
		`, CheckTopDef_Equal("a", []string{"444444"}),
		)
	})

	t.Run("test level2 object", func(t *testing.T) {
		Check(t, `
		obj = {}
		i = 333333
		f = () => {
			j = 444444
			return () => {
				obj.i = j + i
			}
		}
		f1 = f()
		f1()
		a = obj.i
		`, CheckTopDef_Equal("a", []string{"333333", "444444"}),
		)
	})

	t.Run("test level3", func(t *testing.T) {
		Check(t, `
		i = 333333
		f = () => {
			j = 444444
			return () => {
				k = 555555
				return () => {
					i = k
				}
			}
		}
		f1 = f()
		f2 = f1()
		f2()
		a = i
		`, CheckTopDef_Equal("a", []string{"555555"}),
		)
	})

	t.Run("test level3 object", func(t *testing.T) {
		Check(t, `
		obj = {}
		i = 333333
		f = () => {
			j = 444444
			return () => {
				k = 555555
				return () => {
					obj.i = i + j + k
				}
			}
		}
		f1 = f()
		f2 = f1()
		f2()
		a = obj.i
		`, CheckTopDef_Equal("a", []string{"333333", "444444", "555555"}),
		)
	})

}
