package generator

import (
	"bytes"
	"fmt"
	"go-xlsx-protobuf/utils"
	"log"
	"os"
)

var (
	generators = make(map[string]Generator, 0)
)

type Generator interface {
	Generate() *bytes.Buffer
}

func Regist(name string, g Generator) {
	generators[name] = g
}

func Build(typeName, outfile string) bool {
	utils.CheckPath(outfile)

	if gen, ok := generators[typeName]; ok {
		code := gen.Generate()

		f, err := os.Create(outfile)
		defer f.Close()
		if err != nil {
			fmt.Errorf(err.Error())
		}

		_, err = f.WriteString(code.String())
		if err != nil {
			fmt.Errorf(err.Error())
		}

		return ok
	} else {
		log.Println("[错误] 找不到对应的代码生成器：" + typeName)
	}
	return false
}
