package generator

import (
	"bytes"
	"fmt"

	"github.com/wingcd/go-xlsx-exporter/serialize"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
)

type protoBytesGenerator struct {
}

func (g *protoBytesGenerator) Generate(output string) (save bool, data *bytes.Buffer) {
	utils.CheckPath(output)

	fd, _ := utils.BuildFileDesc("", true)

	if !serialize.GenDataTables("", fd, output, settings.TABLES) {
		fmt.Printf("[错误] protobuf数据序列化失败，路径：%s \n", output)
	}

	if !serialize.GenLanguageTables("", fd, output, settings.TABLES, settings.LANG_TABLES) {
		fmt.Printf("[错误] protobuf多语言序列化失败，路径：%s \n", output)
	}

	if !serialize.GenDefineTables("", fd, output, settings.CONSTS) {
		fmt.Printf("[错误] protobuf配置序列化失败，路径：%s \n", output)
	}
	return false, nil
}

func init() {
	Regist("proto_bytes", &protoBytesGenerator{})
}
