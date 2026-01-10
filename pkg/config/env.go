package config

import (
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
)

const (
	DATAKEY_ENV = "_data_env_"
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
	c := GetDict(DICTKEY_CONFIG)
	c.Record(DATAKEY_ENV, dest)
	clog.Info(fmt.Sprintf("<pkg.config> put data_key(%s) into dict_key(%s)", palette.SkyBlue(DATAKEY_ENV), palette.SkyBlue(DICTKEY_CONFIG)))
}
