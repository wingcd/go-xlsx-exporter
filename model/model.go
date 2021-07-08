package model

type DefineTableItem struct {
	FieldName    string
	RawValueType string
	ValueType    string
	Value        string
	Desc         string
	IsArray      bool
	Index        int
}

type DefineTableInfo struct {
	Category     string
	TypeName     string
	DefinedTable string
	Items        []*DefineTableItem
}

func (d *DefineTableInfo) IsValid(typeName string) bool {
	for _, item := range d.Items {
		if item.FieldName == typeName {
			return true
		}
	}
	return false
}

type DataTableHeader struct {
	Desc         string
	RawValueType string
	ValueType    string
	ExportClient bool
	ExportServer bool
	FieldName    string
	IsArray      bool
	Index        int
}

type DataTable struct {
	TypeName     string
	Headers      []*DataTableHeader
	DefinedTable string
	Data         [][]string
}

func Struct2Table(info *DefineTableInfo) *DataTable {
	if info.Category != DEFINE_TYPE_STRUCT {
		return nil
	}
	table := DataTable{}
	table.TypeName = info.TypeName
	table.Headers = make([]*DataTableHeader, 0)
	table.DefinedTable = info.DefinedTable

	for _, item := range info.Items {
		header := DataTableHeader{}
		header.Desc = item.Desc
		header.IsArray = item.IsArray
		header.RawValueType = item.RawValueType
		header.ValueType = item.ValueType
		header.FieldName = item.FieldName
		header.Index = item.Index
		table.Headers = append(table.Headers, &header)
	}

	return &table
}
