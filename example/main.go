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
func main() {
	//testNested()
	testCase("")
	time.Sleep(2 * time.Second)
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
db:= conn
query:="select TARRIF_id from TARIFFS"
var bt []*Tariffs2
err:=db.DB.Select(&bt,query)
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
