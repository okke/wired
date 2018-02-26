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
	wired.Autowire

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

	wired.WithWire(func(wire wired.WireContext) {
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
