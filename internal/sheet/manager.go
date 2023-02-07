package sheet

import (
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/log"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/store"
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
	blockStyle = "block"
)

const (
	timeLayout = "2006/01/02 15:04"
)

type Manager struct {
	file               *excelize.File
	members            *store.MemberStore
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

func NewManager(memberStore *store.MemberStore) *Manager {
	m := newBaseManager()
	m.file = excelize.NewFile()
	m.members = memberStore

	m.membersTable = newMembersTable(m.file, m.MembersCount())
	setTablesExceptMembers(m)

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

	m.membersTable = newMembersTable(m.file, 0)
	m.members = loadMembers(m.membersTable)
	setTablesExceptMembers(m)

	m.expenses = loadExpenses(m.expensesFullTable, m.members)
	m.transactions = loadTransactions(m.transactionsTable, m.members)
	m.baseState = loadBaseState(m.baseStateTable, m.members)

	createStyles(m)

	return m, nil
}

func newBaseManager() *Manager {
	return &Manager{
		sheetIndices: make(map[string]int),
		styleIndices: make(map[string]int),
	}
}

func (m *Manager) PrintData() {
	m.members.Range(func(index int, member *model.Member) {
		fmt.Println(*member)
	})

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
	return m.members.Count()
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
	debtMatrix := m.baseState
	m.baseState = emptyMatrix(m.MembersCount())

	for _, expense := range m.expenses {
		payerIndex := m.members.GetIndexByName(expense.PayerName)
		for _, share := range expense.Shares {
			memberIndex := m.members.GetIndexByName(share.MemberName)
			debtMatrix[memberIndex][payerIndex] =
				debtMatrix[memberIndex][payerIndex].Add(
					expense.Amount.Divide(expense.SumOfWeights()).Multiply(share.ShareWeight))
		}
	}

	for _, transaction := range m.transactions {
		receiverIndex := m.members.GetIndexByName(transaction.ReceiverName)
		payerIndex := m.members.GetIndexByName(transaction.PayerName)
		debtMatrix[payerIndex][receiverIndex] =
			debtMatrix[payerIndex][receiverIndex].Sub(transaction.Amount)
	}

	for r := 0; r < m.MembersCount(); r++ {
		for c := r; c < m.MembersCount(); c++ {
			if debtMatrix[r][c].LessThan(debtMatrix[c][r]) {
				debtMatrix[c][r] = debtMatrix[c][r].Sub(debtMatrix[r][c])
				debtMatrix[r][c] = model.AmountZero
			} else {
				debtMatrix[r][c] = debtMatrix[r][c].Sub(debtMatrix[c][r])
				debtMatrix[c][r] = model.AmountZero
			}
		}
	}

	m.debtMatrix = debtMatrix
}

func (m *Manager) storeDebtMatrix() {
	// TODO
}

func setTablesExceptMembers(m *Manager) {
	m.expensesLeftTable = newExpensesLeftTable(m.file)
	m.expensesRightTable = newExpensesRightTable(m.file, m.MembersCount())
	m.expensesFullTable = newExpensesFullTable(m.file, m.MembersCount())
	m.transactionsTable = newTransactionsTable(m.file)
	m.debtMatrixTable = newDebtMatrixTable(m.file, m.MembersCount())
	m.baseStateTable = newBaseStateTable(m.file, m.MembersCount())
}

func createStyles(m *Manager) {
	m.SetStyle(blockStyle, &excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#606060"}, Pattern: 1, Shading: 0},
	})
}

func createSheets(m *Manager) {
	fatalIfNotNil(m.file.SetSheetName(initialSheetName, membersSheet))
	defer initializeMembers(m)

	i, err := m.file.NewSheet(expensesSheet)
	fatalIfNotNil(err)
	m.sheetIndices[expensesSheet] = i
	defer initializeExpenses(m)

	i, err = m.file.NewSheet(transactionsSheet)
	fatalIfNotNil(err)
	m.sheetIndices[transactionsSheet] = i
	defer initializeTransactions(m)

	i, err = m.file.NewSheet(debtMatrixSheet)
	fatalIfNotNil(err)
	m.sheetIndices[debtMatrixSheet] = i
	defer initializeDebtMatrix(m)

	i, err = m.file.NewSheet(baseStateSheet)
	fatalIfNotNil(err)
	m.sheetIndices[baseStateSheet] = i
	defer initializeBaseState(m)

	i, err = m.file.NewSheet(metadataSheet)
	fatalIfNotNil(err)
	m.sheetIndices[metadataSheet] = i
	defer initializeMetadata(m)
}

func initializeMembers(m *Manager) {
	m.membersTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, _ *int) {
			cells[0].Value = "Name"
			cells[1].Value = "Card Number"
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = m.members.RequireMemberByIndex(rowNumber).Name
			cells[1].Value = m.members.RequireMemberByIndex(rowNumber).CardNumber
		},
		ColumnWidth: 32,
	})
}

func initializeMetadata(m *Manager) {}

func initializeBaseState(m *Manager) {
	m.baseStateTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			m.members.Range(func(i int, member *model.Member) {
				cells[i+1].Value = member.Name
			})
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = m.members.RequireMemberByIndex(rowNumber).Name
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
				m.members.Range(func(i int, member *model.Member) {
					cells[i+1].Value = member.Name
				})
				return
			}
			cells[0].Value = m.members.RequireMemberByIndex(rowNumber - 1).Name
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
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = time.Date(2012, time.June, 26, 5, 6, 0, 0, time.Local).
				Format(timeLayout)
			cells[1].Value = m.members.RequireMemberByIndex(0).Name
			cells[2].Value = m.members.RequireMemberByIndex(1).Name
			cells[3].Value = 223000
		},
		ColumnWidth: 18,
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
			cells[0].Value = time.Now().Format(timeLayout)
			cells[1].Value = "food"
			cells[2].Value = m.members.RequireMemberByIndex(0).Name
			cells[3].Value = 300
			totalAmountCell = m.expensesLeftTable.GetCell(rowNumber, 3)
		},
		ColumnWidth: 18,
	})

	var weightCells []string
	m.expensesRightTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			*mergeCount = 2
			m.members.Range(func(i int, member *model.Member) {
				cells[i*2].Value = member.Name
			})
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

func loadMembers(t *table.Table) *store.MemberStore {
	members := store.NewMemberStore()
	t.ReadRows(table.ReadRowsParams{
		RowReader: func(rowNumber int, cells []*table.RCell) {
			err := members.AddMember(&model.Member{
				Name:       strings.TrimSpace(cells[0].Value),
				CardNumber: strings.TrimSpace(cells[1].Value),
			})
			fatalIfNotNil(err)
		},
		IncludeHeader:   false,
		UnknownRowCount: true,
	})

	return members
}

func loadExpenses(t *table.Table, members *store.MemberStore) []*model.Expense {
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
			fatalIfNotNil(err)

			title := cells[1].Value

			payer := cells[2].Value
			requireMemberPresence(members, payer)

			amount, err := model.ParseAmount(cells[3].Value)
			fatalIfNotNil(err)

			ex := &model.Expense{
				Title:     title,
				Time:      theTime,
				PayerName: payer,
				Amount:    amount,
			}

			var shares []model.Share
			for i := 4; i < t.ColumnCount; i += 2 {
				weight, err := strconv.Atoi(cells[i].Value)
				fatalIfNotNil(err)

				memberName := members.RequireMemberByIndex((i - 4) / 2).Name
				shares = append(shares, model.Share{
					MemberName:  memberName,
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

func loadTransactions(t *table.Table, members *store.MemberStore) []*model.Transaction {

	var transactions []*model.Transaction
	t.ReadRows(table.ReadRowsParams{
		RowReader: func(rowNumber int, cells []*table.RCell) {
			theTime, err := time.ParseInLocation(timeLayout, cells[0].Value, time.Local)
			fatalIfNotNil(err)

			receiver := cells[1].Value
			requireMemberPresence(members, receiver)

			payer := cells[2].Value
			requireMemberPresence(members, payer)

			amount, err := model.ParseAmount(cells[3].Value)
			fatalIfNotNil(err)

			transactions = append(transactions, &model.Transaction{
				Time:         theTime,
				ReceiverName: receiver,
				PayerName:    payer,
				Amount:       amount,
			})
		},
		IncludeHeader:   false,
		UnknownRowCount: true,
	})

	return transactions
}

func loadBaseState(t *table.Table, members *store.MemberStore) [][]model.Amount {

	baseState := emptyMatrix(members.Count())
	t.ReadRows(table.ReadRowsParams{
		RowReader: func(rowNumber int, cells []*table.RCell) {
			if rowNumber == -1 {
				for i := 0; i < members.Count(); i++ {
					requireMemberValidity(members, cells[i+1].Value, i)
				}
				return
			}
			requireMemberValidity(members, cells[0].Value, rowNumber)

			for i := 0; i < members.Count(); i++ {
				amount, err := model.ParseAmount(cells[i+1].Value)
				fatalIfNotNil(err)
				baseState[rowNumber][i] = amount
			}
		},
		IncludeHeader:   true,
		UnknownRowCount: false,
	})

	return baseState
}

func requireMemberValidity(members *store.MemberStore, memberName string, index int) {
	if !members.IsValid(memberName, index) {
		log.FatalErrorByCaller(fmt.Errorf("found no member with name %q and index %d", memberName, index))
	}
}

func requireMemberPresence(members *store.MemberStore, memberName string) {
	if !members.IsPresent(memberName) {
		log.FatalErrorByCaller(fmt.Errorf("found no member with name %q", memberName))
	}
}

func fatalIfNotNil(err error) {
	if err != nil {
		log.FatalErrorByCaller(err)
	}
}

func emptyMatrix(length int) [][]model.Amount {
	result := make([][]model.Amount, length)
	for i := range result {
		result[i] = make([]model.Amount, length)
	}
	return result
}
