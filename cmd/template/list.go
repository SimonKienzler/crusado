package template

import (
	"log"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"

	"github.com/spf13/cobra"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "List crusado templates",
		Long:  `Allows you to list available user story templates. You can specify an output format.`,
		Args:  cobra.NoArgs,
		Run:   List,
	}
)

func init() {
	ListCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "define the output format: [wide, yaml, json, jsonpath]")
}

func List(cmd *cobra.Command, args []string) {

	// create profile from example
	profile, err := config.GetProfileFromFile("./example/profile.yaml")
	if err != nil {
		log.Fatalf("Could not read example template file: %q", err)
		return
	}

	err = GetAll(profile, outputFlag)
	if err != nil {
		log.Fatalf("Could not get templates: %q", err)
	}
}

func GetAll(profile *config.Profile, outputFormat string) error {
	printer, err := getPrinter(outputFormat)
	if err != nil {
		return err
	}
	return printer.Fprint(os.Stdout, profile.Templates)
}
