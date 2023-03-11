package table

import (
	"errors"
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/log"
	"github.com/xuri/excelize/v2"
)

// Table rows and columns numbers start from 0. Offset values are added to them to be used as a sheet cell name.
type Table struct {
	File                    *excelize.File
	SheetName               string
	RowOffset, ColumnOffset int
	ColumnCount             int
}

func (t *Table) fatalIfNotNil(err error) {
	if err != nil {
		log.FatalErrorByCaller(fmt.Errorf("%w: in %q", err, t.SheetName))
	}
}

func (t *Table) WriteRows(params WriteRowsParams) {
	cells := make([]*WCell, t.ColumnCount)
	for i := range cells {
		cells[i] = &WCell{}
	}

	if params.ClearBeforeWrite {
		blankSheetIndex, err := t.File.NewSheet("blank")
		t.fatalIfNotNil(err)
		sheetIndex, err := t.File.GetSheetIndex(t.SheetName)
		t.fatalIfNotNil(err)
		err = t.File.CopySheet(blankSheetIndex, sheetIndex)
		t.fatalIfNotNil(err)
		err = t.File.DeleteSheet("blank")
		t.fatalIfNotNil(err)
	}

	if params.ColumnWidth > 0 {
		err := t.File.SetColWidth(t.SheetName,
			t.getColumn(0),
			t.getColumn(t.ColumnCount-1),
			params.ColumnWidth,
		)
		t.fatalIfNotNil(err)
	}

	if params.ColumnStyler != nil {
		for c := 0; c < t.ColumnCount; c++ {
			if style, ok := params.ColumnStyler(c); ok {
				err := t.File.SetColStyle(t.SheetName, t.getColumn(c), style)
				t.fatalIfNotNil(err)
			}
		}
	}

	if params.RowStyler != nil {
		for r := -1; r < params.RowCount; r++ {
			if style, ok := params.RowStyler(r); ok {
				err := t.File.SetCellStyle(t.SheetName, t.GetCell(r, 0), t.GetCell(r, t.ColumnCount-1), style)
				t.fatalIfNotNil(err)
			}
		}
	}

	mergeCount := 1
	if params.HeaderWriter != nil {
		params.HeaderWriter(cells, &mergeCount)
		t.writeRowCells(-1, cells, mergeCount)
		resetWCells(cells)
	}

	for r := 0; params.RowWriter == nil || r < params.RowCount; r++ {
		params.RowWriter(r, cells)
		t.writeRowCells(r, cells, 1)
		resetWCells(cells)
	}

	for _, condStyle := range params.ConditionalStyles {
		rangeRef := fmt.Sprintf("%s:%s",
			t.GetCell(condStyle.StartRow, condStyle.StartCol), t.GetCell(condStyle.EndRow, condStyle.EndCol))
		err := t.File.SetConditionalFormat(t.SheetName, rangeRef, condStyle.Options)
		t.fatalIfNotNil(err)
	}
}

func (t *Table) writeRowCells(row int, cells []*WCell, multiplier int) {
	var err error
	for i := 0; i < len(cells); i++ {
		cell := t.GetCell(row, i)
		if i%multiplier != 0 {
			continue
		}

		err = t.File.MergeCell(t.SheetName, cell, t.GetCell(row, i+multiplier-1))
		t.fatalIfNotNil(err)

		err = t.File.SetCellValue(t.SheetName, cell, cells[i].Value)
		t.fatalIfNotNil(err)

		err = t.File.SetCellFormula(t.SheetName, cell, cells[i].Formula)
		t.fatalIfNotNil(err)

		if cells[i].Style != nil {
			style := *cells[i].Style
			err = t.File.SetCellStyle(t.SheetName, cell, cell, style)
			t.fatalIfNotNil(err)
		}
	}
}

func (t *Table) GetCell(rowN, colN int) string {
	return fmt.Sprintf("%s%d", t.getColumn(colN), t.getRow(rowN))
}

func (t *Table) getRow(rowN int) int {
	rowN += t.RowOffset
	if rowN <= 0 {
		panic(errors.New("row number must be positive"))
	}
	return rowN
}

func (t *Table) getColumn(colN int) string {
	colN += t.ColumnOffset
	if colN <= 0 {
		panic(errors.New("column number must be positive"))
	}

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

func (t *Table) ReadRows(params ReadRowsParams) {
	cells := make([]*RCell, t.ColumnCount)
	for i := range cells {
		cells[i] = &RCell{}
	}

	i := 0
	if params.IncludeHeader {
		i = -1
	}
	for ; params.RowReader == nil || params.UnknownRowCount || i < params.RowCount; i++ {
		t.readRowCells(i, cells)
		if allValuesEmpty(cells) {
			break
		}
		params.RowReader(i, cells)
		resetRCells(cells)
	}
}

func (t *Table) readRowCells(row int, cells []*RCell) {
	for i := 0; i < len(cells); i++ {
		cell := t.GetCell(row, i)

		formula, err := t.File.GetCellFormula(t.SheetName, cell)
		t.fatalIfNotNil(err)
		cells[i].Formula = formula

		var value string
		if formula != "" {
			value, err = t.File.CalcCellValue(t.SheetName, cell)
		} else {
			value, err = t.File.GetCellValue(t.SheetName, cell)
		}
		t.fatalIfNotNil(err)
		cells[i].Value = value
	}
}
