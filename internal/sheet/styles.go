package sheet

import (
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/style"
)

const (
	blockStyle = iota
	moneyStyle
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
)

const (
	blockBGColor         = "#606060"
	headerBGColor        = "#CACACA"
	secondHeaderBGColor  = "#E6E6E6"
	defaultBorderColor   = "#D0CECE"
	borderColor          = "#404040"
	defaultFontColor     = "#000000"
	helpBGColor          = "#BFF1FF"
	helpFontColor        = "#114654"
	alternate0BGColor    = "#FFFFFF"
	alternate1BGColor    = "#E6F9FF"
	alternate2BGColor    = "#C2E8F2"
	alternateBorderColor = "#9E9E9E"
)

func createStyles(m *Manager) {
	builder := style.Empty()

	m.SetStyle(blockStyle, builder.WithBackground(blockBGColor).Build())
	m.SetStyle(moneyStyle, builder.WithMoneyFormat().Build())
	m.SetStyle(boxStyle, builder.WithFullBoarders(borderColor).Build())
	m.SetStyle(secondHeaderBoxStyle, builder.WithBackground(secondHeaderBGColor).
		WithCenterAlignment().WithFullBoarders(borderColor).WithFont(9, false, defaultFontColor).Build())
	m.SetStyle(headerStyle, builder.WithBackground(headerBGColor).Build())
	m.SetStyle(headerBoxStyle, builder.WithBackground(headerBGColor).
		WithFullBoarders(borderColor).WithCenterAlignment().Build())
	m.SetStyle(blockStyle, builder.WithBackground(blockBGColor).Build())
	m.SetStyle(lastUpdateStyle, builder.WithFont(8, false, defaultFontColor).WithCenterAlignment().Build())
	m.SetStyle(helpStyle, builder.WithBackground(helpBGColor).
		WithFont(9, true, helpFontColor).WithCenterAlignment().WithFullBoarders(borderColor).Build())
	m.SetStyle(rightBorderStyle, builder.WithRightBorderOnly(borderColor).Build())
	m.SetStyle(leftBorderStyle, builder.WithLeftBorderOnly(borderColor).Build())

	m.SetCondStyle(alternate0Style, builder.WithBackground(alternate0BGColor).WithFullBoarders(alternateBorderColor).Build())
	m.SetCondStyle(alternate1Style, builder.WithBackground(alternate1BGColor).WithFullBoarders(alternateBorderColor).Build())
	m.SetCondStyle(alternate2Style, builder.WithBackground(alternate2BGColor).WithFullBoarders(alternateBorderColor).Build())
}
