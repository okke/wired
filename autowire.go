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

<<<<<<< HEAD
	// when a field has a value, auto-wiring is not applicable
	//
=======
>>>>>>> 9eb2bcf052d783e6693e9e20d36fc68425158c73
	if originalValue := internal.GetFieldValueByReflection(obj, field, fieldType); originalValue != nil {
		return internal.NilValue, false
	}

<<<<<<< HEAD
	// when wired does not know how to construct a type, let go
	//
=======
>>>>>>> 9eb2bcf052d783e6693e9e20d36fc68425158c73
	if value := wire.ConstructByType(fieldType.Type); value == nil {
		return reflect.ValueOf(nil), false
	}

	// return value that can be auto wired
	//
	return reflect.ValueOf(wire.ConstructByType(fieldType.Type)), true
}
