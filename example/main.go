package main

import (
	"fmt"
	"github.com/HassanAbdelzaher/lama/structs"
	"log"
	"time"

	"github.com/HassanAbdelzaher/lama"
	_ "github.com/HassanAbdelzaher/lama/dialects/mssql"
)

var conn *lama.Lama

type S struct {
	Con *lama.Lama
}

func init() {
	var err error
	cnsStr := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;log=63", "localhost", "giza", "sa", "hcs@mas")
	conn, err = lama.Connect("sqlserver", cnsStr)
	if err != nil {
		log.Println(err)
	} else {
		conn.Debug = true
	}

}

type Server struct {
	Name    string
	ID      int32
	Enabled bool
	Mp map[string]interface{}
}


type Hub struct {
	ID      int32
	*Server
	Time time.Time
}

type HST_HAND struct{
	RECALC_ID int32 `gorm:"primary_key1;column:RECALC_ID;type:INT;default:0;" json:"RECALC_ID" db:"RECALC_ID"`
	*HAND_MH_ST
}

func main() {
	var cid int32=12
	hand:=&HAND_MH_ST{
		CUSTKEY:"121212",
		CYCLE_ID:&cid,
		STATION_NO:&cid,
	}
	hst:=&HST_HAND{
		RECALC_ID:  cid,
		HAND_MH_ST: hand,
	}
	keys,err:=structs.PrimaryKey(hst)
	log.Println(err)
	log.Println(keys)
}
type Tx time.Time
type Car struct{
	Name string
	Stamp *time.Time
}

type Skoda struct{
	*Car
	Addrees string
	Sx *Tx
	Create time.Time
}

func testNested(){
	now:=time.Now()
	car :=Car{Name:"car"}
	skoda :=Skoda{Car:&car}
	skoda.Name="skoda"
	skoda.Addrees="adddd"
	skoda.Create=now
	skoda.Stamp=&now
	var xx Tx =Tx(now)
	skoda.Sx=&xx
	m := structs.Map(&car,structs.MapOptions{
		SkipZeroValue: false,
		UseFieldName:  false,
		SkipUnTaged:   false,
		SkipComputed:  false,
		Flatten:       false,
	})
	m2 := structs.Map(&skoda,structs.MapOptions{
		SkipZeroValue: false,
		UseFieldName:  false,
		SkipUnTaged:   false,
		SkipComputed:  false,
		Flatten:       false,
	})
	log.Println(m)
	log.Println(m2)
}

func testEmbded(id string) {
	log.Println("start " + id)
	db:= conn
	var bt []*Tariffs2
    err:=db.Model(Tariffs2{}).Find(&bt)
    if err!=nil{
    	log.Println(err)
		return
	}
	for _,b:=range bt{
		log.Println(b.TarrifID,b.TariffCode)
	}
	log.Println(len(bt))
	return

}

func testCase(id string) {
	log.Println("start " + id)
	db := conn
	query := `select i.*, h.* from HAND_MH_ST2 H,HH_BCYC b,BILL_ITEMS i where  b.BILLGROUP=h.BILLGROUP and b.BOOK_NO=h.BOOK_NO_C and b.WALK_NO=h.WALK_NO_C 
	and b.CYCLE_ID=h.CYCLE_ID
	and (b.ISCYCLE_COMPLETED_C=0 or b.ISCYCLE_COMPLETED_C is null)
	and b.IS_ALLOWED_C=1
	and h.IS_COLLECTION_ROW=1
	and h.CL_BLNCE>=0
	and i.CUSTKEY=h.CUSTKEY
	and i.CYCLE_ID=h.CYCLE_ID
	and h.EMPID_C=4040111 And (CDB_HH_C=0 or CDB_HH_C is null)  order by h.billgroup,h.book_no_c,h.walk_no_c,h.cycle_id,h.SEQ_NO_C`

	var bt []*HAND_MH_ST
	err := db.DB.Select(&bt, query)
	if err != nil {
		log.Println(err)
		return
	}
	for _, b := range bt {
		if(b.Cl_blnce!=nil){
			log.Println(b.CUSTKEY,*b.Cl_blnce)

		}else {
			log.Println("missing cl_blnce")

		}
	}
	log.Println(len(bt))
	return
}
