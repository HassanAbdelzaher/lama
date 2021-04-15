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

// Tariffs : Is Table In SQL For Model
type Tariffs struct {
	TarrifID    int32     `gorm:"column:TARRIF_ID;type:INT;" db:"TARRIF_ID"`
	TariffCode  string    `gorm:"column:TARIFF_CODE;type:NVARCHAR;size:50;" db:"TARIFF_CODE"`
	EffectDate  time.Time `gorm:"primary_key;column:EFFECT_DATE;type:DATETIME;" db:"EFFECT_DATE"`
	Description string    `gorm:"column:DESCRIPTION;type:NVARCHAR;size:250;" db:"DESCRIPTION"`
}

// TableName sets the insert table name for this Tariffs
func (t *Tariffs) TableName() string {
	return "TARIFFS"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *Tariffs) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *Tariffs) Prepare() {
}
