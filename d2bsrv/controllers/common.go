package controllers

import (
	"fmt"
	"reflect"
)

// structTags return tags for a struct
func structTags(val reflect.Value, tag string) (map[string]int, error) {
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported type: %s needs to be a Struct", val.Kind())
	}

	keys := make(map[string]int)
	for i := 0; i < val.NumField(); i++ {
		t := val.Type().Field(i)
		k := t.Tag.Get(tag)
		keys[k] = i
	}

	return keys, nil
}
