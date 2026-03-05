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

// MountVersion mount the verb `-version` to the command tree.
// For details, see the `command.json` file in the project structure to
// find the structure of the corresponding command.
func MountVersion(rootCmd *cobra.Command) {
	versionCmd := GetCommand(_verb_version, _default_delimiter)
	versionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(os.Stderr, _artfont, _default_intro, _default_version)
		return nil
	}
	rootCmd.AddCommand(versionCmd)
}
