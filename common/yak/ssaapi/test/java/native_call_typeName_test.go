package java

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaklang/yaklang/common/consts"
	"github.com/yaklang/yaklang/common/utils/filesys"
	"github.com/yaklang/yaklang/common/yak/ssaapi"
	"github.com/yaklang/yaklang/common/yak/ssaapi/test/ssatest"
)

func TestNativeCallTypeName(t *testing.T) {
	ssatest.Check(t, XXE_Code, func(prog *ssaapi.Program) error {
		typeName := prog.SyntaxFlowChain(`documentBuilder<typeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "DocumentBuilder")
		typeName = prog.SyntaxFlowChain(`documentBuilder<fullTypeName> as $id;`)[0]
		assert.Contains(t, typeName.Show().String(), "javax.xml.parsers.DocumentBuilder")
		return nil
	}, ssaapi.WithLanguage(ssaapi.JAVA))
}

func TestNativeCallTypeNameWithSCAVersion(t *testing.T) {
	vf := filesys.NewVirtualFs()
	
	vf.AddFile("FastJSONDemoController.java",
		`package com.example.demo.controller.fastjsondemo;

import com.alibaba.fastjson.JSON;
import org.apache.ibatis.annotations.Param;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/fastjson")
public class FastJSONDemoController {
    @GetMapping("/fromId")
    public ResponseEntity<Object> loadFromParam(@RequestParam(name = "id") String id) {
        // This is a FASTJSON Vuln typically.
        Object anyJSON = JSON.parse(id);
     
        return ResponseEntity.ok(anyJSON);
    }
}
`)
	vf.AddFile("pom.xml",
		`<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>3.2.7</version>
        <relativePath/> <!-- lookup parent from repository -->
    </parent>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>0.0.1-SNAPSHOT</version>
    <name>demo</name>
    <description>Demo project for Spring Boot</description>
    <url/>
    <properties>
        <java.version>8</java.version>
    </properties>
    <dependencies>
        <dependency>
            <groupId>com.alibaba</groupId>
            <artifactId>fastjson</artifactId>
            <version>1.2.24</version>
        </dependency>
    </dependencies>
</project>
`)

	ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
		prog := progs[0]
		obj := prog.SyntaxFlowChain(`JSON<fullTypeName>?{have: 'alibaba.fastjson'} as $obj`).Show(false)
		assert.NotNil(t, obj)

		obj = prog.SyntaxFlowChain(`parse*?{<getObject><fullTypeName>?{have: 'alibaba.fastjson'} } as $obj`).Show(false)
		assert.NotNil(t, obj)

		obj = prog.SyntaxFlowChain(`ok()?{<getCaller><getObject><fullTypeName>?{have: 'org.springframework.'} } as $obj`).Show(true)
		assert.NotNil(t, obj)

		typeName := prog.SyntaxFlowChain(`anyJSON<typeName>?{have:'JSON'} as $id;`).Show()
		assert.Contains(t, typeName.String(), "JSON")
		typeName = prog.SyntaxFlowChain(`anyJSON<fullTypeName>?{have:'JSON'} as $id`)
		assert.Contains(t, typeName.String(), "com.alibaba.fastjson.JSON:1.2.24")
		return nil
	}, ssaapi.WithLanguage(consts.JAVA))
}

func TestLocalVariableDeclareTypeName(t *testing.T){
			vf := filesys.NewVirtualFs()
			vf.AddFile("A.java",
				`package com.org.A;
				class A{
					};

		    `)
		vf.AddFile("B.java",
			`package com.example.B;
			import com.org.A.A;
			class B{
				public static void main(String[] args){
					A res1 = "aaa";  
					A res2 = 1;  				
					var res3 = A;  
					var res4 ="a";     
					A res5 = Dog(); 
					A test1 ,test2 = Dog();
				}
			};	
	`)
	ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
		prog := progs[0]
		prog.Show()

		obj := prog.SyntaxFlowChain(`res1<typeName>?{have: 'string' || have: 'A'} as $obj`)
		assert.Equal(t,3,obj.Len())
		obj = prog.SyntaxFlowChain(`res1<fullTypeName>?{have: 'string' || have: 'A'} as $obj`)
		assert.Equal(t,2,obj.Len())
		
		obj = prog.SyntaxFlowChain(`res2<typeName>?{have:'number' || have: 'A'}  as $obj`)
		assert.Equal(t,3,obj.Len())
		obj = prog.SyntaxFlowChain(`res2<fullTypeName>?{have:'number' || have: 'A'}as $obj`)
		assert.Equal(t,2,obj.Len())

		obj = prog.SyntaxFlowChain(`res3<typeName>?{have:'A'} as $obj`)
		assert.Equal(t,2,obj.Len())
		obj = prog.SyntaxFlowChain(`res3<fullTypeName>?{have:'com.org.A.A'} as $obj`)
		assert.Equal(t,1,obj.Len())

		obj = prog.SyntaxFlowChain(`res4<typeName>?{have:'string'} as $obj`)
		assert.Equal(t,1,obj.Len())
		obj = prog.SyntaxFlowChain(`res4<fullTypeName>?{have: 'string'}as $obj`)
		assert.Equal(t,1,obj.Len())

		obj = prog.SyntaxFlowChain(`res5<typeName>?{have:'A'}as $obj`).Show()
		assert.Equal(t,2,obj.Len())
		obj = prog.SyntaxFlowChain(`res5<fullTypeName>?{have:'com.org.A.A'} as $obj`)
		assert.Equal(t,1,obj.Len())

		obj = prog.SyntaxFlowChain(`test2<typeName>?{have:'A'} as $obj`)
		assert.Equal(t,2,obj.Len())
		obj = prog.SyntaxFlowChain(`test2<fullTypeName>?{have:'com.org.A.A'} as $obj`)
		assert.Equal(t,1,obj.Len())
			return nil
		}, ssaapi.WithLanguage(consts.JAVA))
		
}

func TestMemberCallTypeName(t *testing.T){
	vf := filesys.NewVirtualFs()
		vf.AddFile("Dog.java",`package com.org.Dog; class Dog{};`)
		vf.AddFile("A.java",
			`package com.org.A;
			 import com.org.Dog.Dog;
				class A{
					public int existMethod(){return 666;}
					public Dog getDog(){return new Dog();}
					public static Dog staticMethod(){return new Dog();};
					};
		    `)
		vf.AddFile("B.java",
			`package com.example.B;
			import com.org.A.A;
			class B{
				public static void main(String[] args){
					A object = new A();
					var res1 = object.noExistMethod();  // fulltypeName 应该和object一样
					var res2 = object.existMethod();  // fulltypeName 应该为number
					var res3 = object.getDog();  // fulltypeName 应为com.org.Dog.Dog
					var res4 = object.Runtime().exec();  // fulltypeName 应该和object一样
					var res5 = A.staticMethod();  // fulltypeName 应该找到Dog
					var res6 = A.noExistMethod();  // fulltypeName 应该找到A	
				}
			};	
	`)
		
		ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
			prog := progs[0]
			prog.Show()

			obj := prog.SyntaxFlowChain(`res1<typeName> as $obj`)
			assert.Equal(t,2,obj.Len())
			obj = prog.SyntaxFlowChain(`res1<fullTypeName>?{have: 'com.org.A.A'} as $obj`)
			assert.Equal(t,1,obj.Len())
			
			obj = prog.SyntaxFlowChain(`res2<typeName>?{have:'number'} as $obj`)
			assert.Equal(t,1,obj.Len())
			obj = prog.SyntaxFlowChain(`res2<fullTypeName>?{have:'number'} as $obj`)
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain(`res3<typeName>?{have:'Dog'} as $obj`)
			assert.Equal(t,2,obj.Len())
			obj = prog.SyntaxFlowChain(`res3<fullTypeName>?{have:'com.org.Dog.Dog'} as $obj`)
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain(`res4<typeName>?{have:'A'} as $obj`)
			assert.Equal(t,2,obj.Len())
			obj = prog.SyntaxFlowChain(`res4<fullTypeName>?{have: 'com.org.A.A'}as $obj`)
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain(`res5<typeName>?{have:'Dog'} as $obj`)
			assert.Equal(t,2,obj.Len())
			obj = prog.SyntaxFlowChain(`res5<fullTypeName>?{have: 'com.org.Dog.Dog'}as $obj`)
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain(`res6<typeName>?{have:'A'} as $obj`)
			assert.Equal(t,2,obj.Len())
			obj = prog.SyntaxFlowChain(`res6<fullTypeName>?{have: 'com.org.A.A'}as $obj`)
			assert.Equal(t,1,obj.Len())
			
			return nil
		}, ssaapi.WithLanguage(consts.JAVA))
}

func TestParamTypeName(t *testing.T){
	vf := filesys.NewVirtualFs()
		vf.AddFile("A.java",
			`package com.org.A;
				class A{
					};
		    `)
		vf.AddFile("B.java",
			`package com.example.B;
			import com.org.A.A;
			class B{
				public void hello(int param1,A param2,Dog param3){
					var res1 = param1;
					var res2 = param2;
					var res3 = param3; //Dog()为找不到的类，使用自身作为fullTypeName
					var res4 = a;
				}
			};	
	`)
		ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
			prog := progs[0]
			prog.Show()

			obj := prog.SyntaxFlowChain(`param1<typeName>?{have: 'number'} as $obj`)
			assert.Equal(t,1,obj.Len())
			obj = prog.SyntaxFlowChain(`param1<fullTypeName>?{have: 'number'} as $obj`)
			assert.Equal(t,1,obj.Len())
			
			obj = prog.SyntaxFlowChain(`param2<typeName>?{have:'A'} as $obj`)
			assert.Equal(t,2,obj.Len())
			obj = prog.SyntaxFlowChain(`param2<fullTypeName>?{have:'com.org.A.A'} as $obj`)
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain(`param3<typeName>?{have:'Dog'} as $obj`)
			assert.Equal(t,2,obj.Len())
			obj = prog.SyntaxFlowChain(`param3<fullTypeName>?{have:'Dog'} as $obj`)
			assert.Equal(t,1,obj.Len())

		
			

			return nil
		}, ssaapi.WithLanguage(consts.JAVA))
}


func TestTypeNamePriority(t *testing.T){
	vf := filesys.NewVirtualFs()
		vf.AddFile("A.java",
			`package com.org.A;
				class A{
					};
		    `)
		vf.AddFile("B.java",
			`package com.example.B;
			import com.org.A.A;
			class B{
				public void hello(int param1,A param2){
					Object res1 = (A)param1;
					A res2 = (int)param2;
				}
			};	
	`)
		ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
			prog := progs[0]
			prog.Show()

			obj := prog.SyntaxFlowChain(`res1<typeName>?{have: 'A'} as $obj`)
			assert.Equal(t,2,obj.Len())
			obj = prog.SyntaxFlowChain(`res1<fullTypeName>?{have: 'com.org.A.A'} as $obj`)
			assert.Equal(t,1,obj.Len())
			obj = prog.SyntaxFlowChain(`res2<typeName>?{have:'number'} as $obj`)
			assert.Equal(t,1,obj.Len())
			obj = prog.SyntaxFlowChain(`res2<fullTypeName>?{have:'number'} as $obj`)
			assert.Equal(t,1,obj.Len())


			return nil
		}, ssaapi.WithLanguage(consts.JAVA))
}


func TestTypeNameForImportStar(t *testing.T){
	vf := filesys.NewVirtualFs()
		vf.AddFile("A.java",
			`package com.org.A;
				class A{
					};
		    `)
		vf.AddFile("B.java",
			`package com.example.B;
			import com.org.A.A;
			import com.yak.*;
			class B{
				public void hello(int param1,Dog param2){
					var res1 = param2; 
					Cat res2 = Cat();
					var res3 = new Cat();
				}
			};	
	`)
		ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
			prog := progs[0]
			prog.Show()

			typeName := prog.SyntaxFlowChain(`res1<typeName> as $id;`)
			assert.Equal(t,3,typeName.Show(false).Len())
			typeName = prog.SyntaxFlowChain(`res1<fullTypeName> as $id;`)
			assert.Equal(t,2,typeName.Show(false).Len())
			
			typeName = prog.SyntaxFlowChain(`res2<typeName> as $id;`)
			assert.Equal(t,3,typeName.Show(false).Len())
			typeName = prog.SyntaxFlowChain(`res2<fullTypeName> as $id;`)
			assert.Equal(t,2,typeName.Show(false).Len())

			return nil
		}, ssaapi.WithLanguage(consts.JAVA))
}

func TestFullTypeNameWithParentClass1(t *testing.T){
	vf := filesys.NewVirtualFs()
	
	vf.AddFile("C.java",`
	package com.yak;
	class C{};
	`)
	vf.AddFile("A.java",
		`package com.org.A;
			class A {
				};
		`)
	vf.AddFile("B.java",
		`package com.example.B;
		import com.org.A.A;
		import com.yak.C;
		class B extends A implements C{
			public static void main(String[] args){
				var a = new B();
			}
		};	
`)
		ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
			prog := progs[0]
			prog.Show()

			obj := prog.SyntaxFlowChain("a<typeName> as $obj").Show()
			assert.Equal(t,6,obj.Len())
		
			obj = prog.SyntaxFlowChain("a<typeName>?{have:'com.example.B.B'} as $obj")
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain("a<typeName>?{have:'com.org.A.A'} as $obj")
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain("a<typeName>?{have:'com.org.A.A'} as $obj")
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain("a<fullTypeName>?{have:'com.yak.C'}as $obj")
			assert.Equal(t,1,obj.Len())
		
			obj = prog.SyntaxFlowChain("a<fullTypeName>?{have:'com.example.B.B'} as $obj")
			assert.Equal(t,1,obj.Len())

			obj = prog.SyntaxFlowChain("a<fullTypeName>?{have:'com.org.A.A'} as $obj")
			assert.Equal(t,1,obj.Len())
			
			obj = prog.SyntaxFlowChain("a<typeName>?{have:'com.org.A.A'} as $obj")
			assert.Equal(t,1,obj.Len())
			return nil
		}, ssaapi.WithLanguage(consts.JAVA))
}

func TestFullTypeNameWithParentClass2(t *testing.T){
	vf := filesys.NewVirtualFs()

	vf.AddFile("B.java",
		`package com.example.B;
	
		class B extends A implements C{
			public static void main(String[] args){
				var a = new B();
			}
		};	
`)
		ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
			prog := progs[0]
			prog.Show()

			obj := prog.SyntaxFlowChain("a<typeName> as $obj").Show()			
			assert.Equal(t,6,obj.Len())

			obj = prog.SyntaxFlowChain("a<fullTypeName> as $obj").Show()
			assert.Equal(t,3,obj.Len())
			return nil
		}, ssaapi.WithLanguage(consts.JAVA))
}


func TestFullTypeNameForAnnotation(t *testing.T){
	t.Run("test spring framework annotation", func(t *testing.T) {
		vf := filesys.NewVirtualFs()
	vf.AddFile("Test.java",`
	package com.example;
import org.springframework.web.bind.annotation.*;
import jakarta.servlet.http.HttpServletRequest;
@RequestMapping("/fastjson")
public class FastJSONDemoController {

    public ResponseEntity<Object> loadFromParam(@RequestParam(name = "id") int id) {

    }
}`)

	ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
		prog := progs[0]
		prog.Show()
		obj := prog.SyntaxFlowChain("id.annotation.RequestParam<fullTypeName>?{have:'org.springframework.web.bind.annotation.RequestParam'} as $obj")
		assert.Equal(t,1,obj.Len())

		obj = prog.SyntaxFlowChain("FastJSONDemoController.annotation.RequestMapping<fullTypeName>?{have:'org.springframework.web.bind.annotation.RequestMapping'} as $obj")
		assert.Equal(t,1,obj.Len())

		obj = prog.SyntaxFlowChain("*Param.__ref__<fullTypeName>?{have:number} as $obj")
		assert.Equal(t,1,obj.Len())

		return nil
	}, ssaapi.WithLanguage(consts.JAVA))
	})

	t.Run("test no spring framework anntation type name ", func(t *testing.T) {
		vf := filesys.NewVirtualFs()
	vf.AddFile("Test.java",`
	package com.example;
import org.springframework.web.bind.annotation.*;
import jakarta.servlet.http.HttpServletRequest;

public class FastJSONDemoController {

    public ResponseEntity<Object> loadFromParam(@Hello(name = "id") int id) {

    }
}`)

	ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
		prog := progs[0]
		prog.Show()
		obj := prog.SyntaxFlowChain("id.annotation.Hello<fullTypeName> as $obj").Show()
		assert.Equal(t,2,obj.Len())
		return nil
	}, ssaapi.WithLanguage(consts.JAVA))
	})

	t.Run("test servlet annotation1", func(t *testing.T) {
		vf := filesys.NewVirtualFs()
	vf.AddFile("Test.java",`
	package com.example;

import javax.servlet.annotation.*; 
@WebServlet(value = "/Simple") 
public class Simple extends HttpServlet {

   private static final long serialVersionUID = 1L; 

   protected void doGet(HttpServletRequest request, HttpServletResponse response)  
       { 
   }   
}`)

	ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
		prog := progs[0]
		prog.Show()
		obj := prog.SyntaxFlowChain("Simple.annotation.WebServlet<fullTypeName>?{have:'javax.servlet.annotation.WebServlet'} as $obj")
		assert.Equal(t,1,obj.Len())
		return nil
	}, ssaapi.WithLanguage(consts.JAVA))
	})

	t.Run("test servlet annotation2", func(t *testing.T) {
		vf := filesys.NewVirtualFs()
	vf.AddFile("Test.java",`
	package com.example;

import javax.servlet.annotation.WebServlet; 
@WebServlet(value = "/Simple") 
public class Simple extends HttpServlet {

   private static final long serialVersionUID = 1L; 

   protected void doGet(HttpServletRequest request, HttpServletResponse response)  
       { 
   }   
}`)

	ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
		prog := progs[0]
		prog.Show()
		obj := prog.SyntaxFlowChain("Simple.annotation.WebServlet<fullTypeName>?{have:'javax.servlet.annotation.WebServlet'} as $obj")
		assert.Equal(t,1,obj.Len())
		return nil
	}, ssaapi.WithLanguage(consts.JAVA))
	})
}

func TestTypeNameForCreator(t *testing.T){
	code := `package com.example;
	import java.io.FileWriter;
	import java.io.File;
	class A{
		public static main(String[] args){
		FileWriter fw = new FileWriter(new File("a.txt"));
		}
	}`
	ssatest.Check(t, code, func(prog *ssaapi.Program) error {
		prog.Show()
		res := prog.SyntaxFlowChain(`File<typeName>?{have:'File'} as $a;`)
		assert.Equal(t,2,res.Len())

		res = prog.SyntaxFlowChain(`File<fullTypeName>?{have:'java.io.File'} as $a;`)
		assert.Equal(t,1,res.Len())

		res = prog.SyntaxFlowChain(`FileWriter<typeName>?{have:'FileWriter'} as $a;`)
		assert.Equal(t,2,res.Len())
		res = prog.SyntaxFlowChain(`FileWriter<typeName>?{have:'java.io.FileWriter'} as $a;`)
		assert.Equal(t,1,res.Len())

		return nil
	}, ssaapi.WithLanguage(consts.JAVA))
	
	
}