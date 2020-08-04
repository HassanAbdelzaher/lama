package lama

import (
	"testing"
	_ "github.com/denisenkom/go-mssqldb"
	"time"
   "math/rand"
)
type TestTable struct {
	//[ 0] Id                                             INT                  null: false  primary: true   isArray: false  auto: false  col: INT             len: -1      default: []
	ID int32 `gorm:"primary_key;column:ID;type:INT;" db:"ID"`
	//[ 1] Name                                           NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	Name *string `gorm:"column:NAME;type:NVARCHAR;size:100;" db:"NAME"`
	//[ 2] Address                                        NVARCHAR(100)        null: true   primary: false  isArray: false  auto: false  col: NVARCHAR        len: 100     default: []
	STAMP_DATE *time.Time `gorm:"column:STAMP_DATE;type:date;size:100;" json:"STAMP_DATE" db:"STAMP_DATE"`

	COUNTER int64 `gorm:"column:COUNTER;type:INT;" db:"COUNTER"`
}

// TableName sets the insert table name for this struct type
func (t *TestTable) TableName() string {
	return "TestTable"
}

var lama *Lama
func TestConnect(t *testing.T) {
	var err error
	lama, err = Connect("sqlserver", "server=masgate.com;database=lama_test;user id=sa;password=hcs@mas;log=63")
	if err != nil {
		t.Error(err.Error())
	}
	lama.Debug=false
	rand.Seed(time.Now().UnixNano())
}


func TestDelete(t *testing.T){
	TestConnect(t)

	err:=lama.Delete(TestTable{
		ID:         1,
	})
	if err!=nil{
		t.Error(err.Error())
		return
	}
	var count int
	err=lama.Model(TestTable{}).Where(TestTable{ID:1}).Count(&count)
	if err!=nil{
		t.Error(err.Error())
		return
	}
	if count>0 {
		t.Errorf("after delete expected count is 0 while found %v",count)
	}
}

func TestAdd(t *testing.T){
	TestConnect(t)
	rnd:=rand.Int63()
	t.Logf("expected id:%v",rnd)
	now:=time.Now()
	name:="FIRST"
	err:=lama.Add(TestTable{
		ID:         1,
		Name:       &name,
		STAMP_DATE: &now,
		COUNTER:    rnd,
	})
	if err!=nil{
		t.Error(err.Error())
		return
	}
	var count int
	err=lama.Model(TestTable{}).Where(TestTable{ID:1}).Count(&count)
	if err!=nil{
		t.Error(err.Error())
		return
	}
	if count!=1 {
		t.Errorf("after delete expected count is 1 while found %v",count)
	}
}

func TestUpdate(t *testing.T){
	TestConnect(t)
	rnd:=rand.Int63()
	err:=lama.Model(TestTable{}).Where(TestTable{ID:1}).Update(map[string]interface{}{
		"NAME":"hassan",
		"STAMP_DATE":time.Now(),
		"COUNTER":rnd,
	},false)
	if err!=nil{
		t.Error(err.Error())
		return
	}
	var tbl TestTable
	err=lama.Model(TestTable{}).Where(TestTable{ID:1}).First(&tbl)
	if err!=nil{
		t.Error(err.Error())
		return
	}
	if tbl.COUNTER!=rnd{
		t.Errorf("expected %v while found %v",rnd,tbl.COUNTER)
	}
}

func TestSave(t *testing.T)  {
	var tbl TestTable
	err:=lama.Model(TestTable{}).Where(TestTable{ID:1}).First(&tbl)
	if err!=nil{
		t.Error(err.Error())
		return
	}
	counter:=tbl.COUNTER
	tbl.COUNTER=counter+1
	err=lama.Save(&tbl)
	if err!=nil{
		t.Error(err.Error())
		return
	}
	var tbl2 TestTable
	err=lama.Model(TestTable{}).Where(TestTable{ID:1}).First(&tbl2)
	if err!=nil{
		t.Error(err.Error())
		return
	}
	if tbl2.COUNTER!=counter+1{
		t.Errorf("found %v while expected %v",tbl2.COUNTER,counter+1)
	}
}