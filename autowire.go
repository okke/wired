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

// GetValueFor implements the StructDecorationTag interface
//
func (autowire *autowire) GetValueFor(wire Scope, obj reflect.Value, field reflect.Value, fieldType reflect.StructField) (reflect.Value, bool) {

	// when a field has a value, auto-wiring is not applicable
	//
	if originalValue := internal.GetFieldValueByReflection(obj, field, fieldType); originalValue != nil {
		return internal.NilValue, false
	}

	// when wired does not know how to construct a type, let go
	//
	if value := wire.ConstructByType(fieldType.Type); value == nil {
		return reflect.ValueOf(nil), false
	}

	// return value that can be auto wired
	//
	return reflect.ValueOf(wire.ConstructByType(fieldType.Type)), true
}
