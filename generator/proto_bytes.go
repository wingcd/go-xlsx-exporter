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

	if !serialize.GenDataTables("", nil, output, settings.TABLES, settings.LANG_TABLES) {
		fmt.Printf("[错误] protobuf序列化失败，路径：%s \n", output)
	}
	return false, nil
}

func init() {
	Regist("proto_bytes", &protoBytesGenerator{})
}
