package mssql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	// Importing mssql driver package only in dialect file, otherwide not needed
	"github.com/HassanAbdelzaher/lama"
	_ "github.com/denisenkom/go-mssqldb"
)

func init() {
	lama.RegisterDialect("sqlserver", &mssql{})
}

type mssql struct {
	db *sql.DB
}

func (mssql) GetName() string {
	return "sqlserver"
}

func (s *mssql) SetDB(db *sql.DB) {
	s.db = db
}

func (mssql) BindVar(i int) string {
	return "$$$" // ?
}

func (mssql) BindVarStr(i string) string {
	return fmt.Sprintf("@%s", i) // ?
}

func (mssql) Quote(key string) string {
	return fmt.Sprintf(`[%s]`, key)
}

func (s *mssql) DataTypeOf(field *lama.StructField) string {
	var dataValue, sqlType, size, additionalType = lama.ParseFieldStructForDialect(field, s)

	if sqlType == "" {
		switch dataValue.Kind() {
		case reflect.Bool:
			sqlType = "bit"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
			if s.fieldCanAutoIncrement(field) {
				field.TagSettingsSet("AUTO_INCREMENT", "AUTO_INCREMENT")
				sqlType = "int IDENTITY(1,1)"
			} else {
				sqlType = "int"
			}
		case reflect.Int64, reflect.Uint64:
			if s.fieldCanAutoIncrement(field) {
				field.TagSettingsSet("AUTO_INCREMENT", "AUTO_INCREMENT")
				sqlType = "bigint IDENTITY(1,1)"
			} else {
				sqlType = "bigint"
			}
		case reflect.Float32, reflect.Float64:
			sqlType = "float"
		case reflect.String:
			if size > 0 && size < 8000 {
				sqlType = fmt.Sprintf("nvarchar(%d)", size)
			} else {
				sqlType = "nvarchar(max)"
			}
		case reflect.Struct:
			if _, ok := dataValue.Interface().(time.Time); ok {
				sqlType = "datetimeoffset"
			}
		default:
			if lama.IsByteArrayOrSlice(dataValue) {
				if size > 0 && size < 8000 {
					sqlType = fmt.Sprintf("varbinary(%d)", size)
				} else {
					sqlType = "varbinary(max)"
				}
			}
		}
	}

	if sqlType == "" {
		panic(fmt.Sprintf("invalid sql type %s (%s) for mssql", dataValue.Type().Name(), dataValue.Kind().String()))
	}

	if strings.TrimSpace(additionalType) == "" {
		return sqlType
	}
	return fmt.Sprintf("%v %v", sqlType, additionalType)
}

func (s mssql) fieldCanAutoIncrement(field *lama.StructField) bool {
	if value, ok := field.TagSettingsGet("AUTO_INCREMENT"); ok {
		return value != "FALSE"
	}
	return false
	//return field.IsPrimaryKey
}

func (s mssql) HasIndex(tableName string, indexName string) bool {
	var count int
	s.db.QueryRow("SELECT count(*) FROM sys.indexes WHERE name=? AND object_id=OBJECT_ID(?)", indexName, tableName).Scan(&count)
	return count > 0
}

func (s mssql) RemoveIndex(tableName string, indexName string) error {
	_, err := s.db.Exec(fmt.Sprintf("DROP INDEX %v ON %v", indexName, s.Quote(tableName)))
	return err
}

func (s mssql) HasForeignKey(tableName string, foreignKeyName string) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName)
	s.db.QueryRow(`SELECT count(*) 
	FROM sys.foreign_keys as F inner join sys.tables as T on F.parent_object_id=T.object_id 
		inner join information_schema.tables as I on I.TABLE_NAME = T.name 
	WHERE F.name = ? 
		AND T.Name = ? AND I.TABLE_CATALOG = ?;`, foreignKeyName, tableName, currentDatabase).Scan(&count)
	return count > 0
}

func (s mssql) HasTable(tableName string) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName)
	s.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.tables WHERE table_name = ? AND table_catalog = ?", tableName, currentDatabase).Scan(&count)
	return count > 0
}

func (s mssql) HasColumn(tableName string, columnName string) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName)
	s.db.QueryRow("SELECT count(*) FROM information_schema.columns WHERE table_catalog = ? AND table_name = ? AND column_name = ?", currentDatabase, tableName, columnName).Scan(&count)
	return count > 0
}

func (s mssql) ModifyColumn(tableName string, columnName string, typ string) error {
	_, err := s.db.Exec(fmt.Sprintf("ALTER TABLE %v ALTER COLUMN %v %v", tableName, columnName, typ))
	return err
}

func (s mssql) CurrentDatabase() (name string) {
	s.db.QueryRow("SELECT DB_NAME() AS [Current Database]").Scan(&name)
	return
}

func parseInt(value interface{}) (int64, error) {
	return strconv.ParseInt(fmt.Sprint(value), 0, 0)
}

func (mssql) LimitAndOffsetSQL(limit, offset interface{}) (sql string, err error) {
	if offset != nil {
		if parsedOffset, err := parseInt(offset); err != nil {
			return "", err
		} else if parsedOffset >= 0 {
			sql += fmt.Sprintf(" OFFSET %d ROWS", parsedOffset)
		}
	}
	if limit != nil {
		if parsedLimit, err := parseInt(limit); err != nil {
			return "", err
		} else if parsedLimit >= 0 {
			if sql == "" {
				// add default zero offset
				sql += " OFFSET 0 ROWS"
			}
			sql += fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", parsedLimit)
		}
	}
	return
}

func (mssql) SelectFromDummyTable() string {
	return ""
}

func (mssql) LastInsertIDOutputInterstitial(tableName, columnName string, columns []string) string {
	if len(columns) == 0 {
		// No OUTPUT to query
		return ""
	}
	return fmt.Sprintf("OUTPUT Inserted.%v", columnName)
}

func (mssql) LastInsertIDReturningSuffix(tableName, columnName string) string {
	// https://stackoverflow.com/questions/5228780/how-to-get-last-inserted-id
	return "; SELECT SCOPE_IDENTITY()"
}

func (mssql) DefaultValueStr() string {
	return "DEFAULT VALUES"
}

// NormalizeIndexAndColumn returns argument's index name and column name without doing anything
func (mssql) NormalizeIndexAndColumn(indexName, columnName string) (string, string) {
	return indexName, columnName
}

func currentDatabaseAndTable(dialect lama.Dialect, tableName string) (string, string) {
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		return splitStrings[0], splitStrings[1]
	}
	return dialect.CurrentDatabase(), tableName
}

// JSON type to support easy handling of JSON data in character table fields
// using golang json.RawMessage for deferred decoding/encoding
type JSON struct {
	json.RawMessage
}

// Value get value of JSON
func (j JSON) Value() (driver.Value, error) {
	if len(j.RawMessage) == 0 {
		return nil, nil
	}
	return j.MarshalJSON()
}

// Scan scan value into JSON
func (j *JSON) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value (strcast):", value))
	}
	bytes := []byte(str)
	return json.Unmarshal(bytes, j)
}
