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

type RowWriter func(rowNumber int, values, formulas []string)
type HeaderWriter func(values []string)
type StyleFounder func(rowNumber, columnNumber int, value string) (int, bool)

type WriteRowsParams struct {
	HeaderWriter HeaderWriter
	RowWriter    RowWriter
	StyleFounder StyleFounder
}

func (t *Table) WriteRows(params WriteRowsParams) {
	columnValues := make([]string, t.ColumnCount)
	columnFormulas := make([]string, t.ColumnCount)

	err := t.File.SetColWidth(t.SheetName, t.getColumn(1), t.getColumn(t.ColumnCount), t.ColumnWidth)
	if err != nil {
		t.ErrorHandler(err)
	}

	params.HeaderWriter(columnValues)
	t.writeRowCells(-1, columnValues, columnFormulas, params.StyleFounder)
	fillWithZero(columnValues)
	fillWithZero(columnFormulas)

	for r := 0; r < t.RowCount; r++ {
		params.RowWriter(r, columnValues, columnFormulas)
		t.writeRowCells(r, columnValues, columnFormulas, params.StyleFounder)
		fillWithZero(columnValues)
		fillWithZero(columnFormulas)
	}
}

func (t *Table) writeRowCells(row int, values, formulas []string, styleFounder StyleFounder) {
	for i := 0; i < len(values); i++ {
		cell := t.getCell(row, i)

		err := t.File.SetCellValue(t.SheetName, cell, values[i])
		if err != nil {
			t.ErrorHandler(err)
		}

		err = t.File.SetCellFormula(t.SheetName, cell, formulas[i])
		if err != nil {
			t.ErrorHandler(err)
		}

		style, hasStyle := styleFounder(row, i, values[i])
		if hasStyle {
			err = t.File.SetCellStyle(t.SheetName, cell, cell, style)
			if err != nil {
				t.ErrorHandler(err)
			}
		}
	}
}

func fillWithZero(values []string) {
	for i := range values {
		values[i] = ""
	}
}

func (t *Table) getCell(rowN, colN int) string {
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
