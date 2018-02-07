package core

import (
	"reflect"
)

type wireContext struct {
	constructorMapping map[reflect.Type]interface{}
	singletons         map[reflect.Type]interface{}
}

// WireContext is an interface decribing the main functions used to wire objects
//
type WireContext interface {
	Register(constructor interface{})

	Construct(use interface{}) interface{}
	ConstructByType(reflect.Type) interface{}

	FindSingleton(objType reflect.Type) (interface{}, bool)
	RegisterSingleton(objType reflect.Type, value interface{})
}

func newWireContext() WireContext {
	context := &wireContext{
		constructorMapping: make(map[reflect.Type]interface{}, 100),
		singletons:         make(map[reflect.Type]interface{}, 100)}

	return context
}

var globalWireContext = newWireContext()

// Wire will return a context that can be used to register constructors and
// to construct actual objects
//
func Wire() WireContext {
	return globalWireContext
}

// WithWire will create a new WireContext and use the callback to do whatever
// you like to do with the created context
//
func WithWire(f func(WireContext)) {
	f(newWireContext())
}

func ensureConstructorIsAFunction(constructor interface{}) reflect.Type {
	t := reflect.TypeOf(constructor)

	if t.Kind() != reflect.Func {
		panic("constructor is not a function")
	}

	if t.NumOut() != 1 {
		panic("constructor does not return one object")
	}

	return t
}

// Register will register a constructor for a given Type
//
// Best practice is to use the register function in the init()
// of your package so all types your packages exposes, are unknown
// and can be wired
//
// note, the constructor is defined as interface{} so it will
// accept all kind of values. But it actually must be a function
// otherwise Register will panic
//
func (wireContext *wireContext) Register(constructor interface{}) {
	wireContext.constructorMapping[ensureConstructorIsAFunction(constructor).Out(0)] = constructor
}

func (wireContext *wireContext) decorate(obj interface{}) interface{} {
	var result = obj
	for _, decorator := range FindStructDecorationTags(reflect.TypeOf(obj)) {
		result = decorator.Apply(wireContext, result)
	}
	return result
}

func (wireContext *wireContext) FindSingleton(objType reflect.Type) (interface{}, bool) {
	value, found := wireContext.singletons[objType]
	return value, found
}

func (wireContext *wireContext) RegisterSingleton(objType reflect.Type, value interface{}) {
	wireContext.singletons[objType] = value
}

// Construct takes a function and tries to call it by filling
// in the arguments through the execution of registered constructor functions
//
func (wireContext *wireContext) Construct(use interface{}) interface{} {

	constructorType := ensureConstructorIsAFunction(use)

	constructByReflection := func() interface{} {
		in := make([]reflect.Value, constructorType.NumIn())
		for i := range in {
			in[i] = reflect.ValueOf(wireContext.ConstructByType(constructorType.In(i)))
		}

		return wireContext.decorate(reflect.ValueOf(use).Call(in)[0].Interface())
	}

	outType := constructorType.Out(0)

	if tag, found := FindConstructionTag(outType); found {
		return tag.Apply(wireContext, outType, constructByReflection)
	}

	return constructByReflection()
}

// ConstructByType constructs a type by looking up its registered constructor
// function.
//
func (wireContext *wireContext) ConstructByType(objType reflect.Type) interface{} {

	argConstructor, found := wireContext.constructorMapping[objType]
	if !found {
		return nil
	}
	return wireContext.Construct(argConstructor)
}
