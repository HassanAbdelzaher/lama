package main

import (
	"fmt"
	"log"
	"reflect"
)

type S struct {
	Name string
}
func main(){
	var slc []*S
	fmt.Println(doubleSlice(slc))
}

func doubleSlice(s interface{}) interface{} {
	if reflect.TypeOf(s).Kind() != reflect.Slice {
		fmt.Println("The interface is not a slice.")
		return nil
	}

	v := reflect.ValueOf(s)
	newLen := v.Len()
	newCap := (v.Cap() + 1) * 2
	typ := reflect.TypeOf(s).Elem()
	n:=reflect.New(typ.Elem()).Elem().Interface()
	log.Println(n,typ)
	t := reflect.MakeSlice(reflect.SliceOf(typ), newLen, newCap)
	reflect.Copy(t, v)
	return t.Interface()
}
