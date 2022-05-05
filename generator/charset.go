package generator

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/utils"
)

var charsetTemplate = ""

type charsetFileDesc struct {
	commonFileDesc
	Words []string
}

type charsetGenerator struct {
}

func (g *charsetGenerator) Generate(info *BuildInfo) (save bool, data *bytes.Buffer) {
	utils.CheckPath(info.Output)

	if charsetTemplate == "" {
		temp := getTemplate(info, "./template/charset.gtpl")
		log.Printf("[提示] 加载模板: %s \n", temp)

		data, err := ioutil.ReadFile(temp)
		if err != nil {
			log.Println(err)
			return false, nil
		}
		charsetTemplate = string(data)
	}

	tpl, err := template.New("charset").Funcs(funcs).Parse(charsetTemplate)
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	var fd = charsetFileDesc{}
	fd.Version = settings.TOOL_VERSION
	fd.GoProtoVersion = settings.GO_PROTO_VERTION

	lanTables := settings.LANG_TABLES
	if lanTables != nil && len(lanTables) > 0 {
		// 所有语言表数据
		datas := make([][]string, 0)
		for _, table := range lanTables {
			datas = append(datas, table.Data...)
		}

		var tempFileName = info.Output + ".tpl"
		var content = ""
		if ok, _ := utils.PathExists(tempFileName); ok {
			bts, err := os.ReadFile(tempFileName)
			if err == nil {
				content = content + string(bts)
			}
		}

		var tplfile = "./template/charset.tpl"
		if ok, _ := utils.PathExists(tplfile); ok {
			bts, err := os.ReadFile(tplfile)
			if err == nil {
				content = content + string(bts)
			}
		}

		var allChars = make(map[string]bool, 0)
		var cntStrs = strings.Split(content, "")
		for i := 0; i < len(cntStrs); i++ {
			allChars[cntStrs[i]] = true
		}

		for i := 0; i < len(datas); i++ {
			row := datas[i]
			for j := 0; j < len(row); j++ {
				var cell = row[j]
				var strs = strings.Split(cell, "")
				for k := 0; k < len(strs); k++ {
					allChars[strs[k]] = true
				}
			}
		}

		strs := make([]string, 0)
		for k, _ := range allChars {
			strs = append(strs, k)
		}
		sort.Strings(strs)
		fd.Words = strs
	}

	var buf = bytes.NewBufferString("")

	err = tpl.Execute(buf, &fd)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, buf
}

func init() {
	Regist("charset", &charsetGenerator{})
}
