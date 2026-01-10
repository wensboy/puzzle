package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/wendisx/puzzle/pkg/clog"
)

type (
	A struct {
		Env string
		Id  int
	}
	B struct {
		Env      string
		UserName string
	}
)

// test _default_queue
func show(name DictKey) {
	dd := GetDict(name)
	addr, ok := dd.Find("SERVER_ADDR").Value().(string)
	env, ok := dd.Find("SYSTEM_ENV_DEV").Value().(A)
	if !ok {
		clog.Panic("invalid key from data dict")
	}
	clog.Info(fmt.Sprintf("%#v from %s", env, addr))
}

func Test_datadict(t *testing.T) {
	dd := NewDataDict[any]("_system_")
	PutDict(dd.Name(), dd)
	dd.Record("SYSTEM_ENV_DEV", A{Env: "dev", Id: 1 << 10})
	dd.Record("SYSTEM_ENV_PROD", B{Env: "prod", UserName: "puzzler"})
	dd.Record("SERVER_ADDR", "127.0.0.1:3333")
	time.Sleep(1 * time.Second)
	dd.Record("SERVER_ADDR", "0.0.0.0:3333") // update
	time.Sleep(1 * time.Second)
	show(dd.Name())
}

func Test_bad_key(t *testing.T) {
	show("bad_key")
}
