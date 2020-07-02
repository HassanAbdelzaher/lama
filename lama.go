package lama

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Lama struct {
	query Query
	DB    *sqlx.DB
}

func NewLama() *Lama {
	return &Lama{
		query: Query{db: DbConn},
		DB:    DbConn,
	}
}

func nQ() *Query {
	return &Query{db: DbConn, debug: DEBUG}
}

func (l *Lama) OrderBy(by ...string) *Query {
	return nQ().OrderBy(by...)
}
func (l *Lama) Limit(limit int) *Query {
	return nQ().Limit(limit)
}
func (l *Lama) Offset(off int) *Query {
	return nQ().Offset(off)
}
func (l *Lama) Select(cols ...string) *Query {
	return nQ().Select(cols...)
}
func (l *Lama) Count(key string) *Query {
	return nQ().Count(key)
}
func (l *Lama) ColumnsFromStructOrMap(str interface{}, skipUnTaged bool) *Query {
	return nQ().ColumnsFromStructOrMap(str, skipUnTaged)
}
func (l *Lama) Where(query interface{}, args ...sql.NamedArg) *Query {
	return nQ().Where(query, args...)
}
func (l *Lama) WhereIn(key string, values ...interface{}) *Query {
	return nQ().WhereIn(key, values...)
}
func (l *Lama) WhereOr(w ...Where) *Query {
	return nQ().WhereOr(w...)
}
func (l *Lama) Model(model interface{}) *Query {
	return nQ().Model(model)
}
func (l *Lama) Table(table string) *Query {
	return nQ().Table(table)
}

func (l *Lama) Find(dest interface{}) error {
	return nQ().Find(dest)
}
func (l *Lama) Get(dest interface{}) error {
	return nQ().Get(dest)
}
func (l *Lama) Insert(entity interface{}) error {
	return nQ().Insert(entity)
}
