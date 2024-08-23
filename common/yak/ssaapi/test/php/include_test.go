package php

import (
	_ "embed"
	"github.com/stretchr/testify/require"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yak/ssaapi"
	"github.com/yaklang/yaklang/common/yak/ssaapi/test/ssatest"
	"strings"
	"testing"
)

//go:embed phpcode/vul/exec.php
var ExecCode string

func TestInclude(t *testing.T) {
	t.Run("test-include", func(t *testing.T) {
		code := `<?php
//单一入口模式
error_reporting(0); //关闭错误显示
$file=addslashes($_GET['r']); //接收文件名
$action=$file==''?'index':$file; //判断为空或者等于index
include('files/'.$action.'.php'); //载入相应文件
?>`
		ssatest.Check(t, code, func(prog *ssaapi.Program) error {
			prog.Show()
			result := prog.SyntaxFlow(`include(* #-> *<getMembers><name>?{have: get||have: post} as $include)`)
			if strings.Contains(result.String(), ".get.r") {
				return nil
			} else {
				return utils.Error("not match")
			}
		}, ssaapi.WithLanguage(ssaapi.PHP))
	})
	t.Run("test-exec", func(t *testing.T) {
		ssatest.Check(t, ExecCode, func(prog *ssaapi.Program) error {
			results, err := prog.SyntaxFlowWithError(`exec(* #-> * as $param)`)
			require.NoError(t, err)
			var flag bool
			values := results.GetValues("param")
			values.ForEach(func(value *ssaapi.Value) {
				if strings.Contains(value.String(), "request") {
					flag = true
				}
			})
			require.True(t, flag)
			return nil
		}, ssaapi.WithLanguage(ssaapi.PHP))
	})
}
