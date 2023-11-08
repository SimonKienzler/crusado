package template

import (
	"log"
	"os"

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
	err := GetByName(args[0], outputFlag)
	if err != nil {
		log.Fatalf("Could not get template by name '%s':\n%v", args[0], err)
	}
}

func GetByName(name, outputFormat string) error {
	template, err := crusadoService().GetByName(name)
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
