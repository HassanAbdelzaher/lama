package lama

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

type SelectQuery struct {
	Query
}

func (q *SelectQuery) Build() (string, []sql.NamedArg) {
	statment := "select "
	if q.limit > 0 {
		statment = statment + " top " + strconv.Itoa(q.limit) + " "
	}
	if q.model != nil && (q.columns == nil || len(q.columns) == 0) {
		q.ColumnsFromStructOrMap(q.model, false)
	} else {
		if q.model != nil {
			q.AdaptColumnNamesToStruct(q.model, false)
		}
	}
	if q.columns != nil && len(q.columns) > 0 {
		cols := strings.Join(q.columns, ",")
		statment = statment + cols + " "
	}
	frm := q.getFrom()
	if frm != "" {
		statment = statment + " from " + frm + " "
	}
	where := q.buildWhere()
	statment = statment + where
	if q.orderBy != nil && len(q.orderBy) > 0 {
		ordBy := strings.Join(q.orderBy, ",")
		statment = statment + " order by " + ordBy + " "
	}
	if q.offset > 0 {
		statment = statment + " top " + strconv.Itoa(q.limit) + " "
	}
	if q.debug {
		log.Println(statment, q.args)
	}
	return statment, q.iArgs()
}
