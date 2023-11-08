package cmd

import (
	"github.com/simonkienzler/crusado/cmd/template"
	"github.com/simonkienzler/crusado/cmd/version"

	"github.com/spf13/cobra"
)

var (
	crusadoCmd = &cobra.Command{
		Use:   "crusado",
		Short: "crusado allows you to quickly create user stories, bugs and tasks from predefined templates",
		Long: `crusado uses a list of custom, predefined user story and bug templates (including 
their subtasks) to let you quickly create instances of those templates in your
current iteration. You can even use Markdown syntax in your descriptions.`,
		Run: nil,
	}
)

func init() {
	crusadoCmd.AddCommand(version.RootCmd)
	crusadoCmd.AddCommand(template.RootCmd)
}

func Execute() error {
	return crusadoCmd.Execute()
}
