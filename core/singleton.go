package core

import "reflect"

// Singleton is a struct that can be mixed into another struct
// to express this struct must be used as singleton
//
type Singleton struct {
}

var singletonTag = Singleton{}

func (singleton Singleton) isSingleton(objType reflect.Type) bool {

	// TODO: how to implement singleton interfaces ...

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	if objType.Kind() == reflect.Struct {
		numFields := objType.NumField()
		for walk := 0; walk < numFields; walk++ {
			if objType.Field(walk).Type == reflect.TypeOf(singleton) {
				return true
			}
		}
	}

	return false
}
