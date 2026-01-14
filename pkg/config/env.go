package config

import (
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
)

// LoadEnv load all environment variables matching the specified prefix into the
// specified structure, which must be a pointer type. An error log will be displayed
// if a structure configuration problem exists, but the program will not exit.
// Write all environment variables to the configuration dictionary.
func LoadEnv(prefix string, dest any) {
	vofDest := reflect.ValueOf(dest)
	if vofDest.Kind() != reflect.Ptr || vofDest.Elem().Kind() != reflect.Struct {
		clog.Panic(fmt.Sprintf("invalid environmental container type %s", palette.Red(vofDest.Kind().String())))
	}
	err := envconfig.CheckDisallowed(prefix, dest)
	if err != nil {
		clog.Error(err.Error())
	}
	err = envconfig.Process(prefix, dest)
	if err != nil {
		clog.Error(err.Error())
	}
	// put all environments into data dict
	configDict := GetDict(DICTKEY_CONFIG)
	configDict.Record(DATAKEY_ENV, dest)
}

// GetEnv try to get the pointer to environment structure and will panic
// if not exists the structure.
func GetEnv[ES any]() *ES {
	configDict := GetDict(DICTKEY_CONFIG)
	e, ok := configDict.Find(DATAKEY_ENV).Value().(*ES)
	if !ok {
		clog.Panic(fmt.Sprintf("from data_key(%s) assert to type(*[ES any]) fail", palette.Red(DATAKEY_ENV)))
	}
	return e
}
