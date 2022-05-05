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
	"bytes":   "bytes",
	"void":    "",
}

type protoFileDesc struct {
	commonFileDesc

	Version string
	Package string
	Info    *BuildInfo
	Enums   []*model.DefineTableInfo
	Tables  []*model.DataTable
}

type protoGenerator struct {
}

func (g *protoGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	if protoTemplate == "" {
		temp := getTemplate(info, "./template/proto.gtpl")
		log.Printf("[提示] 加载模板: %s \n", temp)

		data, err := ioutil.ReadFile(temp)
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
		Info:    info,
		Enums:   settings.ENUMS[:],
		Tables:  make([]*model.DataTable, 0),
	}
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

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

		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !h.IsEnum && !h.IsStruct {
				if _, ok := supportProtoTypes[h.StandardValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.RawValueType, t.DefinedTable, h.FieldName)
					return false, nil
				}
				h.ValueType = supportProtoTypes[h.StandardValueType]
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
			header.StandardValueType = t.TypeName
			_, header.PBValueType = utils.ToPBType(t.TypeName)
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
	Regist("proto", &protoGenerator{})
}
