package cli

import (
	"github.com/spf13/cobra"
	"github.com/wendisx/puzzle/pkg/clog"
)

var (
	_exists_root       = false
	_verb_root         = "root"
	_default_delimiter = ":"
)

func RootVerb(verb string) {
	_verb_root = verb
}

func DefaultDelimiter(delimiter string) {
	_default_delimiter = delimiter
}

func mountRoot() *cobra.Command {
	versionCmd := mountVersion()
	var rootCmd *cobra.Command
	if !_exists_root {
		rootCmd = &cobra.Command{
			Use:   "puzzle",
			Short: "",
			Long:  "",
		}
		_dict_command.Record(_verb_root, rootCmd)
	} else {
		rootCmd = GetCommand(_verb_root, _default_delimiter)
	}
	rootCmd.AddCommand(versionCmd)
	return rootCmd
}

func Execute() {
	rootCmd := mountRoot()
	if err := rootCmd.Execute(); err != nil {
		clog.Fatal(err.Error())
	}
}
