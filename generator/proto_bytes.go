package generator

import (
	"bytes"
	"fmt"

	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/xlsx"

	"github.com/wingcd/go-xlsx-exporter/serialize"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
)

type protoBytesGenerator struct {
}

func (g *protoBytesGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	utils.CheckPath(info.Output)

	fd, _ := utils.BuildFileDesc("", true)

	for _, setting := range settings.CONSTS {
		xlsx.CheckDefine(setting)
	}

	for _, table := range settings.TABLES {
		xlsx.CheckTable(table)
	}

	if info.Output[len(info.Output)-1] != '/' {
		info.Output += "/"
	}

	tables := make([]*model.DataTable, 0)
	for _, table := range settings.TABLES {
		if table.TableType != model.ETableType_Message {
			tables = append(tables, table)
		}
	}

	if !serialize.GenDataTables("", fd, info.Output, tables) {
		fmt.Printf("[错误] protobuf数据序列化失败，路径：%s \n", info.Output)
	}

	if !serialize.GenLanguageTables("", fd, info.Output, settings.TABLES, settings.LANG_TABLES) {
		fmt.Printf("[错误] protobuf多语言序列化失败，路径：%s \n", info.Output)
	}

	if !serialize.GenDefineTables("", fd, info.Output, settings.CONSTS) {
		fmt.Printf("[错误] protobuf配置序列化失败，路径：%s \n", info.Output)
	}
	return false, nil
}

func init() {
	Regist("proto_bytes", &protoBytesGenerator{})
}
