package lama

import (
	"database/sql"
	"log"
	"reflect"
)

type InsertQuery struct {
	Query
}

func (q *InsertQuery) Build() (string, []interface{}) {
	q.args = make([]sql.NamedArg, 0)
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
		q.args = append(q.args, sql.NamedArg{Name: k, Value: v})
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
