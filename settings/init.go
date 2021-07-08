package settings

import (
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
