package sheet

import (
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/style"
)

const (
	moneyStyle = iota
	secondHeaderBoxStyle
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
	defaultFontColor   = "#000000"
	defaultBorderColor = "#D0CECE"
	defaultCellBGColor = "#FFFFFF"
)

func createStyles(m *Manager) {
	builder := style.Empty()

	m.setStyle(moneyStyle, builder.WithMoneyFormat().Build())
	m.setStyle(secondHeaderBoxStyle, builder.WithBackground(m.theme.SecondHeaderBGColor).
		WithCenterAlignment().WithFullBoarders(m.theme.BorderColor).
		WithFont(9, false, defaultFontColor).Build())
	m.setStyle(headerBoxStyle, builder.WithBackground(m.theme.HeaderBGColor).
		WithFullBoarders(m.theme.BorderColor).WithCenterAlignment().Build())
	m.setStyle(lastUpdateStyle, builder.WithFont(8, false, defaultFontColor).WithCenterAlignment().Build())
	m.setStyle(helpStyle, builder.WithBackground(m.theme.HelpBGColor).
		WithFont(9, true, m.theme.HelpFontColor).WithCenterAlignment().
		WithFullBoarders(m.theme.BorderColor).Build())
	m.setStyle(rightBorderStyle, builder.WithRightBorderOnly(m.theme.BorderColor).Build())
	m.setStyle(leftBorderStyle, builder.WithLeftBorderOnly(m.theme.BorderColor).Build())
	m.setStyle(alternateBlockStyle, builder.WithCrossedBackground(m.theme.AlternateBlockBGColor).
		WithFullBoarders(m.theme.AlternateBorderColor).Build())

	m.setCondStyle(alternate0Style, builder.WithBackground(m.theme.Alternate0BGColor).
		WithFullBoarders(m.theme.AlternateBorderColor).Build())
	m.setCondStyle(alternate1Style, builder.WithBackground(m.theme.Alternate1BGColor).
		WithFullBoarders(m.theme.AlternateBorderColor).Build())
	m.setCondStyle(alternate2Style, builder.WithBackground(m.theme.Alternate2BGColor).
		WithFullBoarders(m.theme.AlternateBorderColor).Build())
}
