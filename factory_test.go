package wired_test

import (
	"reflect"
	"testing"

	"github.com/okke/wired"
)

type pepperFactory struct {
	wired.Factory

	count int
}

type pepperFromFactory struct {
	nr int
}

var pepperType = reflect.TypeOf((*pepperFromFactory)(nil))

func (pepperFactory *pepperFactory) Construct() *pepperFromFactory {
	pepperFactory.count = pepperFactory.count + 1
	return &pepperFromFactory{nr: pepperFactory.count}
}

func newPepperFactory() *pepperFactory {
	return &pepperFactory{count: 0}
}

func TestFactoryConstruction(t *testing.T) {
	wired.Go(func(scope wired.Scope) {
		scope.Register(newPepperFactory)

		for walk := 1; walk < 10; walk++ {
			pepper := scope.ConstructByType(pepperType).(*pepperFromFactory)
			if pepper.nr != walk {
				t.Error("expected pepper nr to be ", walk, "instead of", pepper.nr)
			}
		}

	})
}
