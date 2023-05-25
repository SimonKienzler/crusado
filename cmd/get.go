package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/userstorytemplates"
	"github.com/spf13/cobra"
	"github.com/thediveo/klo"
)

var (
	getCmd = &cobra.Command{
		Use:   "get [template-name]",
		Short: "TODO",
		Long:  `TODO`,
		Args:  cobra.MaximumNArgs(1),
		Run:   Get,
	}
)

var (
	outputFlag string
)

func init() {
	getCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "define the output format: [wide, yaml, json, jsonpath]")

	crusadoCmd.AddCommand(getCmd)
}

func Get(cmd *cobra.Command, args []string) {

	// create profile from example
	profile, err := config.GetProfileFromFile("./example/profile.yaml")
	if err != nil {
		log.Fatalf("Could not read example template file: %q", err)
		return
	}

	if len(args) > 0 {
		err = GetByName(profile, args[0], outputFlag)
		if err != nil {
			log.Fatalf("Could not get template by name '%s': %q", args[0], err)
		}
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

func GetByName(profile *config.Profile, name, outputFormat string) error {
	ustService := userstorytemplates.Service{
		Profile: *profile,
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

func getPrinter(outputFormat string) (klo.ValuePrinter, error) {
	return klo.PrinterFromFlag(outputFormat, &config.TemplatePrinterSpecs)
}

func prettyPrintTemplate(template *config.UserStoryTemplate) {
	fmt.Printf("Name:             %s\n", template.Name)
	fmt.Printf("User Story Title: %s\n", template.StoryTitle)
	fmt.Printf("Number of Tasks:  %d\n", len(template.Tasks))
	fmt.Print("Task Overview:\n")
	for _, task := range template.Tasks {
		fmt.Printf("  - %s\n", task.Title)
	}
}
