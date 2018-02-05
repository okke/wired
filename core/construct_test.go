package core_test

import (
	"fmt"
	"testing"

	"github.com/okke/wires/core"
	"github.com/okke/wires/internal"
)

type emptyStruct struct {
}

type unknownStruct struct {
}

type oneValueStruct struct {
	value *emptyStruct
}

type singletonStruct struct {
	core.Singleton
	count int
}

type structWithSingleton struct {
	S *singletonStruct
}

func newEmptyStruct() *emptyStruct {
	return &emptyStruct{}
}

func newOneValueStruct() *oneValueStruct {
	return &oneValueStruct{value: core.Wire().Construct(newEmptyStruct).(*emptyStruct)}
}

func newOneValueStructWithValue(arg *emptyStruct) *oneValueStruct {
	return &oneValueStruct{value: arg}
}

func newStructWithUnknownArgument(arg *unknownStruct) *emptyStruct {
	return &emptyStruct{}
}

func newStructWithSingleTon(s *singletonStruct) *structWithSingleton {
	return &structWithSingleton{S: s}
}

func newStructWithSingleTonWithoutArguments() *structWithSingleton {
	return &structWithSingleton{}
}

var structCounter = 0

func newSingletonStruct() *singletonStruct {
	structCounter = structCounter + 1
	return &singletonStruct{count: structCounter}
}

func TestCanNotRegisterNonFunctionConstructor(t *testing.T) {

	defer internal.ShouldPanic(t)()

	core.WithWire(func(wire core.WireContext) {
		wire.Register("chipotle")
	})
}

func nothing() {

}
func TestCanNotRegisterNonConstructorReturningNothing(t *testing.T) {

	defer internal.ShouldPanic(t)()

	core.WithWire(func(wire core.WireContext) {
		wire.Register(nothing)
	})
}

func TestCanNotConstructUnknownType(t *testing.T) {
	defer internal.ShouldPanic(t)()

	core.WithWire(func(wire core.WireContext) {
		s := wire.Construct(newStructWithUnknownArgument).(*emptyStruct)
		if s != nil {
			t.Errorf("constructor should not work")
		}
	})

}

func TestConstructWithoutArguments(t *testing.T) {

	core.WithWire(func(wire core.WireContext) {
		empty := wire.Construct(newEmptyStruct).(*emptyStruct)
		if empty == nil {
			t.Error("expected to construct an empty value")
		}

		oneValue := wire.Construct(newOneValueStruct).(*oneValueStruct)
		if oneValue == nil {
			t.Error("expected to construct a one value")
		}
	})
}

func TestConstructWithoutOneArgument(t *testing.T) {

	core.WithWire(func(wire core.WireContext) {

		wire.Register(newEmptyStruct)

		oneValue := wire.Construct(newOneValueStructWithValue).(*oneValueStruct)
		if oneValue == nil {
			t.Error("expected to construct a one value struct")
		}
		if oneValue.value == nil {
			t.Error("expected struct to contain a also constructed value")
		}
	})
}

func TestCannotAddNilDecorator(t *testing.T) {

	core.WithWire(func(wire core.WireContext) {
		defer internal.ShouldPanic(t)()

		wire.AddDecorator(nil)
	})

}

func TestDecoratorsShouldBeCalled(t *testing.T) {
	core.WithWire(func(wire core.WireContext) {
		var changed = false
		wire.AddDecorator(core.DecoratorFunc(func(obj interface{}) interface{} {
			changed = true
			return obj
		}))
		constructed := wire.Construct(newEmptyStruct)
		if constructed == nil {
			t.Error("expected a struct to be constructed")
		}
		if !changed {
			t.Error("expected a decorator to be called")
		}
	})
}

func TestDecoratorShouldNotReturnNil(t *testing.T) {

	fmt.Println("-->")
	core.WithWire(func(wire core.WireContext) {
		defer internal.ShouldPanic(t)()

		wire.AddDecorator(core.DecoratorFunc(func(obj interface{}) interface{} {
			return nil
		}))
		constructed := wire.Construct(newEmptyStruct)

		t.Error("did not expect to construct", constructed)

	})
}

func TestConstructSingleton(t *testing.T) {

	core.WithWire(func(wire core.WireContext) {
		wire.Register(newSingletonStruct)
		first := wire.Construct(newSingletonStruct)
		second := wire.Construct(newSingletonStruct)

		if first != second {
			t.Error("did not construct singletons", first, second)
		}

		// singleton should be used as constructor argument
		//
		firstUsage := wire.Construct(newStructWithSingleTon).(*structWithSingleton)
		if firstUsage.S != first {
			t.Error("did not use singleton", firstUsage.S, first)
		}

		// singleton should be used when a field is wired into a struct
		//
		secondUsage := wire.Construct(newStructWithSingleTonWithoutArguments).(*structWithSingleton)
		if secondUsage.S != first {
			t.Error("did not use singleton", secondUsage.S, first)
		}
	})
}

type deeper struct {
}

type Deeper interface {
	Deep()
}

func (deeper *deeper) Deep() {
}

func newDeeper() Deeper {
	return &deeper{}
}

type oneMethod struct {
	deep Deeper
}

type OneMethod interface {
	Do(t *testing.T)
}

func (oneMethod *oneMethod) SetDeep(d Deeper) {
	oneMethod.deep = d
}

func (oneMethod *oneMethod) GetDeep() Deeper {
	return oneMethod.deep
}

func (oneMethod *oneMethod) Do(t *testing.T) {
	if oneMethod.deep == nil {
		t.Fatal("need a deeper object")
	}
}

func newOneMethod() OneMethod {
	return &oneMethod{}
}

type useOneMethod struct {
	Method OneMethod
}

func newUseOneMethod() *useOneMethod {
	return &useOneMethod{}
}

func TestConstructInterface(t *testing.T) {

	core.WithWire(func(wire core.WireContext) {
		wire.Register(newDeeper)
		wire.Register(newOneMethod)
		c := wire.Construct(newUseOneMethod).(*useOneMethod)
		if c.Method == nil {
			t.Fatal("expected method object to be constructed")
		}
		c.Method.Do(t)
	})
}
