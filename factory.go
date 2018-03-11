package wired

import (
	"reflect"
)

// Factory tags an object to act as a construction factory for other objects
//
type Factory struct {
}

type factory struct {
	// a factory will act as a singleton
	//
	singleton
}

func init() {
	RegisterConstructionTag(reflect.TypeOf((*Factory)(nil)).Elem(), &factory{})
}

// Apply factory creation logic
//
func (factory *factory) Apply(scope Scope, objType reflect.Type, constructor func() interface{}) interface{} {

	constructed := factory.singleton.Apply(scope, objType, constructor)

	factoryMethod := reflect.ValueOf(constructed).MethodByName("Construct")

	if factoryMethod.Kind() != reflect.Invalid {
		scope.RegisterForType(factoryMethod.Type().Out(0), func() interface{} {
			return factoryMethod.Call([]reflect.Value{})[0].Interface()
		})
	}

	return constructed
}

func (factory *factory) ShouldAutoConstruct() bool {
	return true
}
