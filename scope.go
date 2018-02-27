package wired

import (
	"fmt"
	"reflect"

	"github.com/okke/wired/internal"
)

type scope struct {
	constructorMapping map[reflect.Type]interface{}
	singletons         map[reflect.Type]interface{}
	parent             *scope
	top                *scope
}

// Scope is an interface decribing the main functions used to wire objects
//
type Scope interface {
	Register(constructor interface{})

	Construct(use interface{}) interface{}
	ConstructByType(reflect.Type) interface{}

	FindSingleton(objType reflect.Type) (interface{}, bool)
	RegisterSingleton(objType reflect.Type, value interface{})

	Go(f func(Scope))
}

func newScope(parent *scope) Scope {
	var top *scope
	if parent != nil && parent.top != nil {
		top = parent.top
	}

	context := &scope{
		constructorMapping: make(map[reflect.Type]interface{}, 100),
		singletons:         make(map[reflect.Type]interface{}, 100),
		parent:             parent,
		top:                top}

	if context.top == nil {
		context.top = context
	}

	return context
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

	if t.NumOut() != 1 {
		panic("constructor does not return one object")
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
	value, found := scope.top.singletons[objType]
	return value, found
}

func (scope *scope) RegisterSingleton(objType reflect.Type, value interface{}) {
	scope.top.singletons[objType] = value
}

func (scope *scope) registerSliceConstructor(constructor interface{}, constructorType reflect.Type) {

	sliceType := reflect.SliceOf(constructorType)

	var knownConstructor interface{}
	if knownSliceConstructor, foundSliceConstructor := scope.findConstructor(sliceType); foundSliceConstructor {
		knownConstructor = knownSliceConstructor
	} else {
		knownConstructor, _ = scope.findConstructor(constructorType)
	}

	scope.constructorMapping[sliceType] = func() interface{} {
		if knownConstructor == nil {
			return internal.CreateSliceWithValues(
				sliceType,
				scope.Construct(constructor)).Interface()
		}

		return internal.CreateSliceWithValues(
			sliceType,
			scope.Construct(constructor),
			scope.Construct(knownConstructor)).Interface()
	}

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
func (scope *scope) Register(constructor interface{}) {
	constructorType := ensureConstructorIsAFunction(constructor).Out(0)

	scope.registerSliceConstructor(constructor, constructorType)

	scope.constructorMapping[constructorType] = constructor

}

func (scope *scope) decorate(obj interface{}) interface{} {
	var result = obj
	for _, decorator := range FindStructDecorationTags(reflect.TypeOf(obj)) {
		result = decorator.Apply(scope, result)
	}
	return result
}

// Construct takes a function and tries to call it by filling
// in the arguments through the execution of registered constructor functions
//
func (scope *scope) Construct(use interface{}) interface{} {

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

		return scope.decorate(reflect.ValueOf(use).Call(in)[0].Interface())
	}

	outType := constructorType.Out(0)

	if tag, found := FindConstructionTag(outType); found {
		return tag.Apply(scope, outType, constructByReflection)
	}

	return constructByReflection()
}

// ConstructByType constructs a type by looking up its registered constructor
// function.
//
func (scope *scope) ConstructByType(objType reflect.Type) interface{} {

	argConstructor, found := scope.findConstructor(objType)
	if !found {
		return nil
	}
	return scope.Construct(argConstructor)
}
