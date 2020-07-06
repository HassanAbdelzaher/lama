package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/HassanAbdelzaher/lama"
	_ "github.com/denisenkom/go-mssqldb"
)

var conn *lama.Lama

func init() {
	var err error
	cnsStr := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", "localhost", "fayoum", "sa", "hcs@mas")
	conn, err = lama.Connect("mssql", cnsStr)
	if err != nil {
		log.Println(err)
	} else {
		conn.Debug = true
	}
}
func main() {
	do("sync")
	time.Sleep(2 * time.Second)
}

func do(id string) {
	log.Println("start " + id)
	db, err := conn.Begin()
	defer db.Rollback()
	if err != nil {
		log.Fatalln(err)
	}
	noe := time.Now()
	for i := 0; i < 10; i++ {
		s := "lama"
		t := Test{}
		t.Ix = int32(i) + int32(200)
		t.Name = &s
		t.Date = &noe
		t.Time = &noe
		err := db.Add(&t)
		if err != nil {
			log.Println(err)
			db.Tx.Rollback()
			break
		}
	}
	db.Commit()
	var cnt int
	//err = conn.Model(Test{}).Count("*").Get(&cnt)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Count:", cnt)
	}
	log.Println("did:" + id)
	t2 := Test{}
	err = conn.Where("id>:id", sql.NamedArg{Name: `id`, Value: 100}).Last(&t2)
	if err != nil {
		log.Println(err)
	}
	log.Println(t2)
	addr := "shrook"
	t2.Address = &addr
	err = conn.Save(&t2)
	if err != nil {
		log.Println(err)
	}
}
