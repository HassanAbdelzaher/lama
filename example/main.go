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
	do("")
	time.Sleep(2 * time.Second)
}

type Car struct{
	Name string
}

type Skoda struct{
	*Car
	Addrees string
}

func testNested(){
	car :=Car{Name:"car"}
	skoda :=Skoda{Car:&car}
	skoda.Name="skoda"
	skoda.Addrees="adddd"
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

func do(id string) {
	log.Println("start " + id)
	db:= conn
	var bt []BILL_ITEM2
    err:=db.Where(BILL_ITEM2{BILL_ITEM{CUSTKEY:"100021148"}}).Find(&bt)
    if err!=nil{
    	log.Println(err)
		return
	}
	log.Println(len(bt))
    for _,ro:=range bt{
    	log.Println(ro.CUSTKEY,*ro.CYCLE_ID,*ro.WATER_AMT)
	}
	return

}
