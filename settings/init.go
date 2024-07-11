package settings

import (
	"log"
	"sort"
	"strings"

	"github.com/wingcd/go-xlsx-exporter/model"
)

const (
	EXPORT_TYPE_ALL    = 1
	EXPORT_TYPE_CLIENT = 2
	EXPORT_TYPE_SERVER = 3
	EXPORT_TYPE_IGNORE = 4
)

var (
	GO_PROTO_VERTION = "v1.27.1"
	TOOL_VERSION     = "1.4"

	ExportType     = EXPORT_TYPE_ALL
	PackageName    = "PBGen"
	Channel        = ""
	PbBytesFileExt = ".bytes"
	CommentSymbol  = "#"
	ArraySplitChar = ","
	StrictMode     = false

	// just for debug
	GenLanguageType = false

	Rules []*RuleInfo

	Imports []string

	DEFINES     map[string]*model.DefineTableInfo
	ENUMS       []*model.DefineTableInfo
	STRUCTS     []*model.DefineTableInfo
	CONSTS      []*model.DefineTableInfo
	TABLES      []*model.DataTable
	LANG_TABLES []*model.DataTable
)

func GetAllTables() []*model.DataTable {
	tables := make([]*model.DataTable, 0)
	for _, stru := range STRUCTS {
		table := model.Struct2Table(stru)
		if table != nil {
			tables = append(tables, table)
		}
	}

	for _, cost := range CONSTS {
		table := model.Struct2Table(cost)
		if table != nil {
			tables = append(tables, table)
		}
	}

	for _, table := range TABLES {
		tables = append(tables, table)
	}

	return CombineTables(tables)
}

func SetDefines(defines ...map[string]*model.DefineTableInfo) {
	DEFINES = make(map[string]*model.DefineTableInfo)
	for _, ds := range defines {
		for _, d := range ds {
			DEFINES[d.TypeName] = d
		}
	}

	//合并定义项目
	defs := make(map[string]*model.DefineTableInfo)
	for _, df := range DEFINES {
		if define, ok := defs[df.TypeName]; ok {
			define.Items = append(define.Items, df.Items...)
		} else {
			defs[df.TypeName] = df
		}
	}

	ENUMS = make([]*model.DefineTableInfo, 0)
	STRUCTS = make([]*model.DefineTableInfo, 0)
	CONSTS = make([]*model.DefineTableInfo, 0)

	for _, info := range DEFINES {
		if info.Category == model.DEFINE_TYPE_ENUM {
			ENUMS = append(ENUMS, info)
		} else if info.Category == model.DEFINE_TYPE_STRUCT {
			STRUCTS = append(STRUCTS, info)
		} else if info.Category == model.DEFINE_TYPE_CONST {
			CONSTS = append(CONSTS, info)
		}
	}

	sort.Sort(model.DefineTableInfos(ENUMS))
	sort.Sort(model.DefineTableInfos(STRUCTS))
	sort.Sort(model.DefineTableInfos(CONSTS))
}

func AddLanguageTable() *model.DataTable {
	var table = model.DataTable{}
	table.TypeName = "Language"
	table.TableType = model.ETableType_Language
	table.DefinedTable = ""
	table.Headers = make([]*model.DataTableHeader, 0)

	var hId = model.DataTableHeader{}
	hId.FieldName = "ID"
	hId.RawValueType = "string"
	hId.ValueType = "string"
	hId.StandardValueType = hId.ValueType
	hId.PBValueType = hId.ValueType
	hId.Index = 1
	hId.TitleFieldName = "ID"
	hId.ExportClient = true
	hId.ExportServer = true
	table.Headers = append(table.Headers, &hId)

	var hVal = model.DataTableHeader{}
	hVal.FieldName = "Text"
	hVal.RawValueType = "string"
	hVal.ValueType = "string"
	hVal.StandardValueType = hVal.ValueType
	hVal.PBValueType = hVal.ValueType
	hVal.Index = 2
	hVal.TitleFieldName = "Text"
	hVal.ExportClient = true
	hVal.ExportServer = true
	table.Headers = append(table.Headers, &hVal)

	TABLES = append([]*model.DataTable{&table}, TABLES...)

	return &table
}

func SetTables(tables ...*model.DataTable) {
	TABLES = make([]*model.DataTable, 0)
	LANG_TABLES = make([]*model.DataTable, 0)

	for _, table := range tables {
		if table.TableType == model.ETableType_Language {
			table.TypeName = "Language"
			LANG_TABLES = append(LANG_TABLES, table)
		} else {
			TABLES = append(TABLES, table)
		}
	}

	if LANG_TABLES != nil && len(LANG_TABLES) > 0 {
		table := AddLanguageTable()
		sourceFiles := make([]string, 0)
		for _, t := range LANG_TABLES {
			sourceFiles = append(sourceFiles, t.DefinedTable)
		}
		table.DefinedTable = strings.Join(sourceFiles, ";")
	}

	TABLES = CombineTables(TABLES)

	sort.Sort(model.DataTables(TABLES))
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

func CombineTables(tables []*model.DataTable) []*model.DataTable {
	var ts = make(map[string]*model.DataTable, 0)
	var newTables = make([]*model.DataTable, 0)
	for _, t := range tables {
		if table, ok := ts[t.TypeName]; ok {
			table.Data = append(table.Data, t.Data...)
		} else {
			ts[t.TypeName] = t
			newTables = append(newTables, t)
		}
	}
	return newTables
}

func CheckRule(id int, value string) bool {
	var rule *RuleInfo
	for _, r := range Rules {
		if r.ID == id {
			rule = r
			break
		}
	}
	if rule == nil || rule.RRule == nil {
		log.Fatalf("未找到规则Id=%v", id)
		return true
	}

	if rule.Disable {
		return true
	}

	return rule.RRule.MatchString(value)
}
