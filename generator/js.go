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

var jsTemplate = ""

func jsFormatValue(value interface{}, valueType string, isEnum bool, isArray bool) string {
	var ret = ""
	if isArray {
		var arr = value.([]interface{})
		var lst []string
		for _, it := range arr {
			lst = append(lst, jsFormatValue(it, valueType, isEnum, false))
		}
		ret = fmt.Sprintf("[ %s ]", strings.Join(lst, ", "))
	} else if isEnum {
		var enumStr = utils.ToEnumString(valueType, value.(int32))
		if enumStr != "" {
			ret = fmt.Sprintf("%s.%s", valueType, enumStr)
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

var jsGenetatorInited = false

func jsValueDefault(item interface{}) string {
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
		} else if val, ok := defaultJsValue[inst.StandardValueType]; ok {
			return val
		}
	case *model.DataTable:
		return nilType
	case *model.DefineTableInfo:
		return fmt.Sprintf("%s_%s", inst.TypeName, inst.Items[0].FieldName)
	case string:
		if val, ok := defaultJsValue[inst]; ok {
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

func jsValueFormat(value string, item interface{}) string {
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
	return jsFormatValue(val, valueType, isEnum, isArray)
}

func registJSFuncs() {
	if jsGenetatorInited {
		return
	}
	jsGenetatorInited = true

	funcs["value_format"] = jsValueFormat

	funcs["default"] = jsValueDefault

	funcs["get_alias"] = func(alias string) string {
		if alias == "" {
			return "any"
		}
		return alias
	}
}

var supportJSTypes = map[string]string{
	"bool":   "bool",
	"int":    "int",
	"uint":   "uint",
	"int64":  "int64",
	"uint64": "uint64",
	"float":  "float",
	"double": "double",
	"string": "string",
	"bytes":  "Uint8Array",
	"void":   "",
}

var defaultJsValue = map[string]string{
	"bool":   "false",
	"int":    "0",
	"uint":   "0",
	"int64":  "0",
	"uint64": "0",
	"float":  "0",
	"double": "0",
	"string": "\"\"",
	"bytes":  "new Uint8Array(0)",
	"void":   "null",
}

type jsFileDesc struct {
	commonFileDesc

	Namespace string
	Info      *BuildInfo
	Enums     []*model.DefineTableInfo
	Consts    []*model.DefineTableInfo
	Tables    []*model.DataTable
}

type jsGenerator struct {
}

func (g *jsGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	registJSFuncs()

	if jsTemplate == "" {
		temp := getTemplate(info, "./template/js.gtpl")
		log.Printf("[提示] 加载模板: %s \n", temp)

		data, err := ioutil.ReadFile(temp)
		if err != nil {
			log.Println(err)
			return false, nil
		}
		jsTemplate = string(data)
	}

	tpl, err := template.New("js").Funcs(funcs).Parse(jsTemplate)
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	var fd = jsFileDesc{
		Namespace: settings.PackageName,
		Info:      info,
		Enums:     settings.ENUMS[:],
		Consts:    settings.CONSTS[:],
		Tables:    make([]*model.DataTable, 0),
	}
	fd.Version = settings.TOOL_VERSION
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	utils.PreProcessDefines(fd.Consts)

	tables := settings.GetAllTables()
	utils.PreProcessTables(tables)
	for _, t := range tables {
		if t.TableType == model.ETableType_Message {
			fd.HasMessage = true
		}

		// 排除语言类型
		if t.TableType == model.ETableType_Language && !settings.GenLanguageType {
			continue
		}

		// 排除配置
		if t.TableType == model.ETableType_Define {
			continue
		}

		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !h.IsEnum && !h.IsStruct {
				if _, ok := supportJSTypes[h.StandardValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.RawValueType, t.DefinedTable, h.FieldName)
					return false, nil
				}
			}
		}

		if t.NeedAddItems {
			// 添加数组类型
			table := model.DataTable{}
			table.DefinedTable = t.DefinedTable
			table.TypeName = t.TypeName + "_ARRAY"
			table.IsArray = true
			header := model.DataTableHeader{}
			header.Index = 1
			header.FieldName = "Items"
			header.TitleFieldName = header.FieldName
			header.IsArray = true
			header.ValueType = t.TypeName
			header.RawValueType = t.TypeName + "[]"
			header.IsMessage = true
			table.Headers = []*model.DataTableHeader{&header}

			fd.Tables = append(fd.Tables, &table)
		}
	}
	utils.PreProcessTables(fd.Tables)

	var buf = bytes.NewBufferString("")

	err = tpl.Execute(buf, &fd)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, buf
}

func init() {
	Regist("js", &jsGenerator{})
}
