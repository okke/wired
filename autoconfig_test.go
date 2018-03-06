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
