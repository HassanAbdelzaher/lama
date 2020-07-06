package lama

import (
	"log"
)

type DeleteQuery struct {
	Query
}

func (q *DeleteQuery) Build() (string, map[string]interface{}) {
	if q.args == nil {
		q.args = make(map[string]interface{}, 0)
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
