package types

func (rows Rows) Len() int {
	return len(rows)
}
func (a Rows) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Rows) Less(i, j int) bool {

	return a[i].Position > a[j].Position
}
