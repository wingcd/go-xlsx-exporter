package utils

type PureTable struct {
	Data        [][]string
	ColumnCount int
	RowCount    int
}

// colum first
func NewPureTable(data [][]string) *PureTable {
	t := new(PureTable)
	t.Data = data
	return t
}

func (t *PureTable) GetColumn(col int) []string {
	return t.Data[col]
}

func (t *PureTable) GetRow(row int) []string {
	rowValue := make([]string, 0)
	for _, col := range t.Data {
		rowValue = append(rowValue, col[row])
	}
	return rowValue
}

func (t *PureTable) Get(col, row int) string {
	return t.Data[col][row]
}

func (t *PureTable) RemoveColumn(col int) {
	if len(t.Data) <= col || col < 0 {
		return
	}

	t.Data = append(t.Data[:col], t.Data[col+1:]...)
}

func (t *PureTable) RemoveRow(row int) {
	for ci, col := range t.Data {
		if len(col) > row && row >= 0 {
			t.Data[ci] = append(col[:row], col[row+1:]...)
		}
	}
}

func (t *PureTable) ToRows(row int) {

}
