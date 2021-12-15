package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/wingcd/go-xlsx-exporter/generator"
	"github.com/wingcd/go-xlsx-exporter/model"
	"github.com/wingcd/go-xlsx-exporter/settings"
	"github.com/wingcd/go-xlsx-exporter/xlsx"
)

var (
	config *YamlConf

	p_help              bool
	p_version           bool
	p_gen_language_code bool
	p_exports           string

	p_package           string
	p_pb_bytes_file_ext string
	p_comment_symbol    string
	p_config            string
	p_silent            bool
)

func init() {
	flag.BoolVar(&p_help, "h", false, "获取帮助")
	flag.BoolVar(&p_version, "v", false, "获取工具当前版本号")

	flag.StringVar(&p_config, "cfg", "./conf.yaml", "设置配置文件")
	flag.StringVar(&p_package, "pkg", "", "设置导出包名")
	flag.StringVar(&p_pb_bytes_file_ext, "ext", "", "设置二进制数据文件后缀(unity必须为.bytes)")
	flag.StringVar(&p_comment_symbol, "cmt", "#", "设置表格注释符号")
	flag.StringVar(&p_exports, "exports", "", "设置需要导出的配置项，默认为空，全部导出, 参考：1,2,5-7")

	flag.BoolVar(&p_gen_language_code, "lang", false, "是否生成语言类型到代码（仅测试用，默认为false）")
	flag.BoolVar(&p_silent, "silent", false, "是否静默执行（默认为false）")
}

func main() {
	parseParams()
}

func parseParams() {
	flag.Parse()

	config = NewYamlConf(p_config)
	if p_pb_bytes_file_ext != "" {
		config.PBBytesFileExt = p_pb_bytes_file_ext
	}
	if config.PBBytesFileExt == "" {
		config.PBBytesFileExt = ".bytes"
	}
	config.CommentSymbol = p_comment_symbol

	if p_help {
		flag.Usage()
		return
	}

	if p_version {
		fmt.Printf("go-xlsx-exporter version is %v \n", settings.TOOL_VERSION)
		return
	}

	settings.PbBytesFileExt = config.PBBytesFileExt
	settings.CommentSymbol = config.CommentSymbol
	settings.GenLanguageType = p_gen_language_code
	settings.ArraySplitChar = config.ArraySplitChar

	if settings.ArraySplitChar == "" {
		settings.ArraySplitChar = "|"
	}
	settings.StrictMode = config.StrictMode
	if p_silent {
		config.PauseOnEnd = false
	}

	if settings.CommentSymbol == "" {
		log.Fatalln("注释符号不能为空")
	}

	process()

	if !p_silent && config.PauseOnEnd {
		pause()
	}
}

func pause() {
	fmt.Printf("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

// val sample: 1,2-4,6-7
func getIds(val string) map[int]bool {
	contains := make(map[int]bool)
	strs := strings.Split(val, ",")
	for _, v := range strs {
		if strings.Contains(v, "-") {
			innerStrs := strings.Split(v, "-")
			if len(innerStrs) != 2 {
				log.Fatalf("参数错误：%v \n", val)
			}

			start, err := strconv.Atoi(innerStrs[0])
			if err != nil {
				log.Fatalf("参数错误：%v, %v \n", val, err.Error())
			}

			end, err := strconv.Atoi(innerStrs[1])
			if err != nil {
				log.Fatalf("参数错误：%v, %v \n", val, err.Error())
			}

			if start > end {
				log.Fatalf("参数错误：%v, %v \n", val, v)
			}

			for i := start; i <= end; i++ {
				contains[i] = true
			}
		} else if v != "" {
			value, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("参数错误：%v, %v, %v \n", val, v, err.Error())
			}

			contains[value] = true
		}
	}

	return contains
}

func process() {
	exportIds := getIds(p_exports)
	var exports []ExportInfo
	if len(exportIds) == 0 {
		exports = config.Exports[:]
	} else {
		exports = make([]ExportInfo, 0)
		for _, info := range config.Exports {
			if _, ok := exportIds[info.ID]; ok {
				exports = append(exports, info)
			}
		}
	}

	for _, info := range exports {
		doExport(info)
	}
}

func doExport(exportInfo ExportInfo) {
	if exportInfo.Type == "" {
		log.Fatalln("导出类型不能为空")
	}

	if exportInfo.Path == "" {
		log.Fatalln("导出路径不能为空")
	}

	settings.ExportType = exportInfo.ExportType
	if exportInfo.ExportType != 0 {
		settings.ExportType = exportInfo.ExportType
	}
	if settings.ExportType <= 0 {
		settings.ExportType = 1
	}

	settings.PackageName = config.Package
	if p_package != "" {
		settings.PackageName = p_package
	} else if p_package == "" && exportInfo.Package != "" {
		// 当配置被覆盖时，统一使用外部参数，否则可以使用单项配置
		settings.PackageName = exportInfo.Package
	}
	if settings.PackageName == "" {
		settings.PackageName = "Deploy"
	}

	fmt.Printf("执行导出任务，id:%v, 类型：%v, 表：%v, 导出路径：%v, 导出类型：%v \n",
		exportInfo.ID, exportInfo.Type, exportInfo.Sheets, exportInfo.Path, []string{"忽略前后端配置", "仅客户端", "仅服务器"}[settings.ExportType-1])

	sheetsIds := getIds(exportInfo.Sheets)
	var sheets []SheetInfo

	if len(sheetsIds) == 0 {
		sheets = config.Sheets[:]
	} else {
		sheets = make([]SheetInfo, 0)
		for _, info := range config.Sheets {
			if _, ok := sheetsIds[info.ID]; ok {
				sheets = append(sheets, info)
			}
		}
	}

	defineSheets := make([]SheetInfo, 0)
	dataSheets := make([]SheetInfo, 0)
	langSheets := make([]SheetInfo, 0)
	for _, info := range sheets {
		tp := strings.ToLower(info.Type)
		if tp == "define" {
			defineSheets = append(defineSheets, info)
		} else if tp == "table" {
			if !info.IsLang {
				dataSheets = append(dataSheets, info)
			} else {
				langSheets = append(langSheets, info)
			}
		}
	}

	// 定义表
	if len(defineSheets) > 0 {
		defines := make([]string, 0)
		for _, info := range defineSheets {
			defines = append(defines, info.File, info.Sheet)
		}
		settings.SetDefines(xlsx.ParseDefineSheet(defines...))
	}

	tables := make([]*model.DataTable, 0)

	// 数据表
	if len(dataSheets) > 0 {
		var tableMap = make(map[string][]SheetInfo)
		for _, info := range dataSheets {
			if info.TypeName == "" {
				log.Fatalf("数据表类型名不能为空，id：%v \n", info.ID)
			}
			if _, ok := tableMap[info.TypeName]; !ok {
				tableMap[info.TypeName] = make([]SheetInfo, 0)
			}

			tableMap[info.TypeName] = append(tableMap[info.TypeName], info)
		}

		for _, infos := range tableMap {
			defines := make([]string, 0)
			for _, info := range infos {
				defines = append(defines, info.File, info.Sheet)
			}
			var table = xlsx.ParseDataSheet(defines...)
			if table != nil {
				table.Id = infos[0].ID
				table.TypeName = infos[0].TypeName
				tables = append(tables, table)
			}
		}
	}

	// 语言表
	if len(langSheets) > 0 {
		defines := make([]string, 0)
		for _, info := range langSheets {
			defines = append(defines, info.File, info.Sheet)
		}

		table := xlsx.ParseDataSheet(defines...)
		table.Id = langSheets[0].ID
		table.IsLanguage = true
		tables = append(tables, table)
	}

	sort.Sort(model.DataTables(tables))

	settings.SetTables(tables...)

	// 执行导出任务
	generator.Build(exportInfo.Type, exportInfo.Path)
}
