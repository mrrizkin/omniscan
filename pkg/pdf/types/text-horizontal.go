package types

func (rows TextHorizontal) Len() int {
	return len(rows)
}
func (a TextHorizontal) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a TextHorizontal) Less(i, j int) bool {
	return a[i].X < a[j].X
}
