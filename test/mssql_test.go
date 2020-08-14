package test

import (
	lama2 "github.com/HassanAbdelzaher/lama"
	"math/rand"
	"testing"
	"time"
)

func _TestMssqlConnect(t *testing.T) {
	t.Log("testing connect")
	var err error
	lama, err = lama2.Connect("sqlserver", "server=masgate.com;database=lama_test;user id=sa;password=hcs@mas;log=63")
	if err != nil {
		t.Error(err.Error())
	}
	lama.Debug=true
	rand.Seed(time.Now().UnixNano())
}

func TestMssql(t *testing.T)  {
	_TestMssqlConnect(t)
	_TestDelete(t)
	_TestAdd(t)
	_TestUpdate(t)
	_TestSave(t)
	_TestSum(t)
}
