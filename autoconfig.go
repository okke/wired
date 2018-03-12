package wired

import (
	"os"
	"reflect"
	"strings"

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

type configByEnvironment struct {
}

func (configByEnvironment *configByEnvironment) ConfigValue(key string) string {
	return os.Getenv(strings.ToUpper(strings.Replace(key, ".", "_", -1)))
}

func newConfigByEnvironment() Configurator {
	return &configByEnvironment{}
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
	Global().Register(newConfigByEnvironment)
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
