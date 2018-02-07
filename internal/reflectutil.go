package internal

import (
	"reflect"
	"strings"
)

// GetFieldValueByReflection will retrieve a structs field true reflection by either
// accesing it like a public field or by using its getter
//
func GetFieldValueByReflection(obj reflect.Value, field reflect.Value, fieldType reflect.StructField) interface{} {

	if field.CanSet() {
		if field.Kind() == reflect.Ptr && field.IsNil() {
			return nil
		}
		return field.Interface()
	}

	name := strings.Join([]string{"Get", strings.Title(fieldType.Name)}, "")
	getter := obj.Addr().MethodByName(name)
	if getter.Kind() != reflect.Invalid {
		result := getter.Call([]reflect.Value{})
		if result[0].Kind() == reflect.Ptr && result[0].IsNil() {
			return nil
		}
		return result[0].Interface()
	}
	return nil
}

// SetFieldValueByReflection will set a fields value through reflection. Either by
// accessing it as a public field. Or by using its setter method
//
func SetFieldValueByReflection(objValue reflect.Value, field reflect.Value, fieldType reflect.StructField, value reflect.Value) {
	if field.CanSet() {
		field.Set(value)
	} else {
		name := strings.Join([]string{"Set", strings.Title(fieldType.Name)}, "")
		setter := objValue.Addr().MethodByName(name)
		if setter.Kind() != reflect.Invalid {
			setter.Call([]reflect.Value{value})
		}
	}
}
