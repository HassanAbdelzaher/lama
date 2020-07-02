package lama

import (
	"database/sql"
	"log"
)

type UpdateQuery struct {
	Query
}

func (q *UpdateQuery) Build() (string, []interface{}) {
	q.args = make([]sql.NamedArg, 0)
	statment := "update "
	if q.from != "" {
		statment = statment + q.from + " "
	}
	statment = statment + " set "
	itr := 0
	for k, v := range q.values {
		/*isZero:=reflect.ValueOf(v).IsZero()
		if isZero {
			continue
		}*/
		//never skip zero values because you may be need to reset column value
		itr++
		if itr > 1 {
			statment = statment + ","
		}
		statment = statment + k + "=:" + k
		q.args = append(q.args, sql.NamedArg{Name: k, Value: v})
	}
	where := q.buildWhere()
	statment = statment + where
	if q.debug {
		log.Println(statment, q.args)
	}

	return statment, q.iArgs()
}
