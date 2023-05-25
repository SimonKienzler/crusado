package profile

import (
	"log"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/spf13/cobra"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "List crusado profiles",
		Long:  `Allows you to list available profiles. You can specify an output format.`,
		Args:  cobra.NoArgs,
		Run:   List,
	}
)

func init() {
	ListCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "define the output format: [wide, yaml, json, jsonpath]")
}

func List(cmd *cobra.Command, args []string) {
	// TODO read the list of availble profiles from a well-known source, e.g. $HOME/.crusado/profiles.yaml
	profileConfigs := []config.ProfileConfig{
		{
			Name:     "example-profile",
			FilePath: "/somewhere/on/disk/profile.yaml",
		},
		{
			Name:     "other-profile",
			FilePath: "/path/to/other-profile.yaml",
		},
	}

	printer, err := getPrinter(outputFlag)
	if err != nil {
		log.Fatalf("Could not get printer: %q", err)
	}

	printer.Fprint(os.Stdout, profileConfigs)
	if err != nil {
		log.Fatalf("Could not print profiles: %q", err)
	}
}
