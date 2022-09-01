package serialize

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
	gproto "google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func getErrStr(typeName string, ridx, cidx int, val, colName, rowId string) string {
	return fmt.Sprintf("[错误] 值解析失败 类型:%s 列号：%v(Name:%s) 行号:%v(ID:%s) 值：%v \n", typeName, cidx, colName, ridx, rowId, val)
}

func GenDataTables(pbFilename string, fd pref.FileDescriptor, dir string, tables []*model.DataTable) bool {
	if fd == nil {
		f, err := utils.BuildFileDesc(pbFilename, settings.GenLanguageType)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
		fd = f
	}

	for _, table := range tables {
		if table.TableType == model.ETableType_Data {
			if ok, _ := GenDataTable(fd, dir, table, ""); !ok {
				return false
			}
		}
	}

	return true
}

func GenLanguageTables(pbFilename string, fd pref.FileDescriptor, dir string, tables []*model.DataTable, lanTables []*model.DataTable) bool {
	if fd == nil {
		f, err := utils.BuildFileDesc(pbFilename, settings.GenLanguageType)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
		fd = f
	}

	var langTable *model.DataTable = nil
	for _, table := range tables {
		if table.TableType == model.ETableType_Language {
			langTable = table
			break
		}
	}

	// 生成多语言文件
	if lanTables != nil && langTable != nil {
		// 所有语言表组合
		datas := make([][]string, 0)
		// 所有语言类型
		langs := make([]string, 0)
		for i, table := range lanTables {
			datas = append(datas, table.Data...)

			if i == 0 {
				for hi := 1; hi < len(table.Headers); hi++ {
					var header = table.Headers[hi]
					if header.IsVoid {
						continue
					}

					langs = append(langs, strings.ToLower(header.FieldName))
				}
			}
		}

		// 构建数据
		var langTableName = strings.ToLower(langTable.TypeName)
		for i, lan := range langs {
			langTable.Data = make([][]string, 0)
			for _, row := range datas {
				langTable.Data = append(langTable.Data, []string{row[0], row[i+1]})
			}
			filename := fmt.Sprintf("%s.%s", langTableName, lan)
			if ok, _ := GenDataTable(fd, dir, langTable, filename); !ok {
				return false
			}
		}
	}

	return true
}

func GenDefineTables(pbFilename string, fd pref.FileDescriptor, dir string, tables []*model.DefineTableInfo) bool {
	if fd == nil {
		f, err := utils.BuildFileDesc(pbFilename, settings.GenLanguageType)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
		fd = f
	}

	for _, table := range tables {
		if ok, _ := GenDefineTable(fd, dir, table); !ok {
			return false
		}
	}
	return true
}

func GenDefineTable(fd pref.FileDescriptor, dir string, table *model.DefineTableInfo) (bool, string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	utils.CheckPath(dir)

	dtfileName := dir + strings.ToLower(table.TypeName) + settings.PbBytesFileExt

	// 表数据类型
	typeMD := fd.Messages().ByName(pref.Name(table.TypeName))
	fmt.Printf("开始序列化表：%v\n", table.DefinedTable)

	// 单行数据实例
	item := dynamicpb.NewMessage(typeMD)
	for cidx, ditem := range table.Items {
		cellValue := ditem.Value
		ok, value, _ := utils.ParseValue(ditem.RawValueType, cellValue)
		if !ok {
			panic(getErrStr(table.TypeName, 0, cidx, cellValue, ditem.Desc, ditem.FieldName))
		}

		if ditem.IsArray {
			if utils.IsStruct(ditem.ValueType) {
				// to do...
			} else {
				// 创建此变量的list
				subField := typeMD.Fields().ByName(pref.Name(ditem.FieldName))
				// 列表属性
				subList := item.NewField(subField).List()
				for _, valItem := range value.([]interface{}) {
					val, err := utils.Convert2PBValue(ditem.ValueType, valItem)
					if err != nil {
						panic(getErrStr(table.TypeName, 0, cidx, cellValue, ditem.Desc, ditem.FieldName))
					}
					subList.Append(val)
				}
				item.Set(subField, pref.ValueOf(subList))
			}
		} else {
			if utils.IsStruct(ditem.ValueType) {
				// to do...
			} else {
				val, err := utils.Convert2PBValue(ditem.ValueType, value)
				if err != nil {
					panic(getErrStr(table.TypeName, 0, cidx, cellValue, ditem.Desc, ditem.FieldName))
				}
				item.Set(typeMD.Fields().ByName(pref.Name(ditem.FieldName)), val)
			}
		}
	}

	bytes, err := gproto.MarshalOptions{
		Deterministic: true,
	}.Marshal(item)

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

func GenDataTable(fd pref.FileDescriptor, dir string, table *model.DataTable, filename string) (bool, string) {
	var rowData []string
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			if rowData != nil {
				log.Println(rowData)
			}
		}
	}()

	utils.CheckPath(dir)

	dtfileName := ""
	if filename != "" {
		dtfileName = dir + filename + settings.PbBytesFileExt
	} else {
		dtfileName = dir + strings.ToLower(table.TypeName) + settings.PbBytesFileExt
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
		rowData = row

		// 单行数据实例
		item := dynamicpb.NewMessage(typeMD)
		var idx = 0
		for cidx, header := range table.Headers {
			if header.IsVoid {
				continue
			}

			var cellValue string
			if len(row) > idx {
				cellValue = row[idx]
			}
			idx++

			ok, value, _ := utils.ParseValue(header.RawValueType, cellValue)
			if !ok {
				panic(getErrStr(table.TypeName, ridx, cidx, cellValue, header.Desc, row[0]))
			}

			if header.IsArray {
				if utils.IsStruct(header.ValueType) {
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
				if utils.IsStruct(header.ValueType) {
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

	bytes, err := gproto.MarshalOptions{
		Deterministic: true,
	}.Marshal(listItem)

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
