package config

import (
	"fmt"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
	"github.com/wendisx/puzzle/pkg/util"
)

const (
	_default_config_path = "" // just panic
	DATAKEY_CONFIG       = "_data_config_"
)

/*
config -- don't rely on any default configurationw.
*/

type Config struct {
	DBConfig     DBConfig     `yaml:"database"`
	ServerConfig ServerConfig `yaml:"server"`
	EnvConfig    []string     `yaml:"environment"`
}

func Load(path string) *Config {
	c := &Config{
		DBConfig:     initDBConfig(),
		ServerConfig: initServerConfig(),
	}
	if err := util.ParseYamlFile(path, c); err != nil {
		clog.Panic(fmt.Sprintf("%s", err.Error()))
		return c
	}
	// put Config into data dict
	configDict := NewDataDict[any](DICTKEY_CONFIG)
	configDict.Record(DATAKEY_CONFIG, c)
	PutDict(configDict.Name(), configDict)
	clog.Info(fmt.Sprintf("put data_key(%s) into dict_key(%s)", palette.SkyBlue(DATAKEY_CONFIG), palette.SkyBlue(DICTKEY_CONFIG)))
	return c
}

func GetConfig() *Config {
	configDict := GetDict(DICTKEY_CONFIG)
	c, ok := configDict.Find(DATAKEY_CONFIG).Value().(*Config)
	if !ok {
		clog.Panic(fmt.Sprintf("from dict_key(%s) can't find data_key(%s)", palette.SkyBlue(DICTKEY_CONFIG), palette.SkyBlue(DATAKEY_CONFIG)))
	}
	return c
}
