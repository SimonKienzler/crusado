package cmd

import (
	"fmt"

	"github.com/simonkienzler/crusado/cmd/template"

	"github.com/spf13/cobra"
)

var (
	crusadoCmd = &cobra.Command{
		Use:   "crusado",
		Short: "crusado allows you to quickly create user stories and tasks from predefined templates",
		Long: `crusado uses a list of custom, predefined user story templates (including 
their subtasks) to let you quickly create instances of those templates in
your current iteration. You can even use Markdown syntax in your descriptions.`,
		Run: Crusado,
	}
)

func init() {
	crusadoCmd.AddCommand(template.RootCmd)
}

func Execute() error {
	return crusadoCmd.Execute()
}

func Crusado(cmd *cobra.Command, args []string) {
	// TODO make this nice
	fmt.Println("crusado root command - run with --help to get an overview!")
}
