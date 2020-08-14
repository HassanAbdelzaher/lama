package lama


import (
	"database/sql"
	"fmt"
"reflect"
"regexp"
"strconv"
"strings"
"time"
)

var keyNameRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// DefaultForeignKeyNamer contains the default foreign key name generator method
type DefaultForeignKeyNamer struct {
}

type defaultDialect struct {
	DefaultForeignKeyNamer
}

func init() {
	RegisterDialect("common", &defaultDialect{})
}

func (defaultDialect) GetName() string {
	return "common"
}



func (s *defaultDialect) HaveLog() bool {
	return false
}

func (defaultDialect) BindVar(i int) string {
	return "$$$" // ?
}


func (defaultDialect) BindVarStr(i string) string {
	return ":"+i // ?
}

func (defaultDialect) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (s *defaultDialect) fieldCanAutoIncrement(field *StructField) bool {
	if value, ok := field.TagSettingsGet("AUTO_INCREMENT"); ok {
		return strings.ToLower(value) != "false"
	}
	return field.IsPrimaryKey
}

func (s *defaultDialect) DataTypeOf(field *StructField) string {
	var dataValue, sqlType, size, additionalType = ParseFieldStructForDialect(field, s)

	if sqlType == "" {
		switch dataValue.Kind() {
		case reflect.Bool:
			sqlType = "BOOLEAN"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
			if s.fieldCanAutoIncrement(field) {
				sqlType = "INTEGER AUTO_INCREMENT"
			} else {
				sqlType = "INTEGER"
			}
		case reflect.Int64, reflect.Uint64:
			if s.fieldCanAutoIncrement(field) {
				sqlType = "BIGINT AUTO_INCREMENT"
			} else {
				sqlType = "BIGINT"
			}
		case reflect.Float32, reflect.Float64:
			sqlType = "FLOAT"
		case reflect.String:
			if size > 0 && size < 65532 {
				sqlType = fmt.Sprintf("VARCHAR(%d)", size)
			} else {
				sqlType = "VARCHAR(65532)"
			}
		case reflect.Struct:
			if _, ok := dataValue.Interface().(time.Time); ok {
				sqlType = "TIMESTAMP"
			}
		default:
			if _, ok := dataValue.Interface().([]byte); ok {
				if size > 0 && size < 65532 {
					sqlType = fmt.Sprintf("BINARY(%d)", size)
				} else {
					sqlType = "BINARY(65532)"
				}
			}
		}
	}

	if sqlType == "" {
		panic(fmt.Sprintf("invalid sql type %s (%s) for defaultDialect", dataValue.Type().Name(), dataValue.Kind().String()))
	}

	if strings.TrimSpace(additionalType) == "" {
		return sqlType
	}
	return fmt.Sprintf("%v %v", sqlType, additionalType)
}

func (s defaultDialect) HasIndex(tableName string, indexName string,db *sql.DB) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName,db)
	db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema = ? AND table_name = ? AND index_name = ?", currentDatabase, tableName, indexName).Scan(&count)
	return count > 0
}

func (s defaultDialect) RemoveIndex(tableName string, indexName string,db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf("DROP INDEX %v", indexName))
	return err
}

func (s defaultDialect) HasForeignKey(tableName string, foreignKeyName string,db *sql.DB) bool {
	return false
}

func (s defaultDialect) HasTable(tableName string,db *sql.DB) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName,db)
	db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", currentDatabase, tableName).Scan(&count)
	return count > 0
}

func (s defaultDialect) HasColumn(tableName string, columnName string,db *sql.DB) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName,db)
	db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = ? AND table_name = ? AND column_name = ?", currentDatabase, tableName, columnName).Scan(&count)
	return count > 0
}

func (s defaultDialect) ModifyColumn(tableName string, columnName string, typ string,db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf("ALTER TABLE %v ALTER COLUMN %v TYPE %v", tableName, columnName, typ))
	return err
}

func (s defaultDialect) CurrentDatabase(db *sql.DB) (name string) {
	db.QueryRow("SELECT DATABASE()").Scan(&name)
	return
}

// LimitAndOffsetSQL return generated SQL with Limit and Offset
func (s defaultDialect) LimitAndOffsetSQL(statment string,limit, offset *int) (sql string) {
	sql=statment+" "
	if limit != nil && *limit>0 {
		sql += fmt.Sprintf(" LIMIT %d", *limit)
	}
	if offset != nil && *offset>=0 {
		sql += fmt.Sprintf(" OFFSET %d", *offset)
	}
	return
}

func (defaultDialect) SelectFromDummyTable() string {
	return ""
}

func (defaultDialect) LastInsertIDOutputInterstitial(tableName, columnName string, columns []string) string {
	return ""
}

func (defaultDialect) LastInsertIDReturningSuffix(tableName, columnName string) string {
	return ""
}

func (defaultDialect) DefaultValueStr() string {
	return "DEFAULT VALUES"
}

// BuildKeyName returns a valid key name (foreign key, index key) for the given table, field and reference
func (DefaultForeignKeyNamer) BuildKeyName(kind, tableName string, fields ...string) string {
	keyName := fmt.Sprintf("%s_%s_%s", kind, tableName, strings.Join(fields, "_"))
	keyName = keyNameRegex.ReplaceAllString(keyName, "_")
	return keyName
}

// NormalizeIndexAndColumn returns argument's index name and column name without doing anything
func (defaultDialect) NormalizeIndexAndColumn(indexName, columnName string) (string, string) {
	return indexName, columnName
}

func (defaultDialect) parseInt(value interface{}) (int64, error) {
	return strconv.ParseInt(fmt.Sprint(value), 0, 0)
}
