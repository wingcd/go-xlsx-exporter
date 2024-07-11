package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
	lua "github.com/yuin/gopher-lua"
)

var customTemplate = ""

func customFormatValue(value interface{}, valueType string, isEnum bool, isArray bool) string {
	var ret = ""
	if isArray {
		var arr = value.([]interface{})
		var lst []string
		for _, it := range arr {
			lst = append(lst, customFormatValue(it, valueType, isEnum, false))
		}
		ret = fmt.Sprintf("[%s]", strings.Join(lst, ","))
	} else if isEnum {
		var enumStr = utils.ToEnumString(valueType, value.(int32))
		if enumStr != "" {
			ret = fmt.Sprintf("%v", value)
		} else {
			fmt.Printf("[错误] 值解析失败 类型:%s 值：%v \n", valueType, value)
		}
	} else if valueType == "float" {
		ret = fmt.Sprintf("%v", value)
	} else if valueType == "string" {
		ret = fmt.Sprintf("\"%v\"", value)
	} else {
		ret = fmt.Sprintf("%v", value)
	}
	return ret
}

var customGenetatorInited = false

func customValueDefault(item interface{}) string {
	var nilType = "null"
	switch inst := item.(type) {
	case *model.DataTableHeader:
		if inst.IsArray {
			return nilType
		} else if inst.IsEnum {
			var enumInfo = settings.GetEnum(inst.ValueType)
			if enumInfo != nil {
				return fmt.Sprintf("%s.%s", enumInfo.TypeName, enumInfo.Items[0].FieldName)
			}
		} else if inst.IsStruct {
			return nilType
		} else if val, ok := defaultCustomValue[inst.StandardValueType]; ok {
			return val
		}
	case *model.DataTable:
		return nilType
	case *model.DefineTableInfo:
		return fmt.Sprintf("%s_%s", inst.TypeName, inst.Items[0].FieldName)
	case string:
		if val, ok := defaultCustomValue[inst]; ok {
			return val
		} else if utils.IsEnum(inst) {
			var enumInfo = settings.GetEnum(inst)
			if enumInfo != nil {
				return fmt.Sprintf("%s_%s", enumInfo.TypeName, enumInfo.Items[0].FieldName)
			}
		} else if utils.IsTable(inst) || utils.IsStruct(inst) {
			return nilType
		}
	}
	return ""
}

func customValueFormat(value string, item interface{}) string {
	var isEnum = false
	var valueType = ""
	var rawValueType = ""
	var fieldName = ""
	switch inst := item.(type) {
	case *model.DefineTableItem:
		fieldName = inst.FieldName
		isEnum = inst.IsEnum
		valueType = inst.ValueType
		rawValueType = inst.RawValueType
	case *model.DataTableHeader:
		fieldName = inst.FieldName
		isEnum = inst.IsEnum
		valueType = inst.ValueType
		rawValueType = inst.RawValueType
	}

	var ok, val, isArray = utils.ParseValue(rawValueType, value)
	if !ok {
		fmt.Printf("[错误] 值解析失败 字段：%s 类型:%s 值：%v \n", fieldName, valueType, value)
		return value
	}
	return customFormatValue(val, valueType, isEnum, isArray)
}

func registCustomFuncs() {
	if customGenetatorInited {
		return
	}
	customGenetatorInited = true

	funcs["value_format"] = customValueFormat

	funcs["default"] = customValueDefault
}

var rsupportCustomTypes = map[string]string{
	"bool":   "bool",
	"int":    "int",
	"uint":   "uint",
	"int64":  "int64",
	"uint64": "uint64",
	"float":  "float",
	"double": "double",
	"string": "string",
	"bytes":  "string",
	"void":   "",
}

var defaultCustomValue = map[string]string{
	"bool":   "false",
	"int":    "0",
	"uint":   "0",
	"int64":  "0",
	"uint64": "0",
	"float":  "0",
	"double": "0",
	"string": "\"\"",
	"bytes":  "\"\"",
	"void":   "null",
}

type customFileDesc struct {
	commonFileDesc

	Namespace string
	Info      *BuildInfo
	Enum      *model.DefineTableInfo
	Const     *model.DefineTableInfo
	// 单个文件导出时使用
	Table *model.DataTable

	Enums  []*model.DefineTableInfo
	Consts []*model.DefineTableInfo
	Tables []*model.DataTable
}

type customGenerator struct {
}

func convertMapToLuaTable(L *lua.LState, m interface{}) *lua.LTable {
	tb := L.NewTable()

	reflectValue := reflect.ValueOf(m)
	for _, key := range reflectValue.MapKeys() {
		setValueToTable(L, reflectValue.MapIndex(key).Interface(), tb, key.String())
	}
	return tb
}

func convertArrayToLuaTable(L *lua.LState, arr interface{}) *lua.LTable {
	tb := L.NewTable()

	reflectValue := reflect.ValueOf(arr)
	for i := 0; i < reflectValue.Len(); i++ {
		setValueToTable(L, reflectValue.Index(i).Interface(), tb, strconv.Itoa(i+1))
	}
	return tb
}

func structToLuaTable(L *lua.LState, obj interface{}) *lua.LTable {
	if obj == nil {
		return nil
	}

	reflectValue := reflect.ValueOf(obj)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	if reflectValue.IsZero() {
		return nil
	}

	var tb = L.NewTable()

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Type().Field(i)
		if !field.IsExported() {
			continue
		}
		value := reflectValue.Field(i).Interface()
		setValueToTable(L, value, tb, field.Name)
	}

	return tb
}

func convertToLuaTable(L *lua.LState, obj interface{}) *lua.LTable {
	if obj == nil {
		return nil
	}

	reflectValue := reflect.ValueOf(obj)
	if reflectValue.IsZero() {
		return nil
	}

	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.IsZero() {
		return nil
	}

	if reflectValue.Kind() == reflect.Struct {
		return structToLuaTable(L, obj)
	} else if reflectValue.Kind() == reflect.Map {
		return convertMapToLuaTable(L, obj)
	} else if reflectValue.Kind() == reflect.Array || reflectValue.Kind() == reflect.Slice {
		return convertArrayToLuaTable(L, obj)
	} else {
		return nil
	}
}

func convertToLuaValue(L *lua.LState, obj interface{}) lua.LValue {
	if obj == nil {
		return lua.LNil
	}

	reflectValue := reflect.ValueOf(obj)

	switch reflectValue.Kind() {
	case reflect.String:
		return lua.LString(reflectValue.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(reflectValue.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(reflectValue.Uint())
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(reflectValue.Float())
	case reflect.Bool:
		return lua.LBool(reflectValue.Bool())
	}

	if reflectValue.IsZero() {
		return lua.LNil
	}

	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	if reflectValue.IsZero() {
		return lua.LNil
	}

	switch reflectValue.Kind() {
	case reflect.String:
		return lua.LString(reflectValue.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(reflectValue.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(reflectValue.Uint())
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(reflectValue.Float())
	case reflect.Bool:
		return lua.LBool(reflectValue.Bool())
	case reflect.Map:
		return convertMapToLuaTable(L, reflectValue.Interface())
	case reflect.Array:
		return convertArrayToLuaTable(L, reflectValue.Interface())
	case reflect.Slice:
		return convertArrayToLuaTable(L, reflectValue.Interface())
	case reflect.Struct:
		return structToLuaTable(L, obj)
	default:
		return lua.LNil
	}
}

func setValueToTable(L *lua.LState, obj interface{}, table *lua.LTable, fieldName string) *lua.LTable {
	lValue := convertToLuaValue(L, obj)
	table.RawSetString(fieldName, lValue)
	return table
}

func luaValueFormat(L *lua.LState) int {
	var value = L.CheckString(1)
	var item = L.CheckUserData(2)
	var ret = customValueFormat(value, item)
	L.Push(lua.LString(ret))
	return 1
}

func luaDefaultValue(L *lua.LState) int {
	var item = L.CheckUserData(1)
	var ret = customValueDefault(item)
	L.Push(lua.LString(ret))
	return 1
}

func luaIsInterger(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = isInterger(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsLong(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = isLong(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsFloat(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = isFloat(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsNumber(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = isNumber(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsBool(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = isBool(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsString(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = isString(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsBytes(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = isBytes(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsStruct(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = utils.IsStruct(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsValueType(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = isValueType(value)
	L.Push(lua.LBool(ret))
	return 1
}

func luaGetWireType(L *lua.LState) int {
	var value = L.CheckString(1)
	var ret = utils.GetWireType(value)
	L.Push(lua.LString(ret))
	return 1
}

func luaGetWireOffset(L *lua.LState) int {
	var item = L.CheckUserData(1)
	var ret = getWireOffset(item)
	L.Push(lua.LNumber(ret))
	return 1
}

func luaGetEnum(L *lua.LState) int {
	var pbType = L.CheckString(1)
	var ret = utils.GetEnum(pbType)
	if ret == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(convertToLuaValue(L, ret))
	}

	return 1
}

func luaGetEnumDefault(L *lua.LState) int {
	var item = L.CheckString(1)
	var ret = utils.GetEnumDefault(item)
	if ret == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(convertToLuaValue(L, ret))
	}
	return 1
}

func luaGetEnumValues(L *lua.LState) int {
	var item = L.CheckString(1)
	var ret = utils.GetEnumValues(item)
	if ret == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(convertToLuaValue(L, ret))
	}
	return 1
}

func luaGetEnumNames(L *lua.LState) int {
	var item = L.CheckString(1)
	var ret = utils.GetEnumNames(item)
	if ret == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(convertToLuaValue(L, ret))
	}
	return 1
}

func luaGetEnumString(L *lua.LState) int {
	var pbType = L.CheckString(1)
	var value = L.CheckNumber(2)
	var ret = utils.ToEnumString(pbType, int32(value))
	L.Push(lua.LString(ret))
	return 1
}

func luaIsDefineTable(L *lua.LState) int {
	var tableType = L.CheckInt(1)
	var ret = tableType == int(model.ETableType_Define)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsDataTable(L *lua.LState) int {
	var tableType = L.CheckInt(1)
	var ret = tableType == int(model.ETableType_Data)
	L.Push(lua.LBool(ret))
	return 1
}

func luaIsMessageTable(L *lua.LState) int {
	var tableType = L.CheckInt(1)
	var ret = tableType == int(model.ETableType_Message)
	L.Push(lua.LBool(ret))
	return 1
}

var customFunctions = map[string]lua.LGFunction{
	"json_encode":      apiEncode,
	"json_decode":      apiDecode,
	"value_format":     luaValueFormat,
	"default_value":    luaDefaultValue,
	"is_interger":      luaIsInterger,
	"is_long":          luaIsLong,
	"is_float":         luaIsFloat,
	"is_number":        luaIsNumber,
	"is_bool":          luaIsBool,
	"is_string":        luaIsString,
	"is_bytes":         luaIsBytes,
	"is_struct":        luaIsStruct,
	"is_value_type":    luaIsValueType,
	"get_wire_type":    luaGetWireType,
	"calc_wire_offset": luaGetWireOffset,
	"get_enum":         luaGetEnum,
	"get_enum_default": luaGetEnumDefault,
	"get_enum_values":  luaGetEnumValues,
	"get_enum_names":   luaGetEnumNames,
	"get_enum_string":  luaGetEnumString,
	"is_define_table":  luaIsDefineTable,
	"is_data_table":    luaIsDataTable,
	"is_message_table": luaIsMessageTable,
}

func callLuaFunc(L *lua.LState, funcName string, args ...interface{}) string {
	fn := L.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		log.Fatalf("[错误] 函数%s 未定义 \n", funcName)
		return ""
	}
	L.Push(fn)

	for _, arg := range args {
		reflectValue := reflect.ValueOf(arg)
		switch reflectValue.Kind() {
		case reflect.String:
			L.Push(lua.LString(reflectValue.String()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			L.Push(lua.LNumber(reflectValue.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			L.Push(lua.LNumber(reflectValue.Uint()))
		case reflect.Float32, reflect.Float64:
			L.Push(lua.LNumber(reflectValue.Float()))
		case reflect.Bool:
			L.Push(lua.LBool(reflectValue.Bool()))
		case reflect.Map:
			L.Push(convertMapToLuaTable(L, arg))
		case reflect.Array:
			L.Push(convertArrayToLuaTable(L, arg))
		case reflect.Slice:
			L.Push(convertArrayToLuaTable(L, arg))
		case reflect.Struct:
			L.Push(structToLuaTable(L, arg))
		case reflect.Ptr:
			L.Push(structToLuaTable(L, reflectValue.Elem().Interface()))
		default:
			L.Push(lua.LNil)
		}
	}

	if err := L.PCall(len(args), 1, nil); err != nil {
		log.Fatalf("[错误] %v \n", err)
		return ""
	}

	ret := L.Get(-1)
	L.Pop(1)

	if ret.Type() == lua.LTString {
		return ret.String()
	}

	return ""
}

func genByLua(fd *customFileDesc, info *BuildInfo, luaFile string) {
	L := lua.NewState()
	defer L.Close()

	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exePath) + "/template/?.lua"
	exeDir = strings.Replace(exeDir, "\\", "/", -1)
	L.DoString("package.path = package.path .. ';" + exeDir + "' \n")

	var envInfo = L.NewTable()
	L.SetFuncs(envInfo, customFunctions)
	L.SetField(envInfo, "version", lua.LString(settings.TOOL_VERSION))
	L.SetField(envInfo, "info", convertToLuaTable(L, info))
	L.SetField(envInfo, "fileDesc", convertToLuaTable(L, fd))
	L.SetGlobal("GXE", envInfo)

	if err := L.DoFile(luaFile); err != nil {
		log.Fatalf("[错误] %v \n", err)
		return
	}

	if L.GetGlobal("generate") == lua.LNil {
		return
	}

	var retStr = callLuaFunc(L, "generate")
	if retStr == "" {
		log.Fatalf("[错误] 生成失败 \n")
		return
	}

	var t = fd.Table
	if t == nil {
		var fileName = info.Output
		if err := ioutil.WriteFile(fileName, []byte(retStr), 0644); err != nil {
			log.Println(err)
			return
		}
	} else {
		var fileName = fmt.Sprintf(info.Output, t.TypeName)
		if err := ioutil.WriteFile(fileName, []byte(retStr), 0644); err != nil {
			log.Println(err)
			return
		}
	}
}

func genCustomFile(fd *customFileDesc, info *BuildInfo, tpl *template.Template) bool {
	var t = fd.Table

	if t == nil {
		var buf = bytes.NewBufferString("")

		err := tpl.Execute(buf, &fd)
		if err != nil {
			log.Println(err)
			return false
		}

		var fileName = info.Output
		if err := ioutil.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
			log.Println(err)
			return false
		}
	} else {
		if t.TableType == model.ETableType_Message {
			fd.HasMessage = true
		}

		// 排除语言类型
		if t.TableType == model.ETableType_Language && !settings.GenLanguageType {
			return false
		}

		// 排除配置
		if t.TableType == model.ETableType_Define {
			return false
		}

		// 处理类型
		for _, h := range t.Headers {
			if !h.IsEnum && !h.IsStruct {
				if _, ok := rsupportCustomTypes[h.StandardValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.RawValueType, t.DefinedTable, h.FieldName)
					return false
				}
			}
		}

		utils.PreProcessTable(fd.Table)

		var buf = bytes.NewBufferString("")

		err := tpl.Execute(buf, &fd)
		if err != nil {
			log.Println(err)
			return false
		}

		var fileName = fmt.Sprintf(info.Output, t.TypeName)
		if err := ioutil.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
			log.Println(err)
			return false
		}
	}

	return true
}

func (g *customGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	registCustomFuncs()

	if customTemplate == "" {
		temp := getTemplate(info, "./template/custom.gtpl")
		log.Printf("[提示] 加载模板: %s \n", temp)

		data, err := ioutil.ReadFile(temp)
		if err != nil {
			log.Println(err)
			return false, nil
		}
		customTemplate = string(data)
	}

	tpl, err := template.New("custom").Funcs(funcs).Parse(customTemplate)
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	var fd = customFileDesc{
		Namespace: settings.PackageName,
		Info:      info,
		Table:     nil,

		Enums:  settings.ENUMS[:],
		Consts: settings.CONSTS[:],
		Tables: make([]*model.DataTable, 0),
	}
	fd.Version = settings.TOOL_VERSION
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	utils.PreProcessTables(fd.Tables)
	for _, t := range settings.TABLES {
		if t.TableType == model.ETableType_Message {
			fd.HasMessage = true
		}

		// 排除语言类型
		if t.TableType == model.ETableType_Language && !settings.GenLanguageType {
			continue
		}

		// 排除配置
		if t.TableType == model.ETableType_Define {
			continue
		}

		fd.Tables = append(fd.Tables, t)

		// 处理类型
		for _, h := range t.Headers {
			if !h.IsEnum && !h.IsStruct {
				if _, ok := supportTSTypes[h.StandardValueType]; !ok {
					log.Printf("[错误] 不支持类型%s 表：%s 列：%s \n", h.RawValueType, t.DefinedTable, h.FieldName)
					return false, nil
				}
				h.ValueType = supportTSTypes[h.StandardValueType]
			}
		}

		if t.NeedAddItems {
			// 添加数组类型
			table := model.DataTable{}
			table.DefinedTable = t.DefinedTable
			table.TypeName = t.TypeName + "_ARRAY"
			table.IsArray = true
			header := model.DataTableHeader{}
			header.Index = 1
			header.FieldName = "Items"
			header.TitleFieldName = header.FieldName
			header.IsArray = true
			header.ValueType = t.TypeName
			header.RawValueType = t.TypeName + "[]"
			header.IsMessage = true
			table.Headers = []*model.DataTableHeader{&header}

			fd.Tables = append(fd.Tables, &table)
		}
	}
	utils.PreProcessTables(fd.Tables)

	tables := fd.Tables
	allInOne := !strings.Contains(info.Output, "%s")
	var isLua = strings.HasSuffix(info.Template, ".lua")
	if isLua {
		if allInOne {
			genByLua(&fd, info, info.Template)
		} else {
			for _, t := range tables {
				if t.IsArray {
					continue
				}

				fd.Table = t
				genByLua(&fd, info, info.Template)
			}
		}
	} else {
		if allInOne {
			genCustomFile(&fd, info, tpl)
		} else {
			for _, t := range tables {
				if t.IsArray {
					continue
				}

				fd.Table = t
				genCustomFile(&fd, info, tpl)
			}
		}
	}

	return false, nil
}

func init() {
	Regist("custom", &customGenerator{})
}
