package lama

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type Lama struct {
	//query Query
	dbConn *sqlx.DB
	Tx     *sqlx.Tx
	Debug  bool
	sync.Mutex
}

func Connect(driver string, connstr string) (*Lama, error) {
	//cnsStr := fmt.Sprintf("server=%s;database=%s;user=%s;password=%s", config.AppConfig.Server, config.AppConfig.Database, config.AppConfig.User, config.AppConfig.Passeord)
	DbConn, err := sqlx.Connect(driver, connstr)
	if err != nil {
		return nil, err
	}
	DbConn.SetMaxOpenConns(2)
	DbConn.SetMaxIdleConns(1)
	DbConn.SetConnMaxLifetime(30 * time.Minute)
	DbConn.SetConnMaxLifetime(1 * time.Hour)
	return &Lama{
		//query: Query{db: DbConn},
		dbConn: DbConn,
		Debug:  false,
	}, nil
}

func nQ(l *Lama) *Query {
	l.Lock()
	defer l.Unlock()
	query := Query{debug: l.Debug}
	query.args = make(map[string]interface{})
	query.values = make(map[string]interface{})
	if l.Tx == nil {
		tx, err := l.dbConn.Beginx()
		if err != nil {
			log.Println(err)
			query.addError(err)
		} else {
			query.havePrivateTransaction = true
			query.tx = tx
		}
	} else {
		query.tx = l.Tx
		query.havePrivateTransaction = false
	}
	return &query
}

func (l *Lama) OrderBy(by ...string) *Query {
	return nQ(l).OrderBy(by...)
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
func (l *Lama) Count(key string) *Query {
	return nQ(l).Count(key)
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
func (l *Lama) Save(entity interface{}) error {
	er := nQ(l).Save(entity)

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
	if l.dbConn != nil {
		l.dbConn.Close()
	}
}

func (l *Lama) Begin() (*Lama, error) {
	l.Lock()
	defer l.Unlock()
	tx, err := l.dbConn.Beginx()
	if err != nil {
		return l, err
	}
	return &Lama{
		Debug:  l.Debug,
		dbConn: l.dbConn,
		Tx:     tx,
	}, nil
}

func (l *Lama) Commit() error {
	l.Lock()
	defer l.Unlock()
	if l.Tx == nil {
		return nil
	}
	err := l.Tx.Commit()
	l.Tx = nil
	return err
}

func (l *Lama) Rollback() error {
	l.Lock()
	defer l.Unlock()
	if l.Tx == nil {
		return nil
	}
	err := l.Tx.Rollback()
	l.Tx = nil
	return err
}