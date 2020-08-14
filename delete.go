package lama

import (
	"database/sql"
	"log"
)

type DeleteQuery struct {
	Query
}

func (q *DeleteQuery) Build(di Dialect) (string, []sql.NamedArg) {
	if q.args == nil {
		q.args = make([]sql.NamedArg, 0)
	}
	frm := q.getFrom()
	statment := "delete from "
	statment = statment + frm + " "
	where := q.buildWhere()
	statment = statment + " " + where
	if q.debug && !di.HaveLog() {
		log.Println(statment, q.args)
	}
	return statment, q.iArgs()
}
