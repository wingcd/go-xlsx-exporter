package go_xlsx_exporter

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/wingcd/go-xlsx-exporter/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type DataTableMap map[string]IData
type DataItemMap map[string]IData

type DataArray []interface{}
type DataMap map[string]interface{}

var (
	DefaultIndexKey = "ID"
	DataDir         = ""
	BytesFileExt    = ".bytes"

	dataItemMap  = make(DataItemMap, 0)
	dataTableMap = make(DataTableMap, 0)
	typeMap      = make(map[string]protoreflect.MessageType)
)

type IData interface {
	DataType() reflect.Type
}

func init() {
	DataDir = strings.ReplaceAll(DataDir, "\\", "/")
	if strings.LastIndex(DataDir, "/") != len(DataDir)-1 {
		DataDir = DataDir + "/"
	}
}

func Regist(fd protoreflect.FileDescriptor) {
	for i := 0; i < fd.Messages().Len(); i++ {
		var msg = fd.Messages().Get(i)
		var protoType = msg.Options().ProtoReflect().Type()
		typeMap[string(msg.FullName())] = protoType
	}
}

func RegistDataTable(indexKey string, dataType reflect.Type) {
	if _, ok := dataTableMap[dataType.Name()]; !ok {
		var t = NewDataTable(indexKey, dataType)
		dataTableMap[dataType.Name()] = t
	}
}

func RegistDataTableExt(table IData) {
	var dataType = table.DataType()
	if _, ok := dataTableMap[dataType.Name()]; !ok {
		dataTableMap[dataType.Name()] = table
	}
}

func GetDataItem(dataType reflect.Type) *DataItem {
	var dt *DataItem
	if item, ok := dataItemMap[dataType.Name()]; !ok {
		var t = NewDataItem(dataType)
		dataItemMap[dataType.Name()] = t
	} else {
		dt = item.(*DataItem)
	}
	return dt
}

func GetDataTable(dataType reflect.Type) *DataTable {
	var dt *DataTable
	if item, ok := dataTableMap[dataType.Name()]; !ok {
		var t = NewDataTable("", dataType)
		dataTableMap[dataType.Name()] = t
	} else {
		dt = item.(*DataTable)
	}
	return dt
}

type DataItem struct {
	dataType reflect.Type
	data     interface{}
}

func NewDataItem(dataType reflect.Type) *DataItem {
	t := new(DataItem)
	t.dataType = dataType
	return t
}

func (t *DataItem) DataType() reflect.Type {
	return t.dataType
}

func (t *DataItem) Item() interface{} {
	return t.data
}

func (t *DataItem) Clear() {
	t.data = nil
}

func (t *DataItem) load() {
	if t.data != nil {
		return
	}

	var filename = fmt.Sprintf("%s%s%s", DataDir, strings.ToLower(t.dataType.Name()), BytesFileExt)
	if ok, _ := utils.PathExists(filename); !ok {
		fmt.Printf("can not find data file %s", filename)
		return
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("read bytes file failed, %v", err)
		return
	}

	typeName := t.dataType.Name()
	if tp, ok := typeMap[typeName]; ok {
		var data = tp.New().Interface()
		proto.Unmarshal(bytes, data)
		t.data = data
	}
}

type DataTable struct {
	indexKey string
	dataType reflect.Type
	data     DataArray
	dataMap  DataMap
}

func NewDataTable(indexKey string, dataType reflect.Type) *DataTable {
	if indexKey == "" {
		indexKey = DefaultIndexKey
	}

	t := new(DataTable)
	t.indexKey = indexKey
	t.dataType = dataType
	return t
}

func (t *DataTable) DataType() reflect.Type {
	return t.dataType
}

func (t *DataTable) Items() DataArray {
	if t.data == nil {
		t.load()
	}
	return t.data
}

func (t *DataTable) ItemsMap() DataMap {
	if t.data == nil {
		t.load()
	}
	t.itemsAsMap()
	return t.dataMap
}

func (t *DataTable) Clear() {
	t.data = nil
	t.dataMap = nil
}

func (t *DataTable) GetFilename(typeName string) string {
	return fmt.Sprintf("%s%s%s", DataDir, strings.ToLower(typeName), BytesFileExt)
}

func (t *DataTable) load() {
	if t.data != nil {
		return
	}

	var filename = t.GetFilename(t.dataType.Name())
	if ok, _ := utils.PathExists(filename); !ok {
		fmt.Printf("can not find data file %s", filename)
		return
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("read bytes file failed, %v", err)
		return
	}

	arrTypeName := t.dataType.Name() + "_ARRAY"
	if tp, ok := typeMap[arrTypeName]; ok {
		var arr = tp.New().Interface()
		proto.Unmarshal(bytes, arr)

		var value = reflect.ValueOf(arr)
		var elm = value.Elem()
		if elm.Kind() != reflect.Struct {
			fmt.Printf("type %s kind error", arrTypeName)
			return
		}

		t.data = elm.FieldByName("Items").Interface().([]interface{})[:]
	} else {
		t.data = make(DataArray, 0)
	}
}

func (t *DataTable) itemsAsMap() {
	if t.dataMap != nil {
		return
	}
	t.dataMap = make(DataMap)

	for _, item := range t.data {
		var value = reflect.ValueOf(item)
		var elm = value.Elem()
		var idx = elm.FieldByName(t.indexKey).Interface().(string)
		t.dataMap[idx] = item
	}
}
