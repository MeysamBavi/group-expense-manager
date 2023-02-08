package style

import "github.com/xuri/excelize/v2"

type Builder struct {
	style *excelize.Style
}

func Empty() *Builder {
	return &Builder{
		new(excelize.Style),
	}
}

func (b *Builder) Build() *excelize.Style {
	r := b.style
	b.style = new(excelize.Style)
	return r
}

func (b *Builder) WithBackground(color string) *Builder {
	b.style.Fill.Pattern = 1
	b.style.Fill.Type = "pattern"
	b.style.Fill.Color = []string{color}
	return b
}

func (b *Builder) WithCrossedBackground(color string) *Builder {
	b.style.Fill.Pattern = 8
	b.style.Fill.Type = "pattern"
	b.style.Fill.Color = []string{color}
	return b
}

func (b *Builder) WithRightBorderOnly(color string) *Builder {
	b.style.Border = []excelize.Border{
		{Type: "right", Color: color, Style: 1},
	}
	return b
}

func (b *Builder) WithLeftBorderOnly(color string) *Builder {
	b.style.Border = []excelize.Border{
		{Type: "left", Color: color, Style: 1},
	}
	return b
}

func (b *Builder) WithFullBoarders(borderColor string) *Builder {
	b.style.Border = []excelize.Border{
		{Type: "left", Color: borderColor, Style: 1},
		{Type: "top", Color: borderColor, Style: 1},
		{Type: "bottom", Color: borderColor, Style: 1},
		{Type: "right", Color: borderColor, Style: 1},
	}
	return b
}

func (b *Builder) WithFont(size float64, bold bool, color string) *Builder {
	b.style.Font = &excelize.Font{
		Bold:  bold,
		Size:  size,
		Color: color,
	}
	return b
}

func (b *Builder) WithCenterAlignment() *Builder {
	b.style.Alignment = &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}
	return b
}

func (b *Builder) WithMoneyFormat() *Builder {
	b.style.NumFmt = 3
	return b
}
