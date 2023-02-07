package sheet

import (
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/style"
)

const (
	blockStyle = iota
	moneyStyle
)

const (
	BlockBG            = "#606060"
	HeaderBG           = "#BABABA"
	ShadedCellBG       = "#E6E6E6"
	DefaultBorderColor = "#434343"
)

func createStyles(m *Manager) {
	builder := style.Empty()

	m.SetStyle(blockStyle, builder.WithBackground(BlockBG).Build())
	m.SetStyle(moneyStyle, builder.WithMoneyFormat().Build())
}
