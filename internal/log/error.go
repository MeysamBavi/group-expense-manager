package log

import "fmt"

func CellErrorOf(err error, sheetName, cellName string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: in %q at %q", err, sheetName, cellName)
}

func SheetErrorOf(err error, sheetName string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: in %q", err, sheetName)
}
