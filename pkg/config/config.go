package config

import (
	"fmt"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/util"
)

const (
	DefaultConfigPath = "./config/example.yaml"

	DefaultDsn = "puzzle:af271f8c-9c90-40c9-80f9-9665caedda18@tcp(localhost:63306)/puzzle?"

	DefaultServerHost = "0.0.0.0"
	DefaultServerPort = 3000
)

var globalConfig *ConfigHub

type ConfigHub struct {
	DBConfig     DBConfig     `yaml:"database"`
	ServerConfig ServerConfig `yaml:"server"`
}

func initConfigHub() {
	globalConfig = &ConfigHub{
		DBConfig:     initDBConfig(""),
		ServerConfig: initServerConfig("", 3000),
	}
}

func SetupConfigHub(cover bool, path string) error {
	initConfigHub()
	if err := util.ParseYamlFile(path, globalConfig); err != nil {
		clog.Error(fmt.Sprintf("parse config file fail for %s", err.Error()))
		return err
	}
	if cover {
		// TODO: cover config logic
	}
	return nil
}

func GetDBConfig() *DBConfig {
	return &globalConfig.DBConfig
}

func GetServerConfig() *ServerConfig {
	return &globalConfig.ServerConfig
}
