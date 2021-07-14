package serialize

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wingcd/go-xlsx-protobuf/model"
	"github.com/wingcd/go-xlsx-protobuf/settings"
	"github.com/wingcd/go-xlsx-protobuf/utils"
	gproto "google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func getErrStr(typeName string, ridx, cidx int, val, colName, rowId string) string {
	return fmt.Sprintf("[错误] 值解析失败 类型:%s 列号：%v(Name:%s) 行号:%v(ID:%s) 值：%v \n", typeName, cidx, colName, ridx, rowId, val)
}

func GenDataTables(pbFilename string, fd pref.FileDescriptor, dir string, tables []*model.DataTable) bool {
	if fd == nil {
		f, err := utils.BuildFileDesc(pbFilename)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
		fd = f
	}

	for _, table := range tables {
		if ok, _ := GenDataTable(fd, dir, table); !ok {
			return false
		}
	}
	return true
}

func GenDataTable(fd pref.FileDescriptor, dir string, table *model.DataTable) (bool, string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	utils.CheckPath(dir)

	dtfileName := dir + strings.ToLower(table.TypeName) + settings.PbBytesFileExt
	fd, err := utils.BuildFileDesc("DataModel") // utils.BuildDynamicType([]*model.DataTable{table})
	if err != nil {
		log.Printf("类型构建失败 类型:%s 详情:%s \n", table.TypeName, err.Error())
		return false, dtfileName
	}

	// 表数据类型
	typeMD := fd.Messages().ByName(pref.Name(table.TypeName))
	typeListMD := fd.Messages().ByName(pref.Name(table.TypeName + "_ARRAY"))
	// 列表数据对象
	listItem := dynamicpb.NewMessage(typeListMD)
	lf := typeListMD.Fields().ByName("Items")
	//  创建类型为msg的数组变量Items
	list := listItem.NewField(lf).List()

	for ridx, row := range table.Data {
		// 单行数据实例
		item := dynamicpb.NewMessage(typeMD)
		for cidx, header := range table.Headers {
			cellValue := row[cidx]
			ok, value, _ := utils.ParseValue(header.RawValueType, cellValue)
			if !ok {
				panic(getErrStr(table.TypeName, ridx, cidx, cellValue, header.Desc, row[0]))
			}

			if header.IsArray {
				if settings.IsStruct(header.ValueType) {
					// to do...
				} else {
					// 创建此变量的list
					subField := typeMD.Fields().ByName(pref.Name(header.FieldName))
					// 列表属性
					subList := item.NewField(subField).List()
					for _, valItem := range value.([]interface{}) {
						val, err := utils.Convert2PBValue(header.ValueType, valItem)
						if err != nil {
							panic(getErrStr(table.TypeName, ridx, cidx, cellValue, header.Desc, row[0]))
						}
						subList.Append(val)
					}
					item.Set(subField, pref.ValueOf(subList))
				}
			} else {
				if settings.IsStruct(header.ValueType) {
					// to do...
				} else {
					val, err := utils.Convert2PBValue(header.ValueType, value)
					if err != nil {
						panic(getErrStr(table.TypeName, ridx, cidx, cellValue, header.Desc, row[0]))
					}
					item.Set(typeMD.Fields().ByName(pref.Name(header.FieldName)), val)
				}
			}
		}
		list.Append(pref.ValueOf(item))
	}
	listItem.Set(lf, pref.ValueOf(list))

	bytes, err := gproto.Marshal(listItem)
	if err != nil {
		panic(fmt.Sprintf("[错误] 序列化失败：%s \n", err.Error()))
	}
	f, err := os.Create(dtfileName)
	defer f.Close()
	if err != nil {
		panic(fmt.Sprintf("[错误] 文件创建失败：%s \n", err.Error()))
	}
	_, err = f.Write(bytes)
	if err != nil {
		panic(fmt.Sprintf("[错误] 文件写入失败：%s \n", err.Error()))
	}

	return true, dtfileName
}
