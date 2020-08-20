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
)

type ZeroValueType string

type Query struct {
	from string
	//args                   map[string]interface{}
	args            []sql.NamedArg
	wheres          []Where
	values          map[string]interface{}
	SkipZeroValues  bool
	orderBy         []string
	limit           *int
	offset          *int
	groupBy         []string
	columns         []string
	ReturningColumn string
	//tx                     *sqlx.Tx
	model                  interface{}
	errors                 []error
	debug                  bool
	havePrivateTransaction bool
	lama                   *Lama
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

func (q *Query) iArgs() []sql.NamedArg {
	iar := make([]sql.NamedArg, 0)
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
	for k := range keys {
		w := Where{Expr: k, Value: keys[k], Op: "="}
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
		q.args = make([]sql.NamedArg, 0)
	}
	statment := ""
	if q.wheres != nil && len(q.wheres) > 0 {
		for idx := range q.wheres {
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
			whs, args := q.wheres[idx].Build(q.lama.dialect)
			statment = statment + whs + " "
			if args != nil {
				q.args = append(q.args, args...)
			}
			/*for a, b := range args {
				q.args[a] = b
			}*/
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
	if reflect.TypeOf(val).Kind() == reflect.Struct || reflect.TypeOf(val).Kind() == reflect.Ptr {
		values, err := StructToMap(val, false, true, false)
		if err != nil {
			q.addError(err)
			return q
		}
		appendToMap(q.values, values)
		return q
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
		q.limit = &limit
	}
	return q
}
func (q *Query) Offset(off int) *Query {
	if off > 0 {
		q.offset = &off
	}
	return q
}
func (q *Query) Select(cols ...string) *Query {
	q.columns = make([]string, 0)
	q.columns = append(q.columns, cols...)
	return q
}

func (q *Query) ColumnsFromStructOrMap(str interface{}, skipUnTaged bool) *Query {
	q.columns = make([]string, 0)
	if structs.IsStruct(str) {
		for idx := range structs.Fields(str) {
			v:=structs.Fields(str)[idx]
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
		for i := range q.columns {
			ok, fCol := ContainsStrI(fields, q.columns[i], false)
			if ok {
				nCols = append(nCols, fCol)
			} else {
				if !skipNotMatchedColumns {
					nCols = append(nCols, q.columns[i])
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
	if reflect.TypeOf(query).Kind() == reflect.Struct || reflect.TypeOf(query).Kind() == reflect.Ptr {
		values, err := StructToMap(query, true, false, true)
		if err != nil {
			q.addError(err)
			return q
		}
		return q.whereMap(values)
	}

	return q
}
func (q *Query) WhereIn(key string, values ...interface{}) *Query {
	if values == nil || len(values) == 0 {
		return q
	}
	args := make([]sql.NamedArg, 0)
	ins := make([]string, 0)
	for idx := range values {
		nam := "Arg" + strconv.Itoa(idx) + strconv.Itoa(rand.Int())
		args = append(args, sql.NamedArg{Name: nam, Value: values[idx]})
		//args[nam] = v
		ins = append(ins, ":"+nam)
	}
	stm := " " + key + " in(" + strings.Join(ins, ",") + ")"
	return q._where(Where{Raw: stm, Args: args})
}
func (q *Query) WhereOr(w ...Where) *Query {
	return q._where(Where{Or: w})
}
func (q *Query) Model(model interface{}) *Query {
	q.setModel(model)
	q.from = GetTableName(model)
	return q
}
func (q *Query) Table(table string) *Query {
	q.from = table
	return q
}

////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
//actual database queries
////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
func (q *Query) Find(dest interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	//must be set befoure build
	//tx,rolOrCommi,err:=q.getTx()
	if err != nil {
		return err
	}
	if reflect.TypeOf(dest).Kind() == reflect.Ptr {
		elm := reflect.New(reflect.ValueOf(dest).Elem().Type().Elem()).Interface()
		q.setModel(elm)
		slq := SelectQuery{Query: *q}
		stm, args := slq.Build(q.lama.dialect)
		slice := reflect.ValueOf(dest).Interface()
		namedArgs := make([]interface{}, 0)
		for i := range args {
			namedArgs = append(namedArgs, args[i])
		}
		if q.lama.Tx != nil {
			return q.lama.Tx.Select(slice, stm, namedArgs...)
		} else {
			return q.lama.DB.Select(slice, stm, namedArgs...)
		}
	} else {
		elm := reflect.New(reflect.ValueOf(dest).Type().Elem()).Interface()
		q.setModel(elm)
		slice := reflect.MakeSlice(reflect.TypeOf(dest), 0, 0)
		slq := SelectQuery{Query: *q}
		stm, args := slq.Build(q.lama.dialect)
		namedArgs := make([]interface{}, 0)
		for i := range args {
			namedArgs = append(namedArgs, args[i])
		}
		if q.lama.Tx != nil {
			return q.lama.Tx.Select(slice, stm, namedArgs...)
		} else {
			return q.lama.DB.Select(slice, stm, namedArgs...)
		}
	}
}
func (q *Query) Get(dest interface{}) (err error) {
	defer func() {
		//q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	//must be set befoure build
	q.setModel(dest)
	slq := SelectQuery{Query: *q}
	stm, args := slq.Build(q.lama.dialect)
	namedArgs := make([]interface{}, 0)
	for i:= range args {
		namedArgs = append(namedArgs, args[i])
	}
	if q.lama.Tx != nil {
		return q.lama.Tx.Get(dest, stm, namedArgs...)
	} else {
		return q.lama.DB.Get(dest, stm, namedArgs...)
	}
	//err = q.tx.Get(dest, stm, namedArgs...)
	/*if err == sql.ErrNoRows {
		return nil // it will make dangerous effect
	}*/
	return err
}

func (q *Query) First(dest interface{}) (err error) {
	defer func() {
		//q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	if len(q.orderBy) == 0 {
		keys, err := primaryKey(dest)
		if err != nil {
			return err
		}
		for k := range keys {
			q.orderBy = append(q.orderBy, k)
		}
	}

	q.Limit(1)
	return q.Get(dest)
}

func (q *Query) Last(dest interface{}) (err error) {
	defer func() {
		//q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	if len(q.orderBy) == 0 {
		keys, err := primaryKey(dest)
		if err != nil {
			return err
		}
		for k := range keys {
			q.orderBy = append(q.orderBy, k+" desc ")
		}
	}
	q.Limit(1)
	return q.Get(dest)
}

func (q *Query) Count() (count *int64, err error) {
	defer func() {
		//q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if len(q.errors) > 0 {
		return nil, errors.New("more than one error occured:" + q.errors[0].Error())
	}
	key := "*"
	if key == "" {
		key = "*"
	}
	q.columns = make([]string, 0)
	q.columns = append(q.columns, "count("+key+") as Count")
	//var res Count
	var count_ int64
	err = q.Get(&count_)
	if err != nil {
		return nil, err
	}
	return &count_, nil
}
func (q *Query) Sum(column string) (sm *float64, err error) {
	defer func() {
		//q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if len(q.errors) > 0 {
		return nil, errors.New("more than one error occured:" + q.errors[0].Error())
	}
	if column == "" {
		return nil, errors.New("invalied column name")
	}
	q.columns = make([]string, 0)
	q.columns = append(q.columns, "SUM("+column+") as SM")
	var _sm *float64
	err = q.Get(&_sm)
	if err != nil {
		return nil, err
	}
	if _sm == nil {
		var zero float64 = 0
		return &zero, nil
	}
	return _sm, err
}

func (q *Query) CountColumn(dest interface{}, key string) (err error) {
	defer func() {
		//q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	if key == "" {
		key = "*"
	}
	q.columns = make([]string, 0)
	q.columns = append(q.columns, "count("+key+") as COUNT")
	return q.Get(dest)
}

//Save update the holl entity
func (q *Query) Save(entity interface{}) (err error) {
	var tx *Lama
	defer func() {
		if tx != nil {
			if err != nil {
				log.Println("private transaction roolback")
				err = tx.Rollback()
			} else {
				log.Println("private transaction commit")
				err = tx.Commit()
			}
		}
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			} else {
				if err == nil {
					err = errors.New(("painc at save"))
				}
			}
		}
	}()
	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	//must be set befour build
	if q.model == nil {
		q.setModel(entity)
	}
	keys, err := primaryKey(entity)
	if err != nil {
		return err
	}
	q.setValues(entity)
	for k:= range q.values {
		for a := range keys {
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
	stm, args := slq.Build(q.lama.dialect)
	var eff int64 = 0
	ar := make([]interface{}, 0)
	if args != nil {
		for i := range args {
			ar = append(ar, args[i])
		}
	}
	//this function must be save and rollback if have private transaction

	if q.lama.Tx != nil {
		r, err := q.lama.Tx.Exec(stm, ar...)
		if err != nil {
			return err
		}
		eff, err = r.RowsAffected()
		if err != nil {
			return err
		}
	} else {
		tx, err = q.lama.Begin()
		if err != nil {
			return err
		}
		r, err := tx.Tx.Exec(stm, ar...)
		if err != nil {
			return err
		}
		eff, err = r.RowsAffected()
		if err != nil {
			return err
		}
	}
	if eff == 0 {
		err = errors.New("no data updated")
	}
	if eff > 1 {
		err = errors.New("more than one entity operation cancelled ")
	}
	log.Println("rows effected:", eff)
	return err
}

//Delere entity from database
func (q *Query) Delete(entity interface{}) (err error) {
	var tx *Lama
	defer func() {
		if tx != nil {
			if err != nil {
				log.Println("private transaction roolback")
				err = tx.Rollback()
			} else {
				log.Println("private transaction commit")
				err = tx.Commit()
			}
		}
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			} else {
				if err == nil {
					err = errors.New(("painc at save"))
				}
			}
		}
	}()
	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	//must be set befour build
	if q.model == nil {
		q.setModel(entity)
	}
	keys, err := primaryKey(entity)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return errors.New("primary key is missing")
	}
	slq := DeleteQuery{Query: *q}
	slq.Where(keys)
	stm, args := slq.Build(q.lama.dialect)
	var eff int64 = 0
	ar := make([]interface{}, 0)
	if args != nil {
		for i := range args {
			ar = append(ar, args[i])
		}
	}
	//this function must be save and rollback if have private transaction

	if q.lama.Tx != nil {
		r, err := q.lama.Tx.Exec(stm, ar...)
		if err != nil {
			return err
		}
		eff, err = r.RowsAffected()
		if err != nil {
			return err
		}
	} else {
		tx, err = q.lama.Begin()
		if err != nil {
			return err
		}
		r, err := tx.Tx.Exec(stm, ar...)
		if err != nil {
			return err
		}
		eff, err = r.RowsAffected()
		if err != nil {
			return err
		}
	}
	if eff == 0 {
		err = errors.New("no data updated")
	}
	if eff > 1 {
		err = errors.New("more than one entity operation cancelled ")
	}
	log.Println("rows effected:", eff)
	return err
}

//Update save partial data to entity to database
func (q *Query) Update(data map[string]interface{}, acceptBulk bool) (err error) {
	defer func() {
		//q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()

	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	//must be set befour build
	q.setValues(data)
	if len(q.wheres) == 0 && !acceptBulk {
		return errors.New("bulk update not allowed")
	}
	slq := UpdateQuery{Query: *q}
	stm, args := slq.Build(q.lama.dialect)
	var eff int64 = 0
	ar := make([]interface{}, 0)
	if args != nil {
		for i:= range args {
			ar = append(ar, args[i])
		}
	}
	if q.lama.Tx != nil {
		r, err := q.lama.Tx.Exec(stm, ar...)
		if err != nil {
			return err
		}
		eff, err = r.RowsAffected()
	} else {
		r, err := q.lama.DB.Exec(stm, ar...)
		if err != nil {
			return err
		}
		eff, err = r.RowsAffected()
	}

	log.Println("rows effected:", eff)
	return err
}
//Add insert new entity into the database
func (q *Query) Add(entity interface{}) (err error) {
	defer func() {
		//q.FinalizeWith(err)
		if r := recover(); r != nil {
			log.Println("panic:", r)
			errr, ok := r.(error)
			if ok {
				err = errr
			}
		}
	}()
	if len(q.errors) > 0 {
		return errors.New("more than one error occured:" + q.errors[0].Error())
	}
	//must be set befour build
	if q.model == nil {
		q.setModel(entity)
	}
	q.setValues(entity)
	slq := InsertQuery{Query: *q}
	stm, args := slq.Build(q.lama.dialect)
	var eff int64 = 0
	ar := make([]interface{}, 0)
	for i:= range args {
		ar = append(ar, args[i])
	}
	if q.lama.Tx != nil {
		r, err := q.lama.Tx.Exec(stm, ar...)
		if err != nil {
			return err
		}
		eff, err = r.RowsAffected()
	} else {
		r, err := q.lama.DB.Exec(stm, ar...)
		if err != nil {
			return err
		}
		eff, err = r.RowsAffected()
	}
	log.Println("rows effected:", eff)
	return err
}
//setModel set the model used to find tablename and  generate colum names
func (q *Query) setModel(dest interface{}) {
	if reflect.TypeOf(dest).Kind() == reflect.Struct {
		q.model = dest
	}
	if reflect.TypeOf(dest).Kind() == reflect.Ptr {
		elm := reflect.ValueOf(dest).Elem().Interface()
		if reflect.TypeOf(elm).Kind() == reflect.Struct {
			q.model = elm
		}
	}
}

/*func (q *Query) Finalize(commit bool) error {
	if q.havePrivateTransaction && q.tx != nil {
		if commit {
			return q.tx.Commit()
		} else {
			return q.tx.Rollback()
		}
	}
	return nil
}

/*func (q *Query) FinalizeWith(err error) error {
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
}*/

/*func (q *Query) getTx() (*sqlx.Tx,func(error),error) {
	if q.lama.Tx == nil {
		tx, err := q.lama.DB.Beginx()
		if err != nil {
			q.addError(err)
			return nil,nil,err
		} else {
			q.havePrivateTransaction = true
			return tx, func(err error) {
				if err!=nil{
					tx.Rollback()
				} else {
					tx.Commit()
				}
			},nil
		}
	} else {
		q.havePrivateTransaction = false
		return q.lama.Tx, func(err error) {

		},nil
	}
}*/
