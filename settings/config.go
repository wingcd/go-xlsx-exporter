package settings

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

const (
	SheetTypeDefine = "define"
	SheetTypeTable  = "table"
)

type SheetInfo struct {
	ID        int    `yaml:"id"`
	Type      string `yaml:"type"`
	File      string `yaml:"file"`
	Sheet     string `yaml:"sheet"`
	TypeName  string `yaml:"type_name"`
	IsLang    bool   `yaml:"is_lang"`
	Transpose bool   `yaml:"transpose"`
}

type ExportInfo struct {
	ID         int      `yaml:"id"`
	Type       string   `yaml:"type"`
	Path       string   `yaml:"path"`
	Includes   string   `yaml:"includes"`
	Excludes   string   `yaml:"excludes"`
	ExportType int      `yaml:"export_type"`
	Package    string   `yaml:"package"`
	Template   string   `yaml:"template"`
	Imports    []string `yaml:"imports"`
}

type RuleInfo struct {
	ID      int    `yaml:"id"`
	Rule    string `yaml:"rule"`
	Desc    string `yaml:"desc"`
	Disable bool   `yaml:"disable"`

	RRule *regexp.Regexp
}

type YamlConf struct {
	Package        string `yaml:"package"`
	Channel        string `yaml:"channel"`
	PBBytesFileExt string `yaml:"pb_bytes_file_ext"`
	CommentSymbol  string `yaml:"comment_symbol"`
	ExportType     int    `yaml:"export_type"`
	ArraySplitChar string `yaml:"array_split_char"`
	PauseOnEnd     bool   `yaml:"pause_on_end"`
	StrictMode     bool   `yaml:"strict_mode"`

	Rules   []*RuleInfo   `yaml:"rules"`
	Exports []*ExportInfo `yaml:"exports"`
	Sheets  []*SheetInfo  `yaml:"sheets"`
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewYamlConf(filename string) *YamlConf {
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fn := dir + "/" + filename
	fn := filename
	if ok, _ := PathExists(fn); !ok {
		log.Fatalf("配置文件不存在：%v", fn)
	}

	c := new(YamlConf)
	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, r := range c.Rules {
		r.RRule, err = regexp.Compile(r.Rule)
		if err != nil {
			log.Fatalf("非法的规则：id=%v, rule=%v, desc=%v", r.ID, r.Rule, r.Desc)
		}
	}
	return c
}
