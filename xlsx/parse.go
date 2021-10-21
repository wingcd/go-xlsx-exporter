package xlsx

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
)

const (
	totalSize = 6
)

// @params filename 文件名,表格名，文件名，表格名...
func ParseDefineSheet(files ...string) (infos map[string]*model.DefineTableInfo) {
	var size = len(files)
	var cnt = size / 2
	if size == 0 || size%2 != 0 {
		log.Print("[错误] 参数错误 \n")
		return
	}

	infos = make(map[string]*model.DefineTableInfo, 0)

	rows := make([][]string, 0)
	for i := 0; i < cnt; i++ {
		filename := files[i*2]
		sheet := files[i*2+1]

		fmt.Printf("parse file %s:%s...\n", filename, sheet)
		f, err := excelize.OpenFile(filename)
		if err != nil {
			fmt.Println(err)
			return
		}

		rs, err := f.GetRows(sheet)
		if err != nil {
			fmt.Println(err)
			return
		}

		rows = append(rows, rs...)
	}

	for ri, row := range rows {
		// 种类	类型名	字段名	值	说明
		if ri == 0 {
			continue
		}

		if row == nil {
			// log.Printf("[警告] 有空定义行 表：%v 第%v行 \n", fmt.Sprintf("%s:%s", files[0], files[1]), ri+1)
			continue
		}

		// 过滤注释
		if utils.IsComment(row[0]) {
			continue
		}

		// 空行过滤
		var fixedColCount = 3
		var catgory = strings.ToLower(row[model.DEFINE_CATEGORY_NAME_INDEX])
		if catgory != model.DEFINE_TYPE_ENUM {
			fixedColCount = 4
		}
		for ci := 0; ci < fixedColCount; ci++ {
			if row[ci] == "" {
				// log.Printf("[警告] 有空定义行 表：%v 第%v行 \n", fmt.Sprintf("%s:%s", files[0], files[1]), ri+1)
				continue
			}
		}

		var typename = row[model.DEFINE_COLUMN_NAME_INDEX]
		var info *model.DefineTableInfo
		if info, ok := infos[typename]; !ok {
			info = new(model.DefineTableInfo)
			info.DefinedTable = ""
			for i := 0; i < cnt; i++ {
				info.DefinedTable += fmt.Sprintf("%s:%s;", files[i*2], files[i*2+1])
			}

			infos[typename] = info

			info.Category = catgory
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

		item.ValueType = utils.ConvertToStandardType(item.ValueType)
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
					v, err := strconv.Atoi(item.Value)
					if err != nil {
						log.Fatalf("[错误] 枚举值类型错误 表：%s 列：%s \n", info.DefinedTable, item.Desc)
					}
					value = v
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
		} else if info.Category == model.DEFINE_TYPE_CONST {
			for i, item := range info.Items {
				item.Index = i + 1
			}
		}
	}

	return
}

// @params filename 文件名,表格名，文件名，表格名...
func ParseDataSheet(files ...string) (table *model.DataTable) {
	var size = len(files)
	var cnt = size / 2
	if size == 0 || size%2 != 0 {
		log.Print("[错误] 参数错误 \n")
		return
	}

	cols := make([][]string, 0)
	rows := make([][]string, 0)

	for i := 0; i < cnt; i++ {
		filename := files[i*2]
		sheet := files[i*2+1]

		fmt.Printf("parse file %s:%s...\n", filename, sheet)

		f, err := excelize.OpenFile(filename)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 只需要第一张表的列数据
		if i == 0 {
			cls, err := f.GetCols(sheet)
			if err != nil {
				fmt.Println(err)
				return
			}

			// 过滤数据项（列）,不管前面有多少注释，过滤后的前四行必须按规则编写
			filterCols := make([][]string, 0)
			for ci, col := range cls {
				// 索引列不能为空，否则过滤掉
				var emptyIndex = col == nil
				if emptyIndex {
					if emptyIndex {
						log.Printf("[警告] 有空数据列 表：%v-%v 第%v行 \n", filename, sheet, ci+1)
					}
					continue
				}
				filterCols = append(filterCols, col)
			}
			cols = append(cols, filterCols...)
		}

		// 获取所有行
		rs, err := f.GetRows(sheet)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 过滤数据项（列）,不管前面有多少注释，过滤后的前四行必须按规则编写
		filterRows := make([][]string, 0)
		for ri, row := range rs {
			// 索引列不能为空，否则过滤掉
			var emptyIndex = row == nil
			if emptyIndex || utils.IsComment(row[0]) {
				if emptyIndex {
					log.Printf("[警告] 有空数据行 表：%v-%v 第%v行 \n", filename, sheet, ri+1)
				}
				continue
			}
			filterRows = append(filterRows, row)
		}

		rows = append(rows, filterRows[4:]...)
	}

	table = new(model.DataTable)
	table.DefinedTable = ""
	for i := 0; i < cnt; i++ {
		table.DefinedTable += fmt.Sprintf("%s:%s;", files[i*2], files[i*2+1])
	}

	table.Headers = make([]*model.DataTableHeader, 0)
	table.IsDataTable = true

	ignoreRows := make(map[int]bool, 0)
	ignoreCols := make(map[int]bool, 0)

	// 过滤标准行前的注释项
	filterCols := make([][]string, 0)
	for ci, col := range cols {
		if utils.IsComment(col[0]) {
			ignoreCols[ci] = true
		}
		if ci == 0 {
			for ri, cellValue := range col {
				if utils.IsComment(cellValue) {
					ignoreRows[ri] = true
				}
			}
		}

		// 不直接在此处过滤列，否则过滤数组索引与行数据不一致
		newCol := make([]string, 0)
		for cii, cellValue := range col {
			if _, ignoreRow := ignoreRows[cii]; !ignoreRow {
				newCol = append(newCol, cellValue)
			}

			// 找到前4个非注释行就退出
			if len(newCol) == 4 {
				break
			}
		}
		filterCols = append(filterCols, newCol)
	}
	cols = filterCols

	var firstColIndex = -1
	for ci, col := range cols {
		if _, ignore := ignoreCols[ci]; ignore {
			continue
		}

		// 处理空列
		if col[model.DATA_ROW_FIELD_INDEX] == "" || col[model.DATA_ROW_TYPE_INDEX] == "" {
			log.Fatalf("[错误] 数据类型或字段名不能为空 表：%v 第%v列 \n", table.DefinedTable, ci+1)
			return
		}

		header := new(model.DataTableHeader)
		cs := strings.ToLower(col[model.DATA_ROW_CS_INDEX])
		header.ExportClient = strings.Contains(cs, "c")
		header.ExportServer = strings.Contains(cs, "s")
		if cs == "" {
			header.ExportClient = true
			header.ExportServer = true
		}
		header.Desc = col[model.DATA_ROW_DESC_INDEX]

		ignore := false
		if settings.ExportType != settings.EXPORT_TYPE_IGNORE {
			if settings.EXPORT_TYPE_CLIENT == settings.ExportType && !header.ExportClient {
				// 当为客户端导出，但此列不支持客户端时，过滤掉
				ignore = true
			} else if settings.EXPORT_TYPE_SERVER == settings.ExportType && !header.ExportServer {
				// 当为后端导出，但此列不支持后端时，过滤掉
				ignore = true
			}
		}
		if utils.IsComment(header.Desc) {
			ignore = true
		}

		if !ignore {
			header.FieldName = col[model.DATA_ROW_FIELD_INDEX]
			header.TitleFieldName = strings.Title(header.FieldName)
			header.RawValueType = col[model.DATA_ROW_TYPE_INDEX]
			header.IsArray = utils.IsArray(header.RawValueType)
			header.ValueType = utils.GetBaseType(header.RawValueType)
			header.Index = len(table.Headers) + 1

			header.ValueType = utils.ConvertToStandardType(header.ValueType)
			table.Headers = append(table.Headers, header)

			if firstColIndex < 0 {
				firstColIndex = ci
			}
		} else {
			ignoreCols[ci] = true
		}
	}

	// 过滤数据项（列）,不管前面有多少注释，过滤后的前四行必须按规则编写
	filterRows := make([][]string, 0)
	for ri, row := range rows {
		// 索引列不能为空，否则过滤掉
		var emptyIndex = row == nil
		if emptyIndex || utils.IsComment(row[0]) {
			if emptyIndex {
				log.Printf("[警告] 有空数据行 表：%v 第%v行 \n", table.DefinedTable, ri+1)
			}
			continue
		}
		filterRows = append(filterRows, row)
	}

	// 预处理数据
	// 1. 防止出现空数据
	// 2. 过滤注释数据列
	// 3. 过滤注释数据项（行）
	realHeadSize := len(table.Headers)
	for ri, row := range rows {
		if len(ignoreCols) > 0 {
			newRow := make([]string, 0)
			for ci, cellValue := range row {
				if _, ignore := ignoreCols[ci]; !ignore {
					// 判断是否需要过滤此条数据
					newRow = append(newRow, cellValue)
				}
			}

			// 补齐空数据
			if len(newRow) < realHeadSize {
				for i := 0; i < realHeadSize-len(newRow); i++ {
					newRow = append(newRow, "")
				}
			}

			rows[ri] = newRow
		}

		filterRows = append(filterRows, rows[ri])
	}

	table.Data = rows

	return
}
