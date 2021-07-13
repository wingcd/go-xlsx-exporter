package generator

import (
	"bytes"
	"fmt"
	"go-xlsx-protobuf/utils"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/wingcd/go-xlsx-protobuf/model"
	"github.com/wingcd/go-xlsx-protobuf/settings"
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
		ret = fmt.Sprintf("new %s[]{%s}", valueType, strings.Join(lst, ","))
	} else if isEnum {
		for _, e := range settings.ENUMS {
			if e.TypeName == valueType {
				for _, it := range e.Items {
					if it.Value == value {
						ret = it.FieldName
						break
					}
				}
				break
			}
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

func init() {
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

	Regist("csharp", &csharpGenerator{})
}

var supportCSharpTypes = map[string]string{
	"bool":   "bool",
	"int":    "int",
	"int32":  "int",
	"uint":   "uint",
	"uint32": "uint",
	"int64":  "long",
	"uint64": "ulong",
	"float":  "float",
	"double": "double",
	"string": "string",
}

type csharpFileDesc struct {
	commonFileDesc

	Version   string
	Namespace string
	Enums     []*model.DefineTableInfo
	Structs   []*model.DefineTableInfo
	Consts    []*model.DefineTableInfo
	Tables    []*model.DataTable
}

type csharpGenerator struct {
}

func (g *csharpGenerator) SetOutput(output string) {

}

func (g *csharpGenerator) Generate() *bytes.Buffer {
	if csharpTemplate == "" {
		data, err := ioutil.ReadFile("./template/csharp.gtpl")
		if err != nil {
			log.Println(err)
			return nil
		}
		csharpTemplate = string(data)
	}

	tpl, err := template.New("csharp").Funcs(funcs).Parse(csharpTemplate)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	var fd = csharpFileDesc{
		Version:   settings.TOOL_VERSION,
		Namespace: settings.PackageName,
		Enums:     make([]*model.DefineTableInfo, 0),
		Structs:   make([]*model.DefineTableInfo, 0),
		Consts:    make([]*model.DefineTableInfo, 0),
		Tables:    make([]*model.DataTable, 0),
	}
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	for _, e := range settings.ENUMS {
		fd.Enums = append(fd.Enums, e)
	}

	for _, e := range settings.STRUCTS {
		fd.Structs = append(fd.Structs, e)
	}
	settings.PreProcessDefine(fd.Structs)

	for _, e := range settings.CONSTS {
		fd.Consts = append(fd.Consts, e)
	}
	settings.PreProcessDefine(fd.Consts)

	for _, t := range settings.TABLES {
		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !settings.IsEnum(h.ValueType) && !settings.IsStruct(h.ValueType) {
				if _, ok := supportCSharpTypes[h.ValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.ValueType, t.DefinedTable, h.FieldName)
					return nil
				}
				h.ValueType = supportCSharpTypes[h.ValueType]
			}
		}

		// 添加数组类型
		table := model.DataTable{}
		table.DefinedTable = t.DefinedTable
		table.TypeName = t.TypeName + "_ARRAY"
		header := model.DataTableHeader{}
		header.Index = 1
		header.FieldName = "Items"
		header.TitleFieldName = header.FieldName
		header.IsArray = true
		header.ValueType = t.TypeName
		header.RawValueType = t.TypeName + "[]"
		table.Headers = []*model.DataTableHeader{&header}

		fd.Tables = append(fd.Tables, &table)
	}
	settings.PreProcessTable(fd.Tables)

	var buf = bytes.NewBufferString("")

	err = tpl.Execute(buf, &fd)
	if err != nil {
		log.Println(err)
		return nil
	}

	return buf
}
