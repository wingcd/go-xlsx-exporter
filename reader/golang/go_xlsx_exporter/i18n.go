package go_xlsx_exporter

import (
	"fmt"
	"reflect"
	"strings"
)

var i18nInst *i18n

type i18n struct {
	DataTable

	language string
}

func newI18NInstance() *i18n {
	var t = new(i18n)
	t.dataType = reflect.TypeOf(Language{})
	RegistDataTableExt(t)
	return t
}

func SetLanguage(indexKey, lang string) {
	if i18nInst == nil || i18nInst.language != lang {
		i18nInst = newI18NInstance()
		i18nInst.indexKey = indexKey
		i18nInst.language = lang
		i18nInst.fileGen = new(i18nFileGenerator)
		i18nInst.Clear()
	}
}

type i18nFileGenerator struct {
}

func (t *i18nFileGenerator) GetFilename(typeName string) string {
	return fmt.Sprintf("%s%s.%s%s", dataDir, strings.ToLower(typeName), strings.ToLower(i18nInst.language), BytesFileExt)
}

func GetLanguage() string {
	return i18nInst.language
}

func Translate(id string) string {
	var itemsMap = i18nInst.ItemsMap()
	if text, ok := itemsMap[id]; ok {
		return text.(*Language).Text
	}
	fmt.Printf("no %s language item id=%s", i18nInst.language, id)
	return ""
}
