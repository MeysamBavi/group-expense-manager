package sheet

import (
	"errors"
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
)

func panicE(err error) {
	if err != nil {
		panic(err)
	}
}

func cell(rowN, colN int) string {
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
	return fmt.Sprintf("%s!%s", membersSheet, cell(int(id+2), 1))
}
