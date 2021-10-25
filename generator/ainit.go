package generator

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/wingcd/go-xlsx-exporter/utils"
)

var (
	generators = make(map[string]Generator, 0)

	funcs template.FuncMap
)

// Copied from golint
var commonInitialisms = []string{"ACL", "API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "TCP", "TLS", "TTL", "UDP", "UI", "UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XMPP", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer
var uncommonInitialismsReplacer *strings.Replacer

type commonFileDesc struct {
	GoProtoVersion string
}

func init() {
	var commonInitialismsForReplacer []string
	var uncommonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
		uncommonInitialismsForReplacer = append(uncommonInitialismsForReplacer, strings.Title(strings.ToLower(initialism)), initialism)
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
	uncommonInitialismsReplacer = strings.NewReplacer(uncommonInitialismsForReplacer...)

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

	// 驼峰命名
	funcs["camel_case"] = func(name string) string {
		if name == "" {
			return ""
		}

		temp := strings.Split(name, "_")
		var s string
		for _, v := range temp {
			vv := []rune(v)
			if len(vv) > 0 {
				if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
					vv[0] -= 32
				}
				s += string(vv)
			}
		}

		// s = uncommonInitialismsReplacer.Replace(s)
		//smap.Set(name, s)
		return s
	}

	// 下划线命名
	funcs["under_score_case"] = func(name string) string {
		const (
			lower = false
			upper = true
		)

		if name == "" {
			return ""
		}

		var (
			value                                    = name // commonInitialismsReplacer.Replace(name)
			buf                                      = bytes.NewBufferString("")
			lastCase, currCase, nextCase, nextNumber bool
		)

		for i, v := range value[:len(value)-1] {
			nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
			nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

			if i > 0 {
				if currCase == upper {
					if lastCase == upper && (nextCase == upper || nextNumber == upper) {
						buf.WriteRune(v)
					} else {
						if value[i-1] != '_' && value[i+1] != '_' {
							buf.WriteRune('_')
						}
						buf.WriteRune(v)
					}
				} else {
					buf.WriteRune(v)
					if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
						buf.WriteRune('_')
					}
				}
			} else {
				currCase = upper
				buf.WriteRune(v)
			}
			lastCase = currCase
			currCase = nextCase
		}

		buf.WriteByte(value[len(value)-1])

		s := strings.ToLower(buf.String())
		return s
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
		if utils.IsStruct(valueType) || utils.IsTable(valueType) {
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

	fmt.Printf("启动生成器：%s,生成文件：%s...\n", typeName, outfile)

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
