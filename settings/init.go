package settings

import (
	"strings"

	"github.com/wingcd/go-xlsx-protobuf/model"
)

var (
	TOOL_VERSION      = "1.0"
	EXPORT_FOR_CLIENT = false
	PackageName       = "PBGen"

	DEFINES map[string]*model.DefineTableInfo
	ENUMS   []*model.DefineTableInfo
	STRUCTS []*model.DefineTableInfo
	TABLES  []*model.DataTable
)

func GetAllTables() []*model.DataTable {
	tables := make([]*model.DataTable, 0)
	for _, stru := range STRUCTS {
		table := model.Struct2Table(stru)
		if table != nil {
			tables = append(tables, table)
		}
	}

	PreProcessTable(tables)

	for _, table := range TABLES {
		tables = append(tables, table)
	}

	return tables
}

func SetDefines(defines map[string]*model.DefineTableInfo) {
	DEFINES = defines
	ENUMS = make([]*model.DefineTableInfo, 0)
	STRUCTS = make([]*model.DefineTableInfo, 0)
	for _, info := range defines {
		if info.Category == model.DEFINE_TYPE_ENUM {
			ENUMS = append(ENUMS, info)
		} else if info.Category == model.DEFINE_TYPE_STRUCT {
			STRUCTS = append(STRUCTS, info)
		}
	}
}

func SetTables(tables []*model.DataTable) {
	TABLES = make([]*model.DataTable, 0)
	for _, table := range tables {
		TABLES = append(TABLES, table)
	}

	PreProcessTable(tables)
}

func GetEnum(pbType string) *model.DefineTableInfo {
	if DEFINES == nil {
		return nil
	}
	if val, ok := DEFINES[pbType]; ok {
		return val
	}
	return nil
}

func IsEnum(pbType string) bool {
	if DEFINES == nil {
		return false
	}
	if val, ok := DEFINES[pbType]; ok {
		return val.Category == model.DEFINE_TYPE_ENUM && ok
	}
	return false
}

func IsStruct(pbType string) bool {
	if DEFINES == nil {
		return false
	}
	if val, ok := DEFINES[pbType]; ok {
		return val.Category == model.DEFINE_TYPE_STRUCT && ok
	}
	return false
}

func IsTable(pbType string) bool {
	if TABLES == nil {
		return false
	}

	for _, table := range TABLES {
		if table.TypeName == pbType {
			return true
		}
	}
	return false
}

var pbFieldEncodeTypes = map[string]string{
	"bool":   "varint",
	"int":    "varint",
	"int32":  "varint",
	"uint":   "varint",
	"uint32": "varint",
	"int64":  "varint",
	"uint64": "varint",
	"float":  "fixed32",
	"double": "fixed64",
	"string": "bytes",
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

func PreProcessStruct(structs []*model.StructInfo) {
	for _, st := range structs {
		st.EncodeType, st.IsEnum, st.IsStruct = GetEncodeType(st.RawValueType)
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
