package model

type StructInfo struct {
	// 描述（注释）
	Desc string
	// 字段名
	FieldName string
	// 首字母大写字段名
	TitleFieldName string
	// 表中定义的原始类型
	RawValueType string
	// 转换后的标准类型
	StandardValueType string
	// protobuf类型
	PBValueType string
	// 转换后的值类型
	ValueType string
	// 二进制编码方式
	EncodeType string
	// 基础值是否枚举
	IsEnum bool
	// 基础值是否结构
	IsStruct bool
	// 基础值是否数组
	IsArray bool
	// 编号（1开始）
	Index int
	// 数组分隔符，默认为全局配置符号
	ArraySplitChar string
	// 此字段是否可转换对象
	Convertable bool
	// 换换后的类型是否需要缓存
	Cachable bool
	IsVoid   bool
	// 是否消息类型
	IsMessage bool
	// 别名（可转换对象显示类型）
	Alias string
	// 字段限制规则
	Rule int
}

// 定义表数据项（按行）
type DefineTableItem struct {
	StructInfo
	// 值
	Value string
}

type DefineTableItems []*DefineTableItem

func (a DefineTableItems) Len() int           { return len(a) }
func (a DefineTableItems) Less(i, j int) bool { return a[i].Index < a[j].Index }
func (a DefineTableItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 定义表类型（同类型分组）
type DefineTableInfo struct {
	StartID int64
	// 类型（enum/struct）
	Category string
	// 类型名
	TypeName string
	// 表名
	DefinedTable string
	// 类型子项
	Items []*DefineTableItem
}

type DefineTableInfos []*DefineTableInfo

func (a DefineTableInfos) Len() int           { return len(a) }
func (a DefineTableInfos) Less(i, j int) bool { return a[i].StartID < a[j].StartID }
func (a DefineTableInfos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 判断是否为当前定义类型
func (d *DefineTableInfo) IsValid(typeName string) bool {
	for _, item := range d.Items {
		if item.FieldName == typeName {
			return true
		}
	}
	return false
}

// 数据表表头
type DataTableHeader struct {
	// 是否支持客户端导出
	ExportClient bool
	// 是否支持服务器导出
	ExportServer bool

	StructInfo
}

type ETableType int

const (
	ETableType_Define   ETableType = iota + 1 // 数据表
	ETableType_Data                           // 数据表
	ETableType_Language                       // 语言
	ETableType_Message                        // 消息
)

// 数据表
type DataTable struct {
	Id int
	// 类型名
	TypeName string
	// 表头
	Headers []*DataTableHeader
	// 表文件名
	DefinedTable string
	// 数据
	Data [][]string
	// 表类型
	TableType ETableType
	// 是否数组
	IsArray bool
	// 是否需要增加子项
	NeedAddItems bool
}

type DataTables []*DataTable

func (a DataTables) Len() int           { return len(a) }
func (a DataTables) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a DataTables) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 定义的结构体转表类型
func Struct2Table(info *DefineTableInfo) *DataTable {
	if info.Category != DEFINE_TYPE_STRUCT && info.Category != DEFINE_TYPE_CONST {
		return nil
	}
	table := DataTable{}
	table.TypeName = info.TypeName
	table.Headers = make([]*DataTableHeader, 0)
	table.DefinedTable = info.DefinedTable
	table.TableType = ETableType_Define

	for _, item := range info.Items {
		header := DataTableHeader{}
		header.ExportClient = true
		header.ExportServer = true
		header.StructInfo = item.StructInfo
		table.Headers = append(table.Headers, &header)
	}

	return &table
}

// 定义的常量转表类型
func Const2Table(info *DefineTableInfo) *DataTable {
	if info.Category != DEFINE_TYPE_CONST {
		return nil
	}
	table := DataTable{}
	table.TypeName = info.TypeName
	table.Headers = make([]*DataTableHeader, 0)
	table.DefinedTable = info.DefinedTable
	table.TableType = ETableType_Define

	for _, item := range info.Items {
		header := DataTableHeader{}
		header.ExportClient = true
		header.ExportServer = true
		header.StructInfo = item.StructInfo
		table.Headers = append(table.Headers, &header)
	}

	return &table
}

func SetNotSupportExportType() {
	DATA_ROW_FIELD_INDEX = 2
	FixedRowCount = 3
}
