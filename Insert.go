package lama

import (
	"database/sql"
	"log"
	"reflect"
)

type InsertQuery struct {
	Query
}

func (q *InsertQuery) Build(di Dialect) (string, []sql.NamedArg) {
	if q.args == nil {
		q.args = make([]sql.NamedArg, 0)
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
	for k := range q.values {
		v:=q.values[k]
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
		icols = icols + k
		values = values + di.BindVarStr(k)
		//q.args["@"+k] = v
		q.args=append(q.args,sql.NamedArg{Name:k,Value:v})
	}
	icols = icols + ")";
	if q.ReturningColumn != "" {
		icols = icols + di.LastInsertIDOutputInterstitial(q.getFrom(),q.ReturningColumn,q.columns)
	}
	values = values + ")"
	statment = statment + icols + values
	if q.debug && !di.HaveLog() {
		log.Println(statment, q.args)
	}
	return statment, q.iArgs()
}
