package template

import (
	"log"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/templates"
	"github.com/simonkienzler/crusado/pkg/validator"

	"github.com/spf13/cobra"
)

var (
	ShowCmd = &cobra.Command{
		Use:   "show [template-name]",
		Short: "Show a specific crusado template",
		Long: `Get a detailed look at a specific template, listing all relevant fields
and tasks that belong to the template. You can specify an output format.`,
		Args: cobra.ExactArgs(1),
		Run:  Show,
	}
)

func init() {
	ShowCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "define the output format: [wide, yaml, json, jsonpath]")
}

func Show(_ *cobra.Command, args []string) {
	cfg := config.GetConfig()

	// create templateList from example
	templateList, err := config.GetTemplateListFromFile(cfg.ProfileFilePath)
	if err != nil {
		log.Fatalf("Could not read example template file: %q", err)
		return
	}

	err = validator.ValidateTemplateList(templateList)
	if err != nil {
		prettyPrintValidationError(err)
	}

	err = GetByName(templateList, args[0], outputFlag)
	if err != nil {
		log.Fatalf("Could not get template by name '%s': %q", args[0], err)
	}
}

func GetByName(profile *config.TemplateList, name, outputFormat string) error {
	ustService := templates.Service{
		TemplateList: *profile,
	}

	template, err := ustService.GetTemplateFromName(name)
	if err != nil {
		return err
	}

	// pretty-print single template view in default output format
	if outputFormat == "" {
		prettyPrintTemplate(template)
		return nil
	}

	printer, err := getPrinter(outputFormat)
	if err != nil {
		return err
	}
	return printer.Fprint(os.Stdout, template)
}
