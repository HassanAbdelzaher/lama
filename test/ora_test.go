package test

import (
	lama2 "github.com/HassanAbdelzaher/lama"
	"math/rand"
	"testing"
	"time"
)

func _TestOracleConnect(t *testing.T) {
	var err error
	//lama, err = lama2.Connect("sqlserver", "server=masgate.com;database=lama_test;user id=sa;password=hcs@mas;log=63")
	lama, err = lama2.Connect("godror", `user="hcs_edams" password="ashaman" connectString="localhost:1521/giza"`)
	if err != nil {
		t.Error(err.Error())
		return
	}
	lama.Debug=true
	rand.Seed(time.Now().UnixNano())
}

func T1estOracle(t *testing.T)  {
	_TestOracleConnect(t)
	_TestDelete(t)
	_TestAdd(t)
	_TestUpdate(t)
	_TestSave(t)
	_TestSum(t)
	t.Log("testing connect oracle finish ")
	t.Log("==========================================================")
	t.Log("==========================================================")
}
