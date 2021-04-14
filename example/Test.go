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

//using with database version less >= 20
type BILL_ITEM struct {
	//[ 0] STATION_NO                                     INT                  null: false  primary: true   isArray: false  auto: false  col: INT             len: -1      default: []
	STATION_NO *int32 `gorm:"primary_key;column:STATION_NO;type:INT;" json:"STATION_NO" db:"STATION_NO"`
	//[ 1] CUSTKEY                                        NVARCHAR(60)         null: false  primary: true   isArray: false  auto: false  col: NVARCHAR        len: 60      default: []
	CUSTKEY string `gorm:"primary_key;column:CUSTKEY;type:NVARCHAR;size:60;" json:"CUSTKEY" db:"CUSTKEY"`
	//[ 2] CYCLE_ID                                       INT                  null: false  primary: true   isArray: false  auto: false  col: INT             len: -1      default: [0]
	CYCLE_ID *int32 `gorm:"primary_key;column:CYCLE_ID;type:INT;default:0;" json:"CYCLE_ID" db:"CYCLE_ID"`

	WATER_AMT *float64 `gorm:"column:WATER_AMT;type:FLOAT;" json:"WATER_AMT" db:"WATER_AMT"`
	//[130] SEWER_AMT                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	SEWER_AMT *float64 `gorm:"column:SEWER_AMT;type:FLOAT;" json:"SEWER_AMT" db:"SEWER_AMT"`
	//[131] BASIC_AMT                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	BASIC_AMT *float64 `gorm:"column:BASIC_AMT;type:FLOAT;" json:"BASIC_AMT" db:"BASIC_AMT"`
	//[132] TAX_AMT                                        FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	TAX_AMT *float64 `gorm:"column:TAX_AMT;type:FLOAT;" json:"TAX_AMT" db:"TAX_AMT"`
	//[133] INSTALLS_AMT                                   FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	INSTALLS_AMT *float64 `gorm:"column:INSTALLS_AMT;type:FLOAT;" json:"INSTALLS_AMT" db:"INSTALLS_AMT"`
	//[134] DBT_AMT                                        FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	DBT_AMT *float64 `gorm:"column:DBT_AMT;type:FLOAT;" json:"DBT_AMT" db:"DBT_AMT"`
	//[135] CRDT_AMT                                       FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	CRDT_AMT *float64 `gorm:"column:CRDT_AMT;type:FLOAT;" json:"CRDT_AMT" db:"CRDT_AMT"`
	//[136] AGREEM_AMT                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	AGREEM_AMT *float64 `gorm:"column:AGREEM_AMT;type:FLOAT;" json:"AGREEM_AMT" db:"AGREEM_AMT"`
	//[137] OTHER_AMT                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	OTHER_AMT *float64 `gorm:"column:OTHER_AMT;type:FLOAT;" json:"OTHER_AMT" db:"OTHER_AMT"`
	//[138] OTHER_AMT1                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	OTHER_AMT1 *float64 `gorm:"column:OTHER_AMT1;type:FLOAT;" json:"OTHER_AMT1" db:"OTHER_AMT1"`
	//[139] OTHER_AMT2                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	OTHER_AMT2 *float64 `gorm:"column:OTHER_AMT2;type:FLOAT;" json:"OTHER_AMT2" db:"OTHER_AMT2"`
	//[140] OTHER_AMT3                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	OTHER_AMT3 *float64 `gorm:"column:OTHER_AMT3;type:FLOAT;" json:"OTHER_AMT3" db:"OTHER_AMT3"`
	//[141] OTHER_AMT4                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	OTHER_AMT4 *float64 `gorm:"column:OTHER_AMT4;type:FLOAT;" json:"OTHER_AMT4" db:"OTHER_AMT4"`
	//[142] OTHER_AMT5                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	OTHER_AMT5 *float64 `gorm:"column:OTHER_AMT5;type:FLOAT;" json:"OTHER_AMT5" db:"OTHER_AMT5"`
	//[143] TAKAFUL_AMT                                    FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	TAKAFUL_AMT *float64 `gorm:"column:TAKAFUL_AMT;type:FLOAT;" json:"TAKAFUL_AMT" db:"TAKAFUL_AMT"`
	//[144] TANZEEM_AMT                                    FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	TANZEEM_AMT *float64 `gorm:"column:TANZEEM_AMT;type:FLOAT;" json:"TANZEEM_AMT" db:"TANZEEM_AMT"`
	//[145] METER_INSTALLS_AMT                             FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	METER_INSTALLS_AMT *float64 `gorm:"column:METER_INSTALLS_AMT;type:FLOAT;" json:"METER_INSTALLS_AMT" db:"METER_INSTALLS_AMT"`
	//[146] CONN_INSTALLS_AMT                              FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	CONN_INSTALLS_AMT *float64 `gorm:"column:CONN_INSTALLS_AMT;type:FLOAT;" json:"CONN_INSTALLS_AMT" db:"CONN_INSTALLS_AMT"`
	//[148] ROUND_AMT                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	ROUND_AMT *float64 `gorm:"column:ROUND_AMT;type:FLOAT;" json:"ROUND_AMT" db:"ROUND_AMT"`
	//[149] METER_AMT                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	METER_AMT *float64 `gorm:"column:METER_AMT;type:FLOAT;" json:"METER_AMT" db:"METER_AMT"`
	//[150] CONN_AMT                                       FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	CONN_AMT *float64 `gorm:"column:CONN_AMT;type:FLOAT;" json:"CONN_AMT" db:"CONN_AMT"`
	//[151] METER_MAN_AMT                                  FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	METER_MAN_AMT *float64 `gorm:"column:METER_MAN_AMT;type:FLOAT;" json:"METER_MAN_AMT" db:"METER_MAN_AMT"`
	//[152] COMPUTER_AMT                                   FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	COMPUTER_AMT *float64 `gorm:"column:COMPUTER_AMT;type:FLOAT;" json:"COMPUTER_AMT" db:"COMPUTER_AMT"`
	//[153] CONTRACT_AMT                                   FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	CONTRACT_AMT *float64 `gorm:"column:CONTRACT_AMT;type:FLOAT;" json:"CONTRACT_AMT" db:"CONTRACT_AMT"`
	//[154] GOV_AMT                                        FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	GOV_AMT *float64 `gorm:"column:GOV_AMT;type:FLOAT;" json:"GOV_AMT" db:"GOV_AMT"`
	//[155] UNI_AMT                                        FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	UNI_AMT *float64 `gorm:"column:UNI_AMT;type:FLOAT;" json:"UNI_AMT" db:"UNI_AMT"`
	//[216] OTHER_CUR_AMT                                  FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	OTHER_CUR_AMT *float64 `gorm:"column:OTHER_CUR_AMT;type:FLOAT;" json:"OTHER_CUR_AMT" db:"OTHER_CUR_AMT"`
	//[217] CLEAN_AMT                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	CLEAN_AMT *float64 `gorm:"column:CLEAN_AMT;type:FLOAT;" json:"CLEAN_AMT" db:"CLEAN_AMT"`
	//[218] COLLECTION_AMT                                 FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	COLLECTION_AMT *float64 `gorm:"column:COLLECTION_AMT;type:FLOAT;" json:"COLLECTION_AMT" db:"COLLECTION_AMT"`
	//[219] PARTIAL_CUR_AMT                                FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	PARTIAL_CUR_AMT *float64 `gorm:"column:PARTIAL_CUR_AMT;type:FLOAT;" json:"PARTIAL_CUR_AMT" db:"PARTIAL_CUR_AMT"`
}

// TableName sets the insert table name for this struct type
func (h *BILL_ITEM) TableName() string {
	return "BILL_ITEMS"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (h *BILL_ITEM) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (h *BILL_ITEM) Prepare() {
}

type TransCode string

var WATER_AMT TransCode = "WATER_AMT"
var SEWER_AMT TransCode = "SEWER_AMT"
var BASIC_AMT TransCode = "BASIC_AMT"
var TAX_AMT TransCode = "TAX_AMT"
var INSTALLS_AMT TransCode = "INSTALLS_AMT"
var DBT_AMT TransCode = "DBT_AMT"
var CRDT_AMT TransCode = "CRDT_AMT"
var AGREEM_AMT TransCode = "AGREEM_AMT"
var OTHER_AMT TransCode = "OTHER_AMT"
var OTHER_AMT1 TransCode = "OTHER_AMT1"
var OTHER_AMT2 TransCode = "OTHER_AMT2"
var OTHER_AMT3 TransCode = "OTHER_AMT3"
var OTHER_AMT4 TransCode = "OTHER_AMT4"
var OTHER_AMT5 TransCode = "OTHER_AMT5"
var TAKAFUL_AMT TransCode = "TAKAFUL_AMT"
var TANZEEM_AMT TransCode = "TANZEEM_AMT"
var METER_INSTALLS_AMT TransCode = "METER_INSTALLS_AMT"
var CONN_INSTALLS_AMT TransCode = "CONN_INSTALLS_AMT"
var ROUND_AMT TransCode = "ROUND_AMT"
var METER_AMT TransCode = "METER_AMT"
var CONN_AMT TransCode = "CONN_AMT"
var METER_MAN_AMT TransCode = "METER_MAN_AMT"
var COMPUTER_AMT TransCode = "COMPUTER_AMT"
var CONTRACT_AMT TransCode = "CONTRACT_AMT"
var GOV_AMT TransCode = "GOV_AMT"
var UNI_AMT TransCode = "UNI_AMT"
var OTHER_CUR_AMT TransCode = "OTHER_CUR_AMT"
var CLEAN_AMT TransCode = "CLEAN_AMT"
var COLLECTION_AMT TransCode = "COLLECTION_AMT"
var PARTIAL_CUR_AMT TransCode = "PARTIAL_CUR_AMT"

type TransCodes map[TransCode]string

var FinancialTransCodes = map[TransCode]string{
	WATER_AMT:          "مياة",
	SEWER_AMT:          "صرف صحي",
	BASIC_AMT:          "مقابل استدامة خدمة",
	TAX_AMT:            "الضريبة",
	INSTALLS_AMT:       "قسط",
	DBT_AMT:            "اشعار مدين",
	CRDT_AMT:           "اشعار دائن",
	AGREEM_AMT:         "اتفاقية",
	OTHER_AMT:          "اخرى",
	OTHER_AMT1:         "اخرى1",
	OTHER_AMT2:         "اخرى2",
	OTHER_AMT3:         "اخرى3",
	OTHER_AMT4:         "اخرى4",
	OTHER_AMT5:         "اخرى5",
	TAKAFUL_AMT:        "تكافل",
	TANZEEM_AMT:        "تنظيمي",
	METER_INSTALLS_AMT: "قسط عداد",
	CONN_INSTALLS_AMT:  "قسط توصيلة",
	ROUND_AMT:          "فرق تقريب",
	METER_AMT:          "عداد",
	CONN_AMT:           "توصيلة",
	METER_MAN_AMT:      "صيانة عداد",
	COMPUTER_AMT:       "كمبوتر",
	CONTRACT_AMT:       "دمغة عقد",
	GOV_AMT:            "حكومي",
	UNI_AMT:            "جامعة",
	OTHER_CUR_AMT:      "اخرى حالية",
	CLEAN_AMT:          "نظافة",
	COLLECTION_AMT:     "تحصيل",
	PARTIAL_CUR_AMT:    "تحصيل جزئي",
}

func (h *BILL_ITEM) GetItems() map[TransCode]string {
	return FinancialTransCodes
}

func (h *BILL_ITEM) GetDescription(code TransCode) string {
	var des string = string(code)
	des2, ok := FinancialTransCodes[code]
	if ok {
		des = des2
	}
	return des
}
