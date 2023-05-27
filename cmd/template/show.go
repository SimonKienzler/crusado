package template

import (
	"log"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/userstorytemplates"
	"github.com/simonkienzler/crusado/pkg/validator"

	"github.com/spf13/cobra"
)

var (
	ShowCmd = &cobra.Command{
		Use:   "show [template-name]",
		Short: "Show specific crusado template",
		Long:  `TODO`,
		Args:  cobra.ExactArgs(1),
		Run:   Show,
	}
)

func init() {
	ShowCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "define the output format: [wide, yaml, json, jsonpath]")
}

func Show(cmd *cobra.Command, args []string) {
	cfg := config.GetConfig()

	// create templateList from example
	templateList, err := config.GetTemplateListFromFile(cfg.ProfileFilePath)
	if err != nil {
		log.Fatalf("Could not read example template file: %q", err)
		return
	}

	err = validator.ValidateTemplateList(templateList)
	if err != nil {
		log.Fatalf("Invalid profile: %s", err)
	}

	err = GetByName(templateList, args[0], outputFlag)
	if err != nil {
		log.Fatalf("Could not get template by name '%s': %q", args[0], err)
	}
	return
}

func GetByName(profile *config.TemplateList, name, outputFormat string) error {
	ustService := userstorytemplates.Service{
		TemplateList: *profile,
	}

	userStoryTemplate, err := ustService.GetUserStoryTemplateFromName(name)
	if err != nil {
		return err
	}

	// pretty-print single template view in default output format
	if outputFormat == "" {
		prettyPrintTemplate(userStoryTemplate)
		return nil
	}

	printer, err := getPrinter(outputFormat)
	if err != nil {
		return err
	}
	return printer.Fprint(os.Stdout, userStoryTemplate)
}
