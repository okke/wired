package wired

import (
	"fmt"
	"reflect"

	"github.com/okke/wired/internal"
)

type scope struct {
	constructorMapping map[reflect.Type]interface{} // map type to constructor functions
	singletons         map[reflect.Type]interface{} // map type to singleton objects
	parent             *scope
}

var scopeType = reflect.TypeOf((*scope)(nil))

// Scope is an interface describing the main functions used to wire objects
//
type Scope interface {

	// Register a constructor function
	//
	Register(constructor interface{})

	// Construct an object by providing a constructor function. Wired will
	// inject valid function arguments.
	//
	Construct(use interface{}) interface{}

	// Call a function and autowire function arguments
	//
	Inject(use interface{})

	// Construct an object by providing the type that needs to be created.
	//
	ConstructByType(reflect.Type) interface{}

	// Lookup a singleton by type and return it. When the singleton is found
	// the returned bool will be true. Otherwise it will be false.
	//
	FindSingleton(objType reflect.Type) (interface{}, bool)

	// Register a singleton for a given type. Note, using this method
	// directly is not recommended. Use the Singleton tag instead
	//
	RegisterSingleton(objType reflect.Type, value interface{})

	// Construct a sub scope and use it within given function
	//
	Go(f func(Scope))
}

func newScope(parent *scope) Scope {
	return &scope{
		constructorMapping: make(map[reflect.Type]interface{}, 0),
		singletons:         make(map[reflect.Type]interface{}, 0),
		parent:             parent}
}

var globalScope = newScope(nil)

// Global will return a global scope that can be used to register constructors and
// to construct actual objects
//
func Global() Scope {
	return globalScope
}

// Go will create a new Scope and use the callback to do whatever
// you like to do with the created scope
//
func Go(f func(Scope)) {
	f(newScope(nil))
}

func ensureConstructorIsAFunction(constructor interface{}) reflect.Type {
	t := reflect.TypeOf(constructor)

	if t.Kind() != reflect.Func {
		panic("constructor is not a function")
	}

	return t
}

func (scope *scope) Go(f func(Scope)) {
	f(newScope(scope))
}

func (scope *scope) findConstructor(objType reflect.Type) (interface{}, bool) {
	result, found := scope.constructorMapping[objType]
	if !found && scope.parent != nil {
		result, found = scope.parent.findConstructor(objType)
	}
	return result, found
}

func (scope *scope) FindSingleton(objType reflect.Type) (interface{}, bool) {
	if value, found := scope.singletons[objType]; found {
		return value, true
	}

	if scope.parent != nil {
		value, found := scope.parent.FindSingleton(objType)
		return value, found
	}

	return nil, false
}

func (scope *scope) RegisterSingleton(objType reflect.Type, value interface{}) {
	scope.singletons[objType] = value
}

func (wire *scope) registerSliceConstructor(constructor interface{}, constructorType reflect.Type) {

	sliceType := reflect.SliceOf(constructorType)

	knownConstructor, found := wire.findConstructor(sliceType)

	wire.constructorMapping[sliceType] = func(activeScope *scope) interface{} {
		if !found {
			return internal.CreateSliceWithValues(
				sliceType,
				activeScope.Construct(constructor)).Interface()
		}

		return internal.CreateSliceWithValues(
			sliceType,
			activeScope.Construct(constructor),
			activeScope.Construct(knownConstructor)).Interface()
	}

}

func (wire *scope) registerMapConstructor(constructor interface{}, constructorType reflect.Type) {

	keyMethod, found := constructorType.MethodByName("Key")
	if !found {
		return
	}

	if keyMethod.Type.NumOut() == 0 {
		return
	}

	keyType := keyMethod.Type.Out(0)
	mapType := reflect.MapOf(keyType, constructorType)

	knownConstructor, found := wire.findConstructor(mapType)

	wire.constructorMapping[mapType] = func(activeScope *scope) interface{} {
		constructed := activeScope.Construct(constructor)
		constructedKey := reflect.ValueOf(constructed).MethodByName("Key").Call([]reflect.Value{})[0].Interface()
		if !found {
			return internal.CreateMapWithValues(
				mapType,
				constructedKey,
				constructed).Interface()
		}

		return internal.CreateMapWithValues(
			mapType,
			constructedKey,
			constructed,
			activeScope.Construct(knownConstructor)).Interface()
	}
}

func (scope *scope) Register(constructor interface{}) {

	constructorType := ensureConstructorIsAFunction(constructor).Out(0)

	// ensure we know how to construct slices of given type
	//
	scope.registerSliceConstructor(constructor, constructorType)
	scope.registerMapConstructor(constructor, constructorType)

	scope.constructorMapping[constructorType] = constructor

	if constructorTag, found := FindConstructionTag(constructorType); found {
		if constructorTag.ShouldAutoConstruct() {
			scope.ConstructByType(constructorType)
		}
	}
}

func (scope *scope) doDecorateStruct(objValue reflect.Value, objType reflect.Type) {

	decorators := FindStructDecorationTags(objType)

	for walk := 0; walk < objType.NumField(); walk++ {

		field := objValue.Field(walk)
		fieldType := objType.Field(walk)

		for _, decorator := range decorators {

			value, shouldSet := decorator.GetValueFor(scope, objValue, field, fieldType)

			if shouldSet && value.Type().AssignableTo(field.Type()) {
				internal.SetFieldValueByReflection(objValue, field, fieldType, value)
			}
		}
	}
}

func (scope *scope) doDecorate(objValue reflect.Value, objType reflect.Type) {

	if objType.Kind() == reflect.Ptr {
		scope.doDecorate(objValue.Elem(), objType.Elem())
		return
	}

	if objType.Kind() == reflect.Struct {
		scope.doDecorateStruct(objValue, objType)
		return
	}

}

func (scope *scope) decorate(obj interface{}) interface{} {
	scope.doDecorate(reflect.ValueOf(obj), reflect.TypeOf(obj))
	return obj
}

// Inject takes a function and tries to call it by filling
// in the arguments however it won't return a value
// It's mostly used for readability purpose so one can write
//
// scope.Inject(func (a ArgType, y YetAnotherArgtype) {
//
// })
//
func (scope *scope) Inject(use interface{}) {
	scope.construct(use)
}

// Construct takes a function and tries to call it by filling
// in the arguments through the execution of registered constructor functions
//
func (scope *scope) Construct(use interface{}) interface{} {
	return scope.construct(use)
}

func (scope *scope) construct(use interface{}) interface{} {
	constructorType := ensureConstructorIsAFunction(use)

	constructByReflection := func() interface{} {
		in := make([]reflect.Value, constructorType.NumIn())
		for i := range in {
			if arg := scope.ConstructByType(constructorType.In(i)); arg != nil {
				in[i] = reflect.ValueOf(arg)
			} else {
				panic(fmt.Sprintf("do not know how to construct %v", constructorType.In(i)))
			}

		}

		results := reflect.ValueOf(use).Call(in)

		if constructorType.NumOut() > 0 {
			return scope.decorate(results[0].Interface())
		}

		return nil

	}

	if constructorType.NumOut() > 0 {
		outType := constructorType.Out(0)

		if tag, found := FindConstructionTag(outType); found {
			return tag.Apply(scope, outType, constructByReflection)
		}
	}

	return constructByReflection()
}

// ConstructByType constructs a type by looking up its registered constructor
// function.
//
func (scope *scope) ConstructByType(objType reflect.Type) interface{} {

	// when looking for a scope, always return the current scope
	//
	if objType == scopeType {
		return scope
	}

	argConstructor, found := scope.findConstructor(objType)
	if !found {
		if objType.Kind() == reflect.Slice {
			return internal.CreateSliceWithValues(objType).Interface()
		}
		return nil
	}
	return scope.Construct(argConstructor)
}
