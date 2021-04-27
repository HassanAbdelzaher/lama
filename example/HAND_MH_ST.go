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
//using with database version less than 20
type HAND_MH_ST struct {
	//[ 0] STATION_NO                                     INT                  null: false  primary: true   isArray: false  auto: false  col: INT             len: -1      default: []
	STATION_NO *int32 `gorm:"primary_key;column:STATION_NO;type:INT;" json:"STATION_NO" db:"STATION_NO"`
	//[ 1] CUSTKEY                                        NVARCHAR(60)         null: false  primary: true   isArray: false  auto: false  col: NVARCHAR        len: 60      default: []
	CUSTKEY string `gorm:"primary_key;column:CUSTKEY;type:NVARCHAR;size:60;" json:"CUSTKEY" db:"CUSTKEY"`
	//[ 2] CYCLE_ID                                       INT                  null: false  primary: true   isArray: false  auto: false  col: INT             len: -1      default: [0]
	CYCLE_ID *int32 `gorm:"primary_key;column:CYCLE_ID;type:INT;default:0;" json:"CYCLE_ID" db:"CYCLE_ID"`
	//[ 3] READING_DEVICEID                               NVARCHAR(200)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 200     default: []
	READING_DEVICEID *string `gorm:"column:READING_DEVICEID;type:NVARCHAR;size:200;" json:"READING_DEVICEID" db:"READING_DEVICEID"`
	//[ 4] COLLECTION_DEVICEID                            NVARCHAR(200)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 200     default: []
	COLLECTION_DEVICEID *string `gorm:"column:COLLECTION_DEVICEID;type:NVARCHAR;size:200;" json:"COLLECTION_DEVICEID" db:"COLLECTION_DEVICEID"`
	//[ 5] BILNG_DATE                                     DATE                 null: false  primary: false  isArray: false  auto: false  col: DATE            len: -1      default: []
	BILNG_DATE *time.Time `gorm:"column:BILNG_DATE;type:DATE;" json:"BILNG_DATE" db:"BILNG_DATE"`
	//[ 6] BILLGROUP                                      NVARCHAR(60)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 60      default: []
	BILLGROUP *string `gorm:"column:BILLGROUP;type:NVARCHAR;size:60;" json:"BILLGROUP" db:"BILLGROUP"`
	//[ 7] BOOK_NO_C                                      NVARCHAR(60)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 60      default: []
	BOOK_NO_C *string `gorm:"column:BOOK_NO_C;type:NVARCHAR;size:60;" json:"BOOK_NO_C" db:"BOOK_NO_C"`
	//[ 8] WALK_NO_C                                      NVARCHAR(60)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 60      default: []
	WALK_NO_C *string `gorm:"column:WALK_NO_C;type:NVARCHAR;size:60;" json:"WALK_NO_C" db:"WALK_NO_C"`
	//[ 9] SEQ_NO_C                                       INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	SEQ_NO_C *int64 `gorm:"column:SEQ_NO_C;type:INT;" json:"SEQ_NO_C" db:"SEQ_NO_C"`
	//[10] BOOK_NO_R                                      NVARCHAR(60)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 60      default: []
	BOOK_NO_R *string `gorm:"column:BOOK_NO_R;type:NVARCHAR;size:60;" json:"BOOK_NO_R" db:"BOOK_NO_R"`
	//[11] WALK_NO_R                                      NVARCHAR(60)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 60      default: []
	WALK_NO_R *string `gorm:"column:WALK_NO_R;type:NVARCHAR;size:60;" json:"WALK_NO_R" db:"WALK_NO_R"`
	//[12] SEQ_NO_R                                       INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	SEQ_NO_R *int64 `gorm:"column:SEQ_NO_R;type:INT;" json:"SEQ_NO_R" db:"SEQ_NO_R"`
	//[13] tent_name                                      NVARCHAR(300)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 300     default: []
	Tent_name *string `gorm:"column:tent_name;type:NVARCHAR;size:300;" json:"tent_name" db:"tent_name"`
	//[14] c_type                                         NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	C_type *string `gorm:"column:c_type;type:NVARCHAR;size:40;" json:"c_type" db:"c_type"`
	//[15] description                                    NVARCHAR(140)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 140     default: []
	Description *string `gorm:"column:description;type:NVARCHAR;size:140;" json:"description" db:"description"`
	//[20] meter_type                                     NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	Meter_type *string `gorm:"column:meter_type;type:NVARCHAR;size:40;" json:"meter_type" db:"meter_type"`
	//[21] meter_ref                                      NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	Meter_ref *string `gorm:"column:meter_ref;type:NVARCHAR;size:40;" json:"meter_ref" db:"meter_ref"`
	//[22] ua_adress1                                     NVARCHAR(300)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 300     default: []
	Ua_adress1 *string `gorm:"column:ua_adress1;type:NVARCHAR;size:300;" json:"ua_adress1" db:"ua_adress1"`
	//[23] ua_adress2                                     NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	Ua_adress2 *string `gorm:"column:ua_adress2;type:NVARCHAR;size:100;" json:"ua_adress2" db:"ua_adress2"`
	//[24] ua_adress3                                     NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	Ua_adress3 *string `gorm:"column:ua_adress3;type:NVARCHAR;size:100;" json:"ua_adress3" db:"ua_adress3"`
	//[25] ua_adress4                                     NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	Ua_adress4 *string `gorm:"column:ua_adress4;type:NVARCHAR;size:100;" json:"ua_adress4" db:"ua_adress4"`
	//[26] prop_ref                                       NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	Prop_ref *string `gorm:"column:prop_ref;type:NVARCHAR;size:40;" json:"prop_ref" db:"prop_ref"`
	//[27] cr_date                                        DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	Cr_date *time.Time `gorm:"column:cr_date;type:DATETIME;" json:"cr_date" db:"cr_date"`
	//[28] cr_reading                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: [-999999]
	Cr_reading *float64 `gorm:"column:cr_reading;type:FLOAT;default:-999999;" json:"cr_reading" db:"cr_reading"`
	//[29] pr_read1                                       FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: [-999999]
	Pr_read1 *float64 `gorm:"column:pr_read1;type:FLOAT;default:-999999;" json:"pr_read1" db:"pr_read1"`
	//[30] consump                                        FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: [-999999]
	Consump *float64 `gorm:"column:consump;type:FLOAT;default:-999999;" json:"consump" db:"consump"`
	//[31] descrepancy                                    INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: [0]
	Descrepancy *int64 `gorm:"column:descrepancy;type:INT;default:0;" json:"descrepancy" db:"descrepancy"`
	//[32] descr_msg                                      NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	//[33] stamp_date                                     DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	//[34] stamp_user                                     NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	//[35] op_status                                      INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	Op_status *int64 `gorm:"column:op_status;type:INT;" json:"op_status" db:"op_status"`
	//[36] pay_by                                         NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	//[37] cl_blnce                                       FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	Cl_blnce *float64 `gorm:"column:cl_blnce;type:FLOAT;" json:"cl_blnce" db:"cl_blnce"`
	//[38] delivery_st                                    INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: [0]
	Delivery_st *int64 `gorm:"column:delivery_st;type:INT;default:0;" json:"delivery_st" db:"delivery_st"`
	//[39] payment_no                                     NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	Payment_no *string `gorm:"column:payment_no;type:NVARCHAR;size:100;" json:"payment_no" db:"payment_no"`
	//[42] min_consump                                    DECIMAL              null: true   primary: false  isArray: false  auto: false  col: DECIMAL         len: -1      default: []
	Min_consump *float64 `gorm:"column:min_consump;type:DECIMAL;" json:"min_consump" db:"min_consump"`
	//[43] max_consump                                    DECIMAL              null: true   primary: false  isArray: false  auto: false  col: DECIMAL         len: -1      default: []
	Max_consump *float64 `gorm:"column:max_consump;type:DECIMAL;" json:"max_consump" db:"max_consump"`
	//[44] ctypegrp_id                                    NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	Ctypegrp_id *string `gorm:"column:ctypegrp_id;type:NVARCHAR;size:40;" json:"ctypegrp_id" db:"ctypegrp_id"`
	//[45] no_units                                       INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	No_units *int64 `gorm:"column:no_units;type:INT;" json:"no_units" db:"no_units"`
	//[46] serv_aloc                                      NVARCHAR(20)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 20      default: []
	Serv_aloc *string `gorm:"column:serv_aloc;type:NVARCHAR;size:20;" json:"serv_aloc" db:"serv_aloc"`
	//[47] lat_c                                          FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	Lat_c *float64 `gorm:"column:lat_c;type:FLOAT;" json:"lat_c" db:"lat_c"`
	//[48] lng_c                                          FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	Lng_c *float64 `gorm:"column:lng_c;type:FLOAT;" json:"lng_c" db:"lng_c"`
	//[49] lat_r                                          FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	Lat_r *float64 `gorm:"column:lat_r;type:FLOAT;" json:"lat_r" db:"lat_r"`
	//[50] lng_r                                          FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	Lng_r *float64 `gorm:"column:lng_r;type:FLOAT;" json:"lng_r" db:"lng_r"`
	//[55] IS_READING_ROW                                 INT                  null: false  primary: false  isArray: false  auto: false  col: INT             len: -1      default: [0]
	IS_READING_ROW *int32 `gorm:"column:IS_READING_ROW;type:INT;default:0;" json:"IS_READING_ROW" db:"IS_READING_ROW"`
	//[56] IS_COLLECTION_ROW                              INT                  null: false  primary: false  isArray: false  auto: false  col: INT             len: -1      default: [0]
	IS_COLLECTION_ROW *int32 `gorm:"column:IS_COLLECTION_ROW;type:INT;default:0;" json:"IS_COLLECTION_ROW" db:"IS_COLLECTION_ROW"`
	//[57] NO_UNITS_NEW                                   INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	//[58] C_TYPE_NEW                                     NVARCHAR(40)         null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 40      default: []
	//[59] IS_SEWER_ALLOCATE_NEW                          BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	//[60] READING_DATE                                   DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	READING_DATE *time.Time `gorm:"column:READING_DATE;type:DATETIME;" json:"READING_DATE" db:"READING_DATE"`
	//[61] COLLECTION_DATE                                DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	COLLECTION_DATE *time.Time `gorm:"column:COLLECTION_DATE;type:DATETIME;" json:"COLLECTION_DATE" db:"COLLECTION_DATE"`

	//[73] READING_NOTE                                   INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	READING_NOTE *int64 `gorm:"column:READING_NOTE;type:INT;" json:"READING_NOTE" db:"READING_NOTE"`
	//[74] COLLECTION_NOTE                                INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	COLLECTION_NOTE *int64 `gorm:"column:COLLECTION_NOTE;type:INT;" json:"COLLECTION_NOTE" db:"COLLECTION_NOTE"`
	//[75] LOCATION_SOURCE                                INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	LOCATION_SOURCE *int64 `gorm:"column:LOCATION_SOURCE;type:INT;" json:"LOCATION_SOURCE" db:"LOCATION_SOURCE"`
	//[76] LOCATION_DATE                                  DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	LOCATION_DATE *time.Time `gorm:"column:LOCATION_DATE;type:DATETIME;" json:"LOCATION_DATE" db:"LOCATION_DATE"`
	//[77] EMPID_C                                        INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	EMPID_C *int64 `gorm:"column:EMPID_C;type:INT;" json:"EMPID_C" db:"EMPID_C"`
	//[78] EMPID_R                                        INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	EMPID_R *int64 `gorm:"column:EMPID_R;type:INT;" json:"EMPID_R" db:"EMPID_R"`
	//[79] HHUSER_C                                       NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	HHUSER_C *string `gorm:"column:HHUSER_C;type:NVARCHAR;size:100;" json:"HHUSER_C" db:"HHUSER_C"`
	//[80] HHUSER_R                                       NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	HHUSER_R *string `gorm:"column:HHUSER_R;type:NVARCHAR;size:100;" json:"HHUSER_R" db:"HHUSER_R"`
	//[81] ACCURECY_C                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	ACCURECY_C *float64 `gorm:"column:ACCURECY_C;type:FLOAT;" json:"ACCURECY_C" db:"ACCURECY_C"`
	//[82] ACCURECY_R                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	ACCURECY_R *float64 `gorm:"column:ACCURECY_R;type:FLOAT;" json:"ACCURECY_R" db:"ACCURECY_R"`
	//[83] DEPOSITID                                      BIGINT               null: true   primary: false  isArray: false  auto: false  col: BIGINT          len: -1      default: []
	//[84] OLD_KEY                                        NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	OLD_KEY *string `gorm:"column:OLD_KEY;type:NVARCHAR;size:100;" json:"OLD_KEY" db:"OLD_KEY"`
	//[91] IS_CANCELLED_C                                 BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	IS_CANCELLED_C *bool `gorm:"column:IS_CANCELLED_C;type:BIT;" json:"IS_CANCELLED_C" db:"IS_CANCELLED_C"`
	//[92] IS_CANCELLED_R                                 BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	IS_CANCELLED_R *bool `gorm:"column:IS_CANCELLED_R;type:BIT;" json:"IS_CANCELLED_R" db:"IS_CANCELLED_R"`
	//[93] CANCEL_DATE_C                                  DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	IS_LOCKED_C *bool `gorm:"column:IS_LOCKED_C;type:BIT;" json:"IS_LOCKED_C" db:"IS_LOCKED_C"`
	//[98] IS_LOCKED_R                                    BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	IS_LOCKED_R *bool `gorm:"column:IS_LOCKED_R;type:BIT;" json:"IS_LOCKED_R" db:"IS_LOCKED_R"`
	//[100] NOTE_R                                         NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	NOTE_R *string `gorm:"column:NOTE_R;type:NVARCHAR;size:100;" json:"NOTE_R" db:"NOTE_R"`
	NOTE_C *string `gorm:"column:NOTE_C;type:NVARCHAR;size:100;" json:"NOTE_C" db:"NOTE_C"`
	//[118] HHDEVICE_VERSION                               NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	HHDEVICE_VERSION *string `gorm:"column:HHDEVICE_VERSION;type:NVARCHAR;size:100;" json:"HHDEVICE_VERSION" db:"HHDEVICE_VERSION"`
	//[119] LAT_REF                                        FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	LAT_REF *float64 `gorm:"column:LAT_REF;type:FLOAT;" json:"LAT_REF" db:"LAT_REF"`
	//[120] LNG_REF                                        FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	LNG_REF *float64 `gorm:"column:LNG_REF;type:FLOAT;" json:"LNG_REF" db:"LNG_REF"`
	//[121] IS_REF                                         BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	IS_REF *bool `gorm:"column:IS_REF;type:BIT;" json:"IS_REF" db:"IS_REF"`
	//[122] DISTANCE_REF                                   FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	DISTANCE_REF *float64 `gorm:"column:DISTANCE_REF;type:FLOAT;" json:"DISTANCE_REF" db:"DISTANCE_REF"`
	//[123] OPERATION_FLAGE                                INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	OPERATION_FLAGE *int64 `gorm:"column:OPERATION_FLAGE;type:INT;" json:"OPERATION_FLAGE" db:"OPERATION_FLAGE"`
	//[124] READING_DISTANCE                               FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	READING_DISTANCE *float64 `gorm:"column:READING_DISTANCE;type:FLOAT;" json:"READING_DISTANCE" db:"READING_DISTANCE"`

	OP_BLNCE *float64 `gorm:"column:OP_BLNCE;type:FLOAT;" json:"OP_BLNCE" db:"OP_BLNCE"`
	//[129] WATER_AMT                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
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
	//[147] INSTALMENT                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	INSTALMENT *float64 `gorm:"column:INSTALMENT;type:FLOAT;" json:"INSTALMENT" db:"INSTALMENT"`
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
	//[156] CUR_CHARGES                                    FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	CUR_CHARGES *float64 `gorm:"column:CUR_CHARGES;type:FLOAT;" json:"CUR_CHARGES" db:"CUR_CHARGES"`
	//[157] CUR_PAYMNTS                                    FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	CUR_PAYMNTS *float64 `gorm:"column:CUR_PAYMNTS;type:FLOAT;" json:"CUR_PAYMNTS" db:"CUR_PAYMNTS"`
	//[158] BILL_READY                                     BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	BILL_READY *bool `gorm:"column:BILL_READY;type:BIT;" json:"BILL_READY" db:"BILL_READY"`
	//[159] PR_READ2                                       FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	//[160] PR_CONS                                        FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	PR_CONS *float64 `gorm:"column:PR_CONS;type:FLOAT;" json:"PR_CONS" db:"PR_CONS"`
	//[161] CALC_TYPE                                      NCHAR(60)            null: true   primary: false  isArray: false  auto: false  col: NCHAR           len: 60      default: []
	CALC_TYPE *string `gorm:"column:CALC_TYPE;type:NCHAR;size:60;" json:"CALC_TYPE" db:"CALC_TYPE"`
	//[162] PR_READ_TYPE_TYPE                              NCHAR(60)            null: true   primary: false  isArray: false  auto: false  col: NCHAR           len: 60      default: []
	//[163] PR_READ_TYPE                                   NCHAR(60)            null: true   primary: false  isArray: false  auto: false  col: NCHAR           len: 60      default: []
	PR_READ_TYPE *string `gorm:"column:PR_READ_TYPE;type:NCHAR;size:60;" json:"PR_READ_TYPE" db:"PR_READ_TYPE"`
	//[164] GARD                                           BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	GARD *bool `gorm:"column:GARD;type:BIT;" json:"GARD" db:"GARD"`
	//[165] PAY_PRINT_COUNT                                INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	PAY_PRINT_COUNT *int64 `gorm:"column:PAY_PRINT_COUNT;type:INT;" json:"PAY_PRINT_COUNT" db:"PAY_PRINT_COUNT"`
	//[166] PRV_PRINT_COUNT                                INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	PRV_PRINT_COUNT *int64 `gorm:"column:PRV_PRINT_COUNT;type:INT;" json:"PRV_PRINT_COUNT" db:"PRV_PRINT_COUNT"`
	//[168] IS_HAFZA_PRINTED                               BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	IS_HAFZA_PRINTED *bool `gorm:"column:IS_HAFZA_PRINTED;type:BIT;" json:"IS_HAFZA_PRINTED" db:"IS_HAFZA_PRINTED"`
	//[169] HAFZA_PRINT_DATE                               DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	//[170] INSTALMENT_RATIO                               NCHAR(100)           null: true   primary: false  isArray: false  auto: false  col: NCHAR           len: 100     default: []
	INSTALMENT_RATIO *string `gorm:"column:INSTALMENT_RATIO;type:NCHAR;size:100;" json:"INSTALMENT_RATIO" db:"INSTALMENT_RATIO"`
	//[172] READ_TYPE                                      INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	READ_TYPE *int64 `gorm:"column:READ_TYPE;type:INT;" json:"READ_TYPE" db:"READ_TYPE"`
	//[173] S_METER_TYPE                                   VARCHAR(30)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 30      default: []
	S_METER_TYPE *string `gorm:"column:S_METER_TYPE;type:VARCHAR;size:30;" json:"S_METER_TYPE" db:"S_METER_TYPE"`
	//[178] S_CR_READING                                   FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	S_CR_READING *float64 `gorm:"column:S_CR_READING;type:FLOAT;" json:"S_CR_READING" db:"S_CR_READING"`
	//[179] S_PR_READING                                   FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	S_PR_READING *float64 `gorm:"column:S_PR_READING;type:FLOAT;" json:"S_PR_READING" db:"S_PR_READING"`
	//[180] S_CONSUMP                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	S_CONSUMP *float64 `gorm:"column:S_CONSUMP;type:FLOAT;" json:"S_CONSUMP" db:"S_CONSUMP"`
	//[184] NO_DIALS                                       INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	NO_DIALS *int64 `gorm:"column:NO_DIALS;type:INT;" json:"NO_DIALS" db:"NO_DIALS"`
	//[185] CLOCK_OVER                                     INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	CLOCK_OVER *int64 `gorm:"column:CLOCK_OVER;type:INT;" json:"CLOCK_OVER" db:"CLOCK_OVER"`
	//[188] SUPER_LOCK_R                                   BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	SUPER_LOCK_R *bool `gorm:"column:SUPER_LOCK_R;type:BIT;" json:"SUPER_LOCK_R" db:"SUPER_LOCK_R"`

	STATM_NO *int64 `gorm:"column:STATM_NO;type:INT;" json:"STATM_NO" db:"STATM_NO"`
	IS_ACC_CONSUMP *bool `gorm:"column:IS_ACC_CONSUMP;type:BIT;" json:"IS_ACC_CONSUMP" db:"IS_ACC_CONSUMP"`
	ACC_DATE *time.Time `gorm:"column:ACC_DATE;type:DATE;" json:"ACC_DATE" db:"ACC_DATE"`
	//[215] IS_MULTI_CTYPES                                BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	IS_MULTI_CTYPES *bool `gorm:"column:IS_MULTI_CTYPES;type:BIT;" json:"IS_MULTI_CTYPES" db:"IS_MULTI_CTYPES"`
	//[216] OTHER_CUR_AMT                                  FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	OTHER_CUR_AMT *float64 `gorm:"column:OTHER_CUR_AMT;type:FLOAT;" json:"OTHER_CUR_AMT" db:"OTHER_CUR_AMT"`
	//[217] CLEAN_AMT                                      FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	CLEAN_AMT *float64 `gorm:"column:CLEAN_AMT;type:FLOAT;" json:"CLEAN_AMT" db:"CLEAN_AMT"`
	//[218] COLLECTION_AMT                                 FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	COLLECTION_AMT *float64 `gorm:"column:COLLECTION_AMT;type:FLOAT;" json:"COLLECTION_AMT" db:"COLLECTION_AMT"`
	//[219] PARTIAL_CUR_AMT                                FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	PARTIAL_CUR_AMT *float64 `gorm:"column:PARTIAL_CUR_AMT;type:FLOAT;" json:"PARTIAL_CUR_AMT" db:"PARTIAL_CUR_AMT"`
	//[220] AMOUNT_COLLECTED                               FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	AMOUNT_COLLECTED *float64 `gorm:"column:AMOUNT_COLLECTED;type:FLOAT;" json:"AMOUNT_COLLECTED" db:"AMOUNT_COLLECTED"`
	//[221] BILL_AMOUNT                                    FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	BILL_AMOUNT *float64 `gorm:"column:BILL_AMOUNT;type:FLOAT;" json:"BILL_AMOUNT" db:"BILL_AMOUNT"`
	//[222] DUE_AMOUNT                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	DUE_AMOUNT *float64 `gorm:"column:DUE_AMOUNT;type:FLOAT;" json:"DUE_AMOUNT" db:"DUE_AMOUNT"`
	//[226] PRINT_READY                                    BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	PRINT_READY *bool `gorm:"column:PRINT_READY;type:BIT;" json:"PRINT_READY" db:"PRINT_READY"`
	//[227] BILL_ADJ_AMOUNT                                FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	BILL_ADJ_AMOUNT *float64 `gorm:"column:BILL_ADJ_AMOUNT;type:FLOAT;" json:"BILL_ADJ_AMOUNT" db:"BILL_ADJ_AMOUNT"`
	//[229] INSTALMENT_DATE                                DATETIME             null: true   primary: false  isArray: false  auto: false  col: DATETIME        len: -1      default: []
	INSTALMENT_DATE *time.Time `gorm:"column:INSTALMENT_DATE;type:DATETIME;" json:"INSTALMENT_DATE" db:"INSTALMENT_DATE"`
	//[230] STOP_ISSUE                                     BIT                  null: true   primary: false  isArray: false  auto: false  col: BIT             len: -1      default: []
	STOP_ISSUE *bool `gorm:"column:STOP_ISSUE;type:BIT;" json:"STOP_ISSUE" db:"STOP_ISSUE"`
	//[231] MIN_CHARGE                                     FLOAT                null: true   primary: false  isArray: false  auto: false  col: FLOAT           len: -1      default: []
	MIN_CHARGE *float64 `gorm:"column:MIN_CHARGE;type:FLOAT;" json:"MIN_CHARGE" db:"MIN_CHARGE"`
	//[237] READING_CYCLEID                                INT                  null: true   primary: false  isArray: false  auto: false  col: INT             len: -1      default: []
	READING_CYCLEID *int64 `gorm:"column:READING_CYCLEID;type:INT;" json:"READING_CYCLEID" db:"READING_CYCLEID"`
	//[238] READING_BILNGDATE                              DATE                 null: true   primary: false  isArray: false  auto: false  col: DATE            len: -1      default: []
	READING_BILNGDATE *time.Time `gorm:"column:READING_BILNGDATE;type:DATE;" json:"READING_BILNGDATE" db:"READING_BILNGDATE"`
	//[239] READ_BY                                        NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	READ_BY *string `gorm:"column:READ_BY;type:NVARCHAR;size:100;" json:"READ_BY" db:"READ_BY"`
	//[240] READ_METHOD                                    NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	READ_METHOD *string `gorm:"column:READ_METHOD;type:NVARCHAR;size:100;" json:"READ_METHOD" db:"READ_METHOD"`
	RECEIPT_TYPE *int64 `gorm:"column:RECEIPT_TYPE;type:INT;" json:"RECEIPT_TYPE" db:"RECEIPT_TYPE"`

}

// TableName sets the insert table name for this struct type
func (h *HAND_MH_ST) TableName() string {
	return "HAND_MH_ST"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (h *HAND_MH_ST) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (h *HAND_MH_ST) Prepare() {
}

