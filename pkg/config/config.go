package config

import (
	"fmt"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
	"github.com/wendisx/puzzle/pkg/util"
)

var (
	_default_config = "./vendor/puzzle/demo/example.yaml"
)

/*
config -- don't rely on any default configurationw.
*/

type Config struct {
	DBConfig     DBConfig     `yaml:"database"`
	ServerConfig ServerConfig `yaml:"server"`
	EnvConfig    []string     `yaml:"environment"`
}

func DefaultConfig(path string) {
	_default_config = path
}

func LoadConfig(path string) *Config {
	c := &Config{
		DBConfig:     initDBConfig(),
		ServerConfig: initServerConfig(),
	}
	if path == "" {
		path = _default_config
	}
	if err := util.ParseYamlFile(path, c); err != nil {
		clog.Panic(fmt.Sprintf("%s", err.Error()))
		return c
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
		clog.Panic(fmt.Sprintf("from data_key(%s) assert to type(*Config) fail", palette.Red(DATAKEY_CONFIG)))
	}
	return c
}
