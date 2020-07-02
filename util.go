package lama

import (
	"encoding/json"
	"reflect"
)

func GetTableName(inx interface{}) string {
	tbleName := ""
	val := reflect.ValueOf(inx)
	switch val.Kind() {
	case reflect.Struct:
		concerte := reflect.New(reflect.TypeOf(inx)).Interface()
		kind := reflect.ValueOf(concerte).Kind()
		if kind == reflect.Ptr {
			return GetTableName(concerte)
		} else {
			return ""
		}
	case reflect.Ptr:
		elm := val.Elem()
		tblFuc := reflect.ValueOf(inx).MethodByName("TableName")
		if tblFuc.IsValid() {
			in := make([]reflect.Value, tblFuc.Type().NumIn())
			tblN := tblFuc.Call(in)
			if len(tblN) == 1 {
				str, ok := tblN[0].Interface().(string)
				if ok && len(str) > 0 {
					return str
				}
			}
		}
		tbleName = elm.Type().Name()
	default:
		tbleName = ""
	}
	return tbleName
}

func StructToMap(inx interface{}, skipZeroValues bool) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	j, _ := json.Marshal(inx)
	err := json.Unmarshal(j, &m)
	if err == nil {
		for k, v := range m {
			if v == nil || reflect.ValueOf(v).IsZero() {
				if skipZeroValues {
					delete(m, k)
				}
			}
		}
	}
	return m, err
}
