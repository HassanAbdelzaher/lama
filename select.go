package lama

import (
	"database/sql"
	"log"
	"strings"
)

type SelectQuery struct {
	Query
}

func (q *SelectQuery) Build(di Dialect) (string, []sql.NamedArg) {
	statment := "select "

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
	statment = di.LimitAndOffsetSQL(statment,q.limit,q.offset)
	if q.debug &&!di.HaveLog() {
		log.Println(statment, q.args)
	}
	return statment, q.iArgs()
}
