package config

import (
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
)

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

func GetEnv[ES any]() *ES {
	configDict := GetDict(DICTKEY_CONFIG)
	e, ok := configDict.Find(DATAKEY_ENV).Value().(*ES)
	if !ok {
		clog.Panic(fmt.Sprintf("from data_key(%s) assert to type(*[ES any]) fail", palette.Red(DATAKEY_ENV)))
	}
	return e
}
