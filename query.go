package lama

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/HassanAbdelzaher/lama/structs"
	"log"
	"math/rand"
	"reflect"
)

type Having struct {
	Key  string
	Args []sql.NamedArg
}

type ZeroValueType string


func getArgName(key string) string{
	nam:=fmt.Sprintf(`%s%d`,key,rand.Int31n(1000000))
	return nam;
}

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
	havings                []Having
	selectedZeroValues []string
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
		_idx:=0
		for id := range q.wheres {
			wh:=q.wheres[id]
			if wh.Fake{
				continue
			}
			if _idx == 0 {
				statment = statment + " where  "
			} else {
				statment = statment + " and  "
			}
			_idx++
			whs, args := wh.Build(q.lama.dialect)
			statment = statment + whs + " "
			if args != nil {
				q.args = append(q.args, args...)
			}
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
		values, err := StructToMap(val, false, true, false, false,q.selectedZeroValues)
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
func (q *Query) SelectedZeroValues(col ...string) *Query {
	if q.selectedZeroValues == nil {
		q.selectedZeroValues = make([]string, 0)
	}
	q.selectedZeroValues = append(q.selectedZeroValues, col...)
	return q
}
func (q *Query) GroupBy(by ...string) *Query {
	if q.groupBy == nil {
		q.groupBy = make([]string, 0)
	}
	q.groupBy = append(q.groupBy, by...)
	return q
}
func (q *Query) Having(expr string, args ...sql.NamedArg) *Query {
	if q.havings == nil {
		q.havings = make([]Having, 0)
	}
	q.havings = append(q.havings, Having{Key: expr, Args: args})
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
		mp := structs.New(str, structs.MapOptions{
			SkipZeroValue: false,
			UseFieldName:  false,
			SkipUnTaged:   false,
			SkipComputed:  false,
			Flatten:       true,
		}).Map()
		for idx := range mp {
			q.columns = append(q.columns, idx)
			/*v:=structs.Fields(str)[idx]
			tag := v.Tag("db")
			name := v.Name()
			if tag != "" {
				q.columns = append(q.columns, tag)
			} else {
				if !skipUnTaged {
					q.columns = append(q.columns, name)
				}
			}*/
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
		fields := Map(structs.Fields(str, structs.MapOptions{}), func(i interface{}) interface{} {
			f := i.(*structs.Field)
			return f.Tag("grom")
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
		//q.addError(errors.New("invalied null query"))
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
	val:=reflect.ValueOf(query)
	if val.Kind()==reflect.Ptr{
		val=val.Elem()
	}
	if reflect.TypeOf(val).Kind() == reflect.Struct {
		wh,ok:=val.Interface().(Where)
		if ok{
			q._where(wh)
			return q
		}
		values, err := StructToMap(query, true, false, true, true,q.selectedZeroValues)
		if err != nil {
			q.addError(err)
			return q
		}
		return q.whereMap(values)
	}

	return q
}
func (q *Query) WhereOr(w ...Where) *Query {
	return q._where(Where{Or: w})
}

func (q *Query) In(key string,valus ...interface{}) *Query {
	return q.Where(In(q.lama.dialect,key,valus))
}
func (q *Query) Between(key string,valu1 interface{},valu2 interface{}) *Query {
	return q.Where(Between(q.lama.dialect,key,valu1,valu2))
}
func (q *Query) Like(key string,value string) *Query {
	return q.Where(Like(key,value))
}
func (q *Query) StartsWith(key string,value string) *Query {
	return q.Where(StartsWith(key,value))
}
func (q *Query) EndsWith(key string,value string) *Query {
	return q.Where(EndsWith(key,value))
}
func (q *Query) Gt(key string,value interface{}) *Query {
	return q.Where(Gt(key,value))
}
func (q *Query) Gte(key string,value interface{}) *Query {
	return q.Where(Gte(key,value))
}
func (q *Query) Lt(key string,value interface{}) *Query {
	return q.Where(Lt(key,value))
}
func (q *Query) Lte(key string,value interface{}) *Query {
	return q.Where(Lte(key,value))
}
func (q *Query) Eq(key string,value interface{}) *Query {
	return q.Where(Eq(key,value))
}
func (q *Query) NotEq(key string,value interface{}) *Query {
	return q.Where(NotEq(key,value))
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
	var sliceType reflect.Type
	isPointer :=false
	if reflect.TypeOf(dest).Kind() == reflect.Ptr{
		isPointer=true
		sliceType=reflect.TypeOf(dest).Elem()
	}else{
		sliceType=reflect.TypeOf(dest)
	}
	if sliceType.Kind()!=reflect.Slice {
		return errors.New("lama:destination must be a slice")
	}
	//set model
	if q.model==nil{
		inSlc:=reflect.TypeOf(dest).Elem()
		if inSlc.Kind()!=reflect.Slice {
			return errors.New("lama:destination must be a slice")
		}
		slcItm:=inSlc.Elem()
		if slcItm.Kind()==reflect.Ptr{
			elm := reflect.New(slcItm.Elem()).Interface()
			q.setModel(elm)
		}else{
			elm := reflect.New(slcItm).Interface()
			q.setModel(elm)
		}
	}
	slq := SelectQuery{Query: *q}
	stm, args := slq.Build(q.lama.dialect)
	namedArgs := make([]interface{}, 0)
	for i := range args {
		namedArgs = append(namedArgs, args[i])
	}
	var nwDest interface{}
	if isPointer{
		//nwDest=reflect.ValueOf(dest).Interface()
		nwDest=dest
	}else{
		nwDest = reflect.MakeSlice(reflect.TypeOf(dest), 0, 0)
	}
	if q.lama.Tx != nil {
		return q.lama.Tx.Select(nwDest, stm, namedArgs...)
	} else {
		return q.lama.DB.Select(nwDest, stm, namedArgs...)
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
	for i := range args {
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
		keys, err := structs.PrimaryKey(dest)
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			q.orderBy = append(q.orderBy, "1")
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
		keys, err := structs.PrimaryKey(dest)
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
	var tx *LamaTx
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
	keys, err := structs.PrimaryKey(entity)
	if err != nil {
		return err
	}
	q.setValues(entity)
	for k := range q.values {
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
		err = errors.New("more than one entity match primary key: operation cancelled ")
	}
	log.Println("rows effected:", eff)
	return err
}

func (q *Query) Upsert(entity interface{}) (err error) {
	defer func() {
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
	keys, err := structs.PrimaryKey(entity)
	if err != nil {
		return err
	}
	cq:=*q
	cnt,err:=cq.Where(keys).Count()
	if err!=nil{
		return err
	}
	if cnt==nil || *cnt==0{
		return q.Add(entity)
	}else {
		if *cnt>1{
			return errors.New("more than on item match the primary key")
		}
		return q.Save(entity)
	}
	return err
}

func (q *Query) AddIfNotExists(entity interface{}) (err error) {
	defer func() {
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
	keys, err := structs.PrimaryKey(entity)
	if err != nil {
		return err
	}
	cq:=*q
	cnt,err:=cq.Where(keys).Count()
	if err!=nil{
		return err
	}
	if cnt==nil || *cnt==0{
		return q.Add(entity)
	}else {
		//do nothing so never update if entity exsists
	}
	return err
}

//Delete delete entity match where from database

func (q *Query) Delete(entity interface{}) (err error) {
	var tx *LamaTx
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
	keys, err := structs.PrimaryKey(entity)
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
//DeleteAll  entities match where creteria from database
func (q *Query) DeleteAll() (err error) {
	var tx *LamaTx
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
	slq := DeleteQuery{Query: *q}
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
		for i := range args {
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
	for i := range args {
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

//setModel set the model used to find tablename and  generate column names
func (q *Query) setModel(_dest interface{}) {
	elm := _dest
	if reflect.TypeOf(elm).Kind() == reflect.Ptr {
		elm = reflect.ValueOf(elm).Elem().Interface()
	}
	ekind := reflect.TypeOf(elm).Kind()
	if ekind == reflect.Struct {
		q.model = elm
		return
	}
	/*if ekind == reflect.Slice {
		m:=reflect.New(reflect.TypeOf(elm).Elem()).Elem().Interface()
		q.model =m
	}*/
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
