package cmd

import (
	"log"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/spf13/cobra"
	"github.com/thediveo/klo"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "TODO",
		Long:  `TODO`,
		Args:  listArgs,
		Run:   List,
	}
)

var (
	outputFlag string
)

func init() {
	listCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "define the output format: [wide, yaml, json, jsonpath]")

	crusadoCmd.AddCommand(listCmd)
}

func listArgs(cmd *cobra.Command, args []string) error {
	// TODO validate args
	return nil
}

func List(cmd *cobra.Command, args []string) {
	// create profile from example
	profile, err := config.GetProfileFromFile("./example/profile.yaml")
	if err != nil {
		log.Fatalf("Could not read example template file: %s", err)
		return
	}

	// Create a table printer with custom columns, to be filled from fields
	// of the objects (namely, Name, Foo, and Bar fields).
	myspecs := klo.Specs{
		DefaultColumnSpec: "NAME:{.Name},DESCRIPTION:{.Description}",
		WideColumnSpec:    "NAME:{.Name},DESCRIPTION:{.Description},STORY TITLE:{.StoryTitle},TASKS:{.Tasks[*].Title}",
	}
	prn, err := klo.PrinterFromFlag(outputFlag, &myspecs)
	if err != nil {
		panic(err)
	}
	// Use a table sorter and tell it to sort by the Name field of our column objects.
	prn.Fprint(os.Stdout, profile.Templates)
}
