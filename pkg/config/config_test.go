package config

import (
	"fmt"
	"testing"

	"github.com/wendisx/puzzle/pkg/clog"
)

// test basic configuration [passed]
func Test_basic_print(t *testing.T) {
	configPath := "../../demo/dev.yaml"
	_, _ = Load(configPath)
	print_config()
}

func print_config() {
	c := GetConfig()
	clog.Info(fmt.Sprintf("%#v", c))
}

// test panic GetConfig() [passed]
func Test_get_panic(t *testing.T) {
	_ = GetConfig() // should panic here
}
