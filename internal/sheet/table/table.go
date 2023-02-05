package table

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

// Table rows and columns numbers start from 0. Offset values are added to them to be used as a sheet cell name.
type Table struct {
	File                    *excelize.File
	SheetName               string
	RowOffset, ColumnOffset int
	RowCount, ColumnCount   int
	ColumnWidth             float64
	ErrorHandler            func(error)
}

func (t *Table) callErrorHandler(err error) {
	if err != nil && t.ErrorHandler != nil {
		t.ErrorHandler(err)
	}
}

type Cell struct {
	Value   any
	Style   *int
	Formula string
}

func (c *Cell) reset() {
	c.Value = nil
	c.Style = nil
	c.Formula = ""
}

type RowWriter func(rowNumber int, cells []*Cell)
type HeaderWriter func(cells []*Cell, mergeCount *int)

type WriteRowsParams struct {
	HeaderWriter HeaderWriter
	RowWriter    RowWriter
}

func (w *WriteRowsParams) callHeaderWriter(cells []*Cell, mergeCount *int) {
	if w.HeaderWriter != nil {
		w.HeaderWriter(cells, mergeCount)
	}
}

func (w *WriteRowsParams) callRowWriter(rowNumber int, cells []*Cell) {
	if w.RowWriter != nil {
		w.RowWriter(rowNumber, cells)
	}
}

func (t *Table) WriteRows(params WriteRowsParams) {
	cells := make([]*Cell, t.ColumnCount)
	for i := range cells {
		cells[i] = &Cell{}
	}

	err := t.File.SetColWidth(t.SheetName, t.getColumn(t.ColumnOffset), t.getColumn(t.ColumnCount+t.ColumnOffset-1), t.ColumnWidth)
	t.callErrorHandler(err)

	mergeCount := 1
	params.callHeaderWriter(cells, &mergeCount)
	t.writeRowCells(-1, cells, mergeCount)
	resetCells(cells)

	for r := 0; r < t.RowCount; r++ {
		params.callRowWriter(r, cells)
		t.writeRowCells(r, cells, 1)
		resetCells(cells)
	}
}

func (t *Table) writeRowCells(row int, cells []*Cell, multiplier int) {
	var err error
	for i := 0; i < len(cells); i++ {
		cell := t.GetCell(row, i*multiplier)
		err = t.File.MergeCell(t.SheetName, cell, t.GetCell(row, i*multiplier+multiplier-1))
		t.callErrorHandler(err)

		err = t.File.SetCellValue(t.SheetName, cell, cells[i].Value)
		t.callErrorHandler(err)

		err = t.File.SetCellFormula(t.SheetName, cell, cells[i].Formula)
		t.callErrorHandler(err)

		if cells[i].Style != nil {
			style := *cells[i].Style
			err = t.File.SetCellStyle(t.SheetName, cell, cell, style)
			t.callErrorHandler(err)
		}
	}
}

func resetCells(cells []*Cell) {
	for _, c := range cells {
		c.reset()
	}
}

func (t *Table) GetCell(rowN, colN int) string {
	rowN += t.RowOffset
	colN += t.ColumnOffset

	if rowN <= 0 || colN <= 0 {
		panic(errors.New("row number and column number must be positive"))
	}

	return t.getColumn(colN) + t.getRow(rowN)
}

func (t *Table) getRow(rowN int) string {
	return fmt.Sprintf("%d", rowN)
}

func (t *Table) getColumn(colN int) string {

	colDigits := make([]rune, 0, 1)
	colN -= 1
	colDigits = append(colDigits, 'A'+rune(colN%26))
	colN /= 26

	for colN > 0 {
		colN -= 1
		colDigits = append(colDigits, 'A'+rune(colN%26))
		colN /= 26
	}

	for i := 0; i < len(colDigits)>>1; i++ {
		a, b := i, len(colDigits)-1-i
		colDigits[a], colDigits[b] = colDigits[b], colDigits[a]
	}

	return string(colDigits)
}
