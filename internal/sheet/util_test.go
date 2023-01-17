package sheet

import "testing"

func TestCell(t *testing.T) {
	assert(t, cell(1, 1) == "A1")
	assert(t, cell(2, 5) == "E2")
	assert(t, cell(3, 9) == "I3")
	assert(t, cell(4, 19) == "S4")
	assert(t, cell(100, 27) == "AA100")
	assert(t, cell(12, 28) == "AB1")
	assert(t, cell(1, 52) == "AZ1")
}

func assert(t *testing.T, b bool) {
	if !b {
		t.Failed()
	}
}
