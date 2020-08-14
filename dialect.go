package lama

import (
	"database/sql"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// IsByteArrayOrSlice returns true of the reflected value is an array or slice
func IsByteArrayOrSlice(value reflect.Value) bool {
	return (value.Kind() == reflect.Array || value.Kind() == reflect.Slice) && value.Type().Elem() == reflect.TypeOf(uint8(0))
}

// Dialect interface contains behaviors that differ across SQL database
type Dialect interface {
	HaveLog() bool
	// GetName get dialect's name
	GetName() string
	// BindVar return the placeholder for actual values in SQL statements, in many dbs it is "?", Postgres using $1
	BindVar(i int) string

	BindVarStr(i string) string
	// Quote quotes field name to avoid SQL parsing exceptions by using a reserved word as a field name
	Quote(key string) string
	// DataTypeOf return data's sql type
	DataTypeOf(field *StructField) string
	// HasIndex check has index or not
	HasIndex(tableName string, indexName string,db *sql.DB) bool
	// HasForeignKey check has foreign key or not
	HasForeignKey(tableName string, foreignKeyName string,db *sql.DB) bool
	// RemoveIndex remove index
	RemoveIndex(tableName string, indexName string,db *sql.DB) error
	// HasTable check has table or not
	HasTable(tableName string,db *sql.DB) bool
	// HasColumn check has column or not
	HasColumn(tableName string, columnName string,db *sql.DB) bool
	// ModifyColumn modify column's type
	ModifyColumn(tableName string, columnName string, typ string,db *sql.DB) error

	// LimitAndOffsetSQL return generated SQL with Limit and Offset, as mssql has special case
	LimitAndOffsetSQL(statment string,limit, offset *int) string
	// SelectFromDummyTable return select values, for most dbs, `SELECT values` just works, mysql needs `SELECT value FROM DUAL`
	SelectFromDummyTable() string
	// LastInsertIDOutputInterstitial most dbs support LastInsertId, but mssql needs to use `OUTPUT`
	LastInsertIDOutputInterstitial(tableName, columnName string, columns []string) string
	// LastInsertIdReturningSuffix most dbs support LastInsertId, but postgres needs to use `RETURNING`
	LastInsertIDReturningSuffix(tableName, columnName string) string
	// DefaultValueStr
	DefaultValueStr() string

	// BuildKeyName returns a valid key name (foreign key, index key) for the given table, field and reference
	//BuildKeyName(kind, tableName string, fields ...string) string

	// NormalizeIndexAndColumn returns valid index name and column name depending on each dialect
	NormalizeIndexAndColumn(indexName, columnName string) (string, string)

	// CurrentDatabase return current database name
	CurrentDatabase(db *sql.DB) string
}

var dialectsMap = map[string]Dialect{}

func newDialect(name string) Dialect {
	if value, ok := dialectsMap[name]; ok {
		dialect := reflect.New(reflect.TypeOf(value).Elem()).Interface().(Dialect)
		return dialect
	}
	log.Println("using default dialect")
	return &defaultDialect{
	}
}

// RegisterDialect register new dialect
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect gets the dialect for the specified dialect name
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}

// ParseFieldStructForDialect get field's sql data type
var ParseFieldStructForDialect = func(field *StructField, dialect Dialect) (fieldValue reflect.Value, sqlType string, size int, additionalType string) {
	// Get redirected field type
	var (
		reflectType = field.Struct.Type
		dataType, _ = field.TagSettingsGet("TYPE")
	)

	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	// Get redirected field value
	fieldValue = reflect.Indirect(reflect.New(reflectType))

	if gormDataType, ok := fieldValue.Interface().(interface {
		GormDataType(Dialect) string
	}); ok {
		dataType = gormDataType.GormDataType(dialect)
	}

	// Get scanner's real value
	if dataType == "" {
		var getScannerValue func(reflect.Value)
		getScannerValue = func(value reflect.Value) {
			fieldValue = value
			if _, isScanner := reflect.New(fieldValue.Type()).Interface().(sql.Scanner); isScanner && fieldValue.Kind() == reflect.Struct {
				getScannerValue(fieldValue.Field(0))
			}
		}
		getScannerValue(fieldValue)
	}

	// Default Size
	if num, ok := field.TagSettingsGet("SIZE"); ok {
		size, _ = strconv.Atoi(num)
	} else {
		size = 255
	}

	// Default type from tag setting
	notNull, _ := field.TagSettingsGet("NOT NULL")
	unique, _ := field.TagSettingsGet("UNIQUE")
	additionalType = notNull + " " + unique
	if value, ok := field.TagSettingsGet("DEFAULT"); ok {
		additionalType = additionalType + " DEFAULT " + value
	}

	if value, ok := field.TagSettingsGet("COMMENT"); ok {
		additionalType = additionalType + " COMMENT " + value
	}

	return fieldValue, dataType, size, strings.TrimSpace(additionalType)
}

func currentDatabaseAndTable(dialect Dialect, tableName string,db *sql.DB) (string, string) {
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		return splitStrings[0], splitStrings[1]
	}
	return dialect.CurrentDatabase(db), tableName
}
