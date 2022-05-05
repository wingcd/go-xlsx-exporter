package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
)

var jsonTemplate = ""

func jsonFormatValue(value interface{}, valueType string, isEnum bool, isArray bool) string {
	var ret = ""
	if isArray {
		var arr = value.([]interface{})
		var lst []string
		for _, it := range arr {
			lst = append(lst, jsonFormatValue(it, valueType, isEnum, false))
		}
		ret = fmt.Sprintf("[%s]", strings.Join(lst, ","))
	} else if isEnum {
		var enumStr = utils.ToEnumString(valueType, value.(int32))
		if enumStr != "" {
			ret = fmt.Sprintf("%v", value)
		} else {
			fmt.Printf("[错误] 值解析失败 类型:%s 值：%v \n", valueType, value)
		}
	} else if valueType == "float" {
		ret = fmt.Sprintf("%v", value)
	} else if valueType == "string" {
		ret = fmt.Sprintf("\"%v\"", value)
	} else {
		ret = fmt.Sprintf("%v", value)
	}
	return ret
}

var jsonGenetatorInited = false

func jsonValueDefault(item interface{}) string {
	var nilType = "null"
	switch inst := item.(type) {
	case *model.DataTableHeader:
		if inst.IsArray {
			return nilType
		} else if inst.IsEnum {
			var enumInfo = settings.GetEnum(inst.ValueType)
			if enumInfo != nil {
				return fmt.Sprintf("%s.%s", enumInfo.TypeName, enumInfo.Items[0].FieldName)
			}
		} else if inst.IsStruct {
			return nilType
		} else if val, ok := defaultJsonValue[inst.StandardValueType]; ok {
			return val
		}
	case *model.DataTable:
		return nilType
	case *model.DefineTableInfo:
		return fmt.Sprintf("%s_%s", inst.TypeName, inst.Items[0].FieldName)
	case string:
		if val, ok := defaultJsonValue[inst]; ok {
			return val
		} else if utils.IsEnum(inst) {
			var enumInfo = settings.GetEnum(inst)
			if enumInfo != nil {
				return fmt.Sprintf("%s_%s", enumInfo.TypeName, enumInfo.Items[0].FieldName)
			}
		} else if utils.IsTable(inst) || utils.IsStruct(inst) {
			return nilType
		}
	}
	return ""
}

func jsonValueFormat(value string, item interface{}) string {
	var isEnum = false
	var valueType = ""
	var rawValueType = ""
	var fieldName = ""
	switch inst := item.(type) {
	case *model.DefineTableItem:
		fieldName = inst.FieldName
		isEnum = inst.IsEnum
		valueType = inst.ValueType
		rawValueType = inst.RawValueType
	case *model.DataTableHeader:
		fieldName = inst.FieldName
		isEnum = inst.IsEnum
		valueType = inst.ValueType
		rawValueType = inst.RawValueType
	}

	var ok, val, isArray = utils.ParseValue(rawValueType, value)
	if !ok {
		fmt.Printf("[错误] 值解析失败 字段：%s 类型:%s 值：%v \n", fieldName, valueType, value)
		return value
	}
	return jsonFormatValue(val, valueType, isEnum, isArray)
}

func registJsonFuncs() {
	if jsonGenetatorInited {
		return
	}
	jsonGenetatorInited = true

	funcs["value_format"] = jsonValueFormat

	funcs["default"] = jsonValueDefault
}

var supportJsonTypes = map[string]string{
	"bool":   "bool",
	"int":    "int",
	"uint":   "uint",
	"int64":  "int64",
	"uint64": "uint64",
	"float":  "float",
	"double": "double",
	"string": "string",
	"bytes":  "string",
	"void":   "",
}

var defaultJsonValue = map[string]string{
	"bool":   "false",
	"int":    "0",
	"uint":   "0",
	"int64":  "0",
	"uint64": "0",
	"float":  "0",
	"double": "0",
	"string": "\"\"",
	"bytes":  "\"\"",
	"void":   "null",
}

type jsonFileDesc struct {
	commonFileDesc

	Namespace string
	Info      *BuildInfo
	Enum      *model.DefineTableInfo
	Const     *model.DefineTableInfo
	Table     *model.DataTable
}

type jsonGenerator struct {
}

func genJsonFile(t *model.DataTable, info *BuildInfo, tpl *template.Template) bool {
	var fd = jsonFileDesc{
		Namespace: settings.PackageName,
		Info:     info,
		Table:    t,
	}

	fd.Version = settings.TOOL_VERSION
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	if t.TableType == model.ETableType_Message {
		fd.HasMessage = true
	}

	// 排除语言类型
	if t.TableType == model.ETableType_Language && !settings.GenLanguageType {
		return false
	}

	// 排除配置
	if t.TableType == model.ETableType_Define {
		return false
	}

	// 处理类型
	for _, h := range t.Headers {
		if !h.IsEnum && !h.IsStruct {
			if _, ok := supportJsonTypes[h.StandardValueType]; !ok {
				log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.RawValueType, t.DefinedTable, h.FieldName)
				return false
			}
		}
	}
	
	utils.PreProcessTable(fd.Table)

	var buf = bytes.NewBufferString("")

	err := tpl.Execute(buf, &fd)
	if err != nil {
		log.Println(err)
		return false
	}

	var fileName = fmt.Sprintf("%s%s.json", info.Output, t.TypeName)
	if err := ioutil.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (g *jsonGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	registJsonFuncs()

	if jsonTemplate == "" {		
		temp := getTemplate(info, "./template/json.gtpl")
		log.Printf("[提示] 加载模板: %s \n", temp)

		data, err := ioutil.ReadFile(temp)
		if err != nil {
			log.Println(err)
			return false, nil
		}
		jsonTemplate = string(data)
	}

	tpl, err := template.New("json").Funcs(funcs).Parse(jsonTemplate)
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}
	
	if info.Output[len(info.Output)-1] != '/' {
		info.Output += "/"
	}

	tables := settings.GetAllTables()
	utils.PreProcessTables(tables)
	for _, t := range tables {
		genJsonFile(t, info, tpl)
	}

	return false, nil
}

func init() {
	Regist("json", &jsonGenerator{})
}
