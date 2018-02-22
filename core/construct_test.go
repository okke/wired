package core_test

import (
	"reflect"
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
	core.Autowire
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
	core.Autowire
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

// ------ Test multiple constructors for the same type ----

type Incrementer interface {
	Increment(*int)
}

type firstIncrementer struct {
}

func (firstIncrementer *firstIncrementer) Increment(i *int) {
	*i = *i + 7
}

func newFirstIncrementer() Incrementer {
	return &firstIncrementer{}
}

type secondIncrementer struct {
}

func (secondIncrementer *secondIncrementer) Increment(i *int) {
	*i = *i + 13
}

func newSecondIncrementer() Incrementer {
	return &secondIncrementer{}
}

func newAnotherSecondIncrementer() Incrementer {
	return &secondIncrementer{}
}

// to prevend a recurse loop when resolving, introduce a separate interface
// for the combined incrementer
//
type CombinedIncrementer interface {
	Incrementer
}

type combinedIncrementer struct {
	all []Incrementer
}

func (combinedIncrementer *combinedIncrementer) Increment(i *int) {
	for _, increment := range combinedIncrementer.all {
		increment.Increment(i)
	}
}

func newCombinedIncrementer(combine []Incrementer) CombinedIncrementer {
	return &combinedIncrementer{all: combine}
}

func TestConstructWithMultipleConstructors(t *testing.T) {

	core.WithWire(func(wire core.WireContext) {
		wire.Register(newFirstIncrementer)         // construct Incrementer
		wire.Register(newSecondIncrementer)        // construct Incrementer
		wire.Register(newAnotherSecondIncrementer) // construct Incrementer
		wire.Register(newCombinedIncrementer)      // construct CombinedIncrementer

		giveMeOne := reflect.TypeOf((*Incrementer)(nil)).Elem()
		if incrementer, couldCast := wire.ConstructByType(giveMeOne).(Incrementer); couldCast {
			i := 0
			incrementer.Increment(&i)
			if i != 7 {
				t.Error("first incrementer should increase to 7, not", i)
			}
		} else {
			t.Error("expected an incrementer, not", incrementer)
		}

		lookingFor := reflect.TypeOf((*CombinedIncrementer)(nil)).Elem()
		t.Log(lookingFor)
		combined := wire.ConstructByType(lookingFor).(Incrementer)
		i := 0
		combined.Increment(&i)
		if i != 33 {
			t.Error("expected all incrementers to be called which would result in 33 instead of", i)
		}
	})
}
