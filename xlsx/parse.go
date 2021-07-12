package xlsx

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/wingcd/go-xlsx-protobuf/model"
	"github.com/wingcd/go-xlsx-protobuf/utils"
)

const (
	totalSize = 6
)

func ParseDefineSheet(filename, sheet string) (infos map[string]*model.DefineTableInfo) {
	infos = make(map[string]*model.DefineTableInfo, 0)

	fmt.Printf("parse file %s...\n", filename)

	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println(err)
		return
	}

	for ri, row := range rows {
		// 种类	类型名	字段名	值	说明
		if ri == 0 {
			continue
		}

		var typename = row[model.DEFINE_COLUMN_NAME_INDEX]
		var info *model.DefineTableInfo
		if info, ok := infos[typename]; !ok {
			info = new(model.DefineTableInfo)
			info.DefinedTable = fmt.Sprintf("%s:%s", filename, sheet)
			infos[typename] = info

			info.Category = row[model.DEFINE_CATEGORY_NAME_INDEX]
			info.TypeName = row[model.DEFINE_COLUMN_NAME_INDEX]
			info.Items = make([]*model.DefineTableItem, 0)
		}
		info = infos[typename]

		rowSize := len(row)
		for i := 0; i < model.DEFINE_COLUMN_COUNT-rowSize; i++ {
			row = append(row, "")
		}

		item := new(model.DefineTableItem)
		item.FieldName = row[model.DEFINE_COLUMN_FIELD_INDEX]
		item.TitleFieldName = strings.Title(item.FieldName)
		item.Value = row[model.DEFINE_COLUMN_VALUE_INDEX]
		item.Desc = row[model.DEFINE_COLUMN_COMMENT_INDEX]
		item.RawValueType = row[model.DEFINE_COLUMN_TYPE_INDEX]
		item.IsArray = utils.IsArray(item.RawValueType)
		item.ValueType = utils.GetBaseType(item.RawValueType)
		info.Items = append(info.Items, item)
	}

	// 预处理
	preValue := -1
	for _, info := range infos {
		if info.Category == model.DEFINE_TYPE_ENUM {
			for _, item := range info.Items {
				value := -1
				if item.Value == "" {
					value = preValue + 1
				} else {
					value, err = strconv.Atoi(item.Value)
					if err != nil {
						log.Printf("[错误] 枚举值类型错误 类型：%s 列：%s \n", info.TypeName, item.Desc)
					}
				}
				item.Index = int(value)
				item.Value = strconv.FormatInt(int64(value), 10)
				preValue = value
			}

			// 添加0值
			if len(info.Items) > 0 && (info.Items[0].Value != "" && info.Items[0].Value != "0") {
				item := new(model.DefineTableItem)
				item.FieldName = "UNKNOWN"
				item.TitleFieldName = item.FieldName
				item.ValueType = ""
				item.Value = "0"
				item.Desc = ""
				item.IsArray = false
				item.Index = 0

				info.Items = append([]*model.DefineTableItem{item}, info.Items...)
			}

			sort.Sort(model.DefineTableItems(info.Items))
		} else if info.Category == model.DEFINE_TYPE_STRUCT {
			for i, item := range info.Items {
				item.Index = i + 1
			}
		}
	}

	return
}

func ParseDataSheet(filename, sheet string) (table *model.DataTable) {
	fmt.Printf("parse file %s...\n", filename)

	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	cols, err := f.GetCols(sheet)
	if err != nil {
		fmt.Println(err)
		return
	}

	table = new(model.DataTable)
	table.DefinedTable = fmt.Sprintf("%s:%s", filename, sheet)
	table.Headers = make([]*model.DataTableHeader, 0)
	for i, col := range cols {
		header := new(model.DataTableHeader)
		header.Desc = col[model.DATA_ROW_DESC_INDEX]
		cs := strings.ToLower(col[model.DATA_ROW_CS_INDEX])
		header.ExportClient = strings.Contains(cs, "c")
		header.ExportServer = strings.Contains(cs, "s")
		if !header.ExportClient && !header.ExportServer {
			header.ExportClient = true
			header.ExportServer = true
		}
		header.FieldName = col[model.DATA_ROW_FIELD_INDEX]
		header.TitleFieldName = strings.Title(header.FieldName)
		header.RawValueType = col[model.DATA_ROW_TYPE_INDEX]
		header.IsArray = utils.IsArray(header.RawValueType)
		header.ValueType = utils.GetBaseType(header.RawValueType)
		header.Index = i + 1
		table.Headers = append(table.Headers, header)
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows = rows[4:]

	// 预处理数据，防止出现空数据
	headSize := len(table.Headers)
	for ri, row := range rows {
		if len(row) < headSize {
			for i := 0; i < headSize-len(row); i++ {
				row = append(row, "")
			}
			rows[ri] = row
		}
	}
	table.Data = rows

	return
}
