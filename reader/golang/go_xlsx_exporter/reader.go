package go_xlsx_exporter

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/wingcd/go-xlsx-exporter/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type DataTableMap map[string]IData
type DataItemMap map[string]IData

type DataArray []interface{}
type DataMap map[string]interface{}

var (
	DefaultIndexKey = "ID"
	BytesFileExt    = ".bytes"

	dataDir      = ""
	dataItemMap  = make(DataItemMap, 0)
	dataTableMap = make(DataTableMap, 0)
)

type IFilenameGenerator interface {
	GetFilename(typeName string) string
}

type IData interface {
	DataType() reflect.Type
}

func DataDir() string {
	return dataDir
}

func Initial(dir, defaultKeyName string) {
	dataDir = strings.ReplaceAll(dir, "\\", "/")
	if strings.LastIndex(dataDir, "/") != len(dataDir)-1 {
		dataDir = dataDir + "/"
	}

	if defaultKeyName != "" {
		DefaultIndexKey = defaultKeyName
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
		dt = t
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
		dt = t
	} else {
		dt = item.(*DataTable)
	}
	return dt
}

func getFullname(dataType reflect.Type) string {
	var pkgs = strings.Split(dataType.PkgPath(), "/")
	var pkg = pkgs[len(pkgs)-1]
	var name = dataType.Name()
	return fmt.Sprintf("%s.%s", pkg, name)
}

type DataItem struct {
	dataType reflect.Type
	data     interface{}
	fileGen  IFilenameGenerator
}

func NewDataItem(dataType reflect.Type) *DataItem {
	t := new(DataItem)
	t.dataType = dataType
	t.fileGen = t
	return t
}

func (t *DataItem) DataType() reflect.Type {
	return t.dataType
}

func (t *DataItem) Item() interface{} {
	t.load()
	return t.data
}

func (t *DataItem) Clear() {
	t.data = nil
}

func (t *DataItem) GetFilename(typeName string) string {
	return fmt.Sprintf("%s%s%s", dataDir, strings.ToLower(typeName), BytesFileExt)
}

func (t *DataItem) load() {
	if t.data != nil {
		return
	}

	var filename = t.fileGen.GetFilename(t.dataType.Name())
	if ok, _ := utils.PathExists(filename); !ok {
		fmt.Printf("can not find data file %s", filename)
		return
	}

	fullName := getFullname(t.dataType)
	msgName := protoreflect.FullName(fullName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		fmt.Printf("can not find proto message named %v, %v", fullName, err)
		return
	}
	message := msgType.New().Interface()

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("read bytes file failed, %v", err)
		return
	}

	err = proto.Unmarshal(bytes, message)
	if err != nil {
		fmt.Printf("proto unmarshal failed, %v", err)
		return
	}

	t.data = message
}

type DataTable struct {
	indexKey string
	dataType reflect.Type
	data     DataArray
	dataMap  DataMap
	fileGen  IFilenameGenerator
}

func NewDataTable(indexKey string, dataType reflect.Type) *DataTable {
	if indexKey == "" {
		indexKey = DefaultIndexKey
	}

	t := new(DataTable)
	t.indexKey = indexKey
	t.dataType = dataType
	t.fileGen = t
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

func (t *DataTable) GetItem(id string) interface{} {
	var dataMap = t.ItemsMap()
	if dt, ok := dataMap[id]; ok {
		return dt
	}
	return nil
}

func (t *DataTable) Clear() {
	t.data = nil
	t.dataMap = nil
}

func (t *DataTable) GetFilename(typeName string) string {
	return fmt.Sprintf("%s%s%s", dataDir, strings.ToLower(typeName), BytesFileExt)
}

func (t *DataTable) load() {
	if t.data != nil {
		return
	}

	var filename = t.fileGen.GetFilename(t.dataType.Name())
	if ok, _ := utils.PathExists(filename); !ok {
		fmt.Printf("can not find data file %s", filename)
		return
	}

	fullName := getFullname(t.dataType) + "_ARRAY"
	msgName := protoreflect.FullName(fullName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		fmt.Printf("can not find proto message named %v, %v", fullName, err)
		return
	}
	message := msgType.New().Interface()

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("read bytes file failed, %v", err)
		return
	}

	err = proto.Unmarshal(bytes, message)
	if err != nil {
		fmt.Printf("proto unmarshal failed, %v", err)
		return
	}

	var value = reflect.ValueOf(message)
	var elm = value.Elem()
	if elm.Kind() != reflect.Struct {
		fmt.Printf("type %s kind error", fullName)
		return
	}

	t.data = make(DataArray, 0)
	var items = elm.FieldByName("Items").Interface()
	if reflect.TypeOf(items).Kind() == reflect.Slice {
		s := reflect.ValueOf(items)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			t.data = append(t.data, ele.Interface())
		}
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
		var idx = elm.FieldByName(t.indexKey).Interface()
		t.dataMap[fmt.Sprintf("%v", idx)] = item
	}
}
