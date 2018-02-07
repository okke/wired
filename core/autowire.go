package core

import (
	"reflect"
	"strings"
)

// Autowire is a tag that drives autowiring of struct fields
//
type Autowire struct {
}

type autowire struct {
}

func init() {
	RegisterStructDecorationTag(reflect.TypeOf((*Autowire)(nil)).Elem(), &autowire{})
}

func (autowire *autowire) getFieldValue(obj reflect.Value, field reflect.Value, fieldType reflect.StructField) interface{} {

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

func (autowire *autowire) getValueFor(wire WireContext, obj reflect.Value, field reflect.Value, fieldType reflect.StructField) (reflect.Value, bool) {

	originalValue := autowire.getFieldValue(obj, field, fieldType)

	// TODO: DO WE WANT THIS BEHAVIOUR?
	//
	if originalValue != nil {
		return reflect.ValueOf(autowire.Apply(wire, originalValue)), true
	}

	value := wire.ConstructByType(fieldType.Type)
	if value == nil {
		return reflect.ValueOf(nil), false
	}

	return reflect.ValueOf(wire.ConstructByType(fieldType.Type)), true
}

func (autowire *autowire) doDecorateStruct(wire WireContext, objValue reflect.Value, objType reflect.Type) {

	for walk := 0; walk < objType.NumField(); walk++ {

		field := objValue.Field(walk)
		fieldType := objType.Field(walk)

		value, shouldSet := autowire.getValueFor(wire, objValue, field, fieldType)

		if !shouldSet {
			continue
		}

		if field.CanSet() {
			field.Set(value)
		} else {
			name := strings.Join([]string{"Set", strings.Title(objType.Field(walk).Name)}, "")
			setter := objValue.Addr().MethodByName(name)
			if setter.Kind() != reflect.Invalid {
				setter.Call([]reflect.Value{value})
			}
		}

	}
}

func (autowire *autowire) doDecorate(wire WireContext, objValue reflect.Value, objType reflect.Type) {

	if objType.Kind() == reflect.Ptr {
		autowire.doDecorate(wire, objValue.Elem(), objType.Elem())
		return
	}

	if objType.Kind() == reflect.Struct {
		autowire.doDecorateStruct(wire, objValue, objType)
		return
	}

}

func (autowire *autowire) Apply(context WireContext, object interface{}) interface{} {
	autowire.doDecorate(context, reflect.ValueOf(object), reflect.TypeOf(object))
	return object
}
