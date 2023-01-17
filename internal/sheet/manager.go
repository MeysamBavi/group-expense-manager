package sheet

import (
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	"github.com/xuri/excelize/v2"
)

type Manager struct {
	members []*model.Member
	file    *excelize.File
}

func (m *Manager) SaveAs(name string) error {
	//return m.file.SaveAs(name)
	return nil
}

func NewManager(members []*model.Member) *Manager {
	file := excelize.NewFile()
	m := &Manager{
		members: members,
		file:    file,
	}

	return m
}
