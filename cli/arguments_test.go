package cli_test

import (
	"testing"

	"github.com/okke/wired"
	"github.com/okke/wired/cli"
)

func TestParserShouldAcceptNonflags(t *testing.T) {

	wired.Global().Go(func(scope wired.Scope) {
		scope.Inject(func(parser cli.ArgumentParser) {
			arguments := parser.Parse([]string{"uno", "dos"})

			if len(arguments.NonFlags()) != 2 {
				t.Fatal("expected 2 non flag arguments, not", len(arguments.NonFlags()))
			}

			if arguments.NonFlags()[0] != "uno" {
				t.Error("expected first non flag to be uno, not", arguments.NonFlags()[0])
			}

			if arguments.NonFlags()[1] != "dos" {
				t.Error("expected second non flag to be dos, not", arguments.NonFlags()[1])
			}

		})
	})
}

func TestParserShouldAcceptFlags(t *testing.T) {

	wired.Global().Go(func(scope wired.Scope) {
		scope.Inject(func(parser cli.ArgumentParser) {

			arguments := parser.Parse([]string{"-first", "uno", "--second", "dos"})
			if len(arguments.NonFlags()) != 0 {
				t.Error("expected no non flag arguments, not", len(arguments.NonFlags()))
			}

			if first, found := arguments.Flag("first"); !found {
				t.Error("expected first flag to exist")
			} else {
				if first != "uno" {
					t.Error("expected first to be uno, not", first)
				}
			}

			if second, found := arguments.Flag("second"); !found {
				t.Error("expected second flag to exist")
			} else {
				if second != "dos" {
					t.Error("expected second to be uno, not", second)
				}
			}

		})
	})
}

func TestParserShouldAcceptAMixOfFlagsAndNonFlags(t *testing.T) {

	wired.Global().Go(func(scope wired.Scope) {
		scope.Inject(func(parser cli.ArgumentParser) {

			arguments := parser.Parse([]string{"command", "-first", "uno", "--second", "dos", "chipotle"})
			if len(arguments.NonFlags()) != 2 {
				t.Error("expected two non flag arguments, not", len(arguments.NonFlags()))
			}

			if first, found := arguments.Flag("first"); !found {
				t.Error("expected first flag to exist")
			} else {
				if first != "uno" {
					t.Error("expected first to be uno, not", first)
				}
			}

			if second, found := arguments.Flag("second"); !found {
				t.Error("expected second flag to exist")
			} else {
				if second != "dos" {
					t.Error("expected second to be uno, not", second)
				}
			}

			if arguments.NonFlags()[0] != "command" {
				t.Error("expected first non flag to be command, not", arguments.NonFlags()[0])
			}

			if arguments.NonFlags()[1] != "chipotle" {
				t.Error("expected second non flag to be chipotle, not", arguments.NonFlags()[1])
			}

		})
	})
}

func TestArgumentProviderShouldExist(t *testing.T) {

	wired.Global().Go(func(scope wired.Scope) {
		scope.Inject(func(provider cli.ArgumentProvider, parser cli.ArgumentParser) {
			if provider == nil {
				t.Fatal("expected an argument provider")
			}
			parser.Parse(provider.Arguments())
		})
	})
}

type mockProvider struct {
}

func (mockProvider *mockProvider) Arguments() []string {
	return []string{"knock"}
}

func newArgumentProdiverMock() cli.ArgumentProvider {
	return &mockProvider{}
}
func TestArgumentProviderShouldBeOveridable(t *testing.T) {

	wired.Global().Go(func(scope wired.Scope) {
		scope.Register(newArgumentProdiverMock)

		scope.Inject(func(provider cli.ArgumentProvider, parser cli.ArgumentParser) {
			if provider == nil {
				t.Fatal("expected an argument provider")
			}
			args := provider.Arguments()

			if args[0] != "knock" {
				t.Error("expected first argument to be knock, not", args[0])
			}
		})
	})
}

func TestArgumentsShouldBeSimplyThere(t *testing.T) {
	wired.Global().Go(func(scope wired.Scope) {
		scope.Inject(func(arguments cli.Arguments) {
			if arguments == nil {
				t.Fatal("expected arguments")
			}
		})
	})
}

func TestArgumentsShouldBeSimplyThereAlsoWhenMocked(t *testing.T) {
	wired.Global().Go(func(scope wired.Scope) {
		scope.Register(newArgumentProdiverMock)

		scope.Inject(func(arguments cli.Arguments) {
			if arguments == nil {
				t.Fatal("expected arguments")
			}

			if arguments.NonFlags()[0] != "knock" {
				t.Error("expected first argument to be knock, not", arguments.NonFlags()[0])
			}
		})
	})
}
