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

	expensesRowOffset = 3
	expensesColOffset = 5

	expensesLeftSideRowOffset  = 3
	expensesLeftSideColOffset  = 1
	expensesRightSideRowOffset = 2
	expensesRightSideColOffset = 5

	transactionsRowOffset = 2
	transactionsColOffset = 1

	debtorsMatrixRowOffset = 2
	debtorsMatrixColOffset = 2

	baseStateRowOffset = 2
	baseStateColOffset = 2
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
	file         *excelize.File
	members      []*model.Member
	expenses     []*model.Expense
	transactions []*model.Transaction
	debtMatrix   [][]model.Amount
	baseState    [][]model.Amount
	sheetIndices map[string]int
	styleIndices map[string]int
}

func NewManager(members []*model.Member) *Manager {
	file := excelize.NewFile()

	m := &Manager{
		members:      members,
		file:         file,
		sheetIndices: make(map[string]int),
		styleIndices: make(map[string]int),
	}

	createStyles(m)
	createSheets(m)

	return m
}

func LoadManager(fileName string) (*Manager, error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	members := loadMembers(file)
	m := &Manager{
		file:         file,
		sheetIndices: make(map[string]int),
		styleIndices: make(map[string]int),
		members:      members,
		expenses:     loadExpenses(file, members),
		transactions: loadTransactions(file, members),
		baseState:    loadBaseState(file, members),
	}

	createStyles(m)

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

	return m, nil
}

func loadMembers(file *excelize.File) []*model.Member {
	setOffsets(membersRowOffset, membersColOffset)
	defer resetOffsets()

	members := make([]*model.Member, 0)
	for i := 0; ; i++ {
		name, _ := file.GetCellValue(membersSheet, cell(i, 0))
		if name == "" {
			break
		}

		cardNumber, _ := file.GetCellValue(membersSheet, cell(i, 1))
		members = append(members, &model.Member{
			ID:         model.MID(i),
			Name:       name,
			CardNumber: cardNumber,
		})
	}

	return members
}

func loadExpenses(file *excelize.File, members []*model.Member) []*model.Expense {
	setOffsets(expensesRowOffset, expensesColOffset)
	defer resetOffsets()

	var expenses []*model.Expense
	var err error
	for r := 0; ; r++ {

		theTime, _ := file.GetCellValue(expensesSheet, cell(r, -4), excelize.Options{RawCellValue: false})
		title, _ := file.GetCellValue(expensesSheet, cell(r, -3))
		payer, _ := file.GetCellValue(expensesSheet, cell(r, -2))
		total, _ := file.GetCellValue(expensesSheet, cell(r, -1))

		if total == "" || payer == "" {
			break
		}

		ex := new(model.Expense)
		ex.Title = title

		ex.Time, err = time.Parse(timeLayout, theTime)
		panicE(err)
		ex.Time = ex.Time.Local()

		if mIndex := findMemberIndex(members, payer); mIndex >= 0 {
			ex.PayerID = model.MID(mIndex)
		} else {
			panic(fmt.Errorf("found no member with name %q", payer))
		}

		ex.Amount, err = model.ParseAmount(total)
		panicE(err)

		var shares []model.Share
		for idx, member := range members {
			name, _ := file.GetCellValue(expensesSheet, cell(-2, idx*2))
			if name == "" {
				name, _ = file.CalcCellValue(expensesSheet, cell(-2, idx*2))
			}

			if name != member.Name {
				panic(fmt.Errorf("member names do not match in 'member' and 'expenses': %q != %q", name, member.Name))
			}

			weightStr, _ := file.GetCellValue(expensesSheet, cell(r, idx*2))
			weight, err := strconv.Atoi(weightStr)
			panicE(err)

			shares = append(shares, model.Share{
				MemberID:    member.ID,
				ShareWeight: weight,
			})
		}
		ex.Shares = shares

		expenses = append(expenses, ex)
	}

	return expenses
}

func loadTransactions(file *excelize.File, members []*model.Member) []*model.Transaction {
	var err error
	var transactions []*model.Transaction

	for r := 2; ; r++ {
		timeStr, _ := file.GetCellValue(transactionsSheet, cell(r, 1), excelize.Options{RawCellValue: false})
		receiver, _ := file.GetCellValue(transactionsSheet, cell(r, 2))
		payer, _ := file.GetCellValue(transactionsSheet, cell(r, 3))
		amountStr, _ := file.GetCellValue(transactionsSheet, cell(r, 4))

		if amountStr == "" {
			break
		}

		tr := new(model.Transaction)

		tr.Amount, err = model.ParseAmount(amountStr)
		panicE(err)

		tr.Time, err = time.Parse(timeLayout, timeStr)
		tr.Time = tr.Time.Local()
		panicE(err)

		if mIndex := findMemberIndex(members, payer); mIndex >= 0 {
			tr.PayerID = model.MID(mIndex)
		} else {
			panic(fmt.Errorf("found no member with name %q as payer", payer))
		}

		if mIndex := findMemberIndex(members, receiver); mIndex >= 0 {
			tr.ReceiverID = model.MID(mIndex)
		} else {
			panic(fmt.Errorf("found no member with name %q as receiver", payer))
		}

		transactions = append(transactions, tr)
	}

	return transactions
}

func loadBaseState(file *excelize.File, members []*model.Member) [][]model.Amount {
	setOffsets(baseStateRowOffset, baseStateColOffset)
	defer resetOffsets()

	baseState := make([][]model.Amount, len(members))
	for i := range baseState {
		baseState[i] = make([]model.Amount, len(members))
	}

	for r := 0; r < len(baseState); r++ {
		rowName, _ := file.CalcCellValue(baseStateSheet, cell(r, -1))

		if mIndex := findMemberIndex(members, rowName); mIndex < 0 || r != mIndex {
			panic(fmt.Errorf("found no member with name %q and index %d", rowName, r))
		}

		for c := r + 1; c < len(baseState[r]); c++ {
			colName, _ := file.CalcCellValue(baseStateSheet, cell(-1, c))

			if mIndex := findMemberIndex(members, colName); mIndex < 0 || c != mIndex {
				panic(fmt.Errorf("found no member with name %q and index %d", colName, c))
			}

			amountStr, _ := file.GetCellValue(baseStateSheet, cell(r, c))
			amount, err := model.ParseAmount(amountStr)
			panicE(err)
			baseState[r][c] = amount
		}
	}

	return baseState
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

func initializeMembers(m *Manager) {

	t := table.Table{
		File:         m.file,
		SheetName:    membersSheet,
		RowOffset:    membersRowOffset,
		ColumnOffset: membersColOffset,
		ColumnWidth:  32,
		RowCount:     len(m.members),
		ColumnCount:  2,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}

	t.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.Cell, _ *int) {
			cells[0].Value = "Name"
			cells[1].Value = "Card Number"
		},
		RowWriter: func(rowNumber int, cells []*table.Cell) {
			cells[0].Value = m.members[rowNumber].Name
			cells[1].Value = m.members[rowNumber].CardNumber
		},
	})
}

func initializeMetadata(m *Manager) {}

func initializeBaseState(m *Manager) {
	setOffsets(baseStateRowOffset, baseStateColOffset)
	defer resetOffsets()

	m.file.SetColWidth(baseStateSheet, column(1), column(m.MembersCount()+1), 16)

	for i := 0; i < m.MembersCount(); i++ {
		nameRef := memberNameRef(m.members[i].ID)
		m.file.SetCellFormula(baseStateSheet, cell(i, -1), nameRef)
		m.file.SetCellFormula(baseStateSheet, cell(-1, i), nameRef)
		m.file.SetCellStyle(baseStateSheet, cell(i, i), cell(m.MembersCount()-1, i), m.GetStyle(blockStyle))

		for j := i + 1; j < m.MembersCount(); j++ {
			m.file.SetCellValue(baseStateSheet, cell(i, j), 0)
		}
	}
}

func initializeDebtMatrix(m *Manager) {
	m.file.SetColWidth(debtMatrixSheet, column(1), column(1), 128)
	m.file.SetCellValue(debtMatrixSheet, cell(1, 1), "Run 'update' command to generate debt matrix")
}

func initializeTransactions(m *Manager) {

	t := table.Table{
		File:         m.file,
		SheetName:    transactionsSheet,
		RowOffset:    transactionsRowOffset,
		ColumnOffset: transactionsColOffset,
		RowCount:     0,
		ColumnCount:  4,
		ColumnWidth:  16,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}

	t.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.Cell, _ *int) {
			cells[0].Value = "Time"
			cells[1].Value = "Receiver"
			cells[2].Value = "Payer"
			cells[3].Value = "Amount"
		},
	})
}

func initializeExpenses(m *Manager) {

	t := table.Table{
		File:         m.file,
		SheetName:    expensesSheet,
		RowOffset:    expensesLeftSideRowOffset,
		ColumnOffset: expensesLeftSideColOffset,
		RowCount:     1,
		ColumnCount:  4,
		ColumnWidth:  16,
		ErrorHandler: func(err error) {
			panic(err)
		},
	}

	var totalAmountCell string
	t.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.Cell, mergeCount *int) {
			cells[0].Value = "Time"
			cells[1].Value = "Title"
			cells[2].Value = "Payer"
			cells[3].Value = "Total Amount"
		},
		RowWriter: func(rowNumber int, cells []*table.Cell) {
			cells[0].Value = time.Now()
			cells[1].Value = "food"
			cells[2].Value = m.members[0].Name
			cells[3].Value = 300
			totalAmountCell = t.GetCell(rowNumber, 3)
		},
	})

	t.RowOffset = expensesRightSideRowOffset
	t.ColumnOffset = expensesRightSideColOffset
	t.RowCount = 2
	t.ColumnCount = m.MembersCount() * 2

	var weightCells []string
	t.WriteRows(table.WriteRowsParams{
		HeaderWriter: func(cells []*table.Cell, mergeCount *int) {
			*mergeCount = 2
			for i, v := range m.members {
				cells[i*2].Value = v.Name
			}
		},
		RowWriter: func(rowNumber int, cells []*table.Cell) {
			for i := 0; i < m.MembersCount()*2; i += 2 {
				if rowNumber == 0 {
					cells[i].Value = "Share Weight"
					weightCells = append(weightCells, t.GetCell(1, i))
					cells[i+1].Value = "Share Amount"
				} else if rowNumber == 1 {
					cells[i].Value = i >> 2
					totalWeightsFormula := fmt.Sprintf("SUM(%s)", strings.Join(weightCells, ", "))
					cells[i+1].Formula = fmt.Sprintf("(%s/%s)*%s", t.GetCell(rowNumber, i), totalWeightsFormula, totalAmountCell)
				}
			}
		},
	})
}
