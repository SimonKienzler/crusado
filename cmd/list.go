package cmd

import (
	"fmt"

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
	// TODO implement list
	fmt.Println("hello from the list command")
}
