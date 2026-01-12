package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	_artfont = ` 
┌─┐┬ ┬┌─┐┌─┐┬  ┌─┐
├─┘│ │┌─┘┌─┘│  ├┤ 
┴  └─┘└─┘└─┘┴─┘└─┘
  intro: %s
version: %s
	`
)

var (
	_verb_version = ":version"

	_default_intro   = "-"
	_default_version = "vx.y.z"
)

func mountVersion(rootCmd *cobra.Command) {
	versionCmd := GetCommand(_verb_version, _default_delimiter)
	versionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(os.Stderr, _artfont, _default_intro, _default_version)
		return nil
	}
	rootCmd.AddCommand(versionCmd)
}
