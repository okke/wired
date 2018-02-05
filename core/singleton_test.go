package core_test

import (
	"testing"

	"github.com/okke/wires/core"
)

type singletonStruct struct {
	core.Singleton
	count int
}

type structWithSingleton struct {
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
