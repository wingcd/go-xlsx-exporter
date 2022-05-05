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

var csharpTemplate = ""

func csFormatValue(value interface{}, valueType string, isEnum bool, isArray bool) string {
	var ret = ""
	if isArray {
		var arr = value.([]interface{})
		var lst []string
		for _, it := range arr {
			lst = append(lst, csFormatValue(it, valueType, isEnum, false))
		}
		ret = fmt.Sprintf("new %s[]{ %s }", valueType, strings.Join(lst, ", "))
	} else if isEnum {
		var enumStr = utils.ToEnumString(valueType, value.(int32))
		if enumStr != "" {
			ret = fmt.Sprintf("%s.%s", valueType, enumStr)
		} else {
			fmt.Printf("[错误] 值解析失败 类型:%s 值：%v \n", valueType, value)
		}
	} else if valueType == "float" {
		ret = fmt.Sprintf("%vf", value)
	} else if valueType == "string" {
		ret = fmt.Sprintf("\"%v\"", value)
	} else {
		ret = fmt.Sprintf("%v", value)
	}
	return ret
}

var csGenetatorInited = false

func registCSFuncs() {
	if csGenetatorInited {
		return
	}
	csGenetatorInited = true

	funcs["get_alias"] = func(alias string) string {
		if alias == "" {
			return "object"
		}
		return alias
	}

	funcs["value_format"] = func(value string, item interface{}) string {
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
		return csFormatValue(val, valueType, isEnum, isArray)
	}
}

var supportCSharpTypes = map[string]string{
	"bool":    "bool",
	"int":     "int",
	"int32":   "int",
	"uint":    "uint",
	"uint32":  "uint",
	"int64":   "long",
	"uint64":  "ulong",
	"float":   "float",
	"float32": "float",
	"double":  "double",
	"float64": "double",
	"string":  "string",
	"bytes":   "byte[]",
	"void":    "object",
}

type csharpFileDesc struct {
	commonFileDesc

	Namespace string
	Info      *BuildInfo
	Enums     []*model.DefineTableInfo
	Structs   []*model.DefineTableInfo
	Consts    []*model.DefineTableInfo
	Tables    []*model.DataTable
}

type csharpGenerator struct {
}

func (g *csharpGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	registCSFuncs()

	if csharpTemplate == "" {
		temp := getTemplate(info, "./template/csharp.gtpl")
		log.Printf("[提示] 加载模板: %s \n", temp)

		data, err := ioutil.ReadFile(temp)
		if err != nil {
			log.Println(err)
			return false, nil
		}
		csharpTemplate = string(data)
	}

	tpl, err := template.New("csharp").Funcs(funcs).Parse(csharpTemplate)
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	var fd = csharpFileDesc{
		Info:      info,
		Namespace: settings.PackageName,
		Enums:     settings.ENUMS[:],
		Structs:   settings.STRUCTS[:],
		Consts:    settings.CONSTS[:],
		Tables:    make([]*model.DataTable, 0),
	}
	fd.Version = settings.TOOL_VERSION
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	utils.PreProcessDefines(fd.Structs)
	utils.PreProcessDefines(fd.Consts)

	for _, c := range fd.Structs {
		for _, it := range c.Items {
			if !it.IsEnum && !it.IsStruct {
				it.ValueType = supportCSharpTypes[it.StandardValueType]
			}
		}
	}

	for _, c := range fd.Consts {
		for _, it := range c.Items {
			if !it.IsEnum && !it.IsStruct {
				it.ValueType = supportCSharpTypes[it.StandardValueType]
			}
		}
	}

	utils.PreProcessTables(settings.TABLES)
	for _, t := range settings.TABLES {
		if t.TableType == model.ETableType_Message {
			fd.HasMessage = true
		}

		// 排除语言类型
		if t.TableType == model.ETableType_Language && !settings.GenLanguageType {
			continue
		}

		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !h.IsEnum && !h.IsStruct {
				if _, ok := supportCSharpTypes[h.StandardValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.RawValueType, t.DefinedTable, h.FieldName)
					return false, nil
				}
				h.ValueType = supportCSharpTypes[h.StandardValueType]
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
	Regist("csharp", &csharpGenerator{})
}
