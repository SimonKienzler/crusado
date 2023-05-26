package profile

import (
	"fmt"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/spf13/cobra"
)

var (
	SetDefaultCmd = &cobra.Command{
		Use:   "set-default [profile-name]",
		Short: "Set the default profile for the template subcommands to use",
		Long:  `TODO`,
		Args:  cobra.ExactArgs(1),
		Run:   SetDefault,
	}
)

func SetDefault(cmd *cobra.Command, args []string) {
	_ = config.GetConfig()
	// TODO write the given profile name as the default value in the well-known
	// crusado config file, if it is a valid profile with a retrievable file path
	fmt.Println("this is not yet implemented")
}
