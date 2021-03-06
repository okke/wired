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
	return &oneValueStruct{value: wired.Global().Construct(newEmptyStruct).(*emptyStruct)}
}

func newOneValueStructWithValue(arg *emptyStruct) *oneValueStruct {
	return &oneValueStruct{value: arg}
}

func newStructWithUnknownArgument(arg *unknownStruct) *emptyStruct {
	return &emptyStruct{}
}

func TestCanNotRegisterNonFunctionConstructor(t *testing.T) {

	defer internal.ShouldPanic(t)()

	wired.Go(func(wire wired.Scope) {
		wire.Register("chipotle")
	})
}

func nothing() {

}
func TestCanNotRegisterNonConstructorReturningNothing(t *testing.T) {

	defer internal.ShouldPanic(t)()

	wired.Go(func(wire wired.Scope) {
		wire.Register(nothing)
	})
}

func TestCanNotConstructUnknownType(t *testing.T) {
	defer internal.ShouldPanic(t)()

	wired.Go(func(scope wired.Scope) {
		s := scope.Construct(newStructWithUnknownArgument).(*emptyStruct)
		if s != nil {
			t.Errorf("constructor should not work")
		}
	})

}

func TestConstructWithoutArguments(t *testing.T) {

	wired.Go(func(scope wired.Scope) {
		empty := scope.Construct(newEmptyStruct).(*emptyStruct)
		if empty == nil {
			t.Error("expected to construct an empty value")
		}

		oneValue := scope.Construct(newOneValueStruct).(*oneValueStruct)
		if oneValue == nil {
			t.Error("expected to construct a one value")
		}
	})
}

func TestConstructWithoutOneArgument(t *testing.T) {

	wired.Go(func(scope wired.Scope) {

		scope.Register(newEmptyStruct)

		oneValue := scope.Construct(newOneValueStructWithValue).(*oneValueStruct)
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
	wired.AutoWire
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
	wired.AutoWire
	Method OneMethod
}

func newUseOneMethod() *useOneMethod {
	return &useOneMethod{}
}

func TestConstructInterface(t *testing.T) {

	wired.Go(func(scope wired.Scope) {
		scope.Register(newDeeper)
		scope.Register(newOneMethod)
		c := scope.Construct(newUseOneMethod).(*useOneMethod)
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

	wired.Go(func(scope wired.Scope) {
		scope.Register(newFirstIncrementer)         // construct Incrementer
		scope.Register(newSecondIncrementer)        // construct Incrementer
		scope.Register(newAnotherSecondIncrementer) // construct Incrementer
		scope.Register(newCombinedIncrementer)      // construct CombinedIncrementer

		incrementer := scope.ConstructByType(IncrementerType).(Incrementer)
		if i := incrementer.Increment(0); i != 13 {
			t.Error("first incrementer should increase to 7, not", i)
		}

		combined := scope.ConstructByType(CombinedIncrementerType).(Incrementer)
		if i := combined.Increment(0); i != 33 {
			t.Error("expected all incrementers to be called which would result in 33 instead of", i)
		}
	})

	wired.Go(func(scope wired.Scope) {
		scope.Register(newFirstIncrementer)    // construct Incrementer
		scope.Register(newCombinedIncrementer) // construct CombinedIncrementer

		combined := scope.ConstructByType(CombinedIncrementerType).(Incrementer)
		if i := combined.Increment(0); i != 7 {
			t.Error("expected all incrementers to be called which would result in 7 instead of", i)
		}
	})
}

func TestConstructWithMultipleScopes(t *testing.T) {
	wired.Go(func(scope wired.Scope) {
		scope.Register(newFirstIncrementer)    // construct Incrementer
		scope.Register(newCombinedIncrementer) // construct CombinedIncrementer

		scope.Go(func(inner wired.Scope) {
			inner.Register(newSecondIncrementer)

			incrementer := inner.ConstructByType(IncrementerType).(Incrementer)
			if i := incrementer.Increment(0); i != 13 {
				t.Error("first incrementer should increase to 13, not", i)
			}

			combined := inner.ConstructByType(CombinedIncrementerType).(Incrementer)
			if i := combined.Increment(0); i != 20 {
				t.Error("expected all incrementers to be called which would result in 20 instead of", i)
			}

		})

		incrementer := scope.ConstructByType(IncrementerType).(Incrementer)
		if i := incrementer.Increment(0); i != 7 {
			t.Error("first incrementer should increase to 7, not", i)
		}

		combined := scope.ConstructByType(CombinedIncrementerType).(Incrementer)
		if i := combined.Increment(0); i != 7 {
			t.Error("expected all incrementers to be called which would result in 7 instead of", i)
		}

	})
}

func TestConstructWithoutReturnValue(t *testing.T) {
	wired.Go(func(scope wired.Scope) {

		scope.Register(newEmptyStruct)

		scope.Construct(func(e *emptyStruct) {
			if e == nil {
				t.Error("expected an empty struct, not nil")
			}
		})

		scope.Inject(func(e *emptyStruct) {
			if e == nil {
				t.Error("expected an empty struct, not nil")
			}
		})

	})
}

// ------ Test map construction for the same type ----

type driverWithInvalidKeyField struct {
}

func (driverWithInvalidKeyField *driverWithInvalidKeyField) Key() {
}

func newInvalidDriver() *driverWithInvalidKeyField {
	return &driverWithInvalidKeyField{}
}

type driver struct {
	name string
}

func (driver *driver) Key() string {
	return driver.name
}

func newAWSDriver() *driver {
	return &driver{name: "aws"}
}

func newAzureDriver() *driver {
	return &driver{name: "azure"}
}

func newGoogleDriver() *driver {
	return &driver{name: "google"}
}

func TestConstructMap(t *testing.T) {

	wired.Go(func(scope wired.Scope) {
		scope.Register(newAWSDriver)
		scope.Register(newAzureDriver)
		scope.Register(newGoogleDriver)

		scope.Inject(func(drivers map[string]*driver) {

			for _, driverName := range []string{"aws", "azure", "google"} {

				if driver, found := drivers[driverName]; found {
					if driver.name != driverName {
						t.Error("found different driver for", driverName)
					}
				} else {
					t.Error("expected driver", driverName)
				}
			}

		})
	})
}

func TestConstructMapWithInvalidKeysShouldPanic(t *testing.T) {

	defer internal.ShouldPanic(t)()

	wired.Go(func(scope wired.Scope) {
		scope.Register(newInvalidDriver)

		scope.Inject(func(drivers map[string]*driverWithInvalidKeyField) {
			// will panic before even getting here
		})

	})
}

// ------ Test map and list construction with mocking ----

type mockedOrNot interface {
	IsMock() bool
}

type isTheRealThing struct {
}

func (isTheRealThing *isTheRealThing) IsMock() bool {
	return false
}

func newRealThing() mockedOrNot {
	return &isTheRealThing{}
}

type isFakeButUseful struct {
}

func (isFakeButUseful *isFakeButUseful) IsMock() bool {
	return true
}

func newFakeThing() mockedOrNot {
	return &isFakeButUseful{}
}

type usesMockOrNot struct {
	wired.AutoWire

	Thing mockedOrNot
}

func (usesMockOrNot *usesMockOrNot) Key() string {
	return "MockOrNot"
}

func newMockOrNot() *usesMockOrNot {
	return &usesMockOrNot{}
}

func TestShouldUsesMockInSlicesAndMaps(t *testing.T) {

	wired.Go(func(scope wired.Scope) {
		scope.Register(newRealThing)
		scope.Register(newMockOrNot)

		scope.Inject(func(check *usesMockOrNot) {
			if check.Thing.IsMock() {
				t.Error("did expect the real thing")
			}
		})

		scope.Go(func(inner wired.Scope) {
			inner.Register(newFakeThing)

			inner.Inject(func(check *usesMockOrNot) {
				if !check.Thing.IsMock() {
					t.Error("did expect the fake thing")
				}
			})

			inner.Inject(func(check []*usesMockOrNot) {
				if !check[0].Thing.IsMock() {
					t.Error("did expect the fake thing in slice")
				}
			})

			inner.Inject(func(check map[string]*usesMockOrNot) {
				if !check["MockOrNot"].Thing.IsMock() {
					t.Error("did expect the fake thing in map")
				}
			})

		})

	})
}
