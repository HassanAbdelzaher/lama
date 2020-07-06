package lama

import (
	"errors"
	"reflect"
	"strings"
)

const TagName string = "gorm" // we use same as gorm to use the previous model generator
// gorm tag are rich

type GormTag struct {
	ColumnName     string
	PRIMARY_KEY    bool
	AUTO_INCREMENT bool
	Type           string
	Size           string
	DEFAULT        string
	UNIQUE_INDEX   bool
	Computed       bool
}

func FieldTagExtractor(entity interface{}, fieldName string) (*GormTag, error) {
	tf := reflect.TypeOf(entity)
	if tf.Kind() == reflect.Struct {
		fld, ok := tf.FieldByName(fieldName)
		if !ok {
			return nil, errors.New("field not found:" + fieldName)
		}
		tag := fld.Tag.Get(TagName)
		return TagExtractor(tag)
	} else {
		return nil, errors.New("invalied argument:can not convert map")
	}
}

func TagExtractor(tag string) (*GormTag, error) {
	gTag := GormTag{}
	if len(tag) > 0 {
		parts := strings.Split(tag, `;`)
		for _, p := range parts {
			p = strings.TrimSpace(p)
			up := strings.TrimSpace(strings.ToUpper(p))
			if strings.Index(up, "PRIMARY_KEY") == 0 {
				gTag.PRIMARY_KEY = true
			}
			if strings.Index(up, "UNIQUE") == 0 {
				gTag.UNIQUE_INDEX = true
			}
			if strings.Index(up, "AUTO_INCREMENT") == 0 {
				gTag.AUTO_INCREMENT = true
			}
			if strings.Index(up, "COLUMN") == 0 {
				if len(strings.Split(p, ":")) > 1 {
					cName := strings.Split(p, ":")[1]
					gTag.ColumnName = cName
				}
			}
			if strings.Index(up, "TYPE") == 0 {
				if len(strings.Split(p, ":")) > 1 {
					ty := strings.Split(p, ":")[1]
					gTag.Type = ty
				}
			}
			if strings.Index(up, "DEFAULT") == 0 {
				if len(strings.Split(p, ":")) > 1 {
					ty := strings.Split(p, ":")[1]
					gTag.DEFAULT = ty
				}
			}
			if strings.Index(up, "SIZE") == 0 {
				if len(strings.Split(p, ":")) > 1 {
					ty := strings.Split(p, ":")[1]
					gTag.Size = ty
				}
			}
		}
	}

	return &gTag, nil
}

func primaryKey(entity interface{}) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	if entity == nil {
		return ret, nil
	}
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Struct {
		nf := typ.NumField()
		for i := 0; i < nf; i++ {
			fld := typ.Field(i)
			tag, err := FieldTagExtractor(entity, fld.Name)
			if err != nil {
				return nil, err
			}
			if tag != nil && tag.PRIMARY_KEY {
				nam := tag.ColumnName
				if len(nam) == 0 {
					nam = fld.Name
				}
				ret[nam] = reflect.ValueOf(entity).Field(i).Interface()
			}
		}

	} else if typ.Kind() == reflect.Ptr {
		strct := reflect.ValueOf(entity).Elem().Interface()
		if reflect.TypeOf(strct).Kind() == reflect.Struct {
			return primaryKey(strct)
		} else {
			return nil, errors.New("can not found primarykey for non struct types")
		}
	} else {
		return nil, errors.New("can not found primarykey for non struct types")
	}
	return ret, nil
}
