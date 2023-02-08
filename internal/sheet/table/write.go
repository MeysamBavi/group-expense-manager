package table

type StylerFunc func(n int) (int, bool)

type WriteRowsParams struct {
	HeaderWriter func(cells []*WCell, mergeCount *int)
	RowWriter    func(rowNumber int, cells []*WCell)
	ColumnWidth  float64
	RowCount     int
	ColumnStyler StylerFunc
	RowStyler    StylerFunc
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
