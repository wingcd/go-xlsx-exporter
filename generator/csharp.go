package generator

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/wingcd/go-xlsx-protobuf/model"
	"github.com/wingcd/go-xlsx-protobuf/settings"
)

var csharpTemplate = ""

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
	Version   string
	Namespace string
	Enums     []*model.DefineTableInfo
	Structs   []*model.DefineTableInfo
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
		Tables:    make([]*model.DataTable, 0),
	}

	for _, e := range settings.ENUMS {
		fd.Enums = append(fd.Enums, e)
	}

	for _, e := range settings.STRUCTS {
		fd.Structs = append(fd.Structs, e)
	}

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
		settings.PreProcessTable([]*model.DataTable{&table})

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

func init() {
	Regist("csharp", &csharpGenerator{})
}
