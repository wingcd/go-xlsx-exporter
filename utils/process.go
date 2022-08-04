package utils

import (
	"strconv"
	"strings"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
)

func IsEnum(pbType string) bool {
	if settings.DEFINES == nil {
		return false
	}
	if val, ok := settings.DEFINES[pbType]; ok {
		return val.Category == model.DEFINE_TYPE_ENUM && ok
	}
	return false
}

func GetEnum(pbType string) *model.DefineTableInfo {
	if settings.DEFINES == nil {
		return nil
	}
	if val, ok := settings.DEFINES[pbType]; ok {
		return val
	}
	return nil
}

func GetEnumDefault(pbType string) *model.DefineTableItem {
	if settings.DEFINES == nil {
		return nil
	}
	if val, ok := settings.DEFINES[pbType]; ok {
		return val.Items[0]
	}
	return nil
}

func GetEnumValues(pbType string) []int {
	if settings.DEFINES == nil {
		return []int{}
	}
	if val, ok := settings.DEFINES[pbType]; ok {
		if ok {
			values := make([]int, 0)
			for _, item := range val.Items {
				val, _ := strconv.Atoi(item.Value)
				values = append(values, val)
			}
			return values
		}
	}
	return []int{}
}

func GetEnumNames(pbType string) []string {
	if settings.DEFINES == nil {
		return []string{}
	}
	if val, ok := settings.DEFINES[pbType]; ok {
		if ok {
			values := make([]string, 0)
			for _, item := range val.Items {
				values = append(values, item.FieldName)
			}
			return values
		}
	}
	return []string{}
}

func IsStruct(pbType string) bool {
	if settings.DEFINES == nil {
		return false
	}
	if val, ok := settings.DEFINES[pbType]; ok {
		return val.Category == model.DEFINE_TYPE_STRUCT && ok
	}
	return false
}

func IsTable(pbType string) bool {
	if settings.TABLES == nil {
		return false
	}

	for _, table := range settings.TABLES {
		if table.TypeName == pbType {
			return true
		}
	}
	return false
}

var pbFieldEncodeTypes = map[string]string{
	"bool":    "varint",
	"int":     "varint",
	"int32":   "varint",
	"uint":    "varint",
	"uint32":  "varint",
	"int64":   "varint",
	"uint64":  "varint",
	"float":   "float",
	"float32": "float",
	"double":  "double",
	"float64": "double",
	"string":  "bytes",
}

// 获取编码类型
// 返回值： 编码类型，是否枚举, 是否结构体
func GetEncodeType(valueType string) (string, bool, bool) {
	valueType = strings.Replace(valueType, " ", "", -1)
	finfo := CompileValueType(valueType)

	var isEnum = IsEnum(finfo.ValueType)
	var isStruct = IsStruct(finfo.ValueType) || IsTable(finfo.ValueType)
	if finfo.IsArray {
		return "bytes", isEnum, isStruct
	}
	if tp, ok := pbFieldEncodeTypes[finfo.ValueType]; ok {
		return tp, isEnum, isStruct
	} else if isEnum {
		return "varint", isEnum, isStruct
	}
	return "", isEnum, isStruct
}

func PreProcessDefine(define *model.DefineTableInfo) {
	if define == nil {
		return
	}
	
	for _, st := range define.Items {
		st.EncodeType, st.IsEnum, st.IsStruct = GetEncodeType(st.RawValueType)
	}
}

func PreProcessDefines(defines []*model.DefineTableInfo) {
	for _, d := range defines {
		PreProcessDefine(d)
	}
}

func PreProcessHeader(header *model.DataTableHeader) {
	if header == nil {
		return
	}
	
	header.EncodeType, header.IsEnum, header.IsStruct = GetEncodeType(header.RawValueType)
}

func PreProcessTable(table *model.DataTable) {
	if table == nil {
		return
	}

	for _, header := range table.Headers {
		PreProcessHeader(header)
	}
}

func PreProcessTables(tables []*model.DataTable) {
	for _, table := range tables {
		PreProcessTable(table)
	}
}

func IsVoid(value string) bool {
	value = strings.ToLower(value) 
	return value == "void" || value == "void!"
}

func IsComment(value string) bool {
	if settings.CommentSymbol == "" {
		return false
	}
	return strings.Index(strings.Trim(value, " "), settings.CommentSymbol) == 0
}
