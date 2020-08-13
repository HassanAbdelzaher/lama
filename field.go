package lama

import (
	"reflect"
	"sync"
)

type StructField struct {
	DBName          string
	Name            string
	Names           []string
	IsPrimaryKey    bool
	IsNormal        bool
	IsIgnored       bool
	IsScanner       bool
	HasDefaultValue bool
	Tag             reflect.StructTag
	TagSettings     map[string]string
	Struct          reflect.StructField
	IsForeignKey    bool

	tagSettingsLock sync.RWMutex
}

// TagSettingsSet Sets a tag in the tag settings map
func (sf *StructField) TagSettingsSet(key, val string) {
	sf.tagSettingsLock.Lock()
	defer sf.tagSettingsLock.Unlock()
	sf.TagSettings[key] = val
}

// TagSettingsGet returns a tag from the tag settings
func (sf *StructField) TagSettingsGet(key string) (string, bool) {
	sf.tagSettingsLock.RLock()
	defer sf.tagSettingsLock.RUnlock()
	val, ok := sf.TagSettings[key]
	return val, ok
}

// TagSettingsDelete deletes a tag
func (sf *StructField) TagSettingsDelete(key string) {
	sf.tagSettingsLock.Lock()
	defer sf.tagSettingsLock.Unlock()
	delete(sf.TagSettings, key)
}

func (sf *StructField) clone() *StructField {
	clone := &StructField{
		DBName:          sf.DBName,
		Name:            sf.Name,
		Names:           sf.Names,
		IsPrimaryKey:    sf.IsPrimaryKey,
		IsNormal:        sf.IsNormal,
		IsIgnored:       sf.IsIgnored,
		IsScanner:       sf.IsScanner,
		HasDefaultValue: sf.HasDefaultValue,
		Tag:             sf.Tag,
		TagSettings:     map[string]string{},
		Struct:          sf.Struct,
		IsForeignKey:    sf.IsForeignKey,
	}
	// copy the struct field tagSettings, they should be read-locked while they are copied
	sf.tagSettingsLock.Lock()
	defer sf.tagSettingsLock.Unlock()
	for key, value := range sf.TagSettings {
		clone.TagSettings[key] = value
	}

	return clone
}
