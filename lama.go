package lama

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/HassanAbdelzaher/lama/sqlx"
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

type LamaTx struct {
	*Lama
}

func (l *Lama) Dialect() Dialect {
	return l.dialect
}

func Connect(driver string, connstr string) (*Lama, error) {
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
	DbConn.SetMaxOpenConns(5)
	DbConn.SetMaxIdleConns(1)
	DbConn.SetConnMaxLifetime(1 * time.Hour)
	dialect := newDialect(driver)
	return &Lama{
		DB:      DbConn,
		Debug:   false,
		dialect: dialect,
	}, nil
}

func newQuery(l *Lama) *Query {
	l.Lock()
	defer l.Unlock()
	query := Query{debug: l.Debug, lama: l}
	query.args = make([]sql.NamedArg, 0)
	query.values = make(map[string]interface{})
	query.context = context.Background() //default contxt
	//createing new transaction with new query
	//make connection leak
	//so transaction must be create inside the actual function
	return &query
}

func (l *Lama) Exce(stm string, ar ...interface{}) (eff int64, err error) {
	return newQuery(l).Exce(stm, ar...)
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
func (l *Lama) Max(column string) (*float64, error) {
	return newQuery(l).Max(column)
}
func (l *Lama) Min(column string) (*float64, error) {
	return newQuery(l).Min(column)
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

func (l *Lama) In(key string, valus ...interface{}) *Query {
	return newQuery(l).In(key, valus)
}
func (l *Lama) Between(key string, valu1 interface{}, valu2 interface{}) *Query {
	return newQuery(l).Between(key, valu1, valu2)
}
func (l *Lama) Like(key string, value string) *Query {
	return newQuery(l).Like(key, value)
}
func (l *Lama) StartsWith(key string, value string) *Query {
	return newQuery(l).StartsWith(key, value)
}
func (l *Lama) EndsWith(key string, value string) *Query {
	return newQuery(l).EndsWith(key, value)
}
func (l *Lama) Gt(key string, value interface{}) *Query {
	return newQuery(l).Gt(key, value)
}
func (l *Lama) Gte(key string, value interface{}) *Query {
	return newQuery(l).Gte(key, value)
}
func (l *Lama) Lt(key string, value interface{}) *Query {
	return newQuery(l).Lt(key, value)
}
func (l *Lama) Lte(key string, value interface{}) *Query {
	return newQuery(l).Lte(key, value)
}
func (l *Lama) Eq(key string, value interface{}) *Query {
	return newQuery(l).Eq(key, value)
}
func (l *Lama) NotEq(key string, value interface{}) *Query {
	return newQuery(l).NotEq(key, value)
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
		return &LamaTx{Lama: l}, err
	}
	return &LamaTx{Lama: &Lama{
		Debug:   l.Debug,
		DB:      l.DB,
		dialect: l.dialect,
		Tx:      tx,
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
