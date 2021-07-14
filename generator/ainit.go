package generator

import (
	"bytes"
	"fmt"
	"go-xlsx-protobuf/settings"
	"go-xlsx-protobuf/utils"
	"log"
	"os"
	"text/template"
)

var (
	generators = make(map[string]Generator, 0)

	funcs template.FuncMap
)

type commonFileDesc struct {
	GoProtoVersion string
}

func init() {
	funcs = make(template.FuncMap)

	funcs["upperF"] = func(str string) string {
		if len(str) < 1 {
			return ""
		}
		strArry := []rune(str)
		if strArry[0] >= 97 && strArry[0] <= 122 {
			strArry[0] -= 32
		}
		return string(strArry)
	}

	funcs["space"] = func() string {
		return " "
	}

	funcs["table"] = func() string {
		return "	"
	}

	funcs["add"] = func(a, b int) int {
		return a + b
	}

	funcs["sub"] = func(a, b int) int {
		return a - b
	}

	funcs["join"] = func(strs ...string) string {
		var ret = ""
		for _, str := range strs {
			ret += str
		}
		return ret
	}

	funcs["is_value_type"] = func(valueType string) bool {
		if settings.IsStruct(valueType) || settings.IsTable(valueType) {
			return false
		}
		return true
	}
}

type Generator interface {
	Generate(output string) (save bool, data *bytes.Buffer)
}

func Regist(name string, g Generator) {
	generators[name] = g
}

func Build(typeName, outfile string) bool {
	utils.CheckPath(outfile)

	if gen, ok := generators[typeName]; ok {
		save, code := gen.Generate(outfile)

		if save {
			f, err := os.Create(outfile)
			defer f.Close()
			if err != nil {
				fmt.Errorf(err.Error())
				return false
			}

			_, err = f.WriteString(code.String())
			if err != nil {
				fmt.Errorf(err.Error())
				return false
			}
		}

		return ok
	} else {
		log.Println("[错误] 找不到对应的代码生成器：" + typeName)
	}
	return false
}
