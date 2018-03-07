package wired_test

import (
	"testing"

	"github.com/okke/wired"
)

type needConfig struct {
	wired.AutoConfig

	Unknown  string
	Name     string `autoconfig:"${pepper}"`
	NotFound string `autoconfig:"${unknown}"`
	Default  string `autoconfig:"${unknown:default}"`

	NumberoInt   int   `autoconfig:"-8"`
	NumberoInt8  int8  `autoconfig:"-8"`
	NumberoInt16 int16 `autoconfig:"-16"`
	NumberoInt32 int32 `autoconfig:"-32"`
	NumberoInt64 int64 `autoconfig:"-64"`

	NumberoUint   uint   `autoconfig:"8"`
	NumberoUint8  uint8  `autoconfig:"8"`
	NumberoUint16 uint16 `autoconfig:"16"`
	NumberoUint32 uint32 `autoconfig:"32"`
	NumberoUint64 uint64 `autoconfig:"64"`

	NumberoFloat32 float32 `autoconfig:"32.0"`
	NumberoFloat64 float64 `autoconfig:"64.0"`

	DefinitelyTrue bool `autoconfig:"true"`
}

func NewNeedConfig() *needConfig {
	return &needConfig{}
}

type testConfig struct {
	config map[string]string
}

func (testConfig *testConfig) ConfigValue(key string) string {
	return testConfig.config[key]
}

func newTestConfig() wired.Configurator {
	config := make(map[string]string, 0)
	config["pepper"] = "habanero"
	return &testConfig{config: config}
}

func TestSimpleStringConfig(t *testing.T) {
	wired.Go(func(scope wired.Scope) {

		scope.Register(newTestConfig)

		need := scope.Construct(NewNeedConfig).(*needConfig)

		if need.Name != "habanero" {
			t.Error("expected habanero, not", need.Name)
		}

		if need.NotFound != "" {
			t.Error("did not expect a value but got", need.NotFound)
		}

		if need.Default != "default" {
			t.Error("expected default, not", need.Default)
		}

		if need.NumberoInt != -8 {
			t.Error("expected number -8, not", need.NumberoInt)
		}

		if need.NumberoInt8 != -8 {
			t.Error("expected number -8, not", need.NumberoInt8)
		}

		if need.NumberoInt16 != -16 {
			t.Error("expected number -16, not", need.NumberoInt16)
		}

		if need.NumberoInt32 != -32 {
			t.Error("expected number -32, not", need.NumberoInt32)
		}

		if need.NumberoInt64 != -64 {
			t.Error("expected number -64, not", need.NumberoInt64)
		}

		if need.NumberoUint != 8 {
			t.Error("expected number 8, not", need.NumberoUint)
		}

		if need.NumberoUint8 != 8 {
			t.Error("expected number 8, not", need.NumberoUint8)
		}

		if need.NumberoUint16 != 16 {
			t.Error("expected number 16, not", need.NumberoUint16)
		}

		if need.NumberoUint32 != 32 {
			t.Error("expected number 32, not", need.NumberoUint32)
		}

		if need.NumberoUint64 != 64 {
			t.Error("expected number 64, not", need.NumberoUint64)
		}

		if need.NumberoFloat32 != 32.0 {
			t.Error("expected number 32.0, not", need.NumberoFloat32)
		}

		if need.NumberoFloat64 != 64.0 {
			t.Error("expected number 64.0, not", need.NumberoFloat64)
		}

		if !need.DefinitelyTrue {
			t.Error("expected true")
		}
	})
}

func TestAutoConfigWithoutCongurators(t *testing.T) {
	wired.Go(func(scope wired.Scope) {
		need := scope.Construct(NewNeedConfig).(*needConfig)
		if need == nil {
			t.Error("expected to construct an object")
		}
	})
}
