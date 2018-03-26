package cli

import (
	"os"
	"strings"

	"github.com/okke/wired"
)

// Arguments represents flags and (file)names that can be retrieved from
// a line of input
//
type Arguments interface {
	Flag(string) (string, bool)
	NonFlags() []string
}

// ArgumentParser can parse a line of input into an Argument object
//
type ArgumentParser interface {
	Parse(arguments []string) Arguments
}

type parsedArguments struct {
	wired.Singleton

	flags    map[string]string
	nonflags []string
}

func (parsedArguments *parsedArguments) Flag(name string) (string, bool) {
	result, found := parsedArguments.flags[name]
	return result, found
}

func (parsedArguments *parsedArguments) NonFlags() []string {
	return parsedArguments.nonflags
}

type argumentParser struct {
	wired.Singleton
}

func (argumentParser *argumentParser) Parse(arguments []string) Arguments {
	parsed := &parsedArguments{
		flags:    make(map[string]string),
		nonflags: make([]string, 0)}

	var flagName string
	for _, arg := range arguments {
		if flagName == "" {
			if strings.HasPrefix(arg, "-") {
				flagName = strings.TrimLeft(arg, "-")
			} else {
				parsed.nonflags = append(parsed.nonflags, arg)
			}
		} else {
			// found flag/value pair
			parsed.flags[flagName] = arg
			flagName = ""

		}
	}
	return parsed
}

// ArgumentProvider provides arguments that can be parsed
//
type ArgumentProvider interface {
	Arguments() []string
}

type argumentProvider struct {
	wired.Singleton
}

func (argumentProvider *argumentProvider) Arguments() []string {
	return os.Args
}

func newArgumentParser() ArgumentParser {
	return &argumentParser{}
}

func newArgumentProvider() ArgumentProvider {
	return &argumentProvider{}
}

func newArguments(provider ArgumentProvider, parser ArgumentParser) Arguments {
	return parser.Parse(provider.Arguments())
}

// register the argument parser at global scope
//
func init() {
	wired.Global().Register(newArgumentParser)
	wired.Global().Register(newArgumentProvider)
	wired.Global().Register(newArguments)
}
