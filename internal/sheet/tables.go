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

	settlementsRowOffset = 2
	settlementsColOffset = 1

	baseStateRowOffset = 2
	baseStateColOffset = 1
)

func newMembersTable(file *excelize.File) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    membersSheet,
		RowOffset:    membersRowOffset,
		ColumnOffset: membersColOffset,
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
		ColumnCount:  membersCount + 1,
		ErrorHandler: fatalIfNotNil,
	}
}

func newSettlementsTable(file *excelize.File) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    settlementsSheet,
		RowOffset:    settlementsRowOffset,
		ColumnOffset: settlementsColOffset,
		ColumnCount:  3,
		ErrorHandler: fatalIfNotNil,
	}
}

func newBaseStateTable(file *excelize.File, membersCount int) *table.Table {
	return &table.Table{
		File:         file,
		SheetName:    baseStateSheet,
		RowOffset:    baseStateRowOffset,
		ColumnOffset: baseStateColOffset,
		ColumnCount:  membersCount + 1,
		ErrorHandler: fatalIfNotNil,
	}
}
