#### 自定义导出
导出类型为`type:custom`,且导出脚本配置为`template:对应的lua脚本`时，程序将执行对应的lua脚本生成数据：
* 项目示例导出csv脚本为:[csv_export](../template/data-gen.lua)
* 项目示例导出csv对应dts脚本为:[dts_export](../template/dts-gen.lua)

脚本中可以通过打印json查看所有数据结构：
```lua
print(GXE.json_encode(GXE))
```

lua中可以使用的入口为GXE，一下为数据表简单描述：
```golang
type BuildInfo struct {
	Imports  []string
	Output   string
	Template string
}

type commonFileDesc struct {
	Version        string
	GoProtoVersion string
	HasMessage     bool
}

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

type customFileDesc struct {
	commonFileDesc

	Namespace string
	Info      *BuildInfo
    // 单个文件导出时使用
	Table     *model.DataTable

	Enums  []*model.DefineTableInfo
	Consts []*model.DefineTableInfo
	Tables []*model.DataTable
}

// 导出到lua的数据结构GXE
type LuaExportInfo struct {
    version     string
    info        *BuildInfo
    fileDesc    *customFileDesc
}
```