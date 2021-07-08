package generator

import (
	"bytes"
	"io/ioutil"
	"log"
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
}

type goGenerator struct {
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

	tpl, err := template.New("golang").Parse(goTemplate)
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
		header.IsArray = true
		header.ValueType = t.TypeName
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

func init() {
	Regist("golang", &goGenerator{})
}
