package wired

import "reflect"

// ConstructionTag can be used to apply construction logic
//
type ConstructionTag interface {
	Apply(scope Scope, objType reflect.Type, constructor func() interface{}) interface{}
	ShouldAutoConstruct() bool
}

// StructDecorationTag can be used to initialize a struct after it has been constructed
//
type StructDecorationTag interface {

	// GetValueFor will be called to determine field values. It should return a valid value and true
	// When no valid value could be found, return reflect.ValueOf(nil) and false
	//
	GetValueFor(wire Scope, obj reflect.Value, field reflect.Value, fieldType reflect.StructField) (reflect.Value, bool)
}

var constructionTags = make(map[reflect.Type]ConstructionTag, 10)
var structDecorationTags = make(map[reflect.Type]StructDecorationTag, 10)

// RegisterConstructionTag connects a struct tag to its constructions logic
//
func RegisterConstructionTag(objType reflect.Type, tag ConstructionTag) {
	constructionTags[objType] = tag
}

// RegisterStructDecorationTag connects a struct tag to its decoration logic
//
func RegisterStructDecorationTag(objType reflect.Type, tag StructDecorationTag) {
	structDecorationTags[objType] = tag
}

// FindConstructionTag looks at the fields of a given type and returns
// the first constructiontag
//
func FindConstructionTag(objType reflect.Type) (ConstructionTag, bool) {

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	if objType.Kind() == reflect.Struct {
		numFields := objType.NumField()
		for walk := 0; walk < numFields; walk++ {

			if tag, found := constructionTags[objType.Field(walk).Type]; found {
				return tag, true
			}

		}
	}

	return nil, false
}

// FindStructDecorationTags looks for all struct tags that are known to implements
// the struct decoration interface
//
func FindStructDecorationTags(objType reflect.Type) []StructDecorationTag {

	result := make([]StructDecorationTag, 0, 0)

	if objType.Kind() == reflect.Struct {
		numFields := objType.NumField()
		for walk := 0; walk < numFields; walk++ {

			if tag, found := structDecorationTags[objType.Field(walk).Type]; found {
				result = append(result, tag)
			}

		}
	}

	return result
}
