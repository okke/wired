package core

import (
	"reflect"
	"strings"
)

type fieldDecorator struct {
}

func (fieldDecorator *fieldDecorator) getFieldValue(obj reflect.Value, field reflect.Value, fieldType reflect.StructField) interface{} {

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

func (fieldDecorator *fieldDecorator) getValueFor(wire WireContext, obj reflect.Value, field reflect.Value, fieldType reflect.StructField) (reflect.Value, bool) {

	originalValue := fieldDecorator.getFieldValue(obj, field, fieldType)

	if originalValue != nil {
		return reflect.ValueOf(fieldDecorator.Decorate(wire, originalValue)), true
	}

	value := wire.ConstructByType(fieldType.Type)
	if value == nil {
		return reflect.ValueOf(nil), false
	}

	return reflect.ValueOf(wire.ConstructByType(fieldType.Type)), true
}

func (fieldDecorator *fieldDecorator) doDecorateStruct(wire WireContext, objValue reflect.Value, objType reflect.Type) {

	for walk := 0; walk < objType.NumField(); walk++ {

		field := objValue.Field(walk)
		fieldType := objType.Field(walk)

		value, shouldSet := fieldDecorator.getValueFor(wire, objValue, field, fieldType)

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

func (fieldDecorator *fieldDecorator) doDecorate(wire WireContext, objValue reflect.Value, objType reflect.Type) {

	if objType.Kind() == reflect.Ptr {
		fieldDecorator.doDecorate(wire, objValue.Elem(), objType.Elem())
		return
	}

	if objType.Kind() == reflect.Struct {
		fieldDecorator.doDecorateStruct(wire, objValue, objType)
		return
	}

}

// Decorate a just constructed object (so for example fields can be automagically initialized)
//
func (fieldDecorator *fieldDecorator) Decorate(wire WireContext, obj interface{}) interface{} {

	fieldDecorator.doDecorate(wire, reflect.ValueOf(obj), reflect.TypeOf(obj))
	return obj
}

// NewFieldDecorator constructs a decorator that can inspect the fields of a constructed object
//
func NewFieldDecorator() Decorator {
	return &fieldDecorator{}
}
