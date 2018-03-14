package wired_test

import (
	"testing"

	"github.com/okke/wired"
)

type singletonStruct struct {
	wired.Singleton

	count int
}

type structWithSingleton struct {
	wired.AutoWire

	S *singletonStruct
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

func TestConstructSingleton(t *testing.T) {

	structCounter = 0
	wired.Go(func(scope wired.Scope) {
		scope.Register(newSingletonStruct)
		first := scope.Construct(newSingletonStruct)
		second := scope.Construct(newSingletonStruct)

		if first != second {
			t.Error("did not construct singletons", first, second)
		}

		// singleton should be used as constructor argument
		//
		firstUsage := scope.Construct(newStructWithSingleTon).(*structWithSingleton)
		if firstUsage.S != first {
			t.Error("did not use singleton", firstUsage.S, first)
		}

		// singleton should be used when a field is wired into a struct
		//
		secondUsage := scope.Construct(newStructWithSingleTonWithoutArguments).(*structWithSingleton)
		if secondUsage.S != first {
			t.Error("did not use singleton", secondUsage.S, first)
		}
	})
}

func TestConstructSingletonsWithMultipleScopes(t *testing.T) {

	structCounter = 0
	wired.Go(func(scope wired.Scope) {
		scope.Register(newSingletonStruct)

		var innerUsage *structWithSingleton = nil

		scope.Go(func(inner wired.Scope) {
			innerUsage = inner.Construct(newStructWithSingleTon).(*structWithSingleton)
		})

		outerUsage := scope.Construct(newStructWithSingleTon).(*structWithSingleton)

		if innerUsage.S == outerUsage.S {
			t.Error("singletons should be scoped by default", innerUsage.S, outerUsage.S)
		}

	})

	wired.Go(func(scope wired.Scope) {
		scope.Register(newSingletonStruct)

		var innerUsage *structWithSingleton = nil

		outerUsage := scope.Construct(newStructWithSingleTon).(*structWithSingleton)

		scope.Go(func(inner wired.Scope) {
			innerUsage = inner.Construct(newStructWithSingleTon).(*structWithSingleton)
		})

		if innerUsage.S != outerUsage.S {
			t.Error("outer was created first so should be the same as inner but isn't", innerUsage.S, outerUsage.S)
		}

	})
}
