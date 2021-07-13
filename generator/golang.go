package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wingcd/go-xlsx-protobuf/model"
	"github.com/wingcd/go-xlsx-protobuf/settings"
	"github.com/wingcd/go-xlsx-protobuf/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
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

type goFileDesc struct {
	commonFileDesc

	Version string
	Package string
	Enums   []*model.DefineTableInfo
	Tables  []*model.DataTable

	FileName    string
	FileRawDesc string
	DepIdexs    string
}

type goGenerator struct {
	output string
}

func (g *goGenerator) SetOutput(output string) {
	g.output = output
}

func (f *goFileDesc) genProtoRawDesc() {
	var fd, err = utils.BuildFileDesc(f.FileName)
	if err != nil {
		fmt.Errorf(err.Error())
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
	var b, _ = proto.Marshal(pt)
	if len(b) > 0 {
		// var v = protoimpl.X.CompressGZIP(b)
		var rets = make([]string, 0)
		for i, b := range b {
			if (i%16) == 0 && i != 0 {
				rets = append(rets, "\n")
			}
			rets = append(rets, fmt.Sprintf("0x%02x,", b))
		}
		f.FileRawDesc = strings.Join(rets, "")
	}
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

	tpl, err := template.New("golang").Funcs(funcs).Parse(goTemplate)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	var filename = strings.Split(filepath.Base(g.output), ".")[0]
	var fd = goFileDesc{
		Version:  settings.TOOL_VERSION,
		Package:  settings.PackageName,
		Enums:    make([]*model.DefineTableInfo, 0),
		Tables:   make([]*model.DataTable, 0),
		FileName: filename,
	}
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	fd.genProtoRawDesc()

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
	settings.PreProcessTable(fd.Tables)

	var buf = bytes.NewBufferString("")

	err = tpl.Execute(buf, &fd)
	if err != nil {
		log.Println(err)
		return nil
	}

	return buf
}
