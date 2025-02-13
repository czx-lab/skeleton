package utils

import (
	"reflect"
	"strings"
)

// CaseToCamel 下划线转驼峰(大驼峰)
func CaseToCamel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// LowerCamelCase 转换为小驼峰
func LowerCamelCase(name string) string {
	name = CaseToCamel(name)
	return strings.ToLower(name[:1]) + name[1:]
}

func CheckFieldExistence(obj interface{}, name string) bool {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false
	}
	return val.FieldByName(name).IsValid()
}

func StructHasTag(obj interface{}, name string) bool {
	var typ reflect.Type
	switch IsPointer(obj) {
	case true:
		typ = reflect.TypeOf(obj).Elem()
	default:
		typ = reflect.TypeOf(obj)
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(name)
		if len(tag) > 0 {
			return true
		}
	}
	return false
}

func IsPointer(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Ptr
}
