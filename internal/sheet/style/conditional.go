package style

import (
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/table"
	"github.com/xuri/excelize/v2"
)

type AlternateConditional struct {
	modOffset                            int
	condStyles                           []int
	omitDiagonal                         bool
	diagonalRowOffset, diagonalColOffset int
	startRow, startCol                   int
	endRow, endCol                       int
}

func Alternate(styles ...int) *AlternateConditional {
	ac := new(AlternateConditional)
	ac.condStyles = styles
	return ac
}

func (ac *AlternateConditional) reset() {
	*ac = AlternateConditional{}
}

func (ac *AlternateConditional) mod() int {
	return len(ac.condStyles)
}

func (ac *AlternateConditional) Build() []*table.ConditionalStyle {
	var result []*table.ConditionalStyle
	for i := 0; i < ac.mod(); i++ {
		result = append(result, &table.ConditionalStyle{
			StartRow: ac.startRow,
			StartCol: ac.startCol,
			EndRow:   ac.endRow,
			EndCol:   ac.endCol,
			Options: []excelize.ConditionalFormatOptions{
				{
					Type:     "formula",
					Criteria: ac.createCriteria((i + ac.modOffset) % ac.mod()),
					Format:   ac.condStyles[i%ac.mod()],
				},
			},
		})
	}
	defer ac.reset()
	return result
}

func (ac *AlternateConditional) WithStart(row, column int) *AlternateConditional {
	ac.startRow = row
	ac.startCol = column
	return ac
}

func (ac *AlternateConditional) WithEnd(row, column int) *AlternateConditional {
	ac.endRow = row
	ac.endCol = column
	return ac
}

func (ac *AlternateConditional) WithModOffset(offset int) *AlternateConditional {
	ac.modOffset = offset
	return ac
}

func (ac *AlternateConditional) OmitDiagonal(rowOffset, columnOffset int) *AlternateConditional {
	ac.diagonalRowOffset = rowOffset
	ac.diagonalColOffset = columnOffset
	ac.omitDiagonal = true
	return ac
}

func (ac *AlternateConditional) createCriteria(equals int) string {
	if ac.omitDiagonal {
		return fmt.Sprintf("=AND(MOD(ROW(), %d)=%d,(ROW()-%d)<>(COLUMN()-%d))",
			ac.mod(), equals, ac.diagonalRowOffset, ac.diagonalColOffset)
	}
	return fmt.Sprintf("=MOD(ROW(), %d)=%d", ac.mod(), equals)
}
