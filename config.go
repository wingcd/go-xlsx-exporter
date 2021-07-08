package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	SheetTypeDefine = "define"
	SheetTypeTable  = "table"
)

var (
	Config *YamlConf
)

func init() {
	Config = NewYamlConf("conf.yaml")
}

type SheetInfo struct {
	Type           string `yaml:"type"`
	File           string `yaml:"xls_file"`
	Sheet          string `yaml:"sheet"`
	TargetCodeFile string `yaml:"target_code_file"`
	TargetDataFile string `yaml:"target_data_file"`
	GenCode        bool   `yaml:"gen_code"`
	GenData        bool   `yaml:"gen_data"`
}

type YamlConf struct {
	Package       string
	TargetCodeDir string
	TargetDataDir string
	Sheets        []SheetInfo `yaml:"sheets"`
}

func NewYamlConf(filename string) *YamlConf {
	c := new(YamlConf)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}
