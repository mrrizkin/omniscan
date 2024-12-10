package types

type (
	Text struct {
		Font     string
		FontSize float64
		X        float64
		Y        float64
		S        string
	}

	TextHorizontal []Text

	Row struct {
		Position int64
		Content  TextHorizontal
	}

	Rows []*Row
)
