package utils

import (
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
	repeated := false
	if strings.Contains(valueType, "[]") {
		repeated = true
	}
	var rawType = strings.Replace(valueType, "[]", "", -1)
	var isEnum = IsEnum(rawType)
	var isStruct = IsStruct(rawType) || IsTable(rawType)
	if repeated {
		return "bytes", isEnum, isStruct
	}
	if tp, ok := pbFieldEncodeTypes[rawType]; ok {
		return tp, isEnum, isStruct
	} else if isEnum {
		return "varint", isEnum, isStruct
	}
	return "", isEnum, isStruct
}

func PreProcessDefine(defines []*model.DefineTableInfo) {
	for _, d := range defines {
		for _, st := range d.Items {
			st.EncodeType, st.IsEnum, st.IsStruct = GetEncodeType(st.RawValueType)
		}
	}
}

func PreProcessHeader(header *model.DataTableHeader) {
	header.EncodeType, header.IsEnum, header.IsStruct = GetEncodeType(header.RawValueType)
}

func PreProcessTable(tables []*model.DataTable) {
	for _, table := range tables {
		for _, header := range table.Headers {
			PreProcessHeader(header)
		}
	}
}

func IsComment(value string) bool {
	if settings.CommentSymbol == "" {
		return false
	}
	return strings.Index(strings.Trim(value, " "), settings.CommentSymbol) == 0
}
