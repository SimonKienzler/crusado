package profile

import (
	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/spf13/cobra"
	"github.com/thediveo/klo"
)

var (
	RootCmd = &cobra.Command{
		Use:   "profile",
		Short: "Bundles subcommands that manage crusado profiles",
		Long:  `Use the profile subcommands to list and chose crusado profiles`,
		Args:  cobra.NoArgs,
		Run:   nil,
	}
)

var (
	outputFlag string
)

func init() {
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(SetDefaultCmd)
}

func getPrinter(outputFormat string) (klo.ValuePrinter, error) {
	return klo.PrinterFromFlag(outputFormat, &config.ProfilePrinterSpecs)
}
