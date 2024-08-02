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
        JSON a =(JSON) anyJSON;
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

		typeName := prog.SyntaxFlowChain(`anyJSON<typeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "JSON")
		typeName = prog.SyntaxFlowChain(`anyJSON<fullTypeName> as $id`)[0]
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
		typeName := prog.SyntaxFlowChain(`res1<typeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "string")
		typeName = prog.SyntaxFlowChain(`res1<fullTypeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "string")

		typeName = prog.SyntaxFlowChain(`res2<typeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "number")
		typeName = prog.SyntaxFlowChain(`res2<fullTypeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "number")

		typeName = prog.SyntaxFlowChain(`res3<typeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "A")
		typeName = prog.SyntaxFlowChain(`res3<fullTypeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "com.org.A.A")

		typeName = prog.SyntaxFlowChain(`res4<typeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "string")
		typeName = prog.SyntaxFlowChain(`res4<fullTypeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "string")

		typeName = prog.SyntaxFlowChain(`res5<typeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "A")
		typeName = prog.SyntaxFlowChain(`res5<fullTypeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "com.org.A.A")

		typeName = prog.SyntaxFlowChain(`test2<typeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "A")
		typeName = prog.SyntaxFlowChain(`test2<fullTypeName> as $id;`)[0]
		assert.Contains(t, typeName.String(), "com.org.A.A")
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
					var res5 = A.staticMethod();  // fulltypeName 应该找到A
					var res6 = A.noExistMethod();  // fulltypeName 应该找到A
				}
			};	
	`)
		
		ssatest.CheckWithFS(vf, t, func(progs ssaapi.Programs) error {
			prog := progs[0]
			prog.Show()

			typeName := prog.SyntaxFlowChain(`object<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "A")
			typeName = prog.SyntaxFlowChain(`object<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "com.org.A.A")

			typeName = prog.SyntaxFlowChain(`res1<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "A")
			typeName = prog.SyntaxFlowChain(`res1<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "com.org.A.A")

			typeName = prog.SyntaxFlowChain(`res2<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "number")
			typeName = prog.SyntaxFlowChain(`res2<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "number")

			typeName = prog.SyntaxFlowChain(`res3<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "Dog")
			typeName = prog.SyntaxFlowChain(`res3<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "com.org.Dog.Dog")

			typeName = prog.SyntaxFlowChain(`res4<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "A")
			typeName = prog.SyntaxFlowChain(`res4<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "com.org.A.A")

			typeName = prog.SyntaxFlowChain(`res5<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "Dog")
			typeName = prog.SyntaxFlowChain(`res5<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "com.org.Dog.Dog")

			
			typeName = prog.SyntaxFlowChain(`res6<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "A")
			typeName = prog.SyntaxFlowChain(`res6<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "com.org.A.A")

			
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

			typeName := prog.SyntaxFlowChain(`param1<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "number")
			typeName = prog.SyntaxFlowChain(`param1<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "number")

			typeName = prog.SyntaxFlowChain(`param2<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "A")
			typeName = prog.SyntaxFlowChain(`param2<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "com.org.A.A")

			typeName = prog.SyntaxFlowChain(`param3<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "Dog")
			typeName = prog.SyntaxFlowChain(`param3<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "Dog")

			

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

			typeName := prog.SyntaxFlowChain(`res1<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "A")
			typeName = prog.SyntaxFlowChain(`res1<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "com.org.A.A")

			typeName = prog.SyntaxFlowChain(`res2<typeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "number")
			typeName = prog.SyntaxFlowChain(`res2<fullTypeName> as $id;`)[0]
			assert.Contains(t, typeName.String(), "number")

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
			assert.Equal(t,3,typeName.Show().Len())
			typeName = prog.SyntaxFlowChain(`res1<fullTypeName> as $id;`)
			assert.Equal(t,2,typeName.Show().Len())
			
			typeName = prog.SyntaxFlowChain(`res2<typeName> as $id;`)
			assert.Equal(t,3,typeName.Show().Len())
			typeName = prog.SyntaxFlowChain(`res2<fullTypeName> as $id;`)
			assert.Equal(t,2,typeName.Show().Len())

			return nil
		}, ssaapi.WithLanguage(consts.JAVA))
}
