package sheet

import (
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/style"
)

const (
	moneyStyle = iota
	boxStyle
	secondHeaderBoxStyle
	headerStyle
	headerBoxStyle
	lastUpdateStyle
	helpStyle
	rightBorderStyle
	leftBorderStyle
	alternate0Style
	alternate1Style
	alternate2Style
	alternateBlockStyle
)

const (
	headerBGColor         = "#CACACA"
	secondHeaderBGColor   = "#E6E6E6"
	defaultBorderColor    = "#D0CECE"
	borderColor           = "#404040"
	defaultFontColor      = "#000000"
	helpBGColor           = "#BFF1FF"
	helpFontColor         = "#114654"
	alternate0BGColor     = "#FFFFFF"
	alternate1BGColor     = "#E6F9FF"
	alternate2BGColor     = "#C2E8F2"
	alternateBorderColor  = "#9E9E9E"
	alternateBlockBGColor = "#2C8196"
)

func createStyles(m *Manager) {
	builder := style.Empty()

	m.setStyle(moneyStyle, builder.WithMoneyFormat().Build())
	m.setStyle(boxStyle, builder.WithFullBoarders(borderColor).Build())
	m.setStyle(secondHeaderBoxStyle, builder.WithBackground(secondHeaderBGColor).
		WithCenterAlignment().WithFullBoarders(borderColor).WithFont(9, false, defaultFontColor).Build())
	m.setStyle(headerStyle, builder.WithBackground(headerBGColor).Build())
	m.setStyle(headerBoxStyle, builder.WithBackground(headerBGColor).
		WithFullBoarders(borderColor).WithCenterAlignment().Build())
	m.setStyle(lastUpdateStyle, builder.WithFont(8, false, defaultFontColor).WithCenterAlignment().Build())
	m.setStyle(helpStyle, builder.WithBackground(helpBGColor).
		WithFont(9, true, helpFontColor).WithCenterAlignment().WithFullBoarders(borderColor).Build())
	m.setStyle(rightBorderStyle, builder.WithRightBorderOnly(borderColor).Build())
	m.setStyle(leftBorderStyle, builder.WithLeftBorderOnly(borderColor).Build())
	m.setStyle(alternateBlockStyle, builder.WithCrossedBackground(alternateBlockBGColor).WithFullBoarders(alternateBorderColor).Build())

	m.setCondStyle(alternate0Style, builder.WithBackground(alternate0BGColor).WithFullBoarders(alternateBorderColor).Build())
	m.setCondStyle(alternate1Style, builder.WithBackground(alternate1BGColor).WithFullBoarders(alternateBorderColor).Build())
	m.setCondStyle(alternate2Style, builder.WithBackground(alternate2BGColor).WithFullBoarders(alternateBorderColor).Build())
}
