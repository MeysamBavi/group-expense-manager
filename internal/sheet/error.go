package sheet

import "fmt"

type CellError struct {
	cellName  string
	sheetName string
	err       error
}

func (c CellError) Error() string {
	return fmt.Sprintf("%v: in %q at %q", c.err, c.sheetName, c.cellName)
}

func CellErrorOf(err error, sheetName, cellName string) error {
	if err == nil {
		return nil
	}

	return CellError{err: err, sheetName: sheetName, cellName: cellName}
}
