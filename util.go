package lama

import (
	"github.com/HassanAbdelzaher/lama/structs"
	"reflect"
	"strings"
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
func StructToMap(inx interface{}, skipZeroValue bool, skipComputedColumn bool, useFieldName bool,SkipUnTaged bool,selectedZeroValues []string) (map[string]interface{}, error) {
	mpValues:=structs.New(inx,structs.MapOptions{
		SkipZeroValue: skipZeroValue,
		UseFieldName:  useFieldName,
		SkipUnTaged:   SkipUnTaged,
		SkipComputed:  skipComputedColumn,
		Flatten:       true,
		SelectedZeroValues:selectedZeroValues,
	}).Map()
	return mpValues,nil;
}
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
func Map(data interface{}, mapper func(interface{}) interface{}) []interface{} {
	if data == nil {
		return nil
	}
	_data := InterfaceSlice(data)
	rs := make([]interface{}, len(_data))
	for k := range _data {
		rs[k] = mapper(_data[k])
	}
	return rs
}
func Filter(data interface{}, comparer func(interface{}) bool) []interface{} {
	if data == nil {
		return nil
	}
	_data := InterfaceSlice(data)
	rs := make([]interface{}, 0)
	for i := range _data {
		isMatch := comparer(_data[i])
		if isMatch {
			rs = append(rs, _data[i])
		}
	}
	return rs
}
func ContainsStr(data []string, key string, matchCase bool) (bool, string) {
	if data == nil {
		return false, key
	}
	filtered := Filter(data, func(i interface{}) bool {
		ss := i.(string)
		if matchCase {
			return strings.TrimSpace(ss) == strings.TrimSpace(key)
		} else {
			return strings.TrimSpace(strings.ToUpper(ss)) == strings.TrimSpace(strings.ToUpper(key))
		}
	})
	if len(filtered) > 0 {
		return true, filtered[0].(string)
	}
	return false, key
}
func ContainsStrI(data []interface{}, key string, matchCase bool) (bool, string) {
	if data == nil {
		return false, key
	}
	filtered := Filter(data, func(i interface{}) bool {
		ss := i.(string)
		if matchCase {
			return strings.TrimSpace(ss) == strings.TrimSpace(key)
		} else {
			return strings.TrimSpace(strings.ToUpper(ss)) == strings.TrimSpace(strings.ToUpper(key))
		}
	})
	if len(filtered) > 0 {
		return true, filtered[0].(string)
	}
	return false, key
}
func appendToMap(dest map[string]interface{}, src map[string]interface{}) {
	if dest == nil {
		dest = make(map[string]interface{})
	}
	if src == nil {
		return
	}
	for k := range src {
		dest[k] = src[k]
	}
}
