package sheet

import (
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/log"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/store"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/style"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/table"
	"github.com/xuri/excelize/v2"
	"sort"
	"strings"
	"time"
)

const (
	initialSheetName  = "Sheet1"
	membersSheet      = "members"
	expensesSheet     = "expenses"
	transactionsSheet = "transactions"
	debtMatrixSheet   = "debt matrix"
	settlementsSheet  = "settlements"
	baseStateSheet    = "base state"
	metadataSheet     = "metadata"
)

type Manager struct {
	file               *excelize.File
	members            *store.MemberStore
	expenses           []*model.Expense
	transactions       []*model.Transaction
	debtMatrix         [][]model.Amount
	baseState          [][]model.Amount
	settlements        []*model.Transaction
	styleIndices       map[int]int
	membersTable       *table.Table
	expensesLeftTable  *table.Table
	expensesRightTable *table.Table
	expensesFullTable  *table.Table
	transactionsTable  *table.Table
	debtMatrixTable    *table.Table
	settlementsTable   *table.Table
	baseStateTable     *table.Table
	metadataTable      *table.Table
	theme              *style.Theme
}

func NewManager(memberStore *store.MemberStore, theme *style.Theme) *Manager {
	m := newBaseManager()
	m.file = excelize.NewFile()
	m.members = memberStore
	m.theme = theme

	m.membersTable = newMembersTable(m.file)
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

	m.membersTable = newMembersTable(m.file)
	m.members = loadMembers(m.membersTable)
	setTablesExceptMembers(m)

	m.theme = loadMetadata(m.metadataTable)
	m.expenses = loadExpenses(m.expensesFullTable, m.members)
	m.transactions = loadTransactions(m.transactionsTable, m.members)
	m.baseState = loadBaseState(m.baseStateTable, m.members)

	createStyles(m)

	return m, nil
}

func newBaseManager() *Manager {
	return &Manager{
		styleIndices: make(map[int]int),
	}
}

func (m *Manager) PrintData(summarize bool) {
	fmt.Println("Members:")
	m.members.Range(func(index int, member *model.Member) {
		fmt.Println(*member)
	})

	fmt.Println("Expenses:")
	if summarize {
		fmt.Println("Count:", len(m.expenses))
	} else {
		for _, expense := range m.expenses {
			fmt.Println(*expense)
		}
	}

	fmt.Println("Transactions:")
	if summarize {
		fmt.Println("Count:", len(m.transactions))
	} else {
		for _, transaction := range m.transactions {
			fmt.Println(*transaction)
		}
	}

	fmt.Println("Debt Matrix:")
	fmt.Println(m.debtMatrix)

	fmt.Println("Settlements:")
	for _, settlement := range m.settlements {
		fmt.Println(*settlement)
	}

	fmt.Println("Base State:")
	fmt.Println(m.baseState)
}

func (m *Manager) SaveAs(name string) error {
	err := m.file.SetSheetVisible(metadataSheet, false)
	fatalIfNotNil(err)
	return m.file.SaveAs(name)
}

func (m *Manager) MembersCount() int {
	return m.members.Count()
}

func (m *Manager) UpdateDebtors() {
	m.calculateDebtMatrix()
	m.writeDebtMatrix()
	m.calculateSettlements()
	m.writeSettlements()
}

func (m *Manager) setStyle(key int, value *excelize.Style) {
	si, _ := m.file.NewStyle(value)
	m.styleIndices[key] = si
}

func (m *Manager) setCondStyle(key int, value *excelize.Style) {
	si, _ := m.file.NewConditionalStyle(value)
	m.styleIndices[key] = si
}

func (m *Manager) getStyle(key int) int {
	return m.styleIndices[key]
}

func (m *Manager) calculateDebtMatrix() {
	debtMatrix := copyMatrix(m.baseState)

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
				debtMatrix[r][c] = model.AmountZero()
			} else {
				debtMatrix[r][c] = debtMatrix[r][c].Sub(debtMatrix[c][r])
				debtMatrix[c][r] = model.AmountZero()
			}
		}
	}

	m.debtMatrix = debtMatrix
}

func (m *Manager) writeDebtMatrix() {
	m.debtMatrixTable.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			*mergeCount = m.debtMatrixTable.ColumnCount
			cells[0].Value = "Run 'update' command to update the debt matrix. Person in the row should pay the person in the column."
			cells[0].Style = newInt(m.getStyle(helpStyle))
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			if rowNumber == 0 {
				cells[0].Value = fmt.Sprintf("last update: %s", model.TimeOfGregorian(time.Now()))
				cells[0].Style = newInt(m.getStyle(lastUpdateStyle))
				m.members.Range(func(i int, member *model.Member) {
					cells[i+1].Value = member.Name
					cells[i+1].Style = newInt(m.getStyle(headerBoxStyle))
				})
				return
			}

			memberIndex := rowNumber - 1

			cells[0].Value = m.members.RequireMemberByIndex(memberIndex).Name
			cells[0].Style = newInt(m.getStyle(headerBoxStyle))
			for i := 0; i < m.MembersCount(); i++ {
				amount := m.debtMatrix[memberIndex][i]
				cells[i+1].Value = amount.ToNumeral()
				if amount.IsZero() {
					cells[i+1].Value = ""
				}
				if memberIndex == i {
					cells[i+1].Style = newInt(m.getStyle(alternateBlockStyle))
				} else {
					cells[i+1].Style = newInt(m.getStyle(moneyStyle))
				}
			}
		},
		ColumnWidth: 20,
		RowCount:    m.MembersCount() + 1,
		ConditionalStyles: style.Alternate(m.getStyle(alternate0Style), m.getStyle(alternate1Style)).
			OmitDiagonal(m.debtMatrixTable.RowOffset, m.debtMatrixTable.ColumnOffset).
			WithStart(1, 1).
			WithModOffset(1).
			WithEnd(m.MembersCount(), m.MembersCount()).
			Build(),
	})
}

func (m *Manager) writeBaseState() {
	m.baseStateTable.WriteRows(table.WriteRowsParams{
		RowCount: m.MembersCount(),
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			m.members.Range(func(i int, member *model.Member) {
				cells[i+1].Value = member.Name
				cells[i+1].Style = newInt(m.getStyle(headerBoxStyle))
			})
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = m.members.RequireMemberByIndex(rowNumber).Name
			cells[0].Style = newInt(m.getStyle(headerBoxStyle))
			for i := 0; i < m.MembersCount(); i++ {
				amount := m.baseState[rowNumber][i]
				cells[i+1].Value = amount.ToNumeral()
				if amount.IsZero() {
					cells[i+1].Value = ""
				}
				if rowNumber == i {
					cells[i+1].Style = newInt(m.getStyle(alternateBlockStyle))
				} else {
					cells[i+1].Style = newInt(m.getStyle(moneyStyle))
				}
			}
		},
		ColumnWidth: 20,
		ConditionalStyles: style.Alternate(m.getStyle(alternate0Style), m.getStyle(alternate1Style)).
			OmitDiagonal(0, 0).
			WithStart(0, 1).
			WithEnd(m.MembersCount()-1, m.MembersCount()).
			Build(),
	})
}

func (m *Manager) calculateSettlements() {
	type balance struct {
		name   string
		amount model.Amount // positive means debtor
	}

	balances := make([]balance, 0, m.MembersCount())
	m.members.Range(func(memberIndex int, member *model.Member) {
		gives, receives := model.AmountZero(), model.AmountZero()
		for i := 0; i < m.MembersCount(); i++ {
			receives = receives.Add(m.debtMatrix[i][memberIndex])
		}
		for i := 0; i < m.MembersCount(); i++ {
			gives = gives.Add(m.debtMatrix[memberIndex][i])
		}
		balances = append(balances, balance{
			name:   member.Name,
			amount: gives.Sub(receives),
		})
	})

	sort.SliceStable(balances, func(i, j int) bool {
		return balances[i].amount.LessThan(balances[j].amount)
	})

	settlements := make([]*model.Transaction, 0)
	addSettlement := func(receiver, payer int, amount model.Amount) {
		if amount.IsZero() {
			return
		}
		if amount.LessThan(model.AmountZero()) {
			receiver, payer = payer, receiver
			amount = amount.Negative()
		}
		settlements = append(settlements, &model.Transaction{
			ReceiverName: balances[receiver].name,
			PayerName:    balances[payer].name,
			Amount:       amount,
		})
		balances[payer].amount = balances[payer].amount.Sub(amount)
		balances[receiver].amount = balances[receiver].amount.Add(amount)
	}
	lowest, highest := 0, len(balances)-1
	for highest > lowest {
		if balances[lowest].amount.IsZero() {
			lowest++
			continue
		}
		if balances[highest].amount.IsZero() {
			highest--
			continue
		}

		deficit := balances[highest].amount.Add(balances[lowest].amount)
		if deficit.IsPositive() {
			addSettlement(lowest, highest, balances[lowest].amount.Negative())
		} else {
			addSettlement(lowest, highest, balances[highest].amount)
		}
	}

	sort.SliceStable(settlements, func(i, j int) bool {
		return settlements[i].Amount.LessThan(settlements[j].Amount)
	})
	m.settlements = settlements
}

func (m *Manager) writeSettlements() {
	m.settlementsTable.WriteRows(table.WriteRowsParams{
		RowCount: len(m.settlements) + 1,
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			*mergeCount = m.settlementsTable.ColumnCount
			cells[0].Value = "Run 'update' command to generate settlement transactions."
			cells[0].Style = newInt(m.getStyle(helpStyle))
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			if rowNumber == 0 {
				cells[0].Value = "Receiver"
				cells[1].Value = "Payer"
				cells[2].Value = "Amount"
				return
			}
			rowNumber--
			cells[0].Value = m.settlements[rowNumber].ReceiverName
			cells[1].Value = m.settlements[rowNumber].PayerName
			cells[2].Value = m.settlements[rowNumber].Amount.ToNumeral()
			cells[2].Style = newInt(m.getStyle(moneyStyle))
		},
		ColumnWidth: 18,
		RowStyler: func(row int) (int, bool) {
			if row == 0 {
				return m.getStyle(headerBoxStyle), true
			}
			return 0, false
		},
		ConditionalStyles: style.Alternate(m.getStyle(alternate0Style), m.getStyle(alternate1Style), m.getStyle(alternate2Style)).
			WithStart(1, 0).
			WithEnd(len(m.settlements), m.settlementsTable.ColumnCount-1).
			Build(),
		ClearBeforeWrite: true,
	})
}

func setTablesExceptMembers(m *Manager) {
	m.expensesLeftTable = newExpensesLeftTable(m.file)
	m.expensesRightTable = newExpensesRightTable(m.file, m.MembersCount())
	m.expensesFullTable = newExpensesFullTable(m.file, m.MembersCount())
	m.transactionsTable = newTransactionsTable(m.file)
	m.debtMatrixTable = newDebtMatrixTable(m.file, m.MembersCount())
	m.settlementsTable = newSettlementsTable(m.file)
	m.baseStateTable = newBaseStateTable(m.file, m.MembersCount())
	m.metadataTable = newMetadataTable(m.file)
}

func createSheets(m *Manager) {
	fatalIfNotNil(m.file.SetSheetName(initialSheetName, membersSheet))
	defer initializeMembers(m)

	_, err := m.file.NewSheet(expensesSheet)
	fatalIfNotNil(err)
	defer initializeExpenses(m)

	_, err = m.file.NewSheet(transactionsSheet)
	fatalIfNotNil(err)
	defer initializeTransactions(m)

	_, err = m.file.NewSheet(debtMatrixSheet)
	fatalIfNotNil(err)
	defer initializeDebtMatrix(m)

	_, err = m.file.NewSheet(settlementsSheet)
	fatalIfNotNil(err)
	defer initializeSettlements(m)

	_, err = m.file.NewSheet(baseStateSheet)
	fatalIfNotNil(err)
	defer initializeBaseState(m)

	_, err = m.file.NewSheet(metadataSheet)
	fatalIfNotNil(err)
	defer initializeMetadata(m)
}

func initializeMembers(m *Manager) {
	m.membersTable.WriteRows(table.WriteRowsParams{
		RowCount: m.MembersCount(),
		HeaderWriter: func(cells []*table.WCell, _ *int) {
			cells[0].Value = "Name"
			cells[1].Value = "Card Number"
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = m.members.RequireMemberByIndex(rowNumber).Name
			cells[1].Value = m.members.RequireMemberByIndex(rowNumber).CardNumber
		},
		ColumnWidth: 32,
		RowStyler: func(row int) (int, bool) {
			if row == -1 {
				return m.getStyle(headerBoxStyle), true
			}
			return 0, false
		},
		ConditionalStyles: style.Alternate(m.getStyle(alternate0Style), m.getStyle(alternate1Style)).
			WithStart(0, 0).
			WithEnd(m.MembersCount()-1, 1).
			Build(),
	})
}

func initializeBaseState(m *Manager) {
	m.baseState = emptyMatrix(m.MembersCount())
	m.writeBaseState()
}

func initializeDebtMatrix(m *Manager) {
	m.debtMatrix = emptyMatrix(m.MembersCount())
	m.writeDebtMatrix()
}

func initializeTransactions(m *Manager) {

	m.transactionsTable.WriteRows(table.WriteRowsParams{
		RowCount: 1,
		HeaderWriter: func(cells []*table.WCell, _ *int) {
			cells[0].Value = "Time"
			cells[1].Value = "Receiver"
			cells[2].Value = "Payer"
			cells[3].Value = "Amount"
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = model.TimeOfGregorian(time.Date(
				2012,
				time.June,
				26,
				5,
				6,
				0,
				0,
				time.Local)).String()
			cells[1].Value = m.members.RequireMemberByIndex(0).Name
			cells[2].Value = m.members.RequireMemberByIndex(1).Name
			cells[3].Value = 0
			cells[3].Style = newInt(m.getStyle(moneyStyle))
		},
		ColumnWidth: 18,
		RowStyler: func(row int) (int, bool) {
			if row == -1 {
				return m.getStyle(headerBoxStyle), true
			}
			return 0, false
		},
		ConditionalStyles: style.Alternate(m.getStyle(alternate0Style), m.getStyle(alternate1Style), m.getStyle(alternate2Style)).
			WithStart(0, 0).
			WithEnd(0, m.transactionsTable.ColumnCount-1).
			WithModOffset(2).
			Build(),
	})
}

func initializeExpenses(m *Manager) {

	var totalAmountCell string
	m.expensesLeftTable.WriteRows(table.WriteRowsParams{
		RowCount: 1,
		HeaderWriter: func(cells []*table.WCell, mergeCount *int) {
			cells[0].Value = "Time"
			cells[1].Value = "Title"
			cells[2].Value = "Payer"
			cells[3].Value = "Total Amount"
		},
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = model.TimeOfGregorian(time.Date(
				2007,
				time.May,
				13,
				23,
				57,
				0,
				0,
				time.Local)).String()
			cells[1].Value = "example"
			cells[2].Value = m.members.RequireMemberByIndex(0).Name
			cells[3].Value = 0
			cells[3].Style = newInt(m.getStyle(moneyStyle))
			totalAmountCell = m.expensesLeftTable.GetCell(rowNumber, 3)
		},
		ColumnWidth: 16,
		RowStyler: func(row int) (int, bool) {
			if row == -1 {
				return m.getStyle(secondHeaderBoxStyle), true
			}
			return 0, false
		},
		ConditionalStyles: style.Alternate(m.getStyle(alternate0Style), m.getStyle(alternate1Style), m.getStyle(alternate2Style)).
			WithStart(0, 0).
			WithEnd(0, m.expensesLeftTable.ColumnCount+m.expensesRightTable.ColumnCount-1).
			Build(),
	})

	var weightCells []string
	m.expensesRightTable.WriteRows(table.WriteRowsParams{
		RowCount: 2,
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
					wc := fmt.Sprintf("IF(%s=TRUE, 1, %s)", m.expensesRightTable.GetCell(1, i), m.expensesRightTable.GetCell(1, i))
					weightCells = append(weightCells, wc)
					cells[i+1].Value = "Share Amount"
				} else if rowNumber == 1 {
					cells[i].Value = i >> 2
					totalWeightsFormula := fmt.Sprintf("SUM(%s)", strings.Join(weightCells, ", "))
					cells[i+1].Formula = fmt.Sprintf("(%s/%s)*%s", m.expensesRightTable.GetCell(rowNumber, i), totalWeightsFormula, totalAmountCell)
					cells[i+1].Style = newInt(m.getStyle(moneyStyle))
				}
			}
		},
		RowStyler: func(row int) (int, bool) {
			if row == -1 {
				return m.getStyle(headerBoxStyle), true
			}
			if row == 0 {
				return m.getStyle(secondHeaderBoxStyle), true
			}

			return 0, false
		},
		ColumnWidth: 11,
	})
}

func initializeSettlements(m *Manager) {
	m.settlements = make([]*model.Transaction, 0)
	m.writeSettlements()
}

func initializeMetadata(m *Manager) {
	m.metadataTable.WriteRows(table.WriteRowsParams{
		RowCount: 1,
		RowWriter: func(rowNumber int, cells []*table.WCell) {
			cells[0].Value = m.theme.Code()
		},
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
			fatalIfNotNil(log.CellErrorOf(err, t.SheetName, t.GetCell(rowNumber, 0)))
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
					requireMemberValidity(members, cells[i].Value, (i-4)/2, t.SheetName, t.GetCell(rowNumber, i))
				}
				return
			}

			if rowNumber == 0 {
				return
			}

			theTime, timeErr := model.ParseTime(cells[0].Value)
			fatalIfNotNil(log.CellErrorOf(timeErr, t.SheetName, t.GetCell(rowNumber, 0)))

			title := cells[1].Value

			payer := cells[2].Value
			requireMemberPresence(members, payer, t.SheetName, t.GetCell(rowNumber, 2))

			amount, amountErr := model.ParseAmount(cells[3].Value)
			fatalIfNotNil(log.CellErrorOf(amountErr, t.SheetName, t.GetCell(rowNumber, 3)))

			ex := &model.Expense{
				Title:     title,
				Time:      theTime,
				PayerName: payer,
				Amount:    amount,
			}

			var shares []model.Share
			for i := 4; i < t.ColumnCount; i += 2 {
				weight, err := model.ParseShareWeight(cells[i].Value)
				fatalIfNotNil(log.CellErrorOf(err, t.SheetName, t.GetCell(rowNumber, i)))

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
			theTime, err := model.ParseTime(cells[0].Value)
			fatalIfNotNil(log.CellErrorOf(err, t.SheetName, t.GetCell(rowNumber, 0)))

			receiver := cells[1].Value
			requireMemberPresence(members, receiver, t.SheetName, t.GetCell(rowNumber, 1))

			payer := cells[2].Value
			requireMemberPresence(members, payer, t.SheetName, t.GetCell(rowNumber, 2))

			amount, err := model.ParseAmount(cells[3].Value)
			fatalIfNotNil(log.CellErrorOf(err, t.SheetName, t.GetCell(rowNumber, 3)))

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
		RowCount: members.Count(),
		RowReader: func(rowNumber int, cells []*table.RCell) {
			if rowNumber == -1 {
				for i := 0; i < members.Count(); i++ {
					requireMemberValidity(members, cells[i+1].Value, i, t.SheetName, t.GetCell(rowNumber, i+1))
				}
				return
			}
			requireMemberValidity(members, cells[0].Value, rowNumber, t.SheetName, t.GetCell(rowNumber, 0))

			for i := 0; i < members.Count(); i++ {
				amount, err := model.ParseAmount(cells[i+1].Value)
				fatalIfNotNil(log.CellErrorOf(err, t.SheetName, t.GetCell(rowNumber, i+1)))
				baseState[rowNumber][i] = amount
			}
		},
		IncludeHeader:   true,
		UnknownRowCount: false,
	})

	return baseState
}

func loadMetadata(t *table.Table) *style.Theme {
	var theme *style.Theme
	t.ReadRows(table.ReadRowsParams{
		RowCount: 1,
		RowReader: func(rowNumber int, cells []*table.RCell) {
			if rowNumber == 0 {
				theme = style.ThemeFromCode(cells[0].Value)
			}
		},
		IncludeHeader: false,
	})

	return theme
}

func requireMemberValidity(members *store.MemberStore, memberName string, index int, sheetName, cell string) {
	if !members.IsValid(memberName, index) {
		log.FatalErrorByCaller(log.CellErrorOf(fmt.Errorf("found no member with name %q and index %d", memberName, index), sheetName, cell))
	}
}

func requireMemberPresence(members *store.MemberStore, memberName string, sheetName, cell string) {
	if !members.IsPresent(memberName) {
		log.FatalErrorByCaller(log.CellErrorOf(fmt.Errorf("found no member with name %q", memberName), sheetName, cell))
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

func copyMatrix(source [][]model.Amount) [][]model.Amount {
	result := emptyMatrix(len(source))
	for i := range source {
		copy(result[i], source[i])
	}
	return result
}

func newInt(n int) *int {
	p := new(int)
	*p = n
	return p
}
