package test

import (
	"testing"
	"time"

	lama2 "github.com/HassanAbdelzaher/lama"
	_ "github.com/HassanAbdelzaher/lama/dialects/mssql"
	_ "github.com/HassanAbdelzaher/lama/dialects/oracle"
)

type TestTable struct {
	//[ 0] Id                                             INT                  null: false  primary: true   isArray: false  auto: false  col: INT             len: -1      default: []
	ID int32 `gorm:"primary_key;column:ID;type:INT;" db:"ID"`
	//[ 1] Name                                           NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	Name *string `gorm:"column:NAME;type:NVARCHAR;size:100;" db:"NAME"`
	//[ 2] Address                                        NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	STAMP_DATE *time.Time `gorm:"column:STAMP_DATE;type:date;size:100;" json:"STAMP_DATE" db:"STAMP_DATE"`

	COUNTER *int64 `gorm:"column:COUNTER;type:INT;" db:"COUNTER"`
}

var name = "FIRST"
var now = time.Now()
var counter int64 = 987654321
var sample = TestTable{
	ID:         1,
	Name:       &name,
	STAMP_DATE: &now,
	COUNTER:    &counter,
}

// TableName sets the insert table name for this struct type
func (t *TestTable) TableName() string {
	return "TestTable"
}

var DB *lama2.Lama

func _TestSelect(t *testing.T) {
	if DB == nil {
		return
	}
	var data []TestTable
	err := DB.Model(TestTable{}).Find(&data)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if data == nil {
		t.Error("data returned is nil")
		return
	}
	if len(data) == 0 {
		t.Error("empty returned array while expected one record at least")
		return
	}
	var sm *TestTable = nil
	for i := range data {
		if data[i].ID == sample.ID {
			sm = &data[i]
		}
	}
	if sm == nil {
		t.Error("smaple not found")
		return
	}

	if sm.Name == nil {
		t.Errorf("exprected %s while found nil", *sample.Name)
	} else {
		if *sm.Name != *sample.Name {
			t.Errorf("exprected %s while found %s", *sample.Name, *sm.Name)
		}
	}

	if sm.COUNTER == nil {
		t.Errorf("exprected %d while found nil", sample.COUNTER)
		return
	}
	if *sm.COUNTER != *sample.COUNTER {
		t.Errorf("exprected %d while found %d", *sample.COUNTER, *sm.COUNTER)
	}
	//TestConnect(t)
}
func _TestGroupBy(t *testing.T) {
	if DB == nil {
		return
	}
	sm, err := DB.Model(TestTable{}).Having("count(*)>0").GroupBy("NAME").Sum("ID")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("sum is: ", *sm)
}

func _TestDelete(t *testing.T) {
	if DB == nil {
		return
	}
	count, err := DB.Model(TestTable{}).Where(TestTable{ID: sample.ID}).Count()
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = DB.Delete(TestTable{
		ID: sample.ID,
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	count, err = DB.Model(TestTable{}).Where(TestTable{ID: sample.ID}).Count()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if *count > 0 {
		t.Errorf("after delete expected count is 0 while found %v", count)
	}
}

func _TestAdd(t *testing.T) {
	if DB == nil {
		return
	}
	err := DB.Add(sample)
	if err != nil {
		t.Error(err.Error())
		return
	}
	count, err := DB.Model(TestTable{}).Where(TestTable{ID: sample.ID}).Count()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if *count != 1 {
		t.Errorf("after delete expected count is 1 while found %v", count)
	}
}

func _TestUpdate(t *testing.T) {
	if DB == nil {
		return
	}
	err := DB.Model(TestTable{}).Where(TestTable{ID: sample.ID}).Update(map[string]interface{}{
		"NAME":       "hassan",
		"STAMP_DATE": time.Now(),
		"COUNTER":    *sample.COUNTER + 1,
	}, false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	var tbl TestTable
	err = DB.Model(TestTable{}).Where(TestTable{ID: sample.ID}).First(&tbl)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if tbl.COUNTER == nil || *tbl.COUNTER != *sample.COUNTER+1 {
		t.Errorf("expected %v while found %v", *sample.COUNTER+1, tbl.COUNTER)
	}
}

func _TestSave(t *testing.T) {
	if DB == nil {
		return
	}
	var tbl TestTable
	err := DB.Model(TestTable{}).Where(TestTable{ID: sample.ID}).First(&tbl)
	if err != nil {
		t.Error(err.Error())
		return
	}
	var cnt int64 = 0
	if tbl.COUNTER != nil {
		cnt = *tbl.COUNTER
	}
	counter := cnt + 1
	tbl.COUNTER = &counter
	err = DB.Save(&tbl)
	if err != nil {
		t.Error(err.Error())
		return
	}
	var tbl2 TestTable
	err = DB.Model(TestTable{}).Where(TestTable{ID: sample.ID}).First(&tbl2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if tbl2.COUNTER == nil || *tbl2.COUNTER != counter {
		t.Errorf("found %d while expected %d", *tbl2.COUNTER, counter)
	}
}

func _TestSum(t *testing.T) {
	if DB == nil {
		return
	}
	_, err := DB.Model(TestTable{}).Sum("ID")
	if err != nil {
		t.Error(err)
	}
}
