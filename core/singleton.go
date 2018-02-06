package core

import "reflect"

// Singleton is a struct that can be mixed into another struct
// to express this struct must be used as singleton
//
type Singleton struct {
}

var singletonTag = Singleton{}

func init() {
	RegisterConstructionTag(reflect.TypeOf(singletonTag), singletonTag)
}

// Apply applies singleton creation logic
//
func (singleton Singleton) Apply(wireContext WireContext, objType reflect.Type, constructor func() interface{}) interface{} {

	if object, found := wireContext.FindSingleton(objType); found {
		return object
	}

	constructed := constructor()
	wireContext.RegisterSingleton(objType, constructed)
	return constructed
}
