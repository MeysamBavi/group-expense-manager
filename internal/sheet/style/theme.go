package style

type Theme struct {
	HeaderBGColor         string
	BorderColor           string
	SecondHeaderBGColor   string
	AlternateBorderColor  string
	HelpBGColor           string
	HelpFontColor         string
	Alternate0BGColor     string
	Alternate1BGColor     string
	Alternate2BGColor     string
	AlternateBlockBGColor string
}

func ThemeFromCode(code string) *Theme {
	return &Theme{
		HeaderBGColor:         code[0:7],
		BorderColor:           code[7:14],
		SecondHeaderBGColor:   code[14:21],
		AlternateBorderColor:  code[21:28],
		HelpBGColor:           code[28:35],
		HelpFontColor:         code[35:42],
		Alternate0BGColor:     code[42:49],
		Alternate1BGColor:     code[49:56],
		Alternate2BGColor:     code[56:63],
		AlternateBlockBGColor: code[63:70],
	}
}

func (t *Theme) Code() string {
	return t.HeaderBGColor + t.BorderColor + t.SecondHeaderBGColor + t.AlternateBorderColor +
		t.HelpBGColor + t.HelpFontColor + t.Alternate0BGColor + t.Alternate1BGColor +
		t.Alternate2BGColor + t.AlternateBlockBGColor
}

func BlueTheme() *Theme {
	return &Theme{
		HeaderBGColor:         "#CACACA",
		BorderColor:           "#404040",
		SecondHeaderBGColor:   "#E6E6E6",
		AlternateBorderColor:  "#9E9E9E",
		HelpBGColor:           "#BFF1FF",
		HelpFontColor:         "#114654",
		Alternate0BGColor:     "#FFFFFF",
		Alternate1BGColor:     "#E6F9FF",
		Alternate2BGColor:     "#C2E8F2",
		AlternateBlockBGColor: "#2C7F96",
	}
}

func RedTheme() *Theme {
	return &Theme{
		HeaderBGColor:         "#CACACA",
		BorderColor:           "#404040",
		SecondHeaderBGColor:   "#E6E6E6",
		AlternateBorderColor:  "#9E9E9E",
		HelpBGColor:           "#FFBFBF",
		HelpFontColor:         "#541111",
		Alternate0BGColor:     "#FFFFFF",
		Alternate1BGColor:     "#FFE6E6",
		Alternate2BGColor:     "#F2C2C2",
		AlternateBlockBGColor: "#962C2C",
	}
}

func GreenTheme() *Theme {
	return &Theme{
		HeaderBGColor:         "#CACACA",
		BorderColor:           "#404040",
		SecondHeaderBGColor:   "#E6E6E6",
		AlternateBorderColor:  "#9E9E9E",
		HelpBGColor:           "#BFFFD8",
		HelpFontColor:         "#11542B",
		Alternate0BGColor:     "#FFFFFF",
		Alternate1BGColor:     "#E6FFEF",
		Alternate2BGColor:     "#C2F2D4",
		AlternateBlockBGColor: "#2C9654",
	}
}

func YellowTheme() *Theme {
	return &Theme{
		HeaderBGColor:         "#CACACA",
		BorderColor:           "#404040",
		SecondHeaderBGColor:   "#E6E6E6",
		AlternateBorderColor:  "#9E9E9E",
		HelpBGColor:           "#FFFEBF",
		HelpFontColor:         "#545311",
		Alternate0BGColor:     "#FFFFFF",
		Alternate1BGColor:     "#FFFFE6",
		Alternate2BGColor:     "#F2F1C2",
		AlternateBlockBGColor: "#96962C",
	}
}

func PurpleTheme() *Theme {
	return &Theme{
		HeaderBGColor:         "#CACACA",
		BorderColor:           "#404040",
		SecondHeaderBGColor:   "#E6E6E6",
		AlternateBorderColor:  "#9E9E9E",
		HelpBGColor:           "#E6BFFF",
		HelpFontColor:         "#3A1154",
		Alternate0BGColor:     "#FFFFFF",
		Alternate1BGColor:     "#F5E6FF",
		Alternate2BGColor:     "#E0C2F2",
		AlternateBlockBGColor: "#6E2C96",
	}
}
