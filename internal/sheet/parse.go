package sheet

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ParseSheet(fileName string) {
	file, err := excelize.OpenFile(fileName)
	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		return
	}

	value, err := file.GetCellValue("sheet1", "A1")
	if err != nil {
		fmt.Printf("failed to get cell value: %v\n", err)
		return
	}

	fmt.Println(value)
}
