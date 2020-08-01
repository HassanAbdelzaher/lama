package lama

import (
	"database/sql"
	"log"
)

type DeleteQuery struct {
	Query
}

func (q *DeleteQuery) Build() (string, []sql.NamedArg) {
	if q.args == nil {
		q.args = make([]sql.NamedArg, 0)
	}
	statment := "delete from "
	if q.from != "" {
		statment = statment + q.from + " "
	}
	where := q.buildWhere()
	statment = statment + " " + where
	if q.debug {
		log.Println(statment, q.args)
	}
	return statment, q.iArgs()
}
