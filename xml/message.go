package xml

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/utils"
	"github.com/wingcd/go-xlsx-exporter/xlsx"
)

type DefineField struct {
	Field xml.Name `xml:"field"`
	Name  string   `xml:"name,attr"`
	Value string   `xml:"value,attr"`
	Type  string   `xml:"type,attr"`
	Desc  string   `xml:"desc,attr"`
}

type Enum struct {
	Enum  xml.Name      `xml:"enum"`
	Name  string        `xml:"name,attr"`
	Field []DefineField `xml:"field"`
}

type Define struct {
	Define xml.Name `xml:"define"`
	Enum   []Enum   `xml:"enum"`
}

type Field struct {
	Field xml.Name `xml:"field"`
	Name  string   `xml:"name,attr"`
	Type  string   `xml:"type,attr"`
	Desc  string   `xml:"desc,attr"`
}

type Message struct {
	XMLName xml.Name `xml:"message"`
	Name    string   `xml:"name,attr"`
	Id      int      `xml:"id,attr"`
	Field   []Field  `xml:"field"`
}

type Proto struct {
	XMLName xml.Name  `xml:"proto"`
	Message []Message `xml:"message"`
}

func ParseDefine(files ...string) (infos map[string]*model.DefineTableInfo) {
	var size = len(files)
	if size == 0 {
		log.Print("[错误] 参数错误 \n")
		return
	}

	infos = make(map[string]*model.DefineTableInfo, 0)
	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("error:%v", err)
			return
		}

		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("error:%v", err)
			return
		}

		v := Define{}
		err = xml.Unmarshal(data, &v)
		if err != nil {
			fmt.Println("error:%v", err)
			return
		}

		for _, e := range v.Enum {
			info := new(model.DefineTableInfo)
			info.TypeName = e.Name
			info.Category = model.DEFINE_TYPE_ENUM
			infos[info.TypeName] = info

			for _, field := range e.Field {
				item := new(model.DefineTableItem)
				info.Items = append(info.Items, item)

				item.FieldName = field.Name
				item.TitleFieldName = strings.Title(item.FieldName)
				item.Value = field.Value
				item.Desc = field.Desc
				item.RawValueType = field.Value
				finfo := utils.CompileValueType(item.RawValueType)

				item.IsArray = finfo.IsArray
				item.ValueType = finfo.ValueType
				item.ArraySplitChar = finfo.SplitChar
				item.Convertable = finfo.Convertable
				item.Cachable = finfo.Cachable
				item.IsVoid = finfo.IsVoid
				item.Alias = finfo.Alias
				item.Rule = finfo.Rule

				item.ValueType = utils.ConvertToStandardType(item.ValueType)
				item.StandardValueType = item.ValueType
				_, item.PBValueType = utils.ToPBType(item.StandardValueType)
			}
		}
	}

	xlsx.DefinesPreProcess(infos)

	return
}

func Parse(files ...string) (tables []*model.DataTable) {
	var size = len(files)
	if size == 0 {
		log.Print("[错误] 参数错误 \n")
		return
	}

	tables = make([]*model.DataTable, 0)
	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("error:%v", err)
			return
		}

		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("error:%v", err)
			return
		}

		v := Proto{}
		err = xml.Unmarshal(data, &v)
		if err != nil {
			fmt.Println("error:%v", err)
			return
		}

		for _, msg := range v.Message {
			table := new(model.DataTable)
			table.Id = msg.Id
			table.TypeName = msg.Name
			table.DefinedTable = filename
			table.Data = make([][]string, 0)
			table.IsArray = false
			table.TableType = model.ETableType_Message
			table.Headers = make([]*model.DataTableHeader, 0)
			for _, field := range msg.Field {
				header := new(model.DataTableHeader)
				table.Headers = append(table.Headers, header)

				header.FieldName = field.Name
				header.TitleFieldName = strings.Title(header.FieldName)

				header.Desc = field.Desc
				header.ExportClient = true
				header.ExportServer = true

				header.RawValueType = field.Type
				finfo := utils.CompileValueType(header.RawValueType)
				if !finfo.Valiable {
					log.Fatalf("[错误] 字段定义错误 文件：%s 类型：%s 列：%s 描述：%s\n", table.DefinedTable, table.TypeName, header.FieldName, header.Desc)
				}
				header.IsArray = finfo.IsArray
				header.ValueType = finfo.ValueType
				header.ArraySplitChar = finfo.SplitChar
				header.Convertable = finfo.Convertable
				header.Cachable = finfo.Cachable
				header.IsVoid = finfo.IsVoid
				header.Alias = finfo.Alias
				header.Rule = finfo.Rule
				header.Index = len(table.Headers)

				header.ValueType = utils.ConvertToStandardType(header.ValueType)
				header.StandardValueType = header.ValueType
				_, header.PBValueType = utils.ToPBType(header.StandardValueType)
			}

			tables = append(tables, table)
		}
	}
	return tables
}
