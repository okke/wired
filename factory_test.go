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

type sauceFactory struct {
	wired.Factory
}

type pepperFromFactory struct {
	nr int
}

type sauceFromFactory struct {
	pepper *pepperFromFactory
}

var pepperType = reflect.TypeOf((*pepperFromFactory)(nil))
var sauceType = reflect.TypeOf((*sauceFromFactory)(nil))

func (pepperFactory *pepperFactory) Construct() *pepperFromFactory {
	pepperFactory.count = pepperFactory.count + 1
	return &pepperFromFactory{nr: pepperFactory.count}
}

func (sauceFactory *sauceFactory) Construct(pepper *pepperFromFactory) *sauceFromFactory {
	return &sauceFromFactory{pepper: pepper}
}

func newPepperFactory() *pepperFactory {
	return &pepperFactory{count: 0}
}

func newSauceFactory() *sauceFactory {
	return &sauceFactory{}
}

func TestFactoryConstruction(t *testing.T) {
	wired.Go(func(scope wired.Scope) {
		scope.Register(newPepperFactory)
		scope.Register(newSauceFactory)

		for walk := 1; walk < 10; walk++ {
			pepper := scope.ConstructByType(pepperType).(*pepperFromFactory)
			if pepper.nr != walk {
				t.Error("expected pepper nr to be ", walk, "instead of", pepper.nr)
			}
		}

		sauce := scope.ConstructByType(sauceType).(*sauceFromFactory)
		if sauce.pepper == nil {
			t.Error("expected sauce with pepper")
		} else {
			if sauce.pepper.nr != 10 {
				t.Error("expected pepper nr to be 10 not", sauce.pepper.nr)
			}
		}

	})
}
