package main

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type Test struct {
	//[ 0] Id                                             INT                  null: false  primary: true   isArray: false  auto: false  col: INT             len: -1      default: []
	Ix int32 `gorm:"primary_key;column:Id;type:INT;" db:"Id"`
	//[ 1] Name                                           NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	Name *string `gorm:"primary_key;column:Name;type:NVARCHAR;size:100;" db:"Name"`
	//[ 2] Address                                        NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	Address *string `gorm:"column:Address;type:NVARCHAR;size:100;" json:"Address" db:"Address"`
	//[ 3] Tel                                            DECIMAL              null: true   primary: false  isArray: false  auto: false  col: DECIMAL         len: -1      default: []
	Tel *float64 `gorm:"column:Tel;type:DECIMAL;" json:"Tel" db:"Tel"`
	//[ 4] Date                                           DATE                 null: true   primary: false  isArray: false  auto: false  col: DATE            len: -1      default: []
	Date *time.Time `gorm:"column:Date;type:DATE;" json:"Date" db:"Date"`
	//[ 5] Time                                           DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	Time *time.Time `gorm:"column:Time;type:DATETIME;" json:"Time" db:"Time"`
	//[ 6] X_y                                            INT                  null: false  primary: false  isArray: false  auto: true   col: INT             len: -1      default: []
	X_y *int32 `gorm:"AUTO_INCREMENT;column:X_y;type:INT;" json:"X_y" db:"X_y"`
	//[ 7] az_Nm                                          NCHAR(20)            null: true   primary: false  isArray: false  auto: false  col: NCHAR           len: 20      default: []
	Az_Nm *string `gorm:"column:az_Nm;type:NCHAR;size:20;" json:"az_Nm" db:"az_Nm"`
	//[ 8] cdf                                            NCHAR(20)            null: true   primary: false  isArray: false  auto: false  col: NCHAR           len: 20      default: []
	Cdf *string `gorm:"column:cdf;type:NCHAR;size:20;" json:"cdf" db:"cdf"`
}

// TableName sets the insert table name for this struct type
func (t *Test) TableName() string {
	return "Test"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *Test) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *Test) Prepare() {
}
