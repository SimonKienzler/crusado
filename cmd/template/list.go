package template

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "List crusado templates",
		Long:  `Allows you to list available user story and bug templates. You can specify an output format.`,
		Args:  cobra.NoArgs,
		Run:   List,
	}
)

func init() {
	ListCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "define the output format: [wide, yaml, json, jsonpath]")
}

func List(_ *cobra.Command, _ []string) {
	err := GetAll(outputFlag)
	if err != nil {
		log.Fatalf("Could not get templates: %q", err)
	}
}

func GetAll(outputFormat string) error {
	templates, err := templateService().GetAll()
	if err != nil {
		return err
	}

	printer, err := getPrinter(outputFormat)
	if err != nil {
		return err
	}
	return printer.Fprint(os.Stdout, templates)
}
