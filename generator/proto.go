package generator

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
)

var protoTemplate = ""

var supportProtoTypes = map[string]string{
	"bool":    "bool",
	"int":     "int32",
	"int32":   "int32",
	"uint":    "uint32",
	"uint32":  "uint32",
	"int64":   "int64",
	"uint64":  "uint64",
	"float":   "float",
	"float32": "float",
	"double":  "double",
	"float64": "double",
	"string":  "string",
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

func (g *protoGenerator) Generate(output string) (save bool, data *bytes.Buffer) {
	if protoTemplate == "" {
		data, err := ioutil.ReadFile("./template/proto.gtpl")
		if err != nil {
			log.Println(err)
			return false, nil
		}
		protoTemplate = string(data)
	}

	tpl, err := template.New("proto").Funcs(funcs).Parse(protoTemplate)
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	var fd = protoFileDesc{
		Version: settings.TOOL_VERSION,
		Package: settings.PackageName,
		Enums:   settings.ENUMS[:],
		Tables:  make([]*model.DataTable, 0),
	}
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	tables := settings.GetAllTables()
	utils.PreProcessTable(tables)
	for _, t := range tables {
		// 排除语言类型
		if t.IsLanguage && !settings.GenLanguageType {
			continue
		}

		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !h.IsEnum && !h.IsStruct {
				if _, ok := supportProtoTypes[h.ValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.ValueType, t.DefinedTable, h.FieldName)
					return false, nil
				}
				h.ValueType = supportProtoTypes[h.ValueType]
			}
		}

		if t.IsDataTable {
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
			table.Headers = []*model.DataTableHeader{&header}

			fd.Tables = append(fd.Tables, &table)
		}
	}
	utils.PreProcessTable(fd.Tables)

	var buf = bytes.NewBufferString("")

	err = tpl.Execute(buf, &fd)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, buf
}

func init() {
	Regist("proto", &protoGenerator{})
}
