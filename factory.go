package wired

import (
	"reflect"
	"strings"
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
	constructedValue := reflect.ValueOf(constructed)
	constructedType := reflect.TypeOf(constructed)

	for walk := 0; walk < constructedType.NumMethod(); walk++ {
		factoryMethodType := constructedType.Method(walk)
		if strings.HasPrefix(factoryMethodType.Name, "Construct") {
			factoryMethod := constructedValue.MethodByName(factoryMethodType.Name)

			if factoryMethod.Kind() != reflect.Invalid {
				scope.Register(factoryMethod.Interface())
			}
		}
	}

	return constructed
}

func (factory *factory) ShouldAutoConstruct() bool {
	return true
}
