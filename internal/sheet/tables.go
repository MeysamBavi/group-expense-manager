package sheet

import (
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/table"
	"github.com/xuri/excelize/v2"
)

const (
	membersRowOffset = 2
	membersColOffset = 1

	expensesLeftSideRowOffset  = 3
	expensesLeftSideColOffset  = 1
	expensesRightSideRowOffset = 2
	expensesRightSideColOffset = 5

	transactionsRowOffset = 2
	transactionsColOffset = 1

	debtMatrixRowOffset = 2
	debtMatrixColOffset = 1

	baseStateRowOffset = 2
	baseStateColOffset = 1
)

func newMembersTable(file *excelize.File, membersCount int) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    membersSheet,
		RowOffset:    membersRowOffset,
		ColumnOffset: membersColOffset,
		RowCount:     membersCount,
		ColumnCount:  2,
		ErrorHandler: fatalIfNotNil,
	}
}

func newExpensesLeftTable(file *excelize.File) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    expensesSheet,
		RowOffset:    expensesLeftSideRowOffset,
		ColumnOffset: expensesLeftSideColOffset,
		RowCount:     1,
		ColumnCount:  4,
		ErrorHandler: fatalIfNotNil,
	}
}

func newExpensesRightTable(file *excelize.File, membersCount int) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    expensesSheet,
		RowOffset:    expensesRightSideRowOffset,
		ColumnOffset: expensesRightSideColOffset,
		RowCount:     2,
		ColumnCount:  membersCount * 2,
		ErrorHandler: fatalIfNotNil,
	}
}

func newExpensesFullTable(file *excelize.File, membersCount int) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    expensesSheet,
		RowOffset:    expensesRightSideRowOffset,
		ColumnOffset: expensesLeftSideColOffset,
		RowCount:     -1,
		ColumnCount:  4 + membersCount*2,
		ErrorHandler: fatalIfNotNil,
	}
}

func newTransactionsTable(file *excelize.File) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    transactionsSheet,
		RowOffset:    transactionsRowOffset,
		ColumnOffset: transactionsColOffset,
		RowCount:     1,
		ColumnCount:  4,
		ErrorHandler: fatalIfNotNil,
	}
}

func newDebtMatrixTable(file *excelize.File, membersCount int) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    debtMatrixSheet,
		RowOffset:    debtMatrixRowOffset,
		ColumnOffset: debtMatrixColOffset,
		RowCount:     membersCount + 1,
		ColumnCount:  membersCount + 1,
		ErrorHandler: fatalIfNotNil,
	}
}

func newBaseStateTable(file *excelize.File, membersCount int) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    baseStateSheet,
		RowOffset:    baseStateRowOffset,
		ColumnOffset: baseStateColOffset,
		RowCount:     membersCount,
		ColumnCount:  membersCount + 1,
		ErrorHandler: fatalIfNotNil,
	}
}
