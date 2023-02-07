package table_test

import (
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/table"
	"testing"
)

func TestGetCell(t *testing.T) {
	tableStruct := table.Table{}
	assert(t, tableStruct.GetCell(1, 1) == "A1")
	assert(t, tableStruct.GetCell(2, 5) == "E2")
	assert(t, tableStruct.GetCell(3, 9) == "I3")
	assert(t, tableStruct.GetCell(4, 19) == "S4")
	assert(t, tableStruct.GetCell(100, 27) == "AA100")
	assert(t, tableStruct.GetCell(12, 28) == "AB1")
	assert(t, tableStruct.GetCell(1, 52) == "AZ1")
}

func assert(t *testing.T, b bool) {
	if !b {
		t.Failed()
	}
}
