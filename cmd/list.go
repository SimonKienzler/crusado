package cmd

import (
	"fmt"
	"log"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/spf13/cobra"
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

func init() {
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
	}

	fmt.Printf("Listing user story templates from profile '%s'\n\n", profile.Name)

	for i := range profile.Templates {
		template := profile.Templates[i]

		fmt.Printf("- %s: %s (%d tasks)\n", template.Name, template.Description, len(template.Tasks))
	}
}
