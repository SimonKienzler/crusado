package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the crusado version",
		Args:  cobra.NoArgs,
		Run:   Version,
	}
)

var CrusadoVersion = "v0.0.0"

func Version(cmd *cobra.Command, args []string) {
	fmt.Println(CrusadoVersion)
}
