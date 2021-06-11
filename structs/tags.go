package structs

import (
	"errors"
	"reflect"
	"strings"
)

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

const TagName string = "gorm" // we use same as gorm to use the previous model generator

// tagOptions contains a slice of tag options
type tagOptions []string

// Has returns true if the given option is available in tagOptions
func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}

	return false
}

// parseTag splits a struct field's tag into its name and a list of options
// which comes after a name. A tag is in the form of: "name,option1,option2".
// The name can be neglectected.
func _parseTag(tag string) (string, tagOptions) {
	// tag is one of followings:
	// ""
	// "name"
	// "name,opt"
	// "name,opt,opt2"
	// ",opt"

	res := strings.Split(tag, ",")
	return res[0], res[1:]
}



func FieldTagExtractor(entity interface{}, fieldName string) (*GormTag, error) {
	tf := reflect.TypeOf(entity)
	if tf.Kind() == reflect.Struct {
		fld, ok := tf.FieldByName(fieldName)
		if !ok {
			return nil, errors.New("field not found:" + fieldName)
		}
		tag := fld.Tag.Get(TagName)
		return TagExtractor(tag),nil
	} else {
		return nil, errors.New("invalied argument:can not convert map")
	}
}

func TagExtractor(tag string) *GormTag {
	gTag := GormTag{}
	if len(tag) > 0 {
		parts := strings.Split(tag, `;`)
		for i := range parts {
			p := strings.TrimSpace(parts[i])
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

	return &gTag
}

func PrimaryKey(entity interface{}) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	if entity == nil {
		return ret, nil
	}
	typ := reflect.TypeOf(entity)
	val:=reflect.ValueOf(entity)
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
			}else{
				if tag==nil || tag.ColumnName=="" || tag.ColumnName=="embed"{
					fType:=typ.Field(i)
					isEmbed:=false
					fval:=val.Field(i)
					if fType.Type.Kind()==reflect.Struct {
						isEmbed=true
					}
					if fType.Type.Kind()==reflect.Ptr {
						if fType.Type.Elem().Kind()==reflect.Struct  && fval.IsValid() &&!fval.IsZero() && !fval.IsNil(){
							isEmbed=true
						}
					}
					if isEmbed {
						//if ntag=="embed"|| ntag=="_embed" || ntag=="_"{
						keys,err:=PrimaryKey(fval.Interface())
						if err!=nil{
							return nil,err
						}
						if keys!=nil{
							for k:=range keys{
								_,ok:=ret[k]
								if !ok{
									ret[k]=keys[k]
								}
							}
						}
					}
				}

			}
		}

	} else if typ.Kind() == reflect.Ptr {
		strct := reflect.ValueOf(entity).Elem().Interface()
		if reflect.TypeOf(strct).Kind() == reflect.Struct {
			return PrimaryKey(strct)
		} else {
			return nil, errors.New("can not found primarykey for non struct types")
		}
	} else {
		return nil, errors.New("can not found primarykey for non struct types")
	}
	return ret, nil
}
