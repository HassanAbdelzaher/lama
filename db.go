package lama

import (
	"fmt"
	"sync"

	"mas.com/wsdl/config"

	//"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	//"strconv"
	"github.com/jinzhu/gorm"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var DbConn *sqlx.DB
var DEBUG = false
var mutex sync.Mutex
var Gdb *gorm.DB

func init() {
	if DbConn != nil {
		return
	}
	var err error
	cnsStr := fmt.Sprintf("server=%s;database=%s;user=%s;password=%s", config.AppConfig.Server, config.AppConfig.Database, config.AppConfig.User, config.AppConfig.Passeord)
	DbConn, err = sqlx.Connect("mssql", cnsStr)
	if err != nil || DbConn == nil {
		log.Fatalln("can not create connection pool", err)
	}
	DbConn.SetMaxOpenConns(1)
	DbConn.SetMaxIdleConns(1)
	DbConn.SetConnMaxLifetime(11 * time.Minute)
	SetDebug(true)
	/*Gdb,err=gorm.Open("mssql",cnsStr)
	if err != nil || Gdb==nil {
		log.Fatalln("can not create gorm connection pool",err)
	}
	Gdb.LogMode(true)*/

}

func CloseConnectionPool() {
	if DbConn != nil {
		DbConn.Close()
	}
}

func SetDebug(_debug bool) {
	mutex.Lock()
	defer mutex.Unlock()
	DEBUG = _debug
}
