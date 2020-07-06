package lama

import (
	"log"
	"reflect"
)

type InsertQuery struct {
	Query
}

func (q *InsertQuery) Build() (string, map[string]interface{}) {
	if q.args == nil {
		q.args = make(map[string]interface{}, 0)
	}
	statment := "insert into "
	frm := q.getFrom()
	if frm != "" {
		statment = statment + frm
	}
	if q.values == nil {
		q.values = make(map[string]interface{})
	}
	icols := "("
	values := " values("
	itr := 0
	for k, v := range q.values {
		if q.SkipZeroValues && v != nil {
			isZero := reflect.ValueOf(v).IsZero()
			if isZero {
				continue
			}
		}
		itr++
		if itr > 1 {
			icols = icols + ","
			values = values + ","
		}
		icols = icols + " " + k
		values = values + " :" + k
		q.args[k] = v
	}
	if q.ReturningColumn != "" {
		icols = icols + ") OUTPUT Inserted." + q.ReturningColumn + " "

	} else {
		icols = icols + ") "
	}
	values = values + ")"
	statment = statment + icols + values
	if q.debug {
		log.Println(statment, q.args)
	}
	return statment, q.iArgs()
}
