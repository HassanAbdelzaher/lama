package lama

import (
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"mas.com/wsdl/slices"
)

type ZeroValueType string

type Query struct {
	from            string
	args            []sql.NamedArg
	wheres          []Where
	values          map[string]interface{}
	SkipZeroValues  bool
	orderBy         []string
	limit           int
	offset          int
	groupBy         []string
	columns         []string
	ReturningColumn string
	db              *sqlx.DB
	model           interface{}
	errors          []error
	debug           bool
}

func (q *Query) Debug(dbg bool) *Query {
	q.debug = dbg
	return q
}

func (q *Query) getFrom() string {
	if len(q.from) > 0 {
		return q.from
	}
	if q.model != nil {
		q.from = GetTableName(q.model)
		return q.from
	}
	q.addError(errors.New("no table name defined"))
	return ""
}

func (q *Query) iArgs() []interface{} {
	iar := make([]interface{}, 0)
	if q.args == nil {
		return iar
	}
	for _, ar := range q.args {
		iar = append(iar, ar)
	}
	return iar
}

func (q *Query) addError(err error) {
	if q.errors == nil {
		q.errors = make([]error, 0)
	}
	if err != nil {
		q.errors = append(q.errors, err)
	}
}

func (q *Query) whereMap(keys map[string]interface{}) *Query {
	if q.wheres == nil {
		q.wheres = make([]Where, 0)
	}
	for k, v := range keys {
		w := Where{Expr: k, Value: v, Op: "="}
		q._where(w)
	}
	return q
}

func (q *Query) _where(w Where) *Query {
	if q.wheres == nil {
		q.wheres = make([]Where, 0)
	}
	//w:=Where{Expr:key,Value:value,Op:op}
	q.wheres = append(q.wheres, w)
	return q
}

func (q *Query) _values(values map[string]interface{}) *Query {
	q.values = values
	return q
}

func (q *Query) _StructValues(stct interface{}) *Query {
	mp, err := StructToMap(stct, false)
	if err != nil {
		q.addError(err)
		return q
	}
	q._values(mp)
	return q
}
func (q *Query) buildWhere() string {
	statment := ""
	if q.wheres != nil && len(q.wheres) > 0 {
		for idx, w := range q.wheres {
			isZero := false
			if w.Value != nil {
				isZero = reflect.ValueOf(w.Value).IsZero()
			}
			if isZero && q.SkipZeroValues {
				continue
			}
			if idx == 0 {
				statment = statment + " where  "
			} else {
				statment = statment + " and  "
			}
			whs, args := w.Build()
			statment = statment + whs + " "
			q.args = append(q.args, args...)
		}
	}
	return statment
}
func (q *Query) Values(values interface{}, skipZeroValues bool) *Query {
	if values == nil {
		return q
	}
	if reflect.TypeOf(values).Kind() == reflect.Map {
		mValues, ok := values.(map[string]interface{})
		if !ok {
			q.addError(errors.New("map must be map[string]interface"))
			return q
		}
		return q._values(mValues)
	}
	if reflect.TypeOf(values).Kind() == reflect.Struct {
		values, err := StructToMap(values, skipZeroValues)
		if err != nil {
			q.addError(err)
			return q
		}
		return q._values(values)
	}
	if reflect.TypeOf(values).Kind() == reflect.Ptr {
		val := reflect.ValueOf(values).Elem()
		if reflect.TypeOf(val).Kind() == reflect.Ptr {
			q.addError(errors.New("pointer to pointer not supported"))
			return q
		}
		return q.Values(val, skipZeroValues)
	}
	q.addError(errors.New("values must be map or struct"))
	return q
}
func (q *Query) OrderBy(by ...string) *Query {
	if q.orderBy == nil {
		q.orderBy = make([]string, 0)
	}
	q.orderBy = append(q.orderBy, by...)
	return q
}
func (q *Query) Limit(limit int) *Query {
	if limit > 0 {
		q.limit = limit
	}
	return q
}
func (q *Query) Offset(off int) *Query {
	if off > 0 {
		q.offset = off
	}
	return q
}
func (q *Query) Select(cols ...string) *Query {
	q.columns = make([]string, 0)
	q.columns = append(q.columns, cols...)
	return q
}
func (q *Query) Count(key string) *Query {
	if key == "" {
		key = "*"
	}
	q.columns = make([]string, 0)
	q.columns = append(q.columns, "count("+key+") as COUNT")
	return q
}
func (q *Query) ColumnsFromStructOrMap(str interface{}, skipUnTaged bool) *Query {
	q.columns = make([]string, 0)
	if structs.IsStruct(str) {
		for _, v := range structs.Fields(str) {
			tag := v.Tag("db")
			name := v.Name()
			if tag != "" {
				q.columns = append(q.columns, tag)
			} else {
				if !skipUnTaged {
					q.columns = append(q.columns, name)
				}
			}
		}
	} else {
		mp, isMap := str.(map[string]interface{})
		if isMap {
			for k, _ := range mp {
				q.columns = append(q.columns, k)
			}
		}
	}
	return q
}
func (q *Query) AdaptColumnNamesToStruct(str interface{}, skipNotMatchedColumns bool) *Query {
	if str == nil || q.columns == nil {
		return q
	}
	if structs.IsStruct(str) {
		nCols := make([]string, 0)
		fields := slices.Map(structs.Fields(str), func(i interface{}) interface{} {
			f := i.(*structs.Field)
			return f.Tag("db")
		})
		for _, col := range q.columns {
			ok, fCol := slices.ContainsStrI(fields, col, false)
			if ok {
				nCols = append(nCols, fCol)
			} else {
				if !skipNotMatchedColumns {
					nCols = append(nCols, col)
				}
			}
		}
		q.columns = nCols
	}
	return q
}
func (q *Query) Where(query interface{}, args ...sql.NamedArg) *Query {
	if query == nil {
		q.addError(errors.New("invalied null query"))
		return q
	}
	if reflect.TypeOf(query).Kind() == reflect.String {
		str := reflect.ValueOf(query).String()
		return q._where(Where{Raw: str, Args: args})
	}
	if reflect.TypeOf(query).Kind() == reflect.Map {
		values, ok := query.(map[string]interface{})
		if !ok {
			q.addError(errors.New("map must be map[string]interface"))
			return q
		}
		return q.whereMap(values)
	}
	if reflect.TypeOf(query).Kind() == reflect.Struct {
		values, err := StructToMap(query, true)
		if err != nil {
			q.addError(err)
			return q
		}
		return q.whereMap(values)
	}
	if reflect.TypeOf(query).Kind() == reflect.Ptr {
		val := reflect.ValueOf(query).Elem()
		if reflect.TypeOf(val).Kind() == reflect.Ptr {
			q.addError(errors.New("pointer to pointer not supported"))
			return q
		}
		return q.Where(val, args...)
	}
	return q
}
func (q *Query) WhereIn(key string, values ...interface{}) *Query {
	if values == nil || len(values) == 0 {
		return q
	}
	args := make([]sql.NamedArg, 0)
	ins := make([]string, 0)
	for idx, v := range values {
		nam := "Arg" + strconv.Itoa(idx) + strconv.Itoa(rand.Int())
		args = append(args, sql.NamedArg{Name: nam, Value: v})
		ins = append(ins, ":"+nam)
	}
	stm := " " + key + " in(" + strings.Join(ins, ",") + ")"
	log.Println("in", stm, args)
	return q._where(Where{Raw: stm, Args: args})
}
func (q *Query) WhereOr(w ...Where) *Query {
	return q._where(Where{Or: w})
}
func (q *Query) Model(model interface{}) *Query {
	q.model = model
	q.from = GetTableName(model)
	return q
}
func (q *Query) Table(table string) *Query {
	q.from = table
	return q
}
func (q *Query) Find(dest interface{}) error {
	if q.db == nil {
		return errors.New("no database connection defined")
	}
	//must be set befoure build
	if reflect.TypeOf(dest).Kind() == reflect.Ptr {
		elm := reflect.New(reflect.ValueOf(dest).Elem().Type().Elem()).Interface()
		q.model = elm
		slq := SelectQuery{Query: *q}
		stm, args := slq.Build()
		slice := reflect.ValueOf(dest).Interface()
		log.Println("slice:", slice, reflect.TypeOf(slice))
		return q.db.Select(slice, stm, args...)
	} else {
		elm := reflect.New(reflect.ValueOf(dest).Type().Elem()).Interface()
		q.model = elm
		log.Println("reflect.TypeOf(elm)", reflect.TypeOf(elm))
		slice := reflect.MakeSlice(reflect.TypeOf(dest), 0, 0)
		log.Println("slice:", slice, reflect.TypeOf(slice))
		slq := SelectQuery{Query: *q}
		stm, args := slq.Build()
		return q.db.Select(slice, stm, args...)
	}
}
func (q *Query) Get(dest interface{}) error {
	if q.db == nil {
		return errors.New("no database connection defined")
	}
	//must be set befoure build
	q.model = dest
	slq := SelectQuery{Query: *q}
	stm, args := slq.Build()
	return q.db.Get(dest, stm, args...)
}
func (q *Query) Insert(dest interface{}) error {
	if q.db == nil {
		return errors.New("no database connection defined")
	}
	//must be set befour build
	q.model = dest
	q.Values(dest, true)
	slq := InsertQuery{Query: *q}
	stm, _ := slq.Build()
	_, err := q.db.NamedExec(stm, q.values)
	return err
}
