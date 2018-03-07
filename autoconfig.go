package wired

import (
	"reflect"

	"github.com/okke/wired/wtemplate"

	"github.com/okke/wired/internal"
)

// AutoConfig is a tag that drives auto configuration of struct fields
//
type AutoConfig struct {
}

type autoconfig struct {
}

// Configurator defines a method to lookup a configuration value
//
type Configurator interface {
	ConfigValue(key string) string
}

type allConfigs struct {
	all []Configurator
}

func newAllConfigs(all []Configurator) *allConfigs {
	return &allConfigs{all: all}
}

// allConfigs implements wtemplate.Context
//
func (allConfigs *allConfigs) Solve(key string) string {
	for _, config := range allConfigs.all {
		if value := config.ConfigValue(key); value != "" {
			return value
		}
	}
	return ""
}

func init() {
	RegisterStructDecorationTag(reflect.TypeOf((*AutoConfig)(nil)).Elem(), &autoconfig{})
}

func (autoconfig *autoconfig) GetValueFor(wire Scope, obj reflect.Value, field reflect.Value, fieldType reflect.StructField) (reflect.Value, bool) {

	tag := fieldType.Tag.Get("autoconfig")
	if tag != "" {
		config := wire.Construct(newAllConfigs).(*allConfigs)
		if value := internal.ConvertString2Value(fieldType.Type.Kind(), wtemplate.Parse(config, tag)); value != internal.NilValue {
			return value, true
		}
	}

	return internal.NilValue, false
}
