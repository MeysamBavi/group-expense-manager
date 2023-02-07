package sheet

import (
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/table"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
)

const (
	initialSheetName  = "Sheet1"
	membersSheet      = "members"
	expensesSheet     = "expenses"
	transactionsSheet = "transactions"
	debtMatrixSheet   = "debt-matrix"
	baseStateSheet    = "base state"
	metadataSheet     = "metadata (unmodifiable)"
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

const (
	blockStyle = "block"
	timeStyle  = "time"
)

const (
	timeLayout            = "1/2/06 15:04"
	timeLayoutFormatIndex = 22
)

type Manager struct {
	file               *excelize.File
	members            []*model.Member
	expenses           []*model.Expense
	transactions       []*model.Transaction
	debtMatrix         [][]model.Amount
	baseState          [][]model.Amount
	sheetIndices       map[string]int
	styleIndices       map[string]int
	membersTable       *table.Table
	expensesLeftTable  *table.Table
	expensesRightTable *table.Table
	expensesFullTable  *table.Table
	transactionsTable  *table.Table
	debtMatrixTable    *table.Table
	baseStateTable     *table.Table
}

func newBaseManager() *Manager {
	return &Manager{
		sheetIndices: make(map[string]int),
		styleIndices: make(map[string]int),
	}
}

func NewManager(members []*model.Member) *Manager {
	m := newBaseManager()
	m.file = excelize.NewFile()
	m.members = members
	setMembersTable(m, len(members))
	setTables(m, len(members))

	createStyles(m)
	createSheets(m)

	return m
}

func LoadManager(fileName string) (*Manager, error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	m := newBaseManager()
	m.file = file

	setMembersTable(m, 0)
	m.members = loadMembers(m.membersTable)

	setTables(m, len(m.members))
	m.expenses = loadExpenses(m.expensesFullTable, m.members)
	m.transactions = loadTransactions(m.transactionsTable, m.members)
	m.baseState = loadBaseState(m.baseStateTable, m.members)

	createStyles(m)

	return m, nil
}

func (m *Manager) PrintData() {
	for _, member := range m.members {
		fmt.Println(*member)
	}

	for _, expense := range m.expenses {
		fmt.Println(*expense)
	}

	for _, transaction := range m.transactions {
		fmt.Println(*transaction)
	}

	fmt.Println(m.baseState)
}

func (m *Manager) SaveAs(name string) error {
	return m.file.SaveAs(name)
}

func (m *Manager) MembersCount() int {
	return len(m.members)
}

func (m *Manager) Member(id model.MID) *model.Member {
	return m.members[id]
}

func (m *Manager) SetStyle(key string, value *excelize.Style) {
	si, _ := m.file.NewStyle(value)
	m.styleIndices[key] = si
}

func (m *Manager) GetStyle(key string) int {
	return m.styleIndices[key]
}

func (m *Manager) GetSheetIndex(name string) int {
	return m.sheetIndices[name]
}

func (m *Manager) UpdateDebtors() {
	m.calculateDebtMatrix()
	m.storeDebtMatrix()
}

func (m *Manager) calculateDebtMatrix() {
	// TODO
}

func (m *Manager) storeDebtMatrix() {
	// TODO
}

func createStyles(m *Manager) {
	m.SetStyle(blockStyle, &excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#606060"}, Pattern: 1, Shading: 0},
	})
	m.SetStyle(timeStyle, &excelize.Style{
		NumFmt: timeLayoutFormatIndex,
	})
}

func createSheets(m *Manager) {
	panicE(m.file.SetSheetName(initialSheetName, membersSheet))
	defer initializeMembers(m)

	i, err := m.file.NewSheet(expensesSheet)
	panicE(err)
	m.sheetIndices[expensesSheet] = i
	defer initializeExpenses(m)

	i, err = m.file.NewSheet(transactionsSheet)
	panicE(err)
	m.sheetIndices[transactionsSheet] = i
	defer initializeTransactions(m)

	i, err = m.file.NewSheet(debtMatrixSheet)
	panicE(err)
	m.sheetIndices[debtMatrixSheet] = i
	defer initializeDebtMatrix(m)

	i, err = m.file.NewSheet(baseStateSheet)
	panicE(err)
	m.sheetIndices[baseStateSheet] = i
	defer initializeBaseState(m)

	i, err = m.file.NewSheet(metadataSheet)
	panicE(err)
	m.sheetIndices[metadataSheet] = i
	defer initializeMetadata(m)
}

func setMembersTable(m *Manager, membersCount int) {
	m.membersTable = &table.Table{
		File:         m.file,
		SheetName:    membersSheet,
		RowOffset:    membersRowOffset,
		ColumnOffset: membersColOffset,
		RowCount:     membersCount,
		ColumnCount:  2,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}
}

func setTables(m *Manager, membersCount int) {

	m.expensesLeftTable = &table.Table{
		File:         m.file,
		SheetName:    expensesSheet,
		RowOffset:    expensesLeftSideRowOffset,
		ColumnOffset: expensesLeftSideColOffset,
		RowCount:     1,
		ColumnCount:  4,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}

	m.expensesRightTable = &table.Table{
		File:         m.file,
		SheetName:    expensesSheet,
		RowOffset:    expensesRightSideRowOffset,
		ColumnOffset: expensesRightSideColOffset,
		RowCount:     2,
		ColumnCount:  membersCount * 2,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}

	m.expensesFullTable = &table.Table{
		File:         m.file,
		SheetName:    expensesSheet,
		RowOffset:    expensesRightSideRowOffset,
		ColumnOffset: expensesLeftSideColOffset,
		RowCount:     -1,
		ColumnCount:  4 + membersCount*2,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}

	m.transactionsTable = &table.Table{
		File:         m.file,
		SheetName:    transactionsSheet,
		RowOffset:    transactionsRowOffset,
		ColumnOffset: transactionsColOffset,
		RowCount:     0,
		ColumnCount:  4,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}

	m.debtMatrixTable = &table.Table{
		File:         m.file,
		SheetName:    debtMatrixSheet,
		RowOffset:    debtMatrixRowOffset,
		ColumnOffset: debtMatrixColOffset,
		RowCount:     m.MembersCount() + 1,
		ColumnCount:  m.MembersCount() + 1,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}

	m.baseStateTable = &table.Table{
		File:         m.file,
		SheetName:    baseStateSheet,
		RowOffset:    baseStateRowOffset,
		ColumnOffset: baseStateColOffset,
		RowCount:     m.MembersCount(),
		ColumnCount:  m.MembersCount() + 1,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}
}

func initializeMembers(m *Manager) {
	m.membersTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, _ *int) {
			cells[0].Value = "Name"
			cells[1].Value = "Card Number"
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = m.members[rowNumber].Name
			cells[1].Value = m.members[rowNumber].CardNumber
		},
		ColumnWidth: 32,
	})
}

func initializeMetadata(m *Manager) {}

func initializeBaseState(m *Manager) {
	m.baseStateTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			for i := range m.members {
				cells[i+1].Value = m.members[i].Name
			}
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = m.members[rowNumber].Name
			for i := 1; i < len(cells); i++ {
				cells[i].Value = 0
			}
		},
		ColumnWidth: 20,
	})
}

func initializeDebtMatrix(m *Manager) {
	m.debtMatrixTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			*mergeCount = m.MembersCount() + 1
			cells[0].Value = "Run 'update' command to update debt matrix. Person in the row should pay to the person in the column."
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			if rowNumber == 0 {
				for i := range m.members {
					cells[i+1].Value = m.members[i].Name
				}
				return
			}
			cells[0].Value = m.members[rowNumber-1].Name
			for i := 1; i < len(cells); i++ {
				cells[i].Value = 0
			}
		},
		ColumnWidth: 20,
	})
}

func initializeTransactions(m *Manager) {

	m.transactionsTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, _ *int) {
			cells[0].Value = "Time"
			cells[1].Value = "Receiver"
			cells[2].Value = "Payer"
			cells[3].Value = "Amount"
		},
		ColumnWidth: 16,
	})
}

func initializeExpenses(m *Manager) {

	var totalAmountCell string
	m.expensesLeftTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			cells[0].Value = "Time"
			cells[1].Value = "Title"
			cells[2].Value = "Payer"
			cells[3].Value = "Total Amount"
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = time.Now()
			cells[1].Value = "food"
			cells[2].Value = m.members[0].Name
			cells[3].Value = 300
			totalAmountCell = m.expensesLeftTable.GetCell(rowNumber, 3)
		},
		ColumnWidth: 18,
	})

	var weightCells []string
	m.expensesRightTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			*mergeCount = 2
			for i, v := range m.members {
				cells[i*2].Value = v.Name
			}
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			for i := 0; i < m.MembersCount()*2; i += 2 {
				if rowNumber == 0 {
					cells[i].Value = "Share Weight"
					weightCells = append(weightCells, m.expensesRightTable.GetCell(1, i))
					cells[i+1].Value = "Share Amount"
				} else if rowNumber == 1 {
					cells[i].Value = i >> 2
					totalWeightsFormula := fmt.Sprintf("SUM(%s)", strings.Join(weightCells, ", "))
					cells[i+1].Formula = fmt.Sprintf("(%s/%s)*%s", m.expensesRightTable.GetCell(rowNumber, i), totalWeightsFormula, totalAmountCell)
				}
			}
		},
		ColumnWidth: 14,
	})
}

func loadMembers(t *table.Table) []*model.Member {
	var members []*model.Member
	t.ReadRows(table.ReadRowsParams{
		RowReader: func(rowNumber int, cells []*table.RCell) {
			members = append(members, &model.Member{
				ID:         model.MID(rowNumber),
				Name:       strings.TrimSpace(cells[0].Value),
				CardNumber: strings.TrimSpace(cells[1].Value),
			})
		},
		IncludeHeader:   false,
		UnknownRowCount: true,
	})

	return members
}

func loadExpenses(t *table.Table, members []*model.Member) []*model.Expense {
	var expenses []*model.Expense
	t.ReadRows(table.ReadRowsParams{
		RowReader: func(rowNumber int, cells []*table.RCell) {
			if rowNumber == -1 {
				for i := 4; i < t.ColumnCount; i += 2 {
					requireMemberValidity(members, cells[i].Value, (i-4)/2)
				}
				return
			}

			if rowNumber == 0 {
				return
			}

			theTime, err := time.ParseInLocation(timeLayout, cells[0].Value, time.Local)
			panicE(err)

			title := cells[1].Value

			payer := cells[2].Value
			requireMemberPresence(members, payer)
			payerIndex := findMemberIndex(members, payer)

			amount, err := model.ParseAmount(cells[3].Value)
			panicE(err)

			ex := &model.Expense{
				Title:   title,
				Time:    theTime,
				PayerID: model.MID(payerIndex),
				Amount:  amount,
			}

			var shares []model.Share
			for i := 4; i < t.ColumnCount; i += 2 {
				weight, err := strconv.Atoi(cells[i].Value)
				panicE(err)
				shares = append(shares, model.Share{
					MemberID:    model.MID((i - 4) / 2),
					ShareWeight: weight,
				})
			}

			ex.Shares = shares
			expenses = append(expenses, ex)
		},
		IncludeHeader:   true,
		UnknownRowCount: true,
	})

	return expenses
}

func loadTransactions(t *table.Table, members []*model.Member) []*model.Transaction {

	var transactions []*model.Transaction
	t.ReadRows(table.ReadRowsParams{
		RowReader: func(rowNumber int, cells []*table.RCell) {
			theTime, err := time.ParseInLocation(timeLayout, cells[0].Value, time.Local)
			panicE(err)

			receiver := cells[1].Value
			requireMemberPresence(members, receiver)

			payer := cells[2].Value
			requireMemberPresence(members, payer)

			amount, err := model.ParseAmount(cells[3].Value)
			panicE(err)

			transactions = append(transactions, &model.Transaction{
				Time:       theTime,
				ReceiverID: model.MID(findMemberIndex(members, receiver)),
				PayerID:    model.MID(findMemberIndex(members, payer)),
				Amount:     amount,
			})
		},
		IncludeHeader:   false,
		UnknownRowCount: true,
	})

	return transactions
}

func loadBaseState(t *table.Table, members []*model.Member) [][]model.Amount {

	baseState := make([][]model.Amount, len(members))
	t.ReadRows(table.ReadRowsParams{
		RowReader: func(rowNumber int, cells []*table.RCell) {
			if rowNumber == -1 {
				for i := 0; i < len(members); i++ {
					requireMemberValidity(members, cells[i+1].Value, i)
				}
				return
			}
			requireMemberValidity(members, cells[0].Value, rowNumber)

			baseState[rowNumber] = make([]model.Amount, len(members))
			for i := 0; i < len(members); i++ {
				amount, err := model.ParseAmount(cells[i+1].Value)
				panicE(err)
				baseState[rowNumber][i] = amount
			}
		},
		IncludeHeader:   true,
		UnknownRowCount: false,
	})

	return baseState
}

func requireMemberValidity(members []*model.Member, memberName string, index int) {
	mIndex := findMemberIndex(members, memberName)
	if mIndex < 0 || index != mIndex {
		panic(fmt.Errorf("found no member with name %q and index %d", memberName, index))
	}
}

func requireMemberPresence(members []*model.Member, memberName string) {
	if findMemberIndex(members, memberName) < 0 {
		panic(fmt.Errorf("found no member with name %q", memberName))
	}
}
