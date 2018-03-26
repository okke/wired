package cli

import (
	"github.com/okke/wired"
)

type configByArguments struct {
	wired.AutoWire
	wired.Singleton
	ParsedArguments Arguments
}

func (configByArguments *configByArguments) ConfigValue(key string) string {
	if value, found := configByArguments.ParsedArguments.Flag(key); found {
		return value
	}
	return ""
}

func NewConfigByFlags() wired.Configurator {
	return &configByArguments{}
}

func init() {
	wired.Global().Register(NewConfigByFlags)
}
