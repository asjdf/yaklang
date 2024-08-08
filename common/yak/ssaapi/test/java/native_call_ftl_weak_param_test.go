package java

import (
	"testing"

	"github.com/yaklang/yaklang/common/utils/filesys"
	"github.com/yaklang/yaklang/common/yak/ssaapi"
	"github.com/yaklang/yaklang/common/yak/ssaapi/test/ssatest"
	"gotest.tools/v3/assert"
)

func TestNativeCall_FreeMakerXSS(t *testing.T) {
	vf := filesys.NewVirtualFs()

	vf.AddFile("com/example/demo/controller/freemakerdemo/FreeMakerDemo.java",`package com.example.demo.controller.freemakerdemo;
    
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.ui.freemarker.FreeMarkerTemplateUtils;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import freemarker.template.Configuration;
import freemarker.template.Template;

import java.io.PrintWriter;
@Controller
@RequestMapping("/freemarker")
public class FreeMakerDemo {
    @Autowired
    private Configuration freemarkerConfig;

    @GetMapping("/welcome")
    public String welcome(@RequestParam String name, Model model) {
        if (name == null || name.isEmpty()) {
            model.addAttribute("name", "Welcome to Safe FreeMarker Demo, try <code>/freemarker/safe/welcome?name=Hacker<>");
        } else {
            model.addAttribute("name", name);
        }
        return "welcome";
    }

}`)
	vf.AddFile("src/main/resources/application.properties",`spring.application.name=demo
# freemaker
spring.freemarker.template-loader-path=classpath:/templates/
spring.freemarker.suffix=.ftl


spring.sql.init.mode=always`)
	vf.AddFile("welcome.ftl",`<!DOCTYPE html>
<html>
<head>
    <title>Welcome</title>
</head>
<body>
<h1>Welcome ${name1}!</h1>
<h1>${name2?html}!</h1>
</body>
</html>
`)
ssatest.CheckWithFS(vf, t, func(programs ssaapi.Programs) error {
		prog := programs[0]
		sink := prog.SyntaxFlowChain("*Mapping.__ref__<getFunc><getReturns>?{<typeName>?{have:'string'}}<freeMarkerSink>  as  $a")
        assert.Equal(t, 1, sink.Len())
		return nil
	}, ssaapi.WithLanguage(ssaapi.JAVA))
}

func TestNativeCall_FreeMakerXSS_WithNoSuffixConfig(t *testing.T) {
	vf := filesys.NewVirtualFs()

	vf.AddFile("com/example/demo/controller/freemakerdemo/FreeMakerDemo.java",`package com.example.demo.controller.freemakerdemo;
    
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.ui.freemarker.FreeMarkerTemplateUtils;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import freemarker.template.Configuration;
import freemarker.template.Template;

import java.io.PrintWriter;
@Controller
@RequestMapping("/freemarker")
public class FreeMakerDemo {
    @Autowired
    private Configuration freemarkerConfig;

    @GetMapping("/welcome")
    public String welcome(@RequestParam String name, Model model) {
        if (name == null || name.isEmpty()) {
            model.addAttribute("name", "Welcome to Safe FreeMarker Demo, try <code>/freemarker/safe/welcome?name=Hacker<>");
        } else {
            model.addAttribute("name", name);
        }
        return "welcome.ftl";
    }

}`)
	vf.AddFile("src/main/resources/application.properties",`spring.application.name=demo
# freemaker
spring.freemarker.template-loader-path=classpath:/templates/
`)
	vf.AddFile("welcome.ftl",`<!DOCTYPE html>
<html>
<head>
    <title>Welcome</title>
</head>
<body>
<h1>Welcome ${name}!</h1>
</body>
</html>
`)

ssatest.CheckWithFS(vf, t, func(programs ssaapi.Programs) error {
		prog := programs[0]
		sink := prog.SyntaxFlowChain("*Mapping.__ref__<getFunc><getReturns>?{<typeName>?{have:'string'}}<freeMarkerSink>  as  $a")
        assert.Equal(t, 1, sink.Len())
		return nil
	}, ssaapi.WithLanguage(ssaapi.JAVA))
}
