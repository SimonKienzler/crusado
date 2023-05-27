package template

import (
	"log"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/validator"

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
	cfg := config.GetConfig()

	// create templateList from example
	templateList, err := config.GetTemplateListFromFile(cfg.ProfileFilePath)
	if err != nil {
		log.Fatalf("Could not read example template file: %q", err)
	}

	err = validator.ValidateTemplateList(templateList)
	if err != nil {
		log.Fatalf("Invalid profile: %s", err)
	}

	err = GetAll(templateList, outputFlag)
	if err != nil {
		log.Fatalf("Could not get templates: %q", err)
	}
}

func GetAll(profile *config.TemplateList, outputFormat string) error {
	printer, err := getPrinter(outputFormat)
	if err != nil {
		return err
	}
	return printer.Fprint(os.Stdout, profile.Templates)
}
