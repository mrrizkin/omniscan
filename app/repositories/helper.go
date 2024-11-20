package repositories

type (
	wb struct {
		whereCount int
		where      string
		whereArgs  []interface{}
	}

	jb struct {
		join     string
		joinArgs []interface{}
	}

	jcb struct {
		conditionCount int
		condition      string
		conditionArgs  []interface{}
	}
)

func whereBuilder() *wb {
	return &wb{
		whereCount: 0,
		where:      "",
		whereArgs:  make([]interface{}, 0),
	}
}

func (wb *wb) And(where string, args ...interface{}) {
	if wb.whereCount != 0 {
		wb.where += " AND"
	}

	wb.where += " " + where
	wb.whereArgs = append(wb.whereArgs, args...)
	wb.whereCount++
}

func (wb *wb) Or(where string, args ...interface{}) {
	if wb.whereCount != 0 {
		wb.where += " OR"
	}

	wb.where += " " + where
	wb.whereArgs = append(wb.whereArgs, args...)
	wb.whereCount++
}

func (wb *wb) Get() (string, []interface{}) {
	return wb.where, wb.whereArgs
}

func joinCondBuilder() *jcb {
	return &jcb{
		conditionCount: 0,
		condition:      "",
		conditionArgs:  make([]interface{}, 0),
	}
}

func (jcb *jcb) And(condition string, args ...interface{}) {
	if jcb.conditionCount != 0 {
		jcb.condition += " AND"
	}

	jcb.condition += " " + condition
	jcb.conditionArgs = append(jcb.conditionArgs, args...)
	jcb.conditionCount++
}

func (jcb *jcb) Or(condition string, args ...interface{}) {
	if jcb.conditionCount != 0 {
		jcb.condition += " OR"
	}

	jcb.condition += " " + condition
	jcb.conditionArgs = append(jcb.conditionArgs, args...)
	jcb.conditionCount++
}

func (jcb *jcb) Get() (string, []interface{}) {
	return jcb.condition, jcb.conditionArgs
}

func joinBuilder() *jb {
	return &jb{
		join: "",
	}
}

func (jb *jb) InnerJoin(table string, condition string, args ...interface{}) {
	jb.join += " INNER JOIN"
	jb.join += " " + table + " ON " + condition
	jb.joinArgs = append(jb.joinArgs, args...)
}

func (jb *jb) LeftJoin(table string, condition string, args ...interface{}) {
	jb.join += " LEFT JOIN"
	jb.join += " " + table + " ON " + condition
	jb.joinArgs = append(jb.joinArgs, args...)
}

func (jb *jb) RightJoin(table string, condition string, args ...interface{}) {
	jb.join += " RIGHT JOIN"
	jb.join += " " + table + " ON " + condition
	jb.joinArgs = append(jb.joinArgs, args...)
}

func (jb *jb) Get() (string, []interface{}) {
	return jb.join, jb.joinArgs
}

func chunk[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for chunkSize < len(slice) {
		slice, chunks = slice[chunkSize:], append(chunks, slice[0:chunkSize:chunkSize])
	}
	chunks = append(chunks, slice)
	return chunks
}
