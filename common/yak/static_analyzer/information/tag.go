package information

import (
	"github.com/yaklang/yaklang/common/yak/ssaapi"
)

var (
	AI_PLUGIN           = "AI工具"
	DROP_HTTP_PACKET    = "可能丢弃HTTP数据包"
	FORWARD_HTTP_PACKET = "可能修改HTTP数据包"
)

func ParseTags(prog *ssaapi.Program) []string {
	ret := make([]string, 0)
	{
		// check ai
		res, err := prog.SyntaxFlowWithError("ai.Chat() as $target")
		if err == nil {
			if vs := res.GetValues("target"); vs.Len() > 0 {
				ret = append(ret, AI_PLUGIN)
			}
		}
	}

	{
		if res, err := prog.SyntaxFlowWithError(`
		hijackHTTPRequest(, , ,*() as $forward , *() as $drop)
		hijackHTTPResponse(, , ,*() as $forward , *() as $drop)
		hijackHTTPResponseEx(, , ,, *() as $forward , *() as $drop)
		`); err == nil {
			if vs := res.GetValues("forward"); vs.Len() > 0 {
				ret = append(ret, FORWARD_HTTP_PACKET)
			}
			if vs := res.GetValues("drop"); vs.Len() > 0 {
				ret = append(ret, DROP_HTTP_PACKET)
			}
		}
	}
	return ret
}

func GetLastRef(vs ssaapi.Values) *ssaapi.Value {
	var ret *ssaapi.Value
	vs.ForEach(func(v *ssaapi.Value) {
		if ret == nil {
			ret = v
		}
		if v.GetRange().GetOffset() > ret.GetRange().GetOffset() {
			ret = v
		}
	})
	return ret
}
