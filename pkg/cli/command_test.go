package cli

import (
	"fmt"
	"testing"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
)

// test basic load cmd [passed]
func Test_load_cmd(t *testing.T) {
	// test success [passed]
	config.Load("../../demo/dev.yaml")
	cli := LoadCmd("")
	clog.Info(fmt.Sprintf("%#v", cli))
	// test fail [passed]
	cli = LoadCmd("./command.json")
	clog.Info(fmt.Sprintf("%#v", cli))
}

// test get flags [passed]
func Test_get_flags(t *testing.T) {
	config.Load("../../demo/dev.yaml")
	cli := LoadCmd("")
	// test persistent flags [passed]
	delimiter := ":"
	verb := "server:start"
	clog.Info(fmt.Sprintf("%#v", ParsePersistenFlags(verb, delimiter, cli)))
	// test local flags [passed]
	clog.Info(fmt.Sprintf("%#v", ParseLocalFlags(verb, delimiter, cli)))
	// test diff verb with the same sub verb [passed]
	verb = "client:start"
	clog.Info(fmt.Sprintf("%#v", ParsePersistenFlags(verb, delimiter, cli)))
	clog.Info(fmt.Sprintf("%#v", ParseLocalFlags(verb, delimiter, cli)))
	// expected [2,1,0,3] [passed]
	// not exists verb []
	verb = "log:start"
	clog.Info(fmt.Sprintf("%#v", ParsePersistenFlags(verb, delimiter, cli)))
	clog.Info(fmt.Sprintf("%#v", ParseLocalFlags(verb, delimiter, cli)))
	// test empty verb [passed]
	verb = ""
	clog.Info(fmt.Sprintf("%#v", ParsePersistenFlags(verb, delimiter, cli)))
	clog.Info(fmt.Sprintf("%#v", ParseLocalFlags(verb, delimiter, cli)))
}

// test get cli from data dict  []
func print_cli() {
	cli := GetCmd()
	clog.Info(fmt.Sprintf("%#v", cli))
}

func Test_get_cli(t *testing.T) {
	config.Load("../../demo/dev.yaml")
	_ = LoadCmd("")
	print_cli()
}
