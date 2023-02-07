package table

type WriteRowsParams struct {
	HeaderWriter func(cells []*WCell, mergeCount *int)
	RowWriter    func(rowNumber int, cells []*WCell)
	ColumnWidth  float64
	RowCount     int
}

type WCell struct {
	Value   any
	Style   *int
	Formula string
}

func (c *WCell) reset() {
	c.Value = nil
	c.Style = nil
	c.Formula = ""
}

func resetWCells(cells []*WCell) {
	for _, c := range cells {
		c.reset()
	}
}
