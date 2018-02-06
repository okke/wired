package core

import "reflect"

// ConstructionTag can be used to apply construction logic
//
type ConstructionTag interface {
	Apply(wireContext WireContext, objType reflect.Type, constructor func() interface{}) interface{}
}

var constructionTags = make(map[reflect.Type]ConstructionTag, 10)

// RegisterConstructionTag connects a struct tag to its constructions logic
//
func RegisterConstructionTag(objType reflect.Type, tag ConstructionTag) {
	constructionTags[objType] = tag
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
