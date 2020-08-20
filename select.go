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
	if q.groupBy != nil && len(q.groupBy) > 0 {
		gBy := strings.Join(q.groupBy, ",")
		statment = statment + " group by " + gBy + " "
	}
	if q.havings != nil && len(q.havings) > 0 {
		if q.args == nil {
			q.args = make([]sql.NamedArg, 0)
		}
		for i := range q.havings {
			v := q.havings[i]
			statment = statment + " having " + v.Key + " "
			if v.Args != nil && len(v.Args) > 0 {
				for j := range v.Args {
					va := v.Args[j]
					q.args = append(q.args, va)
				}
			}
		}
	}
	if q.orderBy != nil && len(q.orderBy) > 0 {
		ordBy := strings.Join(q.orderBy, ",")
		statment = statment + " order by " + ordBy + " "
	}
	statment = di.LimitAndOffsetSQL(statment, q.limit, q.offset)
	if q.debug && !di.HaveLog() {
		log.Println(statment, q.args)
	}
	return statment, q.iArgs()
}
