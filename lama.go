package lama

import (
	"database/sql"
	"fmt"
	"github.com/HassanAbdelzaher/lama/sqlx"
	"strings"
	"sync"
	"time"

	//"github.com/HassanAbdelzaher/lama/sqlx"
)

type Lama struct {
	//query Query
	Tx      *sqlx.Tx
	DB      *sqlx.DB
	Debug   bool
	dialect Dialect
	sync.Mutex
}
func (l *Lama) Dialect() Dialect{
	return l.dialect
}

type LamaTx struct {
	*Lama
}

func Connect(driver string, connstr string) (*Lama, error) {
	//cnsStr := fmt.Sprintf("server=%s;database=%s;user=%s;password=%s", config.AppConfig.Server, config.AppConfig.Database, config.AppConfig.User, config.AppConfig.Passeord)
	if driver == "oracle" {
		driver = "godror"
	}
	if driver == "mssql" {
		driver = "sqlserver"
	}
	DbConn, err := sqlx.Connect(driver, connstr)
	if err != nil {
		return nil, err
	}
	DbConn.SetMaxOpenConns(5);
	DbConn.SetMaxIdleConns(1);
	DbConn.SetConnMaxLifetime(1 * time.Hour)
	dialect := newDialect(driver)
	return &Lama{
		DB:      DbConn,
		Debug:   false,
		dialect: dialect,
	}, nil
}


///create wheres
// WhereEqual accept nil and zero values
func Eq(key string,value interface{}) Where {
	if value==nil{
		return Where{Raw:fmt.Sprintf(`%s is null`,key)}
	}
	return Where{Expr:key,Op:"=",Value:value}
}

// WhereEqual accept nil and zero values
func NotEq(key string,value interface{}) Where {
	if value==nil{
		return Where{Raw:fmt.Sprintf(`%s is not null`,key)}
	}
	return Where{Expr:key,Op:"<>",Value:value}
}

func Gt(key string,value interface{}) Where {
	if value == nil {
		return Where{Fake:true}
	}
	return Where{Expr:key,Op:">",Value:value}
}

func Gte(key string,value interface{}) Where {
	if value == nil {
		return Where{Fake:true}
	}
	return Where{Expr:key,Op:">=",Value:value}
}

func Lt(key string,value interface{}) Where {
	if value == nil {
		return Where{Fake:true}
	}
	return Where{Expr:key,Op:"<",Value:value}
}

func Lte(key string,value interface{}) Where {
	if value == nil {
		return Where{Fake:true}
	}
	return Where{Expr:key,Op:"<=",Value:value}
}

func In(di Dialect,key string, values ...interface{}) Where {
	if values == nil || len(values) == 0 {
		return Where{Fake:true}
	}
	args := make([]sql.NamedArg, 0)
	ins := make([]string, 0)
	for idx := range values {
		nam :=(getArgName(fmt.Sprintf(`%s%d`,key,idx)))
		args = append(args, sql.NamedArg{Name: nam, Value: values[idx]})
		//args[nam] = v
		ins = append(ins, di.BindVarStr(nam))
	}
	//stm := " " + key + " between (" + strings.Join(ins, ",") + ")"
	stm:=fmt.Sprintf(`%s in (%s)`,key,strings.Join(ins, ","))
	return Where{Raw: stm, Args: args}
}

func Between(di Dialect,key string, value1 interface{},value2 interface{}) Where {
	if value1 == nil && value2==nil {
		return Where{Fake:true}
	}
	if value2==nil{
		return Gte(key,value1)
	}
	if value1==nil{
		return Lte(key,value2)
	}
	args := make([]sql.NamedArg, 0)
	nam1 :=(getArgName(key))
	nam2 := (getArgName(key+"_"))
	args = append(args, sql.NamedArg{Name: nam1, Value: value1})
	args = append(args, sql.NamedArg{Name: nam2, Value: value2})
	stm := fmt.Sprintf(`%s between %s and %s`,key,di.BindVarStr(nam1),di.BindVarStr(nam2))
	return Where{Raw: stm, Args: args}
}

func Like(key string, value1 interface{}) Where {
	if value1 == nil{
		return Where{Fake:true}
	}
	args := make([]sql.NamedArg, 0)
	stm := fmt.Sprintf(`%s like '%%%s%%'`,key,value1)
	return Where{Raw: stm, Args: args}
}

func EndsWith(key string, value1 string) Where {
	args := make([]sql.NamedArg, 0)
	stm := fmt.Sprintf(`%s like '%%%s'`,key,value1)
	return Where{Raw: stm, Args: args}
}

func StartsWith(key string, value1 string) Where {
	args := make([]sql.NamedArg, 0)
	stm := fmt.Sprintf(`%s like '%s%%'`,key,value1)
	return Where{Raw: stm, Args: args}
}


func newQuery(l *Lama) *Query {
	l.Lock()
	defer l.Unlock()
	query := Query{debug: l.Debug, lama: l}
	query.args = make([]sql.NamedArg, 0)
	query.values = make(map[string]interface{})
	//createing new transaction with new query
	//make connection leak
	//so transaction must be create inside the actual function
	return &query
}

func (l *Lama) OrderBy(by ...string) *Query {
	return newQuery(l).OrderBy(by...)
}
func (l *Lama) GroupBy(by ...string) *Query {
	return newQuery(l).GroupBy(by...)
}
func (l *Lama) Having(expr string, args ...sql.NamedArg) *Query {
	return newQuery(l).Having(expr, args...)
}
func (l *Lama) Limit(limit int) *Query {
	return newQuery(l).Limit(limit)
}
func (l *Lama) Offset(off int) *Query {
	return newQuery(l).Offset(off)
}
func (l *Lama) Select(cols ...string) *Query {
	return newQuery(l).Select(cols...)
}
func (l *Lama) Count() (*int64, error) {
	return newQuery(l).Count()
}
func (l *Lama) Sum(column string) (*float64, error) {
	return newQuery(l).Sum(column)
}
func (l *Lama) CountColumn(dest interface{}, expr string) error {
	return newQuery(l).CountColumn(dest, expr)
}
func (l *Lama) ColumnsFromStructOrMap(str interface{}, skipUnTaged bool) *Query {
	return newQuery(l).ColumnsFromStructOrMap(str, skipUnTaged)
}
func (l *Lama) Where(query interface{}, args ...sql.NamedArg) *Query {
	return newQuery(l).Where(query, args...)
}
func (l *Lama) WhereOr(w ...Where) *Query {
	return newQuery(l).WhereOr(w...)
}

func (l *Lama) In(key string,valus ...interface{}) *Query {
	return newQuery(l).In(key,valus)
}
func (l *Lama) Between(key string,valu1 interface{},valu2 interface{}) *Query {
	return newQuery(l).Between(key,valu1,valu2)
}
func (l *Lama) Like(key string,value string) *Query {
	return newQuery(l).Like(key,value)
}
func (l *Lama) StartsWith(key string,value string) *Query {
	return newQuery(l).StartsWith(key,value)
}
func (l *Lama) EndsWith(key string,value string) *Query {
	return newQuery(l).EndsWith(key,value)
}
func (l *Lama) Gt(key string,value interface{}) *Query {
	return newQuery(l).Gt(key,value)
}
func (l *Lama) Gte(key string,value interface{}) *Query {
	return newQuery(l).Gte(key,value)
}
func (l *Lama) Lt(key string,value interface{}) *Query {
	return newQuery(l).Lt(key,value)
}
func (l *Lama) Lte(key string,value interface{}) *Query {
	return newQuery(l).Lte(key,value)
}
func (l *Lama) Eq(key string,value interface{}) *Query {
	return newQuery(l).Eq(key,value)
}
func (l *Lama) NotEq(key string,value interface{}) *Query {
	return newQuery(l).NotEq(key,value)
}
func (l *Lama) Model(model interface{}) *Query {
	return newQuery(l).Model(model)
}
func (l *Lama) Table(table string) *Query {
	return newQuery(l).Table(table)
}
func (l *Lama) Find(dest interface{}) error {
	er := newQuery(l).Find(dest)
	return er
}
func (l *Lama) Get(dest interface{}) error {
	er := newQuery(l).Get(dest)
	return er
}
func (l *Lama) First(dest interface{}) error {
	er := newQuery(l).First(dest)
	return er
}
func (l *Lama) Last(dest interface{}) error {
	er := newQuery(l).Last(dest)
	return er
}
func (l *Lama) Save(entity interface{}) error {
	er := newQuery(l).Save(entity)
	return er
}
func (l *Lama) Upsert(entity interface{}) error {
	er := newQuery(l).Upsert(entity)
	return er
}
func (l *Lama) Delete(entity interface{}) error {
	er := newQuery(l).Delete(entity)

	return er
}
func (l *Lama) Add(entity interface{}) error {
	er := newQuery(l).Add(entity)
	return er
}
func (l *Lama) AddIfNotExists(entity interface{}) error {
	er := newQuery(l).AddIfNotExists(entity)
	return er
}

func (l *Lama) Update(entity map[string]interface{}, bulkUpdate bool) error {
	er := newQuery(l).Update(entity, bulkUpdate)
	return er
}

func (l *Lama) Close() {
	if l.DB != nil {
		l.DB.Close()
	}
}

func (l *Lama) Begin() (*LamaTx, error) {
	l.Lock()
	defer l.Unlock()
	tx, err := l.DB.Beginx()
	if err != nil {
		return &LamaTx{Lama:l}, err
	}
	return &LamaTx{Lama:&Lama{
		Debug:   l.Debug,
		DB:      l.DB,
		dialect: l.dialect,
		Tx:tx,
	}}, nil
}

func (l *LamaTx) Commit() error {
	if l.Tx != nil {
		l.Lock()
		defer l.Unlock()
		err := l.Tx.Commit()
		l.Tx = nil
		return err
	}
	return nil
}

func (l *LamaTx) Rollback() error {
	if l.Tx != nil {
		l.Lock()
		defer l.Unlock()
		err := l.Tx.Rollback()
		l.Tx = nil
		return err
	}
	return nil
}
