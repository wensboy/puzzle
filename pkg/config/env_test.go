package config

import (
	"fmt"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/wendisx/puzzle/pkg/clog"
)

type (
	Env struct {
		Test_1 string
		Test_2 int
		Test_3 string
		Test_4 map[string]int
		Test_5 bool
		Test_6 []string
		Host   string
		Port   int
	}
)

// test usage for envconfig package [passed]
func Test_usage_envconfig(t *testing.T) {
	// 1. process() with prefix [passed]
	prefix := "puzzle"
	var env Env
	err := envconfig.Process(prefix, &env)
	if err != nil {
		clog.Error(err.Error())
	} else {
		clog.Info(fmt.Sprintf("%#v\n", env))
	}
	// 2. test usage() [passed]
	err = envconfig.Usage(prefix, &env)
	if err != nil {
		clog.Error(err.Error())
	} // 3. test checkDisallowed() [passed] (it looks like this is useless...)
	err = envconfig.CheckDisallowed(prefix, &env)
	if err != nil {
		clog.Error(err.Error())
	} else {
		clog.Info(fmt.Sprintf("%#v", env))
	}
}

// test load environments to custom struct
type (
	DevEnv struct {
		prefix       string
		AgentName    string `envconfig:"AGENT_NAME"`
		AgentVersion string `envconfig:"AGENT_VERSION"`
		Host         string `envconfig:"HOST"`
		Port         int    `envconfig:"PORT"`
		Proxy        string `envconfig:"PROXY"`
	}
)

func (de *DevEnv) Prefix() string {
	return de.prefix
}

func print_dev_env() {
	// old test [passed]
	// d := GetDict(DICTKEY_CONFIG)
	// c, ok := d.Find(DATAKEY_ENV).Value().(*DevEnv)
	// if !ok {
	// 	clog.Error(fmt.Sprintf("from dict_key(%s) get data_key(%s) failed", palette.SkyBlue(DICTKEY_CONFIG), palette.Red(DATAKEY_ENV)))
	// }
	// new test [passed]
	c := GetEnv[DevEnv]()
	clog.Info(fmt.Sprintf("%#v", c))
}

func Test_load_env(t *testing.T) {
	var de DevEnv
	de.prefix = "dev"
	path := "../../demo/dev.yaml"
	_ = Load(path)
	LoadEnv(de.Prefix(), &de)
	print_dev_env()
}
