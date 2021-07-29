package go_xlsx_exporter

import (
	"reflect"
)

type DataArray map[string]interface{}

var (
	DataDir = ""

	dataArrays = make(map[string]Table)
)

type Table struct {
	data DataArray
}

func (t *Table) Load(dataType reflect.Type) {
	if t.data != nil {
		return
	}

	t.data = make(DataArray)
}
