package lama

import (
	"errors"
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

func StructToMap(inx interface{}, skipZeroValue bool, skipComputedColumn bool, useFieldName bool) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	if inx == nil {
		return ret, nil
	}
	tf := reflect.TypeOf(inx)
	if tf.Kind() == reflect.Ptr {
		strct := reflect.ValueOf(inx).Elem().Interface()
		if reflect.TypeOf(strct).Kind() == reflect.Struct {
			return StructToMap(strct, skipZeroValue, skipComputedColumn, useFieldName)
		} else {
			return nil, errors.New("invalied argument:must be struct")
		}
	}
	if tf.Kind() == reflect.Map {
		mvals, ok := inx.(map[string]interface{})
		if ok {
			return mvals, nil
		} else {
			return nil, errors.New("invalied argument:can not convert map")
		}
	}
	if tf.Kind() == reflect.Struct {
		nf := tf.NumField()
		for i := 0; i < nf; i++ {
			filed := reflect.TypeOf(inx).Field(i)
			tag, _ := FieldTagExtractor(inx, filed.Name)
			if tag != nil {
				colName := tag.ColumnName
				cVal := reflect.ValueOf(inx).Field(i).Interface()
				isZero := false
				if skipZeroValue {
					if cVal == nil || reflect.ValueOf(cVal).IsZero() || !reflect.ValueOf(cVal).IsValid() {
						isZero = true
					}
				}
				isComputed := false
				if skipComputedColumn {
					isComputed = tag.AUTO_INCREMENT || tag.Computed
				}
				if isZero || isComputed {
					continue
				} else {
					if len(colName) > 0 {
						ret[colName] = reflect.ValueOf(inx).Field(i).Interface()
					}
				}
			} else {
				if useFieldName {
					ret[filed.Name] = reflect.ValueOf(inx).Field(i).Interface()
				}
			}

		}
	} else {
		return nil, errors.New("invalied argument:can not convert map")
	}
	return ret, nil
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
	for k, v := range _data {
		rs[k] = mapper(v)
	}
	return rs
}

func Filter(data interface{}, comparer func(interface{}) bool) []interface{} {
	if data == nil {
		return nil
	}
	_data := InterfaceSlice(data)
	rs := make([]interface{}, 0)
	for _, v := range _data {
		isMatch := comparer(v)
		if isMatch {
			rs = append(rs, v)
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
	for k, v := range src {
		dest[k] = v
	}
}
