package lama

import (
	"log"
)

type UpdateQuery struct {
	Query
}

func (q *UpdateQuery) Build() (string, map[string]interface{}) {
	if q.args == nil {
		q.args = make(map[string]interface{}, 0)
	}
	statment := "update "
	frm := q.getFrom()
	if frm != "" {
		statment = statment + " " + frm + " "
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
		q.args[k] = v
	}
	where := q.buildWhere()
	statment = statment + where
	if q.debug {
		log.Println(statment, q.args)
	}
	sArgs := q.iArgs()
	return statment, sArgs
}
