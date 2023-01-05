package main

import (
	"reflect"
	"time"
)

type TerraformResourceData interface {
	HasChange(string) bool
	GetOkExists(string) (interface{}, bool)
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
	Set(string, interface{}) error
	SetId(string)
	Id() string
	GetProviderMeta(interface{}) error
	Timeout(key string) time.Duration
}

func IsEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func RemoveEmptyOptions(options map[string]interface{}) {
	var emptyKeys []string
	for k, v := range options {
		if v == nil {
			emptyKeys = append(emptyKeys, k)
		}

		s, ok := v.(string)
		if ok && s == "" {
			emptyKeys = append(emptyKeys, k)
		}

		i, ok := v.(int)
		if ok && i == 0 {
			emptyKeys = append(emptyKeys, k)
		}
	}

	for _, k := range emptyKeys {
		delete(options, k)
	}
}
