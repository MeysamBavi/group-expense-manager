package table

type WriteRowsParams struct {
	HeaderWriter func(cells []*WCell, mergeCount *int)
	RowWriter    func(rowNumber int, cells []*WCell)
	ColumnWidth  float64
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

func (w *WriteRowsParams) callHeaderWriter(cells []*WCell, mergeCount *int) {
	if w.HeaderWriter != nil {
		w.HeaderWriter(cells, mergeCount)
	}
}

func (w *WriteRowsParams) callRowWriter(rowNumber int, cells []*WCell) {
	if w.RowWriter != nil {
		w.RowWriter(rowNumber, cells)
	}
}
