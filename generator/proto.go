package generator

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/wingcd/go-xlsx-protobuf/model"
	"github.com/wingcd/go-xlsx-protobuf/settings"
)

var protoTemplate = ""

var supportProtoTypes = map[string]string{
	"bool":   "bool",
	"int":    "int32",
	"int32":  "int32",
	"uint":   "fixed32",
	"uint32": "fixed32",
	"int64":  "int64",
	"uint64": "fixed64",
	"float":  "float",
	"double": "double",
	"string": "string",
}

type protoFileDesc struct {
	commonFileDesc

	Version string
	Package string
	Enums   []*model.DefineTableInfo
	Tables  []*model.DataTable
}

type protoGenerator struct {
}

func (g *protoGenerator) SetOutput(output string) {

}

func (g *protoGenerator) Generate() *bytes.Buffer {
	if protoTemplate == "" {
		data, err := ioutil.ReadFile("./template/proto.gtpl")
		if err != nil {
			log.Println(err)
			return nil
		}
		protoTemplate = string(data)
	}

	tpl, err := template.New("proto").Funcs(funcs).Parse(protoTemplate)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	var fd = protoFileDesc{
		Version: settings.TOOL_VERSION,
		Package: settings.PackageName,
		Enums:   make([]*model.DefineTableInfo, 0),
		Tables:  make([]*model.DataTable, 0),
	}
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	for _, e := range settings.ENUMS {
		fd.Enums = append(fd.Enums, e)
	}

	tables := settings.GetAllTables()
	for _, t := range tables {
		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !settings.IsEnum(h.ValueType) && !settings.IsStruct(h.ValueType) {
				if _, ok := supportProtoTypes[h.ValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.ValueType, t.DefinedTable, h.FieldName)
					return nil
				}
				h.ValueType = supportProtoTypes[h.ValueType]
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
	Regist("proto", &protoGenerator{})
}
