package internal

import (
	"reflect"
	"strings"
)

// NilValue can be used where a Value is required but no value can be determined
//
var NilValue = reflect.ValueOf(nil)

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

// CreateSliceWithValues creates a slice of given type and appends given values
// When a value is a slice itself, all slice values will be added
//
func CreateSliceWithValues(sliceType reflect.Type, values ...interface{}) reflect.Value {
	slice := reflect.MakeSlice(sliceType, 0, len(values))
	for _, value := range values {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			slice = reflect.AppendSlice(slice, reflect.ValueOf(value))
		} else {
			slice = reflect.Append(slice, reflect.ValueOf(value))
		}
	}
	return slice
}
