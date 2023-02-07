package style

import "github.com/xuri/excelize/v2"

type Builder struct {
	style *excelize.Style
}

func Empty() Builder {
	return Builder{
		new(excelize.Style),
	}
}

func (b *Builder) Build() *excelize.Style {
	r := b.style
	b.style = new(excelize.Style)
	return r
}

func (b *Builder) WithBackground(background string) *Builder {
	b.style.Fill.Pattern = 1
	b.style.Fill.Type = "pattern"
	b.style.Fill.Color = []string{background}
	return b
}

func (b *Builder) WithFullBoarders(borderColor string) *Builder {
	b.style.Border = []excelize.Border{
		{Type: "left", Color: borderColor},
		{Type: "top", Color: borderColor},
		{Type: "bottom", Color: borderColor},
		{Type: "right", Color: borderColor},
	}
	return b
}

func (b *Builder) WithFont(size float64, bold bool) *Builder {
	b.style.Font = &excelize.Font{
		Bold: bold,
		Size: size,
	}
	return b
}

func (b *Builder) WithCenterAlignment() *Builder {
	b.style.Alignment = &excelize.Alignment{
		Horizontal: "center",
	}
	return b
}

func (b *Builder) WithMoneyFormat() *Builder {
	b.style.NumFmt = 3
	return b
}
