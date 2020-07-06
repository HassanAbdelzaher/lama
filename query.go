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
)

type ZeroValueType string

type Query struct {
	from                   string
	args                   map[string]interface{}
	wheres                 []Where
	values                 map[string]interface{}
	SkipZeroValues         bool
	orderBy                []string
	limit                  int
	offset                 int
	groupBy                []string
	columns                []string
	ReturningColumn        string
	tx                     *sqlx.Tx
	model                  interface{}
	errors                 []error
	debug                  bool
	havePrivateTransaction bool
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

func (q *Query) iArgs() map[string]interface{} {
	iar := make(map[string]interface{})
	if q.args == nil {
		return iar
	}

	return q.args
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

func (q *Query) buildWhere() string {
	if q.args == nil {
		q.args = make(map[string]interface{})
	}
	statment := ""
	if q.wheres != nil && len(q.wheres) > 0 {
		for idx, w := range q.wheres {
			/*isZero:=false
			if w.Value!=nil{
				isZero=reflect.ValueOf(w.Value).IsZero()
			}
			if isZero && q.SkipZeroValues {
				continue
			}*/
			if idx == 0 {
				statment = statment + " where  "
			} else {
				statment = statment + " and  "
			}
			whs, args := w.Build()
			statment = statment + whs + " "
			for a, b := range args {
				q.args[a] = b
			}
			log.Println(q.args)
			//q.args = append(q.args, args...)
		}
	}
	return statment
}

func (q *Query) setValues(val interface{}) *Query {
	if val == nil {
		return q
	}
	if q.values == nil {
		q.values = make(map[string]interface{})
	}
	if reflect.TypeOf(val).Kind() == reflect.Map {
		mValues, ok := val.(map[string]interface{})
		if !ok {
			q.addError(errors.New("map must be map[string]interface"))
			return q
		}
		appendToMap(q.values, mValues)
		return q
	}
	if reflect.TypeOf(val).Kind() == reflect.Struct {
		values, err := StructToMap(val, false, true, false)
		if err != nil {
			q.addError(err)
			return q
		}
		appendToMap(q.values, values)
		return q
	}
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		strct := reflect.ValueOf(val).Elem().Interface()
		if reflect.TypeOf(strct).Kind() == reflect.Ptr {
			q.addError(errors.New("pointer to pointer not supported"))
			return q
		}
		return q.setValues(strct)
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
			for k := range mp {
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
		fields := Map(structs.Fields(str), func(i interface{}) interface{} {
			f := i.(*structs.Field)
			return f.Tag("db")
		})
		for _, col := range q.columns {
			ok, fCol := ContainsStrI(fields, col, false)
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
		mArgs := make(map[string]interface{})
		for _, b := range args {
			mArgs[b.Name] = b.Value
		}
		return q._where(Where{Raw: str, Args: mArgs})
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
		values, err := StructToMap(query, true, false, true)
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
	args := make(map[string]interface{})
	ins := make([]string, 0)
	for idx, v := range values {
		nam := "Arg" + strconv.Itoa(idx) + strconv.Itoa(rand.Int())
		//args = append(args, sql.NamedArg{Name: nam, Value: v})
		args[nam] = v
		ins = append(ins, ":"+nam)
	}
	stm := " " + key + " in(" + strings.Join(ins, ",") + ")"
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

////////////////////////////////////
////////////////////////////////////////////////
//actual database queries
func (q *Query) Find(dest interface{}) (err error) {
	defer func() {
		q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if q.tx == nil {
		return errors.New("no database connection defined")
	}
	//must be set befoure build
	if reflect.TypeOf(dest).Kind() == reflect.Ptr {
		elm := reflect.New(reflect.ValueOf(dest).Elem().Type().Elem()).Interface()
		q.model = elm
		slq := SelectQuery{Query: *q}
		stm, args := slq.Build()
		slice := reflect.ValueOf(dest).Interface()
		namedArgs := make([]interface{}, 0)
		for a, b := range args {
			namedArgs = append(namedArgs, sql.NamedArg{Name: a, Value: b})
		}
		return q.tx.Select(slice, stm, namedArgs...)
	} else {
		elm := reflect.New(reflect.ValueOf(dest).Type().Elem()).Interface()
		q.model = elm
		slice := reflect.MakeSlice(reflect.TypeOf(dest), 0, 0)
		slq := SelectQuery{Query: *q}
		stm, args := slq.Build()
		namedArgs := make([]interface{}, 0)
		for a, b := range args {
			namedArgs = append(namedArgs, sql.NamedArg{Name: a, Value: b})
		}
		err = q.tx.Select(slice, stm, namedArgs...)
		return err
	}
}
func (q *Query) Get(dest interface{}) (err error) {
	defer func() {
		q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if q.tx == nil {
		return errors.New("no database connection defined")
	}
	//must be set befoure build
	q.model = dest
	slq := SelectQuery{Query: *q}
	stm, args := slq.Build()
	namedArgs := make([]interface{}, 0)
	for a, b := range args {
		namedArgs = append(namedArgs, sql.NamedArg{Name: a, Value: b})
	}
	err = q.tx.Get(dest, stm, namedArgs...)
	return err
}

func (q *Query) First(dest interface{}) (err error) {
	defer func() {
		q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if q.tx == nil {
		return errors.New("no database connection defined")
	}
	if len(q.orderBy) == 0 {
		keys, err := primaryKey(dest)
		if err != nil {
			return err
		}
		for k, _ := range keys {
			q.orderBy = append(q.orderBy, k)
		}
	}

	q.Limit(1)
	return q.Get(dest)
}

func (q *Query) Last(dest interface{}) (err error) {
	defer func() {
		q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if q.tx == nil {
		return errors.New("no database connection defined")
	}
	if len(q.orderBy) == 0 {
		keys, err := primaryKey(dest)
		if err != nil {
			return err
		}
		for k, _ := range keys {
			q.orderBy = append(q.orderBy, k+" desc ")
		}
	}
	q.Limit(1)
	return q.Get(dest)
}

//save entity
func (q *Query) Save(entity interface{}) (err error) {
	defer func() {
		q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if q.tx == nil {
		return errors.New("no database connection defined")
	}
	//must be set befour build
	if q.model == nil {
		q.model = entity
	}
	keys, err := primaryKey(entity)
	if err != nil {
		return err
	}
	log.Println("primary", keys)
	q.setValues(entity)
	for k, _ := range q.values {
		for a, _ := range keys {
			if a == k {
				delete(q.values, k) //delete primary keys
			}
		}
	}
	slq := UpdateQuery{Query: *q}
	slq.Where(keys)
	if len(keys) == 0 {
		return errors.New("primary key is missing")
	}
	stm, args := slq.Build()
	nArgs := make(map[string]interface{})
	for a, b := range args {
		nArgs[a] = b
	}
	r, err := q.tx.NamedExec(stm, args)
	if err != nil {
		return err
	}
	eff, err := r.RowsAffected()
	if eff == 0 {
		err = errors.New("no data updated")
	}
	if eff > 1 {
		err = errors.New("more than one entity operation cancelled ")
	}
	log.Println("rows effected:", eff)
	return err
}

//save entity
func (q *Query) Update(data map[string]interface{}, acceptBulk bool) (err error) {
	defer func() {
		q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if q.tx == nil {
		return errors.New("no database connection defined")
	}
	//must be set befour build
	q.setValues(data)
	if len(q.wheres) == 0 && !acceptBulk {
		return errors.New("bulk update not allowed")
	}
	slq := UpdateQuery{Query: *q}
	stm, args := slq.Build()
	nArgs := make(map[string]interface{})
	for a, b := range args {
		nArgs[a] = b
	}
	r, err := q.tx.NamedExec(stm, args)
	if err != nil {
		return err
	}
	eff, err := r.RowsAffected()
	log.Println("rows effected:", eff)
	return err
}

func (q *Query) Add(entity interface{}) (err error) {
	defer func() {
		q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if q.tx == nil {
		return errors.New("no database connection defined")
	}
	//must be set befour build
	if q.model == nil {
		q.model = entity
	}
	q.setValues(entity)
	slq := InsertQuery{Query: *q}
	stm, args := slq.Build()
	nArgs := make(map[string]interface{})
	for a, b := range args {
		nArgs[a] = b
	}
	r, err := q.tx.NamedExec(stm, args)
	if err != nil {
		return err
	}
	eff, err := r.RowsAffected()
	log.Println("rows effected:", eff)
	return err
}

func (q *Query) Finalize(commit bool) error {
	if q.havePrivateTransaction && q.tx != nil {
		if commit {
			return q.tx.Commit()
		} else {
			return q.tx.Rollback()
		}
	}
	return nil
}

func (q *Query) FinalizeWith(err error) error {
	commit := true
	if err != nil {
		commit = false
	}
	if q.havePrivateTransaction && q.tx != nil {
		defer func() {
			q.tx = nil
		}()
		if commit {
			return q.tx.Commit()
		} else {
			return q.tx.Rollback()
		}
	}
	return nil
}
