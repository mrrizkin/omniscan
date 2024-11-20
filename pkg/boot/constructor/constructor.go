package constructor

import (
	"reflect"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

type Constructor interface {
	Construct() interface{}
}

// LoadNewConstructors automatically finds and loads all functions starting with "New"
// from Go files in the specified directory and its subdirectories
func Load(modules ...Constructor) fx.Option {
	constructors := make([]interface{}, len(modules))
	for i, module := range modules {
		assert.True(nil, reflect.TypeOf(module.Construct()).Kind() == reflect.Func)
		constructors[i] = module.Construct()
	}

	return fx.Options(
		fx.Provide(constructors...),
	)
}
