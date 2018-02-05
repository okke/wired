package core_test

import (
	"fmt"
	"testing"

	"github.com/okke/wires/core"
)

type chipotle struct {
	name string
}

func (chipotle *chipotle) String() string {
	return fmt.Sprintf("chipotle(name:%s)", chipotle.name)
}

func newChipotle() *chipotle {
	return &chipotle{name: "chipotle"}
}

type jalapeno struct {
	name string
}

func newJalapeno() *jalapeno {
	return &jalapeno{name: "jalapeno"}
}

func (jalapeno *jalapeno) String() string {
	return fmt.Sprintf("jalapeno(name:%s)", jalapeno.name)
}

type habanero struct {
	name       string
	wantsToBe  *jalapeno
	RealyHates *chipotle
}

func (habanero *habanero) GetWantsToBe() *jalapeno {
	return habanero.wantsToBe
}

func (habanero *habanero) SetWantsToBe(j *jalapeno) {
	habanero.wantsToBe = j
}

func (habanero *habanero) String() string {
	return fmt.Sprintf("habanero(name:%s, wantsToBe:%v, RealyHates:%v)", habanero.name, habanero.wantsToBe, habanero.RealyHates)
}

func newHabanero() *habanero {
	return &habanero{name: "hotty pepper"}
}

type peppers struct {
	C     *chipotle
	j     *jalapeno
	H     *habanero
	h     *habanero
	price int
	Name  string
}

func (peppers *peppers) SetJ(j *jalapeno) {
	peppers.j = j
}

func (peppers *peppers) SetH(h *habanero) {
	peppers.h = h
}

func (peppers *peppers) SetPrice(price int) {
	peppers.price = price
}

func (peppers *peppers) GetJ() *jalapeno {
	return peppers.j
}

func (peppers *peppers) GetH() *habanero {
	return peppers.h
}

func newPeppers() *peppers {
	return &peppers{Name: "PerfectForSoup", H: &habanero{name: "habanero"}, h: &habanero{name: "small habanero", wantsToBe: &jalapeno{name: "a red jalapeno"}}}
}

func TestNewFieldDecorator(t *testing.T) {
	core.WithWire(func(wire core.WireContext) {
		wire.Register(newChipotle)
		wire.Register(newJalapeno)
		wire.Register(newHabanero)

		p := wire.Construct(newPeppers).(*peppers)
		if p == nil {
			t.Error("should have constructed some peppers")
		}

		// chipotle can be set through public field
		//
		if p.C == nil || p.C.name != "chipotle" {
			t.Error("expected a chipotle to be constructed, not", p.C)
		}

		// jalapeno can be set through setter
		//
		if p.j == nil || p.j.name != "jalapeno" {
			t.Error("expected a jalapeno to be constructed, not", p.j)
		}

		// habanero won't be set because it already has a value that
		// could be retrieved through its public field
		//
		if p.H == nil || p.H.name != "habanero" {
			t.Error("expected a habanero to be constructed, not", p.H)
		}

		// small habanero won't be set because it already has a value
		// that could be retrieved through a getter
		//
		if p.h == nil || p.h.name != "small habanero" {
			t.Error("expected a habanero to be constructed, not", p.h)
		}

		// but habanero should be decorated so unknown fields are wired
		//
		if p.h == nil || p.h.RealyHates == nil {
			t.Error("expected a habanero to be constructed that hates chipotles", p.h)
		}
		if p.h == nil || p.h.wantsToBe == nil {
			t.Error("expected a habanero to be constructed that wants to be a jalapeno", p.h)
		}

		if p.Name != "PerfectForSoup" {
			t.Error("expected to be named PerfectForSoup")
		}
	})
}
