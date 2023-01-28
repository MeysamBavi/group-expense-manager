package sheet

import (
	"errors"
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
)

var (
	offsetIsEnable = true
	rowOffset      = 0
	colOffset      = 0
	rowOffsetTemp  = 0
	colOffsetTemp  = 0
)

func setOffsets(r, c int) {
	rowOffset = r
	colOffset = c
	offsetIsEnable = true
}

func resetOffsets() {
	rowOffset = 0
	colOffset = 0
}

func disableOffsets() {
	if !offsetIsEnable {
		return
	}
	rowOffsetTemp, colOffsetTemp = rowOffset, colOffset
	resetOffsets()
	offsetIsEnable = false
}

func enableOffsets() {
	if offsetIsEnable {
		return
	}
	setOffsets(rowOffsetTemp, colOffsetTemp)
	offsetIsEnable = true
}

func panicE(err error) {
	if err != nil {
		panic(err)
	}
}

func cell(rowN, colN int) string {
	rowN += rowOffset
	colN += colOffset

	if rowN <= 0 || colN <= 0 {
		panic(errors.New("row number and column number must be positive"))
	}

	return column(colN) + row(rowN)
}

func row(rowN int) string {
	if rowN <= 0 {
		panic(errors.New("row number must be positive"))
	}
	return fmt.Sprintf("%d", rowN)
}

func column(colN int) string {
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

func memberNameRef(id model.MID) string {
	disableOffsets()
	defer enableOffsets()
	return fmt.Sprintf("%s!%s", membersSheet, cell(int(id+2), 1))
}

func findMemberIndex(members []*model.Member, name string) int {
	for i, m := range members {
		if m.Name == name {
			return i
		}
	}
	return -1
}
