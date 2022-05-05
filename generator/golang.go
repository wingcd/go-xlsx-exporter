package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
)

var goTemplate = ""

var supportGoTypes = map[string]string{
	"bool":    "bool",
	"int":     "int32",
	"int32":   "int32",
	"uint":    "uint32",
	"uint32":  "uint32",
	"int64":   "int64",
	"uint64":  "uint64",
	"float":   "float32",
	"float32": "float32",
	"double":  "float64",
	"float64": "float64",
	"string":  "string",
	"bytes":   "[]byte",
	"void":    "interface{}",
}

var defaultGoValues = map[string]string{
	"bool":        "false",
	"int":         "0",
	"int32":       "0",
	"uint":        "0",
	"uint32":      "0",
	"int64":       "0",
	"uint64":      "0",
	"float":       "0",
	"float32":     "0",
	"double":      "0",
	"float64":     "0",
	"string":      "\"\"",
	"[]byte":      "[]byte{}",
	"interface{}": "nil",
}

func goFormatValue(value interface{}, valueType string, isEnum bool, isArray bool) string {
	var ret = ""
	if isArray {
		var arr = value.([]interface{})
		var lst []string
		for _, it := range arr {
			lst = append(lst, goFormatValue(it, valueType, isEnum, false))
		}
		ret = fmt.Sprintf("[]%s{ %s }", valueType, strings.Join(lst, ", "))
	} else if isEnum {
		var enumStr = utils.ToEnumString(valueType, value.(int32))
		if enumStr != "" {
			ret = fmt.Sprintf("%s_%s", valueType, enumStr)
		} else {
			fmt.Printf("[错误] 值解析失败 类型:%s 值：%v \n", valueType, value)
		}
	} else if valueType == "float" {
		ret = fmt.Sprintf("float32(%v)", value)
	} else if valueType == "string" {
		ret = fmt.Sprintf("\"%v\"", value)
	} else {
		ret = fmt.Sprintf("%v", value)
	}
	return ret
}

var goGenetatorInited = false

func registGoFuncs() {
	if goGenetatorInited {
		return
	}
	goGenetatorInited = true

	funcs["default"] = func(item interface{}) string {
		var nilType = "nil"
		switch inst := item.(type) {
		case *model.DataTableHeader:
			if inst.IsArray {
				return nilType
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
		return goFormatValue(val, valueType, isEnum, isArray)
	}
}

type goFileDesc struct {
	commonFileDesc

	Package string
	Info    *BuildInfo
	Enums   []*model.DefineTableInfo
	Consts  []*model.DefineTableInfo
	Tables  []*model.DataTable

	FileName    string
	FileRawDesc string
	DepIdexs    string
}

type goGenerator struct {
}

func (f *goFileDesc) genProtoRawDesc() {
	var fd, err = utils.BuildFileDesc(f.FileName, settings.GenLanguageType)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// 生成依赖索引数组
	var depIdxs []string
	var totalIdx = 0

	for i := 0; i < fd.Messages().Len(); i++ {
		var msg = fd.Messages().Get(i)

		for fi := 0; fi < msg.Fields().Len(); fi++ {
			var field = msg.Fields().Get(fi)
			if field.Enum() != nil {
				// 0, // 0: gen.SInfo.DataType:type_name -> gen.EDataType
				depIdxs = append(depIdxs, fmt.Sprintf("%v, // %v: %v:type_name -> %v",
					field.Enum().Index(), totalIdx,
					msg.FullName(),
					field.Enum().FullName(),
				))
				totalIdx++
			} else if field.Message() != nil {
				var baseIdx = fd.Enums().Len()
				depIdxs = append(depIdxs, fmt.Sprintf("%v, // %v: %v:type_name -> %v",
					field.Message().Index()+baseIdx, totalIdx,
					msg.FullName(),
					field.Message().FullName(),
				))
				totalIdx++
			}
		}
	}

	depIdxs = append(depIdxs, fmt.Sprintf("%v, // [%v:%v] is the sub-list for method output_type", totalIdx, totalIdx, totalIdx))
	depIdxs = append(depIdxs, fmt.Sprintf("%v, // [%v:%v] is the sub-list for method input_type", totalIdx, totalIdx, totalIdx))
	depIdxs = append(depIdxs, fmt.Sprintf("%v, // [%v:%v] is the sub-list for extension type_name", totalIdx, totalIdx, totalIdx))
	depIdxs = append(depIdxs, fmt.Sprintf("%v, // [%v:%v] is the sub-list for extension extendee", totalIdx, totalIdx, totalIdx))
	depIdxs = append(depIdxs, fmt.Sprintf("0, // [0:%v] is the sub-list for field type_name", totalIdx))

	f.DepIdexs = strings.Join(depIdxs, "\n")

	// proto.FileDescriptor(XXX)
	// 生成文件描述数据
	pt := protodesc.ToFileDescriptorProto(fd)
	var bts, _ = proto.Marshal(pt)
	if len(bts) > 0 {
		// var v = protoimpl.X.CompressGZIP(b)
		var rets = make([]string, 0)
		for i, b := range bts {
			if (i%16) == 0 && i != 0 {
				rets = append(rets, "\n")
			}
			rets = append(rets, fmt.Sprintf("0x%02x,", b))
		}
		f.FileRawDesc = strings.Join(rets, "")
	}
}

func (g *goGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	registGoFuncs()

	if goTemplate == "" {
		temp := getTemplate(info, "./template/golang.gtpl")
		log.Printf("[提示] 加载模板: %s \n", temp)

		data, err := ioutil.ReadFile(temp)
		if err != nil {
			log.Println(err)
			return false, nil
		}
		goTemplate = string(data)
	}

	tpl, err := template.New("golang").Funcs(funcs).Parse(goTemplate)
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	var filename = strings.Split(filepath.Base(info.Output), ".")[0]
	var fd = goFileDesc{
		Package:  settings.PackageName,
		Info:     info,
		Enums:    settings.ENUMS[:],
		Consts:   settings.CONSTS[:],
		Tables:   make([]*model.DataTable, 0),
		FileName: filename,
	}
	fd.Version = settings.TOOL_VERSION
	fd.GoProtoVersion = settings.GO_PROTO_VERTION
	utils.PreProcessDefines(fd.Consts)

	fd.genProtoRawDesc()

	for _, c := range fd.Consts {
		for _, it := range c.Items {
			if !it.IsEnum && !it.IsStruct {
				it.ValueType = supportGoTypes[it.StandardValueType]
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

		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !h.IsEnum && !h.IsStruct {
				if _, ok := supportGoTypes[h.StandardValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.RawValueType, t.DefinedTable, h.FieldName)
					return false, nil
				}
				h.ValueType = supportGoTypes[h.StandardValueType]
			}
		}

		if t.NeedAddItems {
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
	Regist("golang", &goGenerator{})
}
