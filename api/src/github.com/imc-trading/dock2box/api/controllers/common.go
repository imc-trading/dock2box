package controllers

import (
	"fmt"
	"reflect"
)

// structTags return a conversion map of tags for a struct
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

// structKinds return a conversion map of kind for tags in a struct
func structKinds(val reflect.Value, tag string) (map[string]reflect.Kind, error) {
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported type: %s needs to be a Struct", val.Kind())
	}

	kinds := make(map[string]reflect.Kind)
	for i := 0; i < val.NumField(); i++ {
		t := val.Type().Field(i)
		k := t.Tag.Get(tag)
		kinds[k] = t.Type.Kind()
	}

	return kinds, nil
}
