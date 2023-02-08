package sheet

import (
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/style"
)

const (
	blockStyle = iota
	moneyStyle
	boxStyle
	shadedStyle
	headerStyle
	headerBoxStyle
	lastUpdateStyle
	helpStyle
	rightBorderStyle
	leftBorderStyle
)

const (
	blockBGColor       = "#606060"
	headerBGColor      = "#CACACA"
	shadedCellBGColor  = "#E6E6E6"
	defaultBorderColor = "#202020"
	defaultFontColor   = "#000000"
	helpBGColor        = "#FFFFCC"
	helpFontColor      = "#9C6500"
)

func createStyles(m *Manager) {
	builder := style.Empty()

	m.SetStyle(blockStyle, builder.WithBackground(blockBGColor).Build())
	m.SetStyle(moneyStyle, builder.WithMoneyFormat().Build())
	m.SetStyle(boxStyle, builder.WithFullBoarders(defaultBorderColor).Build())
	m.SetStyle(shadedStyle, builder.WithBackground(shadedCellBGColor).Build())
	m.SetStyle(headerStyle, builder.WithBackground(headerBGColor).Build())
	m.SetStyle(headerBoxStyle, builder.WithBackground(headerBGColor).WithFullBoarders(defaultBorderColor).Build())
	m.SetStyle(blockStyle, builder.WithBackground(blockBGColor).Build())
	m.SetStyle(lastUpdateStyle, builder.WithFont(8, false, defaultFontColor).WithCenterAlignment().Build())
	m.SetStyle(helpStyle, builder.WithBackground(helpBGColor).WithFont(9, true, helpFontColor).Build())
	m.SetStyle(rightBorderStyle, builder.WithRightBorderOnly(defaultBorderColor).Build())
	m.SetStyle(leftBorderStyle, builder.WithLeftBorderOnly(defaultBorderColor).Build())
}
