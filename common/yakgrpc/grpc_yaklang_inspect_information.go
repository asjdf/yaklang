package yakgrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/yaklang/yaklang/common/schema"

	"github.com/jinzhu/gorm"
	"github.com/samber/lo"
	"github.com/yaklang/yaklang/common/consts"
	"github.com/yaklang/yaklang/common/cve/cveresources"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yak/ssaapi"
	pta "github.com/yaklang/yaklang/common/yak/static_analyzer"
	"github.com/yaklang/yaklang/common/yak/static_analyzer/information"
	"github.com/yaklang/yaklang/common/yakgrpc/yakit"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
)

type PluginParamSelect struct {
	Double bool                    `json:"double"`
	Data   []PluginParamSelectData `json:"data"`
}

type PluginParamSelectData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func uiInfo2grpc(info []*information.UIInfo) []*ypb.YakUIInfo {
	ret := make([]*ypb.YakUIInfo, 0, len(info))
	for _, i := range info {
		ret = append(ret, &ypb.YakUIInfo{
			Typ:            i.Typ,
			Effected:       i.Effected,
			WhenExpression: i.WhenExpression,
		})
	}
	return ret
}

func cliParam2grpc(params []*information.CliParameter) []*ypb.YakScriptParam {
	ret := make([]*ypb.YakScriptParam, 0, len(params))

	for _, param := range params {
		defaultValue := ""
		if param.Default != nil {
			defaultValue = fmt.Sprintf("%v", param.Default)
		}
		extra := []byte{}
		if param.Type == "select" {
			paramSelect := &PluginParamSelect{
				Double: param.MultipleSelect,
				Data:   make([]PluginParamSelectData, 0),
			}
			for k, v := range param.SelectOption {
				paramSelect.Data = append(paramSelect.Data, PluginParamSelectData{
					Key:   k,
					Value: v,
				})
			}
			extra, _ = json.Marshal(paramSelect)
		}

		ret = append(ret, &ypb.YakScriptParam{
			Field:        param.Name,
			DefaultValue: string(defaultValue),
			TypeVerbose:  param.Type,
			FieldVerbose: param.NameVerbose,
			Help:         param.Help,
			Required:     param.Required,
			Group:        param.Group,
			ExtraSetting: string(extra),
			MethodType:   param.MethodType,
		})
	}

	return ret
}

func riskInfo2grpc(info []*information.RiskInfo, db *gorm.DB) []*ypb.YakRiskInfo {
	ret := make([]*ypb.YakRiskInfo, 0, len(info))
	for _, i := range info {
		description := i.Description
		solution := i.Solution

		if (description == "" || solution == "") && i.CVE != "" {
			if db != nil {
				cve, err := cveresources.GetCVE(db, i.CVE)
				if err == nil {
					if description == "" {
						description = cve.DescriptionMainZh
					}
					if solution == "" {
						solution = cve.Solution
					}
					if i.Level == "" {
						i.Level = cve.Severity
					}
				}
			}
		}

		ret = append(ret, &ypb.YakRiskInfo{
			Level:       i.Level,
			TypeVerbose: i.TypeVerbose,
			CVE:         i.CVE,
			Description: description,
			Solution:    solution,
		})
	}
	return ret
}

func (s *Server) YaklangInspectInformation(ctx context.Context, req *ypb.YaklangInspectInformationRequest) (*ypb.YaklangInspectInformationResponse, error) {
	ret := &ypb.YaklangInspectInformationResponse{}
	prog, err := ssaapi.Parse(req.YakScriptCode, pta.GetPluginSSAOpt(req.YakScriptType)...)
	if err != nil {
		return nil, errors.New("ssa parse error")
	}
	parameters, uiInfos := information.ParseCliParameter(prog)
	ret.CliParameter = cliParam2grpc(parameters)
	ret.UIInfo = uiInfo2grpc(uiInfos)
	ret.RiskInfo = riskInfo2grpc(information.ParseRiskInfo(prog), consts.GetGormCVEDatabase())
	ret.Tags = information.ParseTags(prog)

	return ret, nil
}

// CompareParameter p1 and p2 all field
// if return true, p2 information is more than p1
func CompareParameter(p1, p2 *ypb.YakScriptParam) bool {
	if len(p1.Field) < len(p2.Field) {
		return true
	}
	if len(p1.FieldVerbose) < len(p2.FieldVerbose) {
		return true
	}
	if len(p1.TypeVerbose) < len(p2.TypeVerbose) {
		return true
	}
	if len(p1.Help) < len(p2.Help) {
		return true
	}
	if len(p1.DefaultValue) < len(p2.DefaultValue) {
		return true
	}
	if len(p1.Group) < len(p2.Group) {
		return true
	}
	if p1.Required != p2.Required {
		return false
	}
	if len(p1.ExtraSetting) < len(p2.ExtraSetting) {
		return true
	}
	return false
}

func getCliCodeFromParam(params []*ypb.YakScriptParam) string {
	code := ""
	for _, para := range params {
		// switch para.
		Option := make([]string, 0)
		cliFunction := ""
		var cliDefault string
		methodType := para.MethodType
		if methodType == "" {
			methodType = para.TypeVerbose
		}
		switch methodType {
		case "string":
			cliFunction = "String"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		case "boolean":
			cliFunction = "Bool"
			if para.DefaultValue != "" {
				b, err := strconv.ParseBool(para.DefaultValue)
				if err == nil {
					cliDefault = fmt.Sprintf("cli.setDefault(%t)", b)
				}
			}
		case "uint":
			cliFunction = "Int"
			if para.DefaultValue != "" {
				i, err := strconv.ParseInt(para.DefaultValue, 10, 64)
				if err == nil {
					cliDefault = fmt.Sprintf("cli.setDefault(%d)", i)
				}
			}
		case "float":
			cliFunction = "Float"
			if para.DefaultValue != "" {
				f, err := strconv.ParseFloat(para.DefaultValue, 64)
				if err == nil {
					cliDefault = fmt.Sprintf("cli.setDefault(%f)", f)
				}
			}
		case "file":
			cliFunction = "File"
		case "file-name":
			cliFunction = "FileNames"
		case "select":
			cliFunction = "StringSlice"
			if para.ExtraSetting != "" {
				var dataSelect *PluginParamSelect
				json.Unmarshal([]byte(para.ExtraSetting), &dataSelect)
				Option = append(Option, fmt.Sprintf(`cli.setMultipleSelect(%t)`, dataSelect.Double))
				for _, v := range dataSelect.Data {
					Option = append(Option, fmt.Sprintf(`cli.setSelectOption(%#v, %#v)`, v.Key, v.Value))
				}
			}
		case "yak":
			cliFunction = "YakCode"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		case "http-packet":
			cliFunction = "HTTPPacket"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		case "text":
			cliFunction = "Text"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		case "urls":
			cliFunction = "Urls"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		case "ports":
			cliFunction = "Ports"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		case "hosts":
			cliFunction = "Hosts"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		case "file_content":
			cliFunction = "FileOrContent"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		case "line_dict":
			cliFunction = "LineDict"
			if para.DefaultValue != "" {
				cliDefault = fmt.Sprintf("cli.setDefault(%#v)", para.DefaultValue)
			}
		default:
			// cliFunction = "Undefine-" + para.TypeVerbose
			continue
		}

		if cliDefault != "" {
			Option = append(Option, cliDefault)
		}

		if para.Help != "" {
			Option = append(Option, fmt.Sprintf(`cli.setHelp(%#v)`, para.Help))
		}
		if para.FieldVerbose != "" && para.FieldVerbose != para.Field {
			Option = append(Option, fmt.Sprintf(`cli.setVerboseName(%#v)`, para.FieldVerbose))
		}
		if para.Group != "" {
			Option = append(Option, fmt.Sprintf(`cli.setCliGroup(%#v)`, para.Group))
		}
		if para.Required {
			Option = append(Option, fmt.Sprintf(`cli.setRequired(%t)`, para.Required))
		}

		var str string
		if len(Option) == 0 {
			str = fmt.Sprintf(`cli.%s(%#v)`, cliFunction, para.Field)
		} else {
			str = fmt.Sprintf(`cli.%s(%#v, %s)`, cliFunction, para.Field, strings.Join(Option, ","))
		}
		code += str + "\n"
	}
	return code
}

func getParameterFromParamJson(j string) ([]*ypb.YakScriptParam, error) {
	params, err := strconv.Unquote(j)
	if err != nil {
		return nil, utils.Wrapf(err, "unquote error")
	}
	var paras []*ypb.YakScriptParam
	err = json.Unmarshal([]byte(params), &paras)
	if err != nil {
		return nil, utils.Wrapf(err, "unmarshal error")
	}
	// return GetCliCodeFromParam(paras), nil
	return paras, nil
}

func getNeedReturn(script *schema.YakScript) ([]*ypb.YakScriptParam, error) {
	prog, err := ssaapi.Parse(script.Content, pta.GetPluginSSAOpt("yak")...)
	if err != nil {
		return nil, errors.New("ssa parse error")
	}
	parameters, _ := information.ParseCliParameter(prog)
	codeParameter := lo.SliceToMap(
		cliParam2grpc(parameters),
		func(ysp *ypb.YakScriptParam) (string, *ypb.YakScriptParam) {
			return ysp.Field, ysp
		},
	)
	databaseParameter, err := getParameterFromParamJson(script.Params)
	if err != nil {
		return nil, utils.Wrapf(err, "get cli code from param json error")
	}

	log.Infof("codeParameter: %d, databaseParameter: %d", len(codeParameter), len(databaseParameter))
	log.Infof("codeParameter: %v", codeParameter)
	log.Infof("databaseParameter: %v", databaseParameter)

	if len(codeParameter) < len(databaseParameter) {
		return databaseParameter, nil
	}

	// compare codeParameter and databaseParameter, and add rest in databaseParameter to needReturn
	for _, databasePara := range databaseParameter {
		codePara, ok := codeParameter[databasePara.Field]
		if !ok {
			return databaseParameter, nil
		}
		if CompareParameter(codePara, databasePara) {
			return databaseParameter, nil
		}
	}
	return nil, nil
}

func (s *Server) YaklangGetCliCodeFromDatabase(ctx context.Context, req *ypb.YaklangGetCliCodeFromDatabaseRequest) (*ypb.YaklangGetCliCodeFromDatabaseResponse, error) {
	name := req.GetScriptName()
	if name == "" {
		return nil, utils.Errorf("script name is empty")
	}
	script, err := yakit.GetYakScriptByName(s.GetProfileDatabase(), name)
	if err != nil {
		return nil, utils.Wrapf(err, "get script %s error", name)
	}
	if script.Type != "yak" {
		return nil, utils.Errorf("script %s is not yak script", name)
	}
	var code string
	if need, err := getNeedReturn(script); err == nil && len(need) != 0 {
		code = fmt.Sprintf(`/*
// this code generated by yaklang from database
%s
*/`, getCliCodeFromParam(need))
	} else {
		log.Error(err)
		code = ""
	}
	return &ypb.YaklangGetCliCodeFromDatabaseResponse{
		Code:       code,
		NeedHandle: len(code) != 0,
	}, nil
}
