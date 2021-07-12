package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/wingcd/go-xlsx-protobuf/model"
	"github.com/wingcd/go-xlsx-protobuf/settings"
)

var goTemplate = ""

var supportGoTypes = map[string]string{
	"bool":   "bool",
	"int":    "int32",
	"int32":  "int32",
	"uint":   "uint32",
	"uint32": "uint32",
	"int64":  "int64",
	"uint64": "uint64",
	"float":  "float32",
	"double": "float64",
	"string": "string",
}

type goFileDesc struct {
	Version string
	Package string
	Enums   []*model.DefineTableInfo
	Tables  []*model.DataTable

	FileRawDesc string
}

type goGenerator struct {
}

var defaultGoValues = map[string]string{
	"bool":   "false",
	"int":    "0",
	"int32":  "0",
	"uint":   "0",
	"uint32": "0",
	"int64":  "0",
	"uint64": "0",
	"float":  "0",
	"double": "0",
	"string": "\"\"",
}

func init() {
	funcs["title"] = strings.Title

	funcs["default"] = func(item interface{}) string {
		var nilType = "nil"
		switch inst := item.(type) {
		case *model.DataTableHeader:
			if inst.IsArray {
				return fmt.Sprintf("make([]*%s, 0)", inst.ValueType)
			} else if inst.IsEnum {
				var enumInfo = settings.GetEnum(inst.ValueType)
				if enumInfo != nil {
					return fmt.Sprintf("%s_%s", enumInfo.TypeName, enumInfo.Items[0].FieldName)
				}
			} else if inst.IsStruct {
				return nilType
			} else if val, ok := defaultGoValues[inst.ValueType]; ok {
				return val
			}
		case *model.DataTable:
			return nilType
		case *model.DefineTableInfo:
			return fmt.Sprintf("%s_%s", inst.TypeName, inst.Items[0].FieldName)
		case string:
			if val, ok := defaultGoValues[inst]; ok {
				return val
			} else if settings.IsEnum(inst) {
				var enumInfo = settings.GetEnum(inst)
				if enumInfo != nil {
					return fmt.Sprintf("%s_%s", enumInfo.TypeName, enumInfo.Items[0].FieldName)
				}
			} else if settings.IsTable(inst) || settings.IsStruct(inst) {
				return nilType
			}
		}
		return ""
	}

	Regist("golang", &goGenerator{})
}

func (g *goGenerator) Generate() *bytes.Buffer {
	if goTemplate == "" {
		data, err := ioutil.ReadFile("./template/golang.gtpl")
		if err != nil {
			log.Println(err)
			return nil
		}
		goTemplate = string(data)
	}

	// 生成proto缓存文件
	Build("proto", "./gen/all.proto")

	tpl, err := template.New("golang").Funcs(funcs).Parse(goTemplate)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	var fd = goFileDesc{
		Version: settings.TOOL_VERSION,
		Package: settings.PackageName,
		Enums:   make([]*model.DefineTableInfo, 0),
		Tables:  make([]*model.DataTable, 0),
	}

	for _, e := range settings.ENUMS {
		fd.Enums = append(fd.Enums, e)
	}

	tables := settings.GetAllTables()
	for _, t := range tables {
		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !settings.IsEnum(h.ValueType) && !settings.IsStruct(h.ValueType) {
				if _, ok := supportGoTypes[h.ValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.ValueType, t.DefinedTable, h.FieldName)
					return nil
				}
				h.ValueType = supportGoTypes[h.ValueType]
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
		header.EncodeType = "bytes"
		header.RawValueType = t.TypeName + "[]"
		table.Headers = []*model.DataTableHeader{&header}

		fd.Tables = append(fd.Tables, &table)
	}

	var buf = bytes.NewBufferString("")

	err = tpl.Execute(buf, &fd)
	if err != nil {
		log.Println(err)
		return nil
	}

	return buf
}
