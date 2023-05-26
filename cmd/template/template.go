package template

import (
	"fmt"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/spf13/cobra"
	"github.com/thediveo/klo"
)

var (
	RootCmd = &cobra.Command{
		Use:     "template",
		Aliases: []string{"t"},
		Short:   "Bundles subcommands that manage crusado templates",
		Long:    `Use the template subcommands to list, show and apply templates.`,
		Args:    cobra.NoArgs,
		Run:     nil,
	}
)

var (
	outputFlag string
)

func init() {
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(ShowCmd)
	RootCmd.AddCommand(ApplyCmd)
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
