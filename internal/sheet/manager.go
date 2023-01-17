package sheet

import (
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	"github.com/xuri/excelize/v2"
	"strings"
	"time"
)

const (
	initialSheetName  = "Sheet1"
	membersSheet      = "members"
	expensesSheet     = "expenses"
	transactionsSheet = "transactions"
	debtMatrixSheet   = "Debt Matrix"
	baseStateSheet    = "base state"
	metadataSheet     = "metadata (unmodifiable)"
)

type Manager struct {
	members           []*model.Member
	file              *excelize.File
	membersIndex      int
	expensesIndex     int
	transactionsIndex int
	debtMatrixIndex   int
	baseStateIndex    int
	metadataIndex     int
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

func NewManager(members []*model.Member) *Manager {
	file := excelize.NewFile()

	m := &Manager{
		members: members,
		file:    file,
	}

	createSheets(m)

	return m
}

func createSheets(m *Manager) {
	panicE(m.file.SetSheetName(initialSheetName, membersSheet))
	defer initializeMembers(m)

	i, err := m.file.NewSheet(expensesSheet)
	panicE(err)
	m.expensesIndex = i
	defer initializeExpenses(m)

	i, err = m.file.NewSheet(transactionsSheet)
	panicE(err)
	m.transactionsIndex = i
	defer initializeTransactions(m)

	i, err = m.file.NewSheet(debtMatrixSheet)
	panicE(err)
	m.debtMatrixIndex = i
	initializeDebtMatrix(m)

	i, err = m.file.NewSheet(baseStateSheet)
	panicE(err)
	m.baseStateIndex = i
	initializeBaseState(m)

	i, err = m.file.NewSheet(metadataSheet)
	panicE(err)
	m.metadataIndex = i
	initializeMetadata(m)
}

func initializeMembers(m *Manager) {
	m.file.SetColWidth(membersSheet, column(1), column(2), 32)

	m.file.SetCellValue(membersSheet, cell(1, 1), "Name")
	m.file.SetCellValue(membersSheet, cell(1, 2), "Card Number")

	for i, member := range m.members {
		m.file.SetCellValue(membersSheet, cell(i+2, 1), member.Name)
		m.file.SetCellValue(membersSheet, cell(i+2, 2), member.CardNumber)
	}
}

func initializeMetadata(m *Manager) {

}

func initializeBaseState(m *Manager) {

}

func initializeDebtMatrix(m *Manager) {

}

func initializeTransactions(m *Manager) {

}

func initializeExpenses(m *Manager) {
	m.file.SetColWidth(expensesSheet, column(1), column(m.MembersCount()*2+4), 16)

	m.file.SetCellValue(expensesSheet, cell(1, 1), "Time")
	m.file.SetCellValue(expensesSheet, cell(1, 2), "Title")
	m.file.SetCellValue(expensesSheet, cell(1, 3), "Payer")
	m.file.SetCellValue(expensesSheet, cell(1, 4), "Total Amount")

	weightCells := make([]string, 0, m.MembersCount())
	for i, member := range m.members {
		m.file.MergeCell(expensesSheet, cell(1, i*2+5), cell(1, i*2+6))
		m.file.SetCellFormula(expensesSheet, cell(1, i*2+5), memberNameRef(member.ID))

		m.file.SetCellValue(expensesSheet, cell(2, i*2+5), "Share Weight")
		m.file.SetCellValue(expensesSheet, cell(2, i*2+6), "Share Amount")

		weightCell := cell(3, i*2+5)
		weightCells = append(weightCells, weightCell)
	}

	m.file.SetCellValue(expensesSheet, cell(3, 1), time.Now())
	m.file.SetCellValue(expensesSheet, cell(3, 2), "food")
	m.file.SetCellValue(expensesSheet, cell(3, 3), "Fred")
	m.file.SetCellValue(expensesSheet, cell(3, 4), 300)

	totalWeightsFormula := fmt.Sprintf("SUM(%s)", strings.Join(weightCells, ", "))
	for i, weightCell := range weightCells {
		shareAmountFormula := fmt.Sprintf("(%s/%s)*%s", weightCell, totalWeightsFormula, cell(3, 4))
		m.file.SetCellFormula(expensesSheet, cell(3, i*2+6), shareAmountFormula)
		m.file.SetCellValue(expensesSheet, weightCell, i>>1)
	}
}
