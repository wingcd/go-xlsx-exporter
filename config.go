package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	SheetTypeDefine = "define"
	SheetTypeTable  = "table"
)

type SheetInfo struct {
	ID       int    `yaml:"id"`
	Type     string `yaml:"type"`
	File     string `yaml:"xls_file"`
	Sheet    string `yaml:"sheet"`
	TypeName string `yaml:"type_name"`
	IsLang   bool   `yaml:"is_lang"`
}

type ExportInfo struct {
	ID         int    `yaml:"id"`
	Type       string `yaml:"type"`
	Path       string `yaml:"path"`
	Sheets     string `yaml:"sheets"`
	ExportType int    `yaml:"export_type"`
	Package    string `yaml:"package"`
}

type YamlConf struct {
	Package        string `yaml:"package"`
	PBBytesFileExt string `yaml:"pb_bytes_file_ext"`
	CommentSymbol  string `yaml:"comment_symbol"`
	ExportType     int    `yaml:"export_type"`
	ArraySplitChar string `yaml:"array_split_char"`
	PauseOnEnd     bool   `yaml:"pause_on_end"`
	StrictMode     bool   `yaml:"strict_mode"`

	Exports []ExportInfo `yaml:"exports"`
	Sheets  []SheetInfo  `yaml:"sheets"`
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	fn := dir + "/" + filename
	if ok, _ := PathExists(fn); !ok {
		log.Fatalf("配置文件不存在：%v", fn)
	}

	c := new(YamlConf)
	yamlFile, err := ioutil.ReadFile(dir + "/" + filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}
