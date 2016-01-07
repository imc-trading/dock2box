package controllers

import (
	"fmt"
	"reflect"
)

// structTags return ac onversion map of tags for a struct
func structTags(val reflect.Value, fromTag string, toTag string) (map[string]string, error) {
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported type: %s needs to be a Struct", val.Kind())
	}

	keys := make(map[string]string)
	for i := 0; i < val.NumField(); i++ {
		t := val.Type().Field(i)
		k := t.Tag.Get(fromTag)
		v := t.Tag.Get(toTag)
		keys[k] = v
	}

	return keys, nil
}
