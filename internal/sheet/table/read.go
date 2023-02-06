package table

type ReadRowsParams struct {
	RowReader       func(rowNumber int, cells []*RCell)
	IncludeHeader   bool
	UnknownRowCount bool
}

type RCell struct {
	Value   string
	Formula string
}

func (c *RCell) reset() {
	c.Value = ""
	c.Formula = ""
}

func allValuesEmpty(cells []*RCell) bool {
	for _, c := range cells {
		if c.Value != "" {
			return false
		}
	}
	return true
}

func resetRCells(cells []*RCell) {
	for _, c := range cells {
		c.reset()
	}
}
