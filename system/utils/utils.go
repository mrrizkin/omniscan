package utils

import "reflect"

type (
	whereBuilder struct {
		whereCount int
		where      string
		whereArgs  []interface{}
	}

	joinBuilder struct {
		join     string
		joinArgs []interface{}
	}

	joinCondBuilder struct {
		conditionCount int
		condition      string
		conditionArgs  []interface{}
	}
)

func NewWhereBuilder() *whereBuilder {
	return &whereBuilder{
		whereCount: 0,
		where:      "",
		whereArgs:  make([]interface{}, 0),
	}
}

func (wb *whereBuilder) And(where string, args ...interface{}) {
	if wb.whereCount != 0 {
		wb.where += " AND"
	}

	wb.where += " " + where
	wb.whereArgs = append(wb.whereArgs, args...)
	wb.whereCount++
}

func (wb *whereBuilder) Or(where string, args ...interface{}) {
	if wb.whereCount != 0 {
		wb.where += " OR"
	}

	wb.where += " " + where
	wb.whereArgs = append(wb.whereArgs, args...)
	wb.whereCount++
}

func (wb *whereBuilder) Get() (string, []interface{}) {
	return wb.where, wb.whereArgs
}

func NewJoinConditionBuilder() *joinCondBuilder {
	return &joinCondBuilder{
		conditionCount: 0,
		condition:      "",
		conditionArgs:  make([]interface{}, 0),
	}
}

func (jcb *joinCondBuilder) And(condition string, args ...interface{}) {
	if jcb.conditionCount != 0 {
		jcb.condition += " AND"
	}

	jcb.condition += " " + condition
	jcb.conditionArgs = append(jcb.conditionArgs, args...)
	jcb.conditionCount++
}

func (jcb *joinCondBuilder) Or(condition string, args ...interface{}) {
	if jcb.conditionCount != 0 {
		jcb.condition += " OR"
	}

	jcb.condition += " " + condition
	jcb.conditionArgs = append(jcb.conditionArgs, args...)
	jcb.conditionCount++
}

func (jcb *joinCondBuilder) Get() (string, []interface{}) {
	return jcb.condition, jcb.conditionArgs
}

func NewJoinBuilder() *joinBuilder {
	return &joinBuilder{
		join: "",
	}
}

func (jb *joinBuilder) InnerJoin(table string, condition string, args ...interface{}) {
	jb.join += " INNER JOIN"
	jb.join += " " + table + " ON " + condition
	jb.joinArgs = append(jb.joinArgs, args...)
}

func (jb *joinBuilder) LeftJoin(table string, condition string, args ...interface{}) {
	jb.join += " LEFT JOIN"
	jb.join += " " + table + " ON " + condition
	jb.joinArgs = append(jb.joinArgs, args...)
}

func (jb *joinBuilder) RightJoin(table string, condition string, args ...interface{}) {
	jb.join += " RIGHT JOIN"
	jb.join += " " + table + " ON " + condition
	jb.joinArgs = append(jb.joinArgs, args...)
}

func (jb *joinBuilder) Get() (string, []interface{}) {
	return jb.join, jb.joinArgs
}

func In_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	default:
		panic("unexpected reflect.Kind")
	}

	return
}

func Contains(val interface{}, array interface{}) bool {
	exists, _ := In_array(val, array)
	return exists
}
