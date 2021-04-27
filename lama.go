package lama

import (
	"database/sql"
	"github.com/HassanAbdelzaher/lama/sqlx"
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

func nQ(l *Lama) *Query {
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
	return nQ(l).OrderBy(by...)
}
func (l *Lama) GroupBy(by ...string) *Query {
	return nQ(l).GroupBy(by...)
}
func (l *Lama) Having(expr string, args ...sql.NamedArg) *Query {
	return nQ(l).Having(expr, args...)
}
func (l *Lama) Limit(limit int) *Query {
	return nQ(l).Limit(limit)
}
func (l *Lama) Offset(off int) *Query {
	return nQ(l).Offset(off)
}
func (l *Lama) Select(cols ...string) *Query {
	return nQ(l).Select(cols...)
}
func (l *Lama) Count() (*int64, error) {
	return nQ(l).Count()
}
func (l *Lama) Sum(column string) (*float64, error) {
	return nQ(l).Sum(column)
}
func (l *Lama) CountColumn(dest interface{}, expr string) error {
	return nQ(l).CountColumn(dest, expr)
}
func (l *Lama) ColumnsFromStructOrMap(str interface{}, skipUnTaged bool) *Query {
	return nQ(l).ColumnsFromStructOrMap(str, skipUnTaged)
}
func (l *Lama) Where(query interface{}, args ...sql.NamedArg) *Query {
	return nQ(l).Where(query, args...)
}
func (l *Lama) WhereIn(key string, values ...interface{}) *Query {
	return nQ(l).WhereIn(key, values...)
}
func (l *Lama) WhereOr(w ...Where) *Query {
	return nQ(l).WhereOr(w...)
}
func (l *Lama) Model(model interface{}) *Query {
	return nQ(l).Model(model)
}
func (l *Lama) Table(table string) *Query {
	return nQ(l).Table(table)
}

func (l *Lama) Find(dest interface{}) error {
	er := nQ(l).Find(dest)
	return er
}
func (l *Lama) Get(dest interface{}) error {
	er := nQ(l).Get(dest)
	return er
}
func (l *Lama) First(dest interface{}) error {
	er := nQ(l).First(dest)
	return er
}
func (l *Lama) Last(dest interface{}) error {
	er := nQ(l).Last(dest)
	return er
}
func (l *Lama) Save(entity interface{}) error {
	er := nQ(l).Save(entity)

	return er
}
func (l *Lama) Delete(entity interface{}) error {
	er := nQ(l).Delete(entity)

	return er
}
func (l *Lama) Add(entity interface{}) error {
	er := nQ(l).Add(entity)
	return er
}

func (l *Lama) Update(entity map[string]interface{}, bulkUpdate bool) error {
	er := nQ(l).Update(entity, bulkUpdate)
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
