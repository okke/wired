package wired

import "reflect"

// Singleton is a struct that can be mixed into another struct
// to express this struct must be used as singleton
//
type Singleton struct {
}

type singleton struct {
}

func init() {
	RegisterConstructionTag(reflect.TypeOf((*Singleton)(nil)).Elem(), &singleton{})
}

// Apply applies singleton creation logic
//
func (singleton *singleton) Apply(scope Scope, objType reflect.Type, constructor func() interface{}) interface{} {

	if object, found := scope.FindSingleton(objType); found {
		return object
	}

	constructed := constructor()
	scope.RegisterSingleton(objType, constructed)
	return constructed
}

func (singleton *singleton) ShouldAutoConstruct() bool {
	return false
}
