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
	debtMatrixSheet   = "debt-matrix"
	baseStateSheet    = "base state"
	metadataSheet     = "metadata (unmodifiable)"
)

const (
	blockStyle = "block"
)

type Manager struct {
	members      []*model.Member
	file         *excelize.File
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
	// TODO
	//file := excelize.OpenFile(fileName)

	return nil, nil
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
	// TODO
}

func createStyles(m *Manager) {
	m.SetStyle(blockStyle, &excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#606060"}, Pattern: 1, Shading: 0},
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
	m.file.SetColWidth(membersSheet, column(1), column(2), 32)

	m.file.SetCellValue(membersSheet, cell(1, 1), "Name")
	m.file.SetCellValue(membersSheet, cell(1, 2), "Card Number")

	for i, member := range m.members {
		m.file.SetCellValue(membersSheet, cell(i+2, 1), member.Name)
		m.file.SetCellValue(membersSheet, cell(i+2, 2), member.CardNumber)
	}
}

func initializeMetadata(m *Manager) {}

func initializeBaseState(m *Manager) {
	m.file.SetColWidth(baseStateSheet, column(1), column(m.MembersCount()+1), 16)

	for i := 0; i < m.MembersCount(); i++ {
		nameRef := memberNameRef(m.members[i].ID)
		m.file.SetCellFormula(baseStateSheet, cell(i+2, 1), nameRef)
		m.file.SetCellFormula(baseStateSheet, cell(1, i+2), nameRef)
		m.file.SetCellStyle(baseStateSheet, cell(i+2, i+2), cell(m.MembersCount()-1+2, i+2), m.GetStyle(blockStyle))

		for j := i + 1; j < m.MembersCount(); j++ {
			m.file.SetCellValue(baseStateSheet, cell(i+2, j+2), 0)
		}
	}
}

func initializeDebtMatrix(m *Manager) {
	m.file.SetColWidth(debtMatrixSheet, column(1), column(1), 128)
	m.file.SetCellValue(debtMatrixSheet, cell(1, 1), "Run 'update' command to generate debt matrix")
}

func initializeTransactions(m *Manager) {
	m.file.SetColWidth(transactionsSheet, column(1), column(4), 16)

	m.file.SetCellValue(transactionsSheet, cell(1, 1), "Time")
	m.file.SetCellValue(transactionsSheet, cell(1, 2), "Receiver")
	m.file.SetCellValue(transactionsSheet, cell(1, 3), "Payer")
	m.file.SetCellValue(transactionsSheet, cell(1, 4), "Amount")
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
