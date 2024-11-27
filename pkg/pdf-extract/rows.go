package pdfextract

type Rows []Row

type Text struct {
	Font     string
	FontSize float64
	X        float64
	Y        float64
	S        string
}

type Row struct {
	Position float64
	Content  []Text
}

func (rows Rows) Len() int {
	return len(rows)
}
func (a Rows) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Rows) Less(i, j int) bool {

	return a[i].Position > a[j].Position
}

type Position struct {
	X float64
	Y float64
}

type TextObject struct {
	FontName     string
	Encoding     string
	ResourceName string
	Text         string
	FontSize     float64
	Position     Position
}
