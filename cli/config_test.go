package cli_test

import (
	"testing"

	"github.com/okke/wired"
	"github.com/okke/wired/cli"
)

type mockProviderForConfig struct {
}

func (mockProviderForConfig *mockProviderForConfig) Arguments() []string {
	return []string{"-pepper", "chipotle"}
}

func newMockProviderForConfig() cli.ArgumentProvider {
	return &mockProviderForConfig{}
}

type needConfig struct {
	wired.AutoConfig

	Unknown string `autoconfig:"${sauce:unknown}"`
	Name    string `autoconfig:"${pepper}"`
}

func newNeedConfig() *needConfig {
	return &needConfig{}
}

func TestConfigByArguments(t *testing.T) {

	wired.Global().Go(func(scope wired.Scope) {
		scope.Register(newMockProviderForConfig)
		scope.Register(newNeedConfig)

		scope.Inject(func(configured *needConfig) {
			if configured == nil {
				t.Fatal("expected configured to be constructed")
			}

			if configured.Name != "chipotle" {
				t.Error("expected name to be chipotle, not", configured.Name)
			}

			if configured.Unknown != "unknown" {
				t.Error("expected unknown to be unknown, not", configured.Unknown)
			}
		})
	})
}
