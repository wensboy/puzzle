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
}

func Load(path string) *Config {
	c := &Config{
		DBConfig:     initDBConfig(),
		ServerConfig: initServerConfig(),
	}
	if err := util.ParseYamlFile(path, c); err != nil {
		clog.Error(fmt.Sprintf("<pkg.config> %s", err.Error()))
	}
	// put Config into data dict
	configDict := NewDataDict[any](DICTKEY_CONFIG)
	configDict.Record(DATAKEY_CONFIG, c)
	PutDict(configDict.Name(), configDict)
	return c
}

func GetConfig() *Config {
	configDict := GetDict(DICTKEY_CONFIG)
	c, ok := configDict.Find(DATAKEY_CONFIG).Value().(*Config)
	if !ok {
		clog.Error(fmt.Sprintf("from dict(%s) find data_key(%s)", palette.SkyBlue(DICTKEY_CONFIG), palette.SkyBlue(DATAKEY_CONFIG)))
	}
	return c
}
