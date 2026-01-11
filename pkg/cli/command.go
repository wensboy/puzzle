package cli

import (
	"fmt"
	"strings"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
	"github.com/wendisx/puzzle/pkg/palette"
	"github.com/wendisx/puzzle/pkg/util"
)

const (
	_default_command_path = "../../command.json"
)

type (
	Flag struct {
		FullName  string `json:"fullName"`
		ShortName string `json:"shortName"`
		Desc      string `json:"desc"`
	}
	Command struct {
		Verb      string `json:"verb"`
		ShortDesc string `json:"shortDesc"`
		LongDesc  string `json:"longDesc"`
		// localFlags Just collect the names and
		// aliases of local parameters and usage descriptionsw.
		PersistentFlags []Flag    `json:"persistentFlags"`
		LocalFlags      []Flag    `json:"localFlags"`
		SubCommand      []Command `json:"subCommands"`
	}
	Cli struct {
		App      string    `json:"app"`
		Entry    []string  `json:"entry"`
		Version  string    `json:"version"`
		Commands []Command `json:"commands"`
	}
)

func LoadCmd(path string) *Cli {
	if path == "" {
		path = _default_command_path
	}
	var cli Cli
	if err := util.ParseJsonFile(path, &cli); err != nil {
		clog.Panic(err.Error())
		return nil
	}
	// put cli into data dict
	configDict := config.GetDict(config.DICTKEY_CONFIG)
	configDict.Record(config.DATAKEY_CLI, &cli)
	return &cli
}

func GetCmd() *Cli {
	configDict := config.GetDict(config.DICTKEY_CONFIG)
	cli, ok := configDict.Find(config.DATAKEY_CLI).Value().(*Cli)
	if !ok {
		clog.Panic(fmt.Sprintf("from data_key(%s) assert to type(*CLi) fail", palette.Red(config.DATAKEY_CLI)))
	}
	return cli
}

func FindCommand(verb string, delimiter string, cmd Command) (Command, bool) {
	if verb == cmd.Verb {
		return cmd, true
	}
	// verb should shrink prefix
	idx := strings.Index(verb, delimiter)
	if idx != -1 && string(verb[:idx]) == cmd.Verb {
		verb = verb[idx+1:]
	}
	var c Command
	find := false
	for i := range cmd.SubCommand {
		c, find = FindCommand(verb, delimiter, cmd.SubCommand[i])
		if find {
			break
		}
	}
	return c, find
}

func ParsePersistenFlags(verb string, delimiter string, entry *Cli) []Flag {
	if entry == nil {
		clog.Error(fmt.Sprintf("from invalid cli entry find verb(%s) persistent flags", palette.Red(verb)))
		return nil
	}
	for i := range entry.Commands {
		if cmd, find := FindCommand(verb, delimiter, entry.Commands[i]); find {
			clog.Info(fmt.Sprintf("parse verb(%s) persistent +[%s] Flags", palette.SkyBlue(verb), palette.Green(len(cmd.PersistentFlags))))
			return cmd.PersistentFlags
		}
	}
	clog.Warn(fmt.Sprintf("not exists verb(%s)", palette.Red(verb)))
	return nil
}

func ParseLocalFlags(verb string, delimiter string, entry *Cli) []Flag {
	if entry == nil {
		clog.Error(fmt.Sprintf("from invalid cli entry find verb(%s) local flags", palette.Red(verb)))
		return nil
	}
	for i := range entry.Commands {
		if cmd, find := FindCommand(verb, delimiter, entry.Commands[i]); find {
			clog.Info(fmt.Sprintf("parse verb(%s) local +[%s] Flags", palette.SkyBlue(verb), palette.Green(len(cmd.LocalFlags))))
			return cmd.LocalFlags
		}
	}
	clog.Warn(fmt.Sprintf("not exists verb(%s)", palette.Red(verb)))
	return nil
}
