package oracle

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

	// Importing oracle driver package only in dialect file, otherwide not needed
	"github.com/HassanAbdelzaher/lama"
	_ "github.com/godror/godror"
)

func init() {
	lama.RegisterDialect("godror", &oracle{})
}

type oracle struct {
}

func (oracle) HaveLog() bool {
	return false
}

func (oracle) GetName() string {
	return "godror"
}


func (oracle) BindVar(i int) string {
	return "?" // ?
}

func (oracle) BindVarStr(i string) string {
	return fmt.Sprintf(":%s", i) // ?
}

func (oracle) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (s *oracle) DataTypeOf(field *lama.StructField) string {
	var dataValue, sqlType, size, additionalType = lama.ParseFieldStructForDialect(field, s)

	if sqlType == "" {
		switch dataValue.Kind() {
		case reflect.Bool:
			sqlType = "number(1,0)"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
			if s.fieldCanAutoIncrement(field) {
				//field.TagSettingsSet("AUTO_INCREMENT", "AUTO_INCREMENT")
				sqlType = "number(12,0)"
			} else {
				sqlType = "number(12,0)"
			}
		case reflect.Int64, reflect.Uint64:
			if s.fieldCanAutoIncrement(field) {
				field.TagSettingsSet("AUTO_INCREMENT", "AUTO_INCREMENT")
				sqlType = "number(19,0)"
			} else {
				sqlType = "number(19,0)"
			}
		case reflect.Float32, reflect.Float64:
			sqlType = "number"
		case reflect.String:
			if size > 0 && size < 8000 {
				sqlType = fmt.Sprintf("nvarchar2(%d)", size)
			} else {
				sqlType = "nvarchar2(max)"
			}
		case reflect.Struct:
			if _, ok := dataValue.Interface().(time.Time); ok {
				sqlType = "date"
			}
		default:
			if lama.IsByteArrayOrSlice(dataValue) {
				if size > 0 && size < 4000 {
					sqlType = fmt.Sprintf("varchar2(%d)", size)
				} else {
					sqlType = "clob"
				}
			}
		}
	}

	if sqlType == "" {
		panic(fmt.Sprintf("invalid sql type %s (%s) for oracle", dataValue.Type().Name(), dataValue.Kind().String()))
	}

	if strings.TrimSpace(additionalType) == "" {
		return sqlType
	}
	return fmt.Sprintf("%v %v", sqlType, additionalType)
}

func (s oracle) fieldCanAutoIncrement(field *lama.StructField) bool {
	if value, ok := field.TagSettingsGet("AUTO_INCREMENT"); ok {
		return value != "FALSE"
	}
	return false
	//return field.IsPrimaryKey
}

func (s oracle) HasIndex(tableName string, indexName string,db *sql.DB) bool {
	var count int
	db.QueryRow("SELECT count(*) FROM sys.indexes WHERE name=? AND object_id=OBJECT_ID(?)", indexName, tableName).Scan(&count)
	return count > 0
}

func (s oracle) RemoveIndex(tableName string, indexName string,db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf("DROP INDEX %v ON %v", indexName, s.Quote(tableName)))
	return err
}

func (s oracle) HasForeignKey(tableName string, foreignKeyName string,db *sql.DB) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName,db)
	db.QueryRow(`SELECT count(*) 
	FROM sys.foreign_keys as F inner join sys.tables as T on F.parent_object_id=T.object_id 
		inner join information_schema.tables as I on I.TABLE_NAME = T.name 
	WHERE F.name = ? 
		AND T.Name = ? AND I.TABLE_CATALOG = ?;`, foreignKeyName, tableName, currentDatabase).Scan(&count)
	return count > 0
}

func (s oracle) HasTable(tableName string,db *sql.DB) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName,db)
	db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.tables WHERE table_name = ? AND table_catalog = ?", tableName, currentDatabase).Scan(&count)
	return count > 0
}

func (s oracle) HasColumn(tableName string, columnName string,db *sql.DB) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName,db)
	db.QueryRow("SELECT count(*) FROM information_schema.columns WHERE table_catalog = ? AND table_name = ? AND column_name = ?", currentDatabase, tableName, columnName).Scan(&count)
	return count > 0
}

func (s oracle) ModifyColumn(tableName string, columnName string, typ string,db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf("ALTER TABLE %v ALTER COLUMN %v %v", tableName, columnName, typ))
	return err
}

func (s oracle) CurrentDatabase(db *sql.DB) (name string) {
	db.QueryRow("SELECT DB_NAME() AS [Current Database]").Scan(&name)
	return
}

func parseInt(value interface{}) (int64, error) {
	return strconv.ParseInt(fmt.Sprint(value), 0, 0)
}

func (oracle) LimitAndOffsetSQL(statment string,limit, offset *int) (sql string) {
	statment=" "+strings.ToUpper(statment)
	sidx:=strings.Index(statment," SELECT ")
	eidx:=strings.Index(statment," FROM ")
	runes := []rune(statment)
	cols := string(runes[sidx:eidx])
	cols=strings.ReplaceAll(cols," SELECT ","")
	if limit != nil && *limit>=0 {
		of:=0
		if (offset!=nil && *offset>0){
			of=*offset
		}
		sql =fmt.Sprintf("select %s from (select a.*,rownum as rn from (%s) a) where rn>%d and rn<=%d",cols,statment,of,(*limit+of))
	}else {
		sql=statment;
	}
	return
}

func (oracle) SelectFromDummyTable() string {
	return ""
}

func (oracle) LastInsertIDOutputInterstitial(tableName, columnName string, columns []string) string {
	if len(columns) == 0 {
		// No OUTPUT to query
		return ""
	}
	return fmt.Sprintf("OUTPUT Inserted.%v", columnName)
}

func (oracle) LastInsertIDReturningSuffix(tableName, columnName string) string {
	// https://stackoverflow.com/questions/5228780/how-to-get-last-inserted-id
	return "; SELECT SCOPE_IDENTITY()"
}

func (oracle) DefaultValueStr() string {
	return "DEFAULT VALUES"
}

// NormalizeIndexAndColumn returns argument's index name and column name without doing anything
func (oracle) NormalizeIndexAndColumn(indexName, columnName string) (string, string) {
	return indexName, columnName
}

func currentDatabaseAndTable(dialect lama.Dialect, tableName string,db *sql.DB) (string, string) {
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		return splitStrings[0], splitStrings[1]
	}
	return dialect.CurrentDatabase(db), tableName
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
