package wired_test

import (
	"reflect"
	"testing"

	"github.com/okke/wired"
	"github.com/okke/wired/internal"
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
	return &oneValueStruct{value: wired.Wire().Construct(newEmptyStruct).(*emptyStruct)}
}

func newOneValueStructWithValue(arg *emptyStruct) *oneValueStruct {
	return &oneValueStruct{value: arg}
}

func newStructWithUnknownArgument(arg *unknownStruct) *emptyStruct {
	return &emptyStruct{}
}

func TestCanNotRegisterNonFunctionConstructor(t *testing.T) {

	defer internal.ShouldPanic(t)()

	wired.WithWire(func(wire wired.WireContext) {
		wire.Register("chipotle")
	})
}

func nothing() {

}
func TestCanNotRegisterNonConstructorReturningNothing(t *testing.T) {

	defer internal.ShouldPanic(t)()

	wired.WithWire(func(wire wired.WireContext) {
		wire.Register(nothing)
	})
}

func TestCanNotConstructUnknownType(t *testing.T) {
	defer internal.ShouldPanic(t)()

	wired.WithWire(func(wire wired.WireContext) {
		s := wire.Construct(newStructWithUnknownArgument).(*emptyStruct)
		if s != nil {
			t.Errorf("constructor should not work")
		}
	})

}

func TestConstructWithoutArguments(t *testing.T) {

	wired.WithWire(func(wire wired.WireContext) {
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

	wired.WithWire(func(wire wired.WireContext) {

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
	wired.Autowire
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
	wired.Autowire
	Method OneMethod
}

func newUseOneMethod() *useOneMethod {
	return &useOneMethod{}
}

func TestConstructInterface(t *testing.T) {

	wired.WithWire(func(wire wired.WireContext) {
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
	Increment(int) int
}

var IncrementerType = reflect.TypeOf((*Incrementer)(nil)).Elem()

type firstIncrementer struct {
}

func (firstIncrementer *firstIncrementer) Increment(i int) int {
	return i + 7
}

func newFirstIncrementer() Incrementer {
	return &firstIncrementer{}
}

type secondIncrementer struct {
}

func (secondIncrementer *secondIncrementer) Increment(i int) int {
	return i + 13
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

var CombinedIncrementerType = reflect.TypeOf((*CombinedIncrementer)(nil)).Elem()

type combinedIncrementer struct {
	all []Incrementer
}

func (combinedIncrementer *combinedIncrementer) Increment(i int) int {
	result := i
	for _, increment := range combinedIncrementer.all {
		result = increment.Increment(result)
	}
	return result
}

func newCombinedIncrementer(combine []Incrementer) CombinedIncrementer {
	return &combinedIncrementer{all: combine}
}

func TestConstructWithMultipleConstructors(t *testing.T) {

	wired.WithWire(func(wire wired.WireContext) {
		wire.Register(newFirstIncrementer)         // construct Incrementer
		wire.Register(newSecondIncrementer)        // construct Incrementer
		wire.Register(newAnotherSecondIncrementer) // construct Incrementer
		wire.Register(newCombinedIncrementer)      // construct CombinedIncrementer

		incrementer := wire.ConstructByType(IncrementerType).(Incrementer)
		if i := incrementer.Increment(0); i != 7 {
			t.Error("first incrementer should increase to 7, not", i)
		}

		combined := wire.ConstructByType(CombinedIncrementerType).(Incrementer)
		if i := combined.Increment(0); i != 33 {
			t.Error("expected all incrementers to be called which would result in 33 instead of", i)
		}
	})
}
