package xlsx

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
	"github.com/xuri/excelize/v2"
)

const (
	totalSize = 6
)

var (
	channelRegex *regexp.Regexp = nil
)

func init() {
	// [channel]xxx
	channelRegex = regexp.MustCompile(`\[(.+)\](.+)`)
}

func transpose(arr [][]string) [][]string {
	cnt1 := len(arr)
	if cnt1 == 0 {
		return arr
	}
	cnt2 := 0
	for i := 0; i < cnt1; i++ {
		cnt2 = int(math.Max(float64(len(arr[i])), float64(cnt2)))
	}

	var arr2 [][]string = make([][]string, cnt2)
	for i := 0; i < cnt2; i++ {
		arr2[i] = make([]string, cnt1)
	}

	//遍历数组并进行转置
	for i := 0; i < cnt1; i++ {
		for j := 0; j < len(arr[i]); j++ {
			arr2[j][i] = arr[i][j]
		}
	}

	return arr2
}

// @params filename 文件名,表格名，文件名，表格名...
func ParseDefineSheet(files ...*settings.SheetInfo) (infos map[string]*model.DefineTableInfo) {
	var cnt = len(files)
	if cnt == 0 {
		log.Print("[错误] 参数错误 \n")
		return
	}

	infos = make(map[string]*model.DefineTableInfo, 0)

	rows := make([][]string, 0)
	for i := 0; i < cnt; i++ {
		file := files[i]
		filename := file.File
		sheet := file.Sheet

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

		if file.Transpose {
			rs = transpose(rs)
		}

		if i == 0 {
			rows = append(rows, rs...)
		} else {
			// 第二张开始，不需要表头
			rows = append(rows, rs[1:]...)
		}
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
			info.StartID = int64(ri)
			info.DefinedTable = ""
			for i := 0; i < cnt; i++ {
				info.DefinedTable += fmt.Sprintf("%s:%s;", files[i].File, files[i].Sheet)
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
		finfo := utils.CompileValueType(item.RawValueType)
		if info.Category != model.DEFINE_TYPE_ENUM && !finfo.Valiable {
			log.Fatalf("[错误] 字段定义错误 表：%s 类型：%s 列：%s 描述：%s\n", info.DefinedTable, info.TypeName, item.FieldName, item.Desc)
		}

		item.IsArray = finfo.IsArray
		item.ValueType = finfo.ValueType
		item.ArraySplitChar = finfo.SplitChar
		item.Convertable = finfo.Convertable
		item.Cachable = finfo.Cachable
		item.IsVoid = finfo.IsVoid
		item.Alias = finfo.Alias
		item.Rule = finfo.Rule
		info.Items = append(info.Items, item)

		item.ValueType = utils.ConvertToStandardType(item.ValueType)
		item.StandardValueType = item.ValueType
		_, item.PBValueType = utils.ToPBType(item.StandardValueType)
	}

	DefinesPreProcess(infos)

	return
}

func DefinesPreProcess(infos map[string]*model.DefineTableInfo) {
	// 预处理
	preValue := -1
	for _, info := range infos {
		if info.Category == model.DEFINE_TYPE_ENUM {
			names := make(map[string]int)
			for _, item := range info.Items {
				if _, ok := names[item.FieldName]; ok {
					log.Fatalf("[错误] 重复定义属性 表：%s 类型：%s 列：%s 描述：%s\n", info.DefinedTable, info.TypeName, item.FieldName, item.Desc)
				}
				names[item.FieldName] = 0

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
				item.StandardValueType = ""
				item.PBValueType = ""
				item.Value = "0"
				item.Desc = ""
				item.IsArray = false
				item.Index = 0

				info.Items = append([]*model.DefineTableItem{item}, info.Items...)
			}

			sort.Sort(model.DefineTableItems(info.Items))
		} else if info.Category == model.DEFINE_TYPE_STRUCT {
			var idx = 0
			names := make(map[string]int)
			for _, item := range info.Items {
				if _, ok := names[item.FieldName]; ok {
					log.Fatalf("[错误] 重复定义属性 表：%s 类型：%s 列：%s 描述：%s\n", info.DefinedTable, info.TypeName, item.FieldName, item.Desc)
				}
				names[item.FieldName] = 0

				if !item.IsVoid {
					idx++
				}
				item.Index = idx
			}
		} else if info.Category == model.DEFINE_TYPE_CONST {
			var newItems = make([]*model.DefineTableItem, 0)
			var channelItems = make([]*model.DefineTableItem, 0)
			// 将channel变量值覆盖到同名的变量上
			for _, item := range info.Items {
				if channelRegex.MatchString(item.FieldName) {
					channelItems = append(channelItems, item)
				} else {
					newItems = append(newItems, item)
				}
			}

			if settings.Channel != "" && len(channelItems) > 0 {
				for _, item := range channelItems {
					matches := channelRegex.FindStringSubmatch(item.FieldName)
					if len(matches) == 3 {
						channelValue := matches[1]
						filedName := matches[2]
						if settings.Channel == channelValue {
							for _, item2 := range newItems {
								if item2.FieldName == filedName {
									item2.Value = item.Value
									break
								}
							}
						}
					}
				}
			}
			info.Items = newItems

			var idx = 0
			names := make(map[string]int)
			for _, item := range info.Items {
				if _, ok := names[item.FieldName]; ok {
					log.Fatalf("[错误] 重复定义属性 表：%s 类型：%s 列：%s 描述：%s\n", info.DefinedTable, info.TypeName, item.FieldName, item.Desc)
				}
				names[item.FieldName] = 0

				if !item.IsVoid {
					idx++
				}
				item.Index = idx
			}
		}
	}
}

// @params filename 文件名,表格名，文件名，表格名...
func ParseDataSheet(files ...*settings.SheetInfo) (table *model.DataTable) {
	var cnt = len(files)
	if cnt == 0 {
		log.Print("[错误] 参数错误 \n")
		return
	}

	cols := make([][]string, 0)
	rows := make([][]string, 0)

	for i := 0; i < cnt; i++ {
		file := files[i]
		filename := file.File
		sheet := file.Sheet

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

			if file.Transpose {
				cls = transpose(cls)
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

		if file.Transpose {
			rs = transpose(rs)
		}

		// 过滤数据项（列）,不管前面有多少注释，过滤后的前四行必须按规则编写
		filterRows := make([][]string, 0)
		for ri, row := range rs {
			// 索引列不能为空，否则过滤掉
			var emptyIndex = false
			var emptyRow = len(row) == 0
			var comment = !emptyRow && utils.IsComment(row[0])
			if !comment && !emptyRow {
				emptyIndex = row == nil || len(row) > 0 && row[0] == ""
				if !emptyIndex && len(row) > 0 {
					emptyRow = true
					for i := 0; i < len(row); i++ {
						if row[i] != "" {
							emptyRow = false
							break
						}
					}
				}
			}

			if emptyIndex || emptyRow || comment {
				if emptyIndex {
					log.Printf("[警告] 有空索引 表：%v-%v 第%v行 \n", filename, sheet, ri+1)
				} else if emptyRow {
					log.Printf("[警告] 有空数据行 表：%v-%v 第%v行 \n", filename, sheet, ri+1)
				}
				continue
			}
			filterRows = append(filterRows, row)
		}

		rows = append(rows, filterRows[model.FixedRowCount:]...)
	}

	table = new(model.DataTable)
	table.DefinedTable = ""
	for i := 0; i < cnt; i++ {
		table.DefinedTable += fmt.Sprintf("%s:%s;", files[i].File, files[i].Sheet)
	}

	table.Headers = make([]*model.DataTableHeader, 0)
	table.TableType = model.ETableType_Data
	table.NeedAddItems = true

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
			if len(newCol) == model.FixedRowCount {
				break
			}
		}
		filterCols = append(filterCols, newCol)
	}
	cols = filterCols

	var firstColIndex = -1
	var idx = 0
	names := make(map[string]int)
	for ci, col := range cols {
		if _, ignore := ignoreCols[ci]; ignore {
			continue
		}

		// 处理空列
		if col[model.DATA_ROW_FIELD_INDEX] == "" || col[model.DATA_ROW_TYPE_INDEX] == "" {
			if settings.StrictMode {
				log.Fatalf("[错误] 数据类型或字段名不能为空 表：%v 第%v列 \n", table.DefinedTable, ci+1)
				return
			} else {
				log.Printf("[警告] 数据类型或字段名不能为空,将跳过此列 表：%v 第%v列 \n", table.DefinedTable, ci+1)
				continue
			}
		}

		// 处理非渠道数据
		fileName := col[model.DATA_ROW_FIELD_INDEX]
		hasChannel := channelRegex.MatchString(fileName)
		if hasChannel && (settings.Channel == "" || settings.Channel != "" && !strings.Contains(fileName, settings.Channel)) {
			continue
		}

		header := new(model.DataTableHeader)
		cs := ""
		if settings.ExportType != settings.EXPORT_TYPE_ALL {
			cs = strings.ToLower(col[model.DATA_ROW_CS_INDEX])
			header.ExportClient = strings.Contains(cs, "c")
			header.ExportServer = strings.Contains(cs, "s")
		}
		if cs == "" {
			header.ExportClient = true
			header.ExportServer = true
		}
		header.Desc = col[model.DATA_ROW_DESC_INDEX]

		ignore := false
		if settings.ExportType != settings.EXPORT_TYPE_ALL {
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
			header.RawValueType = col[model.DATA_ROW_TYPE_INDEX]
			finfo := utils.CompileValueType(header.RawValueType)
			if !settings.StrictMode && !finfo.Valiable {
				ignore = true
				continue
			}

			if !finfo.Valiable {
				log.Fatalf("[错误] 字段定义错误 表：%s 类型：%s 列：%s 描述：%s\n", table.DefinedTable, table.TypeName, header.FieldName, header.Desc)
			}

			if !finfo.IsVoid {
				idx++
			} else {
				// 数据需要过滤此列
				ignoreCols[ci] = true
			}

			header.FieldName = col[model.DATA_ROW_FIELD_INDEX]
			if _, ok := names[header.FieldName]; ok {
				log.Fatalf("[错误] 重复定义属性 表：%s 类型：%s 列：%s 描述：%s\n", table.DefinedTable, table.TypeName, header.FieldName, header.Desc)
			}
			names[header.FieldName] = 0

			header.TitleFieldName = strings.Title(header.FieldName)
			header.IsArray = finfo.IsArray
			header.ValueType = finfo.ValueType
			header.ArraySplitChar = finfo.SplitChar
			header.Convertable = finfo.Convertable
			header.Cachable = finfo.Cachable
			header.IsVoid = finfo.IsVoid
			header.Alias = finfo.Alias
			header.Rule = finfo.Rule
			header.Index = idx

			header.ValueType = utils.ConvertToStandardType(header.ValueType)
			header.StandardValueType = header.ValueType
			_, header.PBValueType = utils.ToPBType(header.StandardValueType)
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

			rows[ri] = newRow
		}

		// 补齐空数据
		newRow := rows[ri]
		if len(newRow) < realHeadSize {
			for realHeadSize-len(newRow) > 0 {
				newRow = append(newRow, "")
			}
			rows[ri] = newRow
		}

		filterRows = append(filterRows, rows[ri])
	}

	table.Data = rows

	DataSheetPreProcess(table)
	return
}

func DataSheetPreProcess(table *model.DataTable) {
	var newHeaders = make([]*model.DataTableHeader, 0)
	var channelHeaders = make([]*model.DataTableHeader, 0)
	var data = table.Data
	// 将channel变量值覆盖到同名的变量上
	for _, header := range table.Headers {
		if channelRegex.MatchString(header.FieldName) {
			channelHeaders = append(channelHeaders, header)
		} else {
			newHeaders = append(newHeaders, header)
		}
	}

	if len(channelHeaders) > 0 {
		var newData = make([][]string, 0)
		var deleteIndices = make(map[int]bool, 0)
		for _, item := range channelHeaders {
			matches := channelRegex.FindStringSubmatch(item.FieldName)
			if len(matches) == 3 {
				deleteIndices[item.Index-1] = true
				channelValue := matches[1]
				filedName := matches[2]
				if settings.Channel == channelValue {
					for _, item2 := range newHeaders {
						if item2.FieldName == filedName {
							for i, row := range data {
								row[item2.Index-1] = row[item.Index-1]
								data[i] = row
							}
							break
						}
					}
				}
			}
		}

		for _, data := range table.Data {
			var newRow = make([]string, 0)
			for ci, cellValue := range data {
				if _, ignore := deleteIndices[ci]; !ignore {
					newRow = append(newRow, cellValue)
				}
			}
			newData = append(newData, newRow)
		}
		table.Data = newData
	}
	table.Headers = newHeaders
}

func CheckTable(table *model.DataTable) {
	for i, row := range table.Data {
		for j, col := range row {
			if j >= len(table.Headers) {
				if settings.StrictMode {
					log.Printf("[警告] 数据列数超出定义 表：%s 类型：%s 行：%v 列：%v 数据：%v\n", table.DefinedTable, table.TypeName, i+1, j+1, col)
				}
				continue
			}

			header := table.Headers[j]
			if header.Rule > 0 {
				if !settings.CheckRule(header.Rule, col) {
					if settings.StrictMode {
						log.Fatalf("[错误] 不满足规则[%v] 表：%s 类型：%s 数据行:%v 列：%s 描述：%s\n", header.Rule, table.DefinedTable, table.TypeName, i, header.FieldName, header.Desc)
					} else {
						log.Printf("[错误] 不满足规则[%v] 表：%s 类型：%s 数据行:%v 列：%s 描述：%s\n", header.Rule, table.DefinedTable, table.TypeName, i, header.FieldName, header.Desc)
					}
				}
			}
		}
	}
}

func CheckDefine(info *model.DefineTableInfo) {
	for _, item := range info.Items {
		if item.Rule > 0 && !settings.CheckRule(item.Rule, item.Value) {
			if settings.StrictMode {
				log.Fatalf("[错误] 不满足规则[%v] 表：%s 类型：%s 列：%s 描述：%s\n", item.Rule, info.DefinedTable, info.TypeName, item.FieldName, item.Desc)
			} else {
				log.Printf("[错误] 不满足规则[%v] 表：%s 类型：%s 列：%s 描述：%s\n", item.Rule, info.DefinedTable, info.TypeName, item.FieldName, item.Desc)
			}
		}
	}
}
