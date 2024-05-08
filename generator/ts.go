package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
)

var tsTemplate = ""

var tsGenetatorInited = false

func tsTypeFormat(standardValueType, valueType string, isArray bool) string {
	if isArray {
		return fmt.Sprintf("(%s[]|null)", valueType)
	}
	return fmt.Sprintf("(%s|null)", valueType)
}

func egistTSFuncs() {
	if tsGenetatorInited {
		return
	}
	tsGenetatorInited = true

	funcs["value_format"] = jsValueFormat

	funcs["type_format"] = tsTypeFormat

	funcs["default"] = jsValueDefault

	funcs["get_alias"] = func(alias string) string {
		if alias == "" {
			return "any"
		}
		return alias
	}
}

var supportTSTypes = map[string]string{
	"bool":   "boolean",
	"int":    "number",
	"uint":   "number",
	"int64":  "number",
	"uint64": "number",
	"float":  "number",
	"double": "number",
	"string": "string",
	"bytes":  "Uint8Array",
	"void":   "any",
}

var sdefaultTSValue = map[string]string{
	"bool":   "false",
	"int":    "0",
	"uint":   "0",
	"int64":  "0",
	"uint64": "0",
	"float":  "0",
	"double": "0",
	"string": "\"\"",
	"bytes":  "new ArrayBuffer(0)",
	"void":   "null",
}

type tsFileDesc struct {
	commonFileDesc

	Namespace string
	Info      *BuildInfo
	Enums     []*model.DefineTableInfo
	Consts    []*model.DefineTableInfo
	Tables    []*model.DataTable
}

type tsGenerator struct {
}

func (g *tsGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	egistTSFuncs()

	if tsTemplate == "" {
		temp := getTemplate(info, "./template/ts.gtpl")
		log.Printf("[提示] 加载模板: %s \n", temp)

		data, err := ioutil.ReadFile(temp)
		if err != nil {
			log.Println(err)
			return false, nil
		}
		tsTemplate = string(data)
	}

	tpl, err := template.New("ts").Funcs(funcs).Parse(tsTemplate)
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	var fd = tsFileDesc{
		Namespace: settings.PackageName,
		Info:      info,
		Enums:     settings.ENUMS[:],
		Consts:    settings.CONSTS[:],
		Tables:    make([]*model.DataTable, 0),
	}
	fd.Version = settings.TOOL_VERSION
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	utils.PreProcessDefines(fd.Consts)
	for _, c := range fd.Consts {
		for _, it := range c.Items {
			if !it.IsEnum && !it.IsStruct {
				it.ValueType = supportTSTypes[it.StandardValueType]
			}
		}
	}

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
				if _, ok := supportTSTypes[h.StandardValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.RawValueType, t.DefinedTable, h.FieldName)
					return false, nil
				}
				h.ValueType = supportTSTypes[h.StandardValueType]
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
	Regist("ts", &tsGenerator{})
}
