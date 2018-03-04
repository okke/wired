package wired

import (
	"reflect"

	"github.com/okke/wired/internal"
)

// AutoWire is a tag that drives autowiring of struct fields
//
type AutoWire struct {
}

type autowire struct {
}

func init() {
	RegisterStructDecorationTag(reflect.TypeOf((*AutoWire)(nil)).Elem(), &autowire{})
}

func (autowire *autowire) GetValueFor(wire Scope, obj reflect.Value, field reflect.Value, fieldType reflect.StructField) (reflect.Value, bool) {

	originalValue := internal.GetFieldValueByReflection(obj, field, fieldType)

	if originalValue != nil {
		return internal.NilValue, false
	}

	value := wire.ConstructByType(fieldType.Type)
	if value == nil {
		return reflect.ValueOf(nil), false
	}

	return reflect.ValueOf(wire.ConstructByType(fieldType.Type)), true
}
