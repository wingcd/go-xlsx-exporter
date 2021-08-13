package go_xlsx_exporter

import (
	"fmt"
	"reflect"
	"strings"
)

var I18NIndexKey = "ID"

func init() {
	Regist(File___inner_language___proto)
}

type I18N struct {
	DataTable

	instance *I18N
	language string
}

func (t *I18N) Instance() *I18N {
	if t.instance == nil {
		t.instance = new(I18N)
		t.indexKey = I18NIndexKey
		t.dataType = reflect.TypeOf(Language{})
		RegistDataTableExt(t)
	}
	return t.instance
}

func (t *I18N) SetLanguage(lang string) {
	if t.language != lang {
		t.language = lang
		t.Clear()
	}
}

func (t *I18N) Language() string {
	return t.language
}

func (t *I18N) GetFilename(typeName string) string {
	return fmt.Sprintf("%s%s.%s%s", DataDir, strings.ToLower(typeName), strings.ToLower(t.language), BytesFileExt)
}

func (t *I18N) Translate(id string) string {
	var itemsMap = t.ItemsMap()
	if text, ok := itemsMap[id]; ok {
		return text.(*Language).Text
	}
	fmt.Printf("no %s language item id=%s", t.language, id)
	return ""
}
